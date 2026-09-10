package events

import (
	"context"
	"encoding/json"

	"github.com/max-fletcher/golang_web_server_boilerplate/internal/queue"
)

type Publisher interface {
	Publish(ctx context.Context, eventType EventType, payload any) error
}

type publisher struct {
	queue  queue.Queue
	stream string
}

func NewPublisher(queueClient queue.Queue, stream string) Publisher {
	return &publisher{
		queue:  queueClient,
		stream: stream,
	}
}

func (publisher *publisher) Publish(ctx context.Context, eventType EventType, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return publisher.queue.Publish(
		ctx,
		publisher.stream,
		queue.Message{
			Type: string(eventType),
			Data: data,
		},
	)
}
