package users

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
)

type CacheKeyPrefixes string

const (
	CacheKeyUsers CacheKeyPrefixes = "users:"
	CacheKeyUser  CacheKeyPrefixes = "user:"
)

type GetAllUsersResult struct {
	Users []db.User `json:"users"`
	Total int       `json:"total"`
}

func GetByIDCacheKey(id uuid.UUID) string {
	return string(CacheKeyUser) + id.String()
}

func GetAllCacheKey(filterString string, limit int, offset int) string {
	return fmt.Sprintf("%slist:%s:%d:%d", CacheKeyUsers, filterString, limit, offset)
}
