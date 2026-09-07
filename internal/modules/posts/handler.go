package posts

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	id_helpers "github.com/max-fletcher/golang_web_server_boilerplate/helpers/ID"
	fileupload "github.com/max-fletcher/golang_web_server_boilerplate/helpers/file-upload"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/formatters"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/pagination"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/requests"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/responses"
	validator "github.com/max-fletcher/golang_web_server_boilerplate/helpers/validation"
)

// same as the handler in internal/handler.go, but will create a new handler instance that is separate from that
type Handler struct {
	service Service // Service that belongs to this/current package by default(i.e defined in service.go)
	baseUrl string
}

func NewHandler(service Service, baseUrl string) *Handler {
	return &Handler{
		service: service,
		baseUrl: baseUrl,
	}
}

func (handler *Handler) Create(w http.ResponseWriter, r *http.Request) error {
	// Removed since we are not decoding json and using formdata instead
	// params := CreatePostRequest{}
	// // Passing a [pointer to params] not [params] directly, else a copy will be passed
	// if err := requests.DecodeJSON(r, &params); err != nil {
	// 	return err
	// }
	err := requests.ParseFormdata(r)
	if err != nil {
		return err
	}
	filenamesToStore := []string{"photo"}                         // files to get/store/get headers from request
	fileHeaders := fileupload.GetFileHeaders(r, filenamesToStore) // extracted file headers
	params := CreatePostRequest{
		Title:   r.FormValue("title"),
		Content: r.FormValue("content"),
		UserId:  r.FormValue("user_id"),
		Photo:   fileHeaders["photo"],
	}

	createPostInput, err := params.ValidateCreatePostData()
	if err != nil {
		return err
	}

	destination := "posts"
	storedFiles, err := fileupload.LocalFileUploader(r, filenamesToStore, destination, handler.baseUrl, true)
	if err != nil {
		return err
	}

	// Since CreatePostInput's Photo field is now a pointer to a string, we are using the block below
	// instead of createPostInput.Photo = &storedFiles["photo"]
	if photoURL, ok := storedFiles["photo"]; ok {
		createPostInput.Photo = &photoURL
	}

	// 1st param: context for the request
	// 2nd param: the struct that we want to pass so it saves the underlying data in DB
	post, err := handler.service.Create(r.Context(), createPostInput)
	if err != nil {
		return err
	}

	responses.RespondWithSuccess(w, http.StatusCreated, "Created successfully", formatters.DatabasePostToPost(post))
	return nil
}

func (handler *Handler) GetAll(w http.ResponseWriter, r *http.Request) error {
	validatedQSData, err := validator.ValidatePaginationQS(r.URL.Query())
	if err != nil {
		return err
	}

	// 1st param: context for the request
	posts, total, err := handler.service.GetAll(r.Context(), validatedQSData.FilterString, validatedQSData.Limit, validatedQSData.Offset)
	if err != nil {
		return err
	}
	formattedPostData := formatters.DatabasePostsToPosts(posts)
	paginatedData := pagination.GeneratePaginationFormat(validatedQSData, total, formattedPostData)

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
	post, err := handler.service.GetByID(r.Context(), id)
	if err != nil {
		return err
	}

	responses.RespondWithSuccess(w, http.StatusOK, "Fetched successfully", formatters.DatabasePostToPost(post))
	return nil
}

func (handler *Handler) Update(w http.ResponseWriter, r *http.Request) error {
	id, err := id_helpers.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		return err
	}

	// Removed since we are not decoding json and using formdata instead
	// params := UpdatePostRequest{}
	// // Passing a [pointer to params] not [params] directly, else a copy will be passed
	// if err := requests.DecodeJSON(r, &params); err != nil {
	// 	return err
	// }
	err = requests.ParseFormdata(r)
	if err != nil {
		return err
	}
	filenamesToStore := []string{"photo"}                         // filenames to get/store from request
	fileHeaders := fileupload.GetFileHeaders(r, filenamesToStore) // extracted file headers
	params := UpdatePostRequest{
		Title:   r.FormValue("title"),
		Content: r.FormValue("content"),
		UserId:  r.FormValue("user_id"),
		Photo:   fileHeaders["photo"],
	}

	updatePostInput, err := params.ValidateUpdatePostData()
	if err != nil {
		return err
	}

	destination := "posts"
	storedFiles, err := fileupload.LocalFileUploader(r, filenamesToStore, destination, handler.baseUrl, true)
	if err != nil {
		return err
	}

	// Since UpdatePostInput's Photo field is now a pointer to a string, we are using the block below
	// instead of UpdatePostInput.Photo = &storedFiles["photo"]
	if photoURL, ok := storedFiles["photo"]; ok {
		updatePostInput.Photo = &photoURL
	}

	// 1st param: context for the request
	// 2nd param: the struct that we want to pass so it saves the underlying data in DB
	post, err := handler.service.Update(r.Context(), id, updatePostInput, handler.baseUrl)
	if err != nil {
		return err
	}

	responses.RespondWithSuccess(w, http.StatusOK, "Updated successfully", formatters.DatabasePostToPost(post))
	return nil
}

func (handler *Handler) Delete(w http.ResponseWriter, r *http.Request) error {
	id, err := id_helpers.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		return err
	}

	// 1st param: context for the request
	// 2nd param: id(type uuid) param
	post, err := handler.service.Delete(r.Context(), id, handler.baseUrl)
	if err != nil {
		return err
	}

	responses.RespondWithSuccess(w, http.StatusOK, "Deleted successfully", formatters.DatabasePostToPost(post))
	return nil
}
