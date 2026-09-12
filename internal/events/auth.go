package events

import "github.com/google/uuid"

// Add more events and add more fields per struct as needed when handling these events
type AuthRegistration struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

type AuthLogin struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}
