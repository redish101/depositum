package test

import (
	"testing"

	"github.com/redish101/depositum/internal/config"
	"github.com/redish101/depositum/pkg/application"
	"github.com/stretchr/testify/require"
)

func NewTestApplication(t *testing.T) application.Application {
	cfg := &config.Config{
		Debug: true,
		Host:  "localhost",
		Port:  "3000",
		DSN:   ":memory:",
	}

	app, err := application.New(cfg)
	require.NotNil(t, app)
	require.NoError(t, err)

	return app
}
