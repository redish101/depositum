package app

import (
	"net/http"
	"testing"
	"time"

	"github.com/redish101/depositum/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestApp(t *testing.T) App {
	cfg := &config.Config{
		Debug: true,
		Host:  "localhost",
		Port:  "0",
		DSN:   ":memory:",
	}

	app, err := New(cfg)
	require.NoError(t, err)
	require.NotNil(t, app)

	return app
}

func TestRun(t *testing.T) {
	app := setupTestApp(t)

	errChan := make(chan error, 1)
	go func() {
		errChan <- app.Run()
	}()

	time.Sleep(100 * time.Millisecond)

	err := app.Shutdown(t.Context())
	assert.NoError(t, err)

	select {
	case err := <-errChan:
		assert.ErrorIs(t, err, http.ErrServerClosed)
	case <-time.After(2 * time.Second):
		assert.Fail(t, "关闭超时")
	}
}
