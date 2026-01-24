# Admin Panel

CasPaste includes a web-based admin panel for managing the server.

## Accessing the Admin Panel

The admin panel is available at `/admin/` when running in private mode (`server.public: false`).

## First-Time Setup

When starting CasPaste in private mode for the first time, admin credentials are auto-generated:

```
╔════════════════════════════════════════════════════════════╗
║  CasPaste                                                  ║
╠════════════════════════════════════════════════════════════╣
║  Mode:        Private (authentication required)            ║
║  Username:    admin                                        ║
║  Password:    eoYBn7I9Z&ZHGqCY                             ║
║  SAVE THESE CREDENTIALS - shown only once!                 ║
╚════════════════════════════════════════════════════════════╝
```

**Important:** Save these credentials - they are only shown once!

## Admin Features

### Dashboard

- Server status overview
- Recent paste activity
- Storage usage
- Quick actions

### Server Settings

Access via `/admin/server/settings`

- Server title and description
- Public/private mode toggle
- FQDN configuration
- Proxy settings
- Timeout configuration

### Database Management

- View database statistics
- Run cleanup operations
- Export/import data

### Backup & Restore

Access via `/admin/server/backup`

- Create manual backups
- View backup history
- Restore from backup
- Configure automatic backups

### User Management

When running in private mode:

- View logged-in sessions
- Force logout users
- Change admin password

## Security

### Brute Force Protection

- 5 failed login attempts triggers 15-minute lockout
- Lockout applies to IP address
- Attempts are logged for security review

### Session Security

- Sessions expire after 24 hours of inactivity
- HttpOnly cookies (not accessible to JavaScript)
- SameSite=Strict (CSRF protection)
- Secure flag when HTTPS detected

### Audit Logging

All admin actions are logged:

- Login attempts (success/failure)
- Configuration changes
- Backup/restore operations
- User management actions

Logs are stored in `{logs_dir}/audit.log`

## CLI Administration

Many admin tasks can also be performed via CLI:

```bash
# Check status
caspaste --status

# Create backup
caspaste --maintenance backup

# Restore backup
caspaste --maintenance restore

# Run cleanup
caspaste --maintenance cleanup
```

## Troubleshooting

### Locked Out

If you're locked out due to failed login attempts:

1. Wait 15 minutes for lockout to expire, OR
2. Restart the server to clear lockouts, OR
3. Delete the lockout data from the database

### Forgot Password

Reset admin password:

```bash
# Stop the server
caspaste --service stop

# Reset password (creates new credentials)
caspaste --maintenance reset-admin

# Start the server
caspaste --service start
```

### Session Issues

Clear all sessions:

```bash
caspaste --maintenance clear-sessions
```
