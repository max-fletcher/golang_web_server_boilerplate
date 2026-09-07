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
	Title   string                `json:"title"`
	Content string                `json:"content"`
	Photo   *multipart.FileHeader `json:"photo"`
	UserId  string                `json:"user_id"`
}

type CreatePostInput struct {
	Title   string    `json:"title"`
	Content *string   `json:"content"`
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

	formattedErrors, hasValidationErrors := validator.FormatValidationErrors(err)
	userID, uuidErr := id_helpers.ParseUUID(params.UserId) // parsing UserId field
	if uuidErr != nil {                                    // if userID is not valid uuid, put it in formattedErrors and set hasValidationErrors to false
		formattedErrors["user_id"] = uuidErr.Error()
		hasValidationErrors = true
	}
	if hasValidationErrors {
		return CreatePostInput{}, common_errors.ErrValidationError{
			Errors: formattedErrors,
		}
	}
	// Construct an instance of createPostInput
	// *IMPORTANT: This is how you dead with nullable fields.(Sources: posts/structs.go(especially CreatePostRequest and CreatePostInput),
	// posts/services.go and posts/formatters.go)
	// If you want createPostInput.Photo to be a string or nil(Hack for if you want a "string or nil" value for any struct field etc.)
	var photo *string
	if params.Photo != nil {
		photoValue := params.Photo.Filename
		photo = &photoValue
	}
	// Same as above
	var content *string
	if params.Content != "" {
		content = &params.Content
	}
	createPostInput := CreatePostInput{
		Title:   params.Title,
		Content: content,
		Photo:   photo,
		UserId:  userID,
	}

	return createPostInput, nil

}

type UpdatePostRequest struct {
	Title   string
	Content string
	Photo   *multipart.FileHeader
	UserId  string
}

type UpdatePostInput struct {
	Title   string    `json:"title"`
	Content *string   `json:"content"`
	Photo   *string   `json:"photo"`
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
		validation.Field( // Validate photo
			&params.Photo,
			// validation.Required.Error("Photo is required"), // if you want file(photo in this case) to be required
			validation.By(validator.ValidatePhoto),
		),
	)

	formattedErrors, hasValidationErrors := validator.FormatValidationErrors(err)
	userID, uuidErr := id_helpers.ParseUUID(params.UserId) // parsing UserId field
	if uuidErr != nil {                                    // if userID is not valid uuid, put it in formattedErrors and set ok to false(ok == false means validation errors exists)
		formattedErrors["user_id"] = uuidErr.Error()
		hasValidationErrors = true
	}
	if hasValidationErrors {

		return UpdatePostInput{}, common_errors.ErrValidationError{
			Errors: formattedErrors,
		}
	}
	// Construct an instance of updatePostInput

	// *IMPORTANT: This is how you dead with nullable fields.(Sources: posts/structs.go(especially CreatePostRequest and CreatePostInput),
	// posts/services.go and posts/formatters.go)
	// If you want createPostInput.Photo to be a string or nil(Hack for if you want a "string or nil" value for any struct field etc.)
	var photo *string
	if params.Photo != nil {
		photoValue := params.Photo.Filename
		photo = &photoValue
	}
	// Same as above
	var content *string
	if params.Content != "" {
		content = &params.Content
	}

	updatePostInput := UpdatePostInput{
		Title:   params.Title,
		Content: content,
		Photo:   photo,
		UserId:  userID,
	}

	return updatePostInput, nil
}
