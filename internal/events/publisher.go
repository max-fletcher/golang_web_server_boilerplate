package events

import (
	"context"
	"encoding/json"

	"github.com/max-fletcher/golang_web_server_boilerplate/internal/queue"
)

// This struct works like a strategy pattern. publisher struct satisfies Publisher interface. It accepts a queue.Queue struct
// and publisher.Publish method just executes the queue.Queue's publish. The difference is the eventType params in
// publisher.Publish is typed EventType, whereas queue.publish's eventType is just a string. If you want to add new EventTypes,
// add them where EventType is declared(i.e internal/events/events.go)

type Publisher interface {
	Publish(ctx context.Context, eventType EventType, payload any) error
}

type publisher struct {
	queue  queue.Queue
	stream string
}

// func NewPublisher() creates a Publisher struct and returns it, so it can be passed down to services where they are used
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
