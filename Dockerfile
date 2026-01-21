# syntax=docker/dockerfile:1

# Build stage
FROM golang:alpine AS builder

WORKDIR /build

RUN apk add --no-cache git curl jq ca-certificates tzdata

COPY go.mod go.sum ./
COPY . .

RUN go mod tidy && go mod download

ARG VERSION=dev
ARG BUILD_DATE
ARG VCS_REF
ARG VCS_URL=https://github.com/casjay-forks/caspaste

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -tags netgo \
    -ldflags "-w -s -X main.Version=${VERSION} -extldflags '-static'" \
    -o /caspaste \
    ./src/cmd/caspaste

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -tags netgo \
    -ldflags "-w -s -X main.Version=${VERSION} -extldflags '-static'" \
    -o /caspaste-cli \
    ./src/cmd/caspaste-cli

# Final stage
FROM alpine:latest

ARG VERSION=dev
ARG BUILD_DATE
ARG VCS_REF
ARG VCS_URL=https://github.com/casjay-forks/caspaste

# OCI Labels
LABEL org.opencontainers.image.title="CasPaste" \
      org.opencontainers.image.description="A simple, fast, and secure self-hosted pastebin service" \
      org.opencontainers.image.authors="CasjaysDev <docker-admin@casjaysdev.pro>" \
      org.opencontainers.image.vendor="CasjaysDev" \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.url="https://github.com/casjay-forks/caspaste" \
      org.opencontainers.image.documentation="https://github.com/casjay-forks/caspaste/blob/main/README.md" \
      org.opencontainers.image.source="https://github.com/casjay-forks/caspaste" \
      org.opencontainers.image.base.name="docker.io/library/alpine:latest" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.revision="${VCS_REF}"

LABEL com.casjaysdev.app.name="CasPaste" \
      com.casjaysdev.app.version="${VERSION}" \
      com.casjaysdev.app.maintainer="CasjaysDev <docker-admin@casjaysdev.pro>" \
      com.casjaysdev.app.vcs-url="${VCS_URL}" \
      com.casjaysdev.app.vcs-ref="${VCS_REF}" \
      com.casjaysdev.app.build-date="${BUILD_DATE}"

RUN apk add --no-cache ca-certificates tzdata bash shadow

COPY --from=builder /caspaste /usr/local/bin/caspaste
COPY --from=builder /caspaste-cli /usr/local/bin/caspaste-cli

RUN mkdir -p /config/caspaste /data/caspaste /data/log/caspaste /data/db/sqlite /data/backups

WORKDIR /data/caspaste

ENV PORT=80 \
    CASPASTE_CONFIG_DIR=/config/caspaste \
    CASPASTE_DATA_DIR=/data/caspaste \
    CASPASTE_LOGS_DIR=/data/log/caspaste \
    CASPASTE_DB_DIR=/data/db/sqlite \
    CASPASTE_BACKUP_DIR=/data/backups \
    TZ=America/New_York

EXPOSE 80

VOLUME ["/config", "/data"]

CMD ["caspaste", "--config", "/config/caspaste", "--data", "/data/caspaste", "--logs", "/data/log/caspaste"]

HEALTHCHECK --interval=30s --timeout=10s --start-period=90s --retries=3 CMD caspaste --status || exit 1
