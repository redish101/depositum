package v1

type Shelf struct {
	ID uint `json:"id"`

	Name        string `json:"name"`
	Description string `json:"description"`
	LibraryID   uint   `json:"libraryID"`
}

type CreateShelfRequest struct {
}
