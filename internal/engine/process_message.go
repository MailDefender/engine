package engine

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"maildefender/engine/internal/configuration"
	"maildefender/engine/internal/constants"
	db "maildefender/engine/internal/database"
	"maildefender/engine/internal/models"
	"maildefender/engine/internal/services"
	"maildefender/engine/internal/templates"
	thirdModels "maildefender/engine/internal/third_party/models"
	"maildefender/engine/internal/utils"
)

type emailValidationParameters struct {
	ValidationUri string
	ExpiryDate    time.Time
}

func processMessage(inMsg thirdModels.Message) {
	senderEmail := inMsg.From[0].Email

	logger := logrus.WithFields(logrus.Fields{
		"message_id": inMsg.MessageID,
		"sender":     senderEmail,
		"subject":    inMsg.Subject,
	})

	tx := db.Instance().Gorm

	msg := models.Message{
		MessageID:   inMsg.MessageID,
		SenderName:  inMsg.From[0].Name,
		SenderEmail: inMsg.From[0].Email,
		Subject:     inMsg.Subject,
		ReceivedAt:  inMsg.Date,
	}

	if err := models.InsertMessage(tx, &msg); err != nil {
		if utils.IsUniqueViolationErr(err) {
			logger.WithError(err).Warn("it seems that this message is already in the database")
		} else {
			logger.WithError(err).Error("cannot save this message")
		}
	}

	if configuration.SkipReputationCheck() {
		MoveMessageAccordingRules(tx, inMsg, string(constants.DefaultMailbox))
		return
	}

	rep, err := models.GetReputationByEmail(tx, senderEmail)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Info("new sender, validation required")

		// Create all object
		token := models.NewValidationToken(senderEmail, false)
		senderReputation := models.Reputation{
			Email:  senderEmail,
			Status: models.PendingReputationStatus,
		}
		pendingEmail := models.PendingMessage{
			Message:     msg,
			SenderEmail: senderEmail,
			Subject:     msg.Subject,
			Mailbox:     string(constants.PendingMailbox),
		}

		if err := tx.Transaction(func(tx *gorm.DB) error {
			if err := models.SaveReputation(tx, &senderReputation); err != nil {
				return err
			}

			if err := models.SaveValidationToken(tx, &token); err != nil {
				return err
			}

			if err := models.InsertPendingMessage(tx, &pendingEmail); err != nil {
				return err
			}

			return nil
		}); err != nil {
			logger.WithError(err).Fatal("cannot apply transaction while creating token and reputation")
			return
		}

		logger := logrus.WithFields(logrus.Fields{
			"message_id":  msg.MessageID,
			"source":      string(constants.DefaultMailbox),
			"destination": string(constants.PendingMailbox),
			"action":      models.MessageMoveAction,
		})
		if err := MoveMessage(tx, inMsg, string(constants.DefaultMailbox), string(constants.PendingMailbox)); err != nil {
			logger.WithError(err).Error("cannot move message")
		} else {
			logger.Info("message moved")
		}

		templateParams := emailValidationParameters{
			ValidationUri: GenerateValidationUri(token.Token),
			ExpiryDate:    token.ExpiryDate,
		}

		var mailOutContent bytes.Buffer
		if err := templates.FirstEmailValidationTemplate().Execute(&mailOutContent, templateParams); err != nil {
			logrus.WithError(err).WithField("template", "first_email_validation").Fatal("cannot execute template")
			return
		}

		statusCode, err := services.SendEmail(services.SendMailIn{
			ReplyTo:     inMsg.MessageID,
			ThreadTopic: inMsg.Subject,
			To:          []string{senderEmail},
			Subject:     fmt.Sprintf("Re: %s", inMsg.Subject),
			Body:        mailOutContent.String(),
		})

		logger.WithField("http_status", statusCode).WithError(err).Info("Mail validation sent")
		return
	}

	logger.WithField("reputation", rep.Status).Info("reputation found")

	switch rep.Status {
	// Reputation == "blacklist" => Move to "blacklist" mailbox
	case models.BlacklistedReputationStatus:
		MoveMessage(tx, inMsg, string(constants.DefaultMailbox), string(constants.BlackListedMailbox))
		break

	// Reputation == "whitelist" => apply matching rule
	case models.WhitelistedReputationStatus:
		// If the 'daily mailbox' is enabled, the message must be moved to this mailbox
		// The message will then be sorted just before the daily recap

		if configuration.DailyMailboxEnabled() {
			MoveMessage(tx, inMsg, string(constants.DefaultMailbox), string(constants.DailyMailbox))
			logger.WithField("mailbox", string(constants.DailyMailbox)).Info("message moved to the daily mailbox")
			break
		}

		// Get rule and move email
		MoveMessageAccordingRules(tx, inMsg, string(constants.DefaultMailbox))
		break

	// Reputation == "Pending" => send 2nd validation email if needed and move to "Pending" mailbox
	case models.PendingReputationStatus:
		pendingEmailCount, err := models.CountPendingMessageForSenderEmail(tx, senderEmail)
		if err != nil {
			logger.WithError(err).Error("cannot count pending emails for this sender")
			return
		}

		// Retrive last token
		token, err := models.GetLastValidationTokenByEmail(tx, senderEmail)

		// We should create a new token
		var emailTemplate *template.Template
		if errors.Is(err, gorm.ErrRecordNotFound) || token.ExpiryDate.Before(time.Now()) {
			token = models.NewValidationToken(senderEmail, false)
			if err := models.SaveValidationToken(tx, &token); err != nil {
				logger.WithError(err).Fatal("cannot create token for this sender")
				return
			}

			emailTemplate = templates.FirstEmailValidationTemplate()
		} else if pendingEmailCount == 1 {
			emailTemplate = templates.SecondEmailValidationTemplate()
		}

		if err := models.InsertPendingMessage(tx, &models.PendingMessage{
			Message:     msg,
			SenderEmail: senderEmail,
			Subject:     msg.Subject,
			Mailbox:     string(constants.PendingMailbox),
		}); err != nil {
			logger.WithError(err).Error("cannot insert pending message into database")
			return
		}

		logger := logrus.WithFields(logrus.Fields{
			"message_id":  msg.MessageID,
			"source":      string(constants.DefaultMailbox),
			"destination": string(constants.PendingMailbox),
			"action":      models.MessageMoveAction,
		})
		if err := MoveMessage(tx, inMsg, string(constants.DefaultMailbox), string(constants.PendingMailbox)); err != nil {
			logger.WithError(err).Error("cannot move message")
		} else {
			logger.Info("message moved")
		}

		if emailTemplate != nil {
			var mailOutContent bytes.Buffer
			if err := emailTemplate.Execute(&mailOutContent, emailValidationParameters{
				ValidationUri: GenerateValidationUri(token.Token),
				ExpiryDate:    token.ExpiryDate,
			}); err != nil {
				logrus.WithError(err).WithField("template", "first_email_validation").Fatal("cannot execute template")
				return
			}
			httpCode, err := services.SendEmail(services.SendMailIn{
				To:      []string{senderEmail},
				Subject: fmt.Sprintf("Re: %s", msg.Subject),
				Body:    mailOutContent.String(),
			})
			logger.WithField("http_status", httpCode).WithError(err).Info("confirmation email sent")
		}
	}
}
