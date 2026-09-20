package user_roles

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
// 1. use "if errors.Is(err, ErrUserRoleNotFound) {...}" for static errors to figure out what type of error it is and determine what status
// code should be conditionally returned with the response from handler/controller/topmost func)
// 2. use (handler/controller/topmost func):
// 	var emailExistsErr users.ErrUserRoleWithEmailAlreadyExists
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
type ErrUserRoleWithIdNotFound struct {
	ID uuid.UUID
}

func (e ErrUserRoleWithIdNotFound) Error() string {
	return fmt.Sprintf("User-role with id %s not found", e.ID)
}

func (e ErrUserRoleWithIdNotFound) StatusCode() int {
	return http.StatusNotFound
}

func (e ErrUserRoleWithIdNotFound) ClientMsg() string {
	return fmt.Sprintf("User-role with id %s not found", e.ID)
}

var _ common_errors.ErrHTTPBaseError = ErrUserRoleWithIdNotFound{}

type ErrUserRoleWithUserIdAndRoleIdNotFound struct {
	RoleID uuid.UUID
	UserID uuid.UUID
}

func (e ErrUserRoleWithUserIdAndRoleIdNotFound) Error() string {
	return fmt.Sprintf("User-role with role id %s and permission id %s not found", e.RoleID, e.UserID)
}

func (e ErrUserRoleWithUserIdAndRoleIdNotFound) StatusCode() int {
	return http.StatusNotFound
}

func (e ErrUserRoleWithUserIdAndRoleIdNotFound) ClientMsg() string {
	return fmt.Sprintf("User-role with role id %s and permission id %s not found", e.RoleID, e.UserID)
}

var _ common_errors.ErrHTTPBaseError = ErrUserRoleWithUserIdAndRoleIdNotFound{}

type ErrUserRoleWithUserIdNotFound struct {
	ID uuid.UUID
}

func (e ErrUserRoleWithUserIdNotFound) Error() string {
	return fmt.Sprintf("User-role with id %s not found", e.ID)
}

func (e ErrUserRoleWithUserIdNotFound) StatusCode() int {
	return http.StatusNotFound
}

func (e ErrUserRoleWithUserIdNotFound) ClientMsg() string {
	return fmt.Sprintf("User-role with id %s not found", e.ID)
}

var _ common_errors.ErrHTTPBaseError = ErrUserRoleWithUserIdNotFound{}

type ErrUserRolesFetchFailed struct {
	FetchErr error
}

func (e ErrUserRolesFetchFailed) Error() string {
	return "Failed to fetch user with user-roles"
}

func (e ErrUserRolesFetchFailed) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrUserRolesFetchFailed) ClientMsg() string {
	return "Failed to fetch user with user-roles"
}

func (e ErrUserRolesFetchFailed) Unwrap() error { // Unwrap shows underlying details of errors
	return e.FetchErr
}

var _ common_errors.ErrHTTPServerError = ErrUserRolesFetchFailed{}

type ErrUserRoleFetchFailed struct {
	FetchErr error
}

func (e ErrUserRoleFetchFailed) Error() string {
	return "Failed to fetch user-roles"
}

func (e ErrUserRoleFetchFailed) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrUserRoleFetchFailed) ClientMsg() string {
	return "Failed to fetch user-roles"
}

func (e ErrUserRoleFetchFailed) Unwrap() error { // Unwrap shows underlying details of errors
	return e.FetchErr
}

var _ common_errors.ErrHTTPServerError = ErrUserRoleFetchFailed{}

type ErrUserRoleWithUserFetchFailed struct {
	FetchErr error
}

func (e ErrUserRoleWithUserFetchFailed) Error() string {
	return "Failed to fetch user with user-roles"
}

func (e ErrUserRoleWithUserFetchFailed) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrUserRoleWithUserFetchFailed) ClientMsg() string {
	return "Failed to fetch user with user-roles"
}

func (e ErrUserRoleWithUserFetchFailed) Unwrap() error { // Unwrap shows underlying details of errors
	return e.FetchErr
}

var _ common_errors.ErrHTTPServerError = ErrUserRoleWithUserFetchFailed{}

type ErrUserRoleCreateFailed struct {
	CreateErr error
}

func (e ErrUserRoleCreateFailed) Error() string {
	return "Failed to assign user-role to user"
}

func (e ErrUserRoleCreateFailed) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrUserRoleCreateFailed) ClientMsg() string {
	return "Failed to assign user-role to user"
}

func (e ErrUserRoleCreateFailed) Unwrap() error {
	return e.CreateErr
}

var _ common_errors.ErrHTTPServerError = ErrUserRoleCreateFailed{}

type ErrUserRoleWithUserIdAndRoleIdAlreadyExists struct {
	RoleID uuid.UUID
	UserID uuid.UUID
	Err    error
}

func (e ErrUserRoleWithUserIdAndRoleIdAlreadyExists) Error() string {
	return fmt.Sprintf("User-role with role ID %v and permission ID %v already exists", e.RoleID, e.UserID)
}

func (e ErrUserRoleWithUserIdAndRoleIdAlreadyExists) StatusCode() int {
	return http.StatusConflict
}

func (e ErrUserRoleWithUserIdAndRoleIdAlreadyExists) ClientMsg() string {
	return fmt.Sprintf("User-role with role ID %v and permission ID %v already exists", e.RoleID, e.UserID)
}

var _ common_errors.ErrHTTPBaseError = ErrUserRoleWithUserIdAndRoleIdAlreadyExists{}

type ErrUserRoleInvaliduserIdOrRoleId struct {
	RoleID uuid.UUID
	UserID uuid.UUID
	Err    error
}

func (e ErrUserRoleInvaliduserIdOrRoleId) Error() string {
	return fmt.Sprintf("Either role ID %v or permission ID %v is invalid", e.RoleID, e.UserID)
}

func (e ErrUserRoleInvaliduserIdOrRoleId) StatusCode() int {
	return http.StatusBadRequest
}

func (e ErrUserRoleInvaliduserIdOrRoleId) ClientMsg() string {
	return fmt.Sprintf("Either role ID %v or permission ID %v is invalid", e.RoleID, e.UserID)
}

var _ common_errors.ErrHTTPBaseError = ErrUserRoleInvaliduserIdOrRoleId{}

type ErrUserRoleDeleteFailed struct {
	DeleteErr error
}

func (e ErrUserRoleDeleteFailed) Error() string {
	return "Failed to remove role from user"
}

func (e ErrUserRoleDeleteFailed) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrUserRoleDeleteFailed) ClientMsg() string {
	return "Failed to remove role from user"
}

func (e ErrUserRoleDeleteFailed) Unwrap() error {
	return e.DeleteErr
}

var _ common_errors.ErrHTTPServerError = ErrUserRoleDeleteFailed{}

// Static errors
// var (
// 	ErrUserRoleWithIdNotFound   = errors.New("UserRole with this id doesn't exist")
// )
