package main

import (
	"log"

	"github.com/labstack/echo/v4"

	"github.com/redish101/depositum/internal/config"
	"github.com/redish101/depositum/internal/db"
	"github.com/redish101/depositum/server/handler"
)

func main() {
	log.Println("Initializing depositum")
	config.Init()

	if config.Debug {
		log.Println("Running as debug mode.")
	}

	db.Init()

	e := echo.New()

	e.HideBanner = true
	e.HidePort = true
	e.HTTPErrorHandler = handler.HTTPErrorHandler

	addr := config.Host + ":" + config.Port

	log.Println("Listening on " + addr)
	e.Start(addr)
}
