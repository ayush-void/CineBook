package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPostgres opens a GORM connection using the given DSN.
func NewPostgres(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
