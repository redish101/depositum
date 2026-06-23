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

func setupTestShelfHandler(t *testing.T) (*echo.Echo, *v1.Shelf, *v1.Library) {
	db, err := db.NewMemoryDB()
	require.NoError(t, err)

	libraryService := service.NewLibraryService(db)
	shelfService := service.NewShelfService(db, libraryService)

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

	handler := NewShelfHandler(shelfService)

	e := echo.New()
	validator, err := validate.NewValidator()
	require.NoError(t, err)
	e.Validator = validator

	api := e.Group("")
	handler.Register(api)

	return e, shelf, library
}

func TestCreateShelf(t *testing.T) {
	h, _, testLibrary := setupTestShelfHandler(t)

	t.Run("Success", func(t *testing.T) {
		reqBody, _ := json.Marshal(v1.CreateShelfRequest{
			Name:        "Test Shelf",
			Description: "Test Shelf",
			LibraryID:   testLibrary.ID,
		})
		req := httptest.NewRequest(http.MethodPost, "/shelves", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("InvalidLibraryID", func(t *testing.T) {
		reqBody, _ := json.Marshal(v1.CreateShelfRequest{
			Name:        "",
			Description: "",
			LibraryID:   1000,
		})
		req := httptest.NewRequest(http.MethodPost, "/shelves", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InvalidParams", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/shelves", bytes.NewReader([]byte(`{"name":""}`)))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestListShelves(t *testing.T) {
	h, testShelf, _ := setupTestShelfHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/shelves", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp v1.PaginationResponse[v1.Shelf]
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, testShelf.Name, resp.Items[0].Name)
	assert.Equal(t, testShelf.Description, resp.Items[0].Description)
	assert.Equal(t, testShelf.LibraryID, resp.Items[0].LibraryID)
}

func TestGetShelf(t *testing.T) {
	h, testShelf, _ := setupTestShelfHandler(t)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/shelves/%d", testShelf.ID), nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var shelf v1.Shelf
	err := json.Unmarshal(rec.Body.Bytes(), &shelf)
	assert.NoError(t, err)

	assert.Equal(t, testShelf.Name, shelf.Name)
	assert.Equal(t, testShelf.Description, shelf.Description)
	assert.Equal(t, testShelf.LibraryID, shelf.LibraryID)
}

func TestUpdateShelf(t *testing.T) {
	h, testShelf, _ := setupTestShelfHandler(t)

	t.Run("Success", func(t *testing.T) {
		updatedShelfName := "Updated Shelf"
		reqBody, _ := json.Marshal(v1.UpdateShelfRequest{
			Name: &updatedShelfName,
		})
		req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/shelves/%d", testShelf.ID), bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var shelf v1.Shelf
		err := json.Unmarshal(rec.Body.Bytes(), &shelf)
		assert.NoError(t, err)

		assert.Equal(t, updatedShelfName, shelf.Name)
	})

	t.Run("NotFound", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/shelves/9999", nil)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("InvalidLibraryID", func(t *testing.T) {
		invalidLibraryID := uint(9999)
		reqBody, _ := json.Marshal(v1.UpdateShelfRequest{
			LibraryID: &invalidLibraryID,
		})
		req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/shelves/%d", testShelf.ID), bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestDeleteShelf(t *testing.T) {
	h, testShelf, _ := setupTestShelfHandler(t)
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/shelves/%d", testShelf.ID), nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	// 再 GET 验证
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/shelves/%d", testShelf.ID), nil)
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
