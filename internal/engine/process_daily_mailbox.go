package engine

import (
	"errors"

	"github.com/sirupsen/logrus"

	"maildefender/engine/internal/constants"
	db "maildefender/engine/internal/database"
	"maildefender/engine/internal/services"
)

func ProcessDailyMailbox() error {
	// Get all email from "INBOX.MOAB._Daily"
	messages, _, err := services.FetchEmails(services.FetchEmailsIn{
		Mailbox: string(constants.DailyMailbox),
	})

	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{"mailbox": string(constants.DailyMailbox)}).Error("cannot fetch emails")
		return err
	}

	tx := db.Instance().Gorm

	// Move them according to rules
	hasError := false
	for _, msg := range messages {
		if moveErr := MoveMessageAccordingRules(tx, msg, string(constants.DailyMailbox)); moveErr != nil {
			logrus.WithFields(logrus.Fields{
				"message_id": msg.MessageID,
				"source":     string(constants.DailyMailbox),
			}).WithError(err).Error("cannot move message according to rules")
			hasError = true
		}
	}

	if hasError {
		return errors.New("at least one error occured")
	}

	return nil
}
