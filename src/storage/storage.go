
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

var (
	ErrNotFoundID = errors.New("db: could not find ID")
)

type DB struct {
	pool       *sql.DB
	backupPool *sql.DB // SQLite backup/cache when using postgres/mysql
	driver     string
}

func NewPool(driverName string, dataSourceName string, maxOpenConns int, maxIdleConns int, dataDir string) (DB, error) {
	var db DB
	var err error

	db.driver = driverName
	db.pool, err = sql.Open(driverName, dataSourceName)
	if err != nil {
		return db, err
	}

	db.pool.SetMaxOpenConns(maxOpenConns)
	db.pool.SetMaxIdleConns(maxIdleConns)

	// Set connection lifetime and idle timeouts to prevent stale connections
	db.pool.SetConnMaxLifetime(3600 * 1000000000) // 1 hour in nanoseconds
	db.pool.SetConnMaxIdleTime(600 * 1000000000)  // 10 minutes in nanoseconds

	// If using postgres/mysql, also open SQLite backup/cache
	// SQLite cache is ALWAYS required for local operations
	if driverName == "postgres" || driverName == "mysql" || driverName == "mariadb" {
		// Determine SQLite cache path - check env var first, then use standard path
		backupPath := getSQLiteCachePath(dataDir)

		db.backupPool, err = sql.Open("sqlite", backupPath)
		if err != nil {
			// Don't fail if backup can't be opened, just log warning
			db.backupPool = nil
		} else {
			db.backupPool.SetMaxOpenConns(10)
			db.backupPool.SetMaxIdleConns(2)
			db.backupPool.SetConnMaxLifetime(3600 * 1000000000)
			db.backupPool.SetConnMaxIdleTime(600 * 1000000000)
			// Initialize backup database schema
			InitDB("sqlite", backupPath)
		}
	}

	return db, nil
}

// getSQLiteCachePath determines the SQLite cache database path
// Priority: CASPASTE_DB_DIR env var > dataDir/db/ > platform-specific default
func getSQLiteCachePath(dataDir string) string {
	// Check environment variable first
	if envDbDir := os.Getenv("CASPASTE_DB_DIR"); envDbDir != "" {
		return envDbDir + "/caspaste.db"
	}
	// Use data directory if provided
	if dataDir != "" {
		return dataDir + "/db/caspaste.db"
	}
	// Fallback to platform-specific default
	return getDefaultDbPath()
}

// getDefaultDbPath returns the platform-specific default database path
func getDefaultDbPath() string {
	switch runtime.GOOS {
	case "windows":
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			return localAppData + "\\CasPaste\\Data\\db\\caspaste.db"
		}
		return os.Getenv("PROGRAMDATA") + "\\CasPaste\\Data\\db\\caspaste.db"
	case "darwin":
		if isRunningAsRoot() {
			return "/var/lib/casjay-forks/caspaste/db/caspaste.db"
		}
		if home := os.Getenv("HOME"); home != "" {
			return home + "/Library/Application Support/CasPaste/db/caspaste.db"
		}
		return "/var/lib/casjay-forks/caspaste/db/caspaste.db"
	// Linux, BSD, etc.
	default:
		if isRunningAsRoot() {
			return "/var/lib/casjay-forks/caspaste/db/caspaste.db"
		}
		if home := os.Getenv("HOME"); home != "" {
			return home + "/.local/share/casjay-forks/caspaste/db/caspaste.db"
		}
		return "/var/lib/casjay-forks/caspaste/db/caspaste.db"
	}
}

// isRunningAsRoot checks if the process is running with root/admin privileges
func isRunningAsRoot() bool {
	return os.Geteuid() == 0
}

func (db DB) Close() error {
	// Close backup pool first if it exists
	if db.backupPool != nil {
		if err := db.backupPool.Close(); err != nil {
			// Log but don't fail on backup close error
			// Continue to close primary pool
		}
	}
	return db.pool.Close()
}

func InitDB(driverName string, dataSourceName string) error {
	// Open DB
	db, err := NewPool(driverName, dataSourceName, 1, 0, "")
	if err != nil {
		return err
	}
	defer db.Close()

	// Create tables
	_, err = db.pool.Exec(`
		CREATE TABLE IF NOT EXISTS pastes (
			id          TEXT    PRIMARY KEY,
			title       TEXT    NOT NULL,
			body        TEXT    NOT NULL,
			syntax      TEXT    NOT NULL,
			create_time INTEGER NOT NULL,
			delete_time INTEGER NOT NULL,
			one_use     BOOL    NOT NULL
		);
	`)
	if err != nil {
		return err
	}

	// Handle database-specific column additions
	// Define allowed columns with validation (prevents SQL injection)
	type columnDef struct {
		name       string
		definition string
	}
	
	var columns []columnDef
	if driverName == "sqlite3" || driverName == "sqlite" {
		// SQLite: ALTER TABLE ADD COLUMN (ignores duplicate errors)
		columns = []columnDef{
			{"author", "TEXT NOT NULL DEFAULT ''"},
			{"author_email", "TEXT NOT NULL DEFAULT ''"},
			{"author_url", "TEXT NOT NULL DEFAULT ''"},
			{"is_file", "BOOL NOT NULL DEFAULT 0"},
			{"file_name", "TEXT NOT NULL DEFAULT ''"},
			{"mime_type", "TEXT NOT NULL DEFAULT ''"},
			{"is_editable", "BOOL NOT NULL DEFAULT 0"},
			{"is_private", "BOOL NOT NULL DEFAULT 0"},
			{"is_url", "BOOL NOT NULL DEFAULT 0"},
			{"original_url", "TEXT NOT NULL DEFAULT ''"},
		}
		for _, col := range columns {
			// Using string formatting is safe here because column name is from hardcoded whitelist
			_, err := db.pool.Exec(fmt.Sprintf(`ALTER TABLE pastes ADD COLUMN %s %s`, col.name, col.definition))
			// Ignore "duplicate column" errors
			if err != nil && !strings.Contains(err.Error(), "duplicate column") {
				return err
			}
		}

	} else if driverName == "mysql" || driverName == "mariadb" {
		// MySQL/MariaDB: Use ALTER TABLE ADD COLUMN IF NOT EXISTS (MariaDB 10.0+)
		columns = []columnDef{
			{"author", "TEXT NOT NULL DEFAULT ''"},
			{"author_email", "TEXT NOT NULL DEFAULT ''"},
			{"author_url", "TEXT NOT NULL DEFAULT ''"},
			{"is_file", "BOOLEAN NOT NULL DEFAULT false"},
			{"file_name", "TEXT NOT NULL DEFAULT ''"},
			{"mime_type", "TEXT NOT NULL DEFAULT ''"},
			{"is_editable", "BOOLEAN NOT NULL DEFAULT false"},
			{"is_private", "BOOLEAN NOT NULL DEFAULT false"},
			{"is_url", "BOOLEAN NOT NULL DEFAULT false"},
			{"original_url", "TEXT NOT NULL DEFAULT ''"},
		}
		for _, col := range columns {
			// Using string formatting is safe here because column name is from hardcoded whitelist
			_, err := db.pool.Exec(fmt.Sprintf(`ALTER TABLE pastes ADD COLUMN IF NOT EXISTS %s %s`, col.name, col.definition))
			if err != nil {
				return err
			}
		}

	} else {
		// PostgreSQL: supports IF NOT EXISTS
		_, err = db.pool.Exec(`
			ALTER TABLE pastes ADD COLUMN IF NOT EXISTS author       TEXT NOT NULL DEFAULT '';
			ALTER TABLE pastes ADD COLUMN IF NOT EXISTS author_email TEXT NOT NULL DEFAULT '';
			ALTER TABLE pastes ADD COLUMN IF NOT EXISTS author_url   TEXT NOT NULL DEFAULT '';
			ALTER TABLE pastes ADD COLUMN IF NOT EXISTS is_file      BOOL NOT NULL DEFAULT false;
			ALTER TABLE pastes ADD COLUMN IF NOT EXISTS file_name    TEXT NOT NULL DEFAULT '';
			ALTER TABLE pastes ADD COLUMN IF NOT EXISTS mime_type    TEXT NOT NULL DEFAULT '';
			ALTER TABLE pastes ADD COLUMN IF NOT EXISTS is_editable  BOOL NOT NULL DEFAULT false;
			ALTER TABLE pastes ADD COLUMN IF NOT EXISTS is_private   BOOL NOT NULL DEFAULT false;
			ALTER TABLE pastes ADD COLUMN IF NOT EXISTS is_url       BOOL NOT NULL DEFAULT false;
			ALTER TABLE pastes ADD COLUMN IF NOT EXISTS original_url TEXT NOT NULL DEFAULT '';
		`)
		if err != nil {
			return err
		}
	}

	return nil
}
