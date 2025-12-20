package models

import "gorm.io/gorm"

type ReputationStatus string

const (
	WhitelistedReputationStatus ReputationStatus = "whitelisted"
	BlacklistedReputationStatus ReputationStatus = "blacklisted"
	PendingReputationStatus     ReputationStatus = "pending"
)

type Reputation struct {
	gorm.Model
	Email  string           `gorm:"column:email"`
	Status ReputationStatus `gorm:"column:status"`
}

type SearchReputationIn struct {
	Email  string
	Status ReputationStatus
}

func SaveReputation(tx *gorm.DB, rep *Reputation) error {
	return tx.Save(rep).Error
}

func GetAllReputations(tx *gorm.DB) ([]Reputation, error) {
	var out []Reputation
	err := tx.Find(&out).Error
	return out, err
}

func GetReputationByEmail(tx *gorm.DB, email string) (Reputation, error) {
	var out Reputation
	err := tx.Where("email = ?", email).Take(&out).Error
	return out, err
}

func SearchReputation(tx *gorm.DB, in SearchReputationIn) ([]Reputation, error) {
	var out []Reputation

	if in.Email != "" {
		tx = tx.Where("email = ?", in.Email)
	}

	if in.Status != "" {
		tx = tx.Where("status = ?", in.Status)
	}

	err := tx.Find(&out).Error
	return out, err
}
