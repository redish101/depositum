package app

import "github.com/redish101/depositum/internal/service"

type Services struct {
	LibraryService service.LibraryService
}

func (app *app) initServices() {
	services := &Services{}

	services.LibraryService = service.NewLibraryService(app.db)

	app.services = services
}
