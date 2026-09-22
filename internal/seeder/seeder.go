package seeder

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/cryptography"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
	modules "github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/acl/module-names"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/acl/permissions"
	role_permissions "github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/acl/role-permissions"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/acl/roles"
	user_roles "github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/acl/user-roles"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/users"
)

type EnumRoles string

const (
	EnumRolesSuperadmin EnumRoles = "superadmin"
	EnumRolesAdmin      EnumRoles = "admin"
	EnumRolesUser       EnumRoles = "user"
)

func Modules(ctx context.Context, database *db.Queries) error {
	moduleNames := []modules.EnumModuleNames{
		modules.EnumModuleACL,
		modules.EnumModuleUsers,
		modules.EnumModulePosts,
		modules.EnumModuleRoles,
		modules.EnumModuleModules,
		modules.EnumModulePermissions,
		modules.EnumModuleUserRoles,
		modules.EnumModuleRolePermissions,
	}
	moduleRepository := modules.NewRepository(database)

	for _, moduleName := range moduleNames {
		createModuleParams := db.CreateModuleParams{
			ID:        uuid.New(),
			Name:      string(moduleName),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
		_, err := moduleRepository.Create(ctx, createModuleParams)
		if err != nil {
			return err
		}
	}

	return nil
}

func Roles(ctx context.Context, database *db.Queries) error {
	roleNames := []string{"superadmin", "admin", "user"}
	roleRepository := roles.NewRepository(database)

	for _, roleName := range roleNames {
		createRoleParams := db.CreateRoleParams{
			ID:        uuid.New(),
			Name:      roleName,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
		_, err := roleRepository.Create(ctx, createRoleParams)
		if err != nil {
			return err
		}
	}

	return nil
}

func Permissions(ctx context.Context, database *db.Queries) error {
	permissionNames := []db.PermissionNamesEnum{
		db.PermissionNamesEnumCreate,
		db.PermissionNamesEnumRead,
		db.PermissionNamesEnumUpdate,
		db.PermissionNamesEnumDelete,
	}
	moduleRepository := modules.NewRepository(database)
	allModules, err := moduleRepository.GetAll(ctx, "", 10000, 0)
	if err != nil {
		return err
	}

	permissionRepository := permissions.NewRepository(database)
	for _, module := range allModules {
		for _, permissionName := range permissionNames {
			createpermissionParams := db.CreatePermissionParams{
				ID:        uuid.New(),
				Name:      permissionName,
				ModuleID:  module.ID,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			}
			_, err := permissionRepository.Create(ctx, createpermissionParams)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func RolePermissions(ctx context.Context, database *db.Queries) error {
	privileges := [][]string{
		// SUPERADMIN PRIVILEGES
		{string(EnumRolesSuperadmin), string(modules.EnumModuleACL), string(db.PermissionNamesEnumCreate)},
		{string(EnumRolesSuperadmin), string(modules.EnumModuleACL), string(db.PermissionNamesEnumRead)},
		{string(EnumRolesSuperadmin), string(modules.EnumModuleACL), string(db.PermissionNamesEnumUpdate)},
		{string(EnumRolesSuperadmin), string(modules.EnumModuleACL), string(db.PermissionNamesEnumDelete)},

		{string(EnumRolesSuperadmin), string(modules.EnumModuleUsers), string(db.PermissionNamesEnumCreate)},
		{string(EnumRolesSuperadmin), string(modules.EnumModuleUsers), string(db.PermissionNamesEnumRead)},
		{string(EnumRolesSuperadmin), string(modules.EnumModuleUsers), string(db.PermissionNamesEnumUpdate)},
		{string(EnumRolesSuperadmin), string(modules.EnumModuleUsers), string(db.PermissionNamesEnumDelete)},

		{string(EnumRolesSuperadmin), string(modules.EnumModulePosts), string(db.PermissionNamesEnumCreate)},
		{string(EnumRolesSuperadmin), string(modules.EnumModulePosts), string(db.PermissionNamesEnumRead)},
		{string(EnumRolesSuperadmin), string(modules.EnumModulePosts), string(db.PermissionNamesEnumUpdate)},
		{string(EnumRolesSuperadmin), string(modules.EnumModulePosts), string(db.PermissionNamesEnumDelete)},

		{string(EnumRolesSuperadmin), string(modules.EnumModuleRoles), string(db.PermissionNamesEnumCreate)},
		{string(EnumRolesSuperadmin), string(modules.EnumModuleRoles), string(db.PermissionNamesEnumRead)},
		{string(EnumRolesSuperadmin), string(modules.EnumModuleRoles), string(db.PermissionNamesEnumUpdate)},
		{string(EnumRolesSuperadmin), string(modules.EnumModuleRoles), string(db.PermissionNamesEnumDelete)},

		{string(EnumRolesSuperadmin), string(modules.EnumModuleModules), string(db.PermissionNamesEnumCreate)},
		{string(EnumRolesSuperadmin), string(modules.EnumModuleModules), string(db.PermissionNamesEnumRead)},
		{string(EnumRolesSuperadmin), string(modules.EnumModuleModules), string(db.PermissionNamesEnumUpdate)},
		{string(EnumRolesSuperadmin), string(modules.EnumModuleModules), string(db.PermissionNamesEnumDelete)},

		{string(EnumRolesSuperadmin), string(modules.EnumModulePermissions), string(db.PermissionNamesEnumCreate)},
		{string(EnumRolesSuperadmin), string(modules.EnumModulePermissions), string(db.PermissionNamesEnumRead)},
		{string(EnumRolesSuperadmin), string(modules.EnumModulePermissions), string(db.PermissionNamesEnumUpdate)},
		{string(EnumRolesSuperadmin), string(modules.EnumModulePermissions), string(db.PermissionNamesEnumDelete)},

		{string(EnumRolesSuperadmin), string(modules.EnumModuleUserRoles), string(db.PermissionNamesEnumCreate)},
		{string(EnumRolesSuperadmin), string(modules.EnumModuleUserRoles), string(db.PermissionNamesEnumRead)},
		{string(EnumRolesSuperadmin), string(modules.EnumModuleUserRoles), string(db.PermissionNamesEnumUpdate)},
		{string(EnumRolesSuperadmin), string(modules.EnumModuleUserRoles), string(db.PermissionNamesEnumDelete)},

		{string(EnumRolesSuperadmin), string(modules.EnumModuleRolePermissions), string(db.PermissionNamesEnumCreate)},
		{string(EnumRolesSuperadmin), string(modules.EnumModuleRolePermissions), string(db.PermissionNamesEnumRead)},
		{string(EnumRolesSuperadmin), string(modules.EnumModuleRolePermissions), string(db.PermissionNamesEnumUpdate)},
		{string(EnumRolesSuperadmin), string(modules.EnumModuleRolePermissions), string(db.PermissionNamesEnumDelete)},

		// ADMIN PRIVILEGES
		{string(EnumRolesAdmin), string(modules.EnumModuleUsers), string(db.PermissionNamesEnumCreate)},
		{string(EnumRolesAdmin), string(modules.EnumModuleUsers), string(db.PermissionNamesEnumRead)},
		{string(EnumRolesAdmin), string(modules.EnumModuleUsers), string(db.PermissionNamesEnumUpdate)},
		{string(EnumRolesAdmin), string(modules.EnumModuleUsers), string(db.PermissionNamesEnumDelete)},

		{string(EnumRolesAdmin), string(modules.EnumModulePosts), string(db.PermissionNamesEnumCreate)},
		{string(EnumRolesAdmin), string(modules.EnumModulePosts), string(db.PermissionNamesEnumRead)},
		{string(EnumRolesAdmin), string(modules.EnumModulePosts), string(db.PermissionNamesEnumUpdate)},
		{string(EnumRolesAdmin), string(modules.EnumModulePosts), string(db.PermissionNamesEnumDelete)},

		// USER PRIVILEGES
		{string(EnumRolesUser), string(modules.EnumModuleUsers), string(db.PermissionNamesEnumCreate)},
		{string(EnumRolesUser), string(modules.EnumModuleUsers), string(db.PermissionNamesEnumRead)},
		{string(EnumRolesUser), string(modules.EnumModuleUsers), string(db.PermissionNamesEnumUpdate)},
		{string(EnumRolesUser), string(modules.EnumModuleUsers), string(db.PermissionNamesEnumDelete)},

		{string(EnumRolesUser), string(modules.EnumModulePosts), string(db.PermissionNamesEnumCreate)},
		{string(EnumRolesUser), string(modules.EnumModulePosts), string(db.PermissionNamesEnumRead)},
		{string(EnumRolesUser), string(modules.EnumModulePosts), string(db.PermissionNamesEnumUpdate)},
		{string(EnumRolesUser), string(modules.EnumModulePosts), string(db.PermissionNamesEnumDelete)},
	}

	roleRepository := roles.NewRepository(database)
	allRoles, err := roleRepository.GetAll(ctx, "", 10000, 0)
	if err != nil {
		return err
	}

	moduleRepository := modules.NewRepository(database)
	allModules, err := moduleRepository.GetAll(ctx, "", 10000, 0)
	if err != nil {
		return err
	}

	permissionRepository := permissions.NewRepository(database)
	allPermissions, err := permissionRepository.GetAll(ctx, "", 10000, 0)
	if err != nil {
		return err
	}

	rolePermissionsRepository := role_permissions.NewRepository(database)
	for _, privilege := range privileges {
		role, err := findRole(allRoles, privilege[0])
		if err != nil {
			return err
		}
		module, err := findModule(allModules, privilege[1])
		if err != nil {
			return err
		}
		permission, err := findPermission(allPermissions, privilege[2], module.ID)
		if err != nil {
			return err
		}

		creatParams := db.CreateRolePermissionParams{
			ID:           uuid.New(),
			RoleID:       role.ID,
			PermissionID: permission.ID,
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
		}

		rolePermissionsRepository.Create(ctx, creatParams)
	}

	return nil
}

func CreateSuperAdmin(ctx context.Context, database *db.Queries) error {
	roleRepository := roles.NewRepository(database)
	allRoles, err := roleRepository.GetAll(ctx, "", 10000, 0)
	if err != nil {
		return err
	}
	superadminRole, err := findRole(allRoles, string(EnumRolesSuperadmin))

	password := "password"
	hashedPassword, err := cryptography.HashPassword(password)
	if err != nil {
		return err
	}

	userRepository := users.NewRepository(database)
	userParams := db.CreateUserParams{
		ID:        uuid.New(),
		Name:      "Superadmin",
		Email:     "superadmin@mail.com",
		Password:  hashedPassword,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	createdAdmin, err := userRepository.Create(ctx, userParams)
	if err != nil {
		return err
	}

	userRoleRepository := user_roles.NewRepository(database)
	userRoleParam := db.CreateUserRoleParams{
		ID:        uuid.New(),
		UserID:    createdAdmin.ID,
		RoleID:    superadminRole.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_, err = userRoleRepository.Create(ctx, userRoleParam)
	if err != nil {
		return err
	}

	return nil
}

func All(ctx context.Context, database *db.Queries) error {
	err := Modules(ctx, database)
	if err != nil {
		return err
	}
	err = Roles(ctx, database)
	if err != nil {
		return err
	}
	err = Permissions(ctx, database)
	if err != nil {
		return err
	}
	err = RolePermissions(ctx, database)
	if err != nil {
		return err
	}
	err = CreateSuperAdmin(ctx, database)
	if err != nil {
		return err
	}

	return nil
}

func findRole(roles []db.Role, name string) (db.Role, error) {
	for _, role := range roles {
		if role.Name == name {
			return role, nil
		}
	}

	return db.Role{}, errors.New("Role not found")
}

func findModule(modules []db.Module, name string) (db.Module, error) {
	for _, module := range modules {
		if module.Name == name {
			return module, nil
		}
	}

	return db.Module{}, errors.New("Module not found")
}

func findPermission(permissions []db.Permission, name string, moduleID uuid.UUID) (db.Permission, error) {
	for _, permission := range permissions {
		if string(permission.Name) == name && permission.ModuleID == moduleID {
			return permission, nil
		}
	}

	return db.Permission{}, errors.New("Permission not found")
}
