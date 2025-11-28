#!/bin/sh
set -e

# Load .env if present (mounted into /app/.env)
if [ -f /app/.env ]; then
  set -a
  . /app/.env
  set +a
fi

DB_PATH="${DB_PATH:-/data/data.db}"
export DB_PATH

if [ "${RUN_MIGRATIONS:-1}" = "1" ]; then
  echo "Running migrations..."
  /app/migrate
fi

if [ "${RUN_SEED:-0}" = "1" ]; then
  echo "Seeding data..."
  /app/seed || true
fi

exec "$@"
