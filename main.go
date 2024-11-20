package main

import (
	"GIN/db"
	"context"
	"log"
	"os"

	"GIN/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	queries := db.New(conn)
	handler := handlers.NewEmployeeHandler(queries)
	r := gin.Default()

	r.GET("/", handler.ListEmployees)
	r.GET("/employee/:id", handler.GetEmployee)
	r.POST("/employee", handler.CreateEmployee)
	r.PUT("/employee/:id", handler.UpdateEmployee)
	r.DELETE("/employee/:id", handler.DeleteEmployee)

	r.Run()
}
