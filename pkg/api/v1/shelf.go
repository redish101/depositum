package v1

type Shelf struct {
	ID uint `json:"id"`

	Name        string `json:"name"`
	Description string `json:"description"`
	LibraryID   uint   `json:"libraryID"`
}

type CreateShelfRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	LibraryID   uint   `json:"libraryID" validate:"required"`
}

type UpdateShelfRequest struct {
	Name        *string `json:"name" validate:"omitempty"`
	Description *string `json:"description" validate:"omitempty"`
	LibraryID   *uint   `json:"libraryID" validate:"omitempty"`
}
