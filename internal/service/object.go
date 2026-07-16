package service

import (
	"context"
	"errors"
	"reflect"

	"github.com/redish101/depositum/internal/db"
	"github.com/redish101/depositum/internal/model"
	v1 "github.com/redish101/depositum/pkg/api/v1"
	"gorm.io/gorm"
)

type ObjectService interface {
	List(ctx context.Context, pageParams *v1.PaginationParams, listParams *v1.ListObjectRequest) (*v1.PaginationResponse[*v1.Object], error)
	Get(ctx context.Context, id uint) (*v1.Object, error)
	Create(ctx context.Context, request *v1.CreateObjectRequest) (*v1.Object, error)
	Update(ctx context.Context, id uint, request *v1.UpdateObjectRequest) (*v1.Object, error)
	Sync(ctx context.Context, id uint) (*v1.Object, error)
	Delete(ctx context.Context, id uint) error
}

var (
	ErrInvalidListObjectsParams = errors.New("invalid list objects params")
	ErrInvalidShelfID           = errors.New("invalid shelf ID")
	ErrObjectNotFound           = errors.New("object not found")
	ErrShelfNotInLibrary        = errors.New("shelf does not belong to library")
)

type objectService struct {
	db             *gorm.DB
	libraryService LibraryService
	shelfService   ShelfService
}

func NewObjectService(db *gorm.DB, libraryService LibraryService, shelfService ShelfService) ObjectService {
	return &objectService{db: db, libraryService: libraryService, shelfService: shelfService}
}

func (s *objectService) List(ctx context.Context, pageParams *v1.PaginationParams, listParams *v1.ListObjectRequest) (*v1.PaginationResponse[*v1.Object], error) {
	if listParams.SyncedOnly && listParams.UnsyncedOnly {
		return nil, ErrInvalidListObjectsParams
	}

	query := func(db *gorm.DB) *gorm.DB {
		if listParams.SyncedOnly {
			db = db.Where("synced = ?", true)
		}

		if listParams.UnsyncedOnly {
			db = db.Where("synced = ?", false)
		}

		if listParams.LibraryID != nil {
			db = db.Where("current_library_id = ?", *listParams.LibraryID)
		}

		if listParams.ShelfID != nil {
			db = db.Where("current_shelf_id = ?", *listParams.ShelfID)
		}

		return db
	}

	objects, err := db.PaginateWithQuery[model.Object](ctx, s.db, pageParams, query)

	if err != nil {
		return nil, err
	}

	resp := make([]*v1.Object, len(objects.Items))
	for i, obj := range objects.Items {
		resp[i] = &v1.Object{
			ID:            obj.ID,
			CreatedAt:     obj.CreatedAt,
			UpdatedAt:     obj.UpdatedAt,
			Name:          obj.Name,
			Description:   obj.Description,
			Synced:        obj.Synced,
			CurrentStatus: v1.ObjectStatus(obj.CurrentStatus),
			DesiredStatus: v1.ObjectStatus(obj.DesiredStatus),
		}
	}

	return &v1.PaginationResponse[*v1.Object]{
		Items:      resp,
		Page:       objects.Page,
		PageSize:   objects.PageSize,
		Total:      objects.Total,
		TotalPages: objects.TotalPages,
		HasNext:    objects.HasNext,
	}, nil
}

func (s *objectService) Get(ctx context.Context, id uint) (*v1.Object, error) {
	var obj model.Object

	err := s.db.WithContext(ctx).First(&obj, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrObjectNotFound
		}
		return nil, err
	}

	return &v1.Object{
		ID:            obj.ID,
		CreatedAt:     obj.CreatedAt,
		UpdatedAt:     obj.UpdatedAt,
		Name:          obj.Name,
		Description:   obj.Description,
		Synced:        obj.Synced,
		CurrentStatus: v1.ObjectStatus(obj.CurrentStatus),
		DesiredStatus: v1.ObjectStatus(obj.DesiredStatus),
	}, nil
}

func isSynced(obj *model.Object) bool {
	return reflect.DeepEqual(obj.CurrentStatus, obj.DesiredStatus)
}

func writeSyncStatus(obj *model.Object) {
	if isSynced(obj) {
		obj.Synced = true
	} else {
		obj.Synced = false
	}
}

func checkShelfExists(ctx context.Context, shelfService ShelfService, libraryID uint, shelfID uint) error {
	shelf, err := shelfService.Get(ctx, shelfID)
	if err != nil {
		if errors.Is(err, ErrShelfNotFound) {
			return ErrInvalidShelfID
		}
		return err
	}

	if shelf.LibraryID != libraryID {
		return ErrShelfNotInLibrary
	}

	return nil
}

func checkLibraryAndShelfExist(ctx context.Context, libraryService LibraryService, shelfService ShelfService, libraryID uint, shelfID uint) error {
	err := checkLibraryExists(ctx, libraryService, libraryID)
	if err != nil {
		return err
	}

	err = checkShelfExists(ctx, shelfService, libraryID, shelfID)
	if err != nil {
		return err
	}

	return nil
}

func (s *objectService) Create(ctx context.Context, request *v1.CreateObjectRequest) (*v1.Object, error) {
	err := checkLibraryAndShelfExist(ctx, s.libraryService, s.shelfService, request.DesiredStatus.LibraryID, request.DesiredStatus.ShelfID)
	if err != nil {
		return nil, err
	}

	obj := model.Object{
		Name:          request.Name,
		Description:   request.Description,
		DesiredStatus: model.ObjectStatus(request.DesiredStatus),
		CurrentStatus: model.ObjectStatus{
			Phase: v1.ObjectPhaseCreated,
		},
	}

	if request.SyncNow {
		obj.CurrentStatus = model.ObjectStatus(request.DesiredStatus)
	}

	writeSyncStatus(&obj)

	err = s.db.WithContext(ctx).Create(&obj).Error
	if err != nil {
		return nil, err
	}

	return &v1.Object{
		ID:            obj.ID,
		CreatedAt:     obj.CreatedAt,
		UpdatedAt:     obj.UpdatedAt,
		Name:          obj.Name,
		Description:   obj.Description,
		Synced:        obj.Synced,
		CurrentStatus: v1.ObjectStatus(obj.CurrentStatus),
		DesiredStatus: v1.ObjectStatus(obj.DesiredStatus),
	}, nil
}

func (s *objectService) Update(ctx context.Context, id uint, request *v1.UpdateObjectRequest) (*v1.Object, error) {
	var obj model.Object

	err := s.db.WithContext(ctx).First(&obj, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrObjectNotFound
		}
		return nil, err
	}

	if request.Name != nil {
		obj.Name = *request.Name
	}
	if request.Description != nil {
		obj.Description = *request.Description
	}
	if request.DesiredStatus != nil {
		err := checkLibraryAndShelfExist(ctx, s.libraryService, s.shelfService, request.DesiredStatus.LibraryID, request.DesiredStatus.ShelfID)
		if err != nil {
			return nil, err
		}
		obj.DesiredStatus = model.ObjectStatus(*request.DesiredStatus)
	}
	if request.CurrentStatus != nil {
		err := checkLibraryAndShelfExist(ctx, s.libraryService, s.shelfService, request.CurrentStatus.LibraryID, request.CurrentStatus.ShelfID)
		if err != nil {
			return nil, err
		}
		obj.CurrentStatus = model.ObjectStatus(*request.CurrentStatus)
	}

	if request.SyncNow != nil && *request.SyncNow {
		obj.CurrentStatus = obj.DesiredStatus
	}

	writeSyncStatus(&obj)

	err = s.db.WithContext(ctx).Save(&obj).Error
	if err != nil {
		return nil, err
	}

	return &v1.Object{
		ID:            obj.ID,
		CreatedAt:     obj.CreatedAt,
		UpdatedAt:     obj.UpdatedAt,
		Name:          obj.Name,
		Description:   obj.Description,
		Synced:        obj.Synced,
		CurrentStatus: v1.ObjectStatus(obj.CurrentStatus),
		DesiredStatus: v1.ObjectStatus(obj.DesiredStatus),
	}, nil
}

func (s *objectService) Sync(ctx context.Context, id uint) (*v1.Object, error) {
	var obj model.Object

	err := s.db.WithContext(ctx).First(&obj, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrObjectNotFound
		}
		return nil, err
	}

	err = checkLibraryAndShelfExist(ctx, s.libraryService, s.shelfService, obj.DesiredStatus.LibraryID, obj.DesiredStatus.ShelfID)
	if err != nil {
		return nil, err
	}

	obj.CurrentStatus = obj.DesiredStatus
	writeSyncStatus(&obj)

	err = s.db.WithContext(ctx).Save(&obj).Error
	if err != nil {
		return nil, err
	}

	return &v1.Object{
		ID:            obj.ID,
		CreatedAt:     obj.CreatedAt,
		UpdatedAt:     obj.UpdatedAt,
		Name:          obj.Name,
		Description:   obj.Description,
		Synced:        obj.Synced,
		CurrentStatus: v1.ObjectStatus(obj.CurrentStatus),
		DesiredStatus: v1.ObjectStatus(obj.DesiredStatus),
	}, nil
}

func (s *objectService) Delete(ctx context.Context, id uint) error {
	err := s.db.WithContext(ctx).Delete(&model.Object{}, id).Error
	return err
}
