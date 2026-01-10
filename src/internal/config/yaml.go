
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// YAMLConfig represents the YAML configuration file structure
type YAMLConfig struct {
	Server struct {
		Port              int    `yaml:"port"`
		Address           string `yaml:"address"`
		AdminName         string `yaml:"admin_name"`
		AdminEmail        string `yaml:"admin_email"`
		RobotsDisallow    bool   `yaml:"robots_disallow"`
		TrustReverseProxy bool   `yaml:"trust_reverse_proxy"`
	} `yaml:"server"`

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
		AboutFile string `yaml:"about_file"`
		RulesFile string `yaml:"rules_file"`
		TermsFile string `yaml:"terms_file"`
	} `yaml:"content"`

	Directories struct {
		Cache string `yaml:"cache"`
		Logs  string `yaml:"logs"`
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
	defaultConfig.Server.Port = 8080
	defaultConfig.Server.Address = ":8080"
	defaultConfig.Server.AdminName = ""
	defaultConfig.Server.AdminEmail = ""
	defaultConfig.Server.RobotsDisallow = false
	defaultConfig.Server.TrustReverseProxy = false

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
	defaultConfig.UI.DefaultTheme = "dracula"
	defaultConfig.UI.ThemesDir = ""

	// Content defaults
	defaultConfig.Content.AboutFile = ""
	defaultConfig.Content.RulesFile = ""
	defaultConfig.Content.TermsFile = ""

	// Directory defaults (platform-specific, empty = auto-detect)
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
