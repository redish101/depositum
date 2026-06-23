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

const (
	testLibraryName    = "Test Library"
	testLibraryAddress = "Test Library Address"
)

func setupTestLibraryHandler(t *testing.T) (*echo.Echo, uint) {
	db, err := db.NewMemoryDB()
	require.NoError(t, err)

	libraryService := service.NewLibraryService(db)

	library, err := libraryService.Create(t.Context(), &v1.CreateLibraryRequest{
		Name:    testLibraryName,
		Address: testLibraryAddress,
	})
	require.NoError(t, err)

	handler := NewLibraryHandler(libraryService)

	e := echo.New()
	validator, err := validate.NewValidator()
	require.NoError(t, err)
	e.Validator = validator

	api := e.Group("")
	handler.Register(api)

	return e, library.ID
}

func TestListLibraries(t *testing.T) {
	e, _ := setupTestLibraryHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/libraries", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp v1.PaginationResponse[v1.Library]
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Len(t, resp.Items, 1)
	assert.Equal(t, testLibraryName, resp.Items[0].Name)
	assert.Equal(t, testLibraryAddress, resp.Items[0].Address)
}

func TestCreateLibrary(t *testing.T) {
	e, _ := setupTestLibraryHandler(t)

	reqBody := v1.CreateLibraryRequest{
		Name:    "New Library",
		Address: "New Library Address",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/libraries", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var library v1.Library
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &library))

	assert.Equal(t, reqBody.Name, library.Name)
	assert.Equal(t, reqBody.Address, library.Address)
}

func TestGetLibrary(t *testing.T) {
	e, id := setupTestLibraryHandler(t)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/libraries/%d", id), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var library v1.Library
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &library))

	assert.Equal(t, testLibraryName, library.Name)
	assert.Equal(t, testLibraryAddress, library.Address)
}

func TestUpdateLibrary(t *testing.T) {
	e, id := setupTestLibraryHandler(t)

	updatedName := "Updated Library"
	updatedAddress := "Updated Address"

	t.Run("Success", func(t *testing.T) {
		body, _ := json.Marshal(v1.UpdateLibraryRequest{
			Name:    &updatedName,
			Address: &updatedAddress,
		})

		req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/libraries/%d", id), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		// 再 GET 验证
		req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/libraries/%d", id), nil)
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		var library v1.Library
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &library))

		assert.Equal(t, updatedName, library.Name)
		assert.Equal(t, updatedAddress, library.Address)
	})

	t.Run("NotFound", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/libraries/9999", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("InvalidParams", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/libraries/%d", id), bytes.NewReader([]byte(`{"name":""}`)))
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestDeleteLibrary(t *testing.T) {
	e, id := setupTestLibraryHandler(t)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/libraries/%d", id), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}
