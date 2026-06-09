package service

import (
	"testing"

	"github.com/redish101/depositum/internal/db"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestDB(t *testing.T) *gorm.DB {
	gormDB, err := gorm.Open(sqlite.Open(":memory:"))
	require.NotNil(t, gormDB)
	require.NoError(t, err)

	db.Migrate(gormDB)

	return gormDB
}
