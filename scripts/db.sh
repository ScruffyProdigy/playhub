#!/bin/bash
# Manage the shared PlayHub PostgreSQL instance (docker compose).
set -e

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

export DATABASE_URL="${DATABASE_URL:-postgres://app:app-pass@127.0.0.1:5432/playhub?sslmode=disable}"
# macOS often resolves localhost to ::1 while Docker publishes Postgres on IPv4.
if [[ "$DATABASE_URL" == *"@localhost:"* ]]; then
  export DATABASE_URL="${DATABASE_URL/@localhost:/@127.0.0.1:}"
fi

require_docker() {
  if ! command -v docker &> /dev/null; then
    echo "Docker is required. Install Docker and try again."
    exit 1
  fi

  if ! docker info >/dev/null 2>&1; then
    echo "Cannot connect to the Docker daemon."
    if command -v colima &> /dev/null; then
      echo ""
      echo "Colima is installed but Docker is not reachable. Try:"
      echo "  colima start"
      echo ""
      echo "If that does not help, restart Colima:"
      echo "  colima stop && colima start"
    else
      echo "Start Docker Desktop or your Docker runtime, then try again."
    fi
    exit 1
  fi
}

wait_for_postgres() {
  echo "Waiting for PostgreSQL..."
  for i in {1..30}; do
    if docker compose exec -T postgres pg_isready -U app -d playhub >/dev/null 2>&1; then
      break
    fi
    if [ "$i" -eq 30 ]; then
      echo "PostgreSQL failed to become ready"
      exit 1
    fi
    sleep 1
  done

  echo "Waiting for PostgreSQL port on host..."
  for i in {1..30}; do
    if nc -z 127.0.0.1 5432 >/dev/null 2>&1; then
      echo "PostgreSQL is ready"
      return 0
    fi
    sleep 1
  done

  echo "PostgreSQL is running in Docker but port 5432 is not reachable on 127.0.0.1"
  exit 1
}

case "${1:-}" in
  up)
    require_docker
    docker compose up -d postgres redis
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
    (cd backend && DATABASE_URL="$DATABASE_URL" make migrate-up)
    ;;
  reset)
    require_docker
    docker compose down -v
    docker compose up -d postgres redis
    wait_for_postgres
    (cd backend && DATABASE_URL="$DATABASE_URL" make migrate-up)
    ;;
  clean-test-data)
    require_docker
    wait_for_postgres
    docker compose exec -T postgres psql -U app -d playhub <<'SQL'
BEGIN;
DELETE FROM magic_links WHERE email LIKE '%@example.com';
DELETE FROM users WHERE email LIKE '%@example.com';
DELETE FROM games WHERE category IS DISTINCT FROM 'demo';
COMMIT;
SQL
    echo "Removed test games, @example.com users, and their magic links."
    echo "Demo games and real user accounts were kept."
    ;;
  url)
    echo "$DATABASE_URL"
    ;;
  *)
    echo "Usage: $0 {up|down|wait|migrate|reset|clean-test-data|url}"
    exit 1
    ;;
esac
