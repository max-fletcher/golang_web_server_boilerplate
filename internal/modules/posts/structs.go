package posts

import (
	"mime/multipart"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	id_helpers "github.com/max-fletcher/golang_web_server_boilerplate/helpers/ID"
	validator "github.com/max-fletcher/golang_web_server_boilerplate/helpers/validation"
	common_errors "github.com/max-fletcher/golang_web_server_boilerplate/internal/errors"
)

// Struct to be validated
type CreatePostRequest struct {
	Title   string
	Content string
	Photo   *multipart.FileHeader
	UserId  string
}

type CreatePostInput struct {
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Photo   *string   `json:"photo"`
	UserId  uuid.UUID `json:"user_id"`
}

// Rules
func (params CreatePostRequest) ValidateCreatePostData() (CreatePostInput, error) {
	err := validation.ValidateStruct(&params,
		validation.Field(
			&params.Title,
			validation.Required.Error("Title is required"),
			validation.Length(2, 100).Error("Title must be between 2 and 100 characters"),
		),
		validation.Field(
			&params.Content,
			validation.Length(2, 100).Error("Content must be between 10 and 200 characters"),
		),
		validation.Field(
			&params.UserId,
			validation.Required.Error("User is required"),
			is.UUID.Error("Not a valid UUID"),
		),
		validation.Field( // Validate photo
			&params.Photo,
			// validation.Required.Error("Photo is required"), // if you want file(photo in this case) to be required
			validation.By(validator.ValidatePhoto),
		),
	)

	// If you want createPostInput.Photo to be a string or nil(Hack for if you want a "string or nil" value for any struct field etc.)
	var photo *string
	if params.Photo != nil {
		photoValue := params.Photo.Filename
		photo = &photoValue
	}

	// parsing UserId field
	userID, uuidErr := id_helpers.ParseUUID(params.UserId)
	createPostInput := CreatePostInput{
		Title:   params.Title,
		Content: params.Content,
		Photo:   photo, // keeping this empty since we will be populating this later when needed
		UserId:  userID,
	}

	formattedErrors, ok := validator.FormatValidationErrors(err)
	if uuidErr != nil { // if userID is not valid uuid, put it in formattedErrors and set ok to false(ok == false means validation errors exists)
		formattedErrors["user_id"] = uuidErr.Error()
		ok = false
	}
	if !ok {
		return createPostInput, nil
	}

	return CreatePostInput{}, common_errors.ErrValidationError{
		Errors: formattedErrors,
	}
}

type UpdatePostRequest struct {
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	Photo   string `json:"photo,omitempty"`
	UserId  string `json:"user_id,omitempty"`
}

type UpdatePostInput struct {
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Photo   string    `json:"photo,omitempty"`
	UserId  uuid.UUID `json:"user_id"`
}

func (params UpdatePostRequest) ValidateUpdatePostData() (UpdatePostInput, error) {
	err := validation.ValidateStruct(&params,
		validation.Field(
			&params.Title,
			validation.Length(2, 100).Error("Title must be between 2 and 100 characters"),
		),
		validation.Field(
			&params.Content,
			validation.Length(2, 100).Error("Content must be between 10 and 200 characters"),
		),
		validation.Field(
			&params.UserId,
			is.UUID.Error("Not a valid UUID"),
		),
	)

	// parsing UserId field
	userID, uuidErr := id_helpers.ParseUUID(params.UserId)
	updatePostInput := UpdatePostInput{
		Title:   params.Title,
		Content: params.Content,
		Photo:   params.Photo,
		UserId:  userID,
	}

	formattedErrors, ok := validator.FormatValidationErrors(err)
	if uuidErr != nil { // if userID is not valid uuid, put it in formattedErrors and set ok to true(ok == true means validation errors exists)
		formattedErrors["user_id"] = uuidErr.Error()
		ok = true
	}
	if !ok {
		return updatePostInput, nil
	}

	return UpdatePostInput{}, common_errors.ErrValidationError{
		Errors: formattedErrors,
	}
}
