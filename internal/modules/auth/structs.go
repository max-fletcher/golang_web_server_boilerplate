package auth

import (
	"mime/multipart"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	validator "github.com/max-fletcher/golang_web_server_boilerplate/helpers/validation"
	common_errors "github.com/max-fletcher/golang_web_server_boilerplate/internal/errors"
)

// Struct to be validated
type UserRegistrationRequest struct {
	Name            string                `json:"name"`
	Email           string                `json:"email"`
	Password        string                `json:"password"`
	ConfirmPassword string                `json:"confirm_password"`
	Avatar          *multipart.FileHeader `json:"avatar"`
}

type UserRegistrationInput struct {
	Name            string  `json:"name"`
	Email           string  `json:"email"`
	Password        string  `json:"password"`
	ConfirmPassword string  `json:"confirm_password"`
	Avatar          *string `json:"avatar"`
}

// Rules
func (params UserRegistrationRequest) ValidateUserRegistrationData() (UserRegistrationInput, error) {
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
	)

	formattedErrors, hasValidationErrors := validator.FormatValidationErrors(err)
	if hasValidationErrors {
		return UserRegistrationInput{}, common_errors.ErrValidationError{
			Errors: formattedErrors,
		}
	}

	var avatar *string
	if params.Avatar != nil {
		avatarValue := params.Avatar.Filename
		avatar = &avatarValue
	}
	createPostInput := UserRegistrationInput{
		Name:            params.Name,
		Email:           params.Email,
		Password:        params.Password,
		ConfirmPassword: params.ConfirmPassword,
		Avatar:          avatar,
	}

	return createPostInput, nil
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (params UserLoginRequest) ValidateUserLoginData() error {
	err := validation.ValidateStruct(&params,
		validation.Field(
			&params.Email,
			is.Email.Error("Email must be a valid email address"),
		),
		validation.Field(
			&params.Password,
			validation.Required.Error("Password is required"),
			validation.Length(8, 100).Error("Password must be at least 8 characters"),
		),
	)

	formattedErrors, hasValidationErrors := validator.FormatValidationErrors(err)
	if hasValidationErrors {
		return common_errors.ErrValidationError{
			Errors: formattedErrors,
		}
	}

	return nil
}
