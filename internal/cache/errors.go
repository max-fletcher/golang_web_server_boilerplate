package cache

import (
	"net/http"

	common_errors "github.com/max-fletcher/golang_web_server_boilerplate/internal/errors"
)

type ErrCacheMiss struct {
	Err error
}

func (e ErrCacheMiss) Error() string {
	return "Cache miss"
}

func (e ErrCacheMiss) StatusCode() int {
	return http.StatusBadGateway
}

func (e ErrCacheMiss) ClientMsg() string {
	return "Cache miss"
}

var _ common_errors.ErrHTTPBaseError = ErrCacheMiss{}
