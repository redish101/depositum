package service

import (
	"testing"

	"github.com/redish101/depositum/internal/db"
	v1 "github.com/redish101/depositum/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestObjectServices(t *testing.T) (ObjectService, LibraryService, ShelfService) {
	t.Helper()

	database, err := db.NewMemoryDB()
	require.NoError(t, err)

	libSvc := NewLibraryService(database)

	shelfSvc := NewShelfService(database, libSvc)
	objSvc := NewObjectService(database, libSvc, shelfSvc)

	return objSvc, libSvc, shelfSvc
}

func setupTestEnvironment(t *testing.T) (ObjectService, LibraryService, ShelfService, *v1.Library, *v1.Shelf) {
	t.Helper()

	objSvc, libSvc, shelfSvc := newTestObjectServices(t)

	library, err := libSvc.Create(t.Context(), &v1.CreateLibraryRequest{
		Name:    "Test Library",
		Address: "Nusquamium",
	})
	require.NoError(t, err)

	shelf, err := shelfSvc.Create(t.Context(), &v1.CreateShelfRequest{
		Name:      "Test Shelf",
		LibraryID: library.ID,
	})
	require.NoError(t, err)

	return objSvc, libSvc, shelfSvc, library, shelf
}

func setupTestObject(t *testing.T) (ObjectService, *v1.Library, *v1.Shelf, *v1.Object) {
	t.Helper()

	objSvc, _, _, library, shelf := setupTestEnvironment(t)

	object, err := objSvc.Create(t.Context(), &v1.CreateObjectRequest{
		Name:        "Test Object",
		Description: "A test object description",
		DesiredStatus: v1.ObjectStatus{
			LibraryID: library.ID,
			ShelfID:   shelf.ID,
		},
		SyncNow: false,
	})
	require.NoError(t, err)
	require.NotNil(t, object)

	return objSvc, library, shelf, object
}

func TestGetObject(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc, _, _, object := setupTestObject(t)

		got, err := svc.Get(t.Context(), object.ID)

		require.NoError(t, err)
		require.NotNil(t, got)

		assert.Equal(t, object.ID, got.ID)
		assert.Equal(t, "Test Object", got.Name)
		assert.Equal(t, "A test object description", got.Description)
		assert.False(t, got.Synced)
	})

	t.Run("NotFound", func(t *testing.T) {
		svc, _, _ := newTestObjectServices(t)

		_, err := svc.Get(t.Context(), 999)

		assert.ErrorIs(t, err, ErrObjectNotFound)
	})
}

func TestListObject(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		svc, _, _ := newTestObjectServices(t)

		resp, err := svc.List(t.Context(), &v1.PaginationParams{
			Page:     1,
			PageSize: 10,
		}, &v1.ListObjectRequest{})

		require.NoError(t, err)
		require.NotNil(t, resp)

		assert.Empty(t, resp.Items)
		assert.Equal(t, int64(0), resp.Total)
	})

	t.Run("Success", func(t *testing.T) {
		svc, library, shelf, object := setupTestObject(t)

		resp, err := svc.List(t.Context(), &v1.PaginationParams{
			Page:     1,
			PageSize: 10,
		}, &v1.ListObjectRequest{})

		require.NoError(t, err)
		require.NotNil(t, resp)

		assert.Len(t, resp.Items, 1)
		assert.Equal(t, int64(1), resp.Total)

		assert.Equal(t, object.ID, resp.Items[0].ID)
		assert.Equal(t, "Test Object", resp.Items[0].Name)
		assert.Equal(t, library.ID, resp.Items[0].DesiredStatus.LibraryID)
		assert.Equal(t, shelf.ID, resp.Items[0].DesiredStatus.ShelfID)
	})

	t.Run("InvalidParams", func(t *testing.T) {
		svc, _, _ := newTestObjectServices(t)

		_, err := svc.List(t.Context(), &v1.PaginationParams{
			Page:     1,
			PageSize: 10,
		}, &v1.ListObjectRequest{
			SyncedOnly:   true,
			UnsyncedOnly: true,
		})

		assert.ErrorIs(t, err, ErrInvalidListObjectsParams)
	})

	t.Run("Filters", func(t *testing.T) {
		svc, library, shelf, object := setupTestObject(t)

		// Create a synced object
		_, err := svc.Create(t.Context(), &v1.CreateObjectRequest{
			Name: "Synced Object",
			DesiredStatus: v1.ObjectStatus{
				LibraryID: library.ID,
				ShelfID:   shelf.ID,
			},
			SyncNow: true,
		})
		require.NoError(t, err)

		syncedResp, err := svc.List(t.Context(), &v1.PaginationParams{Page: 1, PageSize: 10}, &v1.ListObjectRequest{SyncedOnly: true})
		require.NoError(t, err)
		assert.Len(t, syncedResp.Items, 1)
		assert.Equal(t, "Synced Object", syncedResp.Items[0].Name)

		unsyncedResp, err := svc.List(t.Context(), &v1.PaginationParams{Page: 1, PageSize: 10}, &v1.ListObjectRequest{UnsyncedOnly: true})
		require.NoError(t, err)
		assert.Len(t, unsyncedResp.Items, 1)
		assert.Equal(t, object.ID, unsyncedResp.Items[0].ID)

		libResp, err := svc.List(t.Context(), &v1.PaginationParams{Page: 1, PageSize: 10}, &v1.ListObjectRequest{LibraryID: &library.ID})
		require.NoError(t, err)
		assert.Len(t, libResp.Items, 1)
		assert.Equal(t, "Synced Object", libResp.Items[0].Name)
	})
}

func TestCreateObject(t *testing.T) {
	t.Run("Success_SyncNow", func(t *testing.T) {
		svc, _, _, library, shelf := setupTestEnvironment(t)

		object, err := svc.Create(t.Context(), &v1.CreateObjectRequest{
			Name:        "New Object",
			Description: "Desc",
			DesiredStatus: v1.ObjectStatus{
				LibraryID: library.ID,
				ShelfID:   shelf.ID,
			},
			SyncNow: true,
		})

		require.NoError(t, err)
		require.NotNil(t, object)

		assert.NotZero(t, object.ID)
		assert.Equal(t, "New Object", object.Name)
		assert.True(t, object.Synced)
		assert.Equal(t, object.DesiredStatus, object.CurrentStatus)
	})

	t.Run("ShelfNotInLibrary", func(t *testing.T) {
		svc, libSvc, _, _, shelf := setupTestEnvironment(t)

		otherLibrary, err := libSvc.Create(t.Context(), &v1.CreateLibraryRequest{
			Name: "Other Library",
		})
		require.NoError(t, err)

		_, err = svc.Create(t.Context(), &v1.CreateObjectRequest{
			Name: "Mismatched Object",
			DesiredStatus: v1.ObjectStatus{
				LibraryID: otherLibrary.ID,
				ShelfID:   shelf.ID,
			},
		})

		assert.ErrorIs(t, err, ErrShelfNotInLibrary)
	})

	t.Run("InvalidShelfID", func(t *testing.T) {
		svc, _, _, library, _ := setupTestEnvironment(t)

		_, err := svc.Create(t.Context(), &v1.CreateObjectRequest{
			Name: "Invalid Shelf Object",
			DesiredStatus: v1.ObjectStatus{
				LibraryID: library.ID,
				ShelfID:   999,
			},
		})

		assert.ErrorIs(t, err, ErrInvalidShelfID)
	})
}

func TestUpdateObject(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc, library, shelf, object := setupTestObject(t)

		name := "Updated Object"
		desc := "Updated Desc"
		desiredStatus := v1.ObjectStatus{
			LibraryID: library.ID,
			ShelfID:   shelf.ID,
		}

		updated, err := svc.Update(t.Context(), object.ID, &v1.UpdateObjectRequest{
			Name:          &name,
			Description:   &desc,
			DesiredStatus: &desiredStatus,
		})

		require.NoError(t, err)
		require.NotNil(t, updated)

		assert.Equal(t, object.ID, updated.ID)
		assert.Equal(t, "Updated Object", updated.Name)
		assert.Equal(t, "Updated Desc", updated.Description)

		got, err := svc.Get(t.Context(), object.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Object", got.Name)
	})

	t.Run("PartialUpdate", func(t *testing.T) {
		svc, _, _, object := setupTestObject(t)

		name := "Only Name Changed"

		updated, err := svc.Update(t.Context(), object.ID, &v1.UpdateObjectRequest{
			Name: &name,
		})

		require.NoError(t, err)
		require.NotNil(t, updated)

		assert.Equal(t, "Only Name Changed", updated.Name)
		assert.Equal(t, "A test object description", updated.Description) // Should remain unchanged
	})

	t.Run("NotFound", func(t *testing.T) {
		svc, _, _ := newTestObjectServices(t)

		name := "Non-existent"

		_, err := svc.Update(t.Context(), 999, &v1.UpdateObjectRequest{
			Name: &name,
		})

		assert.ErrorIs(t, err, ErrObjectNotFound)
	})
}

func TestSyncObject(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc, _, _, object := setupTestObject(t)

		assert.False(t, object.Synced)
		assert.NotEqual(t, object.DesiredStatus, object.CurrentStatus)

		syncedObj, err := svc.Sync(t.Context(), object.ID)

		require.NoError(t, err)
		require.NotNil(t, syncedObj)

		// Statuses should match and Synced should be true
		assert.True(t, syncedObj.Synced)
		assert.Equal(t, syncedObj.DesiredStatus, syncedObj.CurrentStatus)

		// Verify DB persistence
		got, err := svc.Get(t.Context(), object.ID)
		require.NoError(t, err)
		assert.True(t, got.Synced)
		assert.Equal(t, got.DesiredStatus, got.CurrentStatus)
	})

	t.Run("NotFound", func(t *testing.T) {
		svc, _, _ := newTestObjectServices(t)

		_, err := svc.Sync(t.Context(), 999)

		assert.ErrorIs(t, err, ErrObjectNotFound)
	})
}

func TestDeleteObject(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc, _, _, object := setupTestObject(t)

		err := svc.Delete(t.Context(), object.ID)
		require.NoError(t, err)

		_, err = svc.Get(t.Context(), object.ID)
		assert.ErrorIs(t, err, ErrObjectNotFound)
	})

	t.Run("DeleteNonExistent", func(t *testing.T) {
		svc, _, _ := newTestObjectServices(t)

		err := svc.Delete(t.Context(), 999)
		assert.NoError(t, err)
	})
}
