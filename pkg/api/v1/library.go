package v1

import (
	"time"
)

type Library struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Name    string `json:"name"`
	Address string `json:"address"`
}

type CreateLibraryRequest struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
}

type UpdateLibraryRequest struct {
	Name    *string `json:"name" validate:"omitempty"`
	Address *string `json:"address" validate:"omitempty"`
}

// type DeleteLibraryResponse struct {
// 	DeletedAt time.Time `json:"deletedAt"`
// }
