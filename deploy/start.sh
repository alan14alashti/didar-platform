#!/usr/bin/env bash
set -euo pipefail

# Render injects PORT; our Go server reads HTTP_PORT.
export HTTP_PORT="${PORT:-${HTTP_PORT:-8080}}"
export POSTGRES_USER="${POSTGRES_USER:-didar}"
export POSTGRES_PASSWORD="${POSTGRES_PASSWORD:-didar}"
export POSTGRES_DB="${POSTGRES_DB:-didar}"
export DATABASE_URL="${DATABASE_URL:-postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@127.0.0.1:5432/${POSTGRES_DB}?sslmode=disable}"

# Start PostGIS via the official image entrypoint (background).
docker-entrypoint.sh postgres &

echo "Waiting for Postgres..."
for i in $(seq 1 90); do
  if pg_isready -h 127.0.0.1 -U "$POSTGRES_USER" -d "$POSTGRES_DB" >/dev/null 2>&1; then
    echo "Postgres is ready."
    break
  fi
  sleep 1
done

if ! pg_isready -h 127.0.0.1 -U "$POSTGRES_USER" -d "$POSTGRES_DB" >/dev/null 2>&1; then
  echo "Postgres failed to become ready" >&2
  exit 1
fi

exec /app/server
