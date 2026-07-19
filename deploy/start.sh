#!/usr/bin/env bash
# All-in-one container: PostGIS + auth + main API + nginx (Render free tier).
set -euo pipefail

export POSTGRES_USER="${POSTGRES_USER:-didar}"
export POSTGRES_PASSWORD="${POSTGRES_PASSWORD:-didar}"
export POSTGRES_DB="${POSTGRES_DB:-didar}"
export DATABASE_URL="${DATABASE_URL:-postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@127.0.0.1:5432/${POSTGRES_DB}?sslmode=disable}"
export AUTH_DATABASE_URL="${AUTH_DATABASE_URL:-postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@127.0.0.1:5432/didar_auth?sslmode=disable}"
export AUTH_SERVICE_URL="${AUTH_SERVICE_URL:-http://127.0.0.1:8081}"
export INTERNAL_SERVICE_KEY="${INTERNAL_SERVICE_KEY:-didar-demo-internal-key}"
export JWT_SECRET="${JWT_SECRET:-didar-demo-change-me-in-production}"
export OTP_CHANNEL="${OTP_CHANNEL:-email}"
export OTP_MOCK_MODE="${OTP_MOCK_MODE:-true}"
export DEFAULT_BRAND_SLUG="${DEFAULT_BRAND_SLUG:-didar}"
export CORS_ORIGINS="${CORS_ORIGINS:-https://alan14alashti.github.io,http://localhost:5173}"

# Public port for nginx (Render injects PORT).
PUBLIC_PORT="${PORT:-${HTTP_PORT:-8080}}"

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

export PGPASSWORD="${POSTGRES_PASSWORD}"

# Auth database (idempotent).
psql -h 127.0.0.1 -U "$POSTGRES_USER" -d "$POSTGRES_DB" -tc "SELECT 1 FROM pg_database WHERE datname = 'didar_auth'" | grep -q 1 \
  || psql -h 127.0.0.1 -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "CREATE DATABASE didar_auth;"

# Auth service on 8081
HTTP_PORT=8081 DATABASE_URL="$AUTH_DATABASE_URL" /app/auth &
AUTH_PID=$!

# Main API on 8082
HTTP_PORT=8082 DATABASE_URL="$DATABASE_URL" AUTH_SERVICE_URL="$AUTH_SERVICE_URL" /app/api &
API_PID=$!

# Point nginx at the public port Render expects.
sed -i "s/listen 8080;/listen ${PUBLIC_PORT};/" /etc/nginx/nginx.conf

echo "Starting nginx on :${PUBLIC_PORT} (auth=:8081 api=:8082)"
nginx -g 'daemon off;' &
NGINX_PID=$!

cleanup() {
  kill "$NGINX_PID" "$AUTH_PID" "$API_PID" 2>/dev/null || true
}
trap cleanup EXIT INT TERM

# Exit if any child dies.
while true; do
  if ! kill -0 "$AUTH_PID" 2>/dev/null; then
    echo "auth process exited" >&2
    exit 1
  fi
  if ! kill -0 "$API_PID" 2>/dev/null; then
    echo "api process exited" >&2
    exit 1
  fi
  if ! kill -0 "$NGINX_PID" 2>/dev/null; then
    echo "nginx process exited" >&2
    exit 1
  fi
  sleep 2
done
