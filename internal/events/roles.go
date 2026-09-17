package events

import "github.com/google/uuid"

// Add more events and add more fields per struct as needed when handling these events
type RoleCreated struct {
	ID uuid.UUID `json:"id"`
}

type RoleUpdated struct {
	ID uuid.UUID `json:"id"`
}

type RoleDeleted struct {
	ID uuid.UUID `json:"id"`
}
