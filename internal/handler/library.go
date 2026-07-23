package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/redish101/depositum/internal/common"
	"github.com/redish101/depositum/internal/service"
	v1 "github.com/redish101/depositum/pkg/api/v1"
)

type LibraryHandler interface {
	Register(group *echo.Group)
	Create(c *echo.Context) error
	List(c *echo.Context) error
	Get(c *echo.Context) error
	Update(c *echo.Context) error
	Delete(c *echo.Context) error
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
	g.POST("", l.Create)
	g.GET("", l.List)
	g.GET("/:id", l.Get)
	g.PATCH("/:id", l.Update)
	g.DELETE("/:id", l.Delete)
}

// Create 创建库
//
//	@Summary	创建库
//	@Tags		Library
//	@Accept		json
//	@Produce	json
//	@Param		body	body		v1.CreateLibraryRequest	true	"库"
//	@Success	200		{object}	v1.Library
//	@Failure	400		{object}	v1.ErrorResponse
//	@Failure	500		{object}	v1.ErrorResponse
//	@Router		/libraries [post]
func (l *libraryHandler) Create(c *echo.Context) error {
	var params v1.CreateLibraryRequest
	if err := c.Bind(&params); err != nil {
		return common.WriteError(c, http.StatusBadRequest, err)
	}
	if err := c.Validate(&params); err != nil {
		return common.WriteError(c, http.StatusBadRequest, err)
	}

	library, err := l.libraryService.Create(c.Request().Context(), &params)
	if err != nil {
		return common.WriteError(c, http.StatusInternalServerError, err)
	}

	return common.WriteEntity(c, library)
}

// List 获取库列表
//
//	@Summary	获取库列表
//	@Tags		Library
//	@Produce	json
//	@Param		page		query		int	false	"页码"
//	@Param		pageSize	query		int	false	"每页数量"
//	@Success	200			{object}	v1.PaginationResponse[v1.Library]
//	@Failure	500			{object}	v1.ErrorResponse
//	@Router		/libraries [get]
func (l *libraryHandler) List(c *echo.Context) error {
	paginationParams := common.ReadPaginationParams(c)
	libraries, err := l.libraryService.List(c.Request().Context(), paginationParams)
	if err != nil {
		return common.WriteError(c, http.StatusInternalServerError, err)
	}
	return common.WriteEntity(c, libraries)
}

// Get 获取库详情
//
//	@Summary	获取库详情
//	@Tags		Library
//	@Produce	json
//	@Param		id	path		int	true	"库 ID"
//	@Success	200	{object}	v1.Library
//	@Failure	400	{object}	v1.ErrorResponse
//	@Failure	404	{object}	v1.ErrorResponse
//	@Failure	500	{object}	v1.ErrorResponse
//	@Router		/libraries/{id} [get]
func (l *libraryHandler) Get(c *echo.Context) error {
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

// Update 更新库
//
//	@Summary	更新库
//	@Tags		Library
//	@Accept		json
//	@Produce	json
//	@Param		id		path		int						true	"库 ID"
//	@Param		body	body		v1.UpdateLibraryRequest	true	"库"
//	@Success	200		{object}	v1.Library
//	@Failure	400		{object}	v1.ErrorResponse
//	@Failure	404		{object}	v1.ErrorResponse
//	@Failure	500		{object}	v1.ErrorResponse
//	@Router		/libraries/{id} [patch]
func (l *libraryHandler) Update(c *echo.Context) error {
	idStr := c.Param("id")
	id, err := common.ReadID(idStr)
	if err != nil {
		return common.WriteError(c, http.StatusBadRequest, err)
	}

	var params v1.UpdateLibraryRequest
	if err := c.Bind(&params); err != nil {
		return common.WriteError(c, http.StatusBadRequest, err)
	}
	if err := c.Validate(&params); err != nil {
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

// Delete 删除库
//
//	@Summary	删除库
//	@Tags		Library
//	@Produce	json
//	@Param		id	path	int	true	"库 ID"
//	@Success	204
//	@Failure	400	{object}	v1.ErrorResponse
//	@Failure	500	{object}	v1.ErrorResponse
//	@Router		/libraries/{id} [delete]
func (l *libraryHandler) Delete(c *echo.Context) error {
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
