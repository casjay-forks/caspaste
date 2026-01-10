// Copyright (C) 2021-2023 Leonid Maslakov.

// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE file for details.

package service

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// ServiceConfig holds service configuration
type ServiceConfig struct {
	Name        string
	DisplayName string
	Description string
	Executable  string
	Args        []string
	WorkingDir  string
	User        string
}

// Manager handles service operations
type Manager struct {
	config ServiceConfig
}

// New creates a new service manager
func New(config ServiceConfig) *Manager {
	return &Manager{config: config}
}

// Install installs the service
func (m *Manager) Install() error {
	switch runtime.GOOS {
	case "linux":
		return m.installLinux()
	case "darwin":
		return m.installDarwin()
	case "windows":
		return m.installWindows()
	case "freebsd", "openbsd":
		return m.installBSD()
	default:
		return fmt.Errorf("service installation not supported on %s", runtime.GOOS)
	}
}

// Uninstall removes the service
func (m *Manager) Uninstall() error {
	switch runtime.GOOS {
	case "linux":
		return m.uninstallLinux()
	case "darwin":
		return m.uninstallDarwin()
	case "windows":
		return m.uninstallWindows()
	case "freebsd", "openbsd":
		return m.uninstallBSD()
	default:
		return fmt.Errorf("service uninstall not supported on %s", runtime.GOOS)
	}
}

// Start starts the service
func (m *Manager) Start() error {
	switch runtime.GOOS {
	case "linux":
		return m.controlLinux("start")
	case "darwin":
		return m.controlDarwin("start")
	case "windows":
		return m.controlWindows("start")
	case "freebsd", "openbsd":
		return m.controlBSD("start")
	default:
		return fmt.Errorf("service start not supported on %s", runtime.GOOS)
	}
}

// Stop stops the service
func (m *Manager) Stop() error {
	switch runtime.GOOS {
	case "linux":
		return m.controlLinux("stop")
	case "darwin":
		return m.controlDarwin("stop")
	case "windows":
		return m.controlWindows("stop")
	case "freebsd", "openbsd":
		return m.controlBSD("stop")
	default:
		return fmt.Errorf("service stop not supported on %s", runtime.GOOS)
	}
}

// Restart restarts the service
func (m *Manager) Restart() error {
	switch runtime.GOOS {
	case "linux":
		return m.controlLinux("restart")
	case "darwin":
		return m.controlDarwin("restart")
	case "windows":
		return m.controlWindows("restart")
	case "freebsd", "openbsd":
		return m.controlBSD("restart")
	default:
		return fmt.Errorf("service restart not supported on %s", runtime.GOOS)
	}
}

// Reload reloads the service configuration
func (m *Manager) Reload() error {
	switch runtime.GOOS {
	case "linux":
		return m.controlLinux("reload")
	case "darwin":
		return m.controlDarwin("reload")
	case "windows":
		return fmt.Errorf("reload not supported on Windows, use restart instead")
	case "freebsd", "openbsd":
		return m.controlBSD("reload")
	default:
		return fmt.Errorf("service reload not supported on %s", runtime.GOOS)
	}
}

// Disable disables the service from starting at boot
func (m *Manager) Disable() error {
	switch runtime.GOOS {
	case "linux":
		return m.disableLinux()
	case "darwin":
		return m.disableDarwin()
	case "windows":
		return m.disableWindows()
	case "freebsd", "openbsd":
		return m.disableBSD()
	default:
		return fmt.Errorf("service disable not supported on %s", runtime.GOOS)
	}
}

// Status checks the service status
func (m *Manager) Status() error {
	switch runtime.GOOS {
	case "linux":
		return m.statusLinux()
	case "darwin":
		return m.statusDarwin()
	case "windows":
		return m.statusWindows()
	case "freebsd", "openbsd":
		return m.statusBSD()
	default:
		return fmt.Errorf("service status not supported on %s", runtime.GOOS)
	}
}

// runCommand executes a command and returns error if it fails
func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// buildArgs creates the arguments string for service configuration
func (m *Manager) buildArgs() string {
	return strings.Join(m.config.Args, " ")
}
