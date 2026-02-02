// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

// Package mode provides application mode detection and management
// per AI.md PART 6
package mode

import (
	"os"
	"runtime"
	"strings"

	"github.com/casjay-forks/caspaste/src/validation"
)

var (
	currentMode  = Production
	debugEnabled = false
)

// AppMode represents the application runtime mode
type AppMode int

const (
	// Production is the default mode with optimizations
	Production AppMode = iota
	// Development mode with verbose logging and debug endpoints
	Development
)

// String returns the string representation of the mode
func (m AppMode) String() string {
	switch m {
	case Development:
		return "development"
	default:
		return "production"
	}
}

// SetAppMode sets the application mode
func SetAppMode(m string) {
	switch strings.ToLower(m) {
	case "dev", "development":
		currentMode = Development
	default:
		currentMode = Production
	}
	updateAppModeProfilingSettings()
}

// Set is an alias for SetAppMode
func Set(m string) {
	SetAppMode(m)
}

// SetDebugEnabled enables or disables debug mode
func SetDebugEnabled(enabled bool) {
	debugEnabled = enabled
	updateAppModeProfilingSettings()
}

// SetDebug is an alias for SetDebugEnabled
func SetDebug(enabled bool) {
	SetDebugEnabled(enabled)
}

// updateAppModeProfilingSettings enables/disables profiling based on debug flag
func updateAppModeProfilingSettings() {
	if debugEnabled {
		// Enable profiling when debug is on
		runtime.SetBlockProfileRate(1)
		runtime.SetMutexProfileFraction(1)
	} else {
		// Disable profiling when debug is off
		runtime.SetBlockProfileRate(0)
		runtime.SetMutexProfileFraction(0)
	}
}

// GetCurrentAppMode returns the current application mode
func GetCurrentAppMode() AppMode {
	return currentMode
}

// IsAppModeDev returns true if in development mode
func IsAppModeDev() bool {
	return currentMode == Development
}

// IsAppModeProd returns true if in production mode
func IsAppModeProd() bool {
	return currentMode == Production
}

// IsDebugEnabled returns true if debug mode is enabled (--debug or DEBUG=true)
func IsDebugEnabled() bool {
	return debugEnabled
}

// GetAppModeString returns mode string with debug suffix if enabled
func GetAppModeString() string {
	s := currentMode.String()
	if debugEnabled {
		s += " [debugging]"
	}
	return s
}

// FromEnv sets mode and debug from environment variables
func FromEnv() {
	if m := os.Getenv("MODE"); m != "" {
		Set(m)
	}
	if validation.IsTruthy(os.Getenv("DEBUG")) {
		SetDebug(true)
	}
}
