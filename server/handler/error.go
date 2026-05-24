package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/redish101/depositum/server/response"
)

func HTTPErrorHandler(err error, c echo.Context) {
	if he, ok := err.(*echo.HTTPError); ok {
		code := he.Code

		response.WriteError(c, code, errors.New(he.Message.(string)))

		return
	}
	response.WriteError(c, http.StatusInternalServerError, err)
}
