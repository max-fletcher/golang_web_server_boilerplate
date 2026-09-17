package permissions

type EnumModule string

const (
	ModulePosts       EnumModule = "posts"
	ModuleUsers       EnumModule = "users"
	ModuleRoles       EnumModule = "roles"
	ModulePermissions EnumModule = "permissions"
	ModuleACL         EnumModule = "acl"
)

var ModulesList = []EnumModule{ModulePosts, ModuleUsers, ModuleRoles, ModulePermissions, ModuleACL}

type EnumPermission string

const (
	PermissionCreate EnumPermission = "create"
	PermissionRead   EnumPermission = "read"
	PermissionUpdate EnumPermission = "update"
	PermissionDelete EnumPermission = "delete"
)

var PermissionsList = []EnumPermission{PermissionCreate, PermissionRead, PermissionUpdate, PermissionDelete}
