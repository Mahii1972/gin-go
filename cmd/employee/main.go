package main

import (
	"GIN/db"
	"context"
	"log"
	"os"

	"GIN/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("Error loading .env file")
	}

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	queries := db.New(pool)
	handler := handlers.NewEmployeeHandler(queries)
	r := gin.Default()

	r.GET("/", handler.ListEmployees)
	r.GET("/employee/:id", handler.GetEmployee)
	r.POST("/employee", handler.CreateEmployee)
	r.PUT("/employee/:id", handler.UpdateEmployee)
	r.DELETE("/employee/:id", handler.DeleteEmployee)

	r.Run()
}
