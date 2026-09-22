package modules

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
