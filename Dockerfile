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
ENV CASPASTE_DB_DIR=/data/db/sqlite

# Expose default port
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:80/ || exit 1

# Entrypoint with default arguments
# Container Paths:
#   - Config: /config/caspaste.yml (auto-generated on first run)
#   - Data: /data/ (application data)
#   - Database: /data/db/sqlite/caspaste.db (SQLite)
#   - Backups: /data/backups/
#   - Cache: Auto-detected
#   - Logs: Auto-detected
# Host Mounts (from docker-compose.yml):
#   - ./rootfs/data:/data
#   - ./rootfs/config:/config (optional)
# Privilege Escalation:
#   - Creates user caspaste (UID:GID 642:642)
#   - Binds to port as root
#   - Drops privileges to caspaste user
# Security:
#   - Auto-trusts reverse proxy headers from private IPs
#   - Prevents IP spoofing from public IPs
ENTRYPOINT ["caspaste", "--config", "/config", "--data", "/data"]
