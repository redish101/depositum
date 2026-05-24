package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func WriteWithStatus(c echo.Context, status int, data any) error {
	return c.JSON(status, data)
}

func Write(c echo.Context, data any) error {
	return WriteWithStatus(c, http.StatusOK, data)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteError(c echo.Context, status int, err error) error {
	resp := ErrorResponse{Error: err.Error()}

	return c.JSON(status, resp)
}
