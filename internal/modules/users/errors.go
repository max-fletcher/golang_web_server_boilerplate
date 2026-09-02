package users

import (
	"fmt"

	"github.com/google/uuid"
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
type ErrUserWithEmailAlreadyExists struct {
	Email string
}

func (e ErrUserWithEmailAlreadyExists) Error() string {
	return fmt.Sprintf("Email %s is already taken", e.Email)
}

type ErrUserWithEmailNotFound struct {
	Email string
}

func (e ErrUserWithEmailNotFound) Error() string {
	return fmt.Sprintf("User with email %s not found", e.Email)
}

type ErrUserWithIdNotFound struct {
	ID uuid.UUID
}

func (e ErrUserWithIdNotFound) Error() string {
	return fmt.Sprintf("User with id %s not found", e.ID)
}

type ErrUsersFetchFailed struct {
	fetchErr error
}

func (e ErrUsersFetchFailed) Error() string {
	return "Failed to fetch users"
}

func (e ErrUsersFetchFailed) Unwrap() error { // Unwrap shows underlying details of errors
	return e.fetchErr
}

type ErrUserFetchFailed struct {
	fetchErr error
}

func (e ErrUserFetchFailed) Error() string {
	return "Failed to fetch user"
}

func (e ErrUserFetchFailed) Unwrap() error { // Unwrap shows underlying details of errors
	return e.fetchErr
}

type ErrUserCreateFailed struct {
	createErr error
}

func (e ErrUserCreateFailed) Error() string {
	return "Failed to create user"
}

func (e ErrUserCreateFailed) Unwrap() error {
	return e.createErr
}

type ErrUserUpdateFailed struct {
	updateUser error
}

func (e ErrUserUpdateFailed) Error() string {
	return "Failed to update user"
}

func (e ErrUserUpdateFailed) Unwrap() error {
	return e.updateUser
}

type ErrUserDeleteFailed struct {
	deleteErr error
}

func (e ErrUserDeleteFailed) Error() string {
	return "Failed to delete user"
}

func (e ErrUserDeleteFailed) Unwrap() error {
	return e.deleteErr
}

// Static errors
var (
// ErrUserWithIdNotFound   = errors.New("User with this id doesn't exist")
)
