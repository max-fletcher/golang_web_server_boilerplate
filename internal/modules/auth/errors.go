package auth

import (
	"net/http"

	common_errors "github.com/max-fletcher/golang_web_server_boilerplate/internal/errors"
)

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

// Dynamic errors
type ErrPasswordMismatch struct {
	Err error
}

func (e ErrPasswordMismatch) Error() string {
	return "Password doesn't match"
}

func (e ErrPasswordMismatch) StatusCode() int {
	return http.StatusUnauthorized
}

func (e ErrPasswordMismatch) ClientMsg() string {
	return "Password doesn't match"
}

var _ common_errors.ErrHTTPBaseError = ErrPasswordMismatch{}

type ErrJWTGenerationFailed struct {
	Err error
}

func (e ErrJWTGenerationFailed) Error() string {
	return "Failed to generate JWT"
}

func (e ErrJWTGenerationFailed) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrJWTGenerationFailed) ClientMsg() string {
	return "Failed to generate JWT"
}

func (e ErrJWTGenerationFailed) Unwrap() error { // Unwrap shows underlying details of errors
	return e.Err
}

var _ common_errors.ErrHTTPServerError = ErrJWTGenerationFailed{}

// unexpected signing method
type ErrUnknownSigningMethod struct {
	Err error
}

func (e ErrUnknownSigningMethod) Error() string {
	return "Unexpected signing method"
}

func (e ErrUnknownSigningMethod) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrUnknownSigningMethod) ClientMsg() string {
	return "Unexpected signing method"
}

func (e ErrUnknownSigningMethod) Unwrap() error { // Unwrap shows underlying details of errors
	return e.Err
}

var _ common_errors.ErrHTTPServerError = ErrUnknownSigningMethod{}

type ErrInvalidAuthorizationHeader struct {
	Err error
}

func (e ErrInvalidAuthorizationHeader) Error() string {
	return "Invalid authorization header"
}

func (e ErrInvalidAuthorizationHeader) StatusCode() int {
	return http.StatusUnauthorized
}

func (e ErrInvalidAuthorizationHeader) ClientMsg() string {
	return "Invalid authorization header"
}

var _ common_errors.ErrHTTPBaseError = ErrInvalidAuthorizationHeader{}

type ErrJWTInvalid struct {
	Err error
}

func (e ErrJWTInvalid) Error() string {
	return "Invalid JWT"
}

func (e ErrJWTInvalid) StatusCode() int {
	return http.StatusUnauthorized
}

func (e ErrJWTInvalid) ClientMsg() string {
	return "Invalid JWT"
}

var _ common_errors.ErrHTTPBaseError = ErrJWTInvalid{}

type ErrInvalidRefreshToken struct {
	Err error
}

func (e ErrInvalidRefreshToken) Error() string {
	return "Invalid refresh token"
}

func (e ErrInvalidRefreshToken) StatusCode() int {
	return http.StatusUnauthorized
}

func (e ErrInvalidRefreshToken) ClientMsg() string {
	return "Invalid refresh token"
}

var _ common_errors.ErrHTTPBaseError = ErrInvalidRefreshToken{}
