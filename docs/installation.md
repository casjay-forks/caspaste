# Installation

## Docker (Recommended)

### Basic Docker Run

```bash
docker run -d \
  --name caspaste \
  -p 172.17.0.1:59093:80 \
  -v ./rootfs/config:/config \
  -v ./rootfs/data:/data \
  -v ./rootfs/backups:/data/backups \
  ghcr.io/casjay-forks/caspaste:latest
```

### Docker Compose

```yaml
version: "3.8"
services:
  caspaste:
    image: ghcr.io/casjay-forks/caspaste:latest
    ports:
      - "172.17.0.1:59093:80"
    volumes:
      - ./rootfs/config:/config
      - ./rootfs/data:/data
      - ./rootfs/backups:/data/backups
    environment:
      - TZ=America/New_York
```

### Docker with PostgreSQL

```yaml
version: "3.8"
services:
  caspaste:
    image: ghcr.io/casjay-forks/caspaste:latest
    ports:
      - "172.17.0.1:59093:80"
    volumes:
      - ./rootfs/config:/config
      - ./rootfs/data:/data
    environment:
      - CASPASTE_DB_DRIVER=postgres
      - CASPASTE_DB_SOURCE=postgres://caspaste:changeme@postgres:5432/caspaste?sslmode=disable
    depends_on:
      - postgres

  postgres:
    image: postgres:16-alpine
    environment:
      - POSTGRES_DB=caspaste
      - POSTGRES_USER=caspaste
      - POSTGRES_PASSWORD=changeme
    volumes:
      - postgres-data:/var/lib/postgresql/data

volumes:
  postgres-data:
```

## Binary Installation

### Download

```bash
# Linux (amd64)
wget https://github.com/casjay-forks/caspaste/releases/latest/download/caspaste-linux-amd64
chmod +x caspaste-linux-amd64
sudo mv caspaste-linux-amd64 /usr/local/bin/caspaste

# macOS (arm64)
wget https://github.com/casjay-forks/caspaste/releases/latest/download/caspaste-darwin-arm64
chmod +x caspaste-darwin-arm64
sudo mv caspaste-darwin-arm64 /usr/local/bin/caspaste
```

### Run

```bash
# Auto-generates config on first run
caspaste

# Or specify directories
caspaste --port 8080 --data /var/lib/casjay-forks/caspaste --config /etc/casjay-forks/caspaste
```

## Service Management

### Install as Service

```bash
# Install (auto-detects platform)
sudo caspaste --service install

# Start
sudo caspaste --service start

# Stop
sudo caspaste --service stop

# Restart
sudo caspaste --service restart

# Status
sudo caspaste --service status

# Uninstall
sudo caspaste --service uninstall
```

### Supported Service Managers

| Platform | Service Type |
|----------|--------------|
| Linux | systemd |
| macOS | launchd |
| Windows | Windows Service |
| BSD | rc.d |

## Platform-Specific Directories

| Directory | Linux (root) | Linux (user) | macOS | Windows |
|-----------|--------------|--------------|-------|---------|
| **Config** | `/etc/casjay-forks/caspaste` | `~/.config/casjay-forks/caspaste` | `~/Library/Application Support/CasPaste/Config` | `%LOCALAPPDATA%\CasPaste\Config` |
| **Data** | `/var/lib/casjay-forks/caspaste` | `~/.local/share/casjay-forks/caspaste` | `~/Library/Application Support/CasPaste` | `%LOCALAPPDATA%\CasPaste\Data` |
| **Logs** | `/var/log/casjay-forks/caspaste` | `~/.local/log/casjay-forks/caspaste` | `~/Library/Logs/CasPaste` | `%LOCALAPPDATA%\CasPaste\Logs` |

## Health Check

```bash
# Check if server is running
caspaste --status
# Exit codes: 0=healthy, 1=unhealthy, 2=degraded
```

## Backup & Restore

```bash
# Create backup
caspaste --maintenance backup

# Restore latest
caspaste --maintenance restore

# Restore specific backup
caspaste --maintenance "restore backup-20240101-120000.tar.gz"
```
