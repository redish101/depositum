package service

import (
	"context"
	"errors"

	"github.com/redish101/depositum/internal/db"
	"github.com/redish101/depositum/internal/model"
	v1 "github.com/redish101/depositum/pkg/api/v1"
	"gorm.io/gorm"
)

type LibraryService interface {
	Get(ctx context.Context, id uint) (*v1.Library, error)
	List(ctx context.Context, pageParams *v1.PaginationParams) (*v1.PaginationResponse[v1.Library], error)
	Create(ctx context.Context, request *v1.CreateLibraryRequest) (*v1.Library, error)
	Update(ctx context.Context, id uint, request *v1.UpdateLibraryRequest) (*v1.Library, error)
	Delete(ctx context.Context, id uint) error
}

var (
	ErrLibraryNotFound = errors.New("library not found")
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
		return nil, ErrLibraryNotFound
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

func (l *libraryService) List(ctx context.Context, pageParams *v1.PaginationParams) (*v1.PaginationResponse[v1.Library], error) {
	libraries, err := db.Paginate[model.Library](ctx, l.db, pageParams)
	if err != nil {
		return &v1.PaginationResponse[v1.Library]{}, err
	}

	resp := make([]v1.Library, len(libraries.Items))
	for i, library := range libraries.Items {
		resp[i] = v1.Library{
			ID:        library.ID,
			CreatedAt: library.CreatedAt,
			UpdatedAt: library.UpdatedAt,
			Name:      library.Name,
			Address:   library.Address,
		}
	}

	return &v1.PaginationResponse[v1.Library]{
		Items:      resp,
		Page:       libraries.Page,
		PageSize:   libraries.PageSize,
		Total:      libraries.Total,
		TotalPages: libraries.TotalPages,
		HasNext:    libraries.HasNext,
		HasPrev:    libraries.HasPrev,
	}, nil
}

func (l *libraryService) Create(ctx context.Context, req *v1.CreateLibraryRequest) (*v1.Library, error) {
	library := model.Library{
		Name:    req.Name,
		Address: req.Address,
	}

	err := l.db.WithContext(ctx).Create(&library).Error
	if err != nil {
		return nil, err
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

func (l *libraryService) Update(ctx context.Context, id uint, req *v1.UpdateLibraryRequest) (*v1.Library, error) {
	var library model.Library

	err := l.db.WithContext(ctx).First(&library, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrLibraryNotFound
	}

	if req.Name != nil {
		library.Name = *req.Name
	}
	if req.Address != nil {
		library.Address = *req.Address
	}

	err = l.db.WithContext(ctx).Save(&library).Error
	if err != nil {
		return nil, err
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

func (l *libraryService) Delete(ctx context.Context, id uint) error {
	err := l.db.WithContext(ctx).Delete(&model.Library{}, id).Error

	return err
}
