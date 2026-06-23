package app

import "github.com/redish101/depositum/internal/service"

type Services struct {
	library service.LibraryService
	shelf   service.ShelfService
}

func (app *app) initServices() {
	libraryService := service.NewLibraryService(app.db)
	shelfService := service.NewShelfService(app.db, libraryService)

	app.services = &Services{
		library: libraryService,
		shelf:   shelfService,
	}
}
