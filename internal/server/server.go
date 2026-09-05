package server

import (
	"log/slog"
	"net/http"

	"github.com/max-fletcher/golang_web_server_boilerplate/internal/config"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/handlers"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/logger"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/posts"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/users"
)

// The file and the function named NewServer(below) is for creating a server instance and binding dependencies to it before returning it.

// Struct containing a router instance
type Server struct {
	config        *config.Config
	Router        http.Handler // Reference to router instance
	CommonHandler *handlers.Handler
	UsersHandler  *users.Handler
	PostsHandler  *posts.Handler
	Logger        *slog.Logger
}

// The name "NewServer" is a naming convention for functions that behave like a constructor. This func will create a new server.
// It is creating and passing a pointer to a server struct because remember, functions that return structs actually return copies of the struct
// and not the object itself
func NewServer(database *db.Queries, cfg *config.Config) *Server {
	// we will be sending baseUrl to handlers where we can in turn send it to LocalFileUpload for storing files with full paths to serve as static assets
	var baseUrl string
	if cfg.AppMode == "local" {
		baseUrl = cfg.LocalBaseUrl
	} else if cfg.AppMode == "production" {
		baseUrl = cfg.LiveBaseUrl
	}

	userRepository := users.NewRepository(database)
	userService := users.NewService(userRepository)
	userHandler := users.NewHandler(userService)

	postRepository := posts.NewRepository(database)
	postService := posts.NewService(postRepository, userService) // using DI
	postHandler := posts.NewHandler(postService, baseUrl)

	server := &Server{
		config:        cfg,
		CommonHandler: handlers.New(database),
		UsersHandler:  userHandler,
		PostsHandler:  postHandler,
		Logger:        logger.New(),
	}

	server.Router = server.routes()

	return server
}
