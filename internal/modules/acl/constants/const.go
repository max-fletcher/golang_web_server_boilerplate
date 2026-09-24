package acl_constants

type EnumModuleNames string

const (
	EnumModuleACL             EnumModuleNames = "acl"
	EnumModuleUsers           EnumModuleNames = "users"
	EnumModulePosts           EnumModuleNames = "posts"
	EnumModuleRoles           EnumModuleNames = "roles"
	EnumModuleModules         EnumModuleNames = "modules"
	EnumModulePermissions     EnumModuleNames = "permissions"
	EnumModuleUserRoles       EnumModuleNames = "user-roles"
	EnumModuleRolePermissions EnumModuleNames = "role-permissions"
)

var ModulesList = []EnumModuleNames{
	EnumModuleACL,
	EnumModuleUsers,
	EnumModulePosts,
	EnumModuleRoles,
	EnumModuleModules,
	EnumModulePermissions,
	EnumModuleUserRoles,
	EnumModuleRolePermissions,
}

type EnumPermission string

const (
	PermissionCreate EnumPermission = "create"
	PermissionRead   EnumPermission = "read"
	PermissionUpdate EnumPermission = "update"
	PermissionDelete EnumPermission = "delete"
)

var PermissionsList = []EnumPermission{PermissionCreate, PermissionRead, PermissionUpdate, PermissionDelete}
