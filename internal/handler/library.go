package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/redish101/depositum/internal/response"
	"github.com/redish101/depositum/internal/service"
)

type LibraryHandler interface {
	Get(c echo.Context) error
}

type libraryHandler struct {
	libraryService service.LibraryService
}

func NewLibraryHandler(libraryService service.LibraryService) LibraryHandler {
	return &libraryHandler{
		libraryService,
	}
}

func (l *libraryHandler) Get(c echo.Context) error {
	library, _ := l.libraryService.Get(c.Request().Context(), 1)

	return response.Write(c, library)
}
