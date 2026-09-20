package roles

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
)

type CacheKeyPrefixes string

const (
	CacheKeyRoles CacheKeyPrefixes = "roles:"
	CacheKeyRole  CacheKeyPrefixes = "role:"
)

type GetAllRolesResult struct {
	Roles []db.Role `json:"roles"`
	Total int       `json:"total"`
}

func GetByIDCacheKey(id uuid.UUID) string {
	return string(CacheKeyRole) + id.String()
}

func GetAllCacheKey(filterString string, limit int, offset int) string {
	return fmt.Sprintf("%slist:%s:%d:%d", CacheKeyRoles, filterString, limit, offset)
}
