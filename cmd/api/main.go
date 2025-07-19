package main

import (
	"log"
	"marketplace/internal/product"
	"marketplace/internal/repository/postgres"
	"marketplace/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=marketplace sslmode=disable"
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Printf("Failed to close database connection: %v", err)
		}
	}()
	// Initialize repositories, services, and handlers here
	repo := postgres.NewProductRepository(db)
	svc := product.NewService(repo)
	h := product.NewHandler(svc)

	r := gin.New()

	h.RegisterRoutes(r)

	r.Use(middleware.ErrorLogger())
	r.Use(gin.Recovery())

	log.Fatal(r.Run(":8080"))
}
