package main

import (
	"log"

	"github.com/redish101/depositum/internal/config"
	"github.com/redish101/depositum/pkg/application"
)

func main() {
	cfg := config.FromEnv()

	app, err := application.New(cfg)
	if err != nil {
		panic(err)
	}

	err = app.Run()

	log.Fatalln(err)
}
