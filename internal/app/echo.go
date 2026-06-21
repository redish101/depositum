package app

import (
	"github.com/labstack/echo/v5"
	"github.com/redish101/depositum/internal/handler"
	"github.com/redish101/depositum/internal/validate"
)

func NewEcho() (*echo.Echo, error) {
	e := echo.New()

	validator, err := validate.NewValidator()
	if err != nil {
		return nil, err
	}
	e.Validator = validator

	errorHandler := handler.NewErrorHandler()
	errorHandler.Register(e)

	return e, nil
}
