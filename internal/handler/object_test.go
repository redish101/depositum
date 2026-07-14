package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/redish101/depositum/internal/db"
	"github.com/redish101/depositum/internal/service"
	"github.com/redish101/depositum/internal/validate"
	v1 "github.com/redish101/depositum/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestObjectHandler(t *testing.T) (*echo.Echo, *v1.Object, *v1.Shelf, *v1.Library) {
	database, err := db.NewMemoryDB()
	require.NoError(t, err)

	libraryService := service.NewLibraryService(database)
	shelfService := service.NewShelfService(database, libraryService)
	objectService := service.NewObjectService(database, libraryService, shelfService)

	library, err := libraryService.Create(t.Context(), &v1.CreateLibraryRequest{
		Name:    "Test Library",
		Address: "Test Library Address",
	})
	require.NoError(t, err)

	shelf, err := shelfService.Create(t.Context(), &v1.CreateShelfRequest{
		Name:      "Test Shelf",
		LibraryID: library.ID,
	})
	require.NoError(t, err)

	object, err := objectService.Create(t.Context(), &v1.CreateObjectRequest{
		Name:        "Test Object",
		Description: "Test Object Description",
		DesiredStatus: v1.ObjectStatus{
			LibraryID: library.ID,
			ShelfID:   shelf.ID,
		},
		SyncNow: false,
	})
	require.NoError(t, err)

	handler := NewObjectHandler(objectService)

	e := echo.New()
	validator, err := validate.NewValidator()
	require.NoError(t, err)
	e.Validator = validator

	api := e.Group("")
	handler.Register(api)

	return e, object, shelf, library
}

func TestCreateObject(t *testing.T) {
	h, _, testShelf, testLibrary := setupTestObjectHandler(t)

	t.Run("Success", func(t *testing.T) {
		reqBody, _ := json.Marshal(v1.CreateObjectRequest{
			Name:        "New Test Object",
			Description: "New Description",
			DesiredStatus: v1.ObjectStatus{
				LibraryID: testLibrary.ID,
				ShelfID:   testShelf.ID,
			},
			SyncNow: true,
		})
		req := httptest.NewRequest(http.MethodPost, "/objects", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		require.Equalf(t, http.StatusOK, rec.Code, "Response: %s", rec.Body.String())
	})

	t.Run("InvalidShelfID", func(t *testing.T) {
		reqBody, _ := json.Marshal(v1.CreateObjectRequest{
			Name:        "Bad Shelf Object",
			Description: "Desc",
			DesiredStatus: v1.ObjectStatus{
				LibraryID: testLibrary.ID,
				ShelfID:   1000,
			},
			SyncNow: true,
		})
		req := httptest.NewRequest(http.MethodPost, "/objects", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InvalidParams", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/objects", bytes.NewReader([]byte(`{"name":""}`)))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestListObjects(t *testing.T) {
	h, testObject, _, _ := setupTestObjectHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/objects", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	require.Equalf(t, http.StatusOK, rec.Code, "Response: %s", rec.Body.String())

	var resp v1.PaginationResponse[v1.Object]
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.NotEmpty(t, resp.Items)
	assert.Equal(t, testObject.Name, resp.Items[0].Name)
	assert.Equal(t, testObject.Description, resp.Items[0].Description)
	assert.Equal(t, testObject.DesiredStatus.LibraryID, resp.Items[0].DesiredStatus.LibraryID)
	assert.Equal(t, testObject.DesiredStatus.ShelfID, resp.Items[0].DesiredStatus.ShelfID)
}

func TestGetObject(t *testing.T) {
	h, testObject, _, _ := setupTestObjectHandler(t)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/objects/%d", testObject.ID), nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	require.Equalf(t, http.StatusOK, rec.Code, "Response: %s", rec.Body.String())

	var obj v1.Object
	err := json.Unmarshal(rec.Body.Bytes(), &obj)
	assert.NoError(t, err)

	assert.Equal(t, testObject.Name, obj.Name)
	assert.Equal(t, testObject.Description, obj.Description)
	assert.Equal(t, testObject.DesiredStatus.LibraryID, obj.DesiredStatus.LibraryID)
	assert.Equal(t, testObject.DesiredStatus.ShelfID, obj.DesiredStatus.ShelfID)
}

func TestUpdateObject(t *testing.T) {
	h, testObject, _, _ := setupTestObjectHandler(t)

	t.Run("Success", func(t *testing.T) {
		updatedName := "Updated Object Name"
		reqBody, _ := json.Marshal(v1.UpdateObjectRequest{
			Name: &updatedName,
		})
		req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/objects/%d", testObject.ID), bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		require.Equalf(t, http.StatusOK, rec.Code, "Response: %s", rec.Body.String())

		var obj v1.Object
		err := json.Unmarshal(rec.Body.Bytes(), &obj)
		assert.NoError(t, err)

		assert.Equal(t, updatedName, obj.Name)
	})

	t.Run("NotFound", func(t *testing.T) {
		reqBody, _ := json.Marshal(v1.UpdateObjectRequest{})
		req := httptest.NewRequest(http.MethodPatch, "/objects/9999", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("InvalidShelfID", func(t *testing.T) {
		invalidDesiredStatus := v1.ObjectStatus{
			LibraryID: testObject.DesiredStatus.LibraryID,
			ShelfID:   9999, // non-existent
		}
		reqBody, _ := json.Marshal(v1.UpdateObjectRequest{
			DesiredStatus: &invalidDesiredStatus,
		})
		req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/objects/%d", testObject.ID), bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestSyncObject(t *testing.T) {
	h, testObject, _, _ := setupTestObjectHandler(t)

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/objects/%d/sync", testObject.ID), nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		require.Equalf(t, http.StatusOK, rec.Code, "Response: %s", rec.Body.String())

		var obj v1.Object
		err := json.Unmarshal(rec.Body.Bytes(), &obj)
		assert.NoError(t, err)

		assert.True(t, obj.Synced)
		assert.Equal(t, obj.DesiredStatus, obj.CurrentStatus)
	})

	t.Run("NotFound", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/objects/9999/sync", nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestDeleteObject(t *testing.T) {
	h, testObject, _, _ := setupTestObjectHandler(t)

	// 测试删除
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/objects/%d", testObject.ID), nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	require.Equalf(t, http.StatusOK, rec.Code, "Response: %s", rec.Body.String())

	// 再 GET 验证确实已删除
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/objects/%d", testObject.ID), nil)
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
