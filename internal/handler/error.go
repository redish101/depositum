package handler

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/redish101/depositum/internal/common"
)

type ErrorHandler interface {
	Register(e *echo.Echo)
	All(err error, c echo.Context)
}

type errorHandler struct{}

func NewErrorHandler() ErrorHandler {
	return &errorHandler{}
}

func (e *errorHandler) Register(echo *echo.Echo) {
	echo.HTTPErrorHandler = e.All
}

func (e *errorHandler) All(err error, c echo.Context) {
	httpError := err.(*echo.HTTPError)
	common.WriteError(c, httpError.Code, fmt.Errorf("%v", httpError.Message))
}
