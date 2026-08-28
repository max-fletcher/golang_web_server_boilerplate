package handlers

import "github.com/max-fletcher/golang_web_server_boilerplate/internal/db"

type Handler struct {
	DB *db.Queries // Reference to a DB connection. Will be used to query data form DB.
}

// This function is named New. It is a naming convention for functions that behave like a constructor. This func will create a new server.
// It is creating and passing a pointer to a server struct because remember, functions that return structs actually return copies of the struct
// and not the object itself
func New(db *db.Queries) *Handler {
	return &Handler{
		DB: db,
	}
}
