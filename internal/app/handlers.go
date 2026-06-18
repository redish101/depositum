package app

import "github.com/redish101/depositum/internal/handler"

func (app *app) initHandlers() {
	libraryHandler := handler.NewLibraryHandler(app.services.library)
	libraryHandler.Register(app.container)

	healthzHandler := handler.NewHealthzHandler()
	healthzHandler.Register(app.container)
}
