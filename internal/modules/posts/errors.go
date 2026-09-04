package posts

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
// 1. use "if errors.Is(err, ErrPostNotFound) {...}" for static errors to figure out what type of error it is and determine what status
// code should be conditionally returned with the response from handler/controller/topmost func)
// 2. use (handler/controller/topmost func):
// 	var emailExistsErr users.ErrPostWithEmailAlreadyExists
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
type ErrPostWithIdNotFound struct {
	ID uuid.UUID
}

func (e ErrPostWithIdNotFound) Error() string {
	return fmt.Sprintf("Post with id %s not found", e.ID)
}

func (e ErrPostWithIdNotFound) StatusCode() int {
	return http.StatusNotFound
}

func (e ErrPostWithIdNotFound) ClientMsg() string {
	return fmt.Sprintf("Post with id %s not found", e.ID)
}

var _ common_errors.ErrHTTPBaseError = ErrPostWithIdNotFound{}

type ErrPostsFetchFailed struct {
	fetchErr error
}

func (e ErrPostsFetchFailed) Error() string {
	return "Failed to fetch posts"
}

func (e ErrPostsFetchFailed) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrPostsFetchFailed) ClientMsg() string {
	return "Failed to fetch posts"
}

func (e ErrPostsFetchFailed) Unwrap() error { // Unwrap shows underlying details of errors
	return e.fetchErr
}

var _ common_errors.ErrHTTPServerError = ErrPostsFetchFailed{}

type ErrPostFetchFailed struct {
	fetchErr error
}

func (e ErrPostFetchFailed) Error() string {
	return "Failed to fetch post"
}

func (e ErrPostFetchFailed) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrPostFetchFailed) ClientMsg() string {
	return "Failed to fetch post"
}

func (e ErrPostFetchFailed) Unwrap() error { // Unwrap shows underlying details of errors
	return e.fetchErr
}

var _ common_errors.ErrHTTPServerError = ErrPostFetchFailed{}

type ErrPostCreateFailed struct {
	createErr error
}

func (e ErrPostCreateFailed) Error() string {
	return "Failed to create post"
}

func (e ErrPostCreateFailed) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrPostCreateFailed) ClientMsg() string {
	return "Failed to create post"
}

func (e ErrPostCreateFailed) Unwrap() error {
	return e.createErr
}

var _ common_errors.ErrHTTPServerError = ErrPostCreateFailed{}

type ErrPostUpdateFailed struct {
	updatePost error
}

func (e ErrPostUpdateFailed) Error() string {
	return "Failed to update post"
}

func (e ErrPostUpdateFailed) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrPostUpdateFailed) ClientMsg() string {
	return "Failed to update post"
}

func (e ErrPostUpdateFailed) Unwrap() error {
	return e.updatePost
}

var _ common_errors.ErrHTTPServerError = ErrPostUpdateFailed{}

type ErrPostDeleteFailed struct {
	deleteErr error
}

func (e ErrPostDeleteFailed) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrPostDeleteFailed) ClientMsg() string {
	return "Failed to delete post"
}

func (e ErrPostDeleteFailed) Error() string {
	return "Failed to delete post"
}

func (e ErrPostDeleteFailed) Unwrap() error {
	return e.deleteErr
}

var _ common_errors.ErrHTTPServerError = ErrPostDeleteFailed{}

// Static errors
// var (
// 	ErrPostWithIdNotFound   = errors.New("Post with this id doesn't exist")
// )
