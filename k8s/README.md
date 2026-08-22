# Development Environment Setup

## Pre-Setup Requirements

- Docker
- Minikube
- K9s

## Setup

From the project root directory, execute the following commands

Define MINIKUBE_PROJECTS_PATH:
```bash
export MINIKUBE_PROJECTS_PATH=/home/{...}
```

Start minikube docker container:

First time: 
```bash
minikube start --cpus 4 --memory 8192 --driver=kvm2 \
  --mount \
  --mount-string="$MINIKUBE_PROJECTS_PATH:/app"
```

Onwards: 
```bash
minikube start
```

Stop Minikube:
```bash
minikube stop
```

## Creating Docker image

Create docker image:
```bash
docker build -t micro-app-boilerplate-go:dev -f Dockerfile.dev .
```

Load created docker image into minikube
```bash
minikube image load micro-app-boilerplate-go:dev
```

## Applying infrastructure

Apply Resources:
```bash
kubectl apply -k k8s/base/
```

Any change you make to kubernetes resources must be applied.

Delete Resources:
```bash
kubectl delete -k k8s/base/
```

## Application DB Migrations

Run database migrations:
```bash
kubectl exec -i deployment/micro-app-boilerplate-go -- go run ./cmd/main.go run-db-migrations
```
