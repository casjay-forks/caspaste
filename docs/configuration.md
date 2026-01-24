# Configuration

CasPaste auto-generates configuration on first run. The config file is the source of truth for all settings.

## Config File Location

- **Linux (root):** `/etc/casjay-forks/caspaste/server.yml`
- **Linux (user):** `~/.config/casjay-forks/caspaste/server.yml`
- **macOS:** `~/Library/Application Support/CasPaste/Config/server.yml`
- **Windows:** `%LOCALAPPDATA%\CasPaste\Config\server.yml`
- **Docker:** `/config/caspaste/server.yml`

## Priority Order

Settings are resolved in this order (highest to lowest):

1. **Command-line flags** (highest priority)
2. **Environment variables** (`CASPASTE_*` prefix)
3. **Config file** (`server.yml`)
4. **Platform-specific defaults** (lowest priority)

## Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `CASPASTE_ADDRESS` | Smart address parsing | `:8080`, `paste.example.com:80` |
| `CASPASTE_CONFIG_DIR` | Config directory | `/config/caspaste` |
| `CASPASTE_DATA_DIR` | Data directory | `/data/caspaste` |
| `CASPASTE_DB_DIR` | Database directory | `/data/db/sqlite` |
| `CASPASTE_LOGS_DIR` | Logs directory | `/data/log/caspaste` |
| `CASPASTE_BACKUP_DIR` | Backup directory | `/data/backups` |
| `CASPASTE_PUBLIC` | Public instance | `true`, `false` |
| `CASPASTE_DB_DRIVER` | Database driver | `sqlite`, `postgres`, `mysql` |
| `CASPASTE_DB_SOURCE` | Database connection string | `postgres://user:pass@host:5432/db` |
| `PORT` | Port (Docker/PaaS) | `80` |

## Config File Structure

```yaml
server:
  public: true                    # true = open, false = auth required
  fqdn: ""                        # Empty = auto-detect from headers/hostname
  listen: all                     # all, ::, 0.0.0.0, or specific IP
  port: ""                        # Empty = auto-detect available port
  title: CasPaste
  tagline: A simple paste service
  description: CasPaste is a simple, fast, and secure paste service
  proxy:
    allowed: []                   # Additional trusted proxies (appended to defaults)
  administrator:
    name: CasPaste Administrator
    email: administrator@{fqdn}   # {fqdn} replaced at runtime
    from: '"CasPaste" <no-reply@{fqdn}>'
  timeouts:
    read: 15
    write: 15
    idle: 60

database:
  driver: sqlite                  # sqlite, postgres, mysql
  source: caspaste.db             # Connection string or filename
  max_open_conns: 25
  max_idle_conns: 5
  cleanup_period: 1m

web:
  ui:
    default_lifetime: never
    default_theme: dark
    themes_dir: ""                # Empty = {data_dir}/web/themes
  content:
    about: ""                     # Empty = auto-generated
    rules: ""                     # Empty = auto-generated
    terms: ""                     # Empty = auto-generated
    security: ""                  # Empty = auto-generated security.txt
  branding:
    logo: ""                      # Path or URL
    favicon: ""                   # Path or URL
  security:
    contact:
      email: security@{fqdn}
      name: Security Team

directories:
  data: /var/lib/casjay-forks/caspaste
  config: /etc/casjay-forks/caspaste
  db: /var/lib/casjay-forks/caspaste/db
  cache: /var/cache/casjay-forks/caspaste
  logs: /var/log/casjay-forks/caspaste

logging:
  level: info                     # info, warn, error
  access:
    stdout: false
    stderr: false
    format: apache                # apache, nginx, text, json
    file: access.log
  error:
    stdout: false
    stderr: true
    format: text
    file: error.log
  server:
    stdout: true
    stderr: false
    format: text
    file: caspaste.log
```

## Database Configuration

### SQLite (Default)

```bash
caspaste --data /var/lib/casjay-forks/caspaste
# Database: /var/lib/casjay-forks/caspaste/db/caspaste.db
```

### PostgreSQL

```bash
caspaste --db-driver postgres \
  --db-source "postgres://user:pass@localhost:5432/caspaste?sslmode=require"
```

### MariaDB/MySQL

```bash
caspaste --db-driver mysql \
  --db-source "user:pass@tcp(localhost:3306)/caspaste?charset=utf8mb4&parseTime=true"
```

## Authentication

CasPaste is **open and public by default** (`server.public: true`).

To require authentication:

```bash
# Via environment
docker run -d -e CASPASTE_PUBLIC=false ghcr.io/casjay-forks/caspaste:latest

# Via config file
# server:
#   public: false
```

On first start with `public: false`, admin credentials are auto-generated and displayed once.

## Trusted Proxies

Private network ranges are **always trusted** for `X-Forwarded-*` headers:

- `10.0.0.0/8`, `172.16.0.0/12`, `192.168.0.0/16` (RFC1918)
- `127.0.0.0/8`, `::1` (loopback)
- `fc00::/7`, `fe80::/10` (IPv6 private/link-local)

Additional proxies can be added via `server.proxy.allowed`.

## Themes

Built-in themes:

- **Dark:** dracula, nord, gruvbox-dark, tokyo-night, catppuccin-mocha, one-dark
- **Light:** github-light, nord-light, gruvbox-light, catppuccin-latte, solarized-light

```yaml
web:
  ui:
    default_theme: nord
```

## Security Features

| Feature | Description |
|---------|-------------|
| **Argon2id Hashing** | OWASP-recommended, memory-hard algorithm |
| **Brute Force Protection** | 5 failed attempts = 15-minute lockout |
| **Secure Sessions** | HttpOnly, SameSite, auto-detect HTTPS |
| **Session Expiry** | 24-hour auto-expire |
