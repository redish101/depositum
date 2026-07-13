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
	ObjectPhaseLost      = "lost"
)

type ObjectStatus struct {
	Phase     string `json:"phase"`
	LibraryID uint   `json:"libraryID"`
	ShelfID   uint   `json:"shelfID"`
}

type Object struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Name        string `json:"name"`
	Description string `json:"description"`

	Synced        bool `json:"synced"`
	CurrentStatus ObjectStatus `json:"currentStatus"`
	DesiredStatus ObjectStatus `json:"desiredStatus"`
}

type ListObjectRequest struct {
	SyncedOnly bool `json:"syncedOnly" validate:"omitempty"`
	UnsyncedOnly bool `json:"unsyncedOnly" validate:"omitempty"`
	LibraryID *uint `json:"libraryID" validate:"omitempty"`
	ShelfID   *uint `json:"shelfID" validate:"omitempty"`
}

type CreateObjectRequest struct {
	Name string `json:"name" validate:"required"`
	Description string `json:"description" validate:"omitempty"`
	DesiredStatus ObjectStatus `json:"desiredStatus" validate:"omitempty"`
	SyncNow bool `json:"syncnow" validate:"required"`
}

type UpdateObjectRequest struct {
	Name *string `json:"name" validate:"omitempty"`
	Description *string `json:"description" validate:"omitempty"`
	CurrentStatus *ObjectStatus `json:"currentStatus" validate:"omitempty"`
	DesiredStatus *ObjectStatus `json:"desiredStatus" validate:"omitempty"`
}

