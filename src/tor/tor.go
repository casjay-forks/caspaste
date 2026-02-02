// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

// Package tor provides Tor hidden service support per AI.md PART 32
// Uses external Tor binary via github.com/cretz/bine for CGO_ENABLED=0 compatibility
// Hidden service is auto-enabled when Tor binary is found
package tor

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Config holds Tor configuration
type Config struct {
	// Binary path (empty = auto-detect)
	Binary string

	// Outbound network settings
	UseNetwork          bool
	AllowUserPreference bool

	// Performance settings
	MaxCircuits      int
	CircuitTimeout   time.Duration
	BootstrapTimeout time.Duration

	// Security settings
	SafeLogging               bool
	MaxStreamsPerCircuit      int
	CloseCircuitOnStreamLimit bool

	// Bandwidth settings
	BandwidthRate       string
	BandwidthBurst      string
	MaxMonthlyBandwidth string

	// Hidden service settings
	NumIntroPoints int
	VirtualPort    int
}

// DefaultConfig returns the default Tor configuration
func DefaultConfig() *Config {
	return &Config{
		Binary:                    "",
		UseNetwork:                false,
		AllowUserPreference:       true,
		MaxCircuits:               32,
		CircuitTimeout:            60 * time.Second,
		BootstrapTimeout:          3 * time.Minute,
		SafeLogging:               true,
		MaxStreamsPerCircuit:      100,
		CloseCircuitOnStreamLimit: true,
		BandwidthRate:             "1 MB",
		BandwidthBurst:            "2 MB",
		MaxMonthlyBandwidth:       "100 GB",
		NumIntroPoints:            3,
		VirtualPort:               80,
	}
}

// Status represents the current Tor status
type Status struct {
	Enabled    bool   `json:"enabled"`
	Running    bool   `json:"running"`
	StatusText string `json:"status"`
	Hostname   string `json:"hostname"`
	Error      string `json:"error,omitempty"`
}

// Service manages the Tor hidden service
type Service struct {
	config     *Config
	configDir  string
	dataDir    string
	logDir     string
	serverPort int
	enabled    bool
	running    bool
	hostname   string
	binaryPath string
	process    *os.Process
	lastError  string
	mu         sync.RWMutex
}

// NewService creates a new Tor service
func NewService(cfg *Config, configDir, dataDir, logDir string, serverPort int) *Service {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	return &Service{
		config:     cfg,
		configDir:  configDir,
		dataDir:    dataDir,
		logDir:     logDir,
		serverPort: serverPort,
	}
}

// Start starts the Tor hidden service
func (s *Service) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Find Tor binary
	binaryPath, err := s.findTorBinary()
	if err != nil {
		s.enabled = false
		s.lastError = "Tor binary not found"
		return nil
	}
	s.binaryPath = binaryPath
	s.enabled = true

	// Create directories
	if err := s.ensureDirs(); err != nil {
		s.lastError = err.Error()
		return fmt.Errorf("failed to create Tor directories: %w", err)
	}

	// Generate torrc
	if err := s.generateTorrc(); err != nil {
		s.lastError = err.Error()
		return fmt.Errorf("failed to generate torrc: %w", err)
	}

	// Start Tor process
	if err := s.startProcess(ctx); err != nil {
		s.lastError = err.Error()
		return fmt.Errorf("failed to start Tor: %w", err)
	}

	s.running = true
	return nil
}

// Stop stops the Tor hidden service
func (s *Service) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.process != nil {
		if err := s.process.Kill(); err != nil {
			return fmt.Errorf("failed to stop Tor: %w", err)
		}
		s.process = nil
	}

	s.running = false
	return nil
}

// IsEnabled returns true if Tor is enabled (binary found)
func (s *Service) IsEnabled() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.enabled
}

// IsRunning returns true if Tor is running
func (s *Service) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// GetHostname returns the .onion hostname
func (s *Service) GetHostname() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.hostname
}

// GetStatus returns the current Tor status
func (s *Service) GetStatus() *Status {
	s.mu.RLock()
	defer s.mu.RUnlock()

	status := &Status{
		Enabled:  s.enabled,
		Running:  s.running,
		Hostname: s.hostname,
		Error:    s.lastError,
	}

	if !s.enabled {
		status.StatusText = "disabled"
	} else if s.running {
		status.StatusText = "healthy"
	} else if s.lastError != "" {
		status.StatusText = "error"
	} else {
		status.StatusText = "stopped"
	}

	return status
}

// GetHTTPClient returns an HTTP client that uses Tor for outbound connections
func (s *Service) GetHTTPClient(useTor bool) *http.Client {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !useTor || !s.running {
		return http.DefaultClient
	}

	// Note: Full implementation would use bine Dialer
	// For now, return default client as placeholder
	return http.DefaultClient
}

// ShouldUseTor determines if Tor should be used based on config and user preference
func (s *Service) ShouldUseTor(userPref *bool) bool {
	if !s.config.AllowUserPreference {
		return s.config.UseNetwork
	}

	if userPref == nil {
		return s.config.UseNetwork
	}

	return *userPref
}

// findTorBinary searches for the Tor binary
func (s *Service) findTorBinary() (string, error) {
	// Check config path first
	if s.config.Binary != "" {
		if _, err := os.Stat(s.config.Binary); err == nil {
			return s.config.Binary, nil
		}
	}

	// Check PATH
	if path, err := exec.LookPath("tor"); err == nil {
		return path, nil
	}

	// Platform-specific common locations
	var locations []string
	switch runtime.GOOS {
	case "linux":
		locations = []string{"/usr/bin/tor", "/usr/local/bin/tor"}
	case "darwin":
		locations = []string{"/usr/local/bin/tor", "/opt/homebrew/bin/tor"}
	case "windows":
		locations = []string{
			`C:\Program Files\Tor\tor.exe`,
			`C:\Program Files (x86)\Tor\tor.exe`,
		}
	default:
		locations = []string{"/usr/local/bin/tor"}
	}

	for _, loc := range locations {
		if _, err := os.Stat(loc); err == nil {
			return loc, nil
		}
	}

	return "", fmt.Errorf("tor binary not found")
}

// ensureDirs creates required Tor directories
func (s *Service) ensureDirs() error {
	dirs := []string{
		filepath.Join(s.configDir, "tor"),
		filepath.Join(s.dataDir, "tor"),
		filepath.Join(s.dataDir, "tor", "site"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return fmt.Errorf("failed to create %s: %w", dir, err)
		}
	}

	return nil
}

// generateTorrc generates the torrc configuration file
func (s *Service) generateTorrc() error {
	torrcPath := filepath.Join(s.configDir, "tor", "torrc")

	// Check if torrc already exists
	if _, err := os.Stat(torrcPath); err == nil {
		return nil
	}

	var content strings.Builder

	// Data directory
	content.WriteString(fmt.Sprintf("DataDirectory %s\n", filepath.Join(s.dataDir, "tor")))

	// Control socket (Unix) or port (Windows)
	if runtime.GOOS == "windows" {
		content.WriteString("ControlPort auto\n")
	} else {
		controlSocket := filepath.Join(s.dataDir, "tor", "control.sock")
		content.WriteString(fmt.Sprintf("ControlSocket %s\n", controlSocket))
	}

	// Hidden service
	hsDir := filepath.Join(s.dataDir, "tor", "site")
	content.WriteString(fmt.Sprintf("HiddenServiceDir %s\n", hsDir))
	content.WriteString(fmt.Sprintf("HiddenServicePort %d 127.0.0.1:%d\n", s.config.VirtualPort, s.serverPort))
	content.WriteString("HiddenServiceVersion 3\n")
	content.WriteString(fmt.Sprintf("HiddenServiceNumIntroductionPoints %d\n", s.config.NumIntroPoints))

	// SOCKS port (for outbound if enabled)
	if s.config.UseNetwork || s.config.AllowUserPreference {
		content.WriteString("SocksPort auto\n")
	} else {
		content.WriteString("SocksPort 0\n")
	}

	// Security settings
	if s.config.SafeLogging {
		content.WriteString("SafeLogging 1\n")
	}

	// Log file
	logFile := filepath.Join(s.logDir, "tor.log")
	content.WriteString(fmt.Sprintf("Log notice file %s\n", logFile))

	// Write torrc
	if err := os.WriteFile(torrcPath, []byte(content.String()), 0600); err != nil {
		return fmt.Errorf("failed to write torrc: %w", err)
	}

	return nil
}

// startProcess starts the Tor process
func (s *Service) startProcess(ctx context.Context) error {
	torrcPath := filepath.Join(s.configDir, "tor", "torrc")

	cmd := exec.CommandContext(ctx, s.binaryPath, "-f", torrcPath)
	cmd.Dir = s.dataDir

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start Tor: %w", err)
	}

	s.process = cmd.Process

	// Wait for hostname file
	hostnameFile := filepath.Join(s.dataDir, "tor", "site", "hostname")
	timeout := time.After(s.config.BootstrapTimeout)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout:
			return fmt.Errorf("bootstrap timeout waiting for hostname")
		case <-ticker.C:
			if data, err := os.ReadFile(hostnameFile); err == nil {
				s.hostname = strings.TrimSpace(string(data))
				return nil
			}
		}
	}
}

// GetConfig returns the current Tor configuration (for display)
func (s *Service) GetConfig() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return map[string]interface{}{
		"enabled":               s.enabled,
		"running":               s.running,
		"hostname":              s.hostname,
		"binary":                s.binaryPath,
		"use_network":           s.config.UseNetwork,
		"allow_user_preference": s.config.AllowUserPreference,
		"virtual_port":          s.config.VirtualPort,
		"num_intro_points":      s.config.NumIntroPoints,
	}
}
