# CasPaste - AI.md Compliance Tasks

## Status: COMPLIANT

AI.md template configured (2026-02-01). All audit issues resolved.

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
- Emojis: feat, fix, docs, style, refactor, perf, test, chore, security, remove, deploy, deps
- COMMIT_MESS must reflect actual `git status` changes
- Recreate if stale (mentions files not in git status)

### Session: 2026-02-01

**Tasks Completed:**
- [x] Copied TEMPLATE.md to AI.md, replaced all placeholders
- [x] Read PART 0-5 of AI.md
- [x] Created .claude/rules/ directory with all 14 rule files
- [x] Committed critical rules to memory
- [x] Verified and deleted AUDIT.AI.md (all issues resolved)

### Compliance Matrix (Current State)

**All Mandatory PARTs (0-33) Implemented:**
- PART 0-5: Rules, Structure, Paths, Config - ✅
- PART 6-8: Modes, Binary, CLI - ✅ (pprof, debug endpoints, --daemon)
- PART 9: Error/Caching - ✅ (ETag support)
- PART 10: Database - ✅ (query timeouts)
- PART 11: Security/Logging - ✅ (Request ID, Security Headers, CSRF, Audit)
- PART 12: Server Config - ✅ (YAML config)
- PART 13: Health - ✅ (`/healthz`, `/api/v1/healthz`)
- PART 14: API - ✅ (OpenAPI/Swagger, GraphQL)
- PART 15: SSL - ✅ (ACME/Let's Encrypt auto-issuance)
- PART 16: Frontend - ✅ (Toast notifications, themes)
- PART 17: Admin Panel - ✅ (Full route hierarchy, dark theme)
- PART 18: Email - ✅ (SMTP support, auto-detect, TLS)
- PART 19: Scheduler - ✅ (Full cron parsing, task management)
- PART 20: GeoIP - ✅ (Country blocking)
- PART 21: Metrics - ✅ (`/metrics` Prometheus)
- PART 22: Backup - ✅ (`--maintenance backup/restore`)
- PART 23: Update - ✅ (`--update check/yes/branch`)
- PART 24-25: Privilege/Service - ✅ (`--service`, privilege drop)
- PART 26: Makefile - ✅ (All targets)
- PART 27: Docker - ✅ (OCI labels, tini, healthcheck)
- PART 28: Workflows - ✅ (release, beta, daily, docker)
- PART 29: Testing - ✅ (tests/ with 3 scripts)
- PART 30: ReadTheDocs - ✅ (docs/, mkdocs.yml)
- PART 31: I18N - ✅ (4 locales)
- PART 32: Tor - ✅ (Hidden service support)
- PART 33: Client - ✅ (CLI + TUI modes)

**Optional PARTs (34-36):** Not implemented (not required)
