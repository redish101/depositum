package db

import (
	"testing"

	"github.com/redish101/depositum/internal/config"
	v1 "github.com/redish101/depositum/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type TestModel struct {
	gorm.Model

	Name string
}

func createTestData(t *testing.T, db *gorm.DB) {
	testData := []TestModel{
		{Name: "Alice"}, {Name: "Bob"}, {Name: "Charlie"},
		{Name: "David"}, {Name: "Eve"}, {Name: "Frank"},
		{Name: "Grace"}, {Name: "Hannah"}, {Name: "Ivy"}, {Name: "Jack"},
	}
	err := db.Create(&testData).Error
	require.NoError(t, err)
}

func TestPaginate(t *testing.T) {
	ctx := t.Context()

	cfg := &config.Config{
		DSN: ":memory:",
	}

	db, err := NewDB(cfg)
	require.NoError(t, err)
	require.NotNil(t, db)

	err = db.AutoMigrate(&TestModel{})
	require.NoError(t, err)

	createTestData(t, db)

	// 测试正常分页
	paginationParams := v1.PaginationParams{
		Page:     2,
		PageSize: 3,
	}

	resp, err := PaginateWithQuery[TestModel](ctx, db, &paginationParams, func(db *gorm.DB) *gorm.DB { return db })
	assert.NoError(t, err)
	require.NotNil(t, resp)

	expectedNames := []string{"David", "Eve", "Frank"}
	actualNames := make([]string, len(resp.Data))
	for i, model := range resp.Data {
		actualNames[i] = model.Name
	}
	assert.Equal(t, expectedNames, actualNames)
	assert.Equal(t, 2, resp.Page)
	assert.Equal(t, 3, resp.PageSize)
	assert.Equal(t, int64(10), resp.Total)
	assert.Equal(t, 4, resp.TotalPages)
	assert.True(t, resp.HasNext)
	assert.True(t, resp.HasPrev)

	// 测试 pageSize 超过最大限制
	paginationParams = v1.PaginationParams{
		Page:     1,
		PageSize: MaximumPageSize + 1,
	}

	resp, err = PaginateWithQuery[TestModel](ctx, db, &paginationParams, func(db *gorm.DB) *gorm.DB { return db })
	assert.NoError(t, err)
	assert.Equal(t, MaximumPageSize, resp.PageSize)
}

func TestPaginateWithQuery(t *testing.T) {
	ctx := t.Context()

	cfg := &config.Config{
		DSN: ":memory:",
	}

	db, err := NewDB(cfg)
	require.NoError(t, err)
	require.NotNil(t, db)

	err = db.AutoMigrate(&TestModel{})
	require.NoError(t, err)

	createTestData(t, db)

	paginationParams := v1.PaginationParams{
		Page:     1,
		PageSize: 1,
	}

	resp, err := PaginateWithQuery[TestModel](ctx, db, &paginationParams, func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", "Alice")
	})

	assert.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 1, len(resp.Data))
	assert.Equal(t, "Alice", resp.Data[0].Name)
	assert.Equal(t, int64(1), resp.Total)
	assert.Equal(t, 1, resp.TotalPages)
	assert.False(t, resp.HasNext)
	assert.False(t, resp.HasPrev)
}
