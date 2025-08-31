// @title Marketplace API
// @version 1.0
// @description This is a sample server for a marketplace application.
// @BasePath /

package main

import (
	"context"
	"embed"
	"errors"
	"flag"
	"log"
	"marketplace/internal/product"
	"marketplace/internal/repository/postgres"
	"marketplace/middleware"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	_ "marketplace/docs"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/pressly/goose/v3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var embedMigrations embed.FS

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func runMigrations(db *sqlx.DB, dir string) error {
	// мигрируем из файловой системы
	if abs, err := filepath.Abs(dir); err == nil {
		if st, err := os.Stat(abs); err == nil && st.IsDir() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			log.Printf("Running migrations from %s", abs)
			return goose.UpContext(ctx, db.DB, abs)
		}
	}
	// мигрируем из embed.FS
	goose.SetBaseFS(embedMigrations)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Printf("Running migrations from embedded FS")
	return goose.UpContext(ctx, db.DB, "migrations")
}

func main() {
	migrationsPathFlag := flag.String("migrations", env("MIGRATIOINS_DIR", "./migrations"), "path to migrations directory")
	flag.Parse()

	dsn := env("DATABASE_URL", "host=localhost port=5432 user=postgres password=postgres dbname=marketplace sslmode=disable")

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Printf("Failed to close database connection: %v", err)
		}
	}()
	if err = goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Failed to set goose dialect: %v", err)
	}
	if err = runMigrations(db, *migrationsPathFlag); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories, services, and handlers here
	repo := postgres.NewProductRepository(db)
	svc := product.NewService(repo)
	h := product.NewHandler(svc)

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.ErrorLogger())
	r.Use(middleware.ErrorHandler())

	h.RegisterRoutes(r)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("Starting server on port %s", srv.Addr)
		if err = srv.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			log.Fatalf("listen: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
