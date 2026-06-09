package db

import (
	"context"

	v1 "github.com/redish101/depositum/pkg/api/v1"
	"gorm.io/gorm"
)

const (
	MaximumPageSize = 100
)

func normalize(p *v1.PaginationParams) {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if p.PageSize > MaximumPageSize {
		p.PageSize = MaximumPageSize
	}
}

func PaginateWithQuery[T any](ctx context.Context, db *gorm.DB, params *v1.PaginationParams, query func(*gorm.DB) *gorm.DB) (*v1.PaginationResponse[T], error) {
	normalize(params)

	var total int64

	// 应用查询条件并计算总数
	countDB := query(db.Model(new(T)))
	if err := countDB.Count(&total).Error; err != nil {
		return nil, err
	}

	var result []T

	// 应用查询条件并查询数据
	offset := (params.Page - 1) * params.PageSize
	queryDB := query(db)
	if err := queryDB.Offset(offset).Limit(params.PageSize).Find(&result).Error; err != nil {
		return nil, err
	}

	// 计算分页信息
	totalPages := int((total + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &v1.PaginationResponse[T]{
		Items:      result,
		Page:       params.Page,
		PageSize:   params.PageSize,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    params.Page < totalPages,
		HasPrev:    params.Page > 1,
	}, nil
}

func Paginate[T any](ctx context.Context, db *gorm.DB, params *v1.PaginationParams) (*v1.PaginationResponse[T], error) {
	return PaginateWithQuery[T](ctx, db, params, func(db *gorm.DB) *gorm.DB { return db })
}
