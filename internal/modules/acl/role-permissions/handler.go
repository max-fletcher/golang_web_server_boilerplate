package role_permissions

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	id_helpers "github.com/max-fletcher/golang_web_server_boilerplate/helpers/ID"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/formatters"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/pagination"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/requests"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/responses"
	validator "github.com/max-fletcher/golang_web_server_boilerplate/helpers/validation"
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
	params := CreateRolePermissionRequest{}
	// Passing a [pointer to params] not [params] directly, else a copy will be passed
	if err := requests.DecodeJSON(r, &params); err != nil {
		return err
	}

	createRolePermissionInput, err := params.ValidateCreateRolePermissionData()
	if err != nil {
		return err
	}

	// 1st param: context for the request
	// 2nd param: the struct that we want to pass so it saves the underlying data in DB
	rolePermission, err := handler.service.Create(r.Context(), createRolePermissionInput)
	if err != nil {
		return err
	}

	formattedRolePermission := formatters.DatabaseGetRolePermissionByIDRowToRolePermissionNoTimestamp(rolePermission)
	responses.RespondWithSuccess(w, http.StatusCreated, "Role assigned to permission successfully.", formattedRolePermission)
	return nil
}

func (handler *Handler) GetAll(w http.ResponseWriter, r *http.Request) error {
	validatedQSData, err := validator.ValidatePaginationQS(r.URL.Query())
	if err != nil {
		return err
	}

	// 1st param: context for the request
	rolePermissions, total, err := handler.service.GetAll(r.Context(), validatedQSData.FilterString, validatedQSData.Limit, validatedQSData.Offset)
	if err != nil {
		return err
	}

	formattedRolePermissionData := formatters.DatabaseRolePermissionsToRolePermissions(rolePermissions)
	paginatedData := pagination.GeneratePaginationFormat(validatedQSData, total, formattedRolePermissionData)
	responses.RespondWithSuccess(w, http.StatusOK, "Fetched successfully", paginatedData)
	return nil
}

func (handler *Handler) GetByID(w http.ResponseWriter, r *http.Request) error {
	id, err := id_helpers.ParseUUID(chi.URLParam(r, "id"), "id")
	if err != nil {
		return err
	}

	// 1st param: context for the request
	// 2nd param: id(type uuid) param
	rolePermission, err := handler.service.GetByID(r.Context(), id)
	if err != nil {
		return err
	}

	formattedRolePermission := formatters.DatabaseGetRolePermissionByIDRowToRolePermissionNoTimestamp(rolePermission)
	responses.RespondWithSuccess(w, http.StatusOK, "Fetched successfully", formattedRolePermission)
	return nil
}

func (handler *Handler) GetAllUsersWithRolePermissions(w http.ResponseWriter, r *http.Request) error {
	validatedQSData, err := validator.ValidatePaginationQS(r.URL.Query())
	if err != nil {
		return err
	}

	// 1st param: context for the request
	rolePermissions, total, err := handler.service.GetAllUsersWithRolePermissions(r.Context(), validatedQSData.FilterString, validatedQSData.Limit, validatedQSData.Offset)
	if err != nil {
		return err
	}

	formattedRolePermissionData := formatters.DatabaseUserWRolePermissionsToUserWRolePermissions(rolePermissions)
	paginatedData := pagination.GeneratePaginationFormat(validatedQSData, total, formattedRolePermissionData)
	responses.RespondWithSuccess(w, http.StatusOK, "Fetched successfully", paginatedData)
	return nil
}

func (handler *Handler) GetByUserID(w http.ResponseWriter, r *http.Request) error {
	userID, err := id_helpers.ParseUUID(chi.URLParam(r, "userID"), "user ID")
	if err != nil {
		return err
	}

	// 1st param: context for the request
	// 2nd param: id(type uuid) param
	rolePermission, err := handler.service.GetByUserID(r.Context(), userID)
	if err != nil {
		return err
	}

	// #TODO: SEE HOW TO FORMAT DATA INSIDE SERVICE AND HERE(1 LINE BELOW), AND SEND IT BACK AS RESPONSE
	fmt.Printf("ACL Handler. Data: %v \n", rolePermission)
	formattedRolePermission := formatters.DatabaseUserWRolePermissionToUserWRolePermission((rolePermission))
	responses.RespondWithSuccess(w, http.StatusOK, "Fetched successfully", formattedRolePermission)
	return nil
}

func (handler *Handler) Delete(w http.ResponseWriter, r *http.Request) error {
	id, err := id_helpers.ParseUUID(chi.URLParam(r, "id"), "rolePermission ID")
	if err != nil {
		return err
	}

	// 1st param: context for the request
	// 2nd param: id(type uuid) param
	rolePermission, err := handler.service.Delete(r.Context(), id)
	if err != nil {
		return err
	}

	formattedRolePermission := formatters.DatabaseGetRolePermissionByIDRowToRolePermissionNoTimestamp(rolePermission)
	responses.RespondWithSuccess(w, http.StatusOK, "Deleted successfully", formattedRolePermission)
	return nil
}

func (handler *Handler) DeleteByRoleIDAndPermissionID(w http.ResponseWriter, r *http.Request) error {
	roleID, err := id_helpers.ParseUUID(chi.URLParam(r, "roleID"), "role ID")
	if err != nil {
		return err
	}

	permissionID, err := id_helpers.ParseUUID(chi.URLParam(r, "permissionID"), "permission ID")
	if err != nil {
		return err
	}

	// 1st param: context for the request
	// 2nd param: id(type uuid) param
	rolePermission, err := handler.service.DeleteByRoleIDAndPermissionID(r.Context(), roleID, permissionID)
	if err != nil {
		return err
	}

	formattedRolePermission := formatters.DatabaseGetRolePermissionByRoleIDAndPermissionIDRowToRolePermission(rolePermission)
	responses.RespondWithSuccess(w, http.StatusOK, "Deleted successfully", formattedRolePermission)
	return nil
}
