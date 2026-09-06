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

The app runs as two separate Deployments/Services, sharing the same image but
each running a single entrypoint command:

- `micro-app-boilerplate-go-http` / `go-http-k8s-service` — runs `http-server`,
  exposed on port 8080 (nodePort 30080).
- `micro-app-boilerplate-go-grpc` / `go-grpc-k8s-service` — runs `grpc-server`,
  exposed on port 9090 (nodePort 30090).

Port-forward either one to your machine:
```bash
kubectl port-forward svc/go-http-k8s-service 8085:8080
kubectl port-forward svc/go-grpc-k8s-service 9095:9090
```
(or `make -C k8s port-forward/http` / `make -C k8s port-forward/grpc`)

Since both are `NodePort` services, minikube already exposes them directly on
the node's IP — no `kubectl port-forward` (and its proxy process/terminal)
needed:
```bash
minikube ip                                    # e.g. 192.168.39.193
curl http://$(minikube ip):30080/health
grpcurl -plaintext $(minikube ip):30090 list
```
or let minikube resolve the URL for you: `minikube service go-http-k8s-service --url`.
This is also the better target for the benchmark scripts
(`ADDR=$(minikube ip):30080` / `:30090`) — `kubectl port-forward` proxies
every request through an extra local process, which adds latency and a
throughput ceiling that direct NodePort access doesn't have. Measured on the
same 300-request HTTP run: port-forward gave p50 35.8ms / p99 239.6ms at
~520 req/s, direct NodePort gave p50 17.3ms / p99 29ms at ~1390 req/s — so
port-forward numbers should not be trusted as the server's real capacity.

## Application DB Migrations

Run database migrations (either deployment works — they share the same
database, this just uses the HTTP one):
```bash
kubectl exec -i deployment/micro-app-boilerplate-go-http -- go run ./cmd/main.go run-db-migrations
```
