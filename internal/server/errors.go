package server

import (
	"errors"
	"net/http"

	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/responses"
	common_errors "github.com/max-fletcher/golang_web_server_boilerplate/internal/errors"
)

func (server *Server) HandleError(w http.ResponseWriter, err error) {
	// Errors generated from validation or invalid JSON that contains map of key-values
	var httpErrorMapErr common_errors.ErrHTTPWithErrorMap
	if errors.As(err, &httpErrorMapErr) {
		responses.RespondWithDetailedErrors(
			w,
			httpErrorMapErr.StatusCode(),
			httpErrorMapErr.ClientMsg(),
			httpErrorMapErr.ErrorMap(),
		)
		return
	}

	// Errors generated from server that must contain the entire error(logging it for easier debugging)
	var httpServerErr common_errors.ErrHTTPServerError
	if errors.As(err, &httpServerErr) {
		server.Logger.Error(
			"++++"+httpServerErr.Error()+"++++",
			"error", httpServerErr.Error(),
			"cause", httpServerErr.Unwrap(),
		)

		responses.RespondWithError(w, httpServerErr.StatusCode(), httpServerErr.ClientMsg())
		return
	}

	// Generic errors
	var httpBaseErr common_errors.ErrHTTPBaseError
	if errors.As(err, &httpBaseErr) {
		responses.RespondWithError(w, httpBaseErr.StatusCode(), httpBaseErr.ClientMsg())
		return
	}

	server.Logger.Error(
		"++++ Internal server error ++++",
		"error", err.Error(),
		"cause", err,
	)

	responses.InternalServerErrorSWW(w)
}
