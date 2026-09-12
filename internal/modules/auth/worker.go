package auth

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

func (worker *Worker) HandleRegistration(ctx context.Context, message queue.Message) error {
	var event events.UserCreated

	if err := json.Unmarshal(
		message.Data,
		&event,
	); err != nil {
		return fmt.Errorf("decode auth.registration event: %w \n", err)
	}

	// Do something with event.ID.
	fmt.Printf("User registration event fired: %v \n", event)

	return nil
}

func (worker *Worker) HandleLogin(ctx context.Context, message queue.Message) error {
	var event events.UserUpdated

	if err := json.Unmarshal(
		message.Data,
		&event,
	); err != nil {
		return fmt.Errorf("decode auth.login event: %w \n", err)
	}

	// Do something with event.ID.
	fmt.Printf("User login event fired: %v \n", event)

	return nil
}
