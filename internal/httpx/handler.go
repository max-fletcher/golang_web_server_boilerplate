package httpx

import (
	"log/slog"
	"net/http"
)

type Handler struct {
	Logger *slog.Logger
}

// Type of all module handler/controller functions. Any function of this type can be passed to the "Handle" method below
type AppHandler func(http.ResponseWriter, *http.Request) error

// This "Handle" method will accept all handler funcs as long as they satisfy the signature defined by AppHandler.
func (handler *Handler) Handle(appHandler AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := appHandler(w, r); err != nil { // if handler throws error, log it
			handler.HandleError(w, err)
		}
	}
}

func NewHttpxHandler(logger *slog.Logger) *Handler {
	return &Handler{
		Logger: logger,
	}
}
