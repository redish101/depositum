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

	v1Group := e.Group("/v1")
	handler.Register(v1Group)

	return e, library.ID
}

func TestListLibraries(t *testing.T) {
	e, _ := setupTestLibraryHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/v1/libraries", nil)
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

func TestGetLibrary(t *testing.T) {
	e, id := setupTestLibraryHandler(t)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/libraries/%d", id), nil)
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

	body, _ := json.Marshal(v1.UpdateLibraryRequest{
		Name:    &updatedName,
		Address: &updatedAddress,
	})

	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/libraries/%d", id), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	// 再 GET 验证
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/libraries/%d", id), nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	var library v1.Library
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &library))

	assert.Equal(t, updatedName, library.Name)
	assert.Equal(t, updatedAddress, library.Address)
}

func TestDeleteLibrary(t *testing.T) {
	e, id := setupTestLibraryHandler(t)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/libraries/%d", id), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}
