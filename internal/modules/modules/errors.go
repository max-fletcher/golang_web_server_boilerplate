package modules

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
// 1. use "if errors.Is(err, ErrModuleNotFound) {...}" for static errors to figure out what type of error it is and determine what status
// code should be conditionally returned with the response from handler/controller/topmost func)
// 2. use (handler/controller/topmost func):
// 	var emailExistsErr users.ErrModuleWithEmailAlreadyExists
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
type ErrModuleWithIdNotFound struct {
	ID uuid.UUID
}

func (e ErrModuleWithIdNotFound) Error() string {
	return fmt.Sprintf("Module with id %s not found", e.ID)
}

func (e ErrModuleWithIdNotFound) StatusCode() int {
	return http.StatusNotFound
}

func (e ErrModuleWithIdNotFound) ClientMsg() string {
	return fmt.Sprintf("Module with id %s not found", e.ID)
}

var _ common_errors.ErrHTTPBaseError = ErrModuleWithIdNotFound{}

type ErrModulesFetchFailed struct {
	FetchErr error
}

func (e ErrModulesFetchFailed) Error() string {
	return "Failed to fetch modules"
}

func (e ErrModulesFetchFailed) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrModulesFetchFailed) ClientMsg() string {
	return "Failed to fetch modules"
}

func (e ErrModulesFetchFailed) Unwrap() error { // Unwrap shows underlying details of errors
	return e.FetchErr
}

var _ common_errors.ErrHTTPServerError = ErrModulesFetchFailed{}

type ErrModuleFetchFailed struct {
	FetchErr error
}

func (e ErrModuleFetchFailed) Error() string {
	return "Failed to fetch module"
}

func (e ErrModuleFetchFailed) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrModuleFetchFailed) ClientMsg() string {
	return "Failed to fetch module"
}

func (e ErrModuleFetchFailed) Unwrap() error { // Unwrap shows underlying details of errors
	return e.FetchErr
}

var _ common_errors.ErrHTTPServerError = ErrModuleFetchFailed{}

type ErrModuleCreateFailed struct {
	CreateErr error
}

func (e ErrModuleCreateFailed) Error() string {
	return "Failed to create module"
}

func (e ErrModuleCreateFailed) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrModuleCreateFailed) ClientMsg() string {
	return "Failed to create module"
}

func (e ErrModuleCreateFailed) Unwrap() error {
	return e.CreateErr
}

var _ common_errors.ErrHTTPServerError = ErrModuleCreateFailed{}

type ErrModuleUpdateFailed struct {
	updateModule error
}

func (e ErrModuleUpdateFailed) Error() string {
	return "Failed to update module"
}

func (e ErrModuleUpdateFailed) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrModuleUpdateFailed) ClientMsg() string {
	return "Failed to update module"
}

func (e ErrModuleUpdateFailed) Unwrap() error {
	return e.updateModule
}

var _ common_errors.ErrHTTPServerError = ErrModuleUpdateFailed{}

type ErrModuleDeleteFailed struct {
	DeleteErr error
}

func (e ErrModuleDeleteFailed) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrModuleDeleteFailed) ClientMsg() string {
	return "Failed to delete module"
}

func (e ErrModuleDeleteFailed) Error() string {
	return "Failed to delete module"
}

func (e ErrModuleDeleteFailed) Unwrap() error {
	return e.DeleteErr
}

var _ common_errors.ErrHTTPServerError = ErrModuleDeleteFailed{}

// Static errors
// var (
// 	ErrModuleWithIdNotFound   = errors.New("Module with this id doesn't exist")
// )
