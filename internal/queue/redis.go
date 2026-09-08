package queue

import (
	"context"
	"encoding/json"

	redis "github.com/redis/go-redis/v9"
)

// Redis Implementation
type redisQueue struct {
	client *redis.Client
}

func NewRedisQueue(client *redis.Client) Queue {
	return &redisQueue{
		client: client,
	}
}

func (queue *redisQueue) Publish(ctx context.Context, stream StreamType, message Message) error {
	data, err := json.Marshal(message.Data)
	if err != nil {
		return err
	}

	_, err = queue.client.XAdd(ctx, &redis.XAddArgs{ // publish messages/events
		Stream: string(stream), // Works like channels/topics in kafka. This needs to be converted to string for it to work as arg
		Values: map[string]any{
			"type": string(message.Type), // This needs to be converted to string for it to work as arg
			"data": string(data),
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
