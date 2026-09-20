package user_roles

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	id_helpers "github.com/max-fletcher/golang_web_server_boilerplate/helpers/ID"
	validator "github.com/max-fletcher/golang_web_server_boilerplate/helpers/validation"
	common_errors "github.com/max-fletcher/golang_web_server_boilerplate/internal/errors"
)

// Struct to be validated
type CreateUserRoleRequest struct {
	UserID string `json:"user_id"`
	RoleID string `json:"role_id"`
}

type CreateUserRoleInput struct {
	UserID uuid.UUID `json:"user_id"`
	RoleID uuid.UUID `json:"role_id"`
}

// Rules
func (params CreateUserRoleRequest) ValidateCreateUserRoleData() (CreateUserRoleInput, error) {
	err := validation.ValidateStruct(&params,
		validation.Field(
			&params.RoleID,
			validation.Required.Error("User id is required"),
			is.UUID.Error("Not a valid UUID"),
		),
		validation.Field(
			&params.UserID,
			validation.Required.Error("Role id is required"),
			is.UUID.Error("Not a valid UUID"),
		),
	)

	formattedErrors, hasValidationErrors := validator.FormatValidationErrors(err)
	userID, uuidErr := id_helpers.ParseUUID(params.UserID, "user ID") // parsing UserID field
	if uuidErr != nil {                                               // if userID is not valid uuid, put it in formattedErrors and set hasValidationErrors to false
		formattedErrors["user_id"] = uuidErr.Error()
		hasValidationErrors = true
	}
	if hasValidationErrors {
		return CreateUserRoleInput{}, common_errors.ErrValidationError{
			Errors: formattedErrors,
		}
	}
	roleID, uuidErr := id_helpers.ParseUUID(params.RoleID, "role ID") // parsing RoleID field
	if uuidErr != nil {                                               // if roleID is not valid uuid, put it in formattedErrors and set hasValidationErrors to false
		formattedErrors["role_id"] = uuidErr.Error()
		hasValidationErrors = true
	}
	if hasValidationErrors {
		return CreateUserRoleInput{}, common_errors.ErrValidationError{
			Errors: formattedErrors,
		}
	}

	// Construct an instance of createUserRoleInput
	createUserRoleInput := CreateUserRoleInput{
		UserID: userID,
		RoleID: roleID,
	}

	return createUserRoleInput, nil

}
