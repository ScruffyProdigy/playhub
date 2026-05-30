#!/bin/bash
# Manage the shared PlayHub PostgreSQL instance (docker compose).
set -e

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

export DATABASE_URL="${DATABASE_URL:-postgres://app:app-pass@localhost:5432/playhub?sslmode=disable}"

require_docker() {
  if ! command -v docker &> /dev/null; then
    echo "Docker is required. Install Docker and try again."
    exit 1
  fi
}

wait_for_postgres() {
  echo "Waiting for PostgreSQL..."
  for i in {1..30}; do
    if docker compose exec -T postgres pg_isready -U app -d playhub >/dev/null 2>&1; then
      echo "PostgreSQL is ready"
      return 0
    fi
    sleep 1
  done
  echo "PostgreSQL failed to become ready"
  exit 1
}

case "${1:-}" in
  up)
    require_docker
    docker compose up -d postgres
    wait_for_postgres
    ;;
  down)
    require_docker
    docker compose down
    ;;
  wait)
    require_docker
    wait_for_postgres
    ;;
  migrate)
    require_docker
    wait_for_postgres
    (cd backend && make migrate-up)
    ;;
  reset)
    require_docker
    docker compose down -v
    docker compose up -d postgres
    wait_for_postgres
    (cd backend && make migrate-up)
    ;;
  url)
    echo "$DATABASE_URL"
    ;;
  *)
    echo "Usage: $0 {up|down|wait|migrate|reset|url}"
    exit 1
    ;;
esac
