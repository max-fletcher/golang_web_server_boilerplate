package redis

import (
	"context"

	redisClient "github.com/redis/go-redis/v9"
)

func NewClient(redisURL string) (*redisClient.Client, error) { // create and return new redis client
	options, err := redisClient.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redisClient.NewClient(options)

	err = client.Ping(context.Background()).Err()
	if err != nil { // close client and connection if Ping fails
		_ = client.Close()
		return nil, err
	}

	return client, nil
}
