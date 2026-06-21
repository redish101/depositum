package common

import (
	"net/http"

	"github.com/labstack/echo/v5"
	v1 "github.com/redish101/depositum/pkg/api/v1"
)

func WriteError(c *echo.Context, statusCode int, err error) error {
	return WriteEntityWithStatus(c, statusCode, &v1.ErrorResponse{
		Error: err.Error(),
	})
}

func WriteEntity(c *echo.Context, data any) error {
	return WriteEntityWithStatus(c, http.StatusOK, data)
}

func WriteEntityWithStatus(c *echo.Context, statusCode int, data any) error {
	return WriteEntityAndHeader(c, statusCode, map[string]string{}, data)
}

func WriteEntityAndHeader(c *echo.Context, statusCode int, headers map[string]string, data any) error {
	for key, value := range headers {
		c.Response().Header().Set(key, value)
	}
	return c.JSON(statusCode, data)
}
