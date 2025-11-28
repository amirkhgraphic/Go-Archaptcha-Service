#!/bin/sh
set -e

if [ -f /app/.env ]; then
  echo "Loading .env file..."
  set -a
  . /app/.env
  set +a
else
  echo "SKIP loading .env file..."
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
