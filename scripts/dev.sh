#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

API_HOST="${API_HOST:-127.0.0.1}"
API_PORT="${API_PORT:-8080}"
GOCACHE="${GOCACHE:-$ROOT_DIR/.gocache}"
export GOCACHE
mkdir -p "$GOCACHE"

API_PID=""
UI_PID=""

cleanup() {
  local exit_code=$?

  if [[ -n "${API_PID}" ]] && kill -0 "${API_PID}" 2>/dev/null; then
    kill -TERM "${API_PID}" 2>/dev/null || true
  fi

  if [[ -n "${UI_PID}" ]] && kill -0 "${UI_PID}" 2>/dev/null; then
    kill -TERM "${UI_PID}" 2>/dev/null || true
  fi

  if [[ -n "${API_PID}" ]]; then
    wait "${API_PID}" 2>/dev/null || true
  fi

  if [[ -n "${UI_PID}" ]]; then
    wait "${UI_PID}" 2>/dev/null || true
  fi

  exit "${exit_code}"
}

trap cleanup INT TERM EXIT

echo "Starting API on http://${API_HOST}:${API_PORT}"
go run . api --host "${API_HOST}" --port "${API_PORT}" &
API_PID=$!

echo "Starting UI on http://127.0.0.1:5173"
(
  cd ui
  npm run dev
) &
UI_PID=$!

# Keep this parent process alive while both children are running.
# If either process exits unexpectedly, stop the other one.
while true; do
  if ! kill -0 "${API_PID}" 2>/dev/null; then
    echo "API process exited; shutting down UI"
    break
  fi

  if ! kill -0 "${UI_PID}" 2>/dev/null; then
    echo "UI process exited; shutting down API"
    break
  fi

  sleep 1
done
