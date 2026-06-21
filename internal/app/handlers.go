package app

import (
	"github.com/labstack/echo/v4"
	"github.com/redish101/depositum/internal/handler"
)

func (app *app) initHandlers(v1 *echo.Group) {
	libraryHandler := handler.NewLibraryHandler(app.services.library)
	libraryHandler.Register(v1)

	healthzHandler := handler.NewHealthzHandler()
	healthzHandler.Register(v1)
}
