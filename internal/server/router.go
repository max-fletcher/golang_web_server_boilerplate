package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/responses"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/posts"
	"github.com/max-fletcher/golang_web_server_boilerplate/middleware"
)

// This method creates a set of routes and returns it as an http handler function. Added as a method to server(server is in server.go)
func (server *Server) routes() http.Handler {
	router := chi.NewRouter()

	// Global middlewares
	router.Use(middleware.RateLimiter(100, 1))
	router.Use(middleware.CORS())
	router.Use(middleware.MaxBodySizeMiddleware(10))

	router.Handle(
		"/uploads/*",
		http.StripPrefix(
			"/uploads/",
			http.FileServer(http.Dir("./uploads")),
		),
	)

	// Route not found
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		responses.RespondWithJSON(w, http.StatusNotFound, responses.ErrorResponse{
			Code:    http.StatusNotFound, // 404 status code
			Status:  "error",
			Message: "Route not found",
		})
	})

	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		responses.RespondWithJSON(w, http.StatusMethodNotAllowed, responses.ErrorResponse{
			Code:    http.StatusMethodNotAllowed, // 405 status code
			Status:  "error",
			Message: "Method not allowed",
		})
	})

	router.Route("/v1", func(router chi.Router) {
		// Public
		router.Get("/healthz", server.HttpxHandler.Handle(server.CommonHandler.HealthCheck))
		router.Get("/error", server.HttpxHandler.Handle(server.CommonHandler.ErrorResponse))

		router.Route("/auth", func(router chi.Router) {
			router.Post("/register", server.HttpxHandler.Handle(server.Authhandler.UserRegistration))
			router.Post("/login", server.HttpxHandler.Handle(server.Authhandler.UserLogin))
			router.Get("/refresh-token", server.HttpxHandler.Handle(server.Authhandler.RefreshToken))
		})

		router.Group(func(router chi.Router) {
			router.Use(server.AuthMiddleware.Authenticate)

			router.Route("/users", func(router chi.Router) {
				router.Get("/", server.HttpxHandler.Handle(server.UsersHandler.GetAll))
				router.Post("/", server.HttpxHandler.Handle(server.UsersHandler.Create))
				router.Get("/{id}", server.HttpxHandler.Handle(server.UsersHandler.GetByID))
				router.Patch("/{id}", server.HttpxHandler.Handle(server.UsersHandler.Update))
				router.Delete("/{id}", server.HttpxHandler.Handle(server.UsersHandler.Delete))
			})

			posts.RegisterRoutes(router, server.PostsHandler, server.HttpxHandler, server.ACLMiddleware)

			router.Route("/posts-with-user", func(router chi.Router) {
				router.Get("/", server.HttpxHandler.Handle(server.PostsWithUserHandler.GetAll))
				router.Post("/", server.HttpxHandler.Handle(server.PostsWithUserHandler.Create))
				router.Get("/{id}", server.HttpxHandler.Handle(server.PostsWithUserHandler.GetByID))
			})

			router.Route("/acl", func(router chi.Router) {
				router.Route("/roles", func(router chi.Router) {
					router.Get("/", server.HttpxHandler.Handle(server.RolesHandler.GetAll))
					router.Post("/", server.HttpxHandler.Handle(server.RolesHandler.Create))
					router.Get("/{id}", server.HttpxHandler.Handle(server.RolesHandler.GetByID))
					router.Patch("/{id}", server.HttpxHandler.Handle(server.RolesHandler.Update))
					router.Delete("/{id}", server.HttpxHandler.Handle(server.RolesHandler.Delete))
				})

				router.Route("/modules", func(router chi.Router) {
					router.Get("/", server.HttpxHandler.Handle(server.ModulesHandler.GetAll))
					router.Post("/", server.HttpxHandler.Handle(server.ModulesHandler.Create))
					router.Get("/{id}", server.HttpxHandler.Handle(server.ModulesHandler.GetByID))
				})

				router.Route("/permissions", func(router chi.Router) {
					router.Get("/", server.HttpxHandler.Handle(server.PermissionsHandler.GetAll))
					router.Post("/", server.HttpxHandler.Handle(server.PermissionsHandler.Create))
					router.Get("/{id}", server.HttpxHandler.Handle(server.PermissionsHandler.GetByID))
				})

				router.Route("/user-roles", func(router chi.Router) {
					router.Get("/", server.HttpxHandler.Handle(server.UserRoleHandler.GetAll))
					router.Get("/{id}", server.HttpxHandler.Handle(server.UserRoleHandler.GetByID))
					router.Post("/", server.HttpxHandler.Handle(server.UserRoleHandler.Create))
					router.Get("/users", server.HttpxHandler.Handle(server.UserRoleHandler.GetAllUsersWithUserRoles))
					router.Get("/users/{userID}", server.HttpxHandler.Handle(server.UserRoleHandler.GetByUserID))
					router.Delete("/{id}", server.HttpxHandler.Handle(server.UserRoleHandler.Delete))
					router.Delete("/roles/{roleID}/permissions/{permissionID}", server.HttpxHandler.Handle(server.UserRoleHandler.DeleteByUserIDAndRoleID))
				})

				router.Route("/role-permissions", func(router chi.Router) {
					router.Get("/", server.HttpxHandler.Handle(server.RolePermissionsHandler.GetAll))
					router.Get("/{id}", server.HttpxHandler.Handle(server.RolePermissionsHandler.GetByID))
					router.Post("/", server.HttpxHandler.Handle(server.RolePermissionsHandler.Create))
					router.Get("/users", server.HttpxHandler.Handle(server.RolePermissionsHandler.GetAllUsersWithRolePermissions))
					router.Get("/users/{userID}", server.HttpxHandler.Handle(server.RolePermissionsHandler.GetByUserID))
					router.Delete("/{id}", server.HttpxHandler.Handle(server.RolePermissionsHandler.Delete))
					router.Delete("/roles/{roleID}/permissions/{permissionID}", server.HttpxHandler.Handle(server.RolePermissionsHandler.DeleteByRoleIDAndPermissionID))
				})
			})
		})
	})

	return router
}
