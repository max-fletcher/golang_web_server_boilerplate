package users

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
)

type Repository interface {
	Create(ctx context.Context, params db.CreateUserParams) (db.User, error)
	GetAll(ctx context.Context, params db.GetUsersParams) ([]db.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (db.User, error)
	GetByEmail(ctx context.Context, email string) (db.User, error)
	// Delete(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	DB *db.Queries // Reference to a DB connection. Will be used to query data form DB.
}

func NewRepository(db *db.Queries) *repository {
	return &repository{
		DB: db,
	}
}

// No need to pass the entire response writer for these repository functions. Just passing the context from response writer(r.Context()) from where
// this is called from (i.e service) works just fine instead of passing the response writer, then using r.Context() as first param of sqlc function

func (repository *repository) Create(ctx context.Context, params db.CreateUserParams) (db.User, error) {
	return repository.DB.CreateUser(ctx, db.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		Email:     params.Email,
		Password:  params.Password,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
}

func (repository *repository) GetAll(ctx context.Context, params db.GetUsersParams) ([]db.User, error) {
	return repository.DB.GetUsers(ctx, params)
}

func (repository *repository) GetByID(ctx context.Context, id uuid.UUID) (db.User, error) {
	return repository.DB.GetUserById(ctx, id)
}

func (repository *repository) GetByEmail(ctx context.Context, email string) (db.User, error) {
	return repository.DB.GetUserByEmail(ctx, email)
}
