package server

import (
	"errors"
	"net/http"

	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/requests"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/responses"
	common_errors "github.com/max-fletcher/golang_web_server_boilerplate/internal/errors"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/users"
)

// A method for Server struct. HandleError accepts a writer and an error, logs it if its unknown and returns a response
func (server *Server) HandleError(w http.ResponseWriter, err error) {
	// -------- Common errors --------
	// Invalid JSON field type
	var fieldTypeErr requests.ErrInvalidFieldType
	if errors.As(err, &fieldTypeErr) {
		responses.RespondWithDetailedErrors(
			w,
			http.StatusBadRequest,
			"Invalid field type",
			map[string]string{
				fieldTypeErr.Field: fieldTypeErr.Error(),
			},
		)
		return
	}

	// Invalid/malformed JSON
	var invalidJSONErr requests.ErrInvalidJSON
	if errors.As(err, &invalidJSONErr) {
		responses.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	var validationErr common_errors.ValidationError
	if errors.As(err, &validationErr) {
		responses.RespondWithDetailedErrors(
			w,
			http.StatusUnprocessableEntity,
			"Validation failed",
			validationErr.Errors,
		)
		return
	}

	// -------- Users module errors --------
	var emailExistsErr users.ErrUserWithEmailAlreadyExists
	if errors.As(err, &emailExistsErr) {
		responses.ConflictError(w, emailExistsErr.Error())
		return
	}

	var userNotFoundErr users.ErrUserWithIdNotFound
	if errors.As(err, &userNotFoundErr) {
		responses.NotFoundError(w, "User not found")
		return
	}

	var hashingPasswordErr common_errors.ErrHashingPassword // will store error if type matches in errors.As(...)
	if errors.As(err, hashingPasswordErr) {
		server.Logger.Error(
			"Failed to hash password",
			"error", hashingPasswordErr.Error(),
		)

		responses.InternalServerError(
			w,
			"Something went wrong. Please try again.",
		)
		return
	}

	var userFetchFailedErr users.ErrUserFetchFailed

	if errors.As(err, &userFetchFailedErr) {
		server.Logger.Error(
			"Failed to fetch user",
			"error", userFetchFailedErr,
		)

		responses.InternalServerError(
			w,
			"Something went wrong. Please try again.",
		)
		return
	}

	var userCreateFailedErr users.ErrUserCreateFailed

	if errors.As(err, &userCreateFailedErr) {
		server.Logger.Error(
			"Failed to create user",
			"error", userCreateFailedErr,
		)

		responses.InternalServerError(
			w,
			"Something went wrong. Please try again.",
		)
		return
	}

	// Unknown/unexpected error.
	server.Logger.Error("Unhandled application error", "error", err)

	responses.InternalServerError(w, "Something went wrong. Please try again.")
}
