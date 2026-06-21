package handler

import (
	"github.com/labstack/echo/v5"
	"github.com/redish101/depositum/internal/common"
	v1 "github.com/redish101/depositum/pkg/api/v1"
)

type HealthzHandler interface {
	Register(v1 *echo.Group)
	Get(c *echo.Context) error
}

type healthzHandler struct{}

func NewHealthzHandler() HealthzHandler {
	return &healthzHandler{}
}

func (h *healthzHandler) Register(v1 *echo.Group) {
	v1.GET("/healthz", h.Get)
}

func (h *healthzHandler) Get(c *echo.Context) error {
	status := v1.HealthzResponse{
		Ok: true,
	}

	return common.WriteEntity(c, status)
}
