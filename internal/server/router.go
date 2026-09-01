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

		router.Route("/users", func(router chi.Router) {
			router.Get("/", server.Handle(server.UsersHandler.GetAll))
			router.Post("/", server.Handle(server.UsersHandler.Create))
			router.Get("/{id}", server.Handle(server.UsersHandler.GetByID))
		})

		// Authenticated
		// v1router.Group(func(v1router chi.Router) {
		// 	v1router.Use(middleware.Authenticated)

		// 	v1router.Get("/users", s.Handler.HandleGetUser)
		// 	v1router.Post("/feeds", s.Handler.HandleCreateFeed)
		// 	v1router.Get("/feeds", s.Handler.HandleGetFeeds)
		// 	v1router.Post("/feed-follows", s.Handler.HandleCreateFeedFollow)
		// })

		// v1router.Route("/users", func(r chi.Router) {
		// 	v1router.Post("/", s.Handler.HandleCreateUser)
		// 	v1router.Get("/", s.Handler.HandleGetUser)
		// })
	})

	return router
}
