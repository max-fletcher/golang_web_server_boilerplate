package modules_names

import (
	acl_constants "github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/acl/constants"
	"github.com/max-fletcher/golang_web_server_boilerplate/middleware"
)

var (
	ACLRead = middleware.AclRequirement{
		Module:     acl_constants.EnumModuleModuleNames,
		Permission: acl_constants.PermissionRead,
	}

	ACLCreate = middleware.AclRequirement{
		Module:     acl_constants.EnumModuleModuleNames,
		Permission: acl_constants.PermissionCreate,
	}
)
