package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
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
// 	return strings.Contains(
// 		err.Error(),
// 		"BUSYGROUP",
// 	)
// }

func (worker *Worker) Run(ctx context.Context) error {
	err := worker.createConsumerGroup(ctx)
	if err != nil {
		return err
	}

	worker.logger.Info(
		"queue worker started",
		"stream", worker.stream,
		"consumer_group", worker.consumerGroup,
		"consumer", worker.consumer,
	)

	for {
		select {
		case <-ctx.Done():
			worker.logger.Info("queue worker stopped")
			return ctx.Err()

		default:
		}

		messages, err := worker.client.XReadGroup(
			ctx,
			&redis.XReadGroupArgs{
				Group:    worker.consumerGroup,
				Consumer: worker.consumer,
				Streams:  []string{worker.stream, ">"},
				Count:    10,
				Block:    30 * time.Second,
			},
		).Result()

		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			if err == redis.Nil {
				continue
			}

			worker.logger.Error(
				"failed to read queue",
				"error", err,
			)

			continue
		}

		for _, stream := range messages {
			for _, message := range stream.Messages {
				if err := worker.process(ctx, message); err != nil {
					worker.logger.Error(
						"failed to process queue message",
						"message_id", message.ID,
						"error", err,
					)

					continue
				}

				_, err = worker.client.XAck(
					ctx,
					worker.stream,
					worker.consumerGroup,
					message.ID,
				).Result()

				if err != nil {
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

func (worker *Worker) createConsumerGroup(ctx context.Context) error {
	err := worker.client.XGroupCreateMkStream(
		ctx,
		worker.stream,
		worker.consumerGroup,
		"0",
	).Err()

	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return err
	}

	return nil
}

func (worker *Worker) process(ctx context.Context, message redis.XMessage) error {
	messageType, ok := message.Values["type"].(string)
	if !ok {
		return fmt.Errorf("invalid message type")
	}

	data, ok := message.Values["data"].(string)
	if !ok {
		return fmt.Errorf("invalid message data")
	}

	var payload json.RawMessage

	if err := json.Unmarshal(
		[]byte(data),
		&payload,
	); err != nil {
		return fmt.Errorf("decode message: %w", err)
	}

	return worker.handler(
		ctx,
		Message{
			Type: messageType,
			Data: payload,
		},
	)
}
