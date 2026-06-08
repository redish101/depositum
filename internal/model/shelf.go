package model

import "gorm.io/gorm"

type Shelf struct {
	gorm.Model

	Name        string
	Description string

	LibraryID uint
	Library   Library
}
