#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$ROOT_DIR"

PORT="${PORT:-19000}"
export PORT

kill_port_process() {
  local port="$1"
  local pid=""

  if command -v lsof >/dev/null 2>&1; then
    pid="$(lsof -ti "tcp:${port}" 2>/dev/null || true)"
  elif command -v ss >/dev/null 2>&1; then
    pid="$(ss -ltnp "sport = :${port}" 2>/dev/null | awk -F'pid=' 'NR>1 {print $2}' | awk -F',' '{print $1}' | head -n 1)"
  fi

  if [ -n "${pid}" ]; then
    echo "[start] port ${port} is in use by pid ${pid}, stopping it..."
    kill "${pid}" 2>/dev/null || true
    sleep 1
    if kill -0 "${pid}" 2>/dev/null; then
      echo "[start] pid ${pid} still running, forcing stop..."
      kill -9 "${pid}" 2>/dev/null || true
    fi
  fi
}

echo "[start] building frontend..."
make build-frontend

echo "[start] building backend..."
go build ./...

kill_port_process "${PORT}"

echo "[start] starting backend on port ${PORT}..."
exec go run main.go --port "${PORT}"
