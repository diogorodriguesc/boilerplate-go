#!/usr/bin/env bash
# Benchmarks the gRPC UsersService.CreateUser entrypoint using ghz
# (https://ghz.sh/), generating a unique username/email per request via ghz's
# built-in call-data templating (`{{.RequestNumber}}`).
#
# Usage:
#   ./scripts/benchmark-grpc-create-user.sh
#   TOTAL=5000 CONCURRENCY=100 ./scripts/benchmark-grpc-create-user.sh
#   DURATION=30s CONCURRENCY=50 ./scripts/benchmark-grpc-create-user.sh
#   ADDR=localhost:9095 ./scripts/benchmark-grpc-create-user.sh   # e.g. against `make -C k8s port-forward/grpc`
#
# Env vars (all optional):
#   GHZ_BIN        path to the ghz binary (default: .bin/ghz, falls back to `ghz` on PATH)
#   ADDR           target host:port (default: localhost:9090)
#   TOTAL          number of requests to send (default: 1000, ignored if DURATION is set)
#   DURATION       run for a fixed duration instead of a fixed count, e.g. "30s" (default: unset)
#   CONCURRENCY    number of concurrent workers (default: 50)
#   CONNECTIONS    number of gRPC connections to spread the workers over (default: 1)
#   OUTPUT_FORMAT  ghz -O value: summary, csv, json, pretty, html, ... (default: summary)

set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."

GHZ_BIN="${GHZ_BIN:-.bin/ghz}"
if [ ! -x "$GHZ_BIN" ]; then
  if command -v ghz >/dev/null 2>&1; then
    GHZ_BIN="ghz"
  else
    echo "ghz not found at '$GHZ_BIN' or on PATH. Run 'make install/ghz' first." >&2
    exit 1
  fi
fi

ADDR=$(minikube ip):30090
TOTAL="${TOTAL:-1000}"
CONCURRENCY="${CONCURRENCY:-50}"
CONNECTIONS="${CONNECTIONS:-1}"
OUTPUT_FORMAT="${OUTPUT_FORMAT:-summary}"
DURATION="${DURATION:-}"

# Unique per run so re-running the script doesn't collide with usernames/emails
# created by a previous run (CreateUser rejects duplicates with AlreadyExists).
RUN_ID="$(date +%s)-$$"

DATA=$(cat <<EOF
{"username":"benchuser-${RUN_ID}-{{.RequestNumber}}","email":"benchuser-${RUN_ID}-{{.RequestNumber}}@example.com"}
EOF
)

ARGS=(
  --insecure
  --proto proto/users/v1/users.proto
  --import-paths proto
  --call users.v1.UsersService/CreateUser
  -d "$DATA"
  -c "$CONCURRENCY"
  --connections "$CONNECTIONS"
  -O "$OUTPUT_FORMAT"
)

if [ -n "$DURATION" ]; then
  ARGS+=(-z "$DURATION")
  echo "Benchmarking CreateUser at $ADDR (duration=$DURATION, concurrency=$CONCURRENCY, connections=$CONNECTIONS)..." >&2
else
  ARGS+=(-n "$TOTAL")
  echo "Benchmarking CreateUser at $ADDR (total=$TOTAL, concurrency=$CONCURRENCY, connections=$CONNECTIONS)..." >&2
fi

exec "$GHZ_BIN" "${ARGS[@]}" "$ADDR"
