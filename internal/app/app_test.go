package app

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/redish101/depositum/internal/config"
	v1 "github.com/redish101/depositum/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestApp(t *testing.T) App {
	cfg := &config.Config{
		Debug: true,
		Host:  "127.0.0.1",
		Port:  "0",
		DSN:   ":memory:",
	}

	app, err := New(cfg)
	require.NoError(t, err)

	return app
}

func TestRun(t *testing.T) {
	app := setupTestApp(t)

	errCh := make(chan error, 1)

	go func() {
		errCh <- app.Run()
	}()

	// 等待服务开始监听
	select {
	case <-app.Started():
	case <-time.After(2 * time.Second):
		t.Fatal("server start timeout")
	}

	// 检查健康检查接口
	resp, err := http.Get("http://" + app.Addr() + "/api/v1/healthz")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var apiResp v1.HealthzResponse
	json.Unmarshal(body, &apiResp)
	assert.Equal(t, true, apiResp.Ok)

	// 优雅关闭
	err = app.Shutdown(t.Context())
	require.NoError(t, err)

	select {
	case err := <-errCh:
		assert.ErrorIs(t, err, http.ErrServerClosed)
	case <-time.After(2 * time.Second):
		t.Fatal("server shutdown timeout")
	}
}
