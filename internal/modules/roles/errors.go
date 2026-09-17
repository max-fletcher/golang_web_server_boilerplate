package roles

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	common_errors "github.com/max-fletcher/golang_web_server_boilerplate/internal/errors"
)

// NOTE: rule of thumb for defining custom errors:
// 1. Define static errors for static strings(see below example)
// 2. Define error interfaces for dynamic errors(errors that need to be constructed using variables)(see below example)
// rule of thumb for handling custom errors:
// 1. use "if errors.Is(err, ErrRoleNotFound) {...}" for static errors to figure out what type of error it is and determine what status
// code should be conditionally returned with the response from handler/controller/topmost func)
// 2. use (handler/controller/topmost func):
// 	var emailExistsErr users.ErrRoleWithEmailAlreadyExists
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

// // Dynamic errors
type ErrRoleWithIdNotFound struct {
	ID uuid.UUID
}

func (e ErrRoleWithIdNotFound) Error() string {
	return fmt.Sprintf("Role with id %s not found", e.ID)
}

func (e ErrRoleWithIdNotFound) StatusCode() int {
	return http.StatusNotFound
}

func (e ErrRoleWithIdNotFound) ClientMsg() string {
	return fmt.Sprintf("Role with id %s not found", e.ID)
}

var _ common_errors.ErrHTTPBaseError = ErrRoleWithIdNotFound{}

type ErrRolesFetchFailed struct {
	FetchErr error
}

func (e ErrRolesFetchFailed) Error() string {
	return "Failed to fetch Roles"
}

func (e ErrRolesFetchFailed) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrRolesFetchFailed) ClientMsg() string {
	return "Failed to fetch Roles"
}

func (e ErrRolesFetchFailed) Unwrap() error { // Unwrap shows underlying details of errors
	return e.FetchErr
}

var _ common_errors.ErrHTTPServerError = ErrRolesFetchFailed{}

type ErrRoleFetchFailed struct {
	FetchErr error
}

func (e ErrRoleFetchFailed) Error() string {
	return "Failed to fetch Role"
}

func (e ErrRoleFetchFailed) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrRoleFetchFailed) ClientMsg() string {
	return "Failed to fetch Role"
}

func (e ErrRoleFetchFailed) Unwrap() error { // Unwrap shows underlying details of errors
	return e.FetchErr
}

var _ common_errors.ErrHTTPServerError = ErrRoleFetchFailed{}

type ErrRoleCreateFailed struct {
	createErr error
}

func (e ErrRoleCreateFailed) Error() string {
	return "Failed to create Role"
}

func (e ErrRoleCreateFailed) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrRoleCreateFailed) ClientMsg() string {
	return "Failed to create Role"
}

func (e ErrRoleCreateFailed) Unwrap() error {
	return e.createErr
}

var _ common_errors.ErrHTTPServerError = ErrRoleCreateFailed{}

type ErrRoleUpdateFailed struct {
	UpdateRole error
}

func (e ErrRoleUpdateFailed) Error() string {
	return "Failed to update Role"
}

func (e ErrRoleUpdateFailed) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrRoleUpdateFailed) ClientMsg() string {
	return "Failed to update Role"
}

func (e ErrRoleUpdateFailed) Unwrap() error {
	return e.UpdateRole
}

var _ common_errors.ErrHTTPServerError = ErrRoleUpdateFailed{}

type ErrRoleDeleteFailed struct {
	deleteErr error
}

func (e ErrRoleDeleteFailed) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrRoleDeleteFailed) ClientMsg() string {
	return "Failed to delete Role"
}

func (e ErrRoleDeleteFailed) Error() string {
	return "Failed to delete Role"
}

func (e ErrRoleDeleteFailed) Unwrap() error {
	return e.deleteErr
}

var _ common_errors.ErrHTTPServerError = ErrRoleDeleteFailed{}

// Static errors
// var (
// 	ErrRoleWithIdNotFound   = errors.New("Role with this id doesn't exist")
// )
