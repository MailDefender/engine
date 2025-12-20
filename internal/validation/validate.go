package validation

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"maildefender/engine/internal/configuration"
	"maildefender/engine/internal/constants"
	"maildefender/engine/internal/engine"
	engineErrors "maildefender/engine/internal/errors"
	"maildefender/engine/internal/models"
	"maildefender/engine/internal/services"
)

func Validate(tx *gorm.DB, token models.ValidationToken) error {
	if token.Validated {
		return engineErrors.ErrAlreadyValidatedToken
	}

	if token.ExpiryDate.Before(time.Now()) {
		return engineErrors.ErrExpiredToken
	}

	// Get user reputation
	rep, err := models.GetReputationByEmail(tx, token.Email)
	if err != nil {
		return err
	}

	rep.Status = models.WhitelistedReputationStatus
	token.Validated = true

	if err := tx.Transaction(func(tx *gorm.DB) error {
		if err := models.SaveValidationToken(tx, &token); err != nil {
			return nil
		}

		if err := models.SaveReputation(tx, &rep); err != nil {
			return nil
		}

		return nil
	}); err != nil {
		return err
	}

	// Even if a task failed, we return a success.
	// A periodic job will free all pending mails with a whitelisted sender.

	// Get pending emails
	pendingEmails, _, err := services.FetchEmails(services.FetchEmailsIn{
		Mailbox: string(constants.PendingMailbox),
		Sender:  &rep.Email,
	})
	if err != nil {
		log.WithError(err).Warn("cannot retrieve pending emails while validating a token")
		return nil
	}

	for _, pe := range pendingEmails {
		var err error
		if configuration.DailyMailboxEnabled() {
			err = engine.MoveMessage(tx, pe, string(constants.PendingMailbox), string(constants.DailyMailbox))
		} else {
			err = engine.MoveMessageAccordingRules(tx, pe, string(constants.PendingMailbox))
		}

		if err != nil &&
			!errors.Is(err, engineErrors.ErrCannotFindMatchingRules) &&
			!errors.Is(err, engineErrors.ErrTooManyRulesFound) &&
			!errors.Is(err, engineErrors.ErrCannotSaveHistory) &&
			!errors.Is(err, engineErrors.ErrNoRuleFound) {
			return err
		}

		if err := models.DeletePendingMessageFromID(tx, pe.MessageID); err != nil {
			log.WithField("message_id", pe.MessageID).WithError(err).Warn("cannot delete pending message from database")
			return err
		}
	}

	return nil
}
