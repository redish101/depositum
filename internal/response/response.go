package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
	v1 "github.com/redish101/depositum/pkg/api/v1"
)

func WriteWithStatus(c echo.Context, status int, data any) error {
	return c.JSON(status, data)
}

func Write(c echo.Context, data any) error {
	return WriteWithStatus(c, http.StatusOK, data)
}

func WriteError(c echo.Context, status int, err error) error {
	resp := v1.ErrorResponse{Error: err.Error()}

	return c.JSON(status, resp)
}
