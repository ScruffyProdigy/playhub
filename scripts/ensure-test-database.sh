#!/bin/bash
# Create playhub_test if missing. Works with docker compose postgres or CI service containers.
set -euo pipefail

PGHOST="${PGHOST:-127.0.0.1}"
PGUSER="${PGUSER:-app}"
PGPASSWORD="${PGPASSWORD:-app-pass}"

export PGPASSWORD
if psql -h "$PGHOST" -U "$PGUSER" -d postgres -tc "SELECT 1 FROM pg_database WHERE datname = 'playhub_test'" | grep -q 1; then
  exit 0
fi
psql -h "$PGHOST" -U "$PGUSER" -d postgres -c "CREATE DATABASE playhub_test OWNER app;"
