package notification

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"maildefender/engine/internal/configuration"
	"maildefender/engine/internal/constants"
	"maildefender/engine/internal/engine"
	"maildefender/engine/internal/models"
	"maildefender/engine/internal/services"
	"maildefender/engine/internal/templates"
)

type messageWrapper struct {
	Entry models.MessageHistory
	Badge struct {
		Type  string
		Label string
	}
	Validation struct {
		AdminUrl string
	}
}

func (w messageWrapper) badgeType() string {
	switch {
	case w.Entry.Destination == string(constants.PendingMailbox):
		return "pending"
	case w.Entry.Source == string(constants.PendingMailbox):
		return "resolved"
	case w.Entry.Destination == string(constants.BlackListedMailbox):
		return "blacklisted"
	default:
		return "new"
	}
}

func (w messageWrapper) badgeLabel() string {
	switch {
	case w.Entry.Destination == string(constants.PendingMailbox):
		return "Pending"
	case w.Entry.Source == string(constants.PendingMailbox):
		return "Resolved"
	case w.Entry.Destination == string(constants.BlackListedMailbox):
		return "Blacklisted"
	default:
		return "New"
	}
}

func (w messageWrapper) Destination() string {
	return strings.TrimPrefix(w.Entry.Destination, string(constants.MailDefenderMailbox)+".")
}

func (w messageWrapper) Source() string {
	return strings.TrimPrefix(w.Entry.Source, string(constants.MailDefenderMailbox)+".")
}

func SendMailRecap(tx *gorm.DB, from, to time.Time) error {
	operations, err := models.GetMessageHistoryBetweenDates(tx, from, to)
	if err != nil {
		return err
	}

	// Get last operation for each email
	uniq := map[string]models.MessageHistory{}
	for _, op := range operations {
		uniq[op.MessageID] = op
	}

	wrappedOperations := make([]messageWrapper, 0)

	for _, val := range uniq {
		w := messageWrapper{Entry: val}
		w.Badge.Label = w.badgeLabel()
		w.Badge.Type = w.badgeType()

		if w.badgeType() == "pending" {
			token, err := models.GetLastAdminValidationTokenByEmail(tx, val.Message.SenderEmail)
			if errors.Is(err, gorm.ErrRecordNotFound) || token.ExpiryDate.Before(time.Now()) {
				token = models.NewValidationToken(val.Message.SenderEmail, true)
				if err := models.SaveValidationToken(tx, &token); err != nil {
					logrus.WithError(err).Fatal("cannot create token for this sender")
				}
				w.Validation.AdminUrl = engine.GenerateValidationUri(token.Token)
			} else if err != nil {
				logrus.WithFields(logrus.Fields{"sender": val.Message.SenderEmail}).Error("cannot retrieve admin token for this sender")
			} else {
				w.Validation.AdminUrl = engine.GenerateValidationUri(token.Token)
			}
		}

		wrappedOperations = append(wrappedOperations, w)
	}

	today := time.Now()

	recapIn := struct {
		TodayDate  time.Time
		Operations []messageWrapper
	}{
		TodayDate:  today,
		Operations: wrappedOperations,
	}

	var recapOut bytes.Buffer
	if err := templates.OperationRecapTemplate().Execute(&recapOut, recapIn); err != nil {
		fmt.Print(err)
		return err
	}

	_, err = services.SendEmail(services.SendMailIn{
		To:      []string{configuration.DailyRecapRecipient()},
		Subject: fmt.Sprintf("MailDefender: Daily recap - %s", time.Now().Format(time.DateOnly)),
		Body:    recapOut.String(),
	})

	if err == nil {
		models.SaveNotification(tx, &models.Notification{
			Type:      models.DailyRecap,
			Channel:   models.EmailChannel,
			Recipient: configuration.DailyRecapRecipient(),
			Content:   recapOut.String(),
		})
	}

	return err
}
