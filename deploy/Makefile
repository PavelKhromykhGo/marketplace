MODULE := $(shell go list -m)
MAIN_PKG := ./cmd/api/main.go
MIGRATIONS_DIR := ./migrations

DATABASE_URL ?= host=localhost port=5432 user=postgres password=postgres dbname=marketplace sslmode=disable

SWAG := go run github.com/swaggo/swag/cmd/swag@latest
GOOSE := go run github.com/pressly/goose/v3/cmd/goose@latest

.PHONY: help deps tidy fmt build run test cover swag docs migrate-up migrate-down migrate-status migrate-reset migrate-create

help:
	@echo "Makefile commands:"
	@echo "  deps            - Install project dependencies"
	@echo "  tidy           - Tidy up go.mod and go.sum files"
	@echo "  fmt             - Format the code using gofmt"
	@echo "  build           - Build the application"
	@echo "  run             - Run the application (uses DATABASE_URL)"
	@echo "  test            - Run tests"
	@echo "  cover           - Run tests with coverage report"
	@echo "  swag            - Generate Swagger documentation ./docs"
	@echo "  docs            - Open Swagger UI in the browser"
	@echo "  migrate-up      - Apply all up migrations"
	@echo "  migrate-down    - Apply all down migrations"
	@echo "  migrate-status  - Show current migration status"
	@echo "  migrate-reset   - Reset the database and apply all migrations"
	@echo "  migrate-create  - Create a new migration file"

deps:
	go mod download
tidy:
	go mod tidy
fmt:
	go fmt ./...
build:
	go build -o bin/marketplace $(MAIN_PKG)
run:
	DATABASE_URL="$(DATABASE_URL)" go run $(MAIN_PKG)
test:
	go test ./... -v
cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"
swag:
	$(SWAG) init \
	--generalInfo $(MAIN_PKG)/main.go \
	--dir $(MAIN_PKG), ./internal \
	--output ./docs

docs:
	@echo "Opening Swagger UI at http://localhost:8080/docs/index.html"
	@echo "Make sure your application is running."
	@xdg-open http://localhost:8080/docs/index.html || open http://localhost:8080/docs/index.html

migrate-up:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" up
migrate-down:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" down
migrate-status:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" status
migrate-reset:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" reset
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" up
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: Please provide a name for the migration. Example: make migrate-create name=add_users_table"; \
		exit 1; \
	fi
	$(GOOSE) -dir $(MIGRATIONS_DIR) create "$(name)" sql
	@echo "Created new migration '$(name)' in $(MIGRATIONS_DIR)"