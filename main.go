package main

// # IMPORTING PACKAGES:
// To import 3rd party modules for go, import and use the package in your file then use
// "go mod tidy"(cleanup unused deps) and then "go mod vendor"(copy only what’s needed).
// Although some packages like "github.com/google/uuid" and "github.com/lib/pq"
// requires the "go get" command to install.
// Here are some other useful commands:
// "go list -m all" to see a list of all packages
// "go mod why github.com/go-chi/chi/v5" to check why a package exists
// "go install github.com/some/tool@latest"; only use "get" when installing binaries
// Remember that this is not the same as downloading packages like "github.com/sqlc-dev/sqlc/cmd/sqlc@latest" or
// "go install github.com/pressly/goose/v3/cmd/goose@latest" that we have
// to use the command "go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest" to get. The reason being that sqlc is a tool
// (a set of binaries to be specific) that are not part of our application i.e nothing is downloaded into vendor and
// when we use "go build -o {filename}.exe", it is not compiled into the application.
import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
	"github.com/max-fletcher/golang_web_server_boilerplate/middleware"
)

// Reference to a DB connection. Will be used to query data form DB.
type apiConfig struct {
	DB *db.Queries
}

func main() {
	fmt.Println("Web Server made in GO. Starting...")

	// We are using this because we need to set a bridge, When we use "os.Getenv" it only checks the os's
	// env variables and not the .env file. In golang .env is not automatically loaded into the os's
	// environment variables.
	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Port is not found in the .env file")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in env")
	}

	// using go's sql package from its standard library to establish connection
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to the database", err)
	}
	err = conn.Ping()
	if err != nil {
		log.Fatal("DB not reachable:", err)
	}
	apiCfg := apiConfig{
		DB: db.New(conn),
	}

	fmt.Printf("PORT: %v", portString)

	router := chi.NewRouter()
	router.Use(middleware.CORS())              // Rate limiter
	router.Use(middleware.RateLimiter(100, 1)) // CORS

	// router.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Hello World!"))
	// })

	// Making a sub-router
	v1router := chi.NewRouter()

	// Another way to use an auth middleware
	// v1router.Use(auth.AuthenticatedMiddleware(apiCfg.DB))
	// v1router.Get("/users2", apiCfg.handlerGetUserByAPIKey2) // route for getting user by apiKey header in DB

	v1router.Get("/healthz", handlerReadiness)

	v1router.Get("/error", handlerError)

	v1router.Post("/users", apiCfg.handlerCreateUser)                                     // route for creating users in DB
	v1router.Get("/users", apiCfg.authenticatedMiddleware(apiCfg.handlerGetUserByAPIKey)) // route for getting user by apiKey header in DB

	v1router.Get("/feeds", apiCfg.handlerGetFeeds)                                    // route for getting all feed in DB
	v1router.Post("/feeds", apiCfg.authenticatedMiddleware(apiCfg.handlerCreateFeed)) // route for creating a feed in DB

	v1router.Post("/feed_follows", apiCfg.authenticatedMiddleware(apiCfg.handlerCreateFeedFollow))                      // route for creating a feed follow in DB
	v1router.Get("/feed_follows_by_user", apiCfg.authenticatedMiddleware(apiCfg.handlerGetFeedFollowsByUserID))         // route for getting all feed follows in DB
	v1router.Delete("/feed_follows/{feedFollowId}", apiCfg.authenticatedMiddleware(apiCfg.handleDeleteFeedFollowsById)) // route for deleting feed follow for a user in DB

	v1router.Get("/posts", apiCfg.authenticatedMiddleware(apiCfg.handlerGetPostsForUser)) // route for getting all posts for a feed if the user is following that feed
	router.Mount("/v1", v1router)

	// Server options like router and port
	// On windows, to start the server, use "go build -o {filename}.exe" then ".\go-rss-agg.exe"
	server := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)
	err = server.ListenAndServe() // initialize server
	if err != nil {               // throws an error if the server fails
		log.Fatal(err)
	}
}
