package service

import (
	"context"
	"errors"

	"github.com/redish101/depositum/internal/db"
	"github.com/redish101/depositum/internal/model"
	v1 "github.com/redish101/depositum/pkg/api/v1"
	"gorm.io/gorm"
)

var ErrShelfNotFound = errors.New("shelf not found")
var ErrInvalidLibraryID = errors.New("invalid library ID")

type ShelfService interface {
	List(ctx context.Context, pageParams *v1.PaginationParams) (*v1.PaginationResponse[v1.Shelf], error)
	Get(ctx context.Context, id uint) (*v1.Shelf, error)
	Create(ctx context.Context, request *v1.CreateShelfRequest) (*v1.Shelf, error)
	Update(ctx context.Context, id uint, request *v1.UpdateShelfRequest) (*v1.Shelf, error)
	Delete(ctx context.Context, id uint) error
}

type shelfService struct {
	db             *gorm.DB
	libraryService LibraryService
}

func NewShelfService(db *gorm.DB, libraryService LibraryService) ShelfService {
	return &shelfService{db: db, libraryService: libraryService}
}

func (s *shelfService) List(ctx context.Context, pageParams *v1.PaginationParams) (*v1.PaginationResponse[v1.Shelf], error) {
	shelves, err := db.Paginate[model.Shelf](ctx, s.db, pageParams)
	if err != nil {
		return &v1.PaginationResponse[v1.Shelf]{}, err
	}

	resp := make([]v1.Shelf, len(shelves.Items))
	for i, shelf := range shelves.Items {
		resp[i] = v1.Shelf{
			ID:          shelf.ID,
			CreatedAt:   shelf.CreatedAt,
			UpdatedAt:   shelf.UpdatedAt,
			Name:        shelf.Name,
			Description: shelf.Description,
			LibraryID:   shelf.LibraryID,
		}
	}

	return &v1.PaginationResponse[v1.Shelf]{
		Page:       shelves.Page,
		PageSize:   shelves.PageSize,
		Total:      shelves.Total,
		TotalPages: shelves.TotalPages,
		Items:      resp,
	}, nil
}

func (s *shelfService) Get(ctx context.Context, id uint) (*v1.Shelf, error) {
	var shelf model.Shelf

	err := s.db.WithContext(ctx).First(&shelf, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrShelfNotFound
	}

	resp := &v1.Shelf{
		ID:          shelf.ID,
		CreatedAt:   shelf.CreatedAt,
		UpdatedAt:   shelf.UpdatedAt,
		Name:        shelf.Name,
		Description: shelf.Description,
		LibraryID:   shelf.LibraryID,
	}

	return resp, nil
}

func checkLibraryExists(ctx context.Context, libraryService LibraryService, libraryID uint) error {
	_, err := libraryService.Get(ctx, libraryID)

	if errors.Is(err, ErrLibraryNotFound) {
		return ErrInvalidLibraryID
	}

	return nil
}

func (s *shelfService) Create(ctx context.Context, request *v1.CreateShelfRequest) (*v1.Shelf, error) {
	err := checkLibraryExists(ctx, s.libraryService, request.LibraryID)
	if err != nil {
		return nil, err
	}

	shelf := model.Shelf{
		Name:        request.Name,
		Description: request.Description,
		LibraryID:   request.LibraryID,
	}

	err = s.db.WithContext(ctx).Create(&shelf).Error
	if err != nil {
		return nil, err
	}

	resp := &v1.Shelf{
		ID:          shelf.ID,
		CreatedAt:   shelf.CreatedAt,
		UpdatedAt:   shelf.UpdatedAt,
		Name:        shelf.Name,
		Description: shelf.Description,
		LibraryID:   shelf.LibraryID,
	}

	return resp, nil
}

func (s *shelfService) Update(ctx context.Context, id uint, request *v1.UpdateShelfRequest) (*v1.Shelf, error) {
	if request.LibraryID != nil {
		err := checkLibraryExists(ctx, s.libraryService, *request.LibraryID)
		if err != nil {
			return nil, err
		}
	}

	var shelf model.Shelf

	err := s.db.WithContext(ctx).First(&shelf, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrShelfNotFound
	}

	if request.Name != nil {
		shelf.Name = *request.Name
	}
	if request.Description != nil {
		shelf.Description = *request.Description
	}
	if request.LibraryID != nil {
		shelf.LibraryID = *request.LibraryID
	}

	err = s.db.WithContext(ctx).Save(&shelf).Error
	if err != nil {
		return nil, err
	}

	resp := &v1.Shelf{
		ID:          shelf.ID,
		CreatedAt:   shelf.CreatedAt,
		UpdatedAt:   shelf.UpdatedAt,
		Name:        shelf.Name,
		Description: shelf.Description,
		LibraryID:   shelf.LibraryID,
	}

	return resp, nil
}

func (s *shelfService) Delete(ctx context.Context, id uint) error {
	err := s.db.WithContext(ctx).Delete(&model.Shelf{}, id).Error

	return err
}
