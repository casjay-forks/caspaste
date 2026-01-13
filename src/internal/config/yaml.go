
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// YAMLConfig represents the YAML configuration file structure
type YAMLConfig struct {
	Server struct {
		Address           string `yaml:"address"` // Public FQDN
		Bind              string `yaml:"bind"`    // Bind address (::, 0.0.0.0, specific IP)
		Port              string `yaml:"port"`    // "8080" or "8080,64453"
		Title             string `yaml:"title"`
		TrustReverseProxy bool   `yaml:"trust_reverse_proxy"`
		Administrator     struct {
			Name  string `yaml:"name"`
			Email string `yaml:"email"`
			From  string `yaml:"from"`
		} `yaml:"administrator"`
	} `yaml:"server"`

	Web struct {
		Security struct {
			Contact struct {
				Email string `yaml:"email"`
				Name  string `yaml:"name"`
			} `yaml:"contact"`
		} `yaml:"security"`
	} `yaml:"web"`

	Site struct {
		Robots struct {
			Allow  string `yaml:"allow"`
			Deny   string `yaml:"deny"`
			Agents struct {
				Deny []string `yaml:"deny"`
			} `yaml:"agents"`
		} `yaml:"robots"`
	} `yaml:"site"`

	Branding struct {
		Logo    string `yaml:"logo"`
		Favicon string `yaml:"favicon"`
	} `yaml:"branding"`

	Database struct {
		Driver        string `yaml:"driver"`
		Source        string `yaml:"source"`
		MaxOpenConns  int    `yaml:"max_open_conns"`
		MaxIdleConns  int    `yaml:"max_idle_conns"`
		CleanupPeriod string `yaml:"cleanup_period"`
	} `yaml:"database"`

	Security struct {
		PasswordFile string `yaml:"password_file"`
	} `yaml:"security"`

	Limits struct {
		TitleMaxLength    int    `yaml:"title_max_length"`
		BodyMaxLength     int    `yaml:"body_max_length"`
		MaxPasteLifetime  string `yaml:"max_paste_lifetime"`
		GetPastesPer5Min  uint   `yaml:"get_pastes_per_5min"`
		GetPastesPer15Min uint   `yaml:"get_pastes_per_15min"`
		GetPastesPer1Hour uint   `yaml:"get_pastes_per_1hour"`
		NewPastesPer5Min  uint   `yaml:"new_pastes_per_5min"`
		NewPastesPer15Min uint   `yaml:"new_pastes_per_15min"`
		NewPastesPer1Hour uint   `yaml:"new_pastes_per_1hour"`
	} `yaml:"limits"`

	UI struct {
		DefaultLifetime string `yaml:"default_lifetime"`
		DefaultTheme    string `yaml:"default_theme"`
		ThemesDir       string `yaml:"themes_dir"`
	} `yaml:"ui"`

	Content struct {
		About    string `yaml:"about"`
		Rules    string `yaml:"rules"`
		Terms    string `yaml:"terms"`
		Security string `yaml:"security"`
	} `yaml:"content"`

	Directories struct {
		Data   string `yaml:"data"`
		Config string `yaml:"config"`
		Db     string `yaml:"db"`
		Cache  string `yaml:"cache"`
		Logs   string `yaml:"logs"`
	} `yaml:"directories"`
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

// GenerateDefaultYAMLConfig generates a default configuration file
func GenerateDefaultYAMLConfig(path string) error {
	defaultConfig := YAMLConfig{}

	// Server defaults
	defaultConfig.Server.Address = ""  // Auto-detected
	defaultConfig.Server.Bind = "::"   // IPv4 + IPv6
	defaultConfig.Server.Port = ""     // Random port in 64xxx range
	defaultConfig.Server.Title = "CasPaste"
	defaultConfig.Server.TrustReverseProxy = false
	defaultConfig.Server.Administrator.Name = "CasPaste"
	defaultConfig.Server.Administrator.Email = "administrator@{fqdn}"
	defaultConfig.Server.Administrator.From = "\"CasPaste\" <no-reply@{fqdn}>"

	// Web defaults
	defaultConfig.Web.Security.Contact.Email = "administrator@{fqdn}"
	defaultConfig.Web.Security.Contact.Name = "Server Administrator"

	// Site defaults
	defaultConfig.Site.Robots.Allow = "*"
	defaultConfig.Site.Robots.Deny = "/settings,/history"
	defaultConfig.Site.Robots.Agents.Deny = []string{
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

	// Branding defaults
	defaultConfig.Branding.Logo = ""
	defaultConfig.Branding.Favicon = ""

	// Database defaults
	// Note: Using "sqlite" (modernc.org/sqlite - pure Go, no CGo)
	// Source is relative - converted to absolute at runtime:
	//   Docker (--data /data): /data/db/sqlite/caspaste.db
	//   Standalone: {data_dir}/db/caspaste.db
	defaultConfig.Database.Driver = "sqlite"
	defaultConfig.Database.Source = "caspaste.db"
	defaultConfig.Database.MaxOpenConns = 25
	defaultConfig.Database.MaxIdleConns = 5
	defaultConfig.Database.CleanupPeriod = "1m"

	// Security defaults
	defaultConfig.Security.PasswordFile = ""

	// Limits defaults
	defaultConfig.Limits.TitleMaxLength = 100
	defaultConfig.Limits.BodyMaxLength = 52428800 // 50MB
	defaultConfig.Limits.MaxPasteLifetime = "never"
	defaultConfig.Limits.GetPastesPer5Min = 50
	defaultConfig.Limits.GetPastesPer15Min = 100
	defaultConfig.Limits.GetPastesPer1Hour = 500
	defaultConfig.Limits.NewPastesPer5Min = 15
	defaultConfig.Limits.NewPastesPer15Min = 30
	defaultConfig.Limits.NewPastesPer1Hour = 40

	// UI defaults
	defaultConfig.UI.DefaultLifetime = "never"
	defaultConfig.UI.DefaultTheme = "dark/dracula"
	defaultConfig.UI.ThemesDir = ""

	// Content defaults
	// Empty = use embedded defaults with variable replacement
	// Set path to override with custom file
	defaultConfig.Content.About = ""
	defaultConfig.Content.Rules = ""
	defaultConfig.Content.Terms = ""
	defaultConfig.Content.Security = ""

	// Directory defaults (platform-specific, empty = auto-detect)
	defaultConfig.Directories.Data = ""
	defaultConfig.Directories.Config = ""
	defaultConfig.Directories.Db = ""
	defaultConfig.Directories.Cache = ""
	defaultConfig.Directories.Logs = ""

	// Marshal to YAML
	data, err := yaml.Marshal(&defaultConfig)
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(path, data, 0644)
}

// SaveYAMLConfig saves configuration to YAML file
// Uses atomic write (temp file + rename) to prevent corruption
func SaveYAMLConfig(path string, cfg *YAMLConfig) error {
	// Marshal to YAML
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to temp file first
	tempPath := path + ".tmp"
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write temp config: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tempPath, path); err != nil {
		os.Remove(tempPath) // Clean up temp file
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}
