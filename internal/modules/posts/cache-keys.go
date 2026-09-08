package posts

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
)

type CacheKeyPrefixes string

const (
	CacheKeyPosts CacheKeyPrefixes = "posts:"
	CacheKeyPost  CacheKeyPrefixes = "post:"
)

type GetAllPostsResult struct {
	Posts []db.Post `json:"posts"`
	Total int       `json:"total"`
}

func GetByIDCacheKey(id uuid.UUID) string {
	return string(CacheKeyPost) + id.String()
}

func GetAllCacheKey(filterString string, limit int, offset int) string {
	return fmt.Sprintf("%slist:%s:%d:%d", CacheKeyPosts, filterString, limit, offset)
}
