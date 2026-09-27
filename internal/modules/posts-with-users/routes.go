package posts_with_users

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
	router.Route("/posts-with-user", func(router chi.Router) {
		router.With(aclMiddlware.RequirePermission(ACLReadPost, ACLReadUser)).Get("/", httpxHandler.Handle(handler.GetAll))
		router.With(aclMiddlware.RequirePermission(ACLCreatePost, ACLCreateUser)).Post("/", httpxHandler.Handle(handler.Create))
		router.With(aclMiddlware.RequirePermission(ACLReadPost, ACLReadUser)).Get("/{id}", httpxHandler.Handle(handler.GetByID))
	})
}
