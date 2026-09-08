package redis

import (
	"context"

	redis "github.com/redis/go-redis/v9"
)

func NewClient(redisURL string) (*redis.Client, error) { // create and return new redis client
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(options)

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
