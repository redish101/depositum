package model

import (
	v1 "github.com/redish101/depositum/pkg/api/v1"
	"gorm.io/gorm"
)

type ObjectStatus struct {
	Phase     v1.ObjectPhase
	LibraryID uint
	ShelfID   uint
}

type Object struct {
	gorm.Model

	Name        string
	Description string

	CurrentStatus ObjectStatus `gorm:"embedded;embeddedPrefix:current_"`
	DesiredStatus ObjectStatus `gorm:"embedded;embeddedPrefix:desired_"`

	Synced bool
}
