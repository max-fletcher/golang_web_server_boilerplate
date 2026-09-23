package permissions

type EnumPermission string

const (
	PermissionCreate EnumPermission = "create"
	PermissionRead   EnumPermission = "read"
	PermissionUpdate EnumPermission = "update"
	PermissionDelete EnumPermission = "delete"
)

var PermissionsList = []EnumPermission{PermissionCreate, PermissionRead, PermissionUpdate, PermissionDelete}
