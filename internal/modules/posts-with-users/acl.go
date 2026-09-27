package posts_with_users

import (
	acl_constants "github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/acl/constants"
	"github.com/max-fletcher/golang_web_server_boilerplate/middleware"
)

var (
	ACLReadUser = middleware.AclRequirement{
		Module:     acl_constants.EnumModuleUsers,
		Permission: acl_constants.PermissionRead,
	}

	ACLCreateUser = middleware.AclRequirement{
		Module:     acl_constants.EnumModuleUsers,
		Permission: acl_constants.PermissionCreate,
	}

	ACLUpdateUser = middleware.AclRequirement{
		Module:     acl_constants.EnumModuleUsers,
		Permission: acl_constants.PermissionUpdate,
	}

	ACLDeleteUser = middleware.AclRequirement{
		Module:     acl_constants.EnumModuleUsers,
		Permission: acl_constants.PermissionDelete,
	}

	ACLReadPost = middleware.AclRequirement{
		Module:     acl_constants.EnumModulePosts,
		Permission: acl_constants.PermissionRead,
	}

	ACLCreatePost = middleware.AclRequirement{
		Module:     acl_constants.EnumModulePosts,
		Permission: acl_constants.PermissionCreate,
	}

	ACLUpdatePost = middleware.AclRequirement{
		Module:     acl_constants.EnumModulePosts,
		Permission: acl_constants.PermissionUpdate,
	}

	ACLDeletePost = middleware.AclRequirement{
		Module:     acl_constants.EnumModulePosts,
		Permission: acl_constants.PermissionDelete,
	}
)
