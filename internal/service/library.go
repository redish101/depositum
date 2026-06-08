package service

import (
	"context"
	"errors"

	"github.com/redish101/depositum/internal/model"
	v1 "github.com/redish101/depositum/pkg/api/v1"
	"gorm.io/gorm"
)

type LibraryService interface {
	Get(ctx context.Context, id uint) (*v1.Library, error)
}

const (
	ErrLibraryNotFound = "library not found"
)

type libraryService struct {
	db *gorm.DB
}

func NewLibraryService(db *gorm.DB) LibraryService {
	return &libraryService{
		db,
	}
}

func (l *libraryService) Get(ctx context.Context, id uint) (*v1.Library, error) {
	var library model.Library

	err := l.db.WithContext(ctx).First(&library, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New(ErrLibraryNotFound)
	}

	resp := &v1.Library{
		ID:        library.ID,
		CreatedAt: library.CreatedAt,
		UpdatedAt: library.UpdatedAt,
		Name:      library.Name,
		Address:   library.Address,
	}

	return resp, nil
}
