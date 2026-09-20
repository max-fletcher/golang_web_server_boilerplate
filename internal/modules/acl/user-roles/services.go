package user_roles

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
	Create(ctx context.Context, createUserRoleInput CreateUserRoleInput) (db.UserRole, error)
	GetAll(ctx context.Context, filterString string, limit int, offset int) ([]db.GetUserRolesRow, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (db.UserRole, error)
	GetByUserID(ctx context.Context, id uuid.UUID) (db.GetUserRolesByUserIDRow, error)
	Delete(ctx context.Context, id uuid.UUID) (db.UserRole, error)
	DeleteByUserIDAndRoleID(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) (db.GetUserRoleByUserIDAndRoleIDRow, error)
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

func (service *service) Create(ctx context.Context, createUserRoleInput CreateUserRoleInput) (db.UserRole, error) {
	// #TODO: CHECK IF ROLE AND PERMISSION EXISTS

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
			return db.UserRole{}, ErrUserRoleWithUserIdAndRoleIdAlreadyExists{
				RoleID: createUserRoleInput.RoleID,
				UserID: createUserRoleInput.UserID,
				Err:    err,
			}
		}

		if pgErr.Code == constants.PGForeignKeyViolationCode {
			return db.UserRole{}, ErrUserRoleInvaliduserIdOrRoleId{
				RoleID: createUserRoleInput.RoleID,
				UserID: createUserRoleInput.UserID,
				Err:    err,
			}
		}

		return db.UserRole{}, ErrUserRoleCreateFailed{
			CreateErr: err,
		}
	}

	userRoleWithDetails, err := service.repository.GetByID(ctx, userRole.ID)
	if err != nil {
		return db.UserRole{}, ErrUserRoleCreateFailed{
			CreateErr: err,
		}
	}

	return userRoleWithDetails, nil
}

func (service *service) GetAll(ctx context.Context, filterString string, limit int, offset int) ([]db.GetUserRolesRow, int, error) {
	// 1st param: context for the request
	userRoles, err := service.repository.GetAll(ctx, filterString, limit, offset)
	if err != nil {
		return []db.GetUserRolesRow{}, 0, ErrUserRolesFetchFailed{
			FetchErr: err,
		}
	}

	total, err := service.repository.GetAllCount(ctx)
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

func (service *service) GetByID(ctx context.Context, id uuid.UUID) (db.UserRole, error) {
	// 1st param: context for the request
	// 2nd param: id(type uuid) param
	userRole, err := service.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // check if error is of type sql.ErrNoRows
			return db.UserRole{}, ErrUserRoleWithUserIdNotFound{
				ID: id,
			}
		}

		return db.UserRole{}, ErrUserRoleWithUserFetchFailed{
			FetchErr: err,
		}
	}

	return userRole, nil
}

func (service *service) GetByUserID(ctx context.Context, id uuid.UUID) ([]db.GetUserRolesByUserIDRow, error) {
	// 1st param: context for the request
	// 2nd param: id(type uuid) param
	userRoles, err := service.repository.GetByUserID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // check if error is of type sql.ErrNoRows
			return []db.GetUserRolesByUserIDRow{}, ErrUserRoleWithUserIdNotFound{
				ID: id,
			}
		}

		return []db.GetUserRolesByUserIDRow{}, ErrUserRoleWithUserFetchFailed{
			FetchErr: err,
		}
	}
	fmt.Println("GetById userRole", userRoles)

	// #TODO: FORMAT userRoles AND RETURN A UNIFIED STRUCT
	return userRoles, nil
}

func (service *service) Delete(ctx context.Context, id uuid.UUID) (db.UserRole, error) {
	existingUserRole, err := service.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // check if error is of type sql.ErrNoRows
			return db.UserRole{}, ErrUserRoleWithIdNotFound{
				ID: id,
			}
		}

		return db.UserRole{}, ErrUserRoleFetchFailed{
			FetchErr: err,
		}
	}

	_, err = service.repository.DeleteByID(ctx, id)
	if err != nil {
		return db.UserRole{}, ErrUserRoleDeleteFailed{
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
