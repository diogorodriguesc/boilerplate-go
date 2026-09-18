#!/usr/bin/env bash
# Benchmarks the HTTP POST /v1/users (CreateUser) entrypoint using vegeta
# (https://github.com/tsenart/vegeta). Counterpart to
# benchmark-grpc-create-user.sh for the gRPC entrypoint.
#
# CreateUser rejects duplicate username/email with 409 Conflict, so every
# request needs a unique body. vegeta has no built-in per-request templating
# (unlike ghz), so this script pre-generates one small body file per request
# and points vegeta at them via its "http" target format (`@bodyfile`) —
# that's a plain bash loop with no subprocess forks per line, so pre-
# generating even tens of thousands of bodies takes a fraction of a second
# and doesn't eat into the timed benchmark window.
#
# Usage:
#   ./scripts/benchmark-http-create-user.sh
#   TOTAL=5000 CONCURRENCY=100 ./scripts/benchmark-http-create-user.sh
#   DURATION=30s CONCURRENCY=50 ./scripts/benchmark-http-create-user.sh
#   ADDR=localhost:8085 ./scripts/benchmark-http-create-user.sh   # e.g. against `make -C k8s port-forward/http`
#
# Env vars (all optional):
#   VEGETA_BIN     path to the vegeta binary (default: .bin/vegeta, falls back to `vegeta` on PATH)
#   ADDR           target host:port (default: localhost:8080)
#   TOTAL          number of requests to send (default: 1000, ignored if DURATION is set)
#   DURATION       run for a fixed duration instead of a fixed count, e.g. "30s" (default: unset)
#   CONCURRENCY    number of concurrent workers (default: 50)
#   POOL_SIZE      number of unique request bodies to pre-generate when DURATION is
#                  set (default: 20000). Must be large enough to cover the whole
#                  duration at whatever throughput is achieved — if it runs out
#                  first, the benchmark stops early (a warning is printed); bump
#                  this if that happens.
#
# Note: with a fixed TOTAL, vegeta stops the instant its target list is
# exhausted, which the report shows as one extra "no targets to attack" (code
# 0) entry — a harmless, well-known vegeta artifact, negligible at any
# realistic TOTAL.

set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."

VEGETA_BIN="${VEGETA_BIN:-.bin/vegeta}"
if [ ! -x "$VEGETA_BIN" ]; then
  if command -v vegeta >/dev/null 2>&1; then
    VEGETA_BIN="vegeta"
  else
    echo "vegeta not found at '$VEGETA_BIN' or on PATH. Run 'make install/vegeta' first." >&2
    exit 1
  fi
fi

ADDR=$(minikube ip):30080
TOTAL="${TOTAL:-1000}"
CONCURRENCY="${CONCURRENCY:-50}"
DURATION="${DURATION:-}"
POOL_SIZE="${POOL_SIZE:-20000}"

WORKDIR="$(mktemp -d)"
trap 'rm -rf "$WORKDIR"' EXIT

# Unique per run so re-running the script doesn't collide with usernames/emails
# created by a previous run.
RUN_ID="$(date +%s)-$$"

if [ -n "$DURATION" ]; then
  COUNT="$POOL_SIZE"
else
  COUNT="$TOTAL"
fi

TARGETS_FILE="$WORKDIR/targets.http"
: > "$TARGETS_FILE"
for ((i = 1; i <= COUNT; i++)); do
  body_file="$WORKDIR/body-$i.json"
  echo "{\"username\":\"benchuser-${RUN_ID}-${i}\",\"email\":\"benchuser-${RUN_ID}-${i}@example.com\"}" > "$body_file"
  {
    echo "POST http://${ADDR}/v1/users"
    echo "Content-Type: application/json"
    echo "@${body_file}"
    echo
  } >> "$TARGETS_FILE"
done

ARGS=(
  -targets="$TARGETS_FILE"
  -lazy
  -rate=0
  -workers="$CONCURRENCY"
  -max-workers="$CONCURRENCY"
)

RESULTS_FILE="$WORKDIR/results.bin"

if [ -n "$DURATION" ]; then
  ARGS+=(-duration="$DURATION")
  echo "Benchmarking CreateUser at $ADDR (duration=$DURATION, concurrency=$CONCURRENCY, pool=$POOL_SIZE)..." >&2
else
  ARGS+=(-duration=0)
  echo "Benchmarking CreateUser at $ADDR (total=$TOTAL, concurrency=$CONCURRENCY)..." >&2
fi

"$VEGETA_BIN" attack "${ARGS[@]}" -output="$RESULTS_FILE"
"$VEGETA_BIN" report "$RESULTS_FILE"

if [ -n "$DURATION" ] && "$VEGETA_BIN" encode < "$RESULTS_FILE" | grep -q "no targets to attack"; then
  echo >&2
  echo "WARNING: ran out of the $POOL_SIZE pre-generated request bodies before the $DURATION duration elapsed." >&2
  echo "Re-run with a larger POOL_SIZE to sustain the full duration, e.g. POOL_SIZE=100000." >&2
fi
