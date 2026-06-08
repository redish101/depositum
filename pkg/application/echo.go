package application

import (
	"github.com/labstack/echo/v4"
	"github.com/redish101/depositum/internal/handler"
)

func NewEcho() *echo.Echo {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true
	e.HTTPErrorHandler = handler.HTTPErrorHandler

	return e
}
