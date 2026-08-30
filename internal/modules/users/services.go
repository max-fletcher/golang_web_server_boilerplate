package users

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/crypto"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
	common_errors "github.com/max-fletcher/golang_web_server_boilerplate/internal/errors"
)

type Service interface {
	Create(ctx context.Context, params CreateUserRequest) (db.User, error)
	GetAll(ctx context.Context, params db.GetUsersParams) ([]db.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (db.User, error)
	GetByEmail(ctx context.Context, email string) (db.User, error)
	// Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repository Repository
}

func NewService(db *db.Queries) *service {
	return &service{
		repository: NewRepository(db),
	}
}

func (service *service) Create(ctx context.Context, params CreateUserRequest) (db.User, error) {
	_, err := service.GetByEmail(ctx, params.Email)
	if err == nil {
		return db.User{}, ErrUserWithEmailAlreadyExists{
			Email: params.Email,
		}
	}

	// If error exists but doesn't match the errors that GetByEmail() sends back
	var userWithEmailNotFoundErr ErrUserWithEmailNotFound
	var userFetchFailedErr ErrUserFetchFailed
	if err != nil && !errors.As(err, &userWithEmailNotFoundErr) && !errors.As(err, &userFetchFailedErr) {
		return db.User{}, ErrUserFetchFailed{
			fetchErr: err,
		}
	}

	hashedPassword, err := crypto.HashPassword(params.Password)
	if err != nil {
		return db.User{}, common_errors.ErrHashingPassword{
			HashErr: err,
		}
	}
	params.Password = hashedPassword

	// 1st param: context for the request
	// 2nd param: the struct that we want to pass so it saves the underlying data in DB
	user, err := service.repository.Create(ctx, db.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		Email:     params.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return db.User{}, ErrUserCreateFailed{
			createErr: err,
		}
	}

	return user, nil
}

// handler method that handles fetching all users. The addition of (handler *Handler)
// turns it into a method for handler.
func (service *service) GetAll(ctx context.Context, params db.GetUsersParams) ([]db.User, error) {
	// 1st param: context for the request
	users, err := service.repository.GetAll(ctx, params)
	if err != nil {
		return []db.User{}, ErrUsersFetchFailed{
			fetchErr: err,
		}
	}

	return users, nil
}

func (service *service) GetByID(ctx context.Context, id uuid.UUID) (db.User, error) {
	// 1st param: context for the request
	// 2nd param: id(type uuid) param
	user, err := service.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // check if error is of type sql.ErrNoRows
			return db.User{}, ErrUserWithIdNotFound{
				ID: id,
			}
		}

		return db.User{}, ErrUserFetchFailed{
			fetchErr: err,
		}
	}

	return user, nil
}

func (service *service) GetByEmail(ctx context.Context, email string) (db.User, error) {
	// 1st param: context for the request
	// 2nd param: id(type uuid) param
	user, err := service.repository.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // check if error is of type sql.ErrNoRows
			return db.User{}, ErrUserWithEmailNotFound{
				Email: email,
			}
		}

		return db.User{}, ErrUserFetchFailed{
			fetchErr: err,
		}
	}

	return user, nil
}
