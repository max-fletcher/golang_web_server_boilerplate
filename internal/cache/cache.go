package cache

import (
	"context"
	"time"
)

type CacheKeyValidPrefixes string

const (
	CacheKeyPosts CacheKeyValidPrefixes = "posts:"
)

type Cache interface {
	Get(ctx context.Context, key string, dest any) error
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	DeleteByPrefix(ctx context.Context, prefix CacheKeyValidPrefixes) error
}
