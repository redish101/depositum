package handler

import (
	"errors"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/redish101/depositum/internal/common"
	"github.com/redish101/depositum/internal/service"
	v1 "github.com/redish101/depositum/pkg/api/v1"
)

type LibraryHandler interface {
	Register(container *restful.Container)
	List(req *restful.Request, resp *restful.Response)
	Get(req *restful.Request, resp *restful.Response)
	Update(req *restful.Request, resp *restful.Response)
	Delete(req *restful.Request, resp *restful.Response)
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
	ws.Route(ws.PATCH("/{id}").To(l.Update))
	ws.Route(ws.DELETE("/{id}").To(l.Delete))

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

func (l *libraryHandler) Update(req *restful.Request, resp *restful.Response) {
	id, err := common.ReadID(req.PathParameter("id"))
	if err != nil {
		common.WriteError(resp, http.StatusBadRequest, err)
		return
	}

	var params v1.UpdateLibraryRequest

	if err := req.ReadEntity(&params); err != nil {
		common.WriteError(resp, http.StatusBadRequest, err)
		return
	}

	updatedLibrary, err := l.libraryService.Update(req.Request.Context(), id, &params)
	if err != nil && errors.Is(err, service.ErrLibraryNotFound) {
		common.WriteError(resp, http.StatusNotFound, err)
		return
	}
	if err != nil {
		common.WriteError(resp, http.StatusInternalServerError, err)
		return
	}
	resp.WriteEntity(updatedLibrary)
}

func (l *libraryHandler) Delete(req *restful.Request, resp *restful.Response) {
	id, err := common.ReadID(req.PathParameter("id"))
	if err != nil {
		common.WriteError(resp, http.StatusBadRequest, err)
		return
	}

	err = l.libraryService.Delete(req.Request.Context(), id)
	if err != nil {
		common.WriteError(resp, http.StatusInternalServerError, err)
		return
	}
	resp.WriteHeader(http.StatusNoContent)
}
