# Boilerplate Golang

## DB Migrations

Run migrations:
```bash
kubectl exec -i deployment/micro-app-boilerplate-go -- go run ./cmd/main.go run-db-migrations
```

## Build Application & Run

```bash
make build
./bin/api http-server
```

```bash
curl --location 'localhost:8081/v1/users?email=email@domain.com'
```
