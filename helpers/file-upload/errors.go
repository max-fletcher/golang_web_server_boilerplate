package fileupload

import (
	"net/http"

	common_errors "github.com/max-fletcher/golang_web_server_boilerplate/internal/errors"
)

// Error storing file on disk
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

func (e ErrFileStorageError) Unwrap() error {
	return e.Err
}

var _ common_errors.ErrHTTPServerError = ErrFileStorageError{}

// Error deleting file
type ErrFileDeleteError struct {
	Err error
}

func (e ErrFileDeleteError) Error() string {
	return "Failed to delete file from storage"
}

func (e ErrFileDeleteError) StatusCode() int {
	return http.StatusInternalServerError
}

func (e ErrFileDeleteError) ClientMsg() string {
	return "Failed to delete file from storage"
}

func (e ErrFileDeleteError) Unwrap() error {
	return e.Err
}

var _ common_errors.ErrHTTPServerError = ErrFileDeleteError{}
