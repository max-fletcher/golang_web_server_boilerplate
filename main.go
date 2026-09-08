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
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/config"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/queue"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/redis"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/server"
)

func main() {
	fmt.Println("Web Server made in GO. Starting...")

	// We are using this because we need to set a bridge, When we use "os.Getenv" it only checks the os's
	// env variables and not the .env file. In golang .env is not automatically loaded into the os's
	// environment variables.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// using Load function to fetch any variables we need from .env(all of them are stored inside cfg struct)
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// using go's sql package from its standard library to establish connection
	// *IMPORTANT: One terminology note: conn/sql.DB is technically a connection pool (*sql.DB), not one specific PostgreSQL connection.
	// BeginTx obtains a connection from that pool and starts the transaction on it.
	conn, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatal("Can't connect to the database", err)
	}
	defer conn.Close()

	err = conn.Ping()
	if err != nil {
		log.Fatal("Database not reachable:", err)
	}

	// Moved this here from server.go due to 2 reasons:
	// 1. The Redis client needs to live for the entire lifetime of the application. If it is inside NewServer, when it returns &server,
	//    the client will close immediately(with or without "defer redisClient.Close()") so it might not run at all. If redis client
	//    is created in main.go, as long as the application lives(with or without the server existing), the connection will live as well.
	// 2. main.go is the better owner of redisClient because redisClient is a long-lived infrastructure resource, and main.go is the
	//    natural place where your application creates, owns, and eventually shuts down those resources. redisClient represents a
	//    connection/pool to an external system (just like the database connection created here) and although several parts of your
	//    application depend on that same client, server.go is supposed to be responsible only for the HTTP-specific
	//    infrastructure(i.e server) of the application.
	redisClient, err := redis.NewClient(cfg.RedisURL)
	if err != nil {
		log.Fatal("Can't connect to Redis:", err)
	}
	defer redisClient.Close()

	// ------ Create a queue worker ------

	//create a handler function to pass to NewWorker(will be bound to Worker struct)
	workerHandler := func(ctx context.Context, message queue.Message) error {
		log.Printf(
			"WORKER RECEIVED: type=%s data=%v",
			message.Type,
			message.Data,
		)

		return nil
	}

	worker := queue.NewWorker( // create a queue worker
		redisClient,
		string(queue.StreamPost),
		"post-workers",
		"worker-1",
		workerHandler,
		slog.Default(),
	)

	go func() { // start the queue worker
		if err := worker.Run(context.Background()); err != nil {
			log.Fatal("Queue worker stopped:", err)
		}
	}()
	// ------ End create a queue worker ------

	database := db.New(conn)                                  // connecting database to sqlc's queries. "database" contains all sqlc queries.
	srv := server.NewServer(database, conn, cfg, redisClient) // Server struct coming from server.go. Create a new server instance

	// Server options like router and port
	// On windows, to run without compiling the server, use "go run ."
	// On windows, to compile(for prod) and run binaries, use "go build -o {filename}.exe" then ".\{filename}.exe"
	log.Printf("Server starting on port %v", cfg.Port)
	err = http.ListenAndServe(":"+cfg.Port, srv.Router) // start/initialize server using router coming from server.go
	if err != nil {                                     // throws an error if the server fails
		log.Fatal(err)
	}
}
