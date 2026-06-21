package app

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/redish101/depositum/internal/config"
	"github.com/redish101/depositum/internal/db"
	"gorm.io/gorm"
)

type App interface {
	Run() error
	Addr() string
	Started() <-chan struct{}
	Shutdown(ctx context.Context) error
}

type app struct {
	config *config.Config

	db *gorm.DB

	server   *http.Server
	echo     *echo.Echo
	listener net.Listener
	started  chan struct{}

	services *Services
}

func New(config *config.Config) (App, error) {
	log.Println("初始化 depositum")

	app := &app{}

	app.config = config

	db, err := db.NewDB(app.config)
	if err != nil {
		return nil, err
	}
	app.db = db

	app.echo = NewEcho()

	// 设置全局前缀 /v1
	v1Group := app.echo.Group("/api/v1")
	app.initServices()
	app.initHandlers(v1Group)

	app.server = &http.Server{
		Handler: app.echo,
	}

	app.started = make(chan struct{})

	return app, nil
}

func (app *app) Run() error {
	addr := app.config.Host + ":" + app.config.Port
	log.Println("监听 " + addr)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	app.listener = ln

	close(app.started)

	err = app.server.Serve(ln)

	return err
}

func (app *app) Addr() string {
	if app.listener != nil {
		return app.listener.Addr().String()
	}
	return ""
}

func (app *app) Started() <-chan struct{} {
	return app.started
}

func (app *app) Shutdown(ctx context.Context) error {
	return app.server.Shutdown(ctx)
}
