package permissions

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

type ModuleService interface { // For DI. Used in validating user in structs.go
	GetByID(ctx context.Context, id uuid.UUID) (db.Module, error)
}

type Service interface {
	Create(ctx context.Context, createPermissionInput CreatePermissionInput) (db.Permission, error)
	GetAll(ctx context.Context, filterString string, limit int, offset int) ([]db.Permission, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (db.Permission, error)
}

type service struct {
	repository    Repository
	moduleService ModuleService // Using DI. See server/server.go where we are passing user service as 2nd param here.
}

func NewService(repository Repository, moduleService ModuleService) *service {
	return &service{
		repository:    repository,
		moduleService: moduleService, // Using DI. See server/server.go where we are passing user service as 2nd param here.
	}
}

func (service *service) Create(ctx context.Context, createPermissionInput CreatePermissionInput) (db.Permission, error) {
	_, err := service.moduleService.GetByID(ctx, createPermissionInput.ModuleId)
	if err != nil {
		return db.Permission{}, err
	}

	// 1st param: context for the request
	// 2nd param: the struct that we want to pass so it saves the underlying data in DB
	permission, err := service.repository.Create(ctx, db.CreatePermissionParams{
		ID:        uuid.New(),
		Name:      db.PermissionNamesEnum(createPermissionInput.Name),
		ModuleID:  createPermissionInput.ModuleId,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		pgErr := common_errors.GetPostgresError(err)
		if pgErr.Code == constants.PGUniqueViolationCode {
			return db.Permission{}, ErrModuleWithPermissionAlreadyExists{
				Err: err,
			}
		}

		return db.Permission{}, ErrPermissionCreateFailed{
			createErr: err,
		}
	}

	return permission, nil
}

func (service *service) GetAll(ctx context.Context, filterString string, limit int, offset int) ([]db.Permission, int, error) {
	// 1st param: context for the request
	permissions, err := service.repository.GetAll(ctx, filterString, limit, offset)
	if err != nil {
		return []db.Permission{}, 0, ErrPermissionFetchFailed{
			FetchErr: err,
		}
	}

	total, err := service.repository.GetAllCount(ctx, filterString)
	if err != nil {
		var bigInt64ToIntError common_errors.ErrBigInt64ToIntError
		if errors.As(err, &bigInt64ToIntError) {
			return []db.Permission{}, 0, fmt.Errorf("Limit and/or offset value out of range")
		}

		return []db.Permission{}, 0, ErrPermissionFetchFailed{
			FetchErr: err,
		}
	}

	return permissions, total, nil
}

func (service *service) GetByID(ctx context.Context, id uuid.UUID) (db.Permission, error) {
	// 1st param: context for the request
	// 2nd param: id(type uuid) param
	permission, err := service.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // check if error is of type sql.ErrNoRows
			return db.Permission{}, ErrPermissionWithIdNotFound{
				ID: id,
			}
		}

		return db.Permission{}, ErrPermissionFetchFailed{
			FetchErr: err,
		}
	}

	return permission, nil
}
