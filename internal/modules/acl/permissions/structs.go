package permissions

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	id_helpers "github.com/max-fletcher/golang_web_server_boilerplate/helpers/ID"
	validator "github.com/max-fletcher/golang_web_server_boilerplate/helpers/validation"
	common_errors "github.com/max-fletcher/golang_web_server_boilerplate/internal/errors"
)

// Struct to be validated
type CreatePermissionRequest struct {
	Name     string    `json:"name"`
	ModuleId uuid.UUID `json:"module_id"`
}

type CreatePermissionInput struct {
	Name     EnumPermission `json:"name"`
	ModuleId uuid.UUID      `json:"module_id"`
}

// Rules
func (params CreatePermissionRequest) ValidateCreatePermissionData() (CreatePermissionInput, error) {
	err := validation.ValidateStruct(&params,
		validation.Field(
			&params.Name,
			validation.Required.Error("Name is required"),
			validation.In(PermissionCreate, PermissionRead, PermissionUpdate, PermissionDelete).Error("Invalid name value"),
		),
		validation.Field(
			&params.ModuleId,
			validation.Required.Error("Module is required"),
			is.UUID.Error("Not a valid UUID"),
		),
	)

	formattedErrors, hasValidationErrors := validator.FormatValidationErrors(err)
	moduleID, uuidErr := id_helpers.ParseUUID(params.ModuleId.String(), "Permission ID") // parsing PermissionId field
	if uuidErr != nil {                                                                  // if PermissionID is not valid uuid, put it in formattedErrors and set hasValidationErrors to false
		formattedErrors["module_id"] = uuidErr.Error()
		hasValidationErrors = true
	}
	if hasValidationErrors {
		return CreatePermissionInput{}, common_errors.ErrValidationError{
			Errors: formattedErrors,
		}
	}
	// Construct an instance of createPermissionInput
	createPermissionInput := CreatePermissionInput{
		Name:     EnumPermission(params.Name),
		ModuleId: moduleID,
	}

	return createPermissionInput, nil
}
