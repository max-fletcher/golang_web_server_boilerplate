package queue

import "context"

// Abstraction -> doesn't know that Redis is being used

type StreamType string

const (
	StreamPost StreamType = "app:post"
	StreamUser StreamType = "app:user"
)

type MessageType string

const (
	QueueMsgPostCreated MessageType = "post.created"
	QueueMsgPostUpdated MessageType = "post.updated"
	QueueMsgPostDeleted MessageType = "post.deleted"

	QueueMsgUserCreated MessageType = "user.created"
	QueueMsgUserUpdated MessageType = "user.updated"
	QueueMsgUserDeleted MessageType = "user.deleted"
)

type Message struct {
	Type MessageType `json:"type"` // Type now has to satisfy "MessageType" above. Works like enum
	Data any         `json:"data"`
}

type Queue interface {
	Publish(ctx context.Context, stream StreamType, message Message) error

	// If you don't want to without streams(#1)
	// Publish(ctx context.Context, message Message) error
}
