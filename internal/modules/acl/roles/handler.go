package roles

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	id_helpers "github.com/max-fletcher/golang_web_server_boilerplate/helpers/ID"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/formatters"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/pagination"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/requests"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/responses"
	validator "github.com/max-fletcher/golang_web_server_boilerplate/helpers/validation"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/auth"
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
	params := CreateRoleRequest{}
	// Passing a [pointer to params] not [params] directly, else a copy will be passed
	if err := requests.DecodeJSON(r, &params); err != nil {
		return err
	}

	createRoleInput, err := params.ValidateCreateRoleData()
	if err != nil {
		return err
	}

	// 1st param: context for the request
	// 2nd param: the struct that we want to pass so it saves the underlying data in DB
	role, err := handler.service.Create(r.Context(), createRoleInput)
	if err != nil {
		return err
	}

	responses.RespondWithSuccess(w, http.StatusCreated, "Created successfully", formatters.DatabaseRoleToRole(role))
	return nil
}

func (handler *Handler) GetAll(w http.ResponseWriter, r *http.Request) error {
	userID, ok := auth.UserIDFromContext(r.Context())
	fmt.Println("Auth user ID:", userID)
	if !ok {
		// This should normally never happen because middleware protects the route.
		return errors.New("User ID not found")
	}

	validatedQSData, err := validator.ValidatePaginationQS(r.URL.Query())
	if err != nil {
		return err
	}

	// 1st param: context for the request
	roles, total, err := handler.service.GetAll(r.Context(), validatedQSData.FilterString, validatedQSData.Limit, validatedQSData.Offset)
	if err != nil {
		return err
	}
	formattedRoleData := formatters.DatabaseRolesToRoles(roles)
	paginatedData := pagination.GeneratePaginationFormat(validatedQSData, total, formattedRoleData)

	responses.RespondWithSuccess(w, http.StatusOK, "Fetched successfully", paginatedData)
	return nil
}

func (handler *Handler) GetByID(w http.ResponseWriter, r *http.Request) error {
	id, err := id_helpers.ParseUUID(chi.URLParam(r, "id"), "role ID")
	if err != nil {
		return err
	}

	// 1st param: context for the request
	// 2nd param: id(type uuid) param
	role, err := handler.service.GetByID(r.Context(), id)
	if err != nil {
		return err
	}

	responses.RespondWithSuccess(w, http.StatusOK, "Fetched successfully", formatters.DatabaseRoleToRole(role))
	return nil
}

func (handler *Handler) Update(w http.ResponseWriter, r *http.Request) error {
	id, err := id_helpers.ParseUUID(chi.URLParam(r, "id"), "role ID")
	if err != nil {
		return err
	}

	params := UpdateRoleRequest{}
	// Passing a [pointer to params] not [params] directly, else a copy will be passed
	if err := requests.DecodeJSON(r, &params); err != nil {
		return err
	}

	updateRoleInput, err := params.ValidateUpdateRoleData()
	if err != nil {
		return err
	}

	// 1st param: context for the request
	// 2nd param: the struct that we want to pass so it saves the underlying data in DB
	role, err := handler.service.Update(r.Context(), id, updateRoleInput)
	if err != nil {
		return err
	}

	responses.RespondWithSuccess(w, http.StatusOK, "Updated successfully", formatters.DatabaseRoleToRole(role))
	return nil
}

func (handler *Handler) Delete(w http.ResponseWriter, r *http.Request) error {
	id, err := id_helpers.ParseUUID(chi.URLParam(r, "id"), "role ID")
	if err != nil {
		return err
	}

	// 1st param: context for the request
	// 2nd param: id(type uuid) param
	role, err := handler.service.Delete(r.Context(), id)
	if err != nil {
		return err
	}

	responses.RespondWithSuccess(w, http.StatusOK, "Deleted successfully", formatters.DatabaseRoleToRole(role))
	return nil
}
