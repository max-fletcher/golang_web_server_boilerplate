package server

import (
	"net/http"
)

// Type of all handler/controller functions. Handle method below will accept all handlers as long as they satisfy this fn signature.
type AppHandler func(http.ResponseWriter, *http.Request) error

func (server *Server) Handle(handler AppHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			server.HandleError(w, err) // defined in /server/errors.go
		}
	})
}
