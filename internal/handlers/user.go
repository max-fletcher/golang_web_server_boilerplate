package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/crypto"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/formatters"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/responses"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
)

func (handler *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body) // decode request body
	params := parameters{}
	err := decoder.Decode(&params) // decode request body(json) and store in params variable, or get an error
	if err != nil {
		responses.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	hashedPassword, err := crypto.HashPassword(params.Password)
	if err != nil {
		responses.RespondWithError(w, 400, fmt.Sprintf("Error hashing password: %v", err))
		return
	}

	// 1st param: context for the request
	// 2nd param: the struct that we want to pass so it saves the underlying data in DB
	user, err := handler.DB.CreateUser(r.Context(), db.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		Email:     params.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		responses.RespondWithError(w, 400, fmt.Sprintf("Couldn't create user: %v", err))
		return
	}

	responses.RespondWithJSON(w, 201, responses.SuccessResponse{
		Code:   201,
		Status: "ok",
		Data:   formatters.DatabaseUserToUser(user),
	})
}

// handler method that handles fetching all users. The addition of (handler *Handler)
// turns it into a method for handler.
func (handler *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// 1st param: context for the request
	users, err := handler.DB.GetUsers(r.Context())
	if err != nil {
		responses.RespondWithError(w, 400, fmt.Sprintf("Failed to fetch users: %v", err))
		return
	}

	responses.RespondWithJSON(w, 200, responses.SuccessResponse{
		Code:   200,
		Status: "ok",
		Data:   formatters.DatabaseUsersToUsers(users),
	})
}
