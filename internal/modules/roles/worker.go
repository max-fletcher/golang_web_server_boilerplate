package roles

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/max-fletcher/golang_web_server_boilerplate/internal/events"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/queue"
)

type Worker struct {
	// Bind any dependency here when initializing roles worker
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
	var event events.RoleCreated

	if err := json.Unmarshal(message.Data, &event); err != nil {
		return fmt.Errorf("decode role.created event: %w \n", err)
	}

	// Do something with event.ID.
	fmt.Printf("Role created event fired: %v \n", event)

	return nil
}

func (worker *Worker) HandleUpdated(ctx context.Context, message queue.Message) error {
	var event events.RoleUpdated

	if err := json.Unmarshal(message.Data, &event); err != nil {
		return fmt.Errorf("decode role.updated event: %w \n", err)
	}

	// Do something with event.ID.
	fmt.Printf("Role updated event fired: %v \n", event)

	return nil
}

func (worker *Worker) HandleDeleted(ctx context.Context, message queue.Message) error {
	var event events.RoleDeleted

	if err := json.Unmarshal(message.Data, &event); err != nil {
		return fmt.Errorf("decode role.deleted event: %w \n", err)
	}

	// Do something with event.ID.
	fmt.Printf("Role deleted event fired: %v \n", event)

	return nil
}
