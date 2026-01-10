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

# Create non-root user and directories
RUN addgroup -g 1000 caspaste && \
    adduser -D -u 1000 -G caspaste caspaste && \
    mkdir -p /data /data/db /config && \
    chown -R caspaste:caspaste /data /config

# Copy binary from builder
COPY --from=builder /caspaste /usr/local/bin/caspaste

# Switch to non-root user
USER caspaste

# Set working directory
WORKDIR /data

# Set environment variables for Docker deployment
# Database directory inside container (defaults to /data/db)
# Backup directory on host OS (must be set via docker-compose or docker run)
ENV CASPASTE_DB_DIR=/data/db

# Expose default port
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:80/ || exit 1

# Entrypoint with default arguments
# Container Paths:
#   - Config: /config/caspaste.yml (auto-generated on first run)
#   - Data: /data/ (application data)
#   - Database: /data/db/caspaste.db (SQLite via CASPASTE_DB_DIR)
#   - Backups: /data/backups/ (via CASPASTE_BACKUP_DIR)
# Host Mounts:
#   - ./rootfs/config:/config
#   - ./rootfs/data:/data
#   - ./rootfs/db:/data/db
#   - ./rootfs/data/backups:/data/backups (→ /mnt/Backups/caspaste on Linux host)
# Privilege Escalation:
#   - Binary prefers system directories when running as root (/var/log, /var/lib, etc.)
#   - Falls back to user directories when running as non-root
# Security:
#   - Auto-trusts reverse proxy headers from private IPs (10.x, 172.16-31.x, 192.168.x, fc00::/7)
#   - Prevents IP spoofing from public IPs
# Usage: docker run -p 172.17.0.1:64365:80 -v ./rootfs/data:/data ...
ENTRYPOINT ["caspaste", "--config", "/config", "--data", "/data"]

# Additional arguments can be passed when running the container
# Example: docker run -p 172.17.0.1:64365:80 caspaste --ui-default-theme nord
