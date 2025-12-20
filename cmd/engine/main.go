package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "ariga.io/atlas-provider-gorm/gormschema"
	"github.com/sirupsen/logrus"

	"maildefender/engine/internal/api"
	"maildefender/engine/internal/configuration"
	"maildefender/engine/internal/constants"
	db "maildefender/engine/internal/database"
	"maildefender/engine/internal/engine"
	"maildefender/engine/internal/models"
	"maildefender/engine/internal/notification"
	"maildefender/engine/internal/rules"
	"maildefender/engine/internal/services"
)

const hourRecap = 20

func createMailboxes(mailboxes []constants.MailboxName) bool {
	for _, mbx := range mailboxes {
		if !createMailbox(mbx) {
			return false
		}
	}
	return true
}

func createMailbox(mailbox constants.MailboxName) bool {

	_, httpCode, err := services.CreateMailbox(string(mailbox))
	logger := logrus.WithFields(logrus.Fields{"mailbox": mailbox, "http_code": httpCode})

	if httpCode == http.StatusCreated {
		logger.Info("mailbox created")
		return true
	}

	if httpCode == http.StatusConflict {
		logger.Info("mailbox already exist")
		return true
	}

	logger.WithError(err).Error("cannot create mailbox")
	return false
}

func main() {
	// logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	logrus.WithFields(logrus.Fields{
		"imap_connector":       configuration.ImapConnectorBaseEndpoint(),
		"notifier":             configuration.NotifierBaseEndpoint(),
		"validator":            configuration.ValidatorPublicBaseEndpoint(),
		"loop_delay":           configuration.LoopDelay(),
		"skip_reputation":      configuration.SkipReputationCheck(),
		"enable_daily_mailbox": configuration.DailyMailboxEnabled(),
		"enable_daily_recap":   configuration.EnableDailyRecap(),
		"rules_directory":      configuration.RulesDirectory(),
	}).Info("Starting engine...")

	_, err := db.Connect(configuration.DatabaseDNS())
	if err != nil {
		logrus.WithError(err).Error("cannot connect to database")
		os.Exit(1)
	}

	if dir := configuration.RulesDirectory(); dir != "" {
		var rulesDirectory string
		if strings.HasPrefix(dir, "/") {
			rulesDirectory = dir
		} else {
			currentDir, err := filepath.Abs("./")
			if err != nil {
				logrus.WithError(err).Error("cannot get current directory")
			} else {
				rulesDirectory = fmt.Sprintf("%s/%s", currentDir, dir)
			}
		}
		if err == nil {
			rulesToIgnore, err := models.GetAllRulesUID(db.Instance().Gorm)
			if err != nil {
				logrus.WithError(err).Error("cannot retrieve rules UUID, no rules to ignore")
				rulesToIgnore = []string{}
			}
			if err := rules.PopulateFromDir(db.Instance().Gorm, rulesDirectory, rulesToIgnore); err != nil {
				logrus.WithError(err).Error("cannot initialize rules")
			}

		}
	}

	// Create default emails folder
	{
		logrus.Info("creating default mailboxes")
		mailboxToCreate := []constants.MailboxName{
			constants.MailDefenderMailbox,
			constants.PendingMailbox,
			constants.BlackListedMailbox,
			constants.UncategorizedMailbox,
		}
		if configuration.DailyMailboxEnabled() {
			mailboxToCreate = append(mailboxToCreate, constants.DailyMailbox)
		}

		mboxesCreate := createMailboxes(mailboxToCreate)
		if !mboxesCreate {
			logrus.Error("default mailboxes not created, quiting...")
			return
		}
		logrus.Info("default mailboxes sucessfully created")
	}

	go func() {
		api.RegisterHandlers()
		api.Run()
	}()

	now := time.Now()

	var lastRecapDate time.Time

	last, err := models.GetLastNotificationByType(db.Instance().Gorm, models.DailyRecap)
	if err == nil {
		lastRecapDate = last.CreatedAt
	} else {
		lastRecapDate = time.Date(now.Year(), now.Month(), now.Day()-1, hourRecap, 0, 0, 0, time.Local)
	}

	nextRecapDate := time.Date(now.Year(), now.Month(), now.Day(), hourRecap, 0, 0, 0, time.Local)

	if now.Hour() >= hourRecap && now.Minute() > 0 {
		nextRecapDate = nextRecapDate.AddDate(0, 0, 1)
	}

	for {
		if err := engine.Process(engine.EngineProcessIn{
			Mailbox: string(constants.DefaultMailbox),
		}); err != nil {
			logrus.WithError(err).Warning("process ended with error")
		}

		if time.Now().After(nextRecapDate) {
			if configuration.DailyMailboxEnabled() {
				if err := engine.ProcessDailyMailbox(); err != nil {
					logrus.WithError(err).Error("cannot process daily mailbox")
				} else {
					logrus.Info("daily mailbox processed")
				}
			}

			if configuration.EnableDailyRecap() {
				err := notification.SendMailRecap(db.Instance().Gorm, lastRecapDate, nextRecapDate)
				logrus.WithError(err).Info("Daily recap sent")
				lastRecapDate, nextRecapDate = nextRecapDate, nextRecapDate.AddDate(0, 0, 1)
			}
		}

		time.Sleep(time.Second * time.Duration(configuration.LoopDelay()))
	}
}
