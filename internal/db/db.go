package db

import (
	"github.com/redish101/depositum/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDB(cfg *config.Config) (*gorm.DB, error) {
	if cfg.Debug {
		return gorm.Open(sqlite.Open("depositum_debug.db"))
	} else {
		return gorm.Open(sqlite.Open("depositum.db"))
	}
}
