package posts

import (
	modules "github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/acl/module-names"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/acl/permissions"
	"github.com/max-fletcher/golang_web_server_boilerplate/middleware"
)

var (
	ACLRead = middleware.AclRequirement{
		Module:     modules.EnumModulePosts,
		Permission: permissions.PermissionRead,
	}

	ACLCreate = middleware.AclRequirement{
		Module:     modules.EnumModulePosts,
		Permission: permissions.PermissionCreate,
	}

	ACLUpdate = middleware.AclRequirement{
		Module:     modules.EnumModulePosts,
		Permission: permissions.PermissionUpdate,
	}

	ACLDelete = middleware.AclRequirement{
		Module:     modules.EnumModulePosts,
		Permission: permissions.PermissionDelete,
	}
)
