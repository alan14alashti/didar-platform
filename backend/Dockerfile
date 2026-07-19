FROM golang:1.23-alpine AS build
WORKDIR /src
RUN apk add --no-cache git
COPY go.work go.work
COPY pkg/authkit/ pkg/authkit/
COPY pkg/platform/ pkg/platform/
COPY backend/go.mod backend/go.sum ./backend/
WORKDIR /src/backend
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=build /out/server ./server
EXPOSE 8080
ENTRYPOINT ["./server"]
