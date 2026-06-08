package v1

import "time"

type Library struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Name    string `json:"name"`
	Address string `json:"address"`
}
