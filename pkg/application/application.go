package application

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/redish101/depositum/internal/config"
	"github.com/redish101/depositum/internal/db"
	"github.com/redish101/depositum/internal/handler"
	"github.com/redish101/depositum/internal/service"
	"gorm.io/gorm"
)

type Application interface {
	Run() error
	Shutdown() error
	NewEchoContext(r *http.Request, w http.ResponseWriter) echo.Context
}

type application struct {
	config *config.Config

	db *gorm.DB

	echo *echo.Echo

	libraryService service.LibraryService

	libraryHandler handler.LibraryHandler
}

func New(config *config.Config) (Application, error) {
	log.Println("初始化 depositum")

	app := &application{}

	app.config = config

	db, err := db.NewDB(app.config)
	if err != nil {
		return nil, err
	}
	app.db = db

	app.echo = NewEcho()

	app.libraryService = service.NewLibraryService(app.db)

	app.libraryHandler = handler.NewLibraryHandler(app.libraryService)

	registerRoutes(app, app.echo)

	return app, nil
}

func (app *application) Run() error {
	addr := app.config.Host + ":" + app.config.Port

	log.Println("监听 " + addr)

	err := app.echo.Start(addr)

	return err
}

func (app *application) Shutdown() error {
	return app.echo.Shutdown(context.Background())
}

func (app *application) NewEchoContext(r *http.Request, w http.ResponseWriter) echo.Context {
	return app.echo.NewContext(r, w)
}
