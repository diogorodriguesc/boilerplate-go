.PHONY: help minikube/start minikube/start-first-time minikube/stop docker/build docker/load apply delete db/migrate port-forward port-forward/http port-forward/grpc

IMAGE_NAME ?= micro-app-boilerplate-go
IMAGE_TAG ?= dev
IMAGE ?= $(IMAGE_NAME):$(IMAGE_TAG)

DOCKERFILE ?= $(CURDIR)/../Dockerfile.dev
BUILD_CONTEXT ?= $(CURDIR)/..

K8S_BASE ?= $(CURDIR)/base

HTTP_DEPLOYMENT ?= deployment/micro-app-boilerplate-go-http
HTTP_SERVICE ?= svc/go-http-k8s-service
HTTP_LOCAL_PORT ?= 8085
HTTP_SERVICE_PORT ?= 8080

GRPC_DEPLOYMENT ?= deployment/micro-app-boilerplate-go-grpc
GRPC_SERVICE ?= svc/go-grpc-k8s-service
GRPC_LOCAL_PORT ?= 9095
GRPC_SERVICE_PORT ?= 9090

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/\-]+:.*?## / {printf "  %-30s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

minikube/start-first-time: ## Start Minikube for the first time, mounting MINIKUBE_PROJECTS_PATH into /app
	@if [ -z "$(MINIKUBE_PROJECTS_PATH)" ]; then \
		echo "MINIKUBE_PROJECTS_PATH is not set. Export it before running this target."; \
		exit 1; \
	fi
	@minikube start --cpus 4 --memory 8192 --driver=kvm2 \
		--mount \
		--mount-string="$(MINIKUBE_PROJECTS_PATH):/app"

minikube/start: ## Start Minikube (after the first-time setup has already run)
	@minikube start

minikube/stop: ## Stop Minikube
	@minikube stop

docker/build: ## Build the development Docker image
	@echo "Building development Docker image..."
	@docker build -t $(IMAGE) -f $(DOCKERFILE) $(BUILD_CONTEXT)

docker/load: ## Load the development Docker image into Minikube
	@echo "Loading $(IMAGE) into Minikube..."
	@minikube image load $(IMAGE)

apply: ## Apply Kubernetes resources, wait for both rollouts, then port-forward the HTTP service to localhost:8085
	@echo "Applying Kubernetes resources..."
	@kubectl apply -k $(K8S_BASE)
	@kubectl rollout status $(HTTP_DEPLOYMENT)
	@kubectl rollout status $(GRPC_DEPLOYMENT)
	@$(MAKE) port-forward/http

delete: ## Delete Kubernetes resources
	@echo "Deleting Kubernetes resources..."
	@kubectl delete -k $(K8S_BASE)

port-forward: port-forward/http ## Alias for port-forward/http (Ctrl+C to stop)

port-forward/http: ## Port-forward the HTTP service to localhost:8085 (Ctrl+C to stop)
	@echo "Port-forwarding $(HTTP_SERVICE) localhost:$(HTTP_LOCAL_PORT) -> $(HTTP_SERVICE_PORT). Press Ctrl+C to stop."
	@kubectl port-forward $(HTTP_SERVICE) $(HTTP_LOCAL_PORT):$(HTTP_SERVICE_PORT)

port-forward/grpc: ## Port-forward the gRPC service to localhost:9095 (Ctrl+C to stop)
	@echo "Port-forwarding $(GRPC_SERVICE) localhost:$(GRPC_LOCAL_PORT) -> $(GRPC_SERVICE_PORT). Press Ctrl+C to stop."
	@kubectl port-forward $(GRPC_SERVICE) $(GRPC_LOCAL_PORT):$(GRPC_SERVICE_PORT)

db/migrate: ## Run database migrations inside the HTTP deployment
	@echo "Running database migrations..."
	@kubectl exec -i $(HTTP_DEPLOYMENT) -- go run ./cmd/main.go run-db-migrations

.DEFAULT_GOAL := help
