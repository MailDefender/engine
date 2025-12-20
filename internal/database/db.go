package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type database struct {
	dsn string

	Gorm     *gorm.DB
	migrator gorm.Migrator
}

var db *database = &database{}

func Instance() *database {
	return db
}

func Connect(dsn string) (*database, error) {
	var err error
	db.dsn = dsn
	db.Gorm, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db, err
}
