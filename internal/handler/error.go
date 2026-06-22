package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/redish101/depositum/internal/common"
)

type ErrorHandler interface {
	Register(e *echo.Echo)
	All(c *echo.Context, err error)
}

type errorHandler struct{}

func NewErrorHandler() ErrorHandler {
	return &errorHandler{}
}

func (e *errorHandler) Register(echo *echo.Echo) {
	echo.HTTPErrorHandler = e.All
}

func (h *errorHandler) All(c *echo.Context, err error) {
	if httpErr, ok := errors.AsType[interface {
		Error() string
		StatusCode() int
	}](err); ok {
		common.WriteError(c, httpErr.StatusCode(), fmt.Errorf("%s", httpErr.Error()))
		return
	}

	common.WriteError(c, http.StatusInternalServerError, err)
}
