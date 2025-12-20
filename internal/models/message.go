package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	MessageID   string    `gorm:"unique;column:message_id"`
	SenderName  string    `gorm:"column:sender_name"`
	SenderEmail string    `gorm:"column:sender_email"`
	Subject     string    `gorm:"column:subject"`
	ReceivedAt  time.Time `gorm:"column:received_at"`
}

func InsertMessage(tx *gorm.DB, entry *Message) error {
	return tx.Create(entry).Error
}
