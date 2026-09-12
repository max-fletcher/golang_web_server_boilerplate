package auth

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/formatters"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/requests"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/responses"
)

type AuthenticatedUser struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Avatar    sql.NullString
	CreatedAt time.Time
	UpdatedAt time.Time
}

// same as the handler in internal/handler.go, but will create a new handler instance that is separate from that
type Handler struct {
	service Service // Service that belongs to this/current package by default(i.e defined in service.go)
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (handler *Handler) UserRegistration(w http.ResponseWriter, r *http.Request) error {
	err := requests.ParseFormdata(r)
	if err != nil {
		return err
	}
	params := UserRegistrationRequest{
		Name:            r.FormValue("name"),
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirm_password"),
	}

	registerUserInput, err := params.ValidateUserRegistrationData()
	if err != nil {
		return err
	}

	// 1st param: context for the request
	// 2nd param: the struct that we want to pass so it saves the underlying data in DB
	authUser, jwt, err := handler.service.UserRegistration(r.Context(), registerUserInput)
	if err != nil {
		return err
	}

	responses.RespondWithSuccess(
		w,
		http.StatusCreated,
		"Registered successfully",
		formatters.ToAuthenticatedUserWithJWT(jwt, formatters.AuthenticatedUser(authUser)))
	return nil
}

func (handler *Handler) UserLogin(w http.ResponseWriter, r *http.Request) error {
	err := requests.ParseFormdata(r)
	if err != nil {
		return err
	}
	params := UserLoginRequest{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	err = params.ValidateUserLoginData()
	if err != nil {
		return err
	}

	// 1st param: context for the request
	// 2nd param: the struct that we want to pass so it saves the underlying data in DB
	authUser, jwt, err := handler.service.UserLogin(r.Context(), params)
	if err != nil {
		return err
	}

	responses.RespondWithSuccess(
		w,
		http.StatusCreated,
		"Login successfully",
		formatters.ToAuthenticatedUserWithJWT(jwt, formatters.AuthenticatedUser(authUser)))
	return nil
}
