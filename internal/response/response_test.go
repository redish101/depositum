package response

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/this-is-a-fake-path", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	data := map[string]string{"message": "Hello, World!"}

	err := Write(c, data)
	assert.NoError(t, err)
}

func TestWriteWithStatus(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/this-is-a-fake-path", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	data := map[string]string{"message": "Hello, World!"}

	err := WriteWithStatus(c, http.StatusCreated, data)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestWriteError(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/this-is-a-fake-path", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	err := WriteError(c, http.StatusBadRequest, assert.AnError)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
