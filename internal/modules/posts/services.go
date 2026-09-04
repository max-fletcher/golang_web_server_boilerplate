package posts

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

type UserExistenceChecker interface { // For DI. Used in validating user in structs.go
	GetByID(ctx context.Context, id uuid.UUID) (db.User, error)
}

type Service interface {
	Create(ctx context.Context, createPostInput CreatePostInput) (db.Post, error)
	GetAll(ctx context.Context, filterString string, limit int, offset int) ([]db.Post, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (db.Post, error)
	Update(ctx context.Context, id uuid.UUID, updatePostInput UpdatePostInput) (db.Post, error)
	Delete(ctx context.Context, id uuid.UUID) (db.Post, error)
}

type service struct {
	repository  Repository
	userChecker UserExistenceChecker // Using DI. See server/server.go where we are passing user service as 2nd param here.
}

// Using different structure so that we can prevent circular dependency
//
//	func NewService(db *db.Queries) *service {
//		return &service{
//			repository: NewRepository(db),
//		}
//	}
func NewService(repository Repository, userChecker UserExistenceChecker) *service {
	return &service{
		repository:  repository,
		userChecker: userChecker, // Using DI. See server/server.go where we are passing user service as 2nd param here.
	}
}

func (service *service) Create(ctx context.Context, createPostInput CreatePostInput) (db.Post, error) {
	_, err := service.userChecker.GetByID(ctx, createPostInput.UserId)
	if err != nil {
		return db.Post{}, err
	}

	// 1st param: context for the request
	// 2nd param: the struct that we want to pass so it saves the underlying data in DB
	post, err := service.repository.Create(ctx, db.CreatePostParams{
		ID:      uuid.New(),
		Title:   createPostInput.Title,
		Content: createPostInput.Content,
		// Photo:     createPostInput.Photo,
		UserID:    createPostInput.UserId,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return db.Post{}, ErrPostCreateFailed{
			createErr: err,
		}
	}

	return post, nil
}

func (service *service) GetAll(ctx context.Context, filterString string, limit int, offset int) ([]db.Post, int, error) {
	// 1st param: context for the request
	users, err := service.repository.GetAll(ctx, filterString, limit, offset)
	if err != nil {
		return []db.Post{}, 0, ErrPostsFetchFailed{
			fetchErr: err,
		}
	}

	total, err := service.repository.GetAllCount(ctx, filterString)
	if err != nil {
		var bigInt64ToIntError common_errors.ErrBigInt64ToIntError
		if errors.As(err, &bigInt64ToIntError) {
			return []db.Post{}, 0, fmt.Errorf("Limit and/or offset value out of range")
		}

		return []db.Post{}, 0, ErrPostsFetchFailed{
			fetchErr: err,
		}
	}

	return users, total, nil
}

func (service *service) GetByID(ctx context.Context, id uuid.UUID) (db.Post, error) {
	// 1st param: context for the request
	// 2nd param: id(type uuid) param
	user, err := service.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // check if error is of type sql.ErrNoRows
			return db.Post{}, ErrPostWithIdNotFound{
				ID: id,
			}
		}

		return db.Post{}, ErrPostFetchFailed{
			fetchErr: err,
		}
	}

	return user, nil
}

func (service *service) Update(ctx context.Context, id uuid.UUID, updatePostInput UpdatePostInput) (db.Post, error) {
	_, err := service.GetByID(ctx, id)
	if err != nil {
		return db.Post{}, err
	}

	// 1st param: context for the request
	// 2nd param: the struct that we want to pass so it saves the underlying data in DB
	user, err := service.repository.Update(ctx, db.UpdatePostParams{
		ID:      id,
		Title:   updatePostInput.Title,
		Content: updatePostInput.Content,
		// Photo:     updatePostInput.Photo,
		UserID:    updatePostInput.UserId,
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return db.Post{}, ErrPostUpdateFailed{
			updatePost: err,
		}
	}

	return user, nil
}

func (service *service) Delete(ctx context.Context, id uuid.UUID) (db.Post, error) {
	user, err := service.GetByID(ctx, id)
	if err != nil {
		return db.Post{}, err
	}

	_, err = service.repository.Delete(ctx, id)
	if err != nil {
		return db.Post{}, ErrPostDeleteFailed{
			deleteErr: err,
		}
	}

	return user, nil
}
