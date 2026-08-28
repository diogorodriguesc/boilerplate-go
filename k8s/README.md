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
make minikube/start-first-time
```

Onwards: 
```bash
make minikube/start
```

Stop Minikube:
```bash
make minikube/stop
```

## Creating Docker image

Create docker image:
```bash
make docker/build
```

Load created docker image into minikube
```bash
make docker/load
```

## Applying infrastructure

Apply Resources:
```bash
make apply
```

Any change you make to kubernetes resources must be applied.

Delete Resources:
```bash
make delete
```

## Application DB Migrations

Run database migrations:
```bash
make db/migrate
```

To get all commands available:
```bash
make help
```