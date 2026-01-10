# CasPaste

A self-hosted, privacy-focused pastebin service for sharing text snippets anonymously. An enhanced fork of [Lenpaste](https://github.com/lcomrade/lenpaste) by Leonid Maslakov.

## About

CasPaste is a modern, secure pastebin service designed for self-hosting. It prioritizes privacy, security, and ease of deployment.

**Key Features:**
- **Privacy-First**: No registration, anonymous sharing, private pastes
- **Secure**: Argon2id hashing, brute force protection, XSS prevention
- **Modern**: Mobile-friendly UI, syntax highlighting, multiple themes
- **File Uploads**: Share images, documents, any file type (50MB max)
- **URL Shortener**: Create short links with QR codes
- **Editable Pastes**: Update pastes after creation
- **API-Ready**: RESTful API with listing, upload, shortening
- **Self-Hosted**: Single binary, SQLite or PostgreSQL
- **Multi-Platform**: Linux, macOS, Windows, BSD (amd64 + arm64)

## Production Deployment

**Note:** CasPaste automatically creates all necessary directories (`db/`, `backups/`, etc.) at startup. No manual setup required.

### Configuration

CasPaste can be configured via:
1. **Config file** (recommended): `caspaste.yml` or `caspaste.yaml`
2. **Command-line flags**: Override config file values

**Config file search locations:**
- `--config` directory (if specified)
- Current directory
- `/etc/caspaste/`

**Auto-generation:** If `--config` directory is specified and no config file exists, a default `caspaste.yml` is created automatically with all settings documented.

### Quick Start

```bash
# Download latest release
wget https://github.com/casjay-forks/caspaste/releases/latest/download/caspaste-linux-amd64

# Make executable
chmod +x caspaste-linux-amd64
mv caspaste-linux-amd64 /usr/local/bin/caspaste

# Option 1: Run with flags (auto-creates all directories)
caspaste --port 8080 --data /var/lib/caspaste

# Option 2: Run with config file (auto-generated if not exists)
caspaste --config /etc/caspaste --data /var/lib/caspaste
# Auto-creates: /etc/caspaste/caspaste.yml with defaults
# Edit /etc/caspaste/caspaste.yml and restart
```

### Docker Deployment

#### Docker Compose (Recommended)

```bash
# Start with SQLite (default)
docker-compose up -d

# Or with PostgreSQL (edit docker-compose.yml first)
# Uncomment postgres service and update caspaste command
docker-compose up -d

# Or with MariaDB (edit docker-compose.yml first)
# Uncomment mariadb service and update caspaste command
docker-compose up -d
```

#### Docker Run (Manual)

```bash
# Pull latest image
docker pull ghcr.io/casjay-forks/caspaste:latest

# Run with SQLite
docker run -d \
  --name caspaste \
  -p 172.17.0.1:64365:80 \
  -v /var/lib/caspaste:/data/caspaste \
  -v /var/lib/caspaste-db:/data/db/sqlite \
  -v /etc/caspaste:/config/caspaste \
  ghcr.io/casjay-forks/caspaste:latest

# Run with PostgreSQL
docker run -d --name caspaste-postgres \
  -e POSTGRES_DB=caspaste \
  -e POSTGRES_USER=caspaste \
  -e POSTGRES_PASSWORD=changeme \
  -v /var/lib/caspaste-postgres:/var/lib/postgresql/data \
  postgres:16-alpine

docker run -d \
  --name caspaste \
  -p 172.17.0.1:64365:80 \
  --link caspaste-postgres:postgres \
  -v /var/lib/caspaste:/data/caspaste \
  -v /etc/caspaste:/config/caspaste \
  ghcr.io/casjay-forks/caspaste:latest \
  --db-driver postgres \
  --db-source "postgres://caspaste:changeme@postgres:5432/caspaste?sslmode=disable"
```

**Docker Configuration:**
- Internal port: `80` (container)
- External port: `172.17.0.1:64365` (Docker bridge)
- Access URL: `http://172.17.0.1:64365`

**Volume Mapping:**
| Container Path | Host Path | Purpose |
|----------------|-----------|---------|
| `/config/caspaste/` | `/etc/caspaste` | Config files |
| `/data/caspaste/` | `/var/lib/caspaste` | Backups, data |
| `/data/db/sqlite/` | `/var/lib/caspaste-db` | SQLite database |

**Files Created:**
- `/etc/caspaste/caspaste.yml` (auto-generated)
- `/var/lib/caspaste-db/caspaste.db` (database)
- `/var/lib/caspaste/backups/` (backups)

### Reverse Proxy Setup (Recommended)

CasPaste works best behind a reverse proxy like nginx or Caddy for TLS termination and advanced features.

#### Nginx Configuration

```nginx
server {
    listen 443 ssl http2;
    server_name paste.example.com;

    ssl_certificate /etc/letsencrypt/live/paste.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/paste.example.com/privkey.pem;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Start CasPaste:

```bash
caspaste --port 8080 --data /var/lib/caspaste
```

#### Caddy Configuration

```caddy
paste.example.com {
    reverse_proxy :8080
}
```

### Production Configuration

Create a systemd service for automatic startup:

```ini
# /etc/systemd/system/caspaste.service
[Unit]
Description=CasPaste Pastebin Service
After=network.target

[Service]
Type=simple
User=caspaste
Group=caspaste
ExecStart=/usr/local/bin/caspaste \
    --port 8080 \
    --data /var/lib/caspaste \
    --config /etc/caspaste \
    --admin-name "Admin" \
    --admin-mail "admin@example.com"
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Install and start:

```bash
# Service auto-creates all directories on startup
sudo systemctl enable caspaste
sudo systemctl start caspaste

# Or use built-in service installer
sudo caspaste --service --install --port 8080 --data /var/lib/caspaste
sudo caspaste --service start
```

### Authentication Setup

Generate password hash:

```bash
go run tools/gen-password/main.go
# Enter username: admin
# Enter password: ********
# Output: admin:$argon2id$v=19$m=65536,t=3,p=4$...
```

Create `/etc/caspaste/passwd`:

```
admin:$argon2id$v=19$m=65536,t=3,p=4$base64salt$base64hash
editor:$argon2id$v=19$m=65536,t=3,p=4$base64salt$base64hash
```

Start with authentication:

```bash
caspaste --port 8080 \
         --data /var/lib/caspaste \
         --config /etc/caspaste \
         --caspasswd-file /etc/caspaste/passwd
```

### Service Management

CasPaste includes built-in service management for automatic startup:

```bash
# Install service (requires sudo/admin)
sudo caspaste --service --install \
     --port 8080 \
     --data /var/lib/caspaste

# Start service
sudo caspaste --service start

# Stop service
sudo caspaste --service stop

# Restart service
sudo caspaste --service restart

# Uninstall service
sudo caspaste --service --uninstall

# Disable from auto-start
sudo caspaste --service --disable
```

**Platform Support:**
- **Linux**: systemd service
- **macOS**: launchd daemon
- **Windows**: Windows Service
- **BSD**: rc.d script

### Health Monitoring

Check service health (no sudo required):

```bash
caspaste --status --data /var/lib/caspaste

# Exit codes:
# 0 = healthy
# 1 = unhealthy
# 2 = degraded
```

### Maintenance Mode

Full disaster recovery backups (entire data directory):

```bash
# Full backup - backs up database + all data files
caspaste --maintenance "backup" --data /var/lib/caspaste
# Creates: /var/lib/caspaste/backups/backup-20260109-143000.tar.gz

# Backup with custom filename
caspaste --maintenance "backup mybackup.tar.gz" --data /var/lib/caspaste

# Restore from latest backup (auto-creates safety backup first)
caspaste --maintenance "restore" --data /var/lib/caspaste

# Restore from specific backup
caspaste --maintenance "restore mybackup.tar.gz" --data /var/lib/caspaste

# Enable maintenance mode (shows 503 to users)
caspaste --maintenance "mode enabled" --data /var/lib/caspaste

# Disable maintenance mode
caspaste --maintenance "mode disabled" --data /var/lib/caspaste
```

**What Gets Backed Up:**
- **Config directory:** All files in `{config_dir}` including `caspaste.yml`
- **Data directory:** All files including `db/caspaste.db` (SQLite primary or cache)
- **External database:** If SQLite database is outside `{data_dir}/db/`, it's included
- **Excludes:** `backups/`, `*.tmp`, `*.lock`, temporary files

**Important:** When using PostgreSQL/MariaDB, the backup includes the synchronized SQLite cache at `db/caspaste.db`, allowing instant disaster recovery without accessing the remote database.

**Directory Structure (auto-created at startup):**
```
/var/lib/caspaste/              # Data directory
├── db/
│   └── caspaste.db             # SQLite primary OR synchronized cache
├── backups/
│   ├── backup-*.tar.gz         # Full disaster recovery backups
│   └── pre-restore-*.tar.gz    # Safety backups before restore
└── .maintenance                # Maintenance mode flag

/etc/caspaste/                  # Config directory (also backed up)
└── caspaste.yml                # Configuration file
```

### Configuration Reference

#### Core Options

| Flag | Description | Default |
|------|-------------|---------|
| `--help` | Show help message | - |
| `--version` | Show version | - |
| `--port` | Port to listen on | `80` |
| `--address` | Full address:port (backward compat.) | `:80` |
| `--data` | Data directory | - |
| `--config` | Config directory | - |
| `--status` | Health check (exit codes: 0/1/2) | - |
| `--service` | Service management (see above) | - |
| `--maintenance` | Maintenance operations (see above) | - |

#### Database

| Flag | Description | Default |
|------|-------------|---------|
| `--db-driver` | Database: `sqlite3` or `postgres` | `sqlite3` |
| `--db-source` | Connection string | **required** |
| `--db-max-open-conns` | Max connections | `25` |
| `--db-cleanup-period` | Cleanup interval | `1m` |

#### Limits

| Flag | Description | Default |
|------|-------------|---------|
| `--body-max-length` | Max paste/file size (bytes) | `52428800` (50MB) |
| `--title-max-length` | Max title length | `100` |
| `--max-paste-lifetime` | Max lifetime | `never` |
| `--get-pastes-per-5min` | View rate limit | `50` |
| `--new-pastes-per-5min` | Create rate limit | `15` |

#### Security

| Flag | Description | Default |
|------|-------------|---------|
| `--caspasswd-file` | Password file for auth | - |
| `--robots-disallow` | Block search engines | `false` |

## API Usage

### Create Text Paste

```bash
curl -X POST https://paste.example.com/api/v1/new \
  -d "body=Hello World" \
  -d "syntax=plaintext" \
  -d "title=My Paste" \
  -d "editable=true" \
  -d "private=false"
```

### Upload File

```bash
curl -X POST https://paste.example.com/api/v1/new \
  -F "file=@image.png" \
  -F "title=My Image"
```

### Create Short URL

```bash
curl -X POST https://paste.example.com/api/v1/new \
  -d "url=true" \
  -d "originalURL=https://example.com/very/long/url" \
  -d "title=Short Link"

# Access via: https://paste.example.com/u/abc12345
```

### Get Paste

```bash
curl https://paste.example.com/api/v1/get?id=abc12345
```

### List Pastes

```bash
curl https://paste.example.com/api/v1/list?limit=10
```

### QR Code

Access QR code for any paste:
```
https://paste.example.com/qr/abc12345
```

Full API documentation: `/docs/apiv1`

## Client Tools

### Command-Line Access

```bash
# Create paste
echo "Hello World" | curl -X POST https://paste.example.com/api/v1/new \
  --data-urlencode body@- \
  -d "syntax=plaintext"

# Get paste
curl https://paste.example.com/api/v1/get?id=abc12345 | jq -r '.body'
```

### Language Libraries

- **Go**: Built-in HTTP client
- **Python**: `requests` library
- **JavaScript**: `fetch` API

See `/docs/api_libs` on your instance for examples.

## Advanced Features

### File Uploads

Upload any file type (images, documents, etc.) up to 50MB:

```bash
# Via web interface: use file upload field
# Via API:
curl -X POST https://paste.example.com/api/v1/new \
  -F "file=@document.pdf" \
  -F "title=My Document"
```

Files are stored in the database and served with correct MIME types.

### URL Shortening

Create short URLs that redirect to any destination:

```bash
curl -X POST https://paste.example.com/api/v1/new \
  -d "url=true" \
  -d "originalURL=https://example.com/very/long/url" \
  -d "title=Short Link"

# Returns: {"id": "abc123", ...}
# Access via: https://paste.example.com/u/abc123
```

### QR Codes

Every paste automatically has a QR code available:

```
https://paste.example.com/qr/{paste-id}
```

Share paste URLs easily via QR codes on mobile devices.

### Editable Pastes

Create pastes that can be updated after creation:

```bash
# Create editable paste
curl -X POST https://paste.example.com/api/v1/new \
  -d "body=Original content" \
  -d "editable=true"

# Edit later (via web interface at /edit/{id})
```

### Private Pastes

Create pastes that don't appear in public listings:

```bash
curl -X POST https://paste.example.com/api/v1/new \
  -d "body=Secret content" \
  -d "private=true"
```

Private pastes are still accessible via direct link but won't show in `/list`.

### Database Backends

CasPaste supports three database backends with **automatic migration** between them:

**SQLite (Default):**
```bash
caspaste --port 8080 --data /var/lib/caspaste
# Database: /var/lib/caspaste/db/caspaste.db
```

**PostgreSQL:**
```bash
caspaste --db-driver postgres \
         --db-source "postgres://caspaste:password@db.example.com:5432/caspaste?sslmode=require" \
         --data /var/lib/caspaste
```

**MariaDB/MySQL:**
```bash
caspaste --db-driver mysql \
         --db-source "caspaste:password@tcp(db.example.com:3306)/caspaste?charset=utf8mb4&parseTime=true" \
         --data /var/lib/caspaste
```

**Automatic Migration:**
When you change database drivers (e.g., SQLite → PostgreSQL), CasPaste automatically:
1. Detects the database driver change
2. Creates a safety backup
3. Migrates all data to the new database
4. Continues operation with the new database

**SQLite Backup/Cache:**
When using PostgreSQL or MariaDB, CasPaste automatically maintains a synchronized SQLite cache at `/data/caspaste/db/caspaste.db` (same path as if SQLite was primary) for:
- **Disaster Recovery:** Full backup always available locally
- **Fast Access:** Read from local SQLite cache
- **Offline Capability:** Continue operations if remote DB is down
- **Automatic Sync:** All writes go to both databases in real-time

**Why this matters for backups:**
- When using PostgreSQL/MariaDB, your disaster recovery backup includes the SQLite cache
- Restore from backup and you can immediately switch back to SQLite if needed
- No need to access the remote database for disaster recovery

### Themes

CasPaste includes 12 built-in themes optimized for mobile-first, readable design:

**Dark Themes:**
- `dracula` (default) - Dracula color scheme
- `nord` - Nord dark palette
- `gruvbox-dark` - Gruvbox dark variant
- `tokyo-night` - Tokyo Night theme
- `catppuccin-mocha` - Catppuccin Mocha
- `one-dark` - One Dark theme
- `dark` - Modern dark theme

**Light Themes:**
- `github-light` - GitHub-inspired light
- `nord-light` - Nord light palette
- `gruvbox-light` - Gruvbox light variant
- `catppuccin-latte` - Catppuccin Latte
- `solarized-light` - Solarized light
- `light` - Modern light theme

**Change default theme:**
```bash
# Via flag
caspaste --ui-default-theme nord

# Via config file
ui:
  default_theme: "nord"
```

**Custom Themes:**
```bash
mkdir -p /etc/caspaste/themes
# Add custom .theme files
caspaste --ui-themes-dir /etc/caspaste/themes
```

### Content Customization

```bash
echo "Welcome to our paste service" > /etc/caspaste/about.txt
echo "Be respectful" > /etc/caspaste/rules.txt

caspaste --server-about /etc/caspaste/about.txt \
         --server-rules /etc/caspaste/rules.txt
```

## Development

### Building from Source

Requirements: Go 1.23+

```bash
git clone https://github.com/casjay-forks/caspaste.git
cd caspaste

# Update dependencies
go mod tidy

# Quick local build
make local

# Build all platforms
make build

# Run tests
make test
```

### Makefile Targets

| Target | Description |
|--------|-------------|
| `make local` | Build for current OS (fast) |
| `make build` | Build all platforms |
| `make release` | Create GitHub release |
| `make docker` | Build/push Docker images |
| `make test` | Run tests |
| `make bump-patch` | Increment version (1.0.0 → 1.0.1) |

### Version Management

Edit `release.txt` to set version:

```bash
echo "1.2.3" > release.txt
make build
```

## Security

CasPaste includes multiple security enhancements:

- **Argon2id Password Hashing**: OWASP-recommended, memory-hard algorithm
- **Brute Force Protection**: 5 failed attempts = 15-minute lockout
- **XSS Prevention**: URL scheme validation, output sanitization
- **Rate Limiting**: Per-IP paste creation and viewing limits
- **Graceful Shutdown**: Proper signal handling on all platforms
- **File Upload Safety**: MIME type validation, size limits
- **Private Pastes**: Control visibility of sensitive content

## Troubleshooting

### Service Won't Start

```bash
# Check status
caspaste --status --db-source /var/lib/caspaste/caspaste.db

# Check logs
journalctl -u caspaste -n 50
```

### Permission Errors

```bash
# Fix data directory ownership (directories auto-created by app)
sudo chown -R caspaste:caspaste /var/lib/caspaste
sudo chmod -R 750 /var/lib/caspaste
```

### Database Issues

```bash
# SQLite: Check file exists and is writable
ls -la /var/lib/caspaste/db/caspaste.db

# PostgreSQL: Test connection
psql "postgres://user:pass@db.example.com:5432/dbname"
```

## Upgrading

```bash
# Create full disaster recovery backup
caspaste --maintenance "backup" --data /var/lib/caspaste

# Stop service
sudo systemctl stop caspaste

# Replace binary
sudo wget https://github.com/casjay-forks/caspaste/releases/latest/download/caspaste-linux-amd64 \
  -O /usr/local/bin/caspaste
sudo chmod +x /usr/local/bin/caspaste

# Start service (auto-creates any new directories needed)
sudo systemctl start caspaste

# If issues occur, restore from backup
# caspaste --maintenance "restore" --data /var/lib/caspaste
```

## License

MIT License - see [LICENSE](LICENSE)

Third-party attributions and original Lenpaste (AGPLv3) attribution - see [LICENSE.md](LICENSE.md)

## Credits

- Original project: [Lenpaste](https://github.com/lcomrade/lenpaste) by Leonid Maslakov
- Fork maintainer: [CasJay](https://github.com/casjay-forks)

## Support

- **Issues**: https://github.com/casjay-forks/caspaste/issues
- **API Docs**: https://paste.example.com/docs/apiv1
- **Changelog**: See [CHANGELOG.md](CHANGELOG.md)
