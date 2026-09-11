package redis_queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/max-fletcher/golang_web_server_boilerplate/internal/queue"
	"github.com/redis/go-redis/v9"
	redisClient "github.com/redis/go-redis/v9"
)

// Worker that will handle redis queue messages

type Worker struct {
	client        *redisClient.Client
	stream        string
	consumerGroup string
	consumer      string
	handler       queue.Handler
	logger        *slog.Logger
}

// creates and returns a new worker struct. We are passing:
// a) a handler(a func that can do what you want) as an arg so it is bound to this worker struct and can be used as long as we pass the client to it
// b) other params and methods that we can use e.g streams are like channels/topics
// c) Run is to run the queue worker and process takes messages and runs the handler against them
// Below: Consumer group = "Which team of workers am I part of?" and Consumer = "Which individual worker am I?". If your application grows
// and you need multiple workers to keep up with handling events(either across different applications or goroutines), you can define say
// consumer = "worker-2" with same consumerGroup name(since we need the new worker to take events from same consumerGroups) and use it so more
// workers can consume events and process them faster
func NewWorker(
	client *redisClient.Client,
	stream string,
	consumerGroup string, // Consumer group = "Which team of workers am I part of?"
	consumer string, // Consumer = "Which individual worker am I?"
	handler queue.Handler,
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
		messages, err := worker.client.XReadGroup(
			ctx,
			&redisClient.XReadGroupArgs{
				Group:    worker.consumerGroup, // Means: Read as part of this consumer group
				Consumer: worker.consumer,      // Means: This particular worker is requesting messages
				Streams:  []string{worker.stream, ">"},
				Count:    10, // Means: this worker can ask Redis for up to 10 messages at once
				// Means: If there are currently no messages, wait for up to 5 seconds for one to arrive
				// Related to the timeout issue faced earlier. If Redis client's network ReadTimeout is shorter than the Redis BLOCK duration, you can get
				// "i/o timeout" even though Redis itself is functioning.
				Block: 30 * time.Second,
			},
		).Result()

		if err != nil {
			if ctx.Err() != nil {
				worker.logger.Info("queue worker stopped")
				return ctx.Err()
			}

			// Retry if when key does not exist. Else, log error below(worker.logger.Error)
			if err == redisClient.Nil {
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

				// signals to redis that the message was processed successfully causing it to be marked as complete
				// and no retries/reissue to other worker happens
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
		queue.Message{
			Type: messageType,
			Data: payload,
		},
	)
}
