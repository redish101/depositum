package service

import (
	"context"

	v1 "github.com/redish101/depositum/pkg/api/v1"
	"gorm.io/gorm"
)

type ShelfService interface {
	Get(ctx context.Context, id uint) (*v1.Library, error)
	List(ctx context.Context, pageParams *v1.PaginationParams) (*v1.PaginationResponse[v1.Library], error)
	Create(ctx context.Context, request *v1.CreateLibraryRequest) (*v1.Library, error)
	Update(ctx context.Context, id uint, request *v1.UpdateLibraryRequest) (*v1.Library, error)
	Delete(ctx context.Context, id uint) error
}

type shelfService struct {
	db *gorm.DB
}
