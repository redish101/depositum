package common

import (
	"strconv"

	"github.com/labstack/echo/v4"
	v1 "github.com/redish101/depositum/pkg/api/v1"
)

func ReadPaginationParams(c echo.Context) *v1.PaginationParams {
	page := 1
	pageSize := 10

	if p := c.QueryParam("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	if ps := c.QueryParam("pageSize"); ps != "" {
		if parsedPageSize, err := strconv.Atoi(ps); err == nil && parsedPageSize > 0 {
			pageSize = parsedPageSize
		}
	}

	return &v1.PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}
}
