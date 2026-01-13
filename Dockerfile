# syntax=docker/dockerfile:1

# Build stage
FROM golang:alpine AS builder

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

# Copy binary from builder
COPY --from=builder /caspaste /usr/local/bin/caspaste

# Set working directory
WORKDIR /data

# Set environment variables for Docker deployment
ENV CASPASTE_DB_DIR=/data/db/sqlite \
    PORT=80

# Expose HTTP port
EXPOSE 80

# Health check - uses builtin --status flag
# Exit codes: 0=healthy, 1=unhealthy, 2=degraded
HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD caspaste --config /config/caspaste --data /data/caspaste --status

# Entrypoint with default arguments
# Container Paths:
#   - Config: /config/caspaste/caspaste.yml (auto-generated)
#   - Data: /data/caspaste/ (application data)
#   - Database: /data/db/sqlite/caspaste.db (SQLite)
#   - Backups: /data/backups/
#   - Cache: /cache/
#   - Logs: /logs/
# Host Mounts (from docker-compose.yml):
#   - ./rootfs/config/caspaste:/config/caspaste
#   - ./rootfs/data/caspaste:/data/caspaste
#   - ./rootfs/data/db/sqlite:/data/db/sqlite
# Privilege Escalation:
#   - Creates user caspaste (UID:GID 642:642)
#   - Binds to port as root
#   - Drops privileges to caspaste user
# Security:
#   - Auto-trusts reverse proxy headers from private IPs
#   - Prevents IP spoofing from public IPs
ENTRYPOINT ["caspaste", "--config", "/config/caspaste", "--data", "/data/caspaste"]
