package server

import (
	"net/http"

	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/handlers"
)

// Struct containing a router instance
type Server struct {
	Router  http.Handler // Reference to router instance
	Handler *handlers.Handler
}

// This function is named New. It is a naming convention for functions that behave like a constructor. This func will create a new server.
// It is creating and passing a pointer to a server struct because remember, functions that return structs actually return copies of the struct
// and not the object itself
func New(db *db.Queries) *Server {
	handler := handlers.New(db)

	server := &Server{
		Handler: handler,
	}

	server.Router = server.routes()

	return server
}
