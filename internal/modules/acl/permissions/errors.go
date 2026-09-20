package permissions

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
// 1. use "if errors.Is(err, ErrPermissionNotFound) {...}" for static errors to figure out what type of error it is and determine what status
// code should be conditionally returned with the response from handler/controller/topmost func)
// 2. use (handler/controller/topmost func):
// 	var emailExistsErr users.ErrPermissionWithEmailAlreadyExists
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
type ErrPermissionWithIdNotFound struct {
	ID uuid.UUID
}

func (e ErrPermissionWithIdNotFound) Error() string {
	return fmt.Sprintf("Permission with id %s not found", e.ID)
}

func (e ErrPermissionWithIdNotFound) StatusCode() int {
	return http.StatusNotFound
}

func (e ErrPermissionWithIdNotFound) ClientMsg() string {
	return fmt.Sprintf("Permission with id %s not found", e.ID)
}

var _ common_errors.ErrHTTPBaseError = ErrPermissionWithIdNotFound{}

type ErrPermissionsFetchFailed struct {
	FetchErr error
}

func (e ErrPermissionsFetchFailed) Error() string {
	return "Failed to fetch permissions"
}

func (e ErrPermissionsFetchFailed) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrPermissionsFetchFailed) ClientMsg() string {
	return "Failed to fetch permissions"
}

func (e ErrPermissionsFetchFailed) Unwrap() error { // Unwrap shows underlying details of errors
	return e.FetchErr
}

var _ common_errors.ErrHTTPServerError = ErrPermissionsFetchFailed{}

type ErrPermissionFetchFailed struct {
	FetchErr error
}

func (e ErrPermissionFetchFailed) Error() string {
	return "Failed to fetch permission"
}

func (e ErrPermissionFetchFailed) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrPermissionFetchFailed) ClientMsg() string {
	return "Failed to fetch permission"
}

func (e ErrPermissionFetchFailed) Unwrap() error { // Unwrap shows underlying details of errors
	return e.FetchErr
}

var _ common_errors.ErrHTTPServerError = ErrPermissionFetchFailed{}

type ErrPermissionCreateFailed struct {
	CreateErr error
}

func (e ErrPermissionCreateFailed) Error() string {
	return "Failed to create permission"
}

func (e ErrPermissionCreateFailed) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrPermissionCreateFailed) ClientMsg() string {
	return "Failed to create permission"
}

func (e ErrPermissionCreateFailed) Unwrap() error {
	return e.CreateErr
}

var _ common_errors.ErrHTTPServerError = ErrPermissionCreateFailed{}

type ErrPermissionUpdateFailed struct {
	updatePermission error
}

func (e ErrPermissionUpdateFailed) Error() string {
	return "Failed to update permission"
}

func (e ErrPermissionUpdateFailed) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrPermissionUpdateFailed) ClientMsg() string {
	return "Failed to update permission"
}

func (e ErrPermissionUpdateFailed) Unwrap() error {
	return e.updatePermission
}

var _ common_errors.ErrHTTPServerError = ErrPermissionUpdateFailed{}

type ErrPermissionDeleteFailed struct {
	DeleteErr error
}

func (e ErrPermissionDeleteFailed) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrPermissionDeleteFailed) ClientMsg() string {
	return "Failed to delete permission"
}

func (e ErrPermissionDeleteFailed) Error() string {
	return "Failed to delete permission"
}

func (e ErrPermissionDeleteFailed) Unwrap() error {
	return e.DeleteErr
}

var _ common_errors.ErrHTTPServerError = ErrPermissionDeleteFailed{}

type ErrModuleWithPermissionAlreadyExists struct {
	Err error
}

func (e ErrModuleWithPermissionAlreadyExists) Error() string {
	return "Module with this permission already exists"
}

func (e ErrModuleWithPermissionAlreadyExists) StatusCode() int {
	return http.StatusConflict
}

func (e ErrModuleWithPermissionAlreadyExists) ClientMsg() string {
	return "Module with this permission already exists"
}

func (e ErrModuleWithPermissionAlreadyExists) Unwrap() error {
	return e.Err
}

var _ common_errors.ErrHTTPServerError = ErrModuleWithPermissionAlreadyExists{}

// Static errors
// var (
// 	ErrPermissionWithIdNotFound   = errors.New("Permission with this id doesn't exist")
// )
