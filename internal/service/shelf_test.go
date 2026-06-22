package service

import (
	"testing"

	"github.com/redish101/depositum/internal/db"
	v1 "github.com/redish101/depositum/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestShelfService(t *testing.T) (ShelfService, *v1.Shelf, *v1.Library) {
	db, err := db.NewMemoryDB()
	require.NoError(t, err)

	libraryService := NewLibraryService(db)
	library, err := libraryService.Create(t.Context(), &v1.CreateLibraryRequest{
		Name:    "Test Library",
		Address: "Test Library",
	})

	svc := NewShelfService(db, libraryService)

	require.NotNil(t, library)
	require.NoError(t, err)

	shelf, err := svc.Create(t.Context(), &v1.CreateShelfRequest{
		Name:        "Test Shelf",
		Description: "Test Shelf",
		LibraryID:   library.ID,
	})

	return svc, shelf, library
}

func TestListShelves(t *testing.T) {
	svc, testShelf, _ := newTestShelfService(t)

	t.Run("Success", func(t *testing.T) {
		resp, err := svc.List(t.Context(), &v1.PaginationParams{
			Page:     1,
			PageSize: 10,
		})
		assert.NoError(t, err)
		assert.Len(t, resp.Items, 1)
		assert.Equal(t, testShelf.Name, resp.Items[0].Name)
	})
}

func TestGetShelf(t *testing.T) {
	svc, testShelf, _ := newTestShelfService(t)

	t.Run("Success", func(t *testing.T) {
		shelf, err := svc.Get(t.Context(), testShelf.ID)
		assert.Equal(t, testShelf.Name, shelf.Name)
		assert.NoError(t, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		_, err := svc.Get(t.Context(), 999)
		assert.ErrorIs(t, err, ErrShelfNotFound)
	})
}

func TestCreateShelf(t *testing.T) {
	svc, _, testLibrary := newTestShelfService(t)

	t.Run("Success", func(t *testing.T) {
		expectedShelf := &v1.Shelf{
			Name:        "Another Test Shelf",
			Description: "Another Test Shelf",
		}
		shelf, err := svc.Create(t.Context(), &v1.CreateShelfRequest{
			Name:        expectedShelf.Name,
			Description: expectedShelf.Description,
			LibraryID:   testLibrary.ID,
		})
		assert.NoError(t, err)
		assert.Equal(t, expectedShelf.Name, shelf.Name)
	})

	t.Run("InvalidLibraryID", func(t *testing.T) {
		_, err := svc.Create(t.Context(), &v1.CreateShelfRequest{
			LibraryID: 1999,
		})
		assert.ErrorIs(t, err, ErrInvalidLibraryID)
	})
}

func TestUpdateShelf(t *testing.T) {
	svc, testShelf, _ := newTestShelfService(t)

	t.Run("Success", func(t *testing.T) {
		expectedShelf := &v1.Shelf{
			Name:        "Updated Test Shelf",
			Description: "Updated Test Shelf",
		}
		shelf, err := svc.Update(t.Context(), testShelf.ID, &v1.UpdateShelfRequest{
			Name:        &expectedShelf.Name,
			Description: &expectedShelf.Description,
		})
		assert.Equal(t, expectedShelf.Name, shelf.Name)
		assert.NoError(t, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		_, err := svc.Update(t.Context(), 999, &v1.UpdateShelfRequest{})
		assert.ErrorIs(t, err, ErrShelfNotFound)
	})

	t.Run("InvalidLibraryID", func(t *testing.T) {
		_, err := svc.Update(t.Context(), testShelf.ID, &v1.UpdateShelfRequest{
			LibraryID: func() *uint { id := uint(1999); return &id }(),
		})
		assert.ErrorIs(t, err, ErrInvalidLibraryID)
	})
}

func TestDeleteShelf(t *testing.T) {
	svc, testShelf, _ := newTestShelfService(t)

	t.Run("Success", func(t *testing.T) {
		err := svc.Delete(t.Context(), testShelf.ID)
		assert.NoError(t, err)

		_, err = svc.Get(t.Context(), testShelf.ID)
		assert.ErrorIs(t, err, ErrShelfNotFound)
	})
}
