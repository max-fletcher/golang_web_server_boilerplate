package server

import (
	"net/http"

	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/responses"
)

type response struct {
	// The part `json:"status"` is called a struct tag
	// It tells fo than when converting this struct to/from JSON, use this name
	// Without this, we will get { "Status": "ok" } instead of { "status": "ok" }
	// This tag has many variations e.g
	// Data interface{} `json:"data,omitempty"` - Field is excluded from JSON
	// Data interface{} `json:"data,omitempty"` - Only include this field if it has a value
	// UserID int `json:"user_id"` - COnvert UserID to user_id when converting to/from JSON(what we are using here)
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// This function will take a payload(struct) and malshal it into a JSON string that will be sent as bytes of data
func (server *Server) handleReadiness(w http.ResponseWriter, r *http.Request) {
	responses.RespondWithJSON(w, 200, response{
		Code:    200,
		Status:  "ok",
		Message: "Server running",
	})
}
