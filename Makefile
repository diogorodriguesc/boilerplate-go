.PHONY: help build/docker-dev-image install/dependencies install/tools tests/unit-tests tests/functional-tests tests/coverage-view tests/coverage-analyze build run

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/\-]+:.*?## / {printf "  %-30s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build/docker-dev-image: ## Build the development Docker image
	@echo "Building development Docker image..."
	@docker buildx build --load -t micro-app-boilerplate-go:dev -f Dockerfile.dev .

install/dependencies: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies downloaded successfully!"

install/tools: ## Install required tools (goose, sqlc)
	@echo "Installing goose..."
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@echo "Installing sqlc..."
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@echo "Tools installed successfully!"

tests/unit-tests: ## Run unit tests
	@echo "Running unit tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Tests completed successfully!"

tests/functional-test: ## Run functional tests
	@echo "Running functional tests..."
	@go test -tags=functional -v ./internal/adapters/chi-server/... -count=1

tests/functional-tests: ## Run functional tests
	@echo "Running functional tests..."
	@go test -tags=functional -v ./...

tests/coverage-view: ## Get code coverage
	@echo "Running tests and checking coverage"
	@go test -coverprofile=cover.out ./...
	@COVERAGE=$$(go tool cover -func=cover.out | grep total | grep -Eo '[0-9]+\.[0-9]+'); \
	echo "Current coverage: $$COVERAGE";

tests/coverage-analyze: ## Analyze code coverage
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

build: ## Build the application
	@echo "Building application..."
	@go build -o bin/api cmd/main.go
	@echo "Build completed successfully!"

run: ## Run the application
	@echo "Running application..."
	@go run cmd/main.go http-server

.DEFAULT_GOAL := help
