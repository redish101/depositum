package app

import (
	"github.com/labstack/echo/v4"
	"github.com/redish101/depositum/internal/handler"
)

func NewEcho() *echo.Echo {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	errorHandler := handler.NewErrorHandler()
	errorHandler.Register(e)

	return e
}
