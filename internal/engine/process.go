package engine

import (
	"github.com/sirupsen/logrus"

	"maildefender/engine/internal/services"
)

type EngineProcessIn struct {
	Mailbox        string
	SkipReputation bool
}

func Process(in EngineProcessIn) error {
	messages, httpCode, err := services.FetchEmails(services.FetchEmailsIn{
		Mailbox: in.Mailbox,
	})
	if err != nil || httpCode != 200 {
		logrus.WithField("http_status", httpCode).WithError(err).Error("failed to fetch emails")
		return err
	}

	logger := logrus.WithField("message_count", len(messages))
	if len(messages) == 0 {
		logger.Info("no message to process")
		return nil
	}

	if err := refreshCachedRules(); err != nil {
		logrus.WithError(err).Error("cannot refresh cached rules")
	}

	logger.Info("processing messages")
	for _, msg := range messages {
		processMessage(msg)
	}
	return nil
}
