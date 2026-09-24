package modules

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	validator "github.com/max-fletcher/golang_web_server_boilerplate/helpers/validation"
	common_errors "github.com/max-fletcher/golang_web_server_boilerplate/internal/errors"
	acl_constants "github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/acl/constants"
)

// Struct to be validated
type CreateModuleRequest struct {
	Name string `json:"name"`
}

type CreateModuleInput struct {
	Name acl_constants.EnumModuleNames `json:"name"`
}

// Rules
func (params CreateModuleRequest) ValidateCreateModuleData() (CreateModuleInput, error) {
	err := validation.ValidateStruct(&params,
		validation.Field(
			&params.Name,
			validation.Required.Error("Name is required"),
			validation.Length(2, 100).Error("Title must be between 2 and 100 characters"),
		),
	)

	formattedErrors, hasValidationErrors := validator.FormatValidationErrors(err)
	if hasValidationErrors {
		return CreateModuleInput{}, common_errors.ErrValidationError{
			Errors: formattedErrors,
		}
	}
	// Construct an instance of createModuleInput
	createModuleInput := CreateModuleInput{
		Name: acl_constants.EnumModuleNames(params.Name),
	}

	return createModuleInput, nil
}
