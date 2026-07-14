package service

import (
	"testing"

	"github.com/redish101/depositum/internal/db"
	v1 "github.com/redish101/depositum/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestLibraryService(t *testing.T) LibraryService {
	t.Helper()

	database, err := db.NewMemoryDB()
	require.NoError(t, err)

	return NewLibraryService(database)
}

func setupTestLibrary(t *testing.T) (LibraryService, *v1.Library) {
	t.Helper()

	svc := newTestLibraryService(t)

	library, err := svc.Create(t.Context(), &v1.CreateLibraryRequest{
		Name:    "Test Library",
		Address: "Nusquamium",
	})
	require.NoError(t, err)
	require.NotNil(t, library)

	return svc, library
}

func TestGetLibrary(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc, library := setupTestLibrary(t)

		got, err := svc.Get(t.Context(), library.ID)

		require.NoError(t, err)
		require.NotNil(t, got)

		assert.Equal(t, library.ID, got.ID)
		assert.Equal(t, "Test Library", got.Name)
		assert.Equal(t, "Nusquamium", got.Address)
	})

	t.Run("NotFound", func(t *testing.T) {
		svc := newTestLibraryService(t)

		_, err := svc.Get(t.Context(), 999)

		assert.ErrorIs(t, err, ErrLibraryNotFound)
	})
}

func TestListLibrary(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		svc := newTestLibraryService(t)

		resp, err := svc.List(t.Context(), &v1.PaginationParams{
			Page:     1,
			PageSize: 10,
		})

		require.NoError(t, err)
		require.NotNil(t, resp)

		assert.Empty(t, resp.Items)
		assert.Equal(t, int64(0), resp.Total)
	})

	t.Run("Success", func(t *testing.T) {
		svc, library := setupTestLibrary(t)

		resp, err := svc.List(t.Context(), &v1.PaginationParams{
			Page:     1,
			PageSize: 10,
		})

		require.NoError(t, err)
		require.NotNil(t, resp)

		assert.Len(t, resp.Items, 1)
		assert.Equal(t, int64(1), resp.Total)

		assert.Equal(t, library.ID, resp.Items[0].ID)
		assert.Equal(t, "Test Library", resp.Items[0].Name)
		assert.Equal(t, "Nusquamium", resp.Items[0].Address)
	})
}

func TestCreateLibrary(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := newTestLibraryService(t)

		library, err := svc.Create(t.Context(), &v1.CreateLibraryRequest{
			Name:    "New Library",
			Address: "Nusquamium",
		})

		require.NoError(t, err)
		require.NotNil(t, library)

		assert.NotZero(t, library.ID)
		assert.Equal(t, "New Library", library.Name)
		assert.Equal(t, "Nusquamium", library.Address)

		got, err := svc.Get(t.Context(), library.ID)

		require.NoError(t, err)
		require.NotNil(t, got)

		assert.Equal(t, library.ID, got.ID)
		assert.Equal(t, "New Library", got.Name)
		assert.Equal(t, "Nusquamium", got.Address)
	})
}

func TestUpdateLibrary(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc, library := setupTestLibrary(t)

		name := "Updated Library"
		address := "Verilocus"

		updated, err := svc.Update(t.Context(), library.ID, &v1.UpdateLibraryRequest{
			Name:    &name,
			Address: &address,
		})

		require.NoError(t, err)
		require.NotNil(t, updated)

		assert.Equal(t, library.ID, updated.ID)
		assert.Equal(t, "Updated Library", updated.Name)
		assert.Equal(t, "Verilocus", updated.Address)

		got, err := svc.Get(t.Context(), library.ID)

		require.NoError(t, err)
		require.NotNil(t, got)

		assert.Equal(t, "Updated Library", got.Name)
		assert.Equal(t, "Verilocus", got.Address)
	})

	t.Run("PartialUpdate", func(t *testing.T) {
		svc, library := setupTestLibrary(t)

		name := "Only Name Changed"

		updated, err := svc.Update(t.Context(), library.ID, &v1.UpdateLibraryRequest{
			Name: &name,
		})

		require.NoError(t, err)
		require.NotNil(t, updated)

		assert.Equal(t, "Only Name Changed", updated.Name)
		assert.Equal(t, "Nusquamium", updated.Address)
	})

	t.Run("NotFound", func(t *testing.T) {
		svc := newTestLibraryService(t)

		name := "Non-existent"
		address := "Nowhere"

		_, err := svc.Update(t.Context(), 999, &v1.UpdateLibraryRequest{
			Name:    &name,
			Address: &address,
		})

		assert.ErrorIs(t, err, ErrLibraryNotFound)
	})
}

func TestDeleteLibrary(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc, library := setupTestLibrary(t)

		err := svc.Delete(t.Context(), library.ID)
		require.NoError(t, err)

		_, err = svc.Get(t.Context(), library.ID)
		assert.ErrorIs(t, err, ErrLibraryNotFound)
	})

	t.Run("DeleteNonExistent", func(t *testing.T) {
		svc := newTestLibraryService(t)

		err := svc.Delete(t.Context(), 999)
		assert.NoError(t, err)
	})
}