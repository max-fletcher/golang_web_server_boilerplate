package role_permissions

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	id_helpers "github.com/max-fletcher/golang_web_server_boilerplate/helpers/ID"
	validator "github.com/max-fletcher/golang_web_server_boilerplate/helpers/validation"
	common_errors "github.com/max-fletcher/golang_web_server_boilerplate/internal/errors"
)

// Struct to be validated
type CreateRolePermissionRequest struct {
	RoleID       string `json:"role_id"`
	PermissionID string `json:"permission_id"`
}

type CreateRolePermissionInput struct {
	RoleID       uuid.UUID `json:"role_id"`
	PermissionID uuid.UUID `json:"permission_id"`
}

// Rules
func (params CreateRolePermissionRequest) ValidateCreateRolePermissionData() (CreateRolePermissionInput, error) {
	err := validation.ValidateStruct(&params,
		validation.Field(
			&params.RoleID,
			validation.Required.Error("Role is required"),
			is.UUID.Error("Not a valid UUID"),
		),
		validation.Field(
			&params.PermissionID,
			validation.Required.Error("Permission is required"),
			is.UUID.Error("Not a valid UUID"),
		),
	)

	formattedErrors, hasValidationErrors := validator.FormatValidationErrors(err)
	roleID, uuidErr := id_helpers.ParseUUID(params.RoleID, "role ID") // parsing RoleID field
	if uuidErr != nil {                                               // if roleID is not valid uuid, put it in formattedErrors and set hasValidationErrors to false
		formattedErrors["role_id"] = uuidErr.Error()
		hasValidationErrors = true
	}
	if hasValidationErrors {
		return CreateRolePermissionInput{}, common_errors.ErrValidationError{
			Errors: formattedErrors,
		}
	}
	permissionID, uuidErr := id_helpers.ParseUUID(params.PermissionID, "permission ID") // parsing PermissionID field
	if uuidErr != nil {                                                                 // if permissionID is not valid uuid, put it in formattedErrors and set hasValidationErrors to false
		formattedErrors["permission_id"] = uuidErr.Error()
		hasValidationErrors = true
	}
	if hasValidationErrors {
		return CreateRolePermissionInput{}, common_errors.ErrValidationError{
			Errors: formattedErrors,
		}
	}

	// Construct an instance of createRolePermissionInput
	createRolePermissionInput := CreateRolePermissionInput{
		RoleID:       roleID,
		PermissionID: permissionID,
	}

	return createRolePermissionInput, nil

}
