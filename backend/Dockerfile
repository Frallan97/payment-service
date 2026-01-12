# Build stage
FROM golang:1.24-alpine AS build

RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -trimpath \
    -o server ./cmd/api

# Final stage
FROM alpine:3.18

RUN apk add --no-cache ca-certificates curl

# Install golang-migrate for migrations
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz \
    | tar xz -C /usr/local/bin migrate

WORKDIR /app

# Copy binary and migrations
COPY --from=build /app/server /usr/local/bin/server
COPY migrations /app/migrations

# Expose port
EXPOSE 8080

# Set entrypoint
ENTRYPOINT ["/usr/local/bin/server"]
