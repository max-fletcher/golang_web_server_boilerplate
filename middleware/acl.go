package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/responses"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
	modules "github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/acl/module-names"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/acl/permissions"
)

type AclRequirement struct {
	Module     modules.EnumModuleNames
	Permission permissions.EnumPermission
}

type RolePermissionsService interface { // For DI. Used in validating user in structs.go
	UserHasPermission(ctx context.Context, moduleName modules.EnumModuleNames, permissionName db.PermissionNamesEnum) (bool, error)
}

type ACLMiddleware struct {
	rolePermissionService RolePermissionsService
}

func NewACLMiddleware(rolePermissionService RolePermissionsService) *ACLMiddleware {
	return &ACLMiddleware{
		rolePermissionService: rolePermissionService,
	}
}

func (middleware *ACLMiddleware) RequirePermission(requirements ...AclRequirement) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, requirement := range requirements {
				allowed, err := middleware.rolePermissionService.UserHasPermission(
					r.Context(),
					requirement.Module,
					db.PermissionNamesEnum(requirement.Permission),
				)
				if err != nil {
					fmt.Println("++++ ACL middleware error: ++++\n", err)
					responses.InternalServerErrorSWW(w)
					return
				}

				if !allowed {
					responses.RespondWithJSON(w, http.StatusNotFound, responses.ErrorResponse{
						Code:    http.StatusForbidden,
						Status:  "error",
						Message: "You don't have the necessary permissions to access this route",
					})
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
