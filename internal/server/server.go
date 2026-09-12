package server

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/max-fletcher/golang_web_server_boilerplate/internal/cache"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/config"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/events"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/handlers"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/logger"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/auth"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/posts"
	posts_with_users "github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/posts-with-users"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/users"
)

// The file and the function named NewServer(below) is for creating a server instance and binding dependencies to it before returning it.

// Struct containing a router instance
type Server struct {
	config               *config.Config
	Router               http.Handler // Reference to router instance
	CommonHandler        *handlers.Handler
	Authhandler          *auth.Handler
	UsersHandler         *users.Handler
	PostsHandler         *posts.Handler
	PostsWithUserHandler *posts_with_users.Handler
	Logger               *slog.Logger
}

// The name "NewServer" is a naming convention for functions that behave like a constructor. This func will create a new server.
// It is creating and passing a pointer to a server struct because remember, functions that return structs actually return copies of the struct
// and not the object itself
func NewServer(database *db.Queries, conn *sql.DB, cfg *config.Config, cache cache.Cache, events events.Publisher) *Server {
	// we will be sending baseUrl to handlers where we can in turn send it to LocalFileUpload for storing files with full paths to serve as static assets
	var baseUrl string

	switch cfg.AppMode {
	case "local":
		baseUrl = cfg.LocalBaseUrl
	case "production":
		baseUrl = cfg.LiveBaseUrl
	}

	// cacheClient := cache.NewRedisCache(cache.client, cfg.CacheActive) // attach redis client to cache struct
	// queueClient := queue.NewRedisQueue(redisClient)                   // attach redis client to queue struct as well

	// (using DI) initialize child structs and passing them so they can be in turn used to initialize parent structs

	userRepository := users.NewRepository(database)
	userService := users.NewService(userRepository, cache, events)
	userHandler := users.NewHandler(userService)

	authService := auth.NewService(userService, events, cfg.JWTSecret, cfg.JWTExpiry)
	authHandler := auth.NewHandler(authService)

	postRepository := posts.NewRepository(database)
	postService := posts.NewService(postRepository, userService, cache, events)
	postHandler := posts.NewHandler(postService, baseUrl)

	postWithUserRepository := posts_with_users.NewRepository(database)
	postWithUserService := posts_with_users.NewService(postWithUserRepository, conn, database, cache)
	postWithUserHandler := posts_with_users.NewHandler(postWithUserService)

	server := &Server{
		config:               cfg,
		CommonHandler:        handlers.New(database),
		Authhandler:          authHandler,
		UsersHandler:         userHandler,
		PostsHandler:         postHandler,
		PostsWithUserHandler: postWithUserHandler,
		Logger:               logger.New(),
	}

	server.Router = server.routes()

	return server
}
