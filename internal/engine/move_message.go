package engine

import (
	// TODO: Change this dependency once initialized with git
	stderr "errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"maildefender/engine/internal/constants"
	"maildefender/engine/internal/errors"
	"maildefender/engine/internal/models"
	"maildefender/engine/internal/rules"
	"maildefender/engine/internal/services"
	third_party "maildefender/engine/internal/third_party/models"
)

func MoveMessageAccordingRules(tx *gorm.DB, msg third_party.Message, sourceMailbox string) error {
	// Refresh cached rules if needed
	if err := refreshCachedRules(); err != nil {
		logrus.WithError(err).Error("cannot refresh cached rules")
		if len(cachedRules.rules) == 0 {
			logrus.Error("no rules in cache, aborting...")
			return stderr.New("no rules in cache")
		}
	}

	llogger := logrus.WithField("message_id", msg.MessageID)

	rules, err := rules.FindMatch(cachedRules.rules, msg)
	if err != nil {
		llogger.WithError(err).Error("cannot find matching rules")
		return errors.ErrCannotFindMatchingRules
	}

	if len(rules) > 1 {
		llogger.WithField("rules_count", len(rules)).Warn("too many matching rules, moving to ", string(constants.UncategorizedMailbox))

		if err := MoveMessage(tx, msg, sourceMailbox, string(constants.UncategorizedMailbox)); err != nil {
			return err
		}

		return nil
	}

	if len(rules) == 0 {
		llogger.Error("no matching rules found, moving to ", string(constants.UncategorizedMailbox))

		if err := MoveMessage(tx, msg, sourceMailbox, string(constants.UncategorizedMailbox)); err != nil {
			return err
		}

		return nil
	}

	destination := fmt.Sprintf("%s.%s", string(constants.MailDefenderMailbox), rules[0].Destination)
	if err := MoveMessage(tx, msg, sourceMailbox, destination); err != nil {
		return err
	}

	return nil
}

func MoveMessage(tx *gorm.DB, msg third_party.Message, source, destination string) error {
	moveOut, httpCode, err := services.MoveEmail(services.MoveEmailIn{
		MessageID:     msg.MessageID,
		Source:        source,
		Destination:   destination,
		CreateMailbox: true,
	})

	if err != nil {
		return err
	}

	if httpCode != 200 {
		return stderr.New(moveOut.Error)
	}

	err = models.SaveMessageHistory(tx, &models.MessageHistory{
		ActionType:  models.MessageMoveAction,
		Message:     models.Message{MessageID: msg.MessageID},
		Source:      source,
		Destination: destination,
	})

	if err != nil {
		return errors.ErrCannotSaveHistory
	}

	return nil
}
