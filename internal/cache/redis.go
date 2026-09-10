package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	redis "github.com/redis/go-redis/v9"
)

// Redis Implementaiton

type redisCache struct {
	client     *redis.Client
	isActive   bool
	defaultExp time.Duration
}

func NewRedisCache(client *redis.Client, isActive bool, defaultExp time.Duration) Cache { // Return a struct that contains an instance of redis client
	return &redisCache{
		client:     client,
		isActive:   isActive,
		defaultExp: defaultExp,
	}
}

// Get value from cache if exists, else, returns an error
// "dest" contains a pointer to the variable you want to store data in
func (cache *redisCache) Get(ctx context.Context, key string, dest any) error {
	if !cache.isActive {
		return nil
	}
	value, err := cache.client.Get(ctx, key).Result() // "value" will contains stringified json data
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrCacheMiss{
				Err: err,
			}
		}

		return err
	}

	// Unmarshal/decode "value"(stringified json data) to dest(struct). If struct and data doesn't match, those fields are safely ignored
	if err := json.Unmarshal([]byte(value), dest); err != nil {
		return err
	}

	return nil
}

// Set value to cache
func (cache *redisCache) Set(ctx context.Context, key string, value any, expiration ...time.Duration) error {
	if !cache.isActive {
		return nil
	}
	data, err := json.Marshal(value) // marshal/encode struct to stringified json
	if err != nil {
		return err
	}

	if len(expiration) > 1 {
		return fmt.Errorf("cache.Set accepts at most one expiration")
	}

	exp := cache.defaultExp
	if len(expiration) == 1 {
		exp = expiration[0]
	}

	return cache.client.Set(ctx, key, data, exp).Err()
}

func (cache *redisCache) Delete(ctx context.Context, key string) error {
	if !cache.isActive {
		return nil
	}
	return cache.client.Del(ctx, key).Err()
}

func (cache *redisCache) DeleteByPrefix(ctx context.Context, prefix CacheKeyValidPrefixes) error {
	if !cache.isActive {
		return nil
	}

	var cursor uint64
	for {
		keys, nextCursor, err := cache.client.Scan( // Scan for keys that contains "prefix*"(prefix, then any trailing chars)
			ctx,
			cursor,
			string(prefix)+"*",
			0,
		).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 { // If keys exist, Unlink i.e delete cache with those keys
			if err := cache.client.Unlink(ctx, keys...).Err(); err != nil {
				return err
			}
		}

		cursor = nextCursor // Scan uses cursor to loop(iterates incrementally) so updates cursor

		if cursor == 0 {
			break
		}
	}

	return nil
}
