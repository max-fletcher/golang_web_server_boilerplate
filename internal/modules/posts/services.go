package posts

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	fileupload "github.com/max-fletcher/golang_web_server_boilerplate/helpers/file-upload"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/formatters"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/cache"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
	common_errors "github.com/max-fletcher/golang_web_server_boilerplate/internal/errors"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/queue"
)

type UserExistenceChecker interface { // For DI. Used in validating user in structs.go
	GetByID(ctx context.Context, id uuid.UUID) (db.User, error)
}

type Service interface {
	Create(ctx context.Context, createPostInput CreatePostInput) (db.Post, error)
	GetAll(ctx context.Context, filterString string, limit int, offset int) ([]db.Post, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (db.Post, error)
	Update(ctx context.Context, id uuid.UUID, updatePostInput UpdatePostInput, baseUrl string) (db.Post, error)
	Delete(ctx context.Context, id uuid.UUID, baseUrl string) (db.Post, error)
}

type service struct {
	repository  Repository
	userChecker UserExistenceChecker // Using DI. See server/server.go where we are passing user service as 2nd param here.
	cache       cache.Cache
	cacheExpiry time.Duration
	queue       queue.Queue
}

func NewService(repository Repository, userChecker UserExistenceChecker, cache cache.Cache, cacheExpiry time.Duration, queue queue.Queue) *service {
	return &service{
		repository:  repository,
		userChecker: userChecker, // Using DI. See server/server.go where we are passing user service as 2nd param here.
		cache:       cache,
		cacheExpiry: cacheExpiry,
		queue:       queue,
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
		ID:        uuid.New(),
		Title:     createPostInput.Title,
		Content:   formatters.StringPointerToNullString(createPostInput.Content),
		Photo:     formatters.StringPointerToNullString(createPostInput.Photo),
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
	cacheKey := GetAllCacheKey(filterString, limit, offset)
	var cachedPosts GetAllPostsResult
	err := service.cache.Get(ctx, cacheKey, &cachedPosts)

	if err == nil { // return if no errros i.e fetched successfully
		fmt.Println("Cache get successful")
		return cachedPosts.Posts, cachedPosts.Total, nil
	}
	// Decide whether cache failure should fail the request. For most caches, I'd allow the request to continue.
	// else {
	// 	return err
	// }

	// 1st param: context for the request
	posts, err := service.repository.GetAll(ctx, filterString, limit, offset)
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

	fmt.Println("Cache miss")
	if err := service.cache.Set( // Set cache
		ctx,
		cacheKey,
		GetAllPostsResult{
			Posts: posts,
			Total: total,
		},
		service.cacheExpiry,
	); err != nil {
		// Failing to set shouldn't fail request since we did get data from database. Just debug why redis is not working
		log.Printf("Failed to store data in redis: %v", err)
	}

	return posts, total, nil
}

func (service *service) GetByID(ctx context.Context, id uuid.UUID) (db.Post, error) {
	cacheKey := GetByIDCacheKey(id)
	var post db.Post
	err := service.cache.Get(ctx, cacheKey, &post)

	if err == nil { // return if no errros i.e fetched successfully
		fmt.Println("Cache get successful")
		return post, nil
	}
	// Decide whether cache failure should fail the request. For most caches, I'd allow the request to continue.
	// else {
	// 	return err
	// }

	// 1st param: context for the request
	// 2nd param: id(type uuid) param
	post, err = service.repository.GetByID(ctx, id)
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

	fmt.Println("Cache miss")

	if err := service.cache.Set( // Set cache
		ctx,
		cacheKey,
		post,
		service.cacheExpiry,
	); err != nil {
		// Failing to set shouldn't fail request since we did get data from database. Just debug why redis is not working
		log.Printf("Failed to store data in redis: %v", err)
	}

	return post, nil
}

func (service *service) Update(ctx context.Context, id uuid.UUID, updatePostInput UpdatePostInput, baseUrl string) (db.Post, error) {
	existingPost, err := service.GetByID(ctx, id)
	if err != nil {
		return db.Post{}, err
	}

	// Check if old file/photo exists and replace old if new file/photo exists
	photoToUpload := existingPost.Photo
	if updatePostInput.Photo != nil {
		photoToUpload = formatters.StringPointerToNullString(updatePostInput.Photo)
	}
	// 1st param: context for the request
	// 2nd param: the struct that we want to pass so it saves the underlying data in DB
	user, err := service.repository.Update(ctx, db.UpdatePostParams{
		ID:        id,
		Title:     updatePostInput.Title,
		Content:   formatters.StringPointerToNullString(updatePostInput.Content),
		Photo:     photoToUpload,
		UserID:    updatePostInput.UserId,
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return db.Post{}, ErrPostUpdateFailed{
			updatePost: err,
		}
	}

	if err := service.cache.Delete( // delete from cache
		ctx,
		"post:"+id.String(),
	); err != nil {
		fmt.Println("Failed to invalidate")
	}

	if existingPost.Photo.Valid { // delete old file/photo if exists
		err = fileupload.DeleteFileUsingURL(existingPost.Photo.String, baseUrl)
		if err != nil {
			return db.Post{}, err
		}
	}

	return user, nil
}

func (service *service) Delete(ctx context.Context, id uuid.UUID, baseUrl string) (db.Post, error) {
	existingPost, err := service.GetByID(ctx, id)
	if err != nil {
		return db.Post{}, err
	}

	_, err = service.repository.Delete(ctx, id)
	if err != nil {
		return db.Post{}, ErrPostDeleteFailed{
			deleteErr: err,
		}
	}

	if err := service.cache.Delete(ctx, "post:"+id.String()); err != nil { // delete from cache
		fmt.Println("Failed to invalidate single")
	}

	if err := service.cache.DeleteByPrefix(ctx, cache.CacheKeyValidPrefixes(CacheKeyPosts)); err != nil { // delete from cache
		fmt.Println("Failed to invalidate by prefix")
	}

	if err := service.queue.Publish(ctx, queue.StreamPost, queue.Message{ // Publish event
		Type: queue.QueueMsgPostDeleted, // using const(sub for enum) here
		Data: map[string]any{
			"post_id": id,
		},
	}); err != nil {
		log.Printf("Failed to publish post deleted event: %v", err)
	}

	if existingPost.Photo.Valid { // delete old file/photo if exists
		err = fileupload.DeleteFileUsingURL(existingPost.Photo.String, baseUrl)
		if err != nil {
			return db.Post{}, err
		}
	}

	return existingPost, nil
}
