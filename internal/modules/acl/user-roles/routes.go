package user_roles

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
	router.Route("/user-roles", func(router chi.Router) {
		router.With(aclMiddlware.RequirePermission(ACLReadUser, ACLReadRole)).Get("/", httpxHandler.Handle(handler.GetAllUsersWithRoles))
		router.With(aclMiddlware.RequirePermission(ACLCreate)).Post("/", httpxHandler.Handle(handler.Create))
		router.With(aclMiddlware.RequirePermission(ACLRead)).Get("/{id}", httpxHandler.Handle(handler.GetByID))
		router.With(aclMiddlware.RequirePermission(ACLDelete)).Delete("/{id}", httpxHandler.Handle(handler.Delete))
		router.With(aclMiddlware.RequirePermission(ACLReadUser, ACLReadRole)).Get("/users/{userID}", httpxHandler.Handle(handler.GetByUserID))
		router.With(aclMiddlware.RequirePermission(ACLDelete)).Delete("/users/{userID}/roles/{roleID}", httpxHandler.Handle(handler.DeleteByUserIDAndRoleID))
	})
}
