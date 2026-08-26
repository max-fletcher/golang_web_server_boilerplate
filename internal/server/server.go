package server

import (
	"net/http"

	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
)

// Struct containing a router instance
type Server struct {
	DB     *db.Queries  // Reference to a DB connection. Will be used to query data form DB.
	Router http.Handler // Reference to router instance
}

// This function is named New. It is a naming convention for functions that behave like a constructor. This func will create a new server
func New(db *db.Queries) *Server {
	server := &Server{
		DB: db,
	}

	server.Router = server.routes()

	return server
}
