package main

import (
	"context"
	"database/sql"
	"flag"
	"log"

	"github.com/joho/godotenv"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/config"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/seeder"
)

func main() {
	// Load config
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Using Load function to fetch any variables we need from .env(all of them are stored inside cfg struct)
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatal("Can't connect to the database", err)
	}
	defer conn.Close()

	err = conn.Ping()
	if err != nil {
		log.Fatal("Database not reachable:", err)
	}

	database := db.New(conn)
	ctx := context.Background()

	seedType := flag.String(
		"seed",
		"",
		"What to seed: modules, permissions, roles, or all",
	)

	flag.Parse()

	switch *seedType {
	case "roles":
		if err := seeder.Roles(ctx, database); err != nil {
			log.Fatal("Failed to seed roles", err)
		}

	case "modules":
		if err := seeder.Modules(ctx, database); err != nil {
			log.Fatal("Failed to seed modules", err)
		}

	case "permissions":
		if err := seeder.Permissions(ctx, database); err != nil {
			log.Fatal("Failed to seed permissions", err)
		}

	case "role-permissions":
		if err := seeder.RolePermissions(ctx, database); err != nil {
			log.Fatal("Failed to seed role-permissions", err)
		}

	case "create-superadmin":
		if err := seeder.CreateSuperAdmin(ctx, database); err != nil {
			log.Fatal("Failed to seed superadmin", err)
		}

	case "all":
		if err := seeder.All(ctx, database); err != nil {
			log.Fatal("Failed to seed all data: ", err)
		}

	default:
		log.Fatal("invalid --seed value")
	}
}
