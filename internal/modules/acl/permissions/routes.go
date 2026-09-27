package permissions

import (
	"github.com/go-chi/chi/v5"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/httpx"
	"github.com/max-fletcher/golang_web_server_boilerplate/middleware"
)

func RegisterRoutes(
	router chi.Router,
	handler *Handler,
	httpxHandler *httpx.Handler,
	aclMiddlware *middleware.ACLMiddleware,
) {
	router.Route("/permissions", func(router chi.Router) {
		router.With(aclMiddlware.RequirePermission(ACLRead)).Get("/", httpxHandler.Handle(handler.GetAll))
		router.With(aclMiddlware.RequirePermission(ACLCreate)).Post("/", httpxHandler.Handle(handler.Create))
		router.With(aclMiddlware.RequirePermission(ACLRead)).Get("/{id}", httpxHandler.Handle(handler.GetByID))
	})
}
