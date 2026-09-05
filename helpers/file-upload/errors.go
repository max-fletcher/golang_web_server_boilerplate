package fileupload

import (
	"net/http"

	common_errors "github.com/max-fletcher/golang_web_server_boilerplate/internal/errors"
)

// Number conversion Errors
type ErrFileStorageError struct {
	Err error
}

func (e ErrFileStorageError) Error() string {
	return "Failed to upload file to storage"
}

func (e ErrFileStorageError) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrFileStorageError) ClientMsg() string {
	return "Failed to upload file to storage"
}

var _ common_errors.ErrHTTPBaseError = ErrFileStorageError{}
