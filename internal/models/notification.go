package models

import "gorm.io/gorm"

type NotificationType string

const (
	DailyRecap NotificationType = "daily_recap"
)

type NotificationChannel string

const (
	EmailChannel NotificationChannel = "email"
)

type Notification struct {
	gorm.Model
	Type      NotificationType    `gorm:"column:type"`
	Channel   NotificationChannel `gorm:"column:channel"`
	Recipient string              `gorm:"column:recipient"`
	Content   string              `gorm:"column:content"`
}

func SaveNotification(tx *gorm.DB, notif *Notification) error {
	return tx.Save(notif).Error
}

func GetLastNotificationByType(tx *gorm.DB, notifType NotificationType) (Notification, error) {
	var out Notification
	err := tx.Where("type = ?", notifType).Last(&out).Error
	return out, err
}
