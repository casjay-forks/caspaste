# CasPaste - AI.md Compliance Tasks

## Status: PRODUCTION-READY

Project is fully implemented and AI.md compliant. Last verified: 2026-02-02.

### Critical Rules Committed to Memory

**NEVER Rules:**
- NEVER guess or assume - ALWAYS ask when uncertain
- NEVER install Go locally - ALL builds use Docker (`make dev`, `make local`, `make build`)
- NEVER run binaries on host - use containers for testing
- NEVER store plaintext passwords (use Argon2id)
- NEVER use `mattn/go-sqlite3` (use `modernc.org/sqlite`)
- NEVER use `strconv.ParseBool()` (use `config.ParseBool()`)
- NEVER use inline comments (comments ABOVE code only)
- NEVER use `.yaml` extension (use `.yml`)
- NEVER modify AI.md PARTS 0-36 (except OPTIONAL->REQUIRED)
- NEVER include AI attribution in code/commits
- NEVER run git add/commit/push (write .git/COMMIT_MESS instead)
- NEVER create forbidden files (SUMMARY.md, COMPLIANCE.md, NOTES.md, etc.)
- NEVER put Dockerfile in root (use docker/Dockerfile)
- NEVER use CGO (CGO_ENABLED=0 always)
- NEVER use external cron (use internal scheduler)
- NEVER create premium/enterprise tiers (all features free)
- NEVER use Makefile in CI/CD (use explicit commands)

**MUST Rules:**
- MUST use CGO_ENABLED=0 for static binaries
- MUST use parameterized SQL queries
- MUST implement CSRF protection on all forms
- MUST normalize and validate ALL paths
- MUST use Argon2id for passwords, SHA-256 for tokens
- MUST re-read spec before implementing (prevent drift)
- MUST verify before claiming completion
- MUST write `.git/COMMIT_MESS` file (AI cannot git commit)
- MUST read file before editing
- MUST search before create
- MUST test before commit
- MUST complete current task before starting next
- MUST use MIT License
- MUST build all 8 platforms (linux/darwin/windows/freebsd x amd64/arm64)
- MUST have comments ABOVE code, never inline

**COMMIT Rules:**
- Write commit message to `.git/COMMIT_MESS` file
- Format: `{emoji} Title (max 64 chars) {emoji}\n\n{description}\n\n- bullets`
- Emojis: ✨ feat, 🐛 fix, 📝 docs, 🎨 style, ♻️ refactor, ⚡ perf, ✅ test, 🔧 chore, 🔒 security, 🗑️ remove, 🚀 deploy, 📦 deps
- COMMIT_MESS must reflect actual `git status` changes
- Recreate if stale (mentions files not in git status)

### Current Session: 2026-02-02

**Completed:**
- [x] Copied TEMPLATE.md to AI.md, replaced all placeholders
- [x] Created .claude/rules/ with 14 rule files
- [x] Read and analyzed entire codebase (109 Go files, 1.3MB)
- [x] Verified compliance with AI.md PARTS 0-33

### Codebase Status (Verified)

**Core Implementation (109 Go files):**
```
src/
├── server/          # Main server (~900 lines)
├── client/          # CLI client with TUI
├── apiv1/           # REST API v1
├── web/             # Web UI, templates, themes, locales
├── graphql/         # GraphQL API
├── swagger/         # OpenAPI/Swagger
├── storage/         # SQLite, PostgreSQL, MySQL
├── config/          # Configuration management
├── admin/           # Admin panel
├── caspasswd/       # Argon2id authentication
├── netshare/        # Rate limiting
├── audit/           # Security audit logging
├── metric/          # Prometheus metrics (singular ✅)
├── path/            # Path utilities (singular ✅)
├── completion/      # Shell completions (singular ✅)
├── scheduler/       # Background task scheduler
├── email/           # Email support
├── geoip/           # GeoIP blocking
├── tor/             # Tor hidden service
├── updater/         # Self-update
├── ssl/             # ACME/Let's Encrypt
├── service/         # systemd/launchd/Windows service
├── privilege/       # UID/GID management
├── tui/             # Terminal UI (bubbletea)
├── display/         # Display mode detection
└── ... (more packages)
```

**Infrastructure:**
- docker/Dockerfile: Multi-stage, alpine, tini, STOPSIGNAL ✅
- tests/: run_tests.sh, docker.sh, incus.sh ✅
- docs/: 7 markdown files for ReadTheDocs ✅
- .github/workflows/: docker, release, beta, daily ✅

**Compliance Matrix:**
- PART 0-5: AI rules, structure, paths, config ✅
- PART 6-8: Modes, binary, CLI (all flags) ✅
- PART 9: Error handling, ETag caching ✅
- PART 10: Database (modernc.org/sqlite) ✅
- PART 11: Security (Argon2id, CSRF, headers) ✅
- PART 12: Server config (YAML) ✅
- PART 13: Health endpoints ✅
- PART 14: API structure (REST, GraphQL, OpenAPI) ✅
- PART 15: SSL/ACME ✅
- PART 16: Frontend (SSR, themes, locales) ✅
- PART 17: Admin panel ✅
- PART 18: Email ✅
- PART 19: Scheduler (internal, not cron) ✅
- PART 20: GeoIP ✅
- PART 21: Metrics (Prometheus) ✅
- PART 22: Backup/restore ✅
- PART 23: Update command ✅
- PART 24-25: Privilege/service ✅
- PART 26: Makefile ✅
- PART 27: Docker (OCI labels, tini) ✅
- PART 28: CI/CD workflows ✅
- PART 29: Testing ✅
- PART 30: ReadTheDocs ✅
- PART 31: I18N (4 locales) ✅
- PART 32: Tor hidden service ✅
- PART 33: Client (CLI + TUI) ✅

**Optional (34-36):** Not implemented (not required)

### Pending Changes (from git status)

The git status shows modified files from previous work session:
- Package renames: completions→completion, metrics→metric, paths→path
- Various source file updates
- These should be committed with appropriate message
