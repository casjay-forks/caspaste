// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package config

import (
	"os"
	"strconv"
	"strings"
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
		// Support format: hostname:port or :port
		// Examples: lp.pste.us:8080, myserver.com:80, :8080

		// If it's a FQDN:port format, convert to :port (bind to all interfaces)
		// The FQDN is used for reverse proxy setups and display purposes
		if strings.Contains(val, ":") {
			parts := strings.Split(val, ":")
			if len(parts) == 2 {
				hostname := parts[0]
				portStr := parts[1]

				// If hostname part contains dots (FQDN) or is not empty, bind to all interfaces
				if strings.Contains(hostname, ".") || (hostname != "" && hostname != "0.0.0.0" && hostname != "::") {
					// FQDN or specific IP - bind to all interfaces on the port
					cfg.Server.Address = ":" + portStr
				} else {
					// Already in :port format
					cfg.Server.Address = val
				}

				// Extract and set port
				if port, err := strconv.Atoi(portStr); err == nil {
					cfg.Server.Port = port
				}
			}
		} else {
			cfg.Server.Address = val
		}
	}
	if val := getEnv("PORT"); val != "" {
		if port, err := strconv.Atoi(val); err == nil {
			cfg.Server.Port = port
		}
	}
	if val := getEnv("ADMIN_NAME"); val != "" {
		cfg.Server.AdminName = val
	}
	if val := getEnv("ADMIN_EMAIL"); val != "" {
		cfg.Server.AdminEmail = val
	}
	if val := getEnv("ADMIN_MAIL"); val != "" { // Alternative name
		cfg.Server.AdminEmail = val
	}
	if val := getEnv("ROBOTS_DISALLOW"); val != "" {
		cfg.Server.RobotsDisallow = strings.ToLower(val) == "true" || val == "1"
	}
	if val := getEnv("TRUST_REVERSE_PROXY"); val != "" {
		cfg.Server.TrustReverseProxy = strings.ToLower(val) == "true" || val == "1"
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
	if val := getEnv("GET_PASTES_PER_5MIN"); val != "" {
		if num, err := strconv.ParseUint(val, 10, 32); err == nil {
			cfg.Limits.GetPastesPer5Min = uint(num)
		}
	}
	if val := getEnv("GET_PASTES_PER_15MIN"); val != "" {
		if num, err := strconv.ParseUint(val, 10, 32); err == nil {
			cfg.Limits.GetPastesPer15Min = uint(num)
		}
	}
	if val := getEnv("GET_PASTES_PER_1HOUR"); val != "" {
		if num, err := strconv.ParseUint(val, 10, 32); err == nil {
			cfg.Limits.GetPastesPer1Hour = uint(num)
		}
	}
	if val := getEnv("NEW_PASTES_PER_5MIN"); val != "" {
		if num, err := strconv.ParseUint(val, 10, 32); err == nil {
			cfg.Limits.NewPastesPer5Min = uint(num)
		}
	}
	if val := getEnv("NEW_PASTES_PER_15MIN"); val != "" {
		if num, err := strconv.ParseUint(val, 10, 32); err == nil {
			cfg.Limits.NewPastesPer15Min = uint(num)
		}
	}
	if val := getEnv("NEW_PASTES_PER_1HOUR"); val != "" {
		if num, err := strconv.ParseUint(val, 10, 32); err == nil {
			cfg.Limits.NewPastesPer1Hour = uint(num)
		}
	}

	// UI settings
	if val := getEnv("UI_DEFAULT_LIFETIME"); val != "" {
		cfg.UI.DefaultLifetime = val
	}
	if val := getEnv("UI_DEFAULT_THEME"); val != "" {
		cfg.UI.DefaultTheme = val
	}
	if val := getEnv("UI_THEMES_DIR"); val != "" {
		cfg.UI.ThemesDir = val
	}

	// Content settings
	if val := getEnv("SERVER_ABOUT"); val != "" {
		cfg.Content.AboutFile = val
	}
	if val := getEnv("ABOUT_FILE"); val != "" { // Alternative name
		cfg.Content.AboutFile = val
	}
	if val := getEnv("SERVER_RULES"); val != "" {
		cfg.Content.RulesFile = val
	}
	if val := getEnv("RULES_FILE"); val != "" { // Alternative name
		cfg.Content.RulesFile = val
	}
	if val := getEnv("SERVER_TERMS"); val != "" {
		cfg.Content.TermsFile = val
	}
	if val := getEnv("TERMS_FILE"); val != "" { // Alternative name
		cfg.Content.TermsFile = val
	}

	// Directory settings
	if val := getEnv("CACHE_DIR"); val != "" {
		cfg.Directories.Cache = val
	}
	if val := getEnv("LOGS_DIR"); val != "" {
		cfg.Directories.Logs = val
	}
}
