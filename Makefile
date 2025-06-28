# Makefile for FeedApp

.PHONY: swagger build test clean

# Swagger documentation generation
swagger:
	swag init -g cmd/server/main.go

# Build the application
build:
	go build -o bin/feedapp ./cmd/server

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Install dependencies
install:
	go mod download
	go install github.com/swaggo/swag/cmd/swag@latest

# Run the application in development mode
dev:
	docker compose up --build

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Docker commands
docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-rebuild:
	docker compose down
	docker compose up --build -d

# Help
help:
	@echo "Available commands:"
	@echo "  swagger        - Generate Swagger documentation"
	@echo "  build          - Build the application"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  fmt            - Format code"
	@echo "  lint           - Lint code"
	@echo "  install        - Install dependencies"
	@echo "  dev            - Run in development mode"
	@echo "  clean          - Clean build artifacts"
	@echo "  docker-up      - Start Docker containers"
	@echo "  docker-down    - Stop Docker containers"
	@echo "  docker-rebuild - Rebuild and restart Docker containers"
	@echo "  help           - Show this help"