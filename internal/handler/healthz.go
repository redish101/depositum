package handler

import (
	"github.com/emicklei/go-restful/v3"
	v1 "github.com/redish101/depositum/pkg/api/v1"
)

type HealthzHandler interface {
	Register(container *restful.Container)
	Get(req *restful.Request, resp *restful.Response)
}

type healthzHandler struct{}

func NewHealthzHandler() HealthzHandler {
	return &healthzHandler{}
}

func (h *healthzHandler) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/healthz").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("").To(h.Get))
	container.Add(ws)
}

func (h *healthzHandler) Get(req *restful.Request, resp *restful.Response) {
	status := v1.HealthzResponse{
		Ok: true,
	}
	resp.WriteEntity(status)
}
