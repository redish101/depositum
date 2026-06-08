package db

import (
	"path/filepath"
	"testing"

	"github.com/redish101/depositum/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestNewDB(t *testing.T) {
	tmpDir := t.TempDir()

	dbPath := filepath.Join(tmpDir, "depositum_test.db")

	cfg := &config.Config{
		DSN: dbPath,
	}

	db, err := NewDB(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, db)
}
