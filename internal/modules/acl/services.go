package acl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	constants "github.com/max-fletcher/golang_web_server_boilerplate/helpers/const"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
	common_errors "github.com/max-fletcher/golang_web_server_boilerplate/internal/errors"
)

type Service interface {
	Create(ctx context.Context, createRolePermissionInput CreateRolePermissionInput) (db.GetRolePermissionByIDRow, error)
	GetAll(ctx context.Context, filterString string, limit int, offset int) ([]db.GetUsersWithRolesAndPermissionsRow, int, error)
	GetByUserID(ctx context.Context, id uuid.UUID) (db.GetUserWithRolesAndPermissionsByUserIDRow, error)
	GetByID(ctx context.Context, id uuid.UUID) (db.GetRolePermissionByIDRow, error)
	Delete(ctx context.Context, id uuid.UUID) (db.GetRolePermissionByIDRow, error)
	DeleteByRoleIDAndPermissionID(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) (db.GetRolePermissionByRoleIDAndPermissionIDRow, error)
}

type service struct {
	repository Repository
}

func NewService(
	repository Repository,
) *service {
	return &service{
		repository: repository,
	}
}

func (service *service) Create(ctx context.Context, createRolePermissionInput CreateRolePermissionInput) (db.GetRolePermissionByIDRow, error) {
	// #TODO: CHECK IF ROLE AND PERMISSION EXISTS

	// 1st param: context for the request
	// 2nd param: the struct that we want to pass so it saves the underlying data in DB
	rolePermission, err := service.repository.Create(ctx, db.CreateRolePermissionParams{
		ID:           uuid.New(),
		RoleID:       createRolePermissionInput.RoleID,
		PermissionID: createRolePermissionInput.PermissionID,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	})
	if err != nil {
		pgErr := common_errors.GetPostgresError(err)
		if pgErr.Code == constants.PGUniqueViolationCode {
			return db.GetRolePermissionByIDRow{}, ErrRolePermissionWithRoleIDAndPermissionIDAlreadyExists{
				RoleID:       createRolePermissionInput.RoleID,
				PermissionID: createRolePermissionInput.PermissionID,
				Err:          err,
			}
		}

		return db.GetRolePermissionByIDRow{}, ErrRolePermissionCreateFailed{
			CreateErr: err,
		}
	}

	// Fetch with details/names(uses JOIN)
	fmt.Println("Creat Here1")
	rolePermissionWithDetails, err := service.repository.GetByID(ctx, rolePermission.ID)
	if err != nil {
		fmt.Println("Creat Here2")
		return db.GetRolePermissionByIDRow{}, ErrRolePermissionCreateFailed{
			CreateErr: err,
		}
	}
	fmt.Println("Creat Here3")

	// #TODO: FORMAT rolePermissionWithDetails AND RETURN A STRUCT
	return rolePermissionWithDetails, nil
}

func (service *service) GetAll(ctx context.Context, filterString string, limit int, offset int) ([]db.GetUsersWithRolesAndPermissionsRow, int, error) {
	// 1st param: context for the request
	usersWithRolesAndPermissions, err := service.repository.GetAll(ctx, filterString, limit, offset)
	if err != nil {
		return []db.GetUsersWithRolesAndPermissionsRow{}, 0, ErrRolePermissionsFetchFailed{
			FetchErr: err,
		}
	}

	total, err := service.repository.GetAllCount(ctx, filterString)
	if err != nil {
		var bigInt64ToIntError common_errors.ErrBigInt64ToIntError
		if errors.As(err, &bigInt64ToIntError) {
			return []db.GetUsersWithRolesAndPermissionsRow{}, 0, fmt.Errorf("Limit and/or offset value out of range")
		}

		return []db.GetUsersWithRolesAndPermissionsRow{}, 0, ErrRolePermissionsFetchFailed{
			FetchErr: err,
		}
	}

	// #TODO: FORMAT rolePermissionWithDetails AND RETURN A STRUCT
	return usersWithRolesAndPermissions, total, nil
}

func (service *service) GetByUserID(ctx context.Context, id uuid.UUID) (db.GetUserWithRolesAndPermissionsByUserIDRow, error) {
	// 1st param: context for the request
	// 2nd param: id(type uuid) param
	usersWithRolesAndPermissions, err := service.repository.GetByUserID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // check if error is of type sql.ErrNoRows
			return db.GetUserWithRolesAndPermissionsByUserIDRow{}, ErrRolePermissionWithUserIdNotFound{
				ID: id,
			}
		}

		return db.GetUserWithRolesAndPermissionsByUserIDRow{}, ErrRolePermissionWithUserFetchFailed{
			FetchErr: err,
		}
	}
	fmt.Println("GetById rolePermission", usersWithRolesAndPermissions)
	if len(usersWithRolesAndPermissions) == 0 {
		return db.GetUserWithRolesAndPermissionsByUserIDRow{}, ErrRolePermissionWithUserIdNotFound{
			ID: id,
		}
	}

	// #TODO: FORMAT rolePermissions AND RETURN A UNIFIED STRUCT
	return usersWithRolesAndPermissions[0], nil
}

func (service *service) GetByID(ctx context.Context, id uuid.UUID) (db.GetRolePermissionByIDRow, error) {
	// 1st param: context for the request
	// 2nd param: id(type uuid) param
	rolePermission, err := service.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // check if error is of type sql.ErrNoRows
			return db.GetRolePermissionByIDRow{}, ErrRolePermissionWithUserIdNotFound{
				ID: id,
			}
		}

		return db.GetRolePermissionByIDRow{}, ErrRolePermissionWithUserFetchFailed{
			FetchErr: err,
		}
	}

	return rolePermission, nil
}

func (service *service) Delete(ctx context.Context, id uuid.UUID) (db.GetRolePermissionByIDRow, error) {
	existingRolePermission, err := service.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // check if error is of type sql.ErrNoRows
			return db.GetRolePermissionByIDRow{}, ErrRolePermissionWithIdNotFound{
				ID: id,
			}
		}

		return db.GetRolePermissionByIDRow{}, ErrRolePermissionFetchFailed{
			FetchErr: err,
		}
	}

	_, err = service.repository.DeleteByID(ctx, id)
	if err != nil {
		return db.GetRolePermissionByIDRow{}, ErrRolePermissionDeleteFailed{
			DeleteErr: err,
		}
	}

	return existingRolePermission, nil
}

func (service *service) DeleteByRoleIDAndPermissionID(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) (db.GetRolePermissionByRoleIDAndPermissionIDRow, error) {
	existingRolePermission, err := service.repository.GetByRoleIDAndPermissionID(ctx, roleID, permissionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // check if error is of type sql.ErrNoRows
			return db.GetRolePermissionByRoleIDAndPermissionIDRow{}, ErrRolePermissionWithRoleIdAndPermissionIDNotFound{
				RoleID:       roleID,
				PermissionID: permissionID,
			}
		}

		return db.GetRolePermissionByRoleIDAndPermissionIDRow{}, ErrRolePermissionFetchFailed{
			FetchErr: err,
		}
	}

	_, err = service.repository.DeleteByRoleIDAndPermissionID(ctx, roleID, permissionID)
	if err != nil {
		return db.GetRolePermissionByRoleIDAndPermissionIDRow{}, ErrRolePermissionDeleteFailed{
			DeleteErr: err,
		}
	}

	return existingRolePermission, nil
}
