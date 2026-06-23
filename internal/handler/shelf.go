package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/redish101/depositum/internal/common"
	"github.com/redish101/depositum/internal/service"
	v1 "github.com/redish101/depositum/pkg/api/v1"
)

type ShelfHandler interface {
	Register(group *echo.Group)
	List(c *echo.Context) error
	Get(c *echo.Context) error
	Update(c *echo.Context) error
	Delete(c *echo.Context) error
}

type shelfHandler struct {
	svc service.ShelfService
}

func NewShelfHandler(svc service.ShelfService) ShelfHandler {
	return &shelfHandler{
		svc: svc,
	}
}

func (h *shelfHandler) Register(group *echo.Group) {
	g := group.Group("/shelves")

	g.GET("", h.List)
	g.GET("/:id", h.Get)
	g.POST("", h.Create)
	g.PATCH("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}

func (h *shelfHandler) Create(c *echo.Context) error {
	var params v1.CreateShelfRequest
	if err := c.Bind(&params); err != nil {
		return common.WriteError(c, http.StatusBadRequest, err)
	}
	if err := c.Validate(&params); err != nil {
		return common.WriteError(c, http.StatusBadRequest, err)
	}

	shelf, err := h.svc.Create(c.Request().Context(), &params)
	if err != nil {
		if errors.Is(err, service.ErrInvalidLibraryID) {
			return common.WriteError(c, http.StatusBadRequest, err)
		}
		return common.WriteError(c, http.StatusInternalServerError, err)
	}

	return common.WriteEntity(c, shelf)
}

func (h *shelfHandler) List(c *echo.Context) error {
	paginationParams := common.ReadPaginationParams(c)
	shelves, err := h.svc.List(c.Request().Context(), paginationParams)
	if err != nil {
		return common.WriteError(c, http.StatusInternalServerError, err)
	}
	return common.WriteEntity(c, shelves)
}

func (h *shelfHandler) Get(c *echo.Context) error {
	id, err := common.ReadID(c.Param("id"))
	if err != nil {
		return common.WriteError(c, http.StatusBadRequest, err)
	}
	shelf, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrShelfNotFound) {
			return common.WriteError(c, http.StatusNotFound, err)
		}
		return common.WriteError(c, http.StatusInternalServerError, err)
	}
	return common.WriteEntity(c, shelf)
}

func (h *shelfHandler) Update(c *echo.Context) error {
	id, err := common.ReadID(c.Param("id"))
	if err != nil {
		return common.WriteError(c, http.StatusBadRequest, err)
	}

	var params v1.UpdateShelfRequest
	if err := c.Bind(&params); err != nil {
		return common.WriteError(c, http.StatusBadRequest, err)
	}
	if err := c.Validate(&params); err != nil {
		return common.WriteError(c, http.StatusBadRequest, err)
	}
	shelf, err := h.svc.Update(c.Request().Context(), id, &params)
	if err != nil {
		if errors.Is(err, service.ErrShelfNotFound) {
			return common.WriteError(c, http.StatusNotFound, err)
		}
		if errors.Is(err, service.ErrInvalidLibraryID) {
			return common.WriteError(c, http.StatusBadRequest, err)
		}
		return common.WriteError(c, http.StatusInternalServerError, err)
	}

	return common.WriteEntity(c, shelf)
}

func (h *shelfHandler) Delete(c *echo.Context) error {
	id, err := common.ReadID(c.Param("id"))
	if err != nil {
		return common.WriteError(c, http.StatusBadRequest, err)
	}

	err = h.svc.Delete(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrShelfNotFound) {
			return common.WriteError(c, http.StatusNotFound, err)
		}
		return common.WriteError(c, http.StatusInternalServerError, err)
	}

	return common.WriteEntity(c, nil)
}
