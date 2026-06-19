package common

import (
	"github.com/emicklei/go-restful/v3"
	v1 "github.com/redish101/depositum/pkg/api/v1"
)

func WriteError(resp *restful.Response, statusCode int, err error) {
	resp.WriteHeaderAndEntity(statusCode, v1.ErrorResponse{
		Error: err.Error(),
	})
}
