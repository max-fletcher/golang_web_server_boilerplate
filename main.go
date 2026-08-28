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

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/server"
)

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

	DBConn := db.New(conn)    // connection database to sqlc's queries
	srv := server.New(DBConn) // Server struct coming from server.go

	// Server options like router and port
	// On windows, to run without compiling the server, use "go run ."
	// On windows, to compile(for prod) and run binaries, use "go build -o {filename}.exe" then ".\{filename}.exe"
	log.Printf("Server startng on port %v", portString)
	err = http.ListenAndServe(":"+portString, srv.Router) // start/initialize server using router coming from server.go
	if err != nil {                                       // throws an error if the server fails
		log.Fatal(err)
	}
	log.Printf("Server started on port %v", portString)
}
