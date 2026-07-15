package resolver

import (
	"github.com/redish101/depositum/internal/service"
	v1 "github.com/redish101/depositum/pkg/api/v1"
	"github.com/redish101/depositum/pkg/graph/model"
)

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

func pageInfo[T any](serviceResponse *v1.PaginationResponse[T]) *model.PageInfo {
	return &model.PageInfo{
		Total:   int32(serviceResponse.Total),
		HasNext: serviceResponse.HasNext,
		HasPrev: serviceResponse.HasPrev,
	}
}

func covertPageParams(input *model.PageParams) *v1.PaginationParams {
	if input == nil {
		return &v1.PaginationParams{
			Page:     1,
			PageSize: 10,
		}
	}
	return &v1.PaginationParams{
		Page:     int(input.Page),
		PageSize: int(input.PageSize),
	}
}
