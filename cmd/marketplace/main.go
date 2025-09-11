// @title Marketplace API
// @version 1.0
// @description This is a sample server for a marketplace application.
// @BasePath /

//securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"context"
	"embed"
	"errors"
	"log"
	"marketplace/internal/auth"
	"marketplace/internal/logger"
	"marketplace/internal/product"
	"marketplace/internal/repository/postgres"
	"marketplace/internal/transport"
	"marketplace/internal/user"
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
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	httpAddr := env("HTTP_ADDR", ":8080")

	dsn := env("DATABASE_URL", "host=localhost port=5432 user=postgres password=postgres dbname=marketplace sslmode=disable")

	migDir := env("MIGRATIONS_DIR", "./migrations")

	if s := env("JWT_SECRET", "your-256-bit-secret"); s != "" {
		auth.SetSecret([]byte(s))
	}

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
	if err = runMigrations(db, migDir); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories, services, and handlers here
	prodRepo := postgres.NewProductRepository(db)
	userRepo := postgres.NewUserRepository(db)

	prodService := product.NewService(prodRepo)
	userService := user.NewService(userRepo)

	logg, err := logger.New(logger.Config{Enviroment: os.Getenv("APP_ENV")})
	if err != nil {
		panic(err)
	}
	defer logg.Sync()

	r := gin.New()

	r.Use(
		middleware.RequestID(),
		middleware.ZapRecovery(),
		middleware.ZapLogger(),
		middleware.ErrorHandler(),
		middleware.Metrics(),
	)

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ok":      true,
			"swagger": "/swagger/index.html",
			"health":  "/healthz",
			"ready":   "/readyz",
		})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	transport.Health{DB: db}.Register(r)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	product.RegisterRoutes(r, prodService)
	user.RegisterRoutes(r, userService)

	srv := &http.Server{
		Addr:              httpAddr,
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
