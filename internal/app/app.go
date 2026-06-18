package app

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/redish101/depositum/internal/config"
	"github.com/redish101/depositum/internal/db"
	v1 "github.com/redish101/depositum/pkg/api/v1"
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

	server    *http.Server
	mux       *http.ServeMux
	container *restful.Container
	listener  net.Listener
	started   chan struct{}

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

	app.mux = http.NewServeMux()

	app.container = restful.NewContainer()

	app.mux.Handle(v1.BasePath+"/", http.StripPrefix(v1.BasePath, app.container))

	app.server = &http.Server{
		Handler: app.mux,
	}

	app.started = make(chan struct{})

	app.initServices()
	app.initHandlers()

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
