package role_permissions

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
	router.Route("/role-permissions", func(router chi.Router) {
		router.With(aclMiddlware.RequirePermission(ACLRead)).Get("/", httpxHandler.Handle(handler.GetAll))
		router.With(aclMiddlware.RequirePermission(ACLCreate)).Post("/", httpxHandler.Handle(handler.Create))
		router.With(aclMiddlware.RequirePermission(ACLRead)).Get("/{id}", httpxHandler.Handle(handler.GetByID))
		router.With(aclMiddlware.RequirePermission(ACLDelete)).Delete("/{id}", httpxHandler.Handle(handler.Delete))
		router.With(aclMiddlware.RequirePermission(ACLReadUser, ACLRead)).Get("/users", httpxHandler.Handle(handler.GetAllUsersWithRolePermissions))
		router.With(aclMiddlware.RequirePermission(ACLReadUser, ACLRead)).Get("/users/{userID}", httpxHandler.Handle(handler.GetByUserID))
		router.With(aclMiddlware.RequirePermission(ACLDelete)).Delete("/roles/{roleID}/permissions/{permissionID}", httpxHandler.Handle(handler.DeleteByRoleIDAndPermissionID))
	})
}
