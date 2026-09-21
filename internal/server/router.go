package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/responses"
	"github.com/max-fletcher/golang_web_server_boilerplate/middleware"
)

// This func creates a server and returns it as an http handler function. Added as a method to server(server is in server.go)
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
		router.Get("/healthz", server.Handle(server.CommonHandler.HealthCheck))
		router.Get("/error", server.Handle(server.CommonHandler.ErrorResponse))

		router.Route("/auth", func(router chi.Router) {
			router.Post("/register", server.Handle(server.Authhandler.UserRegistration))
			router.Post("/login", server.Handle(server.Authhandler.UserLogin))
			router.Get("/refresh-token", server.Handle(server.Authhandler.RefreshToken))
		})

		router.Group(func(router chi.Router) {
			router.Use(server.AuthMiddleware.Authenticate)

			router.Route("/users", func(router chi.Router) {
				router.Get("/", server.Handle(server.UsersHandler.GetAll))
				router.Post("/", server.Handle(server.UsersHandler.Create))
				router.Get("/{id}", server.Handle(server.UsersHandler.GetByID))
				router.Patch("/{id}", server.Handle(server.UsersHandler.Update))
				router.Delete("/{id}", server.Handle(server.UsersHandler.Delete))
			})

			router.Route("/posts", func(router chi.Router) {
				router.Get("/", server.Handle(server.PostsHandler.GetAll))
				router.Post("/", server.Handle(server.PostsHandler.Create))
				router.Get("/{id}", server.Handle(server.PostsHandler.GetByID))
				router.Patch("/{id}", server.Handle(server.PostsHandler.Update))
				router.Delete("/{id}", server.Handle(server.PostsHandler.Delete))
			})

			router.Route("/posts-with-user", func(router chi.Router) {
				router.Get("/", server.Handle(server.PostsWithUserHandler.GetAll))
				router.Post("/", server.Handle(server.PostsWithUserHandler.Create))
				router.Get("/{id}", server.Handle(server.PostsWithUserHandler.GetByID))
			})

			router.Route("/acl", func(router chi.Router) {
				router.Route("/roles", func(router chi.Router) {
					router.Get("/", server.Handle(server.RolesHandler.GetAll))
					router.Post("/", server.Handle(server.RolesHandler.Create))
					router.Get("/{id}", server.Handle(server.RolesHandler.GetByID))
					router.Patch("/{id}", server.Handle(server.RolesHandler.Update))
					router.Delete("/{id}", server.Handle(server.RolesHandler.Delete))
				})

				router.Route("/modules", func(router chi.Router) {
					router.Get("/", server.Handle(server.ModulesHandler.GetAll))
					router.Post("/", server.Handle(server.ModulesHandler.Create))
					router.Get("/{id}", server.Handle(server.ModulesHandler.GetByID))
				})

				router.Route("/permissions", func(router chi.Router) {
					router.Get("/", server.Handle(server.PermissionsHandler.GetAll))
					router.Post("/", server.Handle(server.PermissionsHandler.Create))
					router.Get("/{id}", server.Handle(server.PermissionsHandler.GetByID))
				})

				router.Route("/role-permissions", func(router chi.Router) {
					router.Get("/", server.Handle(server.RolePermissionsHandler.GetAll))
					router.Get("/{id}", server.Handle(server.RolePermissionsHandler.GetByID))
					router.Post("/", server.Handle(server.RolePermissionsHandler.Create))
					router.Get("/users", server.Handle(server.RolePermissionsHandler.GetAllUsersWithRolePermissions))
					router.Get("/users/{userID}", server.Handle(server.RolePermissionsHandler.GetByUserID))
					router.Delete("/{id}", server.Handle(server.RolePermissionsHandler.Delete))
					router.Delete("/roles/{roleID}/permissions/{permissionID}", server.Handle(server.RolePermissionsHandler.DeleteByRoleIDAndPermissionID))
				})

				router.Route("/user-roles", func(router chi.Router) {
					router.Get("/", server.Handle(server.UserRoleHandler.GetAll))
					router.Get("/{id}", server.Handle(server.UserRoleHandler.GetByID))
					router.Post("/", server.Handle(server.UserRoleHandler.Create))
					router.Get("/users", server.Handle(server.UserRoleHandler.GetAllUsersWithUserRoles))
					router.Get("/users/{userID}", server.Handle(server.UserRoleHandler.GetByUserID))
					router.Delete("/{id}", server.Handle(server.UserRoleHandler.Delete))
					router.Delete("/roles/{roleID}/permissions/{permissionID}", server.Handle(server.UserRoleHandler.DeleteByUserIDAndRoleID))
				})
			})
		})
	})

	return router
}
