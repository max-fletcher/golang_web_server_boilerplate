package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	constants "github.com/max-fletcher/golang_web_server_boilerplate/helpers/const"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/crypto"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/cache"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
	common_errors "github.com/max-fletcher/golang_web_server_boilerplate/internal/errors"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/events"
)

type Service interface {
	Create(ctx context.Context, params CreateUserRequest) (db.User, error)
	GetAll(ctx context.Context, filterString string, limit int, offset int) ([]db.User, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (db.User, error)
	GetByEmail(ctx context.Context, email string) (db.User, error)
	Update(ctx context.Context, id uuid.UUID, params UpdateUserRequest) (db.User, error)
	Delete(ctx context.Context, id uuid.UUID) (db.User, error)
}

type service struct {
	repository Repository
	cache      cache.Cache
	events     events.Publisher
}

// Using different structure so that we can prevent circular dependency
//
//	func NewService(db *db.Queries) *service {
//		return &service{
//			repository: NewRepository(db),
//		}
//	}
func NewService(repository Repository, cache cache.Cache, events events.Publisher) *service {
	return &service{
		repository: repository,
		cache:      cache,
		events:     events,
	}
}

func (service *service) Create(ctx context.Context, params CreateUserRequest) (db.User, error) {
	// Not sure if this is needed anymore since I am checking unique constraint violation below on creat
	// and throwing the exact same error
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
		return db.User{}, common_errors.ErrInternalServer{
			Err: err,
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
		pgErr := common_errors.GetPostgresError(err)
		if pgErr.Code == constants.PGUniqueViolationCode {
			return db.User{}, ErrUserWithEmailAlreadyExists{
				Email: params.Email,
			}
		}

		return db.User{}, ErrUserCreateFailed{
			CreateErr: err,
		}
	}

	if err := service.cache.DeleteByPrefix(ctx, cache.CacheKeyValidPrefixes(CacheKeyUsers)); err != nil { // delete from cache
		fmt.Println("Failed to invalidate by prefix")
	}

	err = service.events.Publish(ctx, events.QueueEventUserCreated,
		events.UserCreated{
			ID: user.ID,
		},
	)
	if err != nil {
		log.Printf("Failed to publish user created event: %v", err)
		return db.User{}, err
	}

	return user, nil
}

func (service *service) GetAll(ctx context.Context, filterString string, limit int, offset int) ([]db.User, int, error) {
	cacheKey := GetAllCacheKey(filterString, limit, offset)
	var cachedUsers GetAllUsersResult
	err := service.cache.Get(ctx, cacheKey, &cachedUsers)
	if err == nil { // return if no errros i.e fetched successfully
		fmt.Println("Cache get successful")
		return cachedUsers.Users, cachedUsers.Total, nil
	}
	// Decide whether cache failure should fail the request. For most caches, I'd allow the request to continue.
	// else {
	// 	return err
	// }

	// 1st param: context for the request
	users, err := service.repository.GetAll(ctx, filterString, limit, offset)
	if err != nil {
		return []db.User{}, 0, ErrUsersFetchFailed{
			fetchErr: err,
		}
	}

	total, err := service.repository.GetAllCount(ctx, filterString)
	if err != nil {
		var bigInt64ToIntError common_errors.ErrBigInt64ToIntError
		if errors.As(err, &bigInt64ToIntError) {
			return []db.User{}, 0, fmt.Errorf("Limit and/or offset value out of range")
		}

		return []db.User{}, 0, ErrUsersFetchFailed{
			fetchErr: err,
		}
	}

	fmt.Println("Cache miss")
	if err := service.cache.Set( // Set cache
		ctx,
		cacheKey,
		GetAllUsersResult{
			Users: users,
			Total: total,
		},
	); err != nil {
		// Failing to set shouldn't fail request since we did get data from database. Just debug why redis is not working
		log.Printf("Failed to store data in redis: %v", err)
	}

	return users, total, nil
}

func (service *service) GetByID(ctx context.Context, id uuid.UUID) (db.User, error) {
	cacheKey := GetByIDCacheKey(id)
	var user db.User
	err := service.cache.Get(ctx, cacheKey, &user)
	if err == nil { // return if no errros i.e fetched successfully
		fmt.Println("Cache get successful")
		return user, nil
	}
	// Decide whether cache failure should fail the request. For most caches, I'd allow the request to continue.
	// else {
	// 	return err
	// }

	// 1st param: context for the request
	// 2nd param: id(type uuid) param
	user, err = service.repository.GetByID(ctx, id)
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

	fmt.Println("Cache miss")
	if err := service.cache.Set(ctx, cacheKey, user); err != nil { // Set cache
		// Failing to set shouldn't fail request since we did get data from database. Just debug why redis is not working
		log.Printf("Failed to store data in redis: %v", err)
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

func (service *service) Update(ctx context.Context, id uuid.UUID, params UpdateUserRequest) (db.User, error) {
	_, err := service.GetByID(ctx, id)
	if err != nil {
		return db.User{}, err
	}

	existingUser, err := service.GetByEmail(ctx, params.Email)
	if err == nil && existingUser.ID != id {
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
	user, err := service.repository.Update(ctx, db.UpdateUserParams{
		ID:        id,
		Name:      params.Name,
		Email:     params.Email,
		Password:  hashedPassword,
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return db.User{}, ErrUserUpdateFailed{
			updateUser: err,
		}
	}

	if err := service.cache.Delete(ctx, "user:"+id.String()); err != nil { // delete from cache
		fmt.Println("Failed to invalidate")
	}

	if err := service.cache.DeleteByPrefix(ctx, cache.CacheKeyValidPrefixes(CacheKeyUsers)); err != nil { // delete from cache
		fmt.Println("Failed to invalidate by prefix")
	}

	err = service.events.Publish(ctx, events.QueueEventUserUpdated,
		events.UserUpdated{
			ID: user.ID,
		},
	)
	if err != nil {
		log.Printf("Failed to publish user updated event: %v", err)
		return db.User{}, err
	}

	return user, nil
}

func (service *service) Delete(ctx context.Context, id uuid.UUID) (db.User, error) {
	user, err := service.GetByID(ctx, id)
	if err != nil {
		return db.User{}, err
	}

	_, err = service.repository.Delete(ctx, id)
	if err != nil {
		return db.User{}, ErrUserDeleteFailed{
			deleteErr: err,
		}
	}

	if err := service.cache.Delete(ctx, "user:"+id.String()); err != nil { // delete from cache
		fmt.Println("Failed to invalidate single")
	}

	if err := service.cache.DeleteByPrefix(ctx, cache.CacheKeyValidPrefixes(CacheKeyUsers)); err != nil { // delete from cache
		fmt.Println("Failed to invalidate by prefix")
	}

	err = service.events.Publish(ctx, events.QueueEventUserDeleted,
		events.UserDeleted{
			ID: user.ID,
		},
	)
	if err != nil {
		log.Printf("Failed to publish user deleted event: %v", err)
		return db.User{}, err
	}

	return user, nil
}
