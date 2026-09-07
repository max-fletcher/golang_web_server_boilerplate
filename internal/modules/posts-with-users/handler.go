package posts_with_users

import (
	"net/http"

	fileupload "github.com/max-fletcher/golang_web_server_boilerplate/helpers/file-upload"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/formatters"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/requests"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/responses"
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
	err := requests.ParseFormdata(r)
	if err != nil {
		return err
	}
	filenamesToStore := []string{"photo"}                         // files to get/store/get headers from request
	fileHeaders := fileupload.GetFileHeaders(r, filenamesToStore) // extracted file headers
	params := CreatePostWithUserRequest{
		Name:            r.FormValue("name"),
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirm_password"),
		Title:           r.FormValue("title"),
		Content:         r.FormValue("content"),
		UserId:          r.FormValue("user_id"),
		Photo:           fileHeaders["photo"],
	}

	createPostWithUserInput, err := params.ValidateCreatePostWithUserData()
	if err != nil {
		return err
	}

	// 1st param: context for the request
	// 2nd param: the struct that we want to pass so it saves the underlying data in DB
	user, err := handler.service.Create(r.Context(), createPostWithUserInput)
	if err != nil {
		return err
	}

	responses.RespondWithSuccess(w, http.StatusCreated, "Created successfully", formatters.DatabaseUserToUser(user))
	return nil
}

// func (handler *Handler) GetAll(w http.ResponseWriter, r *http.Request) error {
// 	validatedQSData, err := validator.ValidatePaginationQS(r.URL.Query())
// 	if err != nil {
// 		return err
// 	}

// 	// 1st param: context for the request
// 	users, total, err := handler.service.GetAll(r.Context(), validatedQSData.FilterString, validatedQSData.Limit, validatedQSData.Offset)
// 	if err != nil {
// 		return err
// 	}
// 	formattedUserData := formatters.DatabaseUsersToUsers(users)
// 	paginatedData := pagination.GeneratePaginationFormat(validatedQSData, total, formattedUserData)

// 	responses.RespondWithSuccess(w, http.StatusOK, "Fetched successfully", paginatedData)
// 	return nil
// }

// func (handler *Handler) GetByID(w http.ResponseWriter, r *http.Request) error {
// 	id, err := id_helpers.ParseUUIDRouteParam(chi.URLParam(r, "id"))
// 	if err != nil {
// 		return err
// 	}

// 	// 1st param: context for the request
// 	// 2nd param: id(type uuid) param
// 	user, err := handler.service.GetByID(r.Context(), id)
// 	if err != nil {
// 		return err
// 	}

// 	responses.RespondWithSuccess(w, http.StatusOK, "Fetched successfully", formatters.DatabaseUserToUser(user))
// 	return nil
// }

// func (handler *Handler) Update(w http.ResponseWriter, r *http.Request) error {
// 	id, err := id_helpers.ParseUUIDRouteParam(chi.URLParam(r, "id"))
// 	if err != nil {
// 		return err
// 	}

// 	params := UpdateUserRequest{}
// 	// Passing a [pointer to params] not [params] directly, else a copy will be passed
// 	if err := requests.DecodeJSON(r, &params); err != nil {
// 		return err
// 	}

// 	if err := params.ValidateUpdateUserData(); err != nil {
// 		return err
// 	}

// 	// 1st param: context for the request
// 	// 2nd param: the struct that we want to pass so it saves the underlying data in DB
// 	user, err := handler.service.Update(r.Context(), id, params)
// 	if err != nil {
// 		return err
// 	}

// 	responses.RespondWithSuccess(w, http.StatusOK, "Updated successfully", formatters.DatabaseUserToUser(user))
// 	return nil
// }

// func (handler *Handler) Delete(w http.ResponseWriter, r *http.Request) error {
// 	id, err := id_helpers.ParseUUIDRouteParam(chi.URLParam(r, "id"))
// 	if err != nil {
// 		return err
// 	}

// 	// 1st param: context for the request
// 	// 2nd param: id(type uuid) param
// 	user, err := handler.service.Delete(r.Context(), id)
// 	if err != nil {
// 		return err
// 	}

// 	responses.RespondWithSuccess(w, http.StatusOK, "Deleted successfully", formatters.DatabaseUserToUser(user))
// 	return nil
// }
