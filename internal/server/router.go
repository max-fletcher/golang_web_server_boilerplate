package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/max-fletcher/golang_web_server_boilerplate/middleware"
)

// This func creates a server and returns it as an http handler function. Added as a method to server(server is in server.go)
func (server *Server) routes() http.Handler {
	router := chi.NewRouter()

	// Global middlewares
	router.Use(middleware.RateLimiter(100, 1))
	router.Use(middleware.CORS())

	router.Route("/v1", func(v1router chi.Router) {
		// Public
		v1router.Get("/healthz", server.Handler.HandleReadiness)
		v1router.Get("/error", server.Handler.HandlerError)

		v1router.Post("/users", server.Handler.HandlerCreateUser)

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

	// v1router.Get("/healthz", server.Handler.HandleReadiness)
	// v1router.Get("/error", server.Handler.HandlerError)

	// v1router.Post("/users", server.handleCreateUser)
	// v1router.Get("/users", server.handleGetUserByAPIKey)

	// v1router.Post("/feeds", server.handleCreateFeed)
	////////////////////////
	// v1router.Get("/healthz", handlerReadiness)
	// v1router.Get("/error", handlerError)

	// v1router.Post("/users", apiCfg.handlerCreateUser)                                     // route for creating users in DB
	// v1router.Get("/users", apiCfg.authenticatedMiddleware(apiCfg.handlerGetUserByAPIKey)) // route for getting user by apiKey header in DB

	// v1router.Get("/posts", apiCfg.authenticatedMiddleware(apiCfg.handlerGetPostsForUser)) // route for getting all posts for a feed if the user is following that feed
	// router.Mount("/v1", v1router)

	// router.Mount("/v1", v1router)

	return router
}
