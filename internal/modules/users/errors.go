package users

import (
	"fmt"

	"github.com/google/uuid"
)

// NOTE: rule of thumb for defining custom errors:
// 1. Define error interfaces for dynamic errors(errors that need to be constructed using variables)(see below example)
// 2. Define static errors for static strings(see below example)
// rule of thumb for handling custom errors:
// 1. use "if err != nil{...}" for static one-off errors(error that contain only string and you don't intend to return it
//    conditionally e.g you don't want to determine which status code should be thrown)
// 2. use "if errors.Is(err, ErrUserNotFound) {...}" for error static errors that you want to return conditionally(i.e you
//    want to determine what status code should be returned with response from handler/controller)
// 3. use
// "if err != nil {
// 	var emailExistsErr users.ErrUserWithEmailAlreadyExists
// 	if errors.As(err, &emailExistsErr) {
// 		responses.ConflictError(w, emailExistsErr.Error())
// 		return
// 	}
//  ...
// }"
// for dynamic errors and you want to return response conditionally(i.e determine which status code to throw based on what
// error type). errors.As stores the err into &emailExistsErr if the type matches and from there, you can throw it how you please

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
	return fmt.Sprintf("Failed to fetch users")
}

func (e ErrUsersFetchFailed) Unwrap() error {
	return e.fetchErr
}

type ErrUserFetchFailed struct {
	fetchErr error
}

func (e ErrUserFetchFailed) Error() string {
	return fmt.Sprintf("Failed to fetch user")
}

func (e ErrUserFetchFailed) Unwrap() error {
	return e.fetchErr
}

type ErrUserCreateFailed struct {
	createErr error
}

func (e ErrUserCreateFailed) Error() string {
	return fmt.Sprintf("Failed to create user")
}

func (e ErrUserCreateFailed) Unwrap() error {
	return e.createErr
}

// Static errors
var (
// ErrUserWithIdNotFound   = errors.New("User with this id doesn't exist")
)
