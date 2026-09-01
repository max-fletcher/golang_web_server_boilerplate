package server

import (
	"net/http"
)

// Type of all handler/controller functions. Any function of this type can be passed to the "Handle" method below
type AppHandler func(http.ResponseWriter, *http.Request) error

// This "Handle" method will accept all handler funcs as long as they satisfy the signature defined by AppHandler.
func (server *Server) Handle(handler AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil { // if handler throws error, log it
			server.HandleError(w, err)
		}
	}
}
