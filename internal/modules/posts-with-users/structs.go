package posts_with_users

import (
	"database/sql"
	"mime/multipart"
	"time"

	"github.com/google/uuid"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	id_helpers "github.com/max-fletcher/golang_web_server_boilerplate/helpers/ID"
	validator "github.com/max-fletcher/golang_web_server_boilerplate/helpers/validation"
	common_errors "github.com/max-fletcher/golang_web_server_boilerplate/internal/errors"
)

// Struct to be validated
type CreatePostWithUserRequest struct {
	Name            string                `json:"name"`
	Email           string                `json:"email"`
	Password        string                `json:"password"`
	ConfirmPassword string                `json:"confirm_password"`
	Title           string                `json:"title"`
	Content         string                `json:"content"`
	Photo           *multipart.FileHeader `json:"photo"`
	UserId          string                `json:"user_id"`
}

type CreatePostWithUserInput struct {
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Title    string    `json:"title"`
	Content  *string   `json:"content"`
	Photo    *string   `json:"photo"`
	UserId   uuid.UUID `json:"user_id"`
}

// Rules
func (params CreatePostWithUserRequest) ValidateCreatePostWithUserData() (CreatePostWithUserInput, error) {
	err := validation.ValidateStruct(&params,
		validation.Field(
			&params.Name,
			validation.Required.Error("Name is required"),
			validation.Length(2, 100).Error("Name must be between 2 and 100 characters"),
		),
		validation.Field(
			&params.Email,
			validation.Required.Error("Email is required"),
			is.Email.Error("Email must be a valid email address"),
		),
		validation.Field(
			&params.Password,
			validation.Required.Error("Password is required"),
			validation.Length(8, 100).Error("Password must be at least 8 characters"),
		),
		validation.Field(
			&params.ConfirmPassword,
			validation.Required.Error("Confirm password is required"),
			validator.PasswordsMatch(params.Password),
		),
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
	_, exists := formattedErrors["user_id"]                           // check if err with key "user_id" exists
	userID, uuidErr := id_helpers.ParseUUID(params.UserId, "user ID") // parsing UserId field
	// if userID is not valid uuid and err with key "user_id" doesn't exist, put it in formattedErrors and set hasValidationErrors to false
	if uuidErr != nil && !exists {
		formattedErrors["user_id"] = uuidErr.Error()
		hasValidationErrors = true
	}
	if hasValidationErrors {
		return CreatePostWithUserInput{}, common_errors.ErrValidationError{
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
	createPostWithUserInput := CreatePostWithUserInput{
		Name:     params.Name,
		Email:    params.Email,
		Password: params.Password,
		Title:    params.Title,
		Content:  content,
		Photo:    photo,
		UserId:   userID,
	}

	return createPostWithUserInput, nil
}

type PostWithUser struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string
	Avatar    sql.NullString
	CreatedAt time.Time
	UpdatedAt time.Time
}
