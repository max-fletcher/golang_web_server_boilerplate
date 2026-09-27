package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/httpx"
)

func RegisterRoutes(
	router chi.Router,
	handler *Handler,
	httpxHandler *httpx.Handler,
) {
	router.Route("/auth", func(router chi.Router) {
		router.Post("/register", httpxHandler.Handle(handler.UserRegistration))
		router.Post("/login", httpxHandler.Handle(handler.UserLogin))
		router.Get("/refresh-token", httpxHandler.Handle(handler.RefreshToken))
	})
}
