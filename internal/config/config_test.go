package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	originalEnv := map[string]string{
		"STRING_KEY":       os.Getenv("STRING_KEY"),
		"BOOL_KEY":         os.Getenv("BOOL_KEY"),
		"INT_KEY":          os.Getenv("INT_KEY"),
		"INVALID_INT_KEY":  os.Getenv("INVALID_INT_KEY"),
		"INVALID_BOOL_KEY": os.Getenv("INVALID_BOOL_KEY"),
		"EMPTY_KEY":        os.Getenv("EMPTY_KEY"),
	}
	defer func() {
		for k, v := range originalEnv {
			os.Setenv(k, v)
		}
	}()

	mockEnv := map[string]string{
		"STRING_KEY":       "test",
		"BOOL_KEY":         "true",
		"INVALID_BOOL_KEY": "not_a_bool",
		"INT_KEY":          "42",
		"INVALID_INT_KEY":  "not_an_int",
		"EMPTY_KEY":        "",
	}
	for k, v := range mockEnv {
		os.Setenv(k, v)
	}

	assert.Equal(t, "test", getEnv("STRING_KEY", "default"))
	assert.Equal(t, true, getEnv("BOOL_KEY", false))
	assert.Equal(t, false, getEnv("INVALID_BOOL_KEY", false))
	assert.Equal(t, 42, getEnv("INT_KEY", 0))
	assert.Equal(t, 0, getEnv("INVALID_INT_KEY", 0))
	assert.Equal(t, "default", getEnv("NON_EXISTENT_KEY", "default"))
	assert.Equal(t, "default", getEnv("EMPTY_KEY", "default"))
}
