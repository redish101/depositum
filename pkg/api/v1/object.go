package v1

import "time"

type ObjectOperation struct {
	Verb   string
	Object string
	Date   time.Time
}

type ObjectPhase string

const (
	ObjectPhaseCreated   ObjectPhase = "CREATED"
	ObjectPhaseArchived  ObjectPhase = "ARCHIVED"
	ObjectPhaseWithdrawn ObjectPhase = "WITHDRAWN"
	ObjectPhaseDeleted   ObjectPhase = "DELETED"
	ObjectPhaseLost      ObjectPhase = "LOST"
)

type ObjectStatus struct {
	Phase   ObjectPhase    `json:"phase"`
	Library LibrarySummary `json:"library"`
	Shelf   ShelfSummary   `json:"shelf"`
}

type Object struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Name        string `json:"name"`
	Description string `json:"description"`

	Synced        bool         `json:"synced"`
	CurrentStatus ObjectStatus `json:"currentStatus"`
	DesiredStatus ObjectStatus `json:"desiredStatus"`
}

type ListObjectRequest struct {
	SyncedOnly   bool  `query:"syncedOnly" validate:"omitempty"`
	UnsyncedOnly bool  `query:"unsyncedOnly" validate:"omitempty"`
	LibraryID    *uint `query:"libraryID" validate:"omitempty"`
	ShelfID      *uint `query:"shelfID" validate:"omitempty"`
}

type ObjectStatusInput struct {
	Phase     ObjectPhase `json:"phase"`
	LibraryID uint        `json:"libraryID"`
	ShelfID   uint        `json:"shelfID"`
}

type CreateObjectRequest struct {
	Name          string            `json:"name" validate:"required"`
	Description   string            `json:"description" validate:"omitempty"`
	DesiredStatus ObjectStatusInput `json:"desiredStatus" validate:"omitempty"`
	SyncNow       bool              `json:"syncnow" validate:"required"`
}

type UpdateObjectRequest struct {
	Name          *string            `json:"name" validate:"omitempty"`
	Description   *string            `json:"description" validate:"omitempty"`
	CurrentStatus *ObjectStatusInput `json:"currentStatus" validate:"omitempty"`
	DesiredStatus *ObjectStatusInput `json:"desiredStatus" validate:"omitempty"`
	SyncNow       *bool              `json:"syncnow" validate:"omitempty"`
}
