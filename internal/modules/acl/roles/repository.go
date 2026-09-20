package roles

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/numbers"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
)

type Repository interface {
	Create(ctx context.Context, params db.CreateRoleParams) (db.Role, error)
	GetAll(ctx context.Context, filterString string, limit int, offset int) ([]db.Role, error)
	GetAllCount(ctx context.Context, filterString string) (int, error)
	GetByID(ctx context.Context, id uuid.UUID) (db.Role, error)
	Update(ctx context.Context, params db.UpdateRoleParams) (db.Role, error)
	Delete(ctx context.Context, id uuid.UUID) (int64, error)
}

type repository struct {
	DB *db.Queries // Reference to a DB connection. Will be used to query data form DB.
}

func NewRepository(database *db.Queries) *repository {
	return &repository{
		DB: database,
	}
}

func (repository *repository) Create(ctx context.Context, params db.CreateRoleParams) (db.Role, error) {
	return repository.DB.CreateRole(ctx, db.CreateRoleParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
}

func (repository *repository) GetAll(ctx context.Context, filterString string, limit int, offset int) ([]db.Role, error) {
	params := db.GetRolesParams{
		Column1: filterString,
		Limit:   int32(limit),
		Offset:  int32(offset),
	}
	return repository.DB.GetRoles(ctx, params)
}

func (repository *repository) GetAllCount(ctx context.Context, filterString string) (int, error) {
	countInt64, err := repository.DB.GetRolesCount(ctx, filterString)
	count, err := numbers.BigInt64ToInt(countInt64)
	return count, err
}

func (repository *repository) GetByID(ctx context.Context, id uuid.UUID) (db.Role, error) {
	return repository.DB.GetRoleById(ctx, id)
}

func (repository *repository) Update(ctx context.Context, params db.UpdateRoleParams) (db.Role, error) {
	return repository.DB.UpdateRole(ctx, db.UpdateRoleParams{
		ID:        params.ID,
		Name:      params.Name,
		UpdatedAt: time.Now().UTC(),
	})
}

func (repository *repository) Delete(ctx context.Context, id uuid.UUID) (int64, error) {
	return repository.DB.DeleteRole(ctx, id)
}
