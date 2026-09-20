package modules

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
	params := CreateModuleRequest{}
	// Passing a [pointer to params] not [params] directly, else a copy will be passed
	if err := requests.DecodeJSON(r, &params); err != nil {
		return err
	}

	createModuleInput, err := params.ValidateCreateModuleData()
	if err != nil {
		return err
	}

	// 1st param: context for the request
	// 2nd param: the struct that we want to pass so it saves the underlying data in DB
	module, err := handler.service.Create(r.Context(), createModuleInput)
	if err != nil {
		return err
	}

	responses.RespondWithSuccess(w, http.StatusCreated, "Created successfully", formatters.DatabaseModuleToModule(module))
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
	modules, total, err := handler.service.GetAll(r.Context(), validatedQSData.FilterString, validatedQSData.Limit, validatedQSData.Offset)
	if err != nil {
		return err
	}
	formattedModuleData := formatters.DatabaseModulesToModules(modules)
	paginatedData := pagination.GeneratePaginationFormat(validatedQSData, total, formattedModuleData)

	responses.RespondWithSuccess(w, http.StatusOK, "Fetched successfully", paginatedData)
	return nil
}

func (handler *Handler) GetByID(w http.ResponseWriter, r *http.Request) error {
	id, err := id_helpers.ParseUUID(chi.URLParam(r, "id"), "module ID")
	if err != nil {
		return err
	}

	// 1st param: context for the request
	// 2nd param: id(type uuid) param
	module, err := handler.service.GetByID(r.Context(), id)
	if err != nil {
		return err
	}

	responses.RespondWithSuccess(w, http.StatusOK, "Fetched successfully", formatters.DatabaseModuleToModule(module))
	return nil
}
