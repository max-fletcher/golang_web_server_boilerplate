package users

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/max-fletcher/golang_web_server_boilerplate/internal/events"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/queue"
)

type Worker struct {
	// Bind any dependency here when initializing users worker
	// cache cache.Cache
}

func NewWorker(
// cacheClient cache.Cache
) *Worker {
	return &Worker{
		// cache: cacheClient,
	}
}

func (worker *Worker) HandleCreated(ctx context.Context, message queue.Message) error {
	var event events.UserCreated

	if err := json.Unmarshal(
		message.Data,
		&event,
	); err != nil {
		return fmt.Errorf("decode user.created event: %w \n", err)
	}

	// Do something with event.ID.
	fmt.Printf("User created event fired: %v \n", event)

	return nil
}

func (worker *Worker) HandleUpdated(ctx context.Context, message queue.Message) error {
	var event events.UserUpdated

	if err := json.Unmarshal(
		message.Data,
		&event,
	); err != nil {
		return fmt.Errorf("decode user.updated event: %w \n", err)
	}

	// Do something with event.ID.
	fmt.Printf("User updated event fired: %v \n", event)

	return nil
}

func (worker *Worker) HandleDeleted(ctx context.Context, message queue.Message) error {
	var event events.UserDeleted

	if err := json.Unmarshal(
		message.Data,
		&event,
	); err != nil {
		return fmt.Errorf("decode user.deleted event: %w \n", err)
	}

	// Do something with event.ID.
	fmt.Printf("User deleted event fired: %v \n", event)

	return nil
}
