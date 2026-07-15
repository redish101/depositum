package app

import (
	"github.com/labstack/echo/v5"
	"github.com/redish101/depositum/internal/handler"
)

func (app *app) initHandlers(v1 *echo.Group) {
	libraryHandler := handler.NewLibraryHandler(app.services.Library)
	libraryHandler.Register(v1)

	shelfHandler := handler.NewShelfHandler(app.services.Shelf)
	shelfHandler.Register(v1)

	objectHandler := handler.NewObjectHandler(app.services.Object)
	objectHandler.Register(v1)

	healthzHandler := handler.NewHealthzHandler()
	healthzHandler.Register(v1)
}
