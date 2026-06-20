package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emicklei/go-restful/v3"
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

func setupTestLibraryHandler(t *testing.T) (*restful.Container, uint) {
	db, err := db.NewMemoryDB()
	require.NoError(t, err)

	libraryService := service.NewLibraryService(db)

	library, err := libraryService.Create(t.Context(), &v1.CreateLibraryRequest{
		Name:    testLibraryName,
		Address: testLibraryAddress,
	})
	require.NoError(t, err)

	handler := NewLibraryHandler(libraryService)

	container := restful.NewContainer()
	handler.Register(container)

	return container, library.ID
}

func TestListLibraries(t *testing.T) {
	container, _ := setupTestLibraryHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/libraries", nil)
	rec := httptest.NewRecorder()
	container.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)

	var resp v1.PaginationResponse[v1.Library]
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Len(t, resp.Items, 1)
	assert.Equal(t, testLibraryName, resp.Items[0].Name)
	assert.Equal(t, testLibraryAddress, resp.Items[0].Address)
}

func TestGetLibrary(t *testing.T) {
	container, id := setupTestLibraryHandler(t)

	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/libraries/%d", id),
		nil,
	)

	rec := httptest.NewRecorder()

	container.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var library v1.Library
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &library))

	assert.Equal(t, testLibraryName, library.Name)
	assert.Equal(t, testLibraryAddress, library.Address)
}

func TestUpdateLibrary(t *testing.T) {
	container, id := setupTestLibraryHandler(t)

	updatedName := "Updated Library"
	updatedAddress := "Updated Address"

	body, _ := json.Marshal(v1.UpdateLibraryRequest{
		Name:    &updatedName,
		Address: &updatedAddress,
	})

	req := httptest.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("/libraries/%d", id),
		bytes.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	container.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	// 再 GET 验证
	req = httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/libraries/%d", id),
		nil,
	)

	rec = httptest.NewRecorder()

	container.ServeHTTP(rec, req)

	var library v1.Library
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &library))

	assert.Equal(t, updatedName, library.Name)
	assert.Equal(t, updatedAddress, library.Address)
}

func TestDeleteLibrary(t *testing.T) {
	container, id := setupTestLibraryHandler(t)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/libraries/%d", id), nil)
	rec := httptest.NewRecorder()
	container.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}
