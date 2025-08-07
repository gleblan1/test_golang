.PHONY: help up down logs migrate-up migrate-down test build clean run-api run-worker db-up db-down

# Default target
help:
	@echo "Available commands:"
	@echo "  up          - Start all services with docker-compose"
	@echo "  down        - Stop all services"
	@echo "  logs        - View logs from all services"
	@echo "  migrate-up  - Apply database migrations"
	@echo "  migrate-down- Rollback database migrations"
	@echo "  test        - Run tests"
	@echo "  build       - Build the application"
	@echo "  clean       - Clean build artifacts"
	@echo "  run-api     - Run API server locally"
	@echo "  run-worker  - Run worker locally"
	@echo "  db-up       - Start PostgreSQL only"
	@echo "  db-down     - Stop PostgreSQL"

# Docker Compose commands
up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

# Database commands
db-up:
	docker-compose up -d postgres

db-down:
	docker-compose stop postgres

# Migration commands
migrate-up:
	docker-compose run --rm migrate

migrate-down:
	docker-compose run --rm migrate -path /migrations -database "postgres://postgres:password@postgres:5432/crypto_tracker?sslmode=disable" down

# Development commands
run-api:
	go run cmd/api/main.go

run-worker:
	go run cmd/worker/main.go

# Build commands
build:
	go build -o bin/api cmd/api/main.go
	go build -o bin/worker cmd/worker/main.go

clean:
	rm -rf bin/
	go clean

# Test commands
test:
	go test ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Install dependencies
deps:
	go mod download
	go mod tidy

# Generate swagger docs
swagger:
	swag init -g cmd/api/main.go -o docs 