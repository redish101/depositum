package db

import (
	"fmt"

	"github.com/redish101/depositum/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var err error

	if config.Debug {
		DB, err = gorm.Open(sqlite.Open("depositum_debug.db"))
	} else {
		DB, err = gorm.Open(sqlite.Open("depositum.db"))
	}

	if err != nil {
		panic(fmt.Sprintf("database initialization failed: %v", err))
	}
}
