package role_permissions

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
// 1. use "if errors.Is(err, ErrRolePermissionNotFound) {...}" for static errors to figure out what type of error it is and determine what status
// code should be conditionally returned with the response from handler/controller/topmost func)
// 2. use (handler/controller/topmost func):
// 	var emailExistsErr users.ErrRolePermissionWithEmailAlreadyExists
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
type ErrRolePermissionWithIdNotFound struct {
	ID uuid.UUID
}

func (e ErrRolePermissionWithIdNotFound) Error() string {
	return fmt.Sprintf("Role-permission with id %s not found", e.ID)
}

func (e ErrRolePermissionWithIdNotFound) StatusCode() int {
	return http.StatusNotFound
}

func (e ErrRolePermissionWithIdNotFound) ClientMsg() string {
	return fmt.Sprintf("Role-permission with id %s not found", e.ID)
}

var _ common_errors.ErrHTTPBaseError = ErrRolePermissionWithIdNotFound{}

type ErrRolePermissionWithRoleIdAndPermissionIdNotFound struct {
	RoleID       uuid.UUID
	PermissionID uuid.UUID
}

func (e ErrRolePermissionWithRoleIdAndPermissionIdNotFound) Error() string {
	return fmt.Sprintf("Role-permission with role id %s and permission id %s not found", e.RoleID, e.PermissionID)
}

func (e ErrRolePermissionWithRoleIdAndPermissionIdNotFound) StatusCode() int {
	return http.StatusNotFound
}

func (e ErrRolePermissionWithRoleIdAndPermissionIdNotFound) ClientMsg() string {
	return fmt.Sprintf("Role-permission with role id %s and permission id %s not found", e.RoleID, e.PermissionID)
}

var _ common_errors.ErrHTTPBaseError = ErrRolePermissionWithRoleIdAndPermissionIdNotFound{}

type ErrRolePermissionWithUserIdNotFound struct {
	ID uuid.UUID
}

func (e ErrRolePermissionWithUserIdNotFound) Error() string {
	return fmt.Sprintf("Role-permission with id %s not found", e.ID)
}

func (e ErrRolePermissionWithUserIdNotFound) StatusCode() int {
	return http.StatusNotFound
}

func (e ErrRolePermissionWithUserIdNotFound) ClientMsg() string {
	return fmt.Sprintf("Role-permission with id %s not found", e.ID)
}

var _ common_errors.ErrHTTPBaseError = ErrRolePermissionWithUserIdNotFound{}

type ErrRolePermissionsFetchFailed struct {
	FetchErr error
}

func (e ErrRolePermissionsFetchFailed) Error() string {
	return "Failed to fetch user with role-permissions"
}

func (e ErrRolePermissionsFetchFailed) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrRolePermissionsFetchFailed) ClientMsg() string {
	return "Failed to fetch user with role-permissions"
}

func (e ErrRolePermissionsFetchFailed) Unwrap() error { // Unwrap shows underlying details of errors
	return e.FetchErr
}

var _ common_errors.ErrHTTPServerError = ErrRolePermissionsFetchFailed{}

type ErrRolePermissionFetchFailed struct {
	FetchErr error
}

func (e ErrRolePermissionFetchFailed) Error() string {
	return "Failed to fetch role-permissions"
}

func (e ErrRolePermissionFetchFailed) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrRolePermissionFetchFailed) ClientMsg() string {
	return "Failed to fetch role-permissions"
}

func (e ErrRolePermissionFetchFailed) Unwrap() error { // Unwrap shows underlying details of errors
	return e.FetchErr
}

var _ common_errors.ErrHTTPServerError = ErrRolePermissionFetchFailed{}

type ErrRolePermissionWithUserFetchFailed struct {
	FetchErr error
}

func (e ErrRolePermissionWithUserFetchFailed) Error() string {
	return "Failed to fetch user with role-permissions"
}

func (e ErrRolePermissionWithUserFetchFailed) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrRolePermissionWithUserFetchFailed) ClientMsg() string {
	return "Failed to fetch user with role-permissions"
}

func (e ErrRolePermissionWithUserFetchFailed) Unwrap() error { // Unwrap shows underlying details of errors
	return e.FetchErr
}

var _ common_errors.ErrHTTPServerError = ErrRolePermissionWithUserFetchFailed{}

type ErrRolePermissionCreateFailed struct {
	CreateErr error
}

func (e ErrRolePermissionCreateFailed) Error() string {
	return "Failed to assign role-permission to user"
}

func (e ErrRolePermissionCreateFailed) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrRolePermissionCreateFailed) ClientMsg() string {
	return "Failed to assign role-permission to user"
}

func (e ErrRolePermissionCreateFailed) Unwrap() error {
	return e.CreateErr
}

var _ common_errors.ErrHTTPServerError = ErrRolePermissionCreateFailed{}

type ErrRolePermissionWithRoleIdAndPermissionIdAlreadyExists struct {
	RoleID       uuid.UUID
	PermissionID uuid.UUID
	Err          error
}

func (e ErrRolePermissionWithRoleIdAndPermissionIdAlreadyExists) Error() string {
	return fmt.Sprintf("Role-permission with role ID %v and permission ID %v already exists", e.RoleID, e.PermissionID)
}

func (e ErrRolePermissionWithRoleIdAndPermissionIdAlreadyExists) StatusCode() int {
	return http.StatusConflict
}

func (e ErrRolePermissionWithRoleIdAndPermissionIdAlreadyExists) ClientMsg() string {
	return fmt.Sprintf("Role-permission with role ID %v and permission ID %v already exists", e.RoleID, e.PermissionID)
}

var _ common_errors.ErrHTTPBaseError = ErrRolePermissionWithRoleIdAndPermissionIdAlreadyExists{}

type ErrRolePermissionInvalidRoleIdOrPermissionId struct {
	RoleID       uuid.UUID
	PermissionID uuid.UUID
	Err          error
}

func (e ErrRolePermissionInvalidRoleIdOrPermissionId) Error() string {
	return fmt.Sprintf("Either role ID %v or permission ID %v is invalid", e.RoleID, e.PermissionID)
}

func (e ErrRolePermissionInvalidRoleIdOrPermissionId) StatusCode() int {
	return http.StatusBadRequest
}

func (e ErrRolePermissionInvalidRoleIdOrPermissionId) ClientMsg() string {
	return fmt.Sprintf("Either role ID %v or permission ID %v is invalid", e.RoleID, e.PermissionID)
}

var _ common_errors.ErrHTTPBaseError = ErrRolePermissionInvalidRoleIdOrPermissionId{}

type ErrRolePermissionDeleteFailed struct {
	DeleteErr error
}

func (e ErrRolePermissionDeleteFailed) Error() string {
	return "Failed to remove permission from user"
}

func (e ErrRolePermissionDeleteFailed) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrRolePermissionDeleteFailed) ClientMsg() string {
	return "Failed to remove permission from user"
}

func (e ErrRolePermissionDeleteFailed) Unwrap() error {
	return e.DeleteErr
}

var _ common_errors.ErrHTTPServerError = ErrRolePermissionDeleteFailed{}

// Static errors
// var (
// 	ErrRolePermissionWithIdNotFound   = errors.New("RolePermission with this id doesn't exist")
// )
