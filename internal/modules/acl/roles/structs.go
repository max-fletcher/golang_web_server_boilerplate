package roles

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	validator "github.com/max-fletcher/golang_web_server_boilerplate/helpers/validation"
	common_errors "github.com/max-fletcher/golang_web_server_boilerplate/internal/errors"
)

// Struct to be validated
type CreateRoleRequest struct {
	Name string `json:"name"`
}

type CreateRoleInput struct {
	Name string `json:"name"`
}

// Rules
func (params CreateRoleRequest) ValidateCreateRoleData() (CreateRoleInput, error) {
	err := validation.ValidateStruct(&params,
		validation.Field(
			&params.Name,
			validation.Required.Error("Name is required"),
			validation.Length(2, 100).Error("Name must be between 2 and 100 characters"),
		),
	)

	formattedErrors, hasValidationErrors := validator.FormatValidationErrors(err)
	if hasValidationErrors {
		return CreateRoleInput{}, common_errors.ErrValidationError{
			Errors: formattedErrors,
		}
	}
	// Construct an instance of createRoleInput
	createRoleInput := CreateRoleInput{
		Name: params.Name,
	}

	return createRoleInput, nil
}

type UpdateRoleRequest struct {
	Name string
}

type UpdateRoleInput struct {
	Name string `json:"name"`
}

func (params UpdateRoleRequest) ValidateUpdateRoleData() (UpdateRoleInput, error) {
	err := validation.ValidateStruct(&params,
		validation.Field(
			&params.Name,
			validation.Length(2, 100).Error("Name must be between 2 and 100 characters"),
		),
	)

	formattedErrors, hasValidationErrors := validator.FormatValidationErrors(err)
	if hasValidationErrors {
		return UpdateRoleInput{}, common_errors.ErrValidationError{
			Errors: formattedErrors,
		}
	}
	// Construct an instance of updateRoleInput
	updateRoleInput := UpdateRoleInput{
		Name: params.Name,
	}

	return updateRoleInput, nil
}
