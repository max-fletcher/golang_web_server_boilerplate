package redis_queue

import (
	"context"

	"github.com/max-fletcher/golang_web_server_boilerplate/internal/queue"
	redisClient "github.com/redis/go-redis/v9"
)

// Redis Implementation
type redisQueue struct {
	client *redisClient.Client
}

func NewRedisQueue(client *redisClient.Client) queue.Queue {
	return &redisQueue{
		client: client,
	}
}

func (queue *redisQueue) Publish(ctx context.Context, stream string, message queue.Message) error {
	_, err := queue.client.XAdd(ctx, &redisClient.XAddArgs{ // publish messages/events
		Stream: stream, // Works like channels/topics in kafka. This needs to be converted to string for it to work as arg
		Values: map[string]any{
			"type": message.Type,
			"data": string(message.Data),
		},
	}).Result()

	// If you don't want to without streams(#2)
	// _, err = queue.client.XAdd(ctx, &redis.XAddArgs{
	// 	Values: map[string]any{
	// 		"type": message.Type,
	// 		"data": string(data),
	// 	},
	// }).Result()

	return err
}
