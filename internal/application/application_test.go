package application

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/redish101/depositum/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestApp(t *testing.T) Application {
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

func TestRunAndShutdown(t *testing.T) {
	app := setupTestApp(t)

	errChan := make(chan error, 1)
	go func() {
		errChan <- app.Run()
	}()

	time.Sleep(100 * time.Millisecond)

	err := app.Shutdown()
	assert.NoError(t, err)

	select {
	case err := <-errChan:
		assert.ErrorIs(t, err, http.ErrServerClosed)
	case <-time.After(2 * time.Second):
		assert.Fail(t, "关闭超时")
	}
}

func TestNewEchoContext(t *testing.T) {
	app := setupTestApp(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	ctx := app.NewEchoContext(req, rec)
	assert.NotNil(t, ctx)

	assert.Equal(t, req, ctx.Request())
	assert.Equal(t, rec, ctx.Response().Writer)
}
