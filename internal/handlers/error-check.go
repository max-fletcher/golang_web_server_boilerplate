package handlers

import (
	"net/http"

	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/responses"
)

// This function will return an error response no matter what.
func (handler *Handler) ErrorResponse(w http.ResponseWriter, r *http.Request) error {
	responses.InternalServerErrorSWW(w)

	return nil
}
