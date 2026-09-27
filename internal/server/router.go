package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/responses"
	modules_names "github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/acl/module-names"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/acl/permissions"
	role_permissions "github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/acl/role-permissions"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/acl/roles"
	user_roles "github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/acl/user-roles"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/auth"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/posts"
	posts_with_users "github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/posts-with-users"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/users"
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

		auth.RegisterRoutes(router, server.Authhandler, server.HttpxHandler)
		// The line above is condensed/grouped, but you can use the lines below(kept as example) if you need simplicity instead since these
		// are public routes and don't need ACL
		// router.Route("/auth", func(router chi.Router) {
		// 	router.Post("/register", server.HttpxHandler.Handle(server.Authhandler.UserRegistration))
		// 	router.Post("/login", server.HttpxHandler.Handle(server.Authhandler.UserLogin))
		// 	router.Get("/refresh-token", server.HttpxHandler.Handle(server.Authhandler.RefreshToken))
		// })

		router.Group(func(router chi.Router) {
			router.Use(server.AuthMiddleware.Authenticate)

			users.RegisterRoutes(router, server.UsersHandler, server.HttpxHandler, server.ACLMiddleware)
			posts.RegisterRoutes(router, server.PostsHandler, server.HttpxHandler, server.ACLMiddleware)
			posts_with_users.RegisterRoutes(router, server.PostsWithUserHandler, server.HttpxHandler, server.ACLMiddleware)

			router.Route("/acl", func(router chi.Router) {
				roles.RegisterRoutes(router, server.RolesHandler, server.HttpxHandler, server.ACLMiddleware)
				modules_names.RegisterRoutes(router, server.ModulesHandler, server.HttpxHandler, server.ACLMiddleware)
				permissions.RegisterRoutes(router, server.PermissionsHandler, server.HttpxHandler, server.ACLMiddleware)
				user_roles.RegisterRoutes(router, server.UserRoleHandler, server.HttpxHandler, server.ACLMiddleware)
				role_permissions.RegisterRoutes(router, server.RolePermissionsHandler, server.HttpxHandler, server.ACLMiddleware)
			})
		})
	})

	return router
}
