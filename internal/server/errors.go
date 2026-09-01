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

	var invalidUUIDParamErr requests.ErrInvalidUUIDParam
	if errors.As(err, &invalidUUIDParamErr) {
		responses.BadRequestError(w, invalidUUIDParamErr.Error())
		return
	}

	// Invalid JSON field type
	var fieldTypeErr requests.ErrInvalidFieldType
	if errors.As(err, &fieldTypeErr) {
		responses.RespondWithDetailedErrors(
			w,
			http.StatusBadRequest,
			"Invalid field type",
			fieldTypeErr.FormatInvalidFieldTypeError(),
		)
		return
	}

	// Invalid/malformed JSON
	var invalidJSONErr requests.ErrInvalidJSON
	if errors.As(err, &invalidJSONErr) {
		server.Logger.Error(
			"++++ Invalid JSON Error: ++++",
			"error", invalidJSONErr.Error(),
			"cause", invalidJSONErr.Unwrap(),
		)
		responses.RespondWithError(w, http.StatusBadRequest, "Invalid JSON in request body")
		return
	}

	// Validation error
	var validationErr common_errors.ErrValidationError
	if errors.As(err, &validationErr) {
		responses.RespondWithDetailedErrors(
			w,
			http.StatusUnprocessableEntity,
			validationErr.Error(),
			validationErr.ErrorMap(),
		)
		return
	}

	// Hashing password error
	var passwordHashErr common_errors.ErrHashingPassword
	if errors.As(err, &passwordHashErr) {
		server.Logger.Error(
			"++++ Failed to hash password: ++++",
			"error", passwordHashErr.Error(),
			"cause", passwordHashErr.Unwrap(),
		)
		responses.InternalServerErrorSWW(w)
		return
	}

	// -------- Users module errors --------
	var emailExistsErr users.ErrUserWithEmailAlreadyExists
	if errors.As(err, &emailExistsErr) {
		responses.ConflictError(w, emailExistsErr.Error())
		return
	}

	var userWithIdNotFoundErr users.ErrUserWithIdNotFound
	if errors.As(err, &userWithIdNotFoundErr) {
		responses.NotFoundError(w, userWithIdNotFoundErr.Error())
		return
	}

	var userWithEmailNotFoundErr users.ErrUserWithEmailNotFound
	if errors.As(err, &userWithEmailNotFoundErr) {
		responses.NotFoundError(w, userWithEmailNotFoundErr.Error())
		return
	}

	var usersFetchFailedErr users.ErrUsersFetchFailed
	if errors.As(err, &usersFetchFailedErr) {
		server.Logger.Error(
			"++++ Failed to fetch users ++++",
			"error", usersFetchFailedErr.Error(),
			"cause", usersFetchFailedErr.Unwrap(),
		)

		responses.InternalServerErrorSWW(w)
		return
	}

	var userFetchFailedErr users.ErrUserFetchFailed
	if errors.As(err, &userFetchFailedErr) {
		server.Logger.Error(
			"++++ Failed to fetch single user ++++",
			"error", userFetchFailedErr.Error(),
			"cause", userFetchFailedErr.Unwrap(),
		)

		responses.InternalServerErrorSWW(w)
		return
	}

	var userCreateFailedErr users.ErrUserCreateFailed
	if errors.As(err, &userCreateFailedErr) {
		server.Logger.Error(
			"++++ Failed to create user ++++",
			"error", userCreateFailedErr.Error(),
			"cause", userCreateFailedErr.Unwrap(),
		)

		responses.InternalServerErrorSWW(w)
		return
	}

	// Unknown/unexpected error.
	server.Logger.Error("++++ Unhandled application error ++++", "error", err)
	responses.InternalServerErrorSWW(w)
}
