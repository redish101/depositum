package app

import (
	"context"
	"log"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/redish101/depositum/internal/config"
	"github.com/redish101/depositum/internal/db"
	"gorm.io/gorm"
)

type App interface {
	Run() error
	Shutdown(ctx context.Context) error
}

type app struct {
	config *config.Config

	db *gorm.DB

	server    *http.Server
	container *restful.Container

	services *Services
	handlers *Handlers
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

	app.container = restful.NewContainer()

	addr := app.config.Host + ":" + app.config.Port
	app.server = &http.Server{
		Handler: app.container,
		Addr:    addr,
	}

	app.initServices()
	app.initHandlers()

	return app, nil
}

func (app *app) Run() error {
	addr := app.config.Host + ":" + app.config.Port

	log.Println("监听 " + addr)

	err := app.server.ListenAndServe()

	return err
}

func (app *app) Shutdown(ctx context.Context) error {
	return app.server.Shutdown(ctx)
}
