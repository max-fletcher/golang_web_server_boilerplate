package queue

import (
	"context"
)

// Abstraction -> doesn't know that Redis is being used

type Message struct {
	Type string `json:"type"` // Type now has to satisfy "EventType". Works like enum
	Data []byte `json:"data"`
}

type Queue interface {
	Publish(ctx context.Context, stream string, message Message) error

	// If you don't want to without streams(#1)
	// Publish(ctx context.Context, message Message) error
}
