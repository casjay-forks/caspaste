# CasPaste

## Project Description

CasPaste is a self-hosted, privacy-focused pastebin service with URL shortening capabilities. It provides a fast, secure platform for sharing text snippets, code with syntax highlighting, files, and short URLs. Designed as a single static binary with all assets embedded and no external runtime dependencies.

**Target Users:**
- Developers sharing code snippets and files
- Teams needing private paste hosting
- Privacy-conscious users avoiding public pastebin services
- Self-hosters wanting simple, lightweight paste and URL shortening services
- Marketing/business users needing URL shortening

---

## Project-Specific Features

- **Paste Creation**: Create text pastes with syntax highlighting (Chroma lexers)
- **File Upload**: Upload any file type with automatic MIME detection
- **URL Shortener**: Create short links that redirect to original URLs
- **Privacy Options**: Burn-after-reading (one-use), private pastes, expiration
- **Editable Pastes**: Update paste content after creation
- **Public/Private Mode**: Server-wide authentication for private deployments
- **Themes**: Syntax highlighting themes (dark/light)
- **Localization**: Multi-language support (en, de, bn_IN, ru)
- **PWA**: Progressive Web App support with service worker

---

## Detailed Specification

### Data Models

- **Paste**: Single unified model for text, files, and URLs
  - id: 8-char cryptographically random identifier
  - title: Optional paste title
  - body: Content (text, base64 file data, or empty for URLs)
  - syntax: Language for syntax highlighting
  - create_time: Unix timestamp of creation
  - delete_time: Unix timestamp for expiration (0 = never)
  - one_use: Burn-after-reading flag (deleted after first view)
  - author: Optional author name
  - author_email: Optional author email
  - author_url: Optional author website
  - is_file: True if this is a file upload
  - file_name: Original filename for file uploads
  - mime_type: MIME type for file uploads
  - is_editable: True if paste can be edited after creation
  - is_private: True if paste is not listed publicly
  - is_url: True if this is a URL shortener entry
  - original_url: Destination URL for shortener entries

### Business Rules

- Paste ID: 8-character cryptographically random string
- Title max length: Configurable (default: 120 characters)
- Body max length: Configurable (default: unlimited)
- Paste lifetime options: never (0), 10min, 1hour, 1day, 1week, 1month, 1year
- Burn-after-reading: Paste deleted immediately after first successful view
- Private pastes: Excluded from public listing
- URL shortener: Redirects to original_url when accessed
- Rate limiting: Configurable per-endpoint with 5min/15min/1hour windows
- Brute force protection: 5 failed login attempts = 15 minute lockout
- Cleanup: Background job removes expired pastes
- Public mode: Open access, no authentication required
- Private mode: Argon2id password authentication required

### Features

- Create text paste with optional syntax highlighting
- Create file upload paste (any file type, automatic MIME detection)
- Create URL shortener entry
- View paste with rendered syntax highlighting
- View raw paste (plain text, no formatting)
- Download paste/file as attachment
- Clone paste (create new paste from existing content)
- Edit paste (if marked as editable)
- List pastes (public mode or authenticated)
- Theme switching (dark/light themes)
- Language switching (4 locales)
- QR code generation for paste URLs
- Server information and health status

### Endpoints

- Create paste (text, file, or URL) - see PART 14
- Get paste by ID - see PART 14
- Get raw paste content - see PART 14
- List pastes (with optional filters) - see PART 14
- Get server info and configuration - see PART 14
- Health check endpoint - see PART 13

### External API Compatibility

CasPaste provides create-only compatibility with popular pastebin services, allowing existing tools and scripts to work by simply changing the URL.

**Supported Services:**

| Service | Endpoint | Field Names | Response |
|---------|----------|-------------|----------|
| sprunge.us | POST /sprunge | `sprunge` | Plain text URL |
| ix.io | POST /ix | `f:1`, `f:0`, `f` | Plain text URL |
| termbin | POST /termbin or /nc | Raw body | Plain text URL |
| pastebin.com | POST /api/api_post.php | `api_paste_code`, `api_paste_format`, etc. | Plain text URL |
| stikked/stiqued | POST /api/create | `text`, `code`, `data`, `lang`, `title` | JSON |
| microbin | POST /upload or /p | `content`, `text`, `editordata` | Plain text URL or JSON |
| lenpaste | POST /api/v1/new | `body`, `title`, `syntax`, etc. | JSON |
| Generic | POST /compat or /paste | Multiple field names | Plain text URL or JSON |

**Example Usage:**
```bash
# sprunge-style
echo "hello world" | curl -F 'sprunge=<-' https://yourserver.com/sprunge

# ix.io-style
echo "hello world" | curl -F 'f:1=<-' https://yourserver.com/ix

# termbin-style
echo "hello world" | curl -X POST --data-binary @- https://yourserver.com/termbin

# Generic (accepts multiple field names)
curl -F 'text=hello world' https://yourserver.com/compat
```

**Notes:**
- Anonymous paste creation is supported (no auth required by default)
- All compatibility endpoints use rate limiting
- View/list/delete operations use standard CasPaste API routes

### Data Sources

- Database for paste storage - see PART 10 (SQLite default, PostgreSQL, MySQL)
- SQLite backup/cache pool for resilience
- Embedded templates for web UI - see PART 16
- Embedded static assets (CSS, JS) - see PART 16
- Embedded locale files for i18n - see PART 31
- Chroma lexers for syntax highlighting (embedded)
