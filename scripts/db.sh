#!/bin/bash
# Manage the shared PlayHub PostgreSQL instance (docker compose).
set -e

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

export DATABASE_URL="${DATABASE_URL:-postgres://app:app-pass@127.0.0.1:5432/playhub?sslmode=disable}"
export TEST_DATABASE_URL="${TEST_DATABASE_URL:-postgres://app:app-pass@127.0.0.1:5432/playhub_test?sslmode=disable}"
# macOS often resolves localhost to ::1 while Docker publishes Postgres on IPv4.
if [[ "$DATABASE_URL" == *"@localhost:"* ]]; then
  export DATABASE_URL="${DATABASE_URL/@localhost:/@127.0.0.1:}"
fi
if [[ "$TEST_DATABASE_URL" == *"@localhost:"* ]]; then
  export TEST_DATABASE_URL="${TEST_DATABASE_URL/@localhost:/@127.0.0.1:}"
fi

ensure_test_database() {
  docker compose exec -T postgres psql -U app -d postgres -v ON_ERROR_STOP=1 <<'SQL'
SELECT 'CREATE DATABASE playhub_test OWNER app'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'playhub_test')\gexec
SQL
}

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
  test-url)
    echo "$TEST_DATABASE_URL"
    ;;
  test-migrate)
    require_docker
    wait_for_postgres
    ensure_test_database
    (cd backend && DATABASE_URL="$TEST_DATABASE_URL" make migrate-up)
    echo "Migrations applied to playhub_test (dev database playhub unchanged)."
    ;;
  test-reset)
    require_docker
    wait_for_postgres
    docker compose exec -T postgres psql -U app -d postgres -v ON_ERROR_STOP=1 <<'SQL'
DROP DATABASE IF EXISTS playhub_test WITH (FORCE);
CREATE DATABASE playhub_test OWNER app;
SQL
    (cd backend && DATABASE_URL="$TEST_DATABASE_URL" make migrate-up)
    echo "Reset and migrated playhub_test."
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
  reset-demo-handoff)
    # Dev database only (playhub). Integration tests use playhub_test and restore URLs in cleanup.
    require_docker
    wait_for_postgres
    docker compose exec -T postgres psql -U app -d playhub <<'SQL'
UPDATE games
SET play_url = 'http://localhost:5174',
    api_base_url = 'http://localhost:3001'
WHERE id = 'a1000000-0000-4000-8000-000000000001';
SQL
    echo "Restored demo quick-match handoff URLs (play :5174, API :3001)."
    ;;
  url)
    echo "$DATABASE_URL"
    ;;
  *)
    echo "Usage: $0 {up|down|wait|migrate|reset|test-url|test-migrate|test-reset|clean-test-data|reset-demo-handoff|url}"
    echo ""
    echo "  url / migrate     — development database (playhub)"
    echo "  test-url          — integration test database URL (playhub_test)"
    echo "  test-migrate      — migrate playhub_test only"
    echo "  test-reset        — drop and recreate playhub_test, then migrate"
    exit 1
    ;;
esac
