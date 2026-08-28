.PHONY: help manage/dev-environment install/sqlc generate/sqlc install/dependencies install/tools install/golangci-lint lint tests/unit-tests tests/functional-tests tests/coverage-view tests/coverage-analyze build

SQLC_VERSION ?= v1.31.1
SWAGGER_VERSION ?= v1.16.6
GOLANGCI_LINT_VERSION ?= v1.64.2

TOOLS_BIN ?= $(CURDIR)/.bin
SQLC_BIN ?= $(TOOLS_BIN)/sqlc
SWAGGER_BIN ?= $(TOOLS_BIN)/swag
GOLANGCI_LINT_BIN ?= $(TOOLS_BIN)/golangci-lint

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/\-]+:.*?## / {printf "  %-30s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

manage/dev-environment: ## Manage development environment
	@echo "Go into k8s folder and run make"

install/dependencies: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies downloaded successfully!"

install/sqlc: ## Install sqlc
	@set -e; \
	current=""; \
	if [ -x "$(SQLC_BIN)" ]; then \
		current="$$( $(SQLC_BIN) version 2>/dev/null || true )"; \
	fi; \
	if [ "$$current" != "$(SQLC_VERSION)" ]; then \
		echo "Installing sqlc $(SQLC_VERSION)..."; \
		mkdir -p $(TOOLS_BIN); \
		GOBIN="$(TOOLS_BIN)" go install github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION); \
		echo "sqlc $(SQLC_VERSION) installed successfully!"; \
	fi

install/swagger:
	@set -e; \
	if [ ! -x "$(SWAGGER_BIN)" ] || ! "$(SWAGGER_BIN)" version 2>/dev/null | grep -q "$(SWAGGER_VERSION)"; then \
		echo "Installing swagger $(SWAGGER_VERSION)..."; \
		mkdir -p $(TOOLS_BIN); \
		GOBIN="$(TOOLS_BIN)" go install github.com/swaggo/swag/cmd/swag@$(SWAGGER_VERSION); \
		echo "swagger $(SWAGGER_VERSION) installed successfully!"; \
	fi

generate/sqlc: install/sqlc ## Generate code using sqlc
	@echo "Generating code using sqlc..."
	@$(SQLC_BIN) generate
	@echo "Code generated successfully!"

generate/swagger: ## Generate Swagger documentation
	@echo "Generating Swagger documentation..."
	@$(SWAGGER_BIN) init -g internal/adapters/chi-server/router.go -o docs
	@echo "Swagger documentation generated successfully!"

install/tools: ## Install required tools (goose)
	@echo "Installing goose..."
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@echo "Tools installed successfully!"

install/golangci-lint: ## Install golangci-lint
	@set -e; \
	current=""; \
	if [ -x "$(GOLANGCI_LINT_BIN)" ]; then \
		current="$$( $(GOLANGCI_LINT_BIN) version 2>/dev/null | grep -Eo '[0-9]+\.[0-9]+\.[0-9]+' | head -1 )"; \
	fi; \
	if [ "$$current" != "$(GOLANGCI_LINT_VERSION:v%=%)" ]; then \
		echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION)..."; \
		mkdir -p $(TOOLS_BIN); \
		GOBIN="$(TOOLS_BIN)" go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION); \
		echo "golangci-lint $(GOLANGCI_LINT_VERSION) installed successfully!"; \
	fi

lint: install/golangci-lint ## Run golangci-lint
	@echo "Running golangci-lint..."
	@$(GOLANGCI_LINT_BIN) run ./...
	@echo "Lint completed successfully!"

tests/unit-tests: ## Run unit tests
	@echo "Running unit tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Tests completed successfully!"

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

.DEFAULT_GOAL := help
