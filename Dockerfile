# syntax=docker/dockerfile:1

# Build stage
FROM golang:alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files first for caching
COPY go.mod go.sum ./

# Copy source code
COPY . .

# Download dependencies (go mod tidy updates go.sum if needed)
RUN go mod tidy && go mod download

# Build arguments
ARG VERSION=unknown
ARG BUILD_DATE
ARG VCS_REF
ARG VCS_URL

# Build the server
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -tags netgo \
    -ldflags "-w -s -X main.Version=${VERSION} -extldflags '-static'" \
    -o /caspaste \
    ./src/cmd/caspaste

# Build the CLI client
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -tags netgo \
    -ldflags "-w -s -X main.Version=${VERSION} -extldflags '-static'" \
    -o /caspaste-cli \
    ./src/cmd/caspaste-cli

# Final stage
FROM alpine:latest

# OCI Standard Labels
# https://github.com/opencontainers/image-spec/blob/main/annotations.md
LABEL org.opencontainers.image.title="CasPaste"
LABEL org.opencontainers.image.description="A simple, fast, and secure self-hosted pastebin service with file uploads, syntax highlighting, and burn-after-reading"
LABEL org.opencontainers.image.authors="CasjaysDev <docker-admin@casjaysdev.pro>"
LABEL org.opencontainers.image.vendor="CasjaysDev"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.url="https://github.com/casjay-forks/caspaste"
LABEL org.opencontainers.image.documentation="https://github.com/casjay-forks/caspaste/blob/main/README.md"
LABEL org.opencontainers.image.source="https://github.com/casjay-forks/caspaste"
LABEL org.opencontainers.image.base.name="docker.io/library/alpine:latest"

# Dynamic labels (set via build args)
ARG VERSION=unknown
ARG BUILD_DATE
ARG VCS_REF
ARG VCS_URL="https://github.com/casjay-forks/caspaste"

LABEL org.opencontainers.image.version="${VERSION}"
LABEL org.opencontainers.image.created="${BUILD_DATE}"
LABEL org.opencontainers.image.revision="${VCS_REF}"

# Additional metadata labels
LABEL com.casjaysdev.app.name="CasPaste"
LABEL com.casjaysdev.app.description="Self-hosted pastebin service"
LABEL com.casjaysdev.app.version="${VERSION}"
LABEL com.casjaysdev.app.maintainer="CasjaysDev <docker-admin@casjaysdev.pro>"
LABEL com.casjaysdev.app.support="https://github.com/casjay-forks/caspaste/issues"
LABEL com.casjaysdev.app.license="MIT"
LABEL com.casjaysdev.app.vcs-url="${VCS_URL}"
LABEL com.casjaysdev.app.vcs-ref="${VCS_REF}"
LABEL com.casjaysdev.app.build-date="${BUILD_DATE}"

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata curl

# Copy binaries from builder
COPY --from=builder /caspaste /usr/local/bin/caspaste
COPY --from=builder /caspaste-cli /usr/local/bin/caspaste-cli

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

CMD ["caspaste", "--config", "/config", "--data", "/data"]
HEALTHCHECK --interval=30s --timeout=10s --start-period=90s --retries=3 CMD caspaste --status || exit 1
