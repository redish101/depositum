package main

import (
	"log"

	"github.com/redish101/depositum/internal/app"
	"github.com/redish101/depositum/internal/config"
)

//	@Title			depositum API
//	@Version		1.0
//	@Description	档案管理系统 API
//	@BasePath		/api/v1
//
// the main function
func main() {
	cfg := config.FromEnv()

	app, err := app.New(cfg)
	if err != nil {
		panic(err)
	}

	err = app.Run()

	log.Fatalln(err)
}
