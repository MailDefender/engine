package models

import "gorm.io/gorm"

type PendingMessage struct {
	gorm.Model
	MessageID   string  `gorm:"column:message_id"`
	Message     Message `gorm:"references:MessageID"`
	SenderEmail string  `gorm:"column:sender_email"`
	Subject     string  `gorm:"column:subject"`
	Mailbox     string  `gorm:"column:mailbox"`
}

func InsertPendingMessage(tx *gorm.DB, pe *PendingMessage) error {
	return tx.Create(pe).Error
}

func GetPendingMessageBySenderEmail(tx *gorm.DB, email string) ([]PendingMessage, error) {
	var out []PendingMessage
	err := tx.Where("sender_email = ?", email).Find(&out).Error
	return out, err
}

func CountPendingMessageForSenderEmail(tx *gorm.DB, email string) (int64, error) {
	var out int64
	err := tx.Model(&PendingMessage{}).Where("sender_email = ?", email).Count(&out).Error
	return out, err
}

func DeletePendingMessageFromID(tx *gorm.DB, messageID string) error {
	return tx.Where("message_id LIKE ?", messageID).Delete(&PendingMessage{}).Error
}
