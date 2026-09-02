package common_errors

import "net/http"

// NOTE: rule of thumb for defining custom errors:
// 1. Define static errors for static strings(see below example)
// 2. Define error interfaces for dynamic errors(errors that need to be constructed using variables)(see below example)
// rule of thumb for handling custom errors:
// 1. use "if errors.Is(err, ErrUserNotFound) {...}" for static errors to figure out what type of error it is and determine what status
// code should be conditionally returned with the response from handler/controller/topmost func)
// 2. use (handler/controller/topmost func):
// 	var emailExistsErr users.ErrUserWithEmailAlreadyExists
// 	if errors.As(err, &emailExistsErr) {
//    // log error e.g server.Logger.Error(...)
// 		responses.ConflictError(w, emailExistsErr.Error())
// 		return
// 	}
// for dynamic errors and you want to determine which status code to throw based on error type.
// errors.As stores the err into &emailExistsErr if the type matches and from there, you can throw it how you please.
// You can have multiple methods for this error type and format and log/respond with whatever you want(unlike the
// static error mentioned above).
// In this project, we are handling all errors in HandleError func from internal/server/errors.go

type ErrHTTPBaseError interface { // All errors should implement this
	Error() string
	StatusCode() int
	ClientMsg() string
}

type ErrHTTPServerError interface { // Errors generated from server that must contain stack trace
	ErrHTTPBaseError
	Unwrap() error
}

type ErrHTTPWithErrorMap interface { // Errors generated from validation or invalid JSON that contains map of key-values
	ErrHTTPBaseError
	ErrorMap() map[string]string
}

type ErrValidationError struct {
	Errors map[string]string
}

func (e ErrValidationError) Error() string {
	return "Validation failed"
}

func (e ErrValidationError) StatusCode() int {
	return http.StatusUnprocessableEntity
}

func (e ErrValidationError) ClientMsg() string {
	return "Validation failed"
}

func (e ErrValidationError) ErrorMap() map[string]string {
	return e.Errors
}

var _ ErrHTTPWithErrorMap = ErrValidationError{}

type ErrHashingPassword struct {
	HashErr error
}

func (e ErrHashingPassword) Error() string {
	return "Failed to hash password"
}

func (e ErrHashingPassword) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrHashingPassword) ClientMsg() string {
	return "Failed to hash password"
}

func (e ErrHashingPassword) Unwrap() error {
	return e.HashErr
}

var _ ErrHTTPServerError = ErrHashingPassword{}

var (
// common
// ErrNotLoggedIn            = errors.New("You are not logged in")
// ErrUnauthorized           = errors.New("You are unauthorized to access this resource")
// ErrInsufficientPermission = errors.New("You don't have sufficient permission to access this resource")
// ErrHashingPassword = errors.New("Error processing password")
)
