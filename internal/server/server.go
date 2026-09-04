package server

import (
	"log/slog"
	"net/http"

	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/handlers"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/logger"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/posts"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/users"
)

// The file and the function named NewServer(below) is for creating a server instance and binding dependencies to it before returning it.

// Struct containing a router instance
type Server struct {
	Router        http.Handler // Reference to router instance
	CommonHandler *handlers.Handler
	UsersHandler  *users.Handler
	PostsHandler  *posts.Handler
	Logger        *slog.Logger
}

// The name "NewServer" is a naming convention for functions that behave like a constructor. This func will create a new server.
// It is creating and passing a pointer to a server struct because remember, functions that return structs actually return copies of the struct
// and not the object itself
func NewServer(database *db.Queries) *Server {
	userRepository := users.NewRepository(database)
	userService := users.NewService(userRepository)
	userHandler := users.NewHandler(userService)

	postRepository := posts.NewRepository(database)
	postService := posts.NewService(postRepository, userService) // using DI
	postHandler := posts.NewHandler(postService)

	server := &Server{
		CommonHandler: handlers.New(database),
		UsersHandler:  userHandler,
		PostsHandler:  postHandler,
		Logger:        logger.New(),
	}

	server.Router = server.routes()

	return server
}
