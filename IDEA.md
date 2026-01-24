# CasPaste

## Project Description

CasPaste is a self-hosted, privacy-focused pastebin service. It provides a fast, secure platform for sharing text snippets and code with syntax highlighting, burn-after-reading, and file upload support. Designed as a single static binary with no external runtime dependencies.

**Target Users:**
- Developers sharing code snippets
- Teams needing private paste hosting
- Privacy-conscious users avoiding public pastebin services
- Self-hosters wanting simple, lightweight paste services

---

## Project-Specific Features

- **Paste Creation**: Create text pastes with optional syntax highlighting
- **File Upload**: Upload files (images, documents) as pastes
- **Privacy Options**: Burn-after-reading, password protection, expiration
- **Themes**: 12+ syntax highlighting themes (dark/light)
- **Localization**: Multi-language support (en, de, bn_IN, ru)
- **PWA**: Progressive Web App support for offline access

---

## Detailed Specification

### Data Models

- **Paste**: id, title, content, language, created_at, expires_at, burn_after_read, password_hash, views
- **File**: id, paste_id, filename, mime_type, size, content_hash

### Business Rules

- Maximum paste size: Configurable (default: unlimited)
- Title max length: 120 characters (configurable)
- Paste lifetime options: never, 10min, 1hour, 1day, 1week, 1month
- Burn-after-reading: Paste deleted after first view
- Password protection: Argon2id hashed
- Rate limiting: Configurable per-endpoint (5min, 15min, 1hour windows)
- Cleanup: Scheduled job removes expired pastes

### Features

- Create paste with optional syntax highlighting
- View paste with rendered code highlighting
- Raw paste access (plain text)
- Download paste as file
- Clone paste (create new from existing)
- File upload with automatic MIME detection
- QR code generation for paste URLs
- Paste listing (when public mode enabled)
- Theme switching (12+ themes)
- Language switching (4 languages)

### Endpoints

- Create paste - see PART 14
- Get paste by ID - see PART 14
- Get raw paste - see PART 14
- List pastes (public mode) - see PART 14
- Get server info - see PART 14
- Health check - see PART 14

### Data Sources

- Database for paste storage - see PART 10 (SQLite, PostgreSQL, MySQL)
- Embedded assets for themes and templates - see PART 16
- Embedded locale files for i18n - see PART 31
