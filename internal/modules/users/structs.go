package users

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	validator "github.com/max-fletcher/golang_web_server_boilerplate/helpers/validation"
	common_errors "github.com/max-fletcher/golang_web_server_boilerplate/internal/errors"
)

// Struct to be validated
type CreateUserRequest struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

// Rules
func (params CreateUserRequest) ValidateCreateUserData() error {
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
			passwordsMatch(params.Password),
		),
	)

	formattedErrors, ok := validator.FormatValidationErrors(err)
	if !ok {
		return nil
	}

	return common_errors.ErrValidationError{
		Errors: formattedErrors,
	}
}

type UpdateUserRequest struct {
	Name            string `json:"name,omitempty"`
	Email           string `json:"email,omitempty"`
	Password        string `json:"password,omitempty"`
	ConfirmPassword string `json:"confirm_password,omitempty"`
}

func (params UpdateUserRequest) ValidateUpdateUserData() error {
	err := validation.ValidateStruct(&params,
		validation.Field(
			&params.Name,
			validation.Length(2, 100).Error("name must be between 2 and 100 characters"),
		),
		validation.Field(
			&params.Email,
			is.Email.Error("Email must be a valid email address"),
		),
		validation.Field(
			&params.Password,
			validation.Required.Error("Password is required"),
			validation.Length(8, 100).Error("Password must be at least 8 characters"),
		),
		validation.Field(
			&params.ConfirmPassword,
			ConfirmPasswordRequiredOnPasswordProvided(params.Password),
			passwordsMatch(params.Password),
		),
	)

	formattedErrors, ok := validator.FormatValidationErrors(err)
	if !ok {
		return nil
	}

	return common_errors.ErrValidationError{
		Errors: formattedErrors,
	}
}

// custom validation rule used above
func passwordsMatch(password string) validation.Rule {
	return validation.By(func(value interface{}) error {
		confirmPassword, ok := value.(string)
		if !ok {
			return errors.New("invalid password confirmation")
		}

		if confirmPassword != password {
			return errors.New("passwords do not match")
		}

		return nil
	})
}

func ConfirmPasswordRequiredOnPasswordProvided(password string) validation.Rule {
	return validation.When(password != "", validation.Required.Error("Confirm password is required"))
}
