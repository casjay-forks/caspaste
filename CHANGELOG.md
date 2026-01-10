## 🚀 Changelog: 2026-01-09 - Major Feature Release 🚀

### New Features
- ✅ YAML configuration file support (`caspaste.yml` / `caspaste.yaml`)
- ✅ Paste listing API endpoint (`/api/v1/list`) with web UI (`/list`)
- ✅ File uploads (images, documents, any file type) - 50MB max
- ✅ URL shortener with `/u/{id}` redirect endpoint
- ✅ QR code generation for pastes (`/qr/{id}`)
- ✅ Editable pastes (`/edit/{id}` for updating after creation)
- ✅ Private/public paste toggle (private pastes not listed)
- ✅ Service management (`--service start|stop|restart|reload|--install|--uninstall|--disable`)
- ✅ Maintenance mode (`--maintenance backup|restore|mode`)
- ✅ Health check command (`--status` with exit codes 0/1/2)
- ✅ Auto-versioning system with `release.txt`
- ✅ Modern CLI flags (`--port`, `--data`, `--config`)
- ✅ Cross-platform graceful shutdown (Windows/macOS/BSD/Linux)
- ✅ Auto-directory creation at startup
- ✅ Full disaster recovery backup/restore
- ✅ 12 built-in themes (6 dark, 6 light) - mobile-first design

### Security Enhancements
- ✅ Argon2id password hashing (replaces plain text)
- ✅ Brute force protection (5 attempts = 15min lockout)
- ✅ XSS prevention (Author URL validation)
- ✅ Enhanced reverse proxy support (RFC 7239, Cloudflare, etc.)
- ✅ IP spoofing protection

### Infrastructure
- ✅ Multi-arch Docker images (amd64/arm64) on ghcr.io
- ✅ GitHub Actions for CI/CD with dev/release tags
- ✅ Simplified Makefile for local development
- ✅ Production-first README
- ✅ Clean source archives (no VCS files)

### Database Changes
- 📁 Database filename: `caspaste.db` (renamed from lenpaste.db)
- 📁 Database location: `{dataDir}/db/caspaste.db`
- 📁 Auto-created structure: `db/`, `backups/`
- 📁 New columns: `is_file`, `file_name`, `mime_type`, `is_editable`, `is_private`, `is_url`, `original_url`
- 📁 Default max size: 50MB (increased from 20KB)
- 📁 MariaDB/MySQL support added
- 📁 Automatic migration between SQLite, PostgreSQL, MariaDB
- 📁 SQLite backup/cache: When using PostgreSQL/MariaDB, SQLite at `db/caspaste.db` serves as real-time synchronized cache

### License Change
- 📄 Changed from GNU AGPLv3 to MIT License
- 📄 Original Lenpaste (AGPLv3) attribution maintained in LICENSE.md
- 📄 All 3rd party licenses documented in LICENSE.md

### Breaking Changes
- None (all changes are backward compatible)

### Files Added
- `src/internal/config/yaml.go` - YAML configuration support
- `src/internal/apiv1/api_list.go` - Paste listing API
- `src/internal/web/web_list.go` - Paste listing web handler
- `src/internal/web/web_url.go` - URL shortener redirect handler
- `src/internal/web/web_qr.go` - QR code generation handler
- `src/internal/web/web_edit.go` - Editable paste handler
- `src/src/internal/web/data/list.tmpl` - List page template
- `src/internal/web/middleware.go` - Maintenance mode middleware
- `src/internal/service/service.go` - Service manager
- `src/internal/service/service_linux.go` - systemd support
- `src/internal/service/service_darwin.go` - launchd support
- `src/internal/service/service_windows.go` - Windows Service
- `src/internal/service/service_bsd.go` - rc.d support
- `src/internal/caspasswd/bruteforce.go` - Brute force protection
- `src/tools/gen-password/main.go` - Password hash generator
- `release.txt` - Version file (1.0.0)
- `src/src/internal/web/data/theme/dracula.theme` - Dracula dark theme (default)
- `src/src/internal/web/data/theme/nord.theme` - Nord dark theme
- `src/src/internal/web/data/theme/gruvbox-dark.theme` - Gruvbox dark
- `src/src/internal/web/data/theme/tokyo-night.theme` - Tokyo Night
- `src/src/internal/web/data/theme/catppuccin-mocha.theme` - Catppuccin Mocha
- `src/src/internal/web/data/theme/one-dark.theme` - One Dark
- `src/src/internal/web/data/theme/github-light.theme` - GitHub light
- `src/src/internal/web/data/theme/nord-light.theme` - Nord light
- `src/src/internal/web/data/theme/gruvbox-light.theme` - Gruvbox light
- `src/src/internal/web/data/theme/catppuccin-latte.theme` - Catppuccin Latte
- `src/src/internal/web/data/theme/solarized-light.theme` - Solarized Light
- `LICENSE.md` - Third-party licenses and attributions

### Files Modified
- `LICENSE` - Changed from AGPLv3 to MIT
- `src/cmd/caspaste/caspaste.go` - Major refactor with new commands, YAML config, auto-migration
- `src/internal/cli/cli.go` - Support for `--flag` syntax
- `src/internal/caspasswd/caspasswd.go` - Argon2id + bcrypt support
- `src/internal/netshare/netshare_host.go` - Enhanced proxy headers
- `src/internal/netshare/netshare_paste.go` - URL validation, file upload support
- `src/internal/storage/storage.go` - Dual-database system (primary + SQLite cache), MariaDB support
- `src/internal/storage/storage_paste.go` - PasteList, PasteUpdate, real-time sync to backup
- `src/internal/apiv1/api.go` - Brute force integration
- `src/internal/apiv1/api_new.go` - Authentication improvements
- `src/internal/apiv1/api_error.go` - 429 error handling
- `src/internal/web/web.go` - New routes (list, url, qr, edit)
- `src/src/internal/web/data/theme/dark.theme` - Updated to Modern Dark
- `src/src/internal/web/data/theme/light.theme` - Updated to Modern Light
- `.github/workflows/build.yml` - Docker build job, dev releases, Go 1.23
- `Makefile` - Simplified for local development, Go 1.23
- `README.md` - Complete restructure with all new features
- `Dockerfile` - Updated paths and structure, Go 1.23
- `docker-compose.yml` - Complete rewrite with PostgreSQL/MariaDB options
- `.gitattributes` - Export-ignore for clean archives
- `go.mod` - Added dependencies (yaml, mysql driver, crypto)

----
## 📝 Changelog: 2025-11-23 at 11:00:08 📝

📝 Update codebase 📝  
  
  
README.md  


### 📝 End of changes for 202511231100-git 📝  

----  
## 🔧 Changelog: 2025-11-23 at 10:43:42 🔧  

🔧 Update configuration files 🔧  
  
  
.github/workflows/build.yml  


### 🔧 End of changes for 202511231043-git 🔧  

----  
## 🗃️ Changelog: 2025-11-23 at 10:32:02 🗃️  

🗃️ Update codebase 🗃️  
  
  
Makefile  
README.md  


### 🗃️ End of changes for 202511231032-git 🗃️  

----  
## 🔧 Changelog: 2025-11-23 at 10:00:45 🔧  

🔧 Update configuration files 🔧  
  
  
.gitignore  


### 🔧 End of changes for 202511231000-git 🔧  

----  
## 🔧 Changelog: 2025-11-23 at 09:56:13 🔧  

🔧 Update configuration files 🔧  
  
  
cmd/caspaste/  
cmd/lenpaste/lenpaste.go  
Dockerfile  
entrypoint.sh  
.github/workflows/build.yml  
.github/workflows/release_docker.yml  
.github/workflows/release_sources.yml  
go.mod  
go.sum  
internal/apiv1/api_error.go  
internal/apiv1/api_get.go  
internal/apiv1/api.go  
internal/apiv1/api_main.go  
internal/apiv1/api_new.go  
internal/apiv1/api_server.go  
internal/config/config.go  
internal/logger/logger.go  
internal/netshare/netshare_paste.go  
internal/raw/raw_error.go  
internal/raw/raw.go  
internal/raw/raw_raw.go  
internal/storage/storage.go  
src/internal/web/data/about.tmpl  
src/internal/web/data/authors.tmpl  
src/internal/web/data/base.tmpl  
src/internal/web/data/docs_apiv1.tmpl  
src/internal/web/data/docs.tmpl  
src/internal/web/data/locale/en.json  
src/internal/web/data/main.tmpl  
src/internal/web/data/source_code.tmpl  
src/internal/web/data/style.css  
src/internal/web/data/theme/dark.theme  
src/internal/web/data/theme/light.theme  
internal/web/web_dl.go  
internal/web/web_docs.go  
internal/web/web_embedded.go  
internal/web/web_embedded_help.go  
internal/web/web_error.go  
internal/web/web_get.go  
internal/web/web.go  
internal/web/web_new.go  
internal/web/web_other.go  
internal/web/web_settings.go  
internal/web/web_sitemap.go  
Makefile  
README.md  


### 🔧 End of changes for 202511230956-git 🔧  

----  
## 🐳 Changelog: 2025-11-08 at 18:55:20 🐳  

🐳 Update codebase 🐳  
  
  
Dockerfile  


### 🐳 End of changes for 202511081855-git 🐳  

----  
# Changelog
Semantic versioning is used (https://semver.org/).


## v1.3.1
- Fixed a problem with building Lenpaste from source code.
- Revised documentation.
- Minor improvements were made.

## v1.3
- UI: Added custom themes support. Added light theme.
- UI: Added translations into Bengali and German (thanks Pardesi_Cat and Hiajen).
- UI: Check boxes and spoilers now have a custom design.
- Admin: Added support for `X-Real-IP` header for reverse proxy.
- Admin: Added Server response header (for example: `Lenpaste/1.3`).
- Fix: many bugs and errors.
- Dev: Improved quality of `Dockerfile` and `entrypoint.sh`

## v1.2
- UI: Add history tab.
- UI: Add copy to clipboard button.
- Admin: Rate-limits on paste creation (`LENPASTE_NEW_PASTES_PER_5MIN` or `-new-pastes-per-5min`).
- Admin: Add terms of use support (`/data/terms` or `-server-terms`).
- Admin: Add default paste life time for WEB interface (`LENPASTE_UI_DEFAULT_LIFETIME` or `-ui-default-lifetime`).
- Admin: Private servers - password request to create paste (`/data/caspasswd` or `-caspasswd-file`).
- Fix: **Critical security fix!**
- Fix: not saving cookies.
- Fix: display language name in WEB.
- Fix: compatibility with WebKit (Gnome WEB).
- Dev: Drop Go 1.15 support. Update dependencies.


## v1.1.1
- Fixed: Incorrect operation of the maximum paste life parameter.
- Updated README.


## v1.1
- You can now specify author, author email and author URL for paste.
- Full localization into Russian.
- Added settings menu.
- Paste creation and expiration times are now displayed in the user's time zone.
- Add PostgreSQL DB support.


## v1.0
This is the first stable release of Lenpaste🎉

Compared to the previous unstable versions, everything has been drastically improved:
design, loading speed of the pages, API, work with the database.
Plus added syntax highlighting in the web interface.


## v0.2
Features:
- Paste title
- About server information
- Improved documentation
- Logging and log rotation
- Storage configuration
- Code optimization

Bug fixes:
- Added `./version.json` to Docker image
- Added paste expiration check before opening
- Fixed incorrect error of expired pastes
- API errors now return in JSON


## v0.1
Features:
- Alternative to pastebin.com
- Creating expiration pastes
- Web interface
- API
