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

# Build arguments
ARG VERSION=unknown
ARG BUILD_DATE
ARG VCS_REF
ARG VCS_URL

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -tags netgo \
    -ldflags "-w -s -X main.Version=${VERSION} -extldflags '-static'" \
    -o /caspaste \
    ./src/cmd/caspaste

# Final stage
FROM alpine:latest

# OCI Standard Labels
# https://github.com/opencontainers/image-spec/blob/main/annotations.md
LABEL org.opencontainers.image.title="CasPaste"
LABEL org.opencontainers.image.description="A simple, fast, and secure self-hosted pastebin service with modern UI and 12 themes"
LABEL org.opencontainers.image.authors="CasjaysDev <docker-admin@casjaysdev.pro>"
LABEL org.opencontainers.image.vendor="CasjaysDev"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.url="https://github.com/casjay-forks/caspaste"
LABEL org.opencontainers.image.documentation="https://github.com/casjay-forks/caspaste/blob/main/README.md"
LABEL org.opencontainers.image.source="https://github.com/casjay-forks/caspaste"

# Dynamic labels (set via build args)
ARG VERSION=unknown
ARG BUILD_DATE
ARG VCS_REF
ARG VCS_URL="https://github.com/casjay-forks/caspaste"

LABEL org.opencontainers.image.version="${VERSION}"
LABEL org.opencontainers.image.created="${BUILD_DATE}"
LABEL org.opencontainers.image.revision="${VCS_REF}"

# Additional metadata labels
LABEL com.casjaysdev.app.name="caspaste"
LABEL com.casjaysdev.app.version="${VERSION}"
LABEL com.casjaysdev.app.vcs.ref="${VCS_REF}"
LABEL com.casjaysdev.app.build.date="${BUILD_DATE}"

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata curl

# Copy binary from builder
COPY --from=builder /caspaste /usr/local/bin/caspaste

# Set working directory
WORKDIR /data

# Set environment variables for Docker deployment
ENV PORT=80 \
    CASPASTE_DB_DIR=/data/db/sqlite \
    TZ=America/New_York

# Expose HTTP port
EXPOSE 80

# Use non-root user (created by the app itself)
# The app drops privileges to caspaste:caspaste (UID:GID 642:642)

ENTRYPOINT ["caspaste", "--config", "/config/caspaste", "--data", "/data/caspaste", "--logs", "/var/log/caspaste"]
HEALTHCHECK --interval=30s --timeout=10s --start-period=90s --retries=3 CMD caspaste --status || exit 1
