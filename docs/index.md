# CasPaste

A self-hosted, privacy-focused pastebin service for sharing text snippets, files, and short URLs.

## Features

| Feature | Description |
|---------|-------------|
| **Privacy-First** | No registration required, anonymous sharing, private pastes |
| **Secure** | Argon2id hashing, brute force protection, XSS prevention |
| **Modern UI** | Mobile-friendly, syntax highlighting, 12+ themes |
| **File Uploads** | Share images, documents, any file type |
| **URL Shortener** | Create short links with QR codes |
| **Editable Pastes** | Update pastes after creation |
| **Burn After Reading** | One-time view pastes auto-delete after viewing |
| **API-Ready** | RESTful API with listing, upload, URL shortening |
| **Multi-Database** | SQLite, PostgreSQL, MySQL/MariaDB |
| **Multi-Platform** | Linux, macOS, Windows, BSD (amd64 + arm64) |
| **Single Binary** | Static binary with all assets embedded |

## Quick Start

### Docker (Recommended)

```bash
docker run -d \
  --name caspaste \
  -p 8080:80 \
  -v ./rootfs/config:/config \
  -v ./rootfs/data:/data \
  ghcr.io/casjay-forks/caspaste:latest
```

### Binary

```bash
# Download latest release
wget https://github.com/casjay-forks/caspaste/releases/latest/download/caspaste-linux-amd64
chmod +x caspaste-linux-amd64
./caspaste-linux-amd64
```

## Documentation

- [Installation Guide](installation.md) - Detailed setup instructions
- [Configuration Reference](configuration.md) - All configuration options
- [API Documentation](api.md) - REST API reference
- [Admin Panel](admin.md) - Administration guide
- [CLI Reference](cli.md) - Command-line interface
- [Development](development.md) - Contributing and building

## Links

- **Demo:** [https://lp.pste.us](https://lp.pste.us)
- **Source:** [GitHub](https://github.com/casjay-forks/caspaste)
- **Issues:** [GitHub Issues](https://github.com/casjay-forks/caspaste/issues)

## License

MIT License - see [LICENSE.md](https://github.com/casjay-forks/caspaste/blob/main/LICENSE.md)
