package app

import "github.com/redish101/depositum/internal/handler"

type Handlers struct {
	LibraryHandler handler.LibraryHandler
}

func (app *app) initHandlers() {
	handlers := &Handlers{}

	handlers.LibraryHandler = handler.NewLibraryHandler(app.services.LibraryService)
	handlers.LibraryHandler.Register(app.container)

	app.handlers = handlers
}
