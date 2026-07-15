package resolver

import "github.com/redish101/depositum/internal/service"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	libraryService service.LibraryService
	shelfService   service.ShelfService
	objectService  service.ObjectService
}

func NewResolver(libraryService service.LibraryService, shelfService service.ShelfService, objectService service.ObjectService) *Resolver {
	return &Resolver{
		libraryService: libraryService,
		shelfService:   shelfService,
		objectService:  objectService,
	}
}
