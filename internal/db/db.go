package db

import (
	"github.com/redish101/depositum/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDB(cfg *config.Config) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(cfg.DSN))
}
