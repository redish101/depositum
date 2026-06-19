package handler

import (
	"errors"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/redish101/depositum/internal/common"
	"github.com/redish101/depositum/internal/service"
)

type LibraryHandler interface {
	Register(container *restful.Container)
	List(req *restful.Request, resp *restful.Response)
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

	ws.Route(ws.GET("").To(l.List))
	ws.Route(ws.GET("/{id}").To(l.Get))

	container.Add(ws)
}

func (l *libraryHandler) List(req *restful.Request, resp *restful.Response) {
	paginationParams := common.ReadPaginationParams(req)
	libraries, err := l.libraryService.List(req.Request.Context(), paginationParams)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteEntity(libraries)
}

func (l *libraryHandler) Get(req *restful.Request, resp *restful.Response) {
	idStr := req.PathParameter("id")

	id, err := common.ReadID(idStr)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}

	library, err := l.libraryService.Get(req.Request.Context(), id)

	if err != nil {
		if errors.Is(err, service.ErrLibraryNotFound) {
			common.WriteError(resp, http.StatusNotFound, err)
			return
		}
		common.WriteError(resp, http.StatusInternalServerError, err)
		return
	}

	resp.WriteEntity(library)
}
