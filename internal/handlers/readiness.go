package handlers

import (
	"net/http"

	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/responses"
)

// This function will take a payload(struct) and malshal it into a JSON string that will be sent as bytes of data
func (handler *Handler) HandleReadiness(w http.ResponseWriter, r *http.Request) {
	responses.RespondWithJSON(w, 200, responses.Response{
		Code:    200,
		Status:  "ok",
		Message: "Server Running",
	})
}
