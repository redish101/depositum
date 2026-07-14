package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/redish101/depositum/internal/common"
	"github.com/redish101/depositum/internal/service"
	v1 "github.com/redish101/depositum/pkg/api/v1"
)

type ObjectHandler interface {
	Register(group *echo.Group)
	List(c *echo.Context) error
	Get(c *echo.Context) error
	Create(c *echo.Context) error
	Update(c *echo.Context) error
	Sync(c *echo.Context) error
	Delete(c *echo.Context) error
}

type objectHandler struct {
	svc service.ObjectService
}

func NewObjectHandler(svc service.ObjectService) ObjectHandler {
	return &objectHandler{
		svc: svc,
	}
}

func (h *objectHandler) Register(group *echo.Group) {
	g := group.Group("/objects")

	g.GET("", h.List)
	g.GET("/:id", h.Get)
	g.POST("", h.Create)
	g.PATCH("/:id", h.Update)
	g.POST("/:id/sync", h.Sync)
	g.DELETE("/:id", h.Delete)
}

func (h *objectHandler) Create(c *echo.Context) error {
	var params v1.CreateObjectRequest
	if err := c.Bind(&params); err != nil {
		return common.WriteError(c, http.StatusBadRequest, err)
	}
	if err := c.Validate(&params); err != nil {
		return common.WriteError(c, http.StatusBadRequest, err)
	}

	object, err := h.svc.Create(c.Request().Context(), &params)
	if err != nil {
		if errors.Is(err, service.ErrInvalidShelfID) ||
			errors.Is(err, service.ErrShelfNotInLibrary) ||
			errors.Is(err, service.ErrLibraryNotFound) {
			return common.WriteError(c, http.StatusBadRequest, err)
		}
		return common.WriteError(c, http.StatusInternalServerError, err)
	}

	return common.WriteEntity(c, object)
}

func (h *objectHandler) List(c *echo.Context) error {
	paginationParams := common.ReadPaginationParams(c)

	var listParams v1.ListObjectRequest

	if err := c.Bind(&listParams); err != nil {
		return common.WriteError(c, http.StatusBadRequest, err)
	}
	if err := c.Validate(&listParams); err != nil {
		return common.WriteError(c, http.StatusBadRequest, err)
	}

	objects, err := h.svc.List(c.Request().Context(), paginationParams, &listParams)
	if err != nil {
		if errors.Is(err, service.ErrInvalidListObjectsParams) {
			return common.WriteError(c, http.StatusBadRequest, err)
		}
		return common.WriteError(c, http.StatusInternalServerError, err)
	}

	return common.WriteEntity(c, objects)
}

func (h *objectHandler) Get(c *echo.Context) error {
	id, err := common.ReadID(c.Param("id"))
	if err != nil {
		return common.WriteError(c, http.StatusBadRequest, err)
	}

	object, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrObjectNotFound) {
			return common.WriteError(c, http.StatusNotFound, err)
		}
		return common.WriteError(c, http.StatusInternalServerError, err)
	}

	return common.WriteEntity(c, object)
}

func (h *objectHandler) Update(c *echo.Context) error {
	id, err := common.ReadID(c.Param("id"))
	if err != nil {
		return common.WriteError(c, http.StatusBadRequest, err)
	}

	var params v1.UpdateObjectRequest
	if err := c.Bind(&params); err != nil {
		return common.WriteError(c, http.StatusBadRequest, err)
	}
	if err := c.Validate(&params); err != nil {
		return common.WriteError(c, http.StatusBadRequest, err)
	}

	object, err := h.svc.Update(c.Request().Context(), id, &params)
	if err != nil {
		if errors.Is(err, service.ErrObjectNotFound) {
			return common.WriteError(c, http.StatusNotFound, err)
		}
		if errors.Is(err, service.ErrInvalidShelfID) ||
			errors.Is(err, service.ErrShelfNotInLibrary) ||
			errors.Is(err, service.ErrLibraryNotFound) {
			return common.WriteError(c, http.StatusBadRequest, err)
		}
		return common.WriteError(c, http.StatusInternalServerError, err)
	}

	return common.WriteEntity(c, object)
}

func (h *objectHandler) Sync(c *echo.Context) error {
	id, err := common.ReadID(c.Param("id"))
	if err != nil {
		return common.WriteError(c, http.StatusBadRequest, err)
	}

	object, err := h.svc.Sync(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrObjectNotFound) {
			return common.WriteError(c, http.StatusNotFound, err)
		}
		if errors.Is(err, service.ErrInvalidShelfID) ||
			errors.Is(err, service.ErrShelfNotInLibrary) ||
			errors.Is(err, service.ErrLibraryNotFound) {
			return common.WriteError(c, http.StatusBadRequest, err)
		}
		return common.WriteError(c, http.StatusInternalServerError, err)
	}

	return common.WriteEntity(c, object)
}

func (h *objectHandler) Delete(c *echo.Context) error {
	id, err := common.ReadID(c.Param("id"))
	if err != nil {
		return common.WriteError(c, http.StatusBadRequest, err)
	}

	err = h.svc.Delete(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrObjectNotFound) {
			return common.WriteError(c, http.StatusNotFound, err)
		}
		return common.WriteError(c, http.StatusInternalServerError, err)
	}

	return common.WriteEntity(c, nil)
}
