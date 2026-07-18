# All-in-one free-tier image: PostGIS + Didar API.
# Ephemeral storage is fine for demos — migrations + seed run on every cold start.
FROM golang:1.23-alpine AS build
WORKDIR /src
RUN apk add --no-cache git
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server

FROM postgis/postgis:16-3.4
USER root
RUN apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates curl \
  && rm -rf /var/lib/apt/lists/*

COPY --from=build /out/server /app/server
COPY deploy/start.sh /app/start.sh
RUN chmod +x /app/start.sh /app/server

ENV POSTGRES_USER=didar \
    POSTGRES_PASSWORD=didar \
    POSTGRES_DB=didar \
    APP_ENV=production \
    HTTP_PORT=8080 \
    DATABASE_URL=postgres://didar:didar@127.0.0.1:5432/didar?sslmode=disable \
    OTP_MOCK_MODE=true \
    DEFAULT_BRAND_SLUG=didar \
    JWT_SECRET=didar-demo-change-me-in-production

EXPOSE 8080
CMD ["/app/start.sh"]
