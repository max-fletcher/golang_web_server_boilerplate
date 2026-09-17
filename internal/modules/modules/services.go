package modules

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
	common_errors "github.com/max-fletcher/golang_web_server_boilerplate/internal/errors"
)

type ModuleService interface { // For DI. Used in validating user in structs.go
	GetByID(ctx context.Context, id uuid.UUID) (db.Module, error)
}

type Service interface {
	Create(ctx context.Context, createModuleInput CreateModuleInput) (db.Module, error)
	GetAll(ctx context.Context, filterString string, limit int, offset int) ([]db.Module, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (db.Module, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{
		repository: repository,
	}
}

func (service *service) Create(ctx context.Context, createModuleInput CreateModuleInput) (db.Module, error) {
	// 1st param: context for the request
	// 2nd param: the struct that we want to pass so it saves the underlying data in DB
	permission, err := service.repository.Create(ctx, db.CreateModuleParams{
		ID:        uuid.New(),
		Name:      string(createModuleInput.Name),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return db.Module{}, ErrModuleCreateFailed{
			CreateErr: err,
		}
	}

	return permission, nil
}

func (service *service) GetAll(ctx context.Context, filterString string, limit int, offset int) ([]db.Module, int, error) {
	// 1st param: context for the request
	permissions, err := service.repository.GetAll(ctx, filterString, limit, offset)
	if err != nil {
		return []db.Module{}, 0, ErrModuleFetchFailed{
			FetchErr: err,
		}
	}

	total, err := service.repository.GetAllCount(ctx, filterString)
	if err != nil {
		var bigInt64ToIntError common_errors.ErrBigInt64ToIntError
		if errors.As(err, &bigInt64ToIntError) {
			return []db.Module{}, 0, fmt.Errorf("Limit and/or offset value out of range")
		}

		return []db.Module{}, 0, ErrModuleFetchFailed{
			FetchErr: err,
		}
	}

	return permissions, total, nil
}

func (service *service) GetByID(ctx context.Context, id uuid.UUID) (db.Module, error) {
	// 1st param: context for the request
	// 2nd param: id(type uuid) param
	permission, err := service.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // check if error is of type sql.ErrNoRows
			return db.Module{}, ErrModuleWithIdNotFound{
				ID: id,
			}
		}

		return db.Module{}, ErrModuleFetchFailed{
			FetchErr: err,
		}
	}

	return permission, nil
}
