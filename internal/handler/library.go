package handler

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/redish101/depositum/internal/service"
	v1 "github.com/redish101/depositum/pkg/api/v1"
)

type LibraryHandler interface {
	Register(container *restful.Container)
	Get(req *restful.Request, resp *restful.Response)
}

type libraryHandler struct {
	libraryService service.LibraryService
}

func NewLibraryHandler(libraryService service.LibraryService) LibraryHandler {
	return &libraryHandler{
		libraryService,
	}
}

func (l *libraryHandler) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/libraries").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/{id}").To(l.Get))
	container.Add(ws)
}

func (l *libraryHandler) Get(req *restful.Request, resp *restful.Response) {
	l.libraryService.Create(req.Request.Context(), &v1.CreateLibraryRequest{
		Name:    "1",
		Address: "1",
	})
	library, err := l.libraryService.Get(req.Request.Context(), 1)

	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}

	resp.WriteEntity(library)
}
