package events

import "github.com/google/uuid"

// Add more events and add more fields per struct as needed when handling these events
type UserCreated struct {
	ID uuid.UUID `json:"id"`
}

type UserUpdated struct {
	ID uuid.UUID `json:"id"`
}

type UserDeleted struct {
	ID uuid.UUID `json:"id"`
}
