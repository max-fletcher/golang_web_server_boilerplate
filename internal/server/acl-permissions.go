package server

import (
	"net/http"

	acl_constants "github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/acl/constants"
	"github.com/max-fletcher/golang_web_server_boilerplate/middleware"
)

// func (middleware *middleware.ACLMiddleware) RequirePermission(requirements ...middleware.AclRequirement) func(http.Handler) http.Handler
// func (middleware *ACLMiddleware) RequirePermission(requirements ...AclRequirement) func(http.Handler) http.Handler

func (server *Server) GetPostsACL() func(http.Handler) http.Handler {
	requirements := []middleware.AclRequirement{
		{
			Module:     acl_constants.EnumModulePosts,
			Permission: acl_constants.PermissionRead,
		},
	}

	return server.ACLMiddleware.RequirePermission(requirements...) // IMPORTANT: works like destructuring when you pass a slice to a variadic fn
}
