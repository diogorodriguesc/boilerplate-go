.PHONY: help minikube/start minikube/start-first-time minikube/stop docker/build docker/load apply delete db/migrate port-forward

IMAGE_NAME ?= micro-app-boilerplate-go
IMAGE_TAG ?= dev
IMAGE ?= $(IMAGE_NAME):$(IMAGE_TAG)

DOCKERFILE ?= $(CURDIR)/../Dockerfile.dev
BUILD_CONTEXT ?= $(CURDIR)/..

K8S_BASE ?= $(CURDIR)/base
DEPLOYMENT ?= deployment/micro-app-boilerplate-go
SERVICE ?= svc/go-k8s-service
LOCAL_PORT ?= 8085
SERVICE_PORT ?= 8080

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

apply: ## Apply Kubernetes resources, wait for rollout, then port-forward the service to localhost:8085
	@echo "Applying Kubernetes resources..."
	@kubectl apply -k $(K8S_BASE)
	@kubectl rollout status $(DEPLOYMENT)
	@$(MAKE) port-forward

delete: ## Delete Kubernetes resources
	@echo "Deleting Kubernetes resources..."
	@kubectl delete -k $(K8S_BASE)

port-forward: ## Port-forward go-k8s-service to localhost:8085 (Ctrl+C to stop)
	@echo "Port-forwarding $(SERVICE) localhost:$(LOCAL_PORT) -> $(SERVICE_PORT). Press Ctrl+C to stop."
	@kubectl port-forward $(SERVICE) $(LOCAL_PORT):$(SERVICE_PORT)

db/migrate: ## Run database migrations inside the application deployment
	@echo "Running database migrations..."
	@kubectl exec -i $(DEPLOYMENT) -- go run ./cmd/main.go run-db-migrations

.DEFAULT_GOAL := help
