package v1

import "time"

type ObjectOperation struct {
	Verb   string
	Object string
	Date   time.Time
}

const (
	ObjectPhaseCreated   = "created"
	ObjectPhaseArchived  = "archived"
	ObjectPhaseWithdrawn = "withdrawn"
	ObjectPhaseDeleted   = "deleted"
)

type ObjectStatus struct {
	Phase     string `json:"phase"`
	LibraryID uint   `json:"libraryID"`
	ShelfID   uint   `json:"shelfID"`
}

type Object struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name        string
	Description string

	Synced        bool
	CurrentStatus ObjectStatus `json:"current_status"`
	DesiredStatus ObjectStatus `json:"desired_status"`
}
