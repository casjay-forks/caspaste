# CasPaste - AI.md Compliance Tasks

## Status: AI.md TEMPLATE UPDATED

AI.md has been replaced with the new template from `~/Projects/github/apimgr/TEMPLATE.md`.
Structural compliance tasks previously completed. See compliance matrix below.

### Critical Rules from AI.md (Committed to Memory)

**NEVER Rules:**
- NEVER guess or assume - ALWAYS ask when uncertain
- NEVER install Go locally - ALL builds use Docker (`make dev`, `make local`, `make build`)
- NEVER run binaries on host - use containers for testing
- NEVER store plaintext passwords (use Argon2id)
- NEVER use `mattn/go-sqlite3` (use `modernc.org/sqlite`)
- NEVER use strconv.ParseBool (use `config.ParseBool()`)
- NEVER use inline comments (comments ABOVE code only)
- NEVER use `.yaml` extension (use `.yml`)
- NEVER modify AI.md PARTS 0-36 (except OPTIONAL→REQUIRED)
- NEVER include AI attribution in code/commits

**MUST Rules:**
- MUST use CGO_ENABLED=0 for static binaries
- MUST use parameterized SQL queries
- MUST implement CSRF protection on all forms
- MUST normalize and validate ALL paths
- MUST use Argon2id for passwords, SHA-256 for tokens
- MUST re-read spec before implementing (prevent drift)
- MUST verify before claiming completion
- MUST write `.git/COMMIT_MESS` file (AI cannot git commit)

### Recent Changes (This Session)
- **Updated:** AI.md replaced with new template from TEMPLATE.md
- **Updated:** AI.md PART 3 filled in with project info (caspaste, casjay-forks)
- **Created:** Jenkinsfile per AI.md PART 28 requirement
- **Verified:** .claude/rules/ has all 14 required files
- **Verified:** tests/ has run_tests.sh, docker.sh, incus.sh
- **Verified:** Project structure matches AI.md PART 3
- **Fixed:** Client User-Agent header now set to `caspaste-cli/{version}`
- **Documented:** 8+ compliance gaps identified and documented
- **Added:** Request ID middleware per AI.md PART 11 (`src/web/middleware.go`)
  - Checks X-Request-ID, X-Correlation-ID, X-Trace-ID headers
  - Generates UUID v4 if none provided or invalid
  - Returns X-Request-ID in response headers
  - Provides GetRequestID() helper for logging
- **Added:** Panic Recovery middleware per AI.md PART 6 (`src/web/middleware.go`)
  - Production: Graceful recovery, returns generic 500 error
  - Development (--debug): Verbose response with stack trace
  - Includes request_id in panic logs for tracing
- **Fixed:** JSON API responses now properly formatted per AI.md PART 14
  - All API endpoints use 2-space indented JSON
  - Single trailing newline on all responses
  - Added writeJSON() helper in `src/apiv1/error.go`
- **Fixed:** Text responses now end with trailing newline per AI.md PART 14
  - Updated `src/raw/error.go`
- **Fixed:** Security headers now properly applied per AI.md PART 11
  - Added SecurityHeadersMiddleware to middleware chain
  - Added X-XSS-Protection header support (deprecated but kept for compatibility)
  - Default values: SAMEORIGIN, nosniff, CSP, referrer-policy, permissions-policy, HSTS

### Compliance Check (All PARTs 0-37)

**Mandatory PARTs (0-33):**
- **PART 0-5 (Rules, Structure, Paths, Config):** ✅ Following patterns
- **PART 6-8 (Modes, Binary, CLI):** ⚠️ Missing shell completions (Panic Recovery + pprof debug endpoints implemented)
- **PART 9 (Error/Caching):** ⚠️ Missing ETag, unified response format
- **PART 10 (Database):** ⚠️ Missing query timeouts
- **PART 11 (Security/Logging):** ⚠️ Missing audit.log, CSRF (Request ID + Security Headers implemented)
- **PART 12 (Server Config):** ⚠️ No hot reload, no file watcher
- **PART 13 (Health):** ✅ `/api/healthz` and `/healthz` exist
- **PART 14 (API):** ⚠️ Routes use verbs/camelCase (documented gap)
- **PART 15 (SSL):** ⚠️ Certificate discovery works, no auto-ACME
- **PART 16 (Frontend):** ⚠️ Missing toast notifications (favicon implemented)
- **PART 17 (Admin Panel):** ❌ No src/admin/, no /{admin_path}/ routes
- **PART 18 (Email):** ❌ Not implemented
- **PART 19 (Scheduler):** ⚠️ Basic cleanup only, no task management
- **PART 20 (GeoIP):** ❌ Not implemented
- **PART 21 (Metrics):** ❌ No /metrics, no prometheus support
- **PART 22 (Backup):** ✅ --maintenance backup/restore
- **PART 23 (Update):** ❌ No --update command
- **PART 24-25 (Privilege/Service):** ✅ --service, privilege drop
- **PART 26 (Makefile):** ✅ All targets present
- **PART 27 (Docker):** ✅ OCI labels, tini, healthcheck
- **PART 28 (Workflows):** ✅ release.yml, beta.yml, daily.yml, docker.yml, Jenkinsfile
- **PART 29 (Testing):** ⚠️ tests/ exists with 3 scripts, 100% coverage not enforced
- **PART 30 (ReadTheDocs):** ✅ docs/ has all files, mkdocs.yml, .readthedocs.yaml
- **PART 31 (I18N):** ✅ 4 locales in src/web/data/locale/
- **PART 32 (Tor):** ❌ Not implemented
- **PART 33 (Client):** ⚠️ CLI exists, missing TUI/GUI mode

**Optional PARTs (34-36):** N/A - Not implemented, not required
- PART 34 (Multi-User): Not implemented
- PART 35 (Organizations): Not implemented
- PART 36 (Custom Domains): Not implemented

**Reference Only:**
- PART 37 (IDEA.md Reference): N/A - Read-only reference

## AI.md Mandatory Sections Compliance Matrix

| PART | Section | Status | Notes |
|------|---------|--------|-------|
| 0-1 | AI Rules, Critical Rules | ✅ | Following patterns |
| 2 | License & Attribution | ✅ | MIT license |
| 3 | Project Structure | ✅ | Flat src/, no cmd/internal |
| 4 | OS-Specific Paths | ✅ | {org}/{project} pattern |
| 5 | Configuration | ✅ | server.yml, cli.yml |
| 6 | Application Modes | ✅ | pprof debug endpoints + panic recovery |
| 7 | Binary Requirements | ✅ | CGO_ENABLED=0, embedded assets |
| 8 | Server Binary CLI | ⚠️ | Missing --shell, --pid, --backup, --mode flags |
| 9 | Error Handling & Caching | ⚠️ | Missing ETag support |
| 10 | Database & Cluster | ⚠️ | Missing query timeouts |
| 11 | Security & Logging | ⚠️ | Missing audit.log, CSRF, Request ID |
| 12 | Server Configuration | ⚠️ | YAML config present, no hot reload, no file watcher |
| 13 | Health & Versioning | ✅ | /healthz, /api/healthz |
| 14 | API Structure | ⚠️ | Routes need refactoring (verbs in URLs, camelCase), missing OpenAPI/Swagger |
| 15 | SSL/TLS & Let's Encrypt | ⚠️ | Missing auto ACME issuance |
| 16 | Web Frontend | ⚠️ | Missing toast notifications (favicon implemented) |
| 17 | Admin Panel | ❌ | Not implemented |
| 18 | Email & Notifications | ❌ | Not implemented |
| 19 | Scheduler | ⚠️ | Basic cleanup only |
| 20 | GeoIP | ❌ | Not implemented |
| 21 | Metrics | ❌ | Not implemented |
| 22 | Backup & Restore | ✅ | --maintenance backup/restore |
| 23 | Update Command | ❌ | Not implemented |
| 24-25 | Privilege & Service | ✅ | --service, privilege drop |
| 26 | Makefile | ✅ | All targets present |
| 27 | Docker | ✅ | OCI labels, tini, healthcheck |
| 28 | CI/CD Workflows | ✅ | GitHub Actions + Jenkinsfile present |
| 29 | Testing | ⚠️ | tests/ exists, 100% coverage not verified |
| 30 | ReadTheDocs | ✅ | MkDocs Material |
| 31 | I18N & A11Y | ✅ | 4 locales, WCAG compliant |
| 32 | Tor Hidden Service | ❌ | Not implemented |
| 33 | Client | ⚠️ | CLI exists, missing TUI/GUI mode, shell completions |

**Legend:** ✅ Implemented | ⚠️ Partial | ❌ Not Implemented

**Known gaps** requiring significant new feature work are documented below.

## Phase 1: Source Structure Refactoring - COMPLETED

### 1.1 Move internal/ packages to src/ root - DONE
- [x] `src/internal/config/` → `src/config/`
- [x] `src/internal/storage/` → `src/storage/`
- [x] `src/internal/service/` → `src/service/`
- [x] `src/internal/logger/` → `src/logger/`
- [x] `src/internal/web/` → `src/web/`
- [x] `src/internal/apiv1/` → `src/apiv1/`
- [x] `src/internal/cli/` → `src/cli/`
- [x] `src/internal/netshare/` → `src/netshare/`
- [x] `src/internal/privilege/` → `src/privilege/`
- [x] `src/internal/template/` → `src/template/`
- [x] `src/internal/validation/` → `src/validation/`
- [x] `src/internal/caspasswd/` → `src/caspasswd/`
- [x] `src/internal/raw/` → `src/raw/`
- [x] `src/internal/portutil/` → `src/portutil/`
- [x] `src/internal/lineend/` → `src/lineend/`

### 1.2 Move cmd/ entry points - DONE
- [x] `src/cmd/caspaste/` → `src/server/`
- [x] `src/cmd/caspaste-cli/` → `src/client/`

### 1.3 Update all imports - DONE
- [x] Changed `github.com/casjay-forks/caspaste/src/internal/X` → `github.com/casjay-forks/caspaste/src/X`
- [x] Updated Makefile build paths
- [x] Updated GitHub workflow paths

### 1.4 Remove old directories - DONE
- [x] Removed empty `src/internal/`
- [x] Removed empty `src/cmd/`

## Phase 2: Configuration Changes - COMPLETED

### 2.1 Rename config file - DONE
- [x] Changed `caspaste.yml` → `server.yml` in all references
- [x] Updated config search paths
- [x] Updated CLI config to `cli.yml`

### 2.2 Update config paths - DONE
- [x] Changed `/etc/caspaste/` → `/etc/casjay-forks/caspaste/`
- [x] Changed `~/.config/caspaste/` → `~/.config/casjay-forks/caspaste/`
- [x] Changed `/var/lib/caspaste/` → `/var/lib/casjay-forks/caspaste/`
- [x] Changed `/var/log/caspaste/` → `/var/log/casjay-forks/caspaste/`

## Phase 3: Test Scripts - COMPLETED

### 3.1 Fix tests/ structure - DONE
- [x] Renamed `run-tests.sh` → `run_tests.sh`
- [x] Created `docker.sh` for Docker testing
- [x] Created `incus.sh` for Incus testing
- [x] Updated temp directory to use `${TMPDIR:-/tmp}/casjay-forks/caspaste-XXXXXX`

## Phase 4: Verification - COMPLETED

### 4.1 Build verification - DONE
- [x] `make local` - builds successfully
- [x] Static binary produced (CGO_ENABLED=0)

## Phase 5: Documentation Updates - COMPLETED

### 5.1 Update README.md - DONE
- [x] Updated config paths from `/etc/caspaste/` to `/etc/casjay-forks/caspaste/`
- [x] Updated data paths from `/var/lib/caspaste/` to `/var/lib/casjay-forks/caspaste/`
- [x] Updated log paths from `/var/log/caspaste/` to `/var/log/casjay-forks/caspaste/`
- [x] Updated config file reference from `caspaste.yml` to `server.yml`
- [x] License badge present (MIT)

### 5.2 Update CLAUDE.md - DONE
- [x] Updated project structure (removed cmd/, internal/)
- [x] Updated import paths from `src/internal/X` to `src/X`
- [x] Updated config file reference to `server.yml`
- [x] Updated log paths to use `casjay-forks/caspaste` pattern
- [x] Updated template and theme paths

### 5.3 Verify .claude/ structure - DONE
- [x] `.claude/CLAUDE.md` exists (project memory)
- [x] `.claude/rules/` exists with 14 cheatsheet files

## Phase 6: Additional Fixes - COMPLETED

### 6.1 Fix remaining old path references - DONE
- [x] Fixed `docker/Dockerfile` build paths (src/cmd/* → src/server, src/client)
- [x] Fixed `tests/run_tests.sh` build paths
- [x] Fixed `.claude/rules/frontend-rules.md` embedded asset paths
- [x] Fixed `.claude/rules/testing-rules.md` locale file path
- [x] Fixed `.claude/CLAUDE.md` structure example
- [x] Fixed `src/server/caspaste.go` comment

### 6.2 Update .dockerignore - DONE
- [x] Updated to match AI.md PART 3 requirements
- [x] Excludes .git/, CI/CD files, tests/, docs/, Makefile
- [x] Includes src/, go.mod, go.sum, docker/

### 6.3 Add missing Makefile target - DONE
- [x] Added `make dev` target for quick builds to temp directory
- [x] Per AI.md PART 28 requirements

## Phase 7: ReadTheDocs Documentation - COMPLETED

### 7.1 Create docs/ directory - DONE
- [x] `docs/index.md` - Documentation homepage
- [x] `docs/installation.md` - Installation guide
- [x] `docs/configuration.md` - Configuration reference
- [x] `docs/api.md` - API documentation
- [x] `docs/admin.md` - Admin panel guide
- [x] `docs/cli.md` - CLI reference
- [x] `docs/development.md` - Development/contributing guide
- [x] `docs/requirements.txt` - Python dependencies

### 7.2 Create MkDocs configuration - DONE
- [x] `mkdocs.yml` - MkDocs configuration with Material theme
- [x] `.readthedocs.yaml` - ReadTheDocs build configuration
- [x] Dark/Light/Auto theme toggle
- [x] Navigation structure

## Known Compliance Gaps (Existing Code)

### Inline Comments in Go Code
Per AI.md, comments must be ABOVE code, never inline. The existing codebase has ~220 inline comments, primarily:
- Struct field documentation (Go idiom)
- End-of-line explanations

**Status:** Documented, not fixed (would require significant refactoring of existing code)
**Impact:** Low (code functions correctly, standard Go pattern)

### Built-in Scheduler (PART 19)
Per AI.md PART 19, ALL projects MUST have a built-in scheduler. The current project has a simple ticker-based cleanup but lacks the full scheduler system with:
- Admin panel integration
- Task state persistence
- Catch-up logic for missed tasks
- Full task management

**Status:** Not implemented (significant new feature work)
**Impact:** Medium (cleanup works, but lacks full scheduler features)

### Admin Panel (PART 17)
Per AI.md PART 17, ALL projects MUST have a full admin panel at `/{admin_path}/`. Current project has:
- Basic authentication for private mode
- Auto-generated admin credentials

Missing:
- Full admin panel WebUI
- `src/admin/` package
- Settings management in UI
- Server status dashboard

**Status:** Not implemented (significant new feature work)
**Impact:** Medium (server works, configuration via config file only)

### OpenAPI/Swagger & GraphQL (PART 14)
Per AI.md PART 14, projects MUST have:
- OpenAPI documentation at `/openapi` and `/openapi.json`
- GraphQL endpoint at `/graphql` with `src/graphql/` package
- Both must be in sync with REST API

**Status:** Not implemented
**Impact:** Medium (API works but no interactive docs or GraphQL)

### GeoIP Support (PART 20)
Per AI.md PART 20, ALL projects MUST have built-in GeoIP support using sapics/ip-location-db for:
- Country blocking/allowing
- IP location lookups
- Automatic database downloads and updates

**Status:** Not implemented (significant new feature work)
**Impact:** Low (pastebin functions without GeoIP)

### Prometheus Metrics (PART 21)
Per AI.md PART 21, ALL projects MUST have built-in Prometheus-compatible metrics at `/metrics` including:
- HTTP metrics (requests, duration, size, status codes)
- Database metrics (queries, duration, connections)
- Authentication metrics (attempts, sessions)
- System metrics (CPU, memory, disk, goroutines)

**Status:** Not implemented (significant new feature work)
**Impact:** Medium (no production monitoring capabilities)

### Update Command (PART 23)
Per AI.md PART 23, projects should have `--update` command for self-updates:
- `--update check` - Check for updates
- `--update yes` - Perform update
- `--update branch {name}` - Switch update branch

**Status:** Not implemented
**Impact:** Low (manual updates work, no self-update)

### Automatic SSL/Let's Encrypt (PART 15)
Per AI.md PART 15, ALL projects MUST have built-in Let's Encrypt support with automatic certificate issuance via ACME.

Current project has:
- Can discover and use existing Let's Encrypt certificates from `/etc/letsencrypt/live/`
- Manual certificate path configuration

Missing:
- ACME HTTP-01 challenge handler
- ACME TLS-ALPN-01 challenge handler
- Automatic certificate issuance and renewal
- `autocert` or `certmagic` integration

**Status:** Partial (can use existing certs, no auto-issuance)
**Impact:** Medium (requires external certbot for SSL)

### Tor Hidden Service (PART 32)
Per AI.md PART 32, projects should support Tor hidden services with automatic .onion address generation.

**Status:** Not implemented
**Impact:** Low (clearnet access works, no Tor support)

### Default Favicon (PART 16)
Per AI.md, `/favicon.ico` should be served with an embedded default. Current project uses an empty data URI in templates (`<link rel="icon" href="data:,">`).

**Status:** Partial (prevents 404 but no actual favicon)
**Impact:** Low (browsers don't request missing favicon)

### Debug/Pprof Endpoints (PART 6)
Per AI.md PART 6, `/debug/pprof/` should be available in development mode.

**Status:** Not implemented
**Impact:** Low (development debugging still possible via other means)

### ETag Support (PART 9)
Per AI.md PART 9, ETag support should be implemented for cacheable resources.

**Status:** Not implemented
**Impact:** Low (caching still works via Cache-Control headers)

### Email & Notifications (PART 18)
Per AI.md PART 18, ALL projects MUST have customizable email templates and SMTP support:
- SMTP auto-detection on first run
- Email templates in `{config_dir}/template/email/`
- Password reset, verification emails
- Admin notifications

**Status:** Not implemented (significant new feature work)
**Impact:** Medium (no email notifications, password reset)

### WebUI Toast Notifications (PART 16)
Per AI.md PART 16, projects should use toast notifications instead of JavaScript alerts.

**Status:** Not implemented
**Impact:** Low (current alerts work but not modern UX)

### Audit Logging (PART 11)
Per AI.md PART 11, security events should be logged to `audit.log` in JSON format:
- Login attempts (success/failure)
- Admin actions
- Security-relevant events

Current project has:
- Access logs (multiple formats: apache, nginx, json, text)
- Error logs
- Server logs

Missing:
- Dedicated audit.log for security events
- security.log for security-specific events

**Status:** Partial (logging exists but no dedicated audit log)
**Impact:** Low (events are logged, just not in dedicated audit format)

### Database Query Timeouts (PART 10)
Per AI.md PART 10, ALL database queries MUST have timeouts using context.WithTimeout.

Current implementation uses direct Query/Exec calls without context timeouts.

**Status:** Not implemented
**Impact:** Medium (long-running queries could hang the server)

### CSRF Protection (PART 11)
Per AI.md PART 11, ALL forms MUST have CSRF protection using gorilla/csrf or nosurf.

**Status:** Not implemented
**Impact:** Medium (security vulnerability for form submissions)

### Request ID Tracing (PART 11) - IMPLEMENTED
Per AI.md PART 11, every request MUST have a Request ID (X-Request-ID header) for tracing and debugging.

**Status:** ✅ Implemented in `src/web/middleware.go`
- RequestIDMiddleware checks X-Request-ID, X-Correlation-ID, X-Trace-ID headers
- Generates UUID v4 if none provided or invalid
- Returns X-Request-ID in response headers
- GetRequestID() helper available for logging integration

### Missing CLI Flags (PART 8)
Per AI.md PART 8, these CLI flags are required but missing:
- `--shell completions [SHELL]` - Print shell completion scripts
- `--shell init [SHELL]` - Print shell init command for eval
- `--pid {pid_file}` - Set PID file path
- `--backup {backup_dir}` - Set backup directory
- `--mode {production|development}` - Set application mode

Current implementation has:
- `--maintenance backup` for backup operations (different from --backup dir flag)
- PID file written to fixed location (not configurable)

**Status:** Partial (core flags work, missing shell/pid/backup/mode)
**Impact:** Low (server functions, power user features missing)

### Client TUI/GUI Mode (PART 33)
Per AI.md PART 33, CLI MUST have TUI and GUI modes (NON-NEGOTIABLE):
- TUI mode using bubbletea/lipgloss
- GUI mode using native toolkit (GTK/Cocoa/Win32)
- Setup wizard on first run
- Full server functionality coverage

Current client has:
- Basic CLI commands (new, get, list, login, etc.)
- Config file support (cli.yml)
- User-Agent header (`caspaste-cli/{version}`) - IMPLEMENTED

Missing:
- Interactive TUI mode
- GUI mode
- Setup wizard

**Status:** Partial (CLI works, no TUI/GUI)
**Impact:** Medium (power users have CLI, no interactive mode)

### Shell Completions (PART 8/33)
Per AI.md, ALL binaries (server, agent, client) MUST support shell completions - built into binary, no separate files:
- `--shell completions [SHELL]` - Print completion script to stdout
- `--shell init [SHELL]` - Alias: prints eval-wrapped completions command
- Support for bash, zsh, fish, powershell

**Status:** Not implemented
**Impact:** Medium (users must manually create completion scripts)

### Config Hot Reload (PART 12)
Per AI.md PART 12, config should auto-reload via file watcher (fsnotify):
- Hot-reloadable settings applied automatically on file change
- SIGHUP ignored (file watcher handles reload)
- Admin UI notification for settings requiring restart

**Status:** Not implemented
**Impact:** Low (server restart required for config changes)

### JSON Response Formatting (PART 14) - IMPLEMENTED
Per AI.md PART 14, ALL JSON responses MUST be indented (2 spaces) and end with newline:
```go
data, _ := json.MarshalIndent(response, "", "  ")
w.Write(data)
w.Write([]byte("\n"))
```

**Status:** ✅ Implemented
- Added writeJSON() helper in `src/apiv1/error.go`
- Updated all API endpoints to use 2-space indented JSON
- All responses now end with trailing newline

### API Route Naming Convention (PART 14)
Per AI.md PART 14, API routes MUST follow these rules:
- Versioned: `/api/v1/...`
- Plural nouns: `/api/v1/pastes` not `/api/v1/paste`
- Lowercase with hyphens: `server-info` not `getServerInfo`
- No verbs in URLs: Use HTTP methods instead

Current routes vs AI.md required:
| Current | Should Be | Issue |
|---------|-----------|-------|
| `/api/healthz` | `/api/v1/healthz` | Not versioned |
| `/api/v1/new` | `POST /api/v1/pastes` | Verb in URL |
| `/api/v1/get?id=X` | `GET /api/v1/pastes/{id}` | Verb in URL, query param |
| `/api/v1/list` | `GET /api/v1/pastes` | Verb in URL |
| `/api/v1/getServerInfo` | `GET /api/v1/server/info` | camelCase, verb |

**Status:** Not compliant (breaking API change)
**Impact:** High (requires client updates, deprecation period)

### Unified API Response Format (PART 9)
Per AI.md PART 9, ALL API responses MUST use unified format:

Success:
```json
{"ok": true, "data": {...}}
```

Error:
```json
{"ok": false, "error": "ERROR_CODE", "message": "Human readable message"}
```

Current implementation returns:
- Success: Raw data without wrapper
- Error: `{"code": 400, "error": "Bad Request"}` (different format)

**Status:** Not implemented (breaking API change)
**Impact:** Medium (clients expect current format, migration needed)

### Content Negotiation (PART 14)
Per AI.md PART 14, content negotiation must use Accept header to determine response format:
- `Accept: application/json` → JSON response
- `Accept: text/plain` → Plain text response
- `Accept: text/html` → HTML response
- `.txt` extension → Plain text
- User-Agent based detection for CLI tools

Current implementation:
- Accept-Language header for locale (working)
- No Accept header checking for response format

**Status:** Not implemented
**Impact:** Medium (clients can't request specific formats)

### Panic Recovery Middleware (PART 6) - IMPLEMENTED
Per AI.md, panic recovery middleware MUST be implemented:
- Production: Graceful recovery, logs error, returns 500
- Development: Verbose, full stack in response

**Status:** ✅ Implemented in `src/web/middleware.go`
- PanicRecoveryMiddleware catches all panics
- Production mode: Returns generic "An unexpected error occurred"
- Debug mode (--debug): Returns full stack trace and request ID
- Integrates with Request ID for tracing

### PID File Stale Detection (PART 8)
Per AI.md PART 8, PID file handling requires stale detection:
- CheckPIDFile function to detect stale PID files
- Process alive verification before considering PID valid
- Automatic cleanup of stale PID files

Current implementation:
- Basic PID file write only
- No stale detection

**Status:** Partial (basic write only)
**Impact:** Low (leftover PID files after crash)

## Implemented Features (AI.md Compliance)

### Backup/Restore (PART 22) - IMPLEMENTED
The project has `--maintenance backup` and `--maintenance restore` commands:
- Backup creates tar.gz archives of database and config
- Restore extracts archives to restore server state
- Pre-restore backup is created automatically

### Web Endpoints (PART 16) - IMPLEMENTED
- `/robots.txt` - Dynamically generated with sitemap reference
- `/sitemap.xml` - Dynamically generated
- `/.well-known/security.txt` - RFC 9116 compliant, auto-generated
- `/manifest.json` - PWA manifest embedded

### Accessibility (PART 31) - IMPLEMENTED
- Skip to main content link
- ARIA labels on navigation, forms, buttons
- Role attributes (banner, navigation, main, contentinfo)
- aria-expanded for toggles
- aria-live for dynamic content
- Keyboard navigation support

### Internationalization (PART 31) - IMPLEMENTED
- 4 locales (en, de, bn_IN, ru)
- Accept-Language header support
- Cookie-based language preference
- Fallback to English default

### Service Management (PART 24-25) - IMPLEMENTED
- `--service install/uninstall/start/stop/restart/reload`
- Privilege dropping via `src/privilege/`
- Cross-platform support (Linux, macOS, Windows, BSD)

### SSL/TLS Certificate Discovery (PART 15) - PARTIAL
- Finds existing Let's Encrypt certificates
- Custom cert/key path configuration
- Missing: automatic ACME issuance

### Docker (PART 27) - IMPLEMENTED
- Multi-stage Dockerfile with Alpine base
- OCI labels (all required labels present)
- Tini as init process
- HEALTHCHECK directive
- Proper STOPSIGNAL
- Volume mounts for /config and /data
- Entrypoint script for customization

### Health Check (PART 13) - IMPLEMENTED
- `/healthz` endpoint (HTML for browsers, text for CLI)
- `/api/healthz` endpoint (JSON)
- `--status` CLI flag for health check
- Database and system status reporting

### Security Headers (PART 11) - IMPLEMENTED
- X-Frame-Options: DENY
- X-Content-Type-Options: nosniff
- Content-Security-Policy (comprehensive policy)
- Referrer-Policy: strict-origin-when-cross-origin
- Permissions-Policy (disables geolocation, microphone, camera)
- Strict-Transport-Security (HSTS when SSL enabled)
- CORS configuration

### Logging (PART 11) - IMPLEMENTED
- Multiple log formats (apache, nginx, text, json)
- Log levels (info, warn, error)
- Separate log files (server, error, access, debug)
- Configurable log directory

### Rate Limiting (PART 11) - IMPLEMENTED
- Token bucket rate limiting
- Per-endpoint configuration
- Configurable windows (5min, 15min, 1hour)
- Retry-After header support
