# All-in-one free-tier image: PostGIS + auth service + main API + nginx gateway.
# Ephemeral storage is fine for demos — migrations + seed run on every cold start.
FROM golang:1.23-alpine AS build
WORKDIR /src
RUN apk add --no-cache git
ENV GOWORK=off

COPY pkg/authkit/ pkg/authkit/
COPY pkg/platform/ pkg/platform/
COPY backend/go.mod backend/go.sum backend/
COPY services/auth/go.mod services/auth/go.sum services/auth/

WORKDIR /src/backend
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/api ./cmd/server

WORKDIR /src/services/auth
RUN go mod download
COPY services/auth/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/auth ./cmd/server

FROM postgis/postgis:16-3.4
USER root
RUN apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates curl nginx \
  && rm -rf /var/lib/apt/lists/*

COPY --from=build /out/api /app/api
COPY --from=build /out/auth /app/auth
COPY deploy/nginx.conf /etc/nginx/nginx.conf
COPY deploy/start.sh /app/start.sh
RUN chmod +x /app/start.sh /app/api /app/auth \
  && mkdir -p /run/nginx

ENV POSTGRES_USER=didar \
    POSTGRES_PASSWORD=didar \
    POSTGRES_DB=didar \
    APP_ENV=production \
    HTTP_PORT=8080 \
    DATABASE_URL=postgres://didar:didar@127.0.0.1:5432/didar?sslmode=disable \
    AUTH_DATABASE_URL=postgres://didar:didar@127.0.0.1:5432/didar_auth?sslmode=disable \
    AUTH_SERVICE_URL=http://127.0.0.1:8081 \
    INTERNAL_SERVICE_KEY=didar-demo-internal-key \
    OTP_CHANNEL=email \
    OTP_MOCK_MODE=true \
    DEFAULT_BRAND_SLUG=didar \
    JWT_SECRET=didar-demo-change-me-in-production

EXPOSE 8080
CMD ["/app/start.sh"]
