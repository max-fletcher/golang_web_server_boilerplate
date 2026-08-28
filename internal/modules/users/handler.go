package users

import (
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
)

// same as the handler in internal/handler.go, but will create a new handler instance that is separate from that
type Handler struct {
	DB *db.Queries // Reference to a DB connection. Will be used to query data form DB.
}

func New(db *db.Queries) *Handler {
	return &Handler{
		DB: db,
	}
}
