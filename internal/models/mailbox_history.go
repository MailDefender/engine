package models

import (
	"gorm.io/gorm"
)

type MailboxHistoryAction string

const (
	MailBoxCreateAction MailboxHistoryAction = "create"
)

type MailboxHistory struct {
	gorm.Model
	ActionType MailboxHistoryAction `gorm:"column:action_type"`
	Name       string               `gorm:"column:name"`
}

func SaveMailboxHistory(tx gorm.DB, entry *MailboxHistory) error {
	return tx.Create(entry).Error
}
