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

func (worker *Worker) Run(ctx context.Context) error {
	err := worker.client.XGroupCreateMkStream(
		ctx,
		worker.stream,
		worker.consumerGroup,
		"0",
	).Err()

	if err != nil &&
		!isConsumerGroupExistsError(err) {
		return err
	}

	for {
		select {
		case <-ctx.Done():
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
				Block:    5 * time.Second,
			},
		).Result()

		if err != nil {
			if errors.Is(err, redis.Nil) {
				continue
			}

			return err
		}

		for _, stream := range messages {
			for _, message := range stream.Messages {
				if err := worker.process(ctx, message); err != nil {
					worker.logger.Error(
						"queue message failed",
						"stream", worker.stream,
						"message_id", message.ID,
						"error", err,
					)

					continue
				}

				if err := worker.client.XAck(
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
