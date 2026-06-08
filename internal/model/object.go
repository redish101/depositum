package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	ObjectCreateOperation   = "create"
	ObjectUpdateOperation   = "update"
	ObjectMoveOperation     = "move"
	ObjectWithDrawOperation = "withdraw"
	ObjectRestoreOperation  = "restore"
	ObjectDeleteOperation   = "delete"
)

type ObjectOperation struct {
	Verb   string
	Object string
	Date   time.Time
}

type Object struct {
	gorm.Model

	Name        string
	Description string

	LibraryID uint
	Library   Library

	ShelfID uint
	Shelf   Shelf

	History []ObjectOperation `gorm:"-"`
}
