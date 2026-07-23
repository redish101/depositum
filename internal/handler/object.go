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
	g.PATCH("/:id/sync", h.Sync)
	g.DELETE("/:id", h.Delete)
}

// Create 创建对象
//
//	@Summary	创建对象
//	@Tags		Object
//	@Accept		json
//	@Produce	json
//	@Param		body	body		v1.CreateObjectRequest	true	"对象"
//	@Success	200		{object}	v1.Object
//	@Failure	400		{object}	v1.ErrorResponse
//	@Failure	500		{object}	v1.ErrorResponse
//	@Router		/objects [post]
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

// List 获取对象列表
//
//	@Summary	获取对象列表
//	@Tags		Object
//	@Produce	json
//	@Param		page		query		int						false	"页码"
//	@Param		pageSize	query		int						false	"每页数量"
//	@Param		query		body		v1.ListObjectRequest	false	"查询条件"
//	@Success	200			{object}	v1.PaginationResponse[v1.Object]
//	@Failure	400			{object}	v1.ErrorResponse
//	@Failure	500			{object}	v1.ErrorResponse
//	@Router		/objects [get]
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

// Get 获取对象
//
//	@Summary	获取对象详情
//	@Tags		Object
//	@Produce	json
//	@Param		id	path		int	true	"对象 ID"
//	@Success	200	{object}	v1.Object
//	@Failure	400	{object}	v1.ErrorResponse
//	@Failure	404	{object}	v1.ErrorResponse
//	@Failure	500	{object}	v1.ErrorResponse
//	@Router		/objects/{id} [get]
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

// Update 更新对象
//
//	@Summary	更新对象
//	@Tags		Object
//	@Accept		json
//	@Produce	json
//	@Param		id		path		int						true	"对象 ID"
//	@Param		body	body		v1.UpdateObjectRequest	true	"对象"
//	@Success	200		{object}	v1.Object
//	@Failure	400		{object}	v1.ErrorResponse
//	@Failure	404		{object}	v1.ErrorResponse
//	@Failure	500		{object}	v1.ErrorResponse
//	@Router		/objects/{id} [patch]
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

// Sync 同步对象状态
//
//	@Summary	同步对象状态
//	@Tags		Object
//	@Produce	json
//	@Param		id	path		int	true	"对象 ID"
//	@Success	200	{object}	v1.Object
//	@Failure	400	{object}	v1.ErrorResponse
//	@Failure	404	{object}	v1.ErrorResponse
//	@Failure	500	{object}	v1.ErrorResponse
//	@Router		/objects/{id}/sync [patch]
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

// Delete 删除对象
//
//	@Summary	删除对象
//	@Tags		Object
//	@Produce	json
//	@Param		id	path	int	true	"对象 ID"
//	@Success	200
//	@Failure	400	{object}	v1.ErrorResponse
//	@Failure	404	{object}	v1.ErrorResponse
//	@Failure	500	{object}	v1.ErrorResponse
//	@Router		/objects/{id} [delete]
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
