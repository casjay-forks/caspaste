// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

// Package admin provides the admin panel UI and API per AI.md PART 17
// Admin panel is available at /{admin_path}/ and API at /api/{version}/{admin_path}/
// Route hierarchy:
//   /{admin_path}/ - Dashboard
//   /{admin_path}/profile - Admin's own profile
//   /{admin_path}/preferences - Admin's own preferences
//   /{admin_path}/notifications - Admin's own notifications
//   /{admin_path}/server/* - All server management
package admin

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

// Panel represents the admin panel
type Panel struct {
	basePath   string
	apiPath    string
	apiVersion string
	enabled    bool
	setupDone  bool
	mu         sync.RWMutex
}

// Config holds admin panel configuration
type Config struct {
	// BasePath is the URL path for the admin panel (default: "admin")
	BasePath string
	// APIVersion is the API version prefix (default: "v1")
	APIVersion string
	// Enabled determines if the admin panel is accessible
	Enabled bool
}

// DefaultConfig returns the default admin panel configuration
func DefaultConfig() *Config {
	return &Config{
		BasePath:   "admin",
		APIVersion: "v1",
		Enabled:    true,
	}
}

// ValidAdminRootPaths are the only valid direct children of /{admin_path}/
var ValidAdminRootPaths = map[string]bool{
	"":              true,
	"profile":       true,
	"preferences":   true,
	"notifications": true,
	"server":        true,
}

// ReservedPaths cannot be used as admin_path
var ReservedPaths = []string{
	"api", "static", "assets", "health", "version",
	"metrics", ".well-known", "graphql", "openapi",
}

// New creates a new admin panel
func New(cfg *Config) *Panel {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	return &Panel{
		basePath:   cfg.BasePath,
		apiVersion: cfg.APIVersion,
		apiPath:    "api/" + cfg.APIVersion + "/" + cfg.BasePath,
		enabled:    cfg.Enabled,
	}
}

// IsEnabled returns true if the admin panel is enabled
func (p *Panel) IsEnabled() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.enabled
}

// IsSetupDone returns true if initial setup has been completed
func (p *Panel) IsSetupDone() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.setupDone
}

// SetSetupDone marks the initial setup as complete
func (p *Panel) SetSetupDone(done bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.setupDone = done
}

// BasePath returns the admin panel base path
func (p *Panel) BasePath() string {
	return p.basePath
}

// APIPath returns the admin API base path
func (p *Panel) APIPath() string {
	return p.apiPath
}

// ValidateAdminPath validates a potential admin path
func ValidateAdminPath(path string) error {
	path = strings.ToLower(strings.TrimSpace(path))

	// Check length
	if len(path) < 2 || len(path) > 32 {
		return fmt.Errorf("admin path must be 2-32 characters")
	}

	// Check valid characters
	for _, c := range path {
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-') {
			return fmt.Errorf("admin path can only contain lowercase letters, numbers, and hyphens")
		}
	}

	// Check no leading/trailing hyphens
	if path[0] == '-' || path[len(path)-1] == '-' {
		return fmt.Errorf("admin path cannot start or end with a hyphen")
	}

	// Check reserved paths
	for _, reserved := range ReservedPaths {
		if path == reserved {
			return fmt.Errorf("'%s' is a reserved path", path)
		}
	}

	return nil
}

// ValidateAdminRoute validates that a route follows the admin hierarchy
func ValidateAdminRoute(path string) error {
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return nil
	}

	firstSegment := parts[0]
	if !ValidAdminRootPaths[firstSegment] {
		return fmt.Errorf("invalid admin route: /%s/* - must use /server/* for server management", firstSegment)
	}
	return nil
}

// Handler returns the HTTP handler for the admin panel
func (p *Panel) Handler() http.Handler {
	mux := http.NewServeMux()

	// Admin root - Dashboard
	mux.HandleFunc("/", p.handleDashboard)

	// Admin's own settings (direct children of /{admin_path}/)
	mux.HandleFunc("/profile", p.handleProfile)
	mux.HandleFunc("/preferences", p.handlePreferences)
	mux.HandleFunc("/notifications", p.handleNotifications)

	// Server management routes (all under /server/)
	mux.HandleFunc("/server/", p.handleServerRoot)
	mux.HandleFunc("/server/settings", p.handleServerSettings)
	mux.HandleFunc("/server/ssl", p.handleServerSSL)
	mux.HandleFunc("/server/email", p.handleServerEmail)
	mux.HandleFunc("/server/scheduler", p.handleServerScheduler)
	mux.HandleFunc("/server/logs", p.handleServerLogs)
	mux.HandleFunc("/server/logs/audit", p.handleServerLogsAudit)
	mux.HandleFunc("/server/backup", p.handleServerBackup)
	mux.HandleFunc("/server/updates", p.handleServerUpdates)
	mux.HandleFunc("/server/info", p.handleServerInfo)
	mux.HandleFunc("/server/metrics", p.handleServerMetrics)

	// Network settings
	mux.HandleFunc("/server/network/", p.handleServerNetworkRoot)
	mux.HandleFunc("/server/network/tor", p.handleServerNetworkTor)
	mux.HandleFunc("/server/network/geoip", p.handleServerNetworkGeoIP)

	// Security settings
	mux.HandleFunc("/server/security/", p.handleServerSecurityRoot)
	mux.HandleFunc("/server/security/auth", p.handleServerSecurityAuth)
	mux.HandleFunc("/server/security/tokens", p.handleServerSecurityTokens)
	mux.HandleFunc("/server/security/firewall", p.handleServerSecurityFirewall)

	// User management (if multi-user enabled)
	mux.HandleFunc("/server/users/", p.handleServerUsers)

	return mux
}

// APIHandler returns the HTTP handler for the admin API
func (p *Panel) APIHandler() http.Handler {
	mux := http.NewServeMux()

	// Admin API root
	mux.HandleFunc("/status", p.apiStatus)

	// Admin's own settings API
	mux.HandleFunc("/profile", p.apiProfile)
	mux.HandleFunc("/preferences", p.apiPreferences)

	// Server management API
	mux.HandleFunc("/server/settings", p.apiServerSettings)
	mux.HandleFunc("/server/ssl", p.apiServerSSL)
	mux.HandleFunc("/server/email", p.apiServerEmail)
	mux.HandleFunc("/server/scheduler", p.apiServerScheduler)
	mux.HandleFunc("/server/logs", p.apiServerLogs)
	mux.HandleFunc("/server/backup", p.apiServerBackup)
	mux.HandleFunc("/server/info", p.apiServerInfo)
	mux.HandleFunc("/server/metrics", p.apiServerMetrics)
	mux.HandleFunc("/server/network/geoip", p.apiServerNetworkGeoIP)
	mux.HandleFunc("/server/network/tor", p.apiServerNetworkTor)
	mux.HandleFunc("/server/security/tokens", p.apiServerSecurityTokens)
	mux.HandleFunc("/server/users", p.apiServerUsers)

	return mux
}

// Dashboard handler
func (p *Panel) handleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "" {
		http.NotFound(w, r)
		return
	}
	p.renderPage(w, "Dashboard", p.dashboardContent())
}

// Admin's own profile
func (p *Panel) handleProfile(w http.ResponseWriter, r *http.Request) {
	p.renderPage(w, "Profile", p.profileContent())
}

// Admin's own preferences
func (p *Panel) handlePreferences(w http.ResponseWriter, r *http.Request) {
	p.renderPage(w, "Preferences", p.preferencesContent())
}

// Admin's own notifications
func (p *Panel) handleNotifications(w http.ResponseWriter, r *http.Request) {
	p.renderPage(w, "Notifications", p.notificationsContent())
}

// Server management handlers

func (p *Panel) handleServerRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/server/" || r.URL.Path == "/server" {
		http.Redirect(w, r, "/"+p.basePath+"/server/settings", http.StatusSeeOther)
		return
	}
	http.NotFound(w, r)
}

func (p *Panel) handleServerSettings(w http.ResponseWriter, r *http.Request) {
	p.renderPage(w, "Server Settings", p.serverSettingsContent())
}

func (p *Panel) handleServerSSL(w http.ResponseWriter, r *http.Request) {
	p.renderPage(w, "SSL/TLS", p.serverSSLContent())
}

func (p *Panel) handleServerEmail(w http.ResponseWriter, r *http.Request) {
	p.renderPage(w, "Email Configuration", p.serverEmailContent())
}

func (p *Panel) handleServerScheduler(w http.ResponseWriter, r *http.Request) {
	p.renderPage(w, "Scheduled Tasks", p.serverSchedulerContent())
}

func (p *Panel) handleServerLogs(w http.ResponseWriter, r *http.Request) {
	p.renderPage(w, "Server Logs", p.serverLogsContent())
}

func (p *Panel) handleServerLogsAudit(w http.ResponseWriter, r *http.Request) {
	p.renderPage(w, "Audit Logs", p.serverLogsAuditContent())
}

func (p *Panel) handleServerBackup(w http.ResponseWriter, r *http.Request) {
	p.renderPage(w, "Backup & Restore", p.serverBackupContent())
}

func (p *Panel) handleServerUpdates(w http.ResponseWriter, r *http.Request) {
	p.renderPage(w, "Updates", p.serverUpdatesContent())
}

func (p *Panel) handleServerInfo(w http.ResponseWriter, r *http.Request) {
	p.renderPage(w, "Server Information", p.serverInfoContent())
}

func (p *Panel) handleServerMetrics(w http.ResponseWriter, r *http.Request) {
	p.renderPage(w, "Metrics Dashboard", p.serverMetricsContent())
}

func (p *Panel) handleServerNetworkRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/server/network/" || r.URL.Path == "/server/network" {
		http.Redirect(w, r, "/"+p.basePath+"/server/network/geoip", http.StatusSeeOther)
		return
	}
	http.NotFound(w, r)
}

func (p *Panel) handleServerNetworkTor(w http.ResponseWriter, r *http.Request) {
	p.renderPage(w, "Tor Configuration", p.serverNetworkTorContent())
}

func (p *Panel) handleServerNetworkGeoIP(w http.ResponseWriter, r *http.Request) {
	p.renderPage(w, "GeoIP Settings", p.serverNetworkGeoIPContent())
}

func (p *Panel) handleServerSecurityRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/server/security/" || r.URL.Path == "/server/security" {
		http.Redirect(w, r, "/"+p.basePath+"/server/security/auth", http.StatusSeeOther)
		return
	}
	http.NotFound(w, r)
}

func (p *Panel) handleServerSecurityAuth(w http.ResponseWriter, r *http.Request) {
	p.renderPage(w, "Authentication", p.serverSecurityAuthContent())
}

func (p *Panel) handleServerSecurityTokens(w http.ResponseWriter, r *http.Request) {
	p.renderPage(w, "API Tokens", p.serverSecurityTokensContent())
}

func (p *Panel) handleServerSecurityFirewall(w http.ResponseWriter, r *http.Request) {
	p.renderPage(w, "Firewall Rules", p.serverSecurityFirewallContent())
}

func (p *Panel) handleServerUsers(w http.ResponseWriter, r *http.Request) {
	p.renderPage(w, "User Management", p.serverUsersContent())
}

// renderPage renders an admin page with the common layout
func (p *Panel) renderPage(w http.ResponseWriter, title, content string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en" data-theme="dark">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="robots" content="noindex, nofollow">
    <title>%s - CasPaste Admin</title>
    <style>
        :root {
            --bg-primary: #1a1a2e;
            --bg-secondary: #16213e;
            --bg-tertiary: #0f3460;
            --text-primary: #eaeaea;
            --text-secondary: #b8b8b8;
            --accent: #e94560;
            --success: #4ade80;
            --warning: #fbbf24;
            --error: #ef4444;
            --border: #2d3748;
        }
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: var(--bg-primary);
            color: var(--text-primary);
            min-height: 100vh;
        }
        .admin-layout {
            display: flex;
            min-height: 100vh;
        }
        .sidebar {
            width: 240px;
            background: var(--bg-secondary);
            border-right: 1px solid var(--border);
            padding: 1rem 0;
            flex-shrink: 0;
        }
        .sidebar-header {
            padding: 0.5rem 1rem 1rem;
            border-bottom: 1px solid var(--border);
            margin-bottom: 1rem;
        }
        .sidebar-header h1 {
            font-size: 1.25rem;
            color: var(--accent);
        }
        .sidebar-nav { list-style: none; }
        .sidebar-nav li a {
            display: flex;
            align-items: center;
            padding: 0.75rem 1rem;
            color: var(--text-secondary);
            text-decoration: none;
            transition: background 0.2s, color 0.2s;
        }
        .sidebar-nav li a:hover {
            background: var(--bg-tertiary);
            color: var(--text-primary);
        }
        .sidebar-nav li a.active {
            background: var(--bg-tertiary);
            color: var(--accent);
            border-left: 3px solid var(--accent);
        }
        .sidebar-section {
            margin-top: 1rem;
            padding-top: 1rem;
            border-top: 1px solid var(--border);
        }
        .sidebar-section-title {
            padding: 0.5rem 1rem;
            font-size: 0.75rem;
            text-transform: uppercase;
            color: var(--text-secondary);
            letter-spacing: 0.05em;
        }
        .main-content {
            flex: 1;
            display: flex;
            flex-direction: column;
        }
        .header {
            height: 60px;
            background: var(--bg-secondary);
            border-bottom: 1px solid var(--border);
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding: 0 1.5rem;
        }
        .header-left {
            display: flex;
            align-items: center;
            gap: 1rem;
        }
        .breadcrumb {
            display: flex;
            gap: 0.5rem;
            color: var(--text-secondary);
        }
        .breadcrumb a { color: var(--text-secondary); text-decoration: none; }
        .breadcrumb a:hover { color: var(--text-primary); }
        .header-right {
            display: flex;
            align-items: center;
            gap: 1rem;
        }
        .status-indicator {
            width: 10px;
            height: 10px;
            border-radius: 50%%;
            background: var(--success);
        }
        .user-menu {
            display: flex;
            align-items: center;
            gap: 0.5rem;
            color: var(--text-secondary);
        }
        .btn {
            padding: 0.5rem 1rem;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-size: 0.875rem;
            transition: background 0.2s;
        }
        .btn-primary {
            background: var(--accent);
            color: white;
        }
        .btn-primary:hover { background: #d63d55; }
        .btn-secondary {
            background: var(--bg-tertiary);
            color: var(--text-primary);
        }
        .page-content {
            flex: 1;
            padding: 1.5rem;
        }
        .page-title {
            font-size: 1.5rem;
            margin-bottom: 1.5rem;
        }
        .card {
            background: var(--bg-secondary);
            border: 1px solid var(--border);
            border-radius: 8px;
            padding: 1.5rem;
            margin-bottom: 1rem;
        }
        .card-title {
            font-size: 1rem;
            margin-bottom: 1rem;
            color: var(--text-secondary);
        }
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 1rem;
        }
        .stat-card {
            background: var(--bg-secondary);
            border: 1px solid var(--border);
            border-radius: 8px;
            padding: 1.5rem;
            text-align: center;
        }
        .stat-value {
            font-size: 2rem;
            font-weight: bold;
            color: var(--accent);
        }
        .stat-label {
            color: var(--text-secondary);
            font-size: 0.875rem;
        }
        .footer {
            height: 40px;
            background: var(--bg-secondary);
            border-top: 1px solid var(--border);
            display: flex;
            align-items: center;
            justify-content: center;
            color: var(--text-secondary);
            font-size: 0.75rem;
            gap: 1rem;
        }
        .footer a { color: var(--text-secondary); text-decoration: none; }
        .footer a:hover { color: var(--text-primary); }
        @media (max-width: 768px) {
            .sidebar { display: none; }
            .admin-layout { flex-direction: column; }
        }
    </style>
</head>
<body>
    <div class="admin-layout">
        <nav class="sidebar">
            <div class="sidebar-header">
                <h1>CasPaste</h1>
            </div>
            <ul class="sidebar-nav">
                <li><a href="/%s/">Dashboard</a></li>
            </ul>
            <div class="sidebar-section">
                <div class="sidebar-section-title">Account</div>
                <ul class="sidebar-nav">
                    <li><a href="/%s/profile">Profile</a></li>
                    <li><a href="/%s/preferences">Preferences</a></li>
                    <li><a href="/%s/notifications">Notifications</a></li>
                </ul>
            </div>
            <div class="sidebar-section">
                <div class="sidebar-section-title">Server</div>
                <ul class="sidebar-nav">
                    <li><a href="/%s/server/settings">Settings</a></li>
                    <li><a href="/%s/server/ssl">SSL/TLS</a></li>
                    <li><a href="/%s/server/email">Email</a></li>
                    <li><a href="/%s/server/scheduler">Scheduler</a></li>
                    <li><a href="/%s/server/logs">Logs</a></li>
                    <li><a href="/%s/server/backup">Backup</a></li>
                    <li><a href="/%s/server/info">Info</a></li>
                    <li><a href="/%s/server/metrics">Metrics</a></li>
                </ul>
            </div>
            <div class="sidebar-section">
                <div class="sidebar-section-title">Network</div>
                <ul class="sidebar-nav">
                    <li><a href="/%s/server/network/geoip">GeoIP</a></li>
                    <li><a href="/%s/server/network/tor">Tor</a></li>
                </ul>
            </div>
            <div class="sidebar-section">
                <div class="sidebar-section-title">Security</div>
                <ul class="sidebar-nav">
                    <li><a href="/%s/server/security/auth">Authentication</a></li>
                    <li><a href="/%s/server/security/tokens">API Tokens</a></li>
                    <li><a href="/%s/server/security/firewall">Firewall</a></li>
                </ul>
            </div>
            <div class="sidebar-section">
                <div class="sidebar-section-title">Users</div>
                <ul class="sidebar-nav">
                    <li><a href="/%s/server/users/">Manage Users</a></li>
                </ul>
            </div>
        </nav>
        <div class="main-content">
            <header class="header">
                <div class="header-left">
                    <div class="breadcrumb">
                        <a href="/%s/">Admin</a>
                        <span>/</span>
                        <span>%s</span>
                    </div>
                </div>
                <div class="header-right">
                    <div class="status-indicator" title="Server OK"></div>
                    <div class="user-menu">
                        <span>Admin</span>
                    </div>
                    <button class="btn btn-secondary">Logout</button>
                </div>
            </header>
            <main class="page-content">
                <h1 class="page-title">%s</h1>
                %s
            </main>
            <footer class="footer">
                <span>CasPaste v1.0.0</span>
                <a href="/docs">Documentation</a>
                <span>Status: Running</span>
            </footer>
        </div>
    </div>
</body>
</html>`,
		title,
		p.basePath, p.basePath, p.basePath, p.basePath,
		p.basePath, p.basePath, p.basePath, p.basePath, p.basePath, p.basePath, p.basePath, p.basePath,
		p.basePath, p.basePath,
		p.basePath, p.basePath, p.basePath,
		p.basePath,
		p.basePath, title, title, content)
	w.Write([]byte(html))
}

// Content generators for each page

func (p *Panel) dashboardContent() string {
	return `<div class="stats-grid">
    <div class="stat-card">
        <div class="stat-value">0</div>
        <div class="stat-label">Total Pastes</div>
    </div>
    <div class="stat-card">
        <div class="stat-value">0</div>
        <div class="stat-label">Active Users</div>
    </div>
    <div class="stat-card">
        <div class="stat-value">0 MB</div>
        <div class="stat-label">Storage Used</div>
    </div>
    <div class="stat-card">
        <div class="stat-value">0</div>
        <div class="stat-label">Requests Today</div>
    </div>
</div>
<div class="card" style="margin-top: 1.5rem;">
    <div class="card-title">System Status</div>
    <p>Server is running normally.</p>
</div>`
}

func (p *Panel) profileContent() string {
	return `<div class="card">
    <div class="card-title">Admin Profile</div>
    <p>Manage your admin account settings.</p>
</div>`
}

func (p *Panel) preferencesContent() string {
	return `<div class="card">
    <div class="card-title">Preferences</div>
    <p>Configure your personal admin panel preferences.</p>
</div>`
}

func (p *Panel) notificationsContent() string {
	return `<div class="card">
    <div class="card-title">Notifications</div>
    <p>View and manage your notifications.</p>
</div>`
}

func (p *Panel) serverSettingsContent() string {
	return `<div class="card">
    <div class="card-title">General Settings</div>
    <p>Configure server settings.</p>
</div>`
}

func (p *Panel) serverSSLContent() string {
	return `<div class="card">
    <div class="card-title">SSL/TLS Configuration</div>
    <p>Manage SSL certificates and HTTPS settings.</p>
</div>`
}

func (p *Panel) serverEmailContent() string {
	return `<div class="card">
    <div class="card-title">Email Configuration</div>
    <p>Configure SMTP settings for email notifications.</p>
</div>`
}

func (p *Panel) serverSchedulerContent() string {
	return `<div class="card">
    <div class="card-title">Scheduled Tasks</div>
    <p>View and manage scheduled background tasks.</p>
</div>`
}

func (p *Panel) serverLogsContent() string {
	return `<div class="card">
    <div class="card-title">Server Logs</div>
    <p>View server logs and activity.</p>
</div>`
}

func (p *Panel) serverLogsAuditContent() string {
	return `<div class="card">
    <div class="card-title">Audit Logs</div>
    <p>View security audit logs.</p>
</div>`
}

func (p *Panel) serverBackupContent() string {
	return `<div class="card">
    <div class="card-title">Backup & Restore</div>
    <p>Create backups and restore from previous backups.</p>
</div>`
}

func (p *Panel) serverUpdatesContent() string {
	return `<div class="card">
    <div class="card-title">Updates</div>
    <p>Check for and apply updates.</p>
</div>`
}

func (p *Panel) serverInfoContent() string {
	return `<div class="card">
    <div class="card-title">Server Information</div>
    <p>View server details and system information.</p>
</div>`
}

func (p *Panel) serverMetricsContent() string {
	return `<div class="card">
    <div class="card-title">Metrics Dashboard</div>
    <p>View server metrics and performance data.</p>
</div>`
}

func (p *Panel) serverNetworkTorContent() string {
	return `<div class="card">
    <div class="card-title">Tor Configuration</div>
    <p>Configure Tor hidden service settings.</p>
</div>`
}

func (p *Panel) serverNetworkGeoIPContent() string {
	return `<div class="card">
    <div class="card-title">GeoIP Settings</div>
    <p>Configure GeoIP blocking and location detection.</p>
</div>`
}

func (p *Panel) serverSecurityAuthContent() string {
	return `<div class="card">
    <div class="card-title">Authentication Settings</div>
    <p>Configure authentication and login security.</p>
</div>`
}

func (p *Panel) serverSecurityTokensContent() string {
	return `<div class="card">
    <div class="card-title">API Tokens</div>
    <p>Manage API tokens and access keys.</p>
</div>`
}

func (p *Panel) serverSecurityFirewallContent() string {
	return `<div class="card">
    <div class="card-title">Firewall Rules</div>
    <p>Configure IP blocking and firewall rules.</p>
</div>`
}

func (p *Panel) serverUsersContent() string {
	return `<div class="card">
    <div class="card-title">User Management</div>
    <p>Manage user accounts (if multi-user mode enabled).</p>
</div>`
}

// API Handlers

func (p *Panel) apiStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok": true, "data": {"status": "running"}}` + "\n"))
}

func (p *Panel) apiProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok": true, "data": {"username": "admin"}}` + "\n"))
}

func (p *Panel) apiPreferences(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok": true, "data": {"theme": "dark", "language": "en"}}` + "\n"))
}

func (p *Panel) apiServerSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok": true, "data": {}}` + "\n"))
}

func (p *Panel) apiServerSSL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok": true, "data": {"enabled": false}}` + "\n"))
}

func (p *Panel) apiServerEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok": true, "data": {"enabled": false}}` + "\n"))
}

func (p *Panel) apiServerScheduler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok": true, "data": {"tasks": []}}` + "\n"))
}

func (p *Panel) apiServerLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok": true, "data": {"logs": []}}` + "\n"))
}

func (p *Panel) apiServerBackup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok": true, "data": {"backups": []}}` + "\n"))
}

func (p *Panel) apiServerInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok": true, "data": {"version": "1.0.0"}}` + "\n"))
}

func (p *Panel) apiServerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok": true, "data": {"uptime": 0}}` + "\n"))
}

func (p *Panel) apiServerNetworkGeoIP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok": true, "data": {"enabled": false}}` + "\n"))
}

func (p *Panel) apiServerNetworkTor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok": true, "data": {"enabled": false}}` + "\n"))
}

func (p *Panel) apiServerSecurityTokens(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok": true, "data": {"tokens": []}}` + "\n"))
}

func (p *Panel) apiServerUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok": true, "data": {"users": []}}` + "\n"))
}
