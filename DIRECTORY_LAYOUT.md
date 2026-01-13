# CasPaste Directory Layouts by Platform

This document defines the standard directory layouts for CasPaste across different platforms and deployment methods.

## Docker

```
/config/caspaste/             # Config directory
├── caspaste.yml              # Main configuration file

/data/caspaste/               # Data directory
├── backups/                  # Database backups
├── .db-state                 # Database state tracking
└── .maintenance              # Maintenance mode flag

/data/db/sqlite/              # Database directory
└── caspaste.db               # SQLite database file

/cache/                       # Cache directory
└── (auto-managed)

/logs/                        # Logs directory
└── (auto-managed)
```

**Docker Volume Mounts:**
```yaml
volumes:
  - ./rootfs/config/caspaste:/config/caspaste
  - ./rootfs/data/caspaste:/data/caspaste
  - ./rootfs/data/db/sqlite:/data/db/sqlite
  - ./rootfs/data/backups:/data/backups
  - ./rootfs/cache:/cache
  - ./rootfs/logs:/logs
```

---

## Linux (System/Root Install)

```
/etc/caspaste/                # Config directory
└── caspaste.yml              # Main configuration file

/var/lib/caspaste/            # Data directory
├── backups/                  # Database backups
├── .db-state                 # Database state tracking
└── .maintenance              # Maintenance mode flag

/var/lib/caspaste/db/         # Database directory
└── caspaste.db               # SQLite database file

/var/cache/caspaste/          # Cache directory
└── (auto-managed)

/var/log/caspaste/            # Logs directory
└── (auto-managed)

/mnt/Backups/caspaste/        # Off-site backups
└── backup-YYYYMMDD-HHMMSS.tar.gz
```

**Startup Command:**
```bash
sudo caspaste --config /etc/caspaste --data /var/lib/caspaste
```

---

## Linux (User Install)

```
~/.config/caspaste/           # Config directory
└── caspaste.yml              # Main configuration file

~/.local/share/caspaste/      # Data directory
├── backups/                  # Database backups
├── .db-state                 # Database state tracking
└── .maintenance              # Maintenance mode flag

~/.local/share/caspaste/db/   # Database directory
└── caspaste.db               # SQLite database file

~/.cache/caspaste/            # Cache directory
└── (auto-managed)

~/.local/log/caspaste/        # Logs directory
└── (auto-managed)
```

**Startup Command:**
```bash
caspaste --config ~/.config/caspaste --data ~/.local/share/caspaste
```

---

## macOS (System/Root Install)

```
/etc/caspaste/                # Config directory
└── caspaste.yml              # Main configuration file

/var/lib/caspaste/            # Data directory
├── backups/                  # Database backups
├── .db-state                 # Database state tracking
└── .maintenance              # Maintenance mode flag

/var/lib/caspaste/db/         # Database directory
└── caspaste.db               # SQLite database file

/var/cache/caspaste/          # Cache directory
└── (auto-managed)

/var/log/caspaste/            # Logs directory
└── (auto-managed)

/var/backups/caspaste/        # Off-site backups
└── backup-YYYYMMDD-HHMMSS.tar.gz
```

**Startup Command:**
```bash
sudo caspaste --config /etc/caspaste --data /var/lib/caspaste
```

---

## macOS (User Install)

```
~/Library/Application Support/CasPaste/config/   # Config directory
└── caspaste.yml                                 # Main configuration file

~/Library/Application Support/CasPaste/          # Data directory
├── backups/                                     # Database backups
├── .db-state                                    # Database state tracking
└── .maintenance                                 # Maintenance mode flag

~/Library/Application Support/CasPaste/db/       # Database directory
└── caspaste.db                                  # SQLite database file

~/Library/Caches/CasPaste/                       # Cache directory
└── (auto-managed)

~/Library/Logs/CasPaste/                         # Logs directory
└── (auto-managed)

~/Library/Application Support/CasPaste/Backups/  # Off-site backups
└── backup-YYYYMMDD-HHMMSS.tar.gz
```

**Startup Command:**
```bash
caspaste --config ~/Library/Application\ Support/CasPaste/config \
         --data ~/Library/Application\ Support/CasPaste
```

---

## Windows (Administrator/System Install)

```
C:\ProgramData\CasPaste\config\     # Config directory
└── caspaste.yml                    # Main configuration file

C:\ProgramData\CasPaste\data\       # Data directory
├── backups\                        # Database backups
├── .db-state                       # Database state tracking
└── .maintenance                    # Maintenance mode flag

C:\ProgramData\CasPaste\data\db\    # Database directory
└── caspaste.db                     # SQLite database file

C:\ProgramData\CasPaste\Cache\      # Cache directory
└── (auto-managed)

C:\ProgramData\CasPaste\Logs\       # Logs directory
└── (auto-managed)

C:\ProgramData\CasPaste\Backups\    # Off-site backups
└── backup-YYYYMMDD-HHMMSS.tar.gz
```

**Startup Command (PowerShell as Administrator):**
```powershell
caspaste.exe --config "C:\ProgramData\CasPaste\config" `
             --data "C:\ProgramData\CasPaste\data"
```

---

## Windows (User Install)

```
%APPDATA%\CasPaste\config\          # Config directory
└── caspaste.yml                    # Main configuration file

%APPDATA%\CasPaste\data\            # Data directory
├── backups\                        # Database backups
├── .db-state                       # Database state tracking
└── .maintenance                    # Maintenance mode flag

%APPDATA%\CasPaste\data\db\         # Database directory
└── caspaste.db                     # SQLite database file

%LOCALAPPDATA%\CasPaste\Cache\      # Cache directory
└── (auto-managed)

%LOCALAPPDATA%\CasPaste\Logs\       # Logs directory
└── (auto-managed)

%APPDATA%\CasPaste\Backups\         # Off-site backups
└── backup-YYYYMMDD-HHMMSS.tar.gz
```

**Typical Paths (expands to):**
- `%APPDATA%` = `C:\Users\<username>\AppData\Roaming`
- `%LOCALAPPDATA%` = `C:\Users\<username>\AppData\Local`

**Startup Command:**
```powershell
caspaste.exe --config "%APPDATA%\CasPaste\config" `
             --data "%APPDATA%\CasPaste\data"
```

---

## FreeBSD/OpenBSD (System/Root Install)

```
/etc/caspaste/                # Config directory
└── caspaste.yml              # Main configuration file

/var/lib/caspaste/            # Data directory
├── backups/                  # Database backups
├── .db-state                 # Database state tracking
└── .maintenance              # Maintenance mode flag

/var/lib/caspaste/db/         # Database directory
└── caspaste.db               # SQLite database file

/var/cache/caspaste/          # Cache directory
└── (auto-managed)

/var/log/caspaste/            # Logs directory
└── (auto-managed)

/var/backups/caspaste/        # Off-site backups
└── backup-YYYYMMDD-HHMMSS.tar.gz
```

**Startup Command:**
```bash
doas caspaste --config /etc/caspaste --data /var/lib/caspaste
```

---

## FreeBSD/OpenBSD (User Install)

```
~/.config/caspaste/           # Config directory
└── caspaste.yml              # Main configuration file

~/.local/share/caspaste/      # Data directory
├── backups/                  # Database backups
├── .db-state                 # Database state tracking
└── .maintenance              # Maintenance mode flag

~/.local/share/caspaste/db/   # Database directory
└── caspaste.db               # SQLite database file

~/.cache/caspaste/            # Cache directory
└── (auto-managed)

~/.local/log/caspaste/        # Logs directory
└── (auto-managed)

~/.caspaste/backups/          # Off-site backups
└── backup-YYYYMMDD-HHMMSS.tar.gz
```

**Startup Command:**
```bash
caspaste --config ~/.config/caspaste --data ~/.local/share/caspaste
```

---

## Auto-Detection Logic

CasPaste automatically detects the platform and whether it's running as root/admin:

1. **Platform Detection:** Uses `runtime.GOOS` (linux, darwin, windows, freebsd, openbsd)
2. **Privilege Detection:** Checks if running as root/admin (UID 0 or admin rights)
3. **Directory Selection:** Based on platform + privilege level

**If directories not specified via CLI flags:**
- Uses system directories if root/admin
- Uses user directories if non-root
- Docker detection: if `--data /data/caspaste`, uses Docker paths

---

## Environment Variable Overrides

All directory paths can be overridden via environment variables:

```bash
CASPASTE_CONFIG_DIR=/custom/config
CASPASTE_DATA_DIR=/custom/data
CASPASTE_DB_DIR=/custom/db
CASPASTE_CACHE_DIR=/custom/cache
CASPASTE_LOGS_DIR=/custom/logs
CASPASTE_BACKUP_DIR=/custom/backups
```

---

## Migration Notes

### From LenPaste

If migrating from LenPaste, your data structure may differ. CasPaste maintains backward compatibility:

**Environment Variables:**
- `LENPASTE_*` variables still work
- `CASPASTE_*` takes precedence

**Directory Paths:**
- Old data can stay in place
- Update config file paths to match your setup
- Or use CLI flags to specify custom paths

---

## Summary Table

| Platform | Type | Config | Data | Database |
|----------|------|--------|------|----------|
| Docker | - | /config/caspaste | /data/caspaste | /data/db/sqlite |
| Linux | Root | /etc/caspaste | /var/lib/caspaste | /var/lib/caspaste/db |
| Linux | User | ~/.config/caspaste | ~/.local/share/caspaste | ~/.local/share/caspaste/db |
| macOS | Root | /etc/caspaste | /var/lib/caspaste | /var/lib/caspaste/db |
| macOS | User | ~/Library/Application Support/CasPaste/config | ~/Library/Application Support/CasPaste | ~/Library/Application Support/CasPaste/db |
| Windows | Admin | C:\ProgramData\CasPaste\config | C:\ProgramData\CasPaste\data | C:\ProgramData\CasPaste\data\db |
| Windows | User | %APPDATA%\CasPaste\config | %APPDATA%\CasPaste\data | %APPDATA%\CasPaste\data\db |
| BSD | Root | /etc/caspaste | /var/lib/caspaste | /var/lib/caspaste/db |
| BSD | User | ~/.config/caspaste | ~/.local/share/caspaste | ~/.local/share/caspaste/db |

---

This layout follows FHS (Filesystem Hierarchy Standard) on Unix-like systems and platform conventions on Windows/macOS.
