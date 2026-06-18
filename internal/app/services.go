package app

import "github.com/redish101/depositum/internal/service"

type Services struct {
	library service.LibraryService
}

func (app *app) initServices() {
	libraryService := service.NewLibraryService(app.db)

	app.services = &Services{
		library: libraryService,
	}
}
