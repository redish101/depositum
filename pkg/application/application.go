package application

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/redish101/depositum/internal/config"
	"github.com/redish101/depositum/internal/db"
	"github.com/redish101/depositum/internal/handler"
	"github.com/redish101/depositum/internal/service"
	"gorm.io/gorm"
)

type Application interface {
	Run() error
}

type application struct {
	config *config.Config

	db *gorm.DB

	echo *echo.Echo

	libraryService service.LibraryService

	libraryHandler handler.LibraryHandler
}

func New(config *config.Config) (Application, error) {
	log.Println("Initializing depositum")

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

	log.Println("Listening on " + addr)

	err := app.echo.Start(addr)

	return err
}
