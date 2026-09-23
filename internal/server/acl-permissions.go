package server

import (
	"net/http"

	modules "github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/acl/module-names"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/acl/permissions"
	"github.com/max-fletcher/golang_web_server_boilerplate/middleware"
)

// func (middleware *middleware.ACLMiddleware) RequirePermission(requirements ...middleware.AclRequirement) func(http.Handler) http.Handler
// func (middleware *ACLMiddleware) RequirePermission(requirements ...AclRequirement) func(http.Handler) http.Handler

func (server *Server) GetPostsACL() func(http.Handler) http.Handler {
	requirements := []middleware.AclRequirement{
		{
			Module:     modules.EnumModulePosts,
			Permission: permissions.PermissionRead,
		},
	}

	return server.ACLMiddleware.RequirePermission(requirements...) // IMPORTANT: works like destructuring when you pass a slice to a variadic fn
}
