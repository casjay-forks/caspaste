# 📋 CasPaste

A self-hosted, privacy-focused pastebin service for sharing text snippets, files, and short URLs.

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE.md)
[![Docker](https://img.shields.io/badge/Docker-ghcr.io-2496ED?logo=docker)](https://ghcr.io/casjay-forks/caspaste)
[![Demo](https://img.shields.io/badge/Demo-lp.pste.us-green)](https://lp.pste.us)

## 📖 About

CasPaste is a modern, secure pastebin service designed for self-hosting. It prioritizes privacy, security, and ease of deployment.

### ✨ Key Features

| Feature | Description |
|---------|-------------|
| 🔒 **Privacy-First** | No registration required, anonymous sharing, private pastes |
| 🛡️ **Secure** | Argon2id hashing, brute force protection, XSS prevention |
| 📱 **Modern UI** | Mobile-friendly, syntax highlighting, 12+ themes |
| 📁 **File Uploads** | Share images, documents, any file type (50MB max) |
| 🔗 **URL Shortener** | Create short links with QR codes |
| ✏️ **Editable Pastes** | Update pastes after creation |
| 🔥 **Burn After Reading** | One-time view pastes auto-delete after viewing |
| 🔌 **API-Ready** | RESTful API with listing, upload, URL shortening |
| 💾 **Multi-Database** | SQLite, PostgreSQL, MySQL/MariaDB with auto-migration |
| 🖥️ **Multi-Platform** | Linux, macOS, Windows, BSD (amd64 + arm64) |

---

## 🚀 Quick Start

### 📦 Binary Installation

```bash
# Download latest release
wget https://github.com/casjay-forks/caspaste/releases/latest/download/caspaste-linux-amd64
chmod +x caspaste-linux-amd64
sudo mv caspaste-linux-amd64 /usr/local/bin/caspaste

# Run with flags (auto-creates all directories)
caspaste --port 8080 --data /var/lib/caspaste

# Or run with config file (auto-generated if not exists)
caspaste --config /etc/caspaste --data /var/lib/caspaste
```

### 🐳 Docker Installation

```bash
# Pull and run
docker run -d \
  --name caspaste \
  -p 8080:80 \
  -v ./config:/config \
  -v ./data:/data \
  ghcr.io/casjay-forks/caspaste:latest
```

Access at: `http://localhost:8080`

---

## ⚙️ Configuration

CasPaste can be configured via:

1. **📄 Config file** (recommended): `caspaste.yml` or `caspaste.yaml`
2. **🔧 Command-line flags**: Override config file values
3. **🌍 Environment variables**: `CASPASTE_*` prefix

### 📂 Config File Locations

| Priority | Location |
|----------|----------|
| 1️⃣ | `--config` directory (if specified) |
| 2️⃣ | Current working directory |
| 3️⃣ | `/etc/caspaste/` |

> 💡 **Auto-generation**: If `--config` directory is specified and no config file exists, a default `caspaste.yml` is created automatically.

### 🌍 Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `CASPASTE_ADDRESS` | Smart address parsing (see below) | `:8080`, `paste.example.com:80` |
| `CASPASTE_DB_DIR` | SQLite database directory | `/var/lib/caspaste/db` |
| `CASPASTE_PUBLIC` | Public instance (no auth) | `true`, `false` |
| `CASPASTE_PASSWORD_FILE` | Password file path | `/config/.auth` |

#### 🎯 Smart Address Parsing

`CASPASTE_ADDRESS` intelligently configures the server based on the value format:

| Format | Result |
|--------|--------|
| `:8080` | Port only (listen on all interfaces) |
| `paste.example.com:80` | FQDN + port (validates real domain via Public Suffix List) |
| `127.0.0.1` | Listen address only |
| `172.17.0.1:8091` | Listen address + port |

---

## 🐳 Docker Deployment

### Docker Compose (Recommended)

```yaml
version: "3.8"
services:
  caspaste:
    image: ghcr.io/casjay-forks/caspaste:latest
    ports:
      - "8080:80"
    volumes:
      - ./config:/config
      - ./data:/data
    environment:
      - TZ=America/New_York
```

### Docker Run with PostgreSQL

```bash
# Start PostgreSQL
docker run -d --name caspaste-postgres \
  -e POSTGRES_DB=caspaste \
  -e POSTGRES_USER=caspaste \
  -e POSTGRES_PASSWORD=changeme \
  -v ./postgres-data:/var/lib/postgresql/data \
  postgres:16-alpine

# Start CasPaste
docker run -d \
  --name caspaste \
  -p 8080:80 \
  --link caspaste-postgres:postgres \
  -v ./config:/config \
  -v ./data:/data \
  ghcr.io/casjay-forks/caspaste:latest \
  caspaste --config /config --data /data \
  --db-driver postgres \
  --db-source "postgres://caspaste:changeme@postgres:5432/caspaste?sslmode=disable"
```

### 📁 Volume Mapping

| Container Path | Host Path | Purpose |
|----------------|-----------|---------|
| `/config/` | `./config` | Config files (caspaste.yml) |
| `/data/` | `./data` | Data files, backups |
| `/data/db/sqlite/` | Auto-created | SQLite database |

---

## 🔐 Authentication

CasPaste is **open and public by default** (`server.public: true`). For private instances, set `public: false`.

### 🔑 Private Instance Setup

```bash
# Docker with authentication
docker run -d \
  --name caspaste \
  -p 8080:80 \
  -e CASPASTE_PUBLIC=false \
  -v ./config:/config \
  -v ./data:/data \
  ghcr.io/casjay-forks/caspaste:latest
```

On first start, admin credentials are **auto-generated** and displayed in the logs:

```
╔════════════════════════════════════════════════════════════╗
║  CasPaste                                                  ║
╠════════════════════════════════════════════════════════════╣
║  Mode:        Private (authentication required)            ║
║  Username:    admin                                        ║
║  Password:    eoYBn7I9Z&ZHGqCY                             ║
║  ⚠ SAVE THESE CREDENTIALS - shown only once!              ║
╚════════════════════════════════════════════════════════════╝
```

### 🛡️ Security Features

| Feature | Description |
|---------|-------------|
| 🔐 **Argon2id Hashing** | OWASP-recommended, memory-hard algorithm |
| 🚫 **Brute Force Protection** | 5 failed attempts = 15-minute lockout |
| 🍪 **Secure Sessions** | HttpOnly, SameSite, auto-detect HTTPS for Secure flag |
| ⏰ **Session Expiry** | 24-hour auto-expire |

### 🚪 Protected vs Public Routes

When `server.public: false`:

| 🔒 Protected (require login) | 🌐 Always Public |
|------------------------------|------------------|
| `/` - Create paste | `/about/**` - About pages |
| `/list` - Paste list | `/docs/**` - Documentation |
| `/{id}` - View paste | `/terms` - Terms of use |
| `/edit/{id}` - Edit paste | `/login`, `/logout` |
| `/api/v1/*` - All API | `/healthz`, `/robots.txt` |

---

## 💾 Database Backends

CasPaste supports three database backends with **automatic migration**:

### SQLite (Default)

```bash
caspaste --port 8080 --data /var/lib/caspaste
# Database: /var/lib/caspaste/db/caspaste.db
```

### PostgreSQL

```bash
caspaste --db-driver postgres \
  --db-source "postgres://user:pass@db.example.com:5432/caspaste?sslmode=require" \
  --data /var/lib/caspaste
```

### MariaDB/MySQL

```bash
caspaste --db-driver mysql \
  --db-source "user:pass@tcp(db.example.com:3306)/caspaste?charset=utf8mb4&parseTime=true" \
  --data /var/lib/caspaste
```

### 🔄 SQLite Backup/Cache

When using PostgreSQL or MariaDB, CasPaste automatically maintains a **synchronized SQLite cache** for:

| Benefit | Description |
|---------|-------------|
| 🆘 **Disaster Recovery** | Full backup always available locally |
| ⚡ **Fast Access** | Read from local SQLite cache |
| 📴 **Offline Capability** | Continue operations if remote DB is down |
| 🔄 **Automatic Sync** | All writes go to both databases in real-time |

> 💡 The SQLite cache location can be set via `CASPASTE_DB_DIR` environment variable.

---

## 🔥 Burn After Reading

Create pastes that **auto-delete after being viewed once**:

```bash
# Via API
curl -X POST https://paste.example.com/api/v1/new \
  -d "body=Secret message" \
  -d "oneUse=true"

# Via CLI
echo "Secret" | caspaste-cli new --one-use
```

The paste is **automatically deleted** immediately after the first view - no special parameters needed when viewing.

---

## 🌐 API Usage

### 📝 Create Text Paste

```bash
curl -X POST https://paste.example.com/api/v1/new \
  -d "body=Hello World" \
  -d "syntax=plaintext" \
  -d "title=My Paste"
```

### 📁 Upload File

```bash
curl -X POST https://paste.example.com/api/v1/new \
  -F "file=@image.png" \
  -F "title=My Image"
```

### 🔗 Create Short URL

```bash
curl -X POST https://paste.example.com/api/v1/new \
  -d "url=true" \
  -d "originalURL=https://example.com/very/long/url"
# Access via: https://paste.example.com/u/abc123
```

### 📖 Get Paste

```bash
curl https://paste.example.com/api/v1/get?id=abc12345
```

### 📋 List Pastes

```bash
curl https://paste.example.com/api/v1/list?limit=10
```

### 📱 QR Code

```
https://paste.example.com/qr/{paste-id}
```

> 📚 Full API documentation available at `/docs/apiv1` on your instance.

---

## 💻 CLI Client

CasPaste includes a command-line client for terminal-based usage.

### 📥 Installation

```bash
# From release
wget https://github.com/casjay-forks/caspaste/releases/latest/download/caspaste-cli-linux-amd64
chmod +x caspaste-cli-linux-amd64
sudo mv caspaste-cli-linux-amd64 /usr/local/bin/caspaste-cli

# Configure
caspaste-cli login
```

### 📋 Usage Examples

```bash
# Create paste from stdin
echo "Hello World" | caspaste-cli new

# Create paste from file with syntax highlighting
caspaste-cli new -f script.py -s python -t "My Script"

# Create private one-time paste
cat secret.txt | caspaste-cli new -p --one-use -l 1h

# Get a paste
caspaste-cli get abc123

# List recent pastes
caspaste-cli list -n 20
```

### 🔧 CLI Commands

| Command | Description |
|---------|-------------|
| `login` | Configure server and credentials |
| `new`, `create`, `paste` | Create a new paste |
| `get`, `show`, `view` | Get a paste by ID |
| `list`, `ls` | List pastes |
| `health` | Check server health |

---

## 🎨 Themes

CasPaste includes 12+ built-in themes:

### 🌙 Dark Themes

| Theme | Description |
|-------|-------------|
| `dracula` | Dracula color scheme (default) |
| `nord` | Nord dark palette |
| `gruvbox-dark` | Gruvbox dark variant |
| `tokyo-night` | Tokyo Night theme |
| `catppuccin-mocha` | Catppuccin Mocha |
| `one-dark` | One Dark theme |

### ☀️ Light Themes

| Theme | Description |
|-------|-------------|
| `github-light` | GitHub-inspired light |
| `nord-light` | Nord light palette |
| `gruvbox-light` | Gruvbox light variant |
| `catppuccin-latte` | Catppuccin Latte |
| `solarized-light` | Solarized light |

Change default theme:
```yaml
# caspaste.yml
ui:
  default_theme: "nord"
```

---

## 🔧 Service Management

CasPaste includes built-in service management:

```bash
# Install as service
sudo caspaste --service --install --port 8080 --data /var/lib/caspaste

# Manage service
sudo caspaste --service start
sudo caspaste --service stop
sudo caspaste --service restart

# Uninstall
sudo caspaste --service --uninstall
```

### 🖥️ Platform Support

| Platform | Service Type |
|----------|--------------|
| 🐧 Linux | systemd service |
| 🍎 macOS | launchd daemon |
| 🪟 Windows | Windows Service |
| 😈 BSD | rc.d script |

---

## 🏥 Health Monitoring

```bash
# Check health (exit codes: 0=healthy, 1=unhealthy, 2=degraded)
caspaste --status --data /var/lib/caspaste

# Docker health check (built-in)
docker inspect --format='{{.State.Health.Status}}' caspaste
```

---

## 💾 Backup & Restore

### 📦 Create Backup

```bash
caspaste --maintenance "backup" --data /var/lib/caspaste
# Creates: /var/lib/caspaste/backups/backup-YYYYMMDD-HHMMSS.tar.gz
```

### 📂 Restore Backup

```bash
# Restore from latest backup (auto-creates safety backup first)
caspaste --maintenance "restore" --data /var/lib/caspaste

# Restore from specific backup
caspaste --maintenance "restore mybackup.tar.gz" --data /var/lib/caspaste
```

### 📋 What Gets Backed Up

| Item | Description |
|------|-------------|
| ⚙️ Config directory | All files including `caspaste.yml` |
| 💾 Data directory | All files including SQLite database |
| 🔄 SQLite cache | When using PostgreSQL/MySQL |

> ⚠️ When using PostgreSQL/MariaDB, the backup includes the synchronized SQLite cache, allowing instant disaster recovery without accessing the remote database.

---

## 🔒 Reverse Proxy Setup

### Nginx

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

### Caddy

```caddy
paste.example.com {
    reverse_proxy :8080
}
```

---

## 🛠️ Development

### 📋 Requirements

- Go 1.23+
- Make (optional)

### 🔨 Building from Source

```bash
git clone https://github.com/casjay-forks/caspaste.git
cd caspaste

# Update dependencies
go mod tidy

# Build for current platform
make local

# Build all platforms
make build

# Run tests
make test
```

### 📋 Makefile Targets

| Target | Description |
|--------|-------------|
| `make local` | Build for current OS |
| `make build` | Build all platforms |
| `make docker` | Build Docker images |
| `make test` | Run tests |

---

## 🛡️ Security

| Feature | Description |
|---------|-------------|
| 🔐 **Argon2id Hashing** | OWASP-recommended, memory-hard algorithm |
| 🚫 **Brute Force Protection** | 5 failed attempts = 15-minute lockout |
| 🛡️ **XSS Prevention** | URL scheme validation, output sanitization |
| ⏱️ **Rate Limiting** | Per-IP paste creation and viewing limits |
| 🍪 **Secure Cookies** | HttpOnly, SameSite, auto-detect HTTPS |
| 📁 **File Safety** | MIME type validation, size limits |
| 🔒 **Private Pastes** | Control visibility of sensitive content |

---

## 📚 Configuration Reference

### 🔧 Core Options

| Flag | Description | Default |
|------|-------------|---------|
| `--help` | Show help | - |
| `--version` | Show version | - |
| `--port` | Port to listen on | `80` |
| `--data` | Data directory | - |
| `--config` | Config directory | - |

### 💾 Database Options

| Flag | Description | Default |
|------|-------------|---------|
| `--db-driver` | `sqlite3`, `postgres`, `mysql` | `sqlite3` |
| `--db-source` | Connection string | **required** |
| `--db-max-open-conns` | Max connections | `25` |

### 📏 Limits

| Flag | Description | Default |
|------|-------------|---------|
| `--body-max-length` | Max paste/file size | `52428800` (50MB) |
| `--title-max-length` | Max title length | `100` |
| `--max-paste-lifetime` | Max lifetime | `never` |
| `--get-pastes-per-5min` | View rate limit | `50` |
| `--new-pastes-per-5min` | Create rate limit | `15` |

---

## 🐛 Troubleshooting

### ❌ Service Won't Start

```bash
# Check status
caspaste --status --data /var/lib/caspaste

# Check logs (Linux)
journalctl -u caspaste -n 50
```

### 🔑 Permission Errors

```bash
sudo chown -R caspaste:caspaste /var/lib/caspaste
sudo chmod -R 750 /var/lib/caspaste
```

### 💾 Database Issues

```bash
# SQLite: Check file permissions
ls -la /var/lib/caspaste/db/caspaste.db

# PostgreSQL: Test connection
psql "postgres://user:pass@db.example.com:5432/caspaste"
```

---

## ⬆️ Upgrading

```bash
# 1. Create backup
caspaste --maintenance "backup" --data /var/lib/caspaste

# 2. Stop service
sudo systemctl stop caspaste

# 3. Replace binary
sudo wget https://github.com/casjay-forks/caspaste/releases/latest/download/caspaste-linux-amd64 \
  -O /usr/local/bin/caspaste
sudo chmod +x /usr/local/bin/caspaste

# 4. Start service
sudo systemctl start caspaste
```

---

## 📜 License

MIT License - see [LICENSE.md](LICENSE.md) for details and third-party attributions.

## 👥 Credits

- Maintainer: [CasjaysDev](https://github.com/casjay-forks)

## 🆘 Support

- 🌐 **Demo**: https://lp.pste.us
- 🐛 **Issues**: https://github.com/casjay-forks/caspaste/issues
- 📚 **API Docs**: https://lp.pste.us/docs/apiv1
