package roles

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/cache"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
	common_errors "github.com/max-fletcher/golang_web_server_boilerplate/internal/errors"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/events"
)

type Service interface {
	Create(ctx context.Context, createRoleInput CreateRoleInput) (db.Role, error)
	GetAll(ctx context.Context, filterString string, limit int, offset int) ([]db.Role, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (db.Role, error)
	Update(ctx context.Context, id uuid.UUID, updateRoleInput UpdateRoleInput) (db.Role, error)
	Delete(ctx context.Context, id uuid.UUID) (db.Role, error)
}

type service struct {
	repository Repository
	cache      cache.Cache
	events     events.Publisher
}

func NewService(
	repository Repository,
	cache cache.Cache,
	events events.Publisher,
) *service {
	return &service{
		repository: repository,
		cache:      cache,
		events:     events,
	}
}

func (service *service) Create(ctx context.Context, createRoleInput CreateRoleInput) (db.Role, error) {
	// 1st param: context for the request
	// 2nd param: the struct that we want to pass so it saves the underlying data in DB
	role, err := service.repository.Create(ctx, db.CreateRoleParams{
		ID:        uuid.New(),
		Name:      createRoleInput.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return db.Role{}, ErrRoleCreateFailed{
			CreateErr: err,
		}
	}

	if err := service.cache.DeleteByPrefix(ctx, cache.CacheKeyValidPrefixes(CacheKeyRoles)); err != nil { // delete from cache
		fmt.Println("Failed to invalidate by prefix")
	}

	err = service.events.Publish(ctx, events.QueueEventRoleCreated,
		events.RoleCreated{
			ID: role.ID,
		},
	)
	if err != nil {
		log.Printf("Failed to publish role created event: %v", err)
		return db.Role{}, err
	}

	return role, nil
}

func (service *service) GetAll(ctx context.Context, filterString string, limit int, offset int) ([]db.Role, int, error) {
	cacheKey := GetAllCacheKey(filterString, limit, offset)
	var cachedRoles GetAllRolesResult
	err := service.cache.Get(ctx, cacheKey, &cachedRoles)
	if err == nil { // return if no errros i.e fetched successfully
		fmt.Println("Cache get successful")
		return cachedRoles.Roles, cachedRoles.Total, nil
	}
	// Decide whether cache failure should fail the request. For most caches, I'd allow the request to continue.
	// else {
	// 	return err
	// }

	// 1st param: context for the request
	roles, err := service.repository.GetAll(ctx, filterString, limit, offset)
	if err != nil {
		return []db.Role{}, 0, ErrRolesFetchFailed{
			FetchErr: err,
		}
	}

	total, err := service.repository.GetAllCount(ctx, filterString)
	if err != nil {
		var bigInt64ToIntError common_errors.ErrBigInt64ToIntError
		if errors.As(err, &bigInt64ToIntError) {
			return []db.Role{}, 0, fmt.Errorf("Limit and/or offset value out of range")
		}

		return []db.Role{}, 0, ErrRolesFetchFailed{
			FetchErr: err,
		}
	}

	fmt.Println("Cache miss")
	if err := service.cache.Set( // Set cache
		ctx,
		cacheKey,
		GetAllRolesResult{
			Roles: roles,
			Total: total,
		},
	); err != nil {
		// Failing to set shouldn't fail request since we did get data from database. Just debug why redis is not working
		log.Printf("Failed to store data in redis: %v", err)
	}

	return roles, total, nil
}

func (service *service) GetByID(ctx context.Context, id uuid.UUID) (db.Role, error) {
	cacheKey := GetByIDCacheKey(id)
	var role db.Role
	err := service.cache.Get(ctx, cacheKey, &role)
	if err == nil { // return if no errros i.e fetched successfully
		fmt.Println("Cache get successful")
		return role, nil
	}
	// Decide whether cache failure should fail the request. For most caches, I'd allow the request to continue.
	// else {
	// 	return err
	// }

	// 1st param: context for the request
	// 2nd param: id(type uuid) param
	role, err = service.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // check if error is of type sql.ErrNoRows
			return db.Role{}, ErrRoleWithIdNotFound{
				ID: id,
			}
		}

		return db.Role{}, ErrRoleFetchFailed{
			FetchErr: err,
		}
	}

	fmt.Println("Cache miss")
	if err := service.cache.Set(ctx, cacheKey, role); err != nil { // Set cache
		// Failing to set shouldn't fail request since we did get data from database. Just debug why redis is not working
		log.Printf("Failed to store data in redis: %v", err)
	}

	return role, nil
}

func (service *service) Update(ctx context.Context, id uuid.UUID, updateRoleInput UpdateRoleInput) (db.Role, error) {
	existingRole, err := service.GetByID(ctx, id)
	if err != nil {
		return db.Role{}, err
	}

	// 1st param: context for the request
	// 2nd param: the struct that we want to pass so it saves the underlying data in DB
	role, err := service.repository.Update(ctx, db.UpdateRoleParams{
		ID:        id,
		Name:      updateRoleInput.Name,
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return db.Role{}, ErrRoleUpdateFailed{
			UpdateRole: err,
		}
	}

	if err := service.cache.Delete(ctx, "role:"+id.String()); err != nil { // delete from cache
		fmt.Println("Failed to invalidate")
	}

	if err := service.cache.DeleteByPrefix(ctx, cache.CacheKeyValidPrefixes(CacheKeyRoles)); err != nil { // delete from cache
		fmt.Println("Failed to invalidate by prefix")
	}

	err = service.events.Publish(ctx, events.QueueEventRoleUpdated,
		events.RoleUpdated{
			ID: existingRole.ID,
		},
	)
	if err != nil {
		log.Printf("Failed to publish role updated event: %v", err)
		return db.Role{}, err
	}

	return role, nil
}

func (service *service) Delete(ctx context.Context, id uuid.UUID) (db.Role, error) {
	existingRole, err := service.GetByID(ctx, id)
	if err != nil {
		return db.Role{}, err
	}

	_, err = service.repository.Delete(ctx, id)
	if err != nil {
		return db.Role{}, ErrRoleDeleteFailed{
			DeleteErr: err,
		}
	}

	if err := service.cache.Delete(ctx, "role:"+id.String()); err != nil { // delete from cache
		fmt.Println("Failed to invalidate single")
	}

	if err := service.cache.DeleteByPrefix(ctx, cache.CacheKeyValidPrefixes(CacheKeyRoles)); err != nil { // delete from cache
		fmt.Println("Failed to invalidate by prefix")
	}

	err = service.events.Publish(ctx, events.QueueEventRoleDeleted,
		events.RoleDeleted{
			ID: existingRole.ID,
		},
	)
	if err != nil {
		log.Printf("Failed to publish role deleted event: %v", err)
		return db.Role{}, err
	}

	return existingRole, nil
}
