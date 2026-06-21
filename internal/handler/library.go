package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/redish101/depositum/internal/common"
	"github.com/redish101/depositum/internal/service"
	v1 "github.com/redish101/depositum/pkg/api/v1"
)

type LibraryHandler interface {
	Register(group *echo.Group)
	List(c echo.Context) error
	Get(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
}

type libraryHandler struct {
	libraryService service.LibraryService
}

func NewLibraryHandler(libraryService service.LibraryService) LibraryHandler {
	return &libraryHandler{
		libraryService: libraryService,
	}
}

func (l *libraryHandler) Register(group *echo.Group) {
	g := group.Group("/libraries")
	g.GET("", l.List)
	g.GET("/:id", l.Get)
	g.PATCH("/:id", l.Update)
	g.DELETE("/:id", l.Delete)
}

func (l *libraryHandler) List(c echo.Context) error {
	paginationParams := common.ReadPaginationParams(c)
	libraries, err := l.libraryService.List(c.Request().Context(), paginationParams)
	if err != nil {
		return common.WriteError(c, http.StatusInternalServerError, err)
	}
	return common.WriteEntity(c, libraries)
}

func (l *libraryHandler) Get(c echo.Context) error {
	idStr := c.Param("id")
	id, err := common.ReadID(idStr)
	if err != nil {
		return common.WriteError(c, http.StatusBadRequest, err)
	}

	library, err := l.libraryService.Get(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrLibraryNotFound) {
			return common.WriteError(c, http.StatusNotFound, err)
		}
		return common.WriteError(c, http.StatusInternalServerError, err)
	}
	return common.WriteEntity(c, library)
}

func (l *libraryHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := common.ReadID(idStr)
	if err != nil {
		return common.WriteError(c, http.StatusBadRequest, err)
	}

	var params v1.UpdateLibraryRequest
	if err := c.Bind(&params); err != nil {
		return common.WriteError(c, http.StatusBadRequest, err)
	}

	updatedLibrary, err := l.libraryService.Update(c.Request().Context(), id, &params)
	if err != nil {
		if errors.Is(err, service.ErrLibraryNotFound) {
			return common.WriteError(c, http.StatusNotFound, err)
		}
		return common.WriteError(c, http.StatusInternalServerError, err)
	}
	return common.WriteEntity(c, updatedLibrary)
}

func (l *libraryHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := common.ReadID(idStr)
	if err != nil {
		return common.WriteError(c, http.StatusBadRequest, err)
	}

	err = l.libraryService.Delete(c.Request().Context(), id)
	if err != nil {
		return common.WriteError(c, http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusNoContent) // 无内容响应，无实体，可直接使用 echo 原生方法
}
