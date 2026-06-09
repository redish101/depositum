package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func getEnv[T comparable](key string, defaultValue T) T {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	var result T

	switch any(result).(type) {

	case string:
		result = any(value).(T)

	case bool:
		if value == "true" {
			result = any(true).(T)
		} else {
			result = any(false).(T)
		}

	case int:
		var intValue int
		_, err := fmt.Sscanf(value, "%d", &intValue)
		if err != nil {
			intValue = 0
		}
		result = any(intValue).(T)

	default:
		panic("unimplemented type")
	}

	return result
}

type Config struct {
	Debug bool
	Host  string
	Port  string
	DSN   string
}

func FromEnv() *Config {
	godotenv.Load()

	cfg := &Config{
		Debug: getEnv("DEBUG", false),
		Host:  getEnv("HOST", "localhost"),
		Port:  getEnv("PORT", "3000"),
		DSN:   getEnv("DSN", "depositum.db"),
	}

	return cfg
}
