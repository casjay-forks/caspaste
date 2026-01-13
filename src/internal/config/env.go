// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package config

import (
	"os"
	"strconv"

	"github.com/casjay-forks/caspaste/src/internal/validation"
)

// getEnv tries CASPASTE_* first, then LENPASTE_* for backward compatibility
func getEnv(name string) string {
	if val := os.Getenv("CASPASTE_" + name); val != "" {
		return val
	}
	if val := os.Getenv("LENPASTE_" + name); val != "" {
		return val
	}
	return ""
}

// ApplyEnvironmentOverrides applies environment variables to config
// Environment variables override config file values
func ApplyEnvironmentOverrides(cfg *YAMLConfig) {
	// Server settings
	if val := getEnv("ADDRESS"); val != "" {
		cfg.Server.Address = val
	}
	if val := getEnv("BIND"); val != "" {
		cfg.Server.Bind = val
	}
	if val := getEnv("PORT"); val != "" {
		cfg.Server.Port = val // Now string format: "8080" or "8080,64453"
	}
	if val := getEnv("SERVER_TITLE"); val != "" {
		cfg.Server.Title = val
	}
	if val := getEnv("TITLE"); val != "" { // Alternative
		cfg.Server.Title = val
	}
	if val := getEnv("TRUST_REVERSE_PROXY"); val != "" {
		cfg.Server.TrustReverseProxy = validation.IsTruthy(val)
	}

	// Server administrator
	if val := getEnv("ADMIN_NAME"); val != "" {
		cfg.Server.Administrator.Name = val
	}
	if val := getEnv("SERVER_ADMINISTRATOR_NAME"); val != "" {
		cfg.Server.Administrator.Name = val
	}
	if val := getEnv("ADMIN_EMAIL"); val != "" {
		cfg.Server.Administrator.Email = val
	}
	if val := getEnv("ADMIN_MAIL"); val != "" { // Alternative
		cfg.Server.Administrator.Email = val
	}
	if val := getEnv("SERVER_ADMINISTRATOR_EMAIL"); val != "" {
		cfg.Server.Administrator.Email = val
	}
	if val := getEnv("SERVER_ADMINISTRATOR_FROM"); val != "" {
		cfg.Server.Administrator.From = val
	}

	// Web security contact
	if val := getEnv("WEB_SECURITY_CONTACT_EMAIL"); val != "" {
		cfg.Web.Security.Contact.Email = val
	}
	if val := getEnv("WEB_SECURITY_CONTACT_NAME"); val != "" {
		cfg.Web.Security.Contact.Name = val
	}

	// Site robots -> Web.SEO.Robots
	if val := getEnv("SITE_ROBOTS_ALLOW"); val != "" {
		cfg.Web.SEO.Robots.Allow = val
	}
	if val := getEnv("SITE_ROBOTS_DENY"); val != "" {
		cfg.Web.SEO.Robots.Deny = val
	}
	if val := getEnv("ROBOTS_DISALLOW"); val != "" { // Legacy compatibility
		if validation.IsTruthy(val) {
			cfg.Web.SEO.Robots.Deny = "/"
		}
	}

	// Branding -> Web.Branding
	if val := getEnv("BRANDING_LOGO"); val != "" {
		cfg.Web.Branding.Logo = val
	}
	if val := getEnv("BRANDING_FAVICON"); val != "" {
		cfg.Web.Branding.Favicon = val
	}

	// Database settings
	if val := getEnv("DB_DRIVER"); val != "" {
		cfg.Database.Driver = val
	}
	if val := getEnv("DB_SOURCE"); val != "" {
		cfg.Database.Source = val
	}
	if val := getEnv("DB_MAX_OPEN_CONNS"); val != "" {
		if num, err := strconv.Atoi(val); err == nil {
			cfg.Database.MaxOpenConns = num
		}
	}
	if val := getEnv("DB_MAX_IDLE_CONNS"); val != "" {
		if num, err := strconv.Atoi(val); err == nil {
			cfg.Database.MaxIdleConns = num
		}
	}
	if val := getEnv("DB_CLEANUP_PERIOD"); val != "" {
		cfg.Database.CleanupPeriod = val
	}

	// Security settings
	if val := getEnv("PASSWORD_FILE"); val != "" {
		cfg.Security.PasswordFile = val
	}
	if val := getEnv("CASPASSWD_FILE"); val != "" { // Alternative name
		cfg.Security.PasswordFile = val
	}

	// Limits settings
	if val := getEnv("TITLE_MAX_LENGTH"); val != "" {
		if num, err := strconv.Atoi(val); err == nil {
			cfg.Limits.TitleMaxLength = num
		}
	}
	if val := getEnv("BODY_MAX_LENGTH"); val != "" {
		if num, err := strconv.Atoi(val); err == nil {
			cfg.Limits.BodyMaxLength = num
		}
	}
	if val := getEnv("MAX_PASTE_LIFETIME"); val != "" {
		cfg.Limits.MaxPasteLifetime = val
	}
	// Rate limits - GET pastes
	if val := getEnv("GET_PASTES_PER_5MIN"); val != "" {
		if num, err := strconv.ParseUint(val, 10, 32); err == nil {
			cfg.Limits.RateLimit.GetPastes.Per5Min = uint(num)
		}
	}
	if val := getEnv("GET_PASTES_PER_15MIN"); val != "" {
		if num, err := strconv.ParseUint(val, 10, 32); err == nil {
			cfg.Limits.RateLimit.GetPastes.Per15Min = uint(num)
		}
	}
	if val := getEnv("GET_PASTES_PER_1HOUR"); val != "" {
		if num, err := strconv.ParseUint(val, 10, 32); err == nil {
			cfg.Limits.RateLimit.GetPastes.Per1Hour = uint(num)
		}
	}
	
	// Rate limits - NEW pastes
	if val := getEnv("NEW_PASTES_PER_5MIN"); val != "" {
		if num, err := strconv.ParseUint(val, 10, 32); err == nil {
			cfg.Limits.RateLimit.NewPastes.Per5Min = uint(num)
		}
	}
	if val := getEnv("NEW_PASTES_PER_15MIN"); val != "" {
		if num, err := strconv.ParseUint(val, 10, 32); err == nil {
			cfg.Limits.RateLimit.NewPastes.Per15Min = uint(num)
		}
	}
	if val := getEnv("NEW_PASTES_PER_1HOUR"); val != "" {
		if num, err := strconv.ParseUint(val, 10, 32); err == nil {
			cfg.Limits.RateLimit.NewPastes.Per1Hour = uint(num)
		}
	}

	// UI settings -> Web.UI
	if val := getEnv("UI_DEFAULT_LIFETIME"); val != "" {
		cfg.Web.UI.DefaultLifetime = val
	}
	if val := getEnv("UI_DEFAULT_THEME"); val != "" {
		cfg.Web.UI.DefaultTheme = val
	}
	if val := getEnv("UI_THEMES_DIR"); val != "" {
		cfg.Web.UI.ThemesDir = val
	}

	// Content settings -> Web.Content
	if val := getEnv("CONTENT_ABOUT"); val != "" {
		cfg.Web.Content.About = val
	}
	if val := getEnv("SERVER_ABOUT"); val != "" { // Legacy
		cfg.Web.Content.About = val
	}
	if val := getEnv("CONTENT_RULES"); val != "" {
		cfg.Web.Content.Rules = val
	}
	if val := getEnv("SERVER_RULES"); val != "" { // Legacy
		cfg.Web.Content.Rules = val
	}
	if val := getEnv("CONTENT_TERMS"); val != "" {
		cfg.Web.Content.Terms = val
	}
	if val := getEnv("SERVER_TERMS"); val != "" { // Legacy
		cfg.Web.Content.Terms = val
	}
	if val := getEnv("CONTENT_SECURITY"); val != "" {
		cfg.Web.Content.Security = val
	}

	// Directory settings
	if val := getEnv("CACHE_DIR"); val != "" {
		cfg.Directories.Cache = val
	}
	if val := getEnv("LOGS_DIR"); val != "" {
		cfg.Directories.Logs = val
	}
}
