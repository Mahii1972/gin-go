package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

func main() {
	// Get command from command-line arguments
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go up/down")
	}
	command := os.Args[1]

	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("Error loading .env file")
	}

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Get a *sql.DB connection instead of the connection string
	sqlDB, err := sql.Open("pgx", pool.Config().ConnConfig.ConnString())
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	// Run migrations based on the command
	var migrationErr error
	switch command {
	case "up":
		migrationErr = goose.Up(sqlDB, "../../db/migrations")
	case "down":
		migrationErr = goose.Down(sqlDB, "../../db/migrations")
	default:
		log.Fatalf("Invalid command: %s. Use 'up' or 'down'", command)
	}

	if migrationErr != nil {
		log.Fatalf("Migration error: %v", migrationErr)
	}

	log.Printf("Migrations (%s) completed successfully", command)
}
