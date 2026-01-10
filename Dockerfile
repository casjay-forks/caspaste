# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.21-alpine AS builder

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
    ./cmd/caspaste

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 caspaste && \
    adduser -D -u 1000 -G caspaste caspaste

# Create data directory for SQLite
RUN mkdir -p /data && chown caspaste:caspaste /data

# Copy binary from builder
COPY --from=builder /caspaste /usr/local/bin/caspaste

# Switch to non-root user
USER caspaste

# Set working directory
WORKDIR /data

# Expose default port
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:80/ || exit 1

# Entrypoint with default arguments
# Default: SQLite database in /data, listening on :80
ENTRYPOINT ["caspaste", "-db-source", "/data/lenpaste.db"]

# Additional arguments can be passed when running the container
# Example: docker run caspaste -address :8080 -db-driver postgres -db-source "postgresql://..."
