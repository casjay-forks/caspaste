// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package validation

import (
	"fmt"
	"strings"
)

// DetectDriver auto-detects database driver from connection string
// Supports: sqlite:///path, postgres://..., mysql://..., mariadb://...
// Returns normalized driver name (mariadb → mysql, sqlite3 → sqlite)
func DetectDriver(source string) (string, error) {
	source = strings.ToLower(strings.TrimSpace(source))

	// SQLite variants
	if strings.HasPrefix(source, "sqlite://") {
		return "sqlite", nil
	}
	if strings.HasPrefix(source, "sqlite3://") {
		return "sqlite", nil
	}

	// PostgreSQL
	if strings.HasPrefix(source, "postgres://") || strings.HasPrefix(source, "postgresql://") {
		return "postgres", nil
	}

	// MySQL
	if strings.HasPrefix(source, "mysql://") {
		return "mysql", nil
	}

	// MariaDB (uses MySQL driver)
	if strings.HasPrefix(source, "mariadb://") {
		return "mysql", nil
	}

	// If no scheme, check if it looks like a file path (SQLite)
	if strings.Contains(source, "/") || strings.HasSuffix(source, ".db") {
		return "sqlite", nil
	}

	return "", fmt.Errorf("could not detect database driver from source: %s", source)
}

// NormalizeDriver normalizes driver names
// mariadb → mysql, sqlite3 → sqlite
func NormalizeDriver(driver string) string {
	driver = strings.ToLower(driver)
	if driver == "mariadb" || driver == "maria" {
		return "mysql"
	}
	if driver == "sqlite3" {
		return "sqlite"
	}
	if driver == "postgresql" {
		return "postgres"
	}
	return driver
}

// NormalizeConnectionString normalizes connection strings
// sqlite:///path → /path
// Other drivers: keep as-is
func NormalizeConnectionString(driver, source string) string {
	// Remove sqlite:// prefix
	if driver == "sqlite" && strings.HasPrefix(source, "sqlite://") {
		return strings.TrimPrefix(source, "sqlite://")
	}
	if driver == "sqlite" && strings.HasPrefix(source, "sqlite3://") {
		return strings.TrimPrefix(source, "sqlite3://")
	}

	return source
}
