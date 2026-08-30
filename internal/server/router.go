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
		router.Get("/healthz", server.CommonHandler.HealthCheck)
		router.Get("/error", server.CommonHandler.ErrorResponse)

		router.Route("/users", func(router chi.Router) {
			router.Get("/", server.UsersHandler.GetAll)
			router.Post("/", server.UsersHandler.Create)
			router.Get("/{id}", server.UsersHandler.GetByID)
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

	// V1 API
	// v1router := chi.NewRouter()

	// v1router.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Hello World!"))
	// })

	// v1router.Get("/healthz", server.Handler.HealthCheck)
	// v1router.Get("/error", server.Handler.ErrorResponse)

	// v1router.Post("/users", server.handleCreateUser)
	// v1router.Get("/users", server.handleGetUserByAPIKey)

	// v1router.Post("/feeds", server.handleCreateFeed)
	////////////////////////
	// v1router.Get("/healthz", handlerReadiness)
	// v1router.Get("/error", ErrorResponse)

	// v1router.Post("/users", apiCfg.CreateUser)                                     // route for creating users in DB
	// v1router.Get("/users", apiCfg.authenticatedMiddleware(apiCfg.handlerGetUserByAPIKey)) // route for getting user by apiKey header in DB

	// v1router.Get("/posts", apiCfg.authenticatedMiddleware(apiCfg.handlerGetPostsForUser)) // route for getting all posts for a feed if the user is following that feed
	// router.Mount("/v1", v1router)

	// router.Mount("/v1", v1router)

	return router
}
