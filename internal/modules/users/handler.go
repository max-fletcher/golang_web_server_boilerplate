package users

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	id_helpers "github.com/max-fletcher/golang_web_server_boilerplate/helpers/ID"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/formatters"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/pagination"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/requests"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/responses"
	validator "github.com/max-fletcher/golang_web_server_boilerplate/helpers/validation"
	common_errors "github.com/max-fletcher/golang_web_server_boilerplate/internal/errors"
)

// same as the handler in internal/handler.go, but will create a new handler instance that is separate from that
type Handler struct {
	service Service // Service that belongs to this/current package by default(i.e defined in service.go)
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (handler *Handler) Create(w http.ResponseWriter, r *http.Request) error {
	params := CreateUserRequest{}
	// Passing a [pointer to params] not [params] directly, else a copy will be passed
	if err := requests.DecodeJSON(r, &params); err != nil {
		return err
	}

	if err := params.ValidateCreateUserData(); err != nil {
		return err
	}

	// 1st param: context for the request
	// 2nd param: the struct that we want to pass so it saves the underlying data in DB
	user, err := handler.service.Create(r.Context(), params)
	if err != nil {
		return err
	}

	responses.RespondWithSuccess(w, http.StatusCreated, "Created successfully", formatters.DatabaseUserToUser(user))
	return nil
}

func (handler *Handler) GetAll(w http.ResponseWriter, r *http.Request) error {
	validatedQSData, validationErrors := validator.ValidatePaginationQS(r.URL.Query())
	if len(validationErrors) > 0 {
		return common_errors.ErrValidationError{
			Errors: validationErrors,
		}
	}

	// 1st param: context for the request
	users, total, err := handler.service.GetAll(r.Context(), validatedQSData.FilterString, validatedQSData.Limit, validatedQSData.Offset)
	if err != nil {
		return err
	}
	formattedUserData := formatters.DatabaseUsersToUsers(users)
	paginatedData := pagination.GeneratePaginationFormat(validatedQSData, total, formattedUserData)

	responses.RespondWithSuccess(w, http.StatusOK, "Fetched successfully", paginatedData)
	return nil
}

func (handler *Handler) GetByID(w http.ResponseWriter, r *http.Request) error {
	id, err := id_helpers.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		return err
	}

	// 1st param: context for the request
	// 2nd param: id(type uuid) param
	user, err := handler.service.GetByID(r.Context(), id)
	if err != nil {
		return err
	}

	responses.RespondWithSuccess(w, http.StatusOK, "Fetched successfully", formatters.DatabaseUserToUser(user))
	return nil
}

func (handler *Handler) Update(w http.ResponseWriter, r *http.Request) error {
	id, err := id_helpers.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		return err
	}

	params := UpdateUserRequest{}
	// Passing a [pointer to params] not [params] directly, else a copy will be passed
	if err := requests.DecodeJSON(r, &params); err != nil {
		return err
	}

	if err := params.ValidateUpdateUserData(); err != nil {
		return err
	}

	// 1st param: context for the request
	// 2nd param: the struct that we want to pass so it saves the underlying data in DB
	user, err := handler.service.Update(r.Context(), id, params)
	if err != nil {
		return err
	}

	responses.RespondWithSuccess(w, http.StatusOK, "Updated successfully", formatters.DatabaseUserToUser(user))
	return nil
}

func (handler *Handler) Delete(w http.ResponseWriter, r *http.Request) error {
	id, err := id_helpers.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		return err
	}

	// 1st param: context for the request
	// 2nd param: id(type uuid) param
	user, err := handler.service.Delete(r.Context(), id)
	if err != nil {
		return err
	}

	responses.RespondWithSuccess(w, http.StatusOK, "Deleted successfully", formatters.DatabaseUserToUser(user))
	return nil
}
