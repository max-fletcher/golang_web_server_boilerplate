package events

type StreamType string

const DefaultStream = "app:events"

// use this below instead of the line above if you ever needed more streams
// const (
// 	StreamPost StreamType = "app:post"
// 	StreamUser StreamType = "app:user"
// )

// *IMPORTANT: If you want to add new events, do the following:
//  1. Add a new EventType here
//  2. Add a new event struct inside internal/events/{your_module_name}.go OR make new file and add it there
//  3. Create new handler against the new event you created above inside internal/modules/{your_module}/worker.go and define any logic
//     that should be executed if the event is processed by the queue here in the handler
//  4. Register your EventType against the handler you created inside /internal/workers/workers.go
type EventType string

const (
	QueueEventPostCreated EventType = "post.created"
	QueueEventPostUpdated EventType = "post.updated"
	QueueEventPostDeleted EventType = "post.deleted"

	QueueEventUserCreated EventType = "user.created"
	QueueEventUserUpdated EventType = "user.updated"
	QueueEventUserDeleted EventType = "user.deleted"
)
