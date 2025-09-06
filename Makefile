SHELL := /bin/bash

MODULE := $(shell go list -m)
MAIN_PKG := ./cmd/marketplace/main.go
MIGRATIONS_DIR := ./migrations

# Go env (хост)
GOPATH ?= $(HOME)/go
GOMODCACHE ?= $(GOPATH)/pkg/mod
GOCACHE ?= $(HOME)/.cache/go-build
GOTMPDIR ?= $(HOME)/.cache/go-tmp
GOBIN := $(GOPATH)/bin
export GOPATH
export GOMODCACHE
export GOCACHE
export GOTMPDIR
export GOBIN
export PATH := $(PATH):$(GOBIN)

DATABASE_URL ?= host=localhost port=5432 user=postgres password=postgres dbname=marketplace sslmode=disable

# swag (хост)
SWAG := $(GOBIN)/swag
$(SWAG):
	@echo "Installing swag..."
	@mkdir -p $(dir $(SWAG))
	@mkdir -p $(GOTMPDIR)
	go install github.com/swaggo/swag/cmd/swag@latest

HOST_PWD := $(shell pwd -W 2>/dev/null || pwd)

DOCKER_ENV := MSYS_NO_PATHCONV=1 MSYS2_ARG_CONV_EXCL="*"

DOCKER_GO := $(DOCKER_ENV) docker run --rm \
  -v "$(HOST_PWD):/app" \
  -w /app \
  -e GOPATH=/go \
  -e GOMODCACHE=/go/pkg/mod \
  -e GOCACHE=/root/.cache/go-build \
  golang:1.24-alpine

GOOSE := go run github.com/pressly/goose/v3/cmd/goose@latest

.PHONY: help deps tidy fmt build run test cover swag docs migrate-up migrate-down migrate-status migrate-reset migrate-create docker-build swag-docker

help:
	@echo "Makefile commands:"
	@echo "  deps            - Install project dependencies"
	@echo "  tidy            - Tidy up go.mod and go.sum files"
	@echo "  fmt             - Format the code using gofmt"
	@echo "  build           - Build the application"
	@echo "  run             - Run the application (uses DATABASE_URL)"
	@echo "  test            - Run tests"
	@echo "  cover           - Run tests with coverage report"
	@echo "  swag            - Generate Swagger documentation ./docs"
	@echo "  swag-docker     - Generate Swagger docs in Docker"
	@echo "  docs            - Open Swagger UI in the browser"
	@echo "  migrate-up      - Apply all up migrations"
	@echo "  migrate-down    - Apply all down migrations"
	@echo "  migrate-status  - Show current migration status"
	@echo "  migrate-reset   - Reset DB and apply all migrations"
	@echo "  migrate-create  - Create a new migration file"
	@echo "  docker-build    - Build Docker image for the application"

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

# Генерация док на хосте (ставит swag при необходимости)
swag: $(SWAG)
	$(SWAG) init \
		--generalInfo $(MAIN_PKG) \
		--dir ./,./internal \
		--output ./docs

# Генерация докеризованно (ничего на хост не ставит)
swag-docker:
	$(DOCKER_GO) bash -lc '\
	  go install github.com/swaggo/swag/cmd/swag@latest && \
	  /go/bin/swag init \
	    --generalInfo $(MAIN_PKG) \
	    --dir ./,./internal \
	    --output ./docs \
	'

docs:
	@echo "Opening Swagger UI at http://localhost:8080/docs/index.html"
	@echo "Make sure your application is running."
	@xdg-open http://localhost:8080/docs/index.html || open http://localhost:8080/docs/index.html

docker-build: swag
	docker build -t marketplace .

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
