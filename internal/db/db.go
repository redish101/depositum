package db

import (
	"github.com/redish101/depositum/internal/config"
	"github.com/redish101/depositum/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DSN))
	if err != nil {
		return nil, err
	}

	err = Migrate(db)

	return db, err
}

func NewMemoryDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"))
	if err != nil {
		return nil, err
	}

	err = Migrate(db)

	return db, err
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Object{},
		&model.Library{},
		&model.Shelf{},
	)
}
