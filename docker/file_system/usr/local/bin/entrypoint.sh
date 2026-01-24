#!/usr/bin/env bash
# CasPaste container entrypoint
# Per AI.md PART 27: NEVER modify ENTRYPOINT/CMD - customize here

set -e

# Create directories if they don't exist
mkdir -p /config/caspaste /data/caspaste /data/log/caspaste

# Run the server
exec caspaste \
    --config "${CASPASTE_CONFIG_DIR:-/config/caspaste}" \
    --data "${CASPASTE_DATA_DIR:-/data/caspaste}" \
    --log "${CASPASTE_LOGS_DIR:-/data/log/caspaste}" \
    "$@"
