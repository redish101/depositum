package app

import "github.com/redish101/depositum/internal/service"

type Services struct {
	Library service.LibraryService
	Shelf   service.ShelfService
	Object  service.ObjectService
}

func (app *app) initServices() {
	libraryService := service.NewLibraryService(app.db)
	shelfService := service.NewShelfService(app.db, libraryService)
	objectService := service.NewObjectService(app.db, libraryService, shelfService)

	app.services = &Services{
		Library: libraryService,
		Shelf:   shelfService,
		Object:  objectService,
	}
}
