.PHONY: help install-tools build run test deps

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install-tools: ## Install required tools (goose, sqlc)
	@echo "Installing goose..."
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@echo "Installing sqlc..."
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@echo "Tools installed successfully!"

build: ## Build the application
	@echo "Building application..."
	@go build -o bin/api cmd/main.go
	@echo "Build completed successfully!"

run: ## Run the application
	@echo "Running application..."
	@go run cmd/main.go http-server

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Tests completed successfully!"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies downloaded successfully!"

.DEFAULT_GOAL := help
