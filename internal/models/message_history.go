package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type MessageHistoryAction string

const (
	MessageMoveAction MessageHistoryAction = "move"
)

type MessageHistory struct {
	gorm.Model
	ActionType  MessageHistoryAction `gorm:"column:action_type"`
	MessageID   string               `gorm:"column:message_hash"`
	Message     Message              `gorm:"references:MessageID"`
	Source      string               `gorm:"column:source"`
	Destination string               `gorm:"column:destination"`
}

func SaveMessageHistory(tx *gorm.DB, entry *MessageHistory) error {
	return tx.Create(entry).Error
}

func GetMessageHistoryBetweenDates(tx *gorm.DB, from, to time.Time) ([]MessageHistory, error) {
	var out []MessageHistory
	err := tx.Preload("Message").Where("created_at BETWEEN ? AND ?", from, to).Find(&out).Error
	return out, err
}

func (m MessageHistory) IsPending() bool {
	return strings.Contains(m.Destination, "_Pending")
}
