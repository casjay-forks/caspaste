
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// YAMLConfig represents the YAML configuration file structure
// All configuration is organized into logical top-level sections
type YAMLConfig struct {
	Server struct {
		Public      bool   `yaml:"public"`      // Public instance (default: true = no auth, false = auth required)
		FQDN        string `yaml:"fqdn"`        // Public FQDN for building URLs (empty=auto-detect from headers/hostname, set to override)
		Listen      string `yaml:"listen"`      // Listen address (all, ::, 0.0.0.0, specific IP)
		Port        string `yaml:"port"`        // Port number (empty=auto-detect available port)
		Title       string `yaml:"title"`       // Server title
		TagLine     string `yaml:"tagline"`     // Server tagline (short description)
		Description string `yaml:"description"` // Server description (longer description for meta tags)

		Proxy struct {
			Allowed []string `yaml:"allowed"` // Additional trusted proxy IPs/CIDRs (appended to default private ranges)
		} `yaml:"proxy"`
		
		Administrator struct {
			Name  string `yaml:"name"`  // Admin name
			Email string `yaml:"email"` // Admin email
			From  string `yaml:"from"`  // Email from address
		} `yaml:"administrator"`
		
		Timeouts struct {
			Read  int `yaml:"read"`  // Read timeout in seconds (default: 15)
			Write int `yaml:"write"` // Write timeout in seconds (default: 15)
			Idle  int `yaml:"idle"`  // Idle timeout in seconds (default: 60)
		} `yaml:"timeouts"`
	} `yaml:"server"`

	Database struct {
		Driver        string `yaml:"driver"`         // sqlite, postgres, mysql
		Source        string `yaml:"source"`         // Connection string
		MaxOpenConns  int    `yaml:"max_open_conns"` // Max open connections
		MaxIdleConns  int    `yaml:"max_idle_conns"` // Max idle connections
		CleanupPeriod string `yaml:"cleanup_period"` // Cleanup interval (e.g. "1m", "5m")
	} `yaml:"database"`

	Security struct {
		PasswordFile string `yaml:"password_file"` // Path to password file (auto-generated when server.public=false)
		
		Headers struct {
			XFrameOptions           string `yaml:"x_frame_options"`            // X-Frame-Options header
			XContentTypeOptions     string `yaml:"x_content_type_options"`     // X-Content-Type-Options header
			ContentSecurityPolicy   string `yaml:"content_security_policy"`    // Content-Security-Policy header
			ReferrerPolicy          string `yaml:"referrer_policy"`            // Referrer-Policy header
			PermissionsPolicy       string `yaml:"permissions_policy"`         // Permissions-Policy header
			StrictTransportSecurity string `yaml:"strict_transport_security"`  // Strict-Transport-Security header
		} `yaml:"headers"`
		
		TLS struct {
			MinVersion   string   `yaml:"min_version"`   // Minimum TLS version: 1.0, 1.1, 1.2, 1.3
			CipherSuites []string `yaml:"cipher_suites"` // Allowed cipher suites
			CertFile     string   `yaml:"cert_file"`     // TLS certificate file path (optional, auto-detected)
			KeyFile      string   `yaml:"key_file"`      // TLS key file path (optional, auto-detected)
		} `yaml:"tls"`
		
		Upload struct {
			MaxFileSize int64    `yaml:"max_file_size"`      // Max upload size in bytes
			AllowedMIME []string `yaml:"allowed_mime_types"` // Allowed MIME types
		} `yaml:"upload"`
		
		CORS struct {
			Enabled        bool     `yaml:"enabled"`          // Enable CORS
			AllowedOrigins []string `yaml:"allowed_origins"`  // Allowed origins (* for all)
			AllowedMethods []string `yaml:"allowed_methods"`  // Allowed HTTP methods
			AllowedHeaders []string `yaml:"allowed_headers"`  // Allowed headers
			MaxAge         int      `yaml:"max_age"`          // Preflight cache duration in seconds
		} `yaml:"cors"`
	} `yaml:"security"`

	Web struct {
		UI struct {
			DefaultLifetime string `yaml:"default_lifetime"` // Default paste lifetime
			DefaultTheme    string `yaml:"default_theme"`    // Default theme (e.g. "dracula")
			ThemesDir       string `yaml:"themes_dir"`       // Themes directory (default: {data_dir}/web/themes)
		} `yaml:"ui"`

		Content struct {
			About    string `yaml:"about"`    // Path to custom about page (empty=auto-generated, relative to {data_dir}/web/docs)
			Rules    string `yaml:"rules"`    // Path to custom rules page (empty=auto-generated)
			Terms    string `yaml:"terms"`    // Path to custom terms page (empty=auto-generated)
			Security string `yaml:"security"` // Path to custom security.txt (empty=auto-generated)
		} `yaml:"content"`

		Branding struct {
			Logo    string `yaml:"logo"`    // Logo path or URL (e.g. "/static/logo.png" or "https://example.com/logo.png")
			Favicon string `yaml:"favicon"` // Favicon path or URL (e.g. "/static/favicon.ico" or "https://example.com/favicon.ico")
		} `yaml:"branding"`
		
		Security struct {
			Contact struct {
				Email string `yaml:"email"` // Security contact email
				Name  string `yaml:"name"`  // Security contact name
			} `yaml:"contact"`
		} `yaml:"security"`
		
		SEO struct {
			Robots struct {
				Allow  string `yaml:"allow"` // Paths to allow in robots.txt
				Deny   string `yaml:"deny"`  // Paths to deny in robots.txt
				Agents struct {
					Deny []string `yaml:"deny"` // User agents to deny
				} `yaml:"agents"`
			} `yaml:"robots"`
		} `yaml:"seo"`
	} `yaml:"web"`

	Limits struct {
		TitleMaxLength    int    `yaml:"title_max_length"`   // Max title length
		BodyMaxLength     int    `yaml:"body_max_length"`    // Max paste body length
		MaxPasteLifetime  string `yaml:"max_paste_lifetime"` // Max paste lifetime (e.g. "30d", "never")
		
		RateLimit struct {
			GetPastes struct {
				Per5Min  uint `yaml:"per_5min"`  // GET requests per 5 minutes
				Per15Min uint `yaml:"per_15min"` // GET requests per 15 minutes
				Per1Hour uint `yaml:"per_1hour"` // GET requests per 1 hour
			} `yaml:"get_pastes"`
			
			NewPastes struct {
				Per5Min  uint `yaml:"per_5min"`  // POST requests per 5 minutes
				Per15Min uint `yaml:"per_15min"` // POST requests per 15 minutes
				Per1Hour uint `yaml:"per_1hour"` // POST requests per 1 hour
			} `yaml:"new_pastes"`
		} `yaml:"rate_limit"`
	} `yaml:"limits"`

	Directories struct {
		Data   string `yaml:"data"`   // Data directory
		Config string `yaml:"config"` // Config directory
		Db     string `yaml:"db"`     // Database directory
		Cache  string `yaml:"cache"`  // Cache directory
		Logs   string `yaml:"logs"`   // Logs directory
	} `yaml:"directories"`
	
	Logging struct {
		Level string `yaml:"level"` // Log level: info, warn, error (default: info)
		
		Access struct {
			Stdout bool   `yaml:"stdout"` // Enable access log to stdout (default: true)
			Stderr bool   `yaml:"stderr"` // Enable access log to stderr (default: false)
			Format string `yaml:"format"` // apache, nginx, text, json (default: apache)
			File   string `yaml:"file"`   // Access log file (default: access.log)
		} `yaml:"access"`
		
		Error struct {
			Stdout bool   `yaml:"stdout"` // Enable error log to stdout (default: false)
			Stderr bool   `yaml:"stderr"` // Enable error log to stderr (default: true)
			Format string `yaml:"format"` // text, json (default: text)
			File   string `yaml:"file"`   // Error log file (default: error.log)
		} `yaml:"error"`
		
		Server struct {
			Stdout bool   `yaml:"stdout"` // Enable server log to stdout (default: true)
			Stderr bool   `yaml:"stderr"` // Enable server log to stderr (default: false)
			Format string `yaml:"format"` // text, json (default: text)
			File   string `yaml:"file"`   // Server log file (default: caspaste.log)
		} `yaml:"server"`
		
		Debug struct {
			Stdout bool   `yaml:"stdout"` // Enable debug log to stdout (default: true)
			Stderr bool   `yaml:"stderr"` // Enable debug log to stderr (default: false)
			Format string `yaml:"format"` // text, json (default: text)
			File   string `yaml:"file"`   // Debug log file (default: debug.log)
		} `yaml:"debug"`
	} `yaml:"logging"`
}

// LoadYAMLConfig loads configuration from YAML file
func LoadYAMLConfig(path string) (*YAMLConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg YAMLConfig
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

// SaveYAMLConfig saves configuration to YAML file
func SaveYAMLConfig(path string, cfg *YAMLConfig) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// ResolvePlaceholders replaces placeholder values in the config with actual values
// Placeholders: {fqdn}, {data_dir}, {config_dir}
func ResolvePlaceholders(cfg *YAMLConfig, fqdn, dataDir, configDir string) {
	// Helper function to replace placeholders in a string
	replace := func(s string) string {
		s = strings.ReplaceAll(s, "{fqdn}", fqdn)
		s = strings.ReplaceAll(s, "{data_dir}", dataDir)
		s = strings.ReplaceAll(s, "{config_dir}", configDir)
		return s
	}

	// Server section
	cfg.Server.Administrator.Email = replace(cfg.Server.Administrator.Email)
	cfg.Server.Administrator.From = replace(cfg.Server.Administrator.From)

	// Web section
	cfg.Web.UI.ThemesDir = replace(cfg.Web.UI.ThemesDir)
	cfg.Web.Content.About = replace(cfg.Web.Content.About)
	cfg.Web.Content.Rules = replace(cfg.Web.Content.Rules)
	cfg.Web.Content.Terms = replace(cfg.Web.Content.Terms)
	cfg.Web.Content.Security = replace(cfg.Web.Content.Security)
	cfg.Web.Branding.Logo = replace(cfg.Web.Branding.Logo)
	cfg.Web.Branding.Favicon = replace(cfg.Web.Branding.Favicon)
	cfg.Web.Security.Contact.Email = replace(cfg.Web.Security.Contact.Email)

	// Security section
	cfg.Security.PasswordFile = replace(cfg.Security.PasswordFile)
	cfg.Security.TLS.CertFile = replace(cfg.Security.TLS.CertFile)
	cfg.Security.TLS.KeyFile = replace(cfg.Security.TLS.KeyFile)

	// Database section
	cfg.Database.Source = replace(cfg.Database.Source)

	// Set defaults for empty values that need data_dir
	if cfg.Web.UI.ThemesDir == "" {
		cfg.Web.UI.ThemesDir = dataDir + "/web/themes"
	}
}

// GetDefaultPrivateProxies returns the default trusted proxy CIDR ranges
// These are always trusted for X-Forwarded-* headers
func GetDefaultPrivateProxies() []string {
	return []string{
		"10.0.0.0/8",     // Private Class A
		"172.16.0.0/12",  // Private Class B
		"192.168.0.0/16", // Private Class C
		"127.0.0.0/8",    // Loopback IPv4
		"::1",            // Loopback IPv6
		"fc00::/7",       // Unique Local IPv6
		"fe80::/10",      // Link-Local IPv6
	}
}

// GetAllTrustedProxies returns all trusted proxies (defaults + configured)
func GetAllTrustedProxies(cfg *YAMLConfig) []string {
	proxies := GetDefaultPrivateProxies()
	proxies = append(proxies, cfg.Server.Proxy.Allowed...)
	return proxies
}

// GenerateDefaultYAMLConfig generates a default configuration file with sane defaults
func GenerateDefaultYAMLConfig(path string) error {
	defaultConfig := YAMLConfig{}

	// ============================================================================
	// SERVER CONFIGURATION
	// ============================================================================
	defaultConfig.Server.Public = true  // Default: open/public instance (no auth required)
	defaultConfig.Server.FQDN = ""      // Empty = auto-detect from X-Forwarded-Host (trusted proxies) or hostname; Set to override
	defaultConfig.Server.Listen = "all" // Listen on all interfaces (IPv4 + IPv6)
	defaultConfig.Server.Port = ""      // Empty = auto-detect available port at runtime
	defaultConfig.Server.Title = "CasPaste"
	defaultConfig.Server.TagLine = "A simple paste service"
	defaultConfig.Server.Description = "CasPaste is a simple, fast, and secure paste service for sharing code snippets and text"

	// Additional trusted proxy IPs/CIDRs to append to default private ranges
	// Default private ranges (always trusted): 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16, 127.0.0.0/8, ::1, fc00::/7, fe80::/10
	// Any IPs/CIDRs specified here are APPENDED to these defaults
	defaultConfig.Server.Proxy.Allowed = []string{}
	
	defaultConfig.Server.Administrator.Name = "CasPaste Administrator"
	defaultConfig.Server.Administrator.Email = "administrator@{fqdn}"
	defaultConfig.Server.Administrator.From = "\"CasPaste\" <no-reply@{fqdn}>"
	
	defaultConfig.Server.Timeouts.Read = 15
	defaultConfig.Server.Timeouts.Write = 15
	defaultConfig.Server.Timeouts.Idle = 60

	// ============================================================================
	// DATABASE CONFIGURATION
	// ============================================================================
	// Using modernc.org/sqlite (pure Go, no CGo)
	// Source path is relative - converted to absolute at runtime
	defaultConfig.Database.Driver = "sqlite"
	defaultConfig.Database.Source = "caspaste.db"
	defaultConfig.Database.MaxOpenConns = 25
	defaultConfig.Database.MaxIdleConns = 5
	defaultConfig.Database.CleanupPeriod = "1m"

	// ============================================================================
	// SECURITY CONFIGURATION
	// ============================================================================
	defaultConfig.Security.PasswordFile = "" // Empty = auto-generate when server.public=false
	
	// HTTP Security Headers
	defaultConfig.Security.Headers.XFrameOptions = "DENY"
	defaultConfig.Security.Headers.XContentTypeOptions = "nosniff"
	defaultConfig.Security.Headers.ContentSecurityPolicy = "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; object-src 'none'; base-uri 'self'; form-action 'self'"
	defaultConfig.Security.Headers.ReferrerPolicy = "strict-origin-when-cross-origin"
	defaultConfig.Security.Headers.PermissionsPolicy = "geolocation=(), microphone=(), camera=()"
	defaultConfig.Security.Headers.StrictTransportSecurity = "max-age=31536000; includeSubDomains"
	
	// TLS Configuration
	defaultConfig.Security.TLS.MinVersion = "1.2"
	defaultConfig.Security.TLS.CipherSuites = []string{
		"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
		"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
		"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
		"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
		"TLS_CHACHA20_POLY1305_SHA256",
	}
	defaultConfig.Security.TLS.CertFile = "/etc/caspaste/tls/cert.pem" // Auto-detected from Let's Encrypt
	defaultConfig.Security.TLS.KeyFile = "/etc/caspaste/tls/key.pem"  // Auto-detected from Let's Encrypt
	
	// Upload Security
	defaultConfig.Security.Upload.MaxFileSize = 52428800 // 50MB
	defaultConfig.Security.Upload.AllowedMIME = []string{
		"text/plain",
		"text/markdown",
		"text/html",
		"text/css",
		"text/javascript",
		"application/json",
		"application/xml",
		"application/pdf",
		"image/jpeg",
		"image/png",
		"image/gif",
		"image/svg+xml",
		"image/webp",
	}
	
	// CORS Configuration
	defaultConfig.Security.CORS.Enabled = true
	defaultConfig.Security.CORS.AllowedOrigins = []string{"*"}
	defaultConfig.Security.CORS.AllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	defaultConfig.Security.CORS.AllowedHeaders = []string{"Content-Type", "Authorization", "X-Requested-With"}
	defaultConfig.Security.CORS.MaxAge = 86400 // 24 hours

	// ============================================================================
	// WEB CONFIGURATION
	// ============================================================================

	// UI Settings
	defaultConfig.Web.UI.DefaultLifetime = "never"
	defaultConfig.Web.UI.DefaultTheme = "dark" // Accepts: "dark" (dracula), "light" (github), "auto", or full path like "dark/dracula"
	defaultConfig.Web.UI.ThemesDir = ""        // Empty = {data_dir}/web/themes (resolved at runtime)

	// Content Pages - all empty = auto-generated from embedded defaults
	// If set, paths are relative to {data_dir}/web/docs unless absolute
	defaultConfig.Web.Content.About = ""    // Empty = auto-generated
	defaultConfig.Web.Content.Rules = ""    // Empty = auto-generated
	defaultConfig.Web.Content.Terms = ""    // Empty = auto-generated
	defaultConfig.Web.Content.Security = "" // Empty = auto-generated security.txt

	// Branding - can be local paths or URLs
	defaultConfig.Web.Branding.Logo = ""    // Empty = use embedded default
	defaultConfig.Web.Branding.Favicon = "" // Empty = use embedded default
	
	// Security Contact (for security.txt)
	defaultConfig.Web.Security.Contact.Email = "security@{fqdn}"
	defaultConfig.Web.Security.Contact.Name = "Security Team"
	
	// SEO / Robots
	defaultConfig.Web.SEO.Robots.Allow = "*"
	defaultConfig.Web.SEO.Robots.Deny = "/settings,/history"
	defaultConfig.Web.SEO.Robots.Agents.Deny = []string{
		"GPTBot",
		"ChatGPT-User",
		"Google-Extended",
		"CCBot",
		"anthropic-ai",
		"Claude-Web",
		"cohere-ai",
		"Omgilibot",
		"FacebookBot",
		"Diffbot",
	}

	// ============================================================================
	// LIMITS & RATE LIMITING
	// ============================================================================
	defaultConfig.Limits.TitleMaxLength = 100
	defaultConfig.Limits.BodyMaxLength = 52428800 // 50MB
	defaultConfig.Limits.MaxPasteLifetime = "never"
	
	// Rate limiting for GET requests
	defaultConfig.Limits.RateLimit.GetPastes.Per5Min = 50
	defaultConfig.Limits.RateLimit.GetPastes.Per15Min = 100
	defaultConfig.Limits.RateLimit.GetPastes.Per1Hour = 500
	
	// Rate limiting for POST requests
	defaultConfig.Limits.RateLimit.NewPastes.Per5Min = 15
	defaultConfig.Limits.RateLimit.NewPastes.Per15Min = 30
	defaultConfig.Limits.RateLimit.NewPastes.Per1Hour = 40

	// ============================================================================
	// DIRECTORIES
	// ============================================================================
	// Platform-specific defaults
	defaultConfig.Directories.Data = "/var/lib/caspaste"
	defaultConfig.Directories.Config = "/etc/caspaste"
	defaultConfig.Directories.Db = "/var/lib/caspaste/db"    // Database directory - if under data dir, included in data backup
	defaultConfig.Directories.Cache = "/var/cache/caspaste"
	defaultConfig.Directories.Logs = "/var/log/caspaste"

	// ============================================================================
	// LOGGING
	// ============================================================================
	defaultConfig.Logging.Level = "info" // info, warn, error (default: info)
	
	// Access Log (HTTP requests)
	defaultConfig.Logging.Access.Stdout = false  // Don't clutter console with every request
	defaultConfig.Logging.Access.Stderr = false
	defaultConfig.Logging.Access.Format = "apache" // apache (combined), nginx, text, json
	defaultConfig.Logging.Access.File = "access.log"
	
	// Error Log (ERROR messages)
	defaultConfig.Logging.Error.Stdout = false
	defaultConfig.Logging.Error.Stderr = true // Errors to stderr by default
	defaultConfig.Logging.Error.Format = "text" // text, json
	defaultConfig.Logging.Error.File = "error.log"
	
	// Server Log (INFO messages)
	defaultConfig.Logging.Server.Stdout = true // Show info messages on console
	defaultConfig.Logging.Server.Stderr = false
	defaultConfig.Logging.Server.Format = "text" // text, json
	defaultConfig.Logging.Server.File = "caspaste.log"
	
	// Debug Log (DEBUG messages, only with --debug flag)
	defaultConfig.Logging.Debug.Stdout = true
	defaultConfig.Logging.Debug.Stderr = false
	defaultConfig.Logging.Debug.Format = "text" // text, json
	defaultConfig.Logging.Debug.File = "debug.log"

	// Write to file
	data, err := yaml.Marshal(defaultConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal default config: %w", err)
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write default config: %w", err)
	}

	return nil
}
