package events

import "github.com/google/uuid"

// Add more events and add more fields per struct as needed when handling these events
type PostCreated struct {
	ID uuid.UUID `json:"id"`
}

type PostUpdated struct {
	ID uuid.UUID `json:"id"`
}

type PostDeleted struct {
	ID uuid.UUID `json:"id"`
}
