package posts

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/max-fletcher/golang_web_server_boilerplate/internal/events"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/queue"
)

type Worker struct {
	// Bind any dependency here when initializing posts worker
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
	var event events.PostCreated

	if err := json.Unmarshal(message.Data, &event); err != nil {
		return fmt.Errorf("decode post.created event: %w", err)
	}

	// Do something with event.ID.
	fmt.Printf("Post created event fired: %v", event)

	return nil
}

func (worker *Worker) HandleUpdated(ctx context.Context, message queue.Message) error {
	var event events.PostUpdated

	if err := json.Unmarshal(message.Data, &event); err != nil {
		return fmt.Errorf("decode post.updated event: %w", err)
	}

	// Do something with event.ID.
	fmt.Printf("Post updated event fired: %v", event)

	return nil
}

func (worker *Worker) HandleDeleted(ctx context.Context, message queue.Message) error {
	var event events.PostDeleted

	if err := json.Unmarshal(message.Data, &event); err != nil {
		return fmt.Errorf("decode post.deleted event: %w", err)
	}

	// Do something with event.ID.
	fmt.Printf("Post deleted event fired: %v", event)

	return nil
}
