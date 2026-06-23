package model

import "gorm.io/gorm"

type ObjectStatus struct {
	Phase     string
	LibraryID uint
	ShelfID   uint
}

type Object struct {
	gorm.Model

	Name        string
	Description string

	CurrentStatus ObjectStatus `gorm:"embedded;embeddedPrefix:current_"`
	DesiredStatus ObjectStatus `gorm:"embedded;embeddedPrefix:desired_"`
}
