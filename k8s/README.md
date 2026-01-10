# Development Environment Setup

## Pre-Setup Requirements

- Docker
- Minikube
- K9s

## Setup

From project root directory, execute the following commands

Define MINIKUBE_PROJECTS_PATH:
```bash
export MINIKUBE_PROJECTS_PATH=/home/{...}
```

Start minikube docker container:

First time: 
```bash
minikube start --driver=kvm2 \
  --mount \
  --mount-string="$MINIKUBE_PROJECTS_PATH:/app"
```

Onwards: 
```bash
minikube start
```

Create docker image:
```bash
docker build -t micro-app-boilerplate-go:dev -f Dockerfile.dev .
```

Load created docker image into minikube
```bash
minikube image load micro-app-boilerplate-go:dev
```

Apply Resources:
```bash
kubectl apply -k k8s/base/
```

Any change you make to kubernetes resources must be applied.

Delete Resources:
```bash
kubectl delete -k k8s/base/
```

Stop Minikube:
```bash
minikube stop
```