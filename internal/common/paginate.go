package common

import (
	"strconv"

	"github.com/emicklei/go-restful/v3"
	v1 "github.com/redish101/depositum/pkg/api/v1"
)

func ReadPaginationParams(req *restful.Request) *v1.PaginationParams {
	page := 1
	pageSize := 10

	if p := req.QueryParameter("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	if ps := req.QueryParameter("pageSize"); ps != "" {
		if parsedPageSize, err := strconv.Atoi(ps); err == nil && parsedPageSize > 0 {
			pageSize = parsedPageSize
		}
	}

	return &v1.PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}
}
