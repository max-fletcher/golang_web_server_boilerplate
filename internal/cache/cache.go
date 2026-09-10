package cache

import (
	"context"
	"time"
)

type CacheKeyValidPrefixes string

const (
	CacheKeyPosts CacheKeyValidPrefixes = "posts:"
	CacheKeyUsers CacheKeyValidPrefixes = "users:"
)

type Cache interface {
	Get(ctx context.Context, key string, dest any) error
	Set(ctx context.Context, key string, value any, expiration ...time.Duration) error // making expiration optional by making it variadic
	Delete(ctx context.Context, key string) error
	DeleteByPrefix(ctx context.Context, prefix CacheKeyValidPrefixes) error
}
