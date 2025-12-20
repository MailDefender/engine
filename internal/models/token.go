package models

import (
	"time"

	"gorm.io/gorm"

	"maildefender/engine/internal/utils"
)

type ValidationToken struct {
	gorm.Model
	Email          string    `gorm:"column:email"`
	Token          string    `gorm:"unique;column:token"`
	AdminToken     bool      `gorm:"column:admin_token"`
	GenerationDate time.Time `gorm:"column:generation_date"`
	ExpiryDate     time.Time `gorm:"column:expiry_date"`
	Validated      bool      `gorm:"column:validated"`
}

func NewValidationToken(email string, adminToken bool) ValidationToken {
	return ValidationToken{
		Email:          email,
		Token:          utils.RandomUuid(),
		AdminToken:     adminToken,
		GenerationDate: time.Now(),
		ExpiryDate:     time.Now().AddDate(0, 0, 7),
		Validated:      false,
	}
}

func SaveValidationToken(tx *gorm.DB, token *ValidationToken) error {
	return tx.Save(token).Error
}

func GetLastValidationTokenByEmail(tx *gorm.DB, email string) (ValidationToken, error) {
	var out ValidationToken
	err := tx.Where("email = ?", email).
		Last(&out).Error
	return out, err
}

func GetLastAdminValidationTokenByEmail(tx *gorm.DB, email string) (ValidationToken, error) {
	var out ValidationToken
	err := tx.Where("email = ?", email).Where("admin_token = ?", true).
		Last(&out).Error
	return out, err
}

func GetValidationTokenByToken(tx *gorm.DB, token string) (ValidationToken, error) {
	var out ValidationToken
	err := tx.Where("token = ?", token).First(&out).Error
	return out, err
}
