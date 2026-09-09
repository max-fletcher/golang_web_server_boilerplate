package posts_with_users

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
)

type CacheKeyPrefixes string

const (
	CacheKeyPostsWithUser CacheKeyPrefixes = "postsWithUser:"
	CacheKeyPostWithUser  CacheKeyPrefixes = "postWithUser:"
)

type GetAllPostsWithUsersResult struct {
	Posts []db.GetPostsWithUserRow `json:"posts"`
	Total int                      `json:"total"`
}

func GetByIDCacheKey(id uuid.UUID) string {
	return string(CacheKeyPostWithUser) + id.String()
}

func GetAllCacheKey(filterString string, limit int, offset int) string {
	return fmt.Sprintf("%slist:%s:%d:%d", CacheKeyPostWithUser, filterString, limit, offset)
}
