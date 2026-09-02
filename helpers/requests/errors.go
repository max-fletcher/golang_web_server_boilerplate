package requests

import (
	"fmt"
	"net/http"

	common_errors "github.com/max-fletcher/golang_web_server_boilerplate/internal/errors"
)

type ErrInvalidUUIDParam struct {
	Param string
	Err   error
}

func (e ErrInvalidUUIDParam) Error() string {
	return fmt.Sprintf("Invalid UUID for route param %v(%v)", e.Param, e.Err.Error())
}

func (e ErrInvalidUUIDParam) StatusCode() int {
	return http.StatusBadRequest
}

func (e ErrInvalidUUIDParam) ClientMsg() string {
	return fmt.Sprintf("Invalid UUID for route param %v(%v)", e.Param, e.Err.Error())
}

var _ common_errors.ErrHTTPBaseError = ErrInvalidUUIDParam{}

type ErrInvalidJSON struct {
	Err error
}

func (e ErrInvalidJSON) Error() string {
	return fmt.Sprintf("Invalid JSON(%v)", e.Err)
}

func (e ErrInvalidJSON) StatusCode() int {
	return http.StatusBadRequest
}

func (e ErrInvalidJSON) ClientMsg() string {
	return fmt.Sprintf("Invalid JSON(%v)", e.Err)
}

func (e ErrInvalidJSON) Unwrap() error {
	return e.Err
}

var _ common_errors.ErrHTTPServerError = ErrInvalidJSON{}

type ErrInvalidFieldType struct {
	Field string
	Type  string
}

func (e ErrInvalidFieldType) Error() string {
	return fmt.Sprintf(
		"%s must be of type %s",
		e.Field,
		e.Type,
	)
}

func (e ErrInvalidFieldType) StatusCode() int {
	return http.StatusBadRequest
}

func (e ErrInvalidFieldType) ClientMsg() string {
	return fmt.Sprintf(
		"%s must be of type %s",
		e.Field,
		e.Type,
	)
}

func (e ErrInvalidFieldType) ErrorMap() map[string]string {
	return map[string]string{
		e.Field: e.Error(),
	}
}

var _ common_errors.ErrHTTPWithErrorMap = ErrInvalidFieldType{}
