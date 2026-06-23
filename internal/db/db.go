package db

import (
	"path/filepath"
	"testing"

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
	db, err := gorm.Open(sqlite.Open(":memory:"))
	if err != nil {
		return nil, err
	}

	err = Migrate(db)

	return db, err
}

func NewTempDB(tb testing.TB) (*gorm.DB, error) {
	tb.Helper()

	dir := tb.TempDir()
	dsn := filepath.Join(dir, "benchmark.db")

	db, err := gorm.Open(sqlite.Open(dsn))
	if err != nil {
		return nil, err
	}

	if err := Migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Object{},
		&model.Library{},
		&model.Shelf{},
	)
}
