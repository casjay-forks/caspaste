// Copyright (C) 2021-2023 Leonid Maslakov.

// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE file for details.

//go:build windows
// +build windows

package service

import (
	"fmt"
	"strings"
)

func (m *Manager) installWindows() error {
	binPath := fmt.Sprintf("\"%s\" %s", m.config.Executable, m.buildArgs())

	args := []string{
		"create",
		m.config.Name,
		"binPath=", binPath,
		"DisplayName=", m.config.DisplayName,
		"start=", "auto",
	}

	if err := runCommand("sc", args...); err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	// Set description
	if m.config.Description != "" {
		descArgs := []string{
			"description",
			m.config.Name,
			m.config.Description,
		}
		runCommand("sc", descArgs...) // Ignore error as description is optional
	}

	fmt.Printf("Service %s installed successfully\n", m.config.Name)
	return nil
}

func (m *Manager) uninstallWindows() error {
	// Stop service first
	runCommand("sc", "stop", m.config.Name)

	if err := runCommand("sc", "delete", m.config.Name); err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}

	fmt.Printf("Service %s uninstalled successfully\n", m.config.Name)
	return nil
}

func (m *Manager) controlWindows(action string) error {
	var scAction string
	switch action {
	case "start":
		scAction = "start"
	case "stop":
		scAction = "stop"
	case "restart":
		if err := runCommand("sc", "stop", m.config.Name); err != nil {
			return err
		}
		scAction = "start"
	default:
		return fmt.Errorf("unknown action: %s", action)
	}

	return runCommand("sc", scAction, m.config.Name)
}

func (m *Manager) disableWindows() error {
	return runCommand("sc", "config", m.config.Name, "start=", "disabled")
}

func (m *Manager) statusWindows() error {
	return runCommand("sc", "query", m.config.Name)
}
