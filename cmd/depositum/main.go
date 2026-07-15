package main

import (
	"log"

	"github.com/redish101/depositum/internal/app"
	"github.com/redish101/depositum/internal/config"
)

func main() {
	cfg := config.FromEnv()

	app, err := app.New(cfg)
	if err != nil {
		panic(err)
	}

	err = app.Run()

	log.Fatalln(err)
}
