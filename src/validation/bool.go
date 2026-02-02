// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package validation

import (
	"strings"
)

// ParseBool parses a boolean value with support for various truthy/falsey strings
// Truthy:  yes, true, enable, enabled, 1, on
// Falsey:  no, false, disable, disabled, 0, off, (empty string)
// Returns: (value, wasSet)
func ParseBool(s string) (bool, bool) {
	if s == "" {
		return false, false
	}

	s = strings.ToLower(strings.TrimSpace(s))

	// Truthy values
	truthy := []string{"yes", "true", "enable", "enabled", "1", "on", "y", "t"}
	for _, val := range truthy {
		if s == val {
			return true, true
		}
	}

	// Falsey values
	falsey := []string{"no", "false", "disable", "disabled", "0", "off", "n", "f"}
	for _, val := range falsey {
		if s == val {
			return false, true
		}
	}

	// Unknown value, treat as not set
	return false, false
}

// IsTruthy returns true if the string represents a truthy value
func IsTruthy(s string) bool {
	value, wasSet := ParseBool(s)
	return wasSet && value
}

// IsFalsey returns true if the string represents a falsey value
func IsFalsey(s string) bool {
	value, wasSet := ParseBool(s)
	return wasSet && !value
}
