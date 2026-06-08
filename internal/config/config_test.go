package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	originalEnv := map[string]string{
		"STRING_KEY": os.Getenv("STRING_KEY"),
		"BOOL_KEY":   os.Getenv("BOOL_KEY"),
		"INT_KEY":    os.Getenv("INT_KEY"),
		"EMPTY_KEY":  os.Getenv("EMPTY_KEY"),
	}
	defer func() {
		for k, v := range originalEnv {
			os.Setenv(k, v)
		}
	}()

	mockEnv := map[string]string{
		"STRING_KEY": "test",
		"BOOL_KEY":   "true",
		"INT_KEY":    "42",
		"EMPTY_KEY":  "",
	}
	for k, v := range mockEnv {
		os.Setenv(k, v)
	}

	cfg := FromEnv()

	assert.Equal(t, cfg.Host, "localhost")
	assert.Equal(t, cfg.Port, "3000")
	assert.Equal(t, cfg.Debug, false)
	assert.Equal(t, cfg.DSN, "depositum.db")
	assert.Equal(t, getEnv("NON_EXISTENT_KEY", "default"), "default")
}
