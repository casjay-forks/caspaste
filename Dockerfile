# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
ARG VERSION=unknown
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -tags netgo \
    -ldflags "-w -s -X main.Version=${VERSION} -extldflags '-static'" \
    -o /caspaste \
    ./src/cmd/caspaste

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user and directories
RUN addgroup -g 1000 caspaste && \
    adduser -D -u 1000 -G caspaste caspaste && \
    mkdir -p /data/caspaste /data/db/sqlite /config/caspaste && \
    chown -R caspaste:caspaste /data /config

# Copy binary from builder
COPY --from=builder /caspaste /usr/local/bin/caspaste

# Switch to non-root user
USER caspaste

# Set working directory
WORKDIR /data/caspaste

# Expose default port
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:80/ || exit 1

# Entrypoint with default arguments
# Container Paths:
#   - Config: /config/caspaste/caspaste.yml (auto-generated on first run)
#   - Data: /data/caspaste/ (backups and other data)
#   - Database: /data/db/sqlite/caspaste.db (SQLite database)
# Default Ports:
#   - Internal: 80
#   - External: Map to 172.17.0.1:64365
# Usage: docker run -p 172.17.0.1:64365:80 -v ...
ENTRYPOINT ["caspaste", "--config", "/config/caspaste", "--data", "/data/caspaste", "--db-source", "/data/db/sqlite/caspaste.db"]

# Additional arguments can be passed when running the container
# Example: docker run -p 172.17.0.1:64365:80 caspaste --ui-default-theme nord
