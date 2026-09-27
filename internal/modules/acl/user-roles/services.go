package user_roles

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	constants "github.com/max-fletcher/golang_web_server_boilerplate/helpers/const"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/formatters"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
	common_errors "github.com/max-fletcher/golang_web_server_boilerplate/internal/errors"
)

type UserService interface {
	GetByID(ctx context.Context, id uuid.UUID) (db.User, error)
}

type RoleService interface {
	GetByID(ctx context.Context, id uuid.UUID) (db.Role, error)
}

type Service interface {
	Create(ctx context.Context, createUserRoleInput CreateUserRoleInput) (formatters.SingleUserWRoles, error)
	GetAll(ctx context.Context, filterString string, limit int, offset int) ([]db.GetUserRolesRow, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (db.GetUserRoleByIdRow, error)
	GetByUserID(ctx context.Context, id uuid.UUID) ([]db.GetUserRolesByUserIDRow, error)
	Delete(ctx context.Context, id uuid.UUID) (db.GetUserRoleByIdRow, error)
	DeleteByUserIDAndRoleID(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) (db.GetUserRoleByUserIDAndRoleIDRow, error)
}

type service struct {
	repository  Repository
	userService UserService
	roleService RoleService
}

func NewService(repository Repository, userService UserService, roleService RoleService) *service {
	return &service{
		repository:  repository,
		userService: userService,
		roleService: roleService,
	}
}

func (service *service) Create(ctx context.Context, createUserRoleInput CreateUserRoleInput) (formatters.SingleUserWRoles, error) {
	_, err := service.userService.GetByID(ctx, createUserRoleInput.UserID)
	if err != nil {
		return formatters.SingleUserWRoles{}, err
	}
	_, err = service.roleService.GetByID(ctx, createUserRoleInput.RoleID)
	if err != nil {
		return formatters.SingleUserWRoles{}, err
	}

	// 1st param: context for the request
	// 2nd param: the struct that we want to pass so it saves the underlying data in DB
	userRole, err := service.repository.Create(ctx, db.CreateUserRoleParams{
		ID:        uuid.New(),
		UserID:    createUserRoleInput.UserID,
		RoleID:    createUserRoleInput.RoleID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		pgErr := common_errors.GetPostgresError(err)
		if pgErr.Code == constants.PGUniqueViolationCode {
			return formatters.SingleUserWRoles{}, ErrUserRoleWithUserIdAndRoleIdAlreadyExists{
				RoleID: createUserRoleInput.RoleID,
				UserID: createUserRoleInput.UserID,
				Err:    err,
			}
		}

		if pgErr.Code == constants.PGForeignKeyViolationCode {
			return formatters.SingleUserWRoles{}, ErrUserRoleInvalidUserIdOrRoleId{
				RoleID: createUserRoleInput.RoleID,
				UserID: createUserRoleInput.UserID,
				Err:    err,
			}
		}

		return formatters.SingleUserWRoles{}, ErrUserRoleCreateFailed{
			CreateErr: err,
		}
	}

	userRoleWithDetails, err := service.repository.GetByID(ctx, userRole.ID)
	if err != nil {
		return formatters.SingleUserWRoles{}, ErrUserRoleCreateFailed{
			CreateErr: err,
		}
	}

	formattedData := formatters.DatabaseUserWRoleToUserWRole(userRoleWithDetails)
	return formattedData, nil
}

func (service *service) GetAll(ctx context.Context, filterString string, limit int, offset int) ([]db.GetUserRolesRow, int, error) {
	// 1st param: context for the request
	userRoles, err := service.repository.GetAll(ctx, filterString, limit, offset)
	if err != nil {
		return []db.GetUserRolesRow{}, 0, ErrUserRolesFetchFailed{
			FetchErr: err,
		}
	}

	total, err := service.repository.GetAllCount(ctx, filterString)
	if err != nil {
		var bigInt64ToIntError common_errors.ErrBigInt64ToIntError
		if errors.As(err, &bigInt64ToIntError) {
			return []db.GetUserRolesRow{}, 0, fmt.Errorf("Limit and/or offset value out of range")
		}

		return []db.GetUserRolesRow{}, 0, ErrUserRolesFetchFailed{
			FetchErr: err,
		}
	}

	return userRoles, total, nil
}

func (service *service) GetByID(ctx context.Context, id uuid.UUID) (db.GetUserRoleByIdRow, error) {
	// 1st param: context for the request
	// 2nd param: id(type uuid) param
	userRole, err := service.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // check if error is of type sql.ErrNoRows
			return db.GetUserRoleByIdRow{}, ErrUserRoleWithUserIdNotFound{
				ID: id,
			}
		}

		return db.GetUserRoleByIdRow{}, ErrUserRoleWithUserFetchFailed{
			FetchErr: err,
		}
	}

	return userRole, nil
}

func (service *service) GetByUserID(ctx context.Context, userID uuid.UUID) ([]db.GetUserRolesByUserIDRow, error) {
	// 1st param: context for the request
	// 2nd param: id(type uuid) param
	userRoles, err := service.repository.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // check if error is of type sql.ErrNoRows
			return []db.GetUserRolesByUserIDRow{}, ErrUserRoleWithUserIdNotFound{
				ID: userID,
			}
		}

		return []db.GetUserRolesByUserIDRow{}, ErrUserRoleWithUserFetchFailed{
			FetchErr: err,
		}
	}

	// #TODO: FORMAT userRoles AND RETURN A UNIFIED STRUCT
	return userRoles, nil
}

func (service *service) Delete(ctx context.Context, id uuid.UUID) (db.GetUserRoleByIdRow, error) {
	existingUserRole, err := service.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // check if error is of type sql.ErrNoRows
			return db.GetUserRoleByIdRow{}, ErrUserRoleWithIdNotFound{
				ID: id,
			}
		}

		return db.GetUserRoleByIdRow{}, ErrUserRoleFetchFailed{
			FetchErr: err,
		}
	}

	_, err = service.repository.DeleteByID(ctx, id)
	if err != nil {
		return db.GetUserRoleByIdRow{}, ErrUserRoleDeleteFailed{
			DeleteErr: err,
		}
	}

	return existingUserRole, nil
}

func (service *service) DeleteByUserIDAndRoleID(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) (db.GetUserRoleByUserIDAndRoleIDRow, error) {
	existingUserRole, err := service.repository.GetUserRoleByUserIDAndRoleID(ctx, userID, roleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // check if error is of type sql.ErrNoRows
			return db.GetUserRoleByUserIDAndRoleIDRow{}, ErrUserRoleWithUserIdAndRoleIdNotFound{
				RoleID: roleID,
				UserID: userID,
			}
		}

		return db.GetUserRoleByUserIDAndRoleIDRow{}, ErrUserRoleFetchFailed{
			FetchErr: err,
		}
	}

	_, err = service.repository.DeleteByUserIDAndRoleID(ctx, userID, roleID)
	if err != nil {
		return db.GetUserRoleByUserIDAndRoleIDRow{}, ErrUserRoleDeleteFailed{
			DeleteErr: err,
		}
	}

	return existingUserRole, nil
}
