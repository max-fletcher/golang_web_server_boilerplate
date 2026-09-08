package queue

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"time"

	redis "github.com/redis/go-redis/v9"
)

// Worker that will handle redis queue messages

type Handler func(ctx context.Context, message Message) error

type Worker struct {
	client        *redis.Client
	stream        string
	consumerGroup string
	consumer      string
	handler       Handler
	logger        *slog.Logger
}

// creates and returns a new worker struct. We are passing:
// a) a handler(a func that can do what you want) as an arg so it is bound to this worker struct and can be used as long as we pass the client to it
// b) other params and methods that we can use e.g streams are like channels/topics
// c) Run is to run the queue worker and process takes messages and runs the handler against them
func NewWorker(
	client *redis.Client,
	stream string,
	consumerGroup string,
	consumer string,
	handler Handler,
	logger *slog.Logger,
) *Worker {
	return &Worker{
		client:        client,
		stream:        stream,
		consumerGroup: consumerGroup,
		consumer:      consumer,
		handler:       handler,
		logger:        logger,
	}
}

// *IMPORTANT
// Current Queue flow:
// XReadGroup()
//  ↓
// process()
//  ↓
// SUCCESS
//  ↓
// XAck()

func (worker *Worker) Run(ctx context.Context) error {
	err := worker.client.XGroupCreateMkStream( // When the worker starts, it needs to make sure the group exists
		ctx,
		worker.stream,
		worker.consumerGroup,
		"0",
	).Err()

	if err != nil &&
		!isConsumerGroupExistsError(err) {
		return err
	}

	worker.logger.Info( // check if queue started or not
		"queue worker started",
		"stream", worker.stream,
		"consumer_group", worker.consumerGroup,
		"consumer", worker.consumer,
	)

	for {
		messages, err := worker.client.XReadGroup( // the worker starts reading
			ctx,
			&redis.XReadGroupArgs{
				Group:    worker.consumerGroup,         // Which consumer group is reading
				Consumer: worker.consumer,              // Which worker within the group is reading
				Streams:  []string{worker.stream, ">"}, // input messages that have never been delivered to another consumer in this group
				Count:    10,
				Block:    5 * time.Second, // makes sure the worker doesn't constantly hammer Redis
			},
		).Result()

		if err != nil {
			if errors.Is(err, redis.Nil) {
				continue
			}

			if ctx.Err() != nil {
				return ctx.Err()
			}

			return err
		}

		worker.logger.Info( // check if XReadGroup received the message
			"queue messages received",
			"count", len(messages),
		)

		for _, stream := range messages {
			for _, message := range stream.Messages {

				worker.logger.Info( // for checking if messages are coming upto here before being processes
					"processing queue message",
					"message_id", message.ID,
					"type", message.Values["type"],
				)

				if err := worker.process(ctx, message); err != nil { // if processing a message failed, throw err
					worker.logger.Error(
						"queue message failed",
						"stream", worker.stream,
						"message_id", message.ID,
						"error", err,
					)

					continue
				}

				if err := worker.client.XAck( // runs after successfully processing message(i.e acknowledged)
					ctx,
					worker.stream,
					worker.consumerGroup,
					message.ID,
				).Err(); err != nil {
					worker.logger.Error(
						"failed to acknowledge queue message",
						"message_id", message.ID,
						"error", err,
					)
				}
			}
		}
	}
}

func (worker *Worker) process(
	ctx context.Context,
	message redis.XMessage,
) error {
	messageType, ok := message.Values["type"].(string)
	if !ok {
		return errors.New("queue message type missing")
	}

	data, ok := message.Values["data"].(string)
	if !ok {
		return errors.New("queue message data missing")
	}

	var decoded any

	if err := json.Unmarshal([]byte(data), &decoded); err != nil {
		return err
	}

	return worker.handler(
		ctx,
		Message{
			Type: MessageType(messageType), // Casting.converting it to "MessageType" type(see internal/queue/queue.go)
			Data: decoded,
		},
	)
}

func isConsumerGroupExistsError(err error) bool {
	return strings.Contains(
		err.Error(),
		"BUSYGROUP",
	)
}
