package ui

import (
	"embed"
	"io/fs"
)

//go:embed client/*
var files embed.FS

func FS() (fs.FS, error) {
	sub, err := fs.Sub(files, "client")
	return sub, err
}
