package service

import (
	"testing"

	v1 "github.com/redish101/depositum/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestLibraryService(t *testing.T) LibraryService {
	return NewLibraryService(newTestDB(t))
}

func setupTestLibrary(t *testing.T) LibraryService {
	svc := newTestLibraryService(t)

	library, err := svc.Create(t.Context(), &v1.CreateLibraryRequest{
		Name:    "Test Library",
		Address: "Nusquamium",
	})
	require.NotNil(t, library)
	require.NoError(t, err)

	return svc
}

func TestGetLibrary(t *testing.T) {
	svc := setupTestLibrary(t)

	t.Run("Success", func(t *testing.T) {
		library, err := svc.Get(t.Context(), 1)
		assert.NotNil(t, library)
		assert.NoError(t, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		_, err := svc.Get(t.Context(), 999)
		assert.ErrorIs(t, err, ErrLibraryNotFound)
	})
}

func TestListLibrary(t *testing.T) {
	svc := setupTestLibrary(t)

	t.Run("Success", func(t *testing.T) {
		resp, err := svc.List(t.Context(), &v1.PaginationParams{
			Page:     1,
			PageSize: 10,
		})
		assert.NoError(t, err)
		assert.Len(t, resp.Items, 1)
	})
}

func TestCreateLibrary(t *testing.T) {
	svc := newTestLibraryService(t)

	t.Run("Success", func(t *testing.T) {
		library, err := svc.Create(t.Context(), &v1.CreateLibraryRequest{
			Name:    "New Library",
			Address: "Nusquamium",
		})
		assert.NotNil(t, library)
		assert.NoError(t, err)
	})
}

func TestUpdateLibrary(t *testing.T) {
	svc := setupTestLibrary(t)

	t.Run("Success", func(t *testing.T) {
		library, err := svc.Update(t.Context(), 1, &v1.UpdateLibraryRequest{
			Name:    "Updated Library",
			Address: "Verilocus",
		})
		assert.NotNil(t, library)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Library", library.Name)
	})

	t.Run("NotFound", func(t *testing.T) {
		_, err := svc.Update(t.Context(), 999, &v1.UpdateLibraryRequest{
			Name:    "Non-existent Library",
			Address: "Nusquamium",
		})
		assert.ErrorIs(t, err, ErrLibraryNotFound)
	})
}

func TestDeleteLibrary(t *testing.T) {
	svc := setupTestLibrary(t)

	t.Run("Success", func(t *testing.T) {
		err := svc.Delete(t.Context(), 1)
		assert.NoError(t, err)

		_, err = svc.Get(t.Context(), 1)
		assert.ErrorIs(t, err, ErrLibraryNotFound)
	})

	// t.Run("NotFound", func(t *testing.T) {
	// 	err := svc.Delete(t.Context(), 999)
	// 	assert.ErrorIs(t, err, ErrLibraryNotFound)
	// })
}
