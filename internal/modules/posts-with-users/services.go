package posts_with_users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	constants "github.com/max-fletcher/golang_web_server_boilerplate/helpers/const"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/crypto"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/formatters"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
	common_errors "github.com/max-fletcher/golang_web_server_boilerplate/internal/errors"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/posts"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/users"
)

type Service interface {
	Create(ctx context.Context, params CreatePostWithUserInput) (db.User, error)
	GetAll(ctx context.Context, filterString string, limit int, offset int) ([]db.GetPostsWithUserRow, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (db.GetPostWithUserByIdRow, error)
}

type service struct {
	repository Repository
	sqlDB      *sql.DB
	DB         *db.Queries // We are importing(using DI) db.Queries here so we can use database transactions
}

func NewService(repository Repository, conn *sql.DB, database *db.Queries) *service {
	return &service{
		repository: repository,
		sqlDB:      conn,
		DB:         database, // We are importing(using DI) db.Queries here so we can use database transactions
	}
}

func (service *service) Create(ctx context.Context, createPostWithUserInput CreatePostWithUserInput) (db.User, error) {
	hashedPassword, err := crypto.HashPassword(createPostWithUserInput.Password)
	if err != nil {
		return db.User{}, common_errors.ErrHashingPassword{
			HashErr: err,
		}
	}
	createPostWithUserInput.Password = hashedPassword

	tx, err := service.sqlDB.BeginTx(ctx, nil) // begin transaction
	if err != nil {
		return db.User{}, err
	}
	defer tx.Rollback()

	// IMPORTANT: txDB contains all the queries from sqlc. So initializing txDB like this and using it to make queries will cause
	// all queries to be part of one transaction
	txDB := service.DB.WithTx(tx)

	// we can do this because both txDB and "database" param(1st param) in NewRepository is of type *db.Queries. This means
	// since we passed txDB into these repositories, any query executed from inside these are now part of this database transaction.
	// *IMPORTANT: replace these later with DI(see posts service that uses the UserExistenceChecker interface to import usersService)
	userRepository := users.NewRepository(txDB)
	postRepository := posts.NewRepository(txDB)

	// 1st param: context for the request
	// 2nd param: the struct that we want to pass so it saves the underlying data in DB
	user, err := userRepository.Create(ctx, db.CreateUserParams{
		ID:        uuid.New(),
		Name:      createPostWithUserInput.Name,
		Email:     createPostWithUserInput.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		pgErr := common_errors.GetPostgresError(err)
		if pgErr.Code == constants.PGUniqueViolationCode { // handle duplicate constraint error i.e unique in email
			return db.User{}, ErrUserWithEmailAlreadyExists{
				Email: createPostWithUserInput.Email,
			}
		}

		return db.User{}, ErrUserCreateFailed{
			createErr: err,
		}
	}

	_, err = postRepository.Create(ctx, db.CreatePostParams{
		ID:        uuid.New(),
		Title:     createPostWithUserInput.Title,
		Content:   formatters.StringPointerToNullString(createPostWithUserInput.Content),
		Photo:     formatters.StringPointerToNullString(createPostWithUserInput.Photo),
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return db.User{}, ErrPostCreateFailed{
			createErr: err,
		}
	}

	if err := tx.Commit(); err != nil {
		return db.User{}, err
	}

	return user, nil
}

func (service *service) GetAll(ctx context.Context, filterString string, limit int, offset int) ([]db.GetPostsWithUserRow, int, error) {
	// 1st param: context for the request
	users, err := service.repository.GetAll(ctx, filterString, limit, offset)
	if err != nil {
		return []db.GetPostsWithUserRow{}, 0, ErrUsersFetchFailed{
			fetchErr: err,
		}
	}

	postRepository := posts.NewRepository(service.DB)
	total, err := postRepository.GetAllCount(ctx, filterString)
	if err != nil {
		var bigInt64ToIntError common_errors.ErrBigInt64ToIntError
		if errors.As(err, &bigInt64ToIntError) {
			return []db.GetPostsWithUserRow{}, 0, fmt.Errorf("Limit and/or offset value out of range")
		}

		return []db.GetPostsWithUserRow{}, 0, ErrUsersFetchFailed{
			fetchErr: err,
		}
	}

	return users, total, nil
}

func (service *service) GetByID(ctx context.Context, id uuid.UUID) (db.GetPostWithUserByIdRow, error) {
	// 1st param: context for the request
	// 2nd param: id(type uuid) param
	user, err := service.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // check if error is of type sql.ErrNoRows
			return db.GetPostWithUserByIdRow{}, ErrUserWithIdNotFound{
				ID: id,
			}
		}

		return db.GetPostWithUserByIdRow{}, ErrUserFetchFailed{
			fetchErr: err,
		}
	}

	return user, nil
}
