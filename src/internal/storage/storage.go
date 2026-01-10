// Copyright (C) 2021-2023 Leonid Maslakov.

// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE file for details.

package storage

import (
	"database/sql"
	"errors"

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

	// If using postgres/mysql, also open SQLite backup/cache
	if driverName == "postgres" || driverName == "mysql" || driverName == "mariadb" {
		// SQLite backup at standard path (same as if SQLite was primary)
		backupPath := dataDir + "/db/caspaste.db"
		if dataDir == "" {
			backupPath = "./data/db/caspaste.db"
		}

		db.backupPool, err = sql.Open("sqlite3", backupPath)
		if err != nil {
			// Don't fail if backup can't be opened, just log warning
			db.backupPool = nil
		} else {
			db.backupPool.SetMaxOpenConns(10)
			db.backupPool.SetMaxIdleConns(2)
			// Initialize backup database schema
			InitDB("sqlite3", backupPath)
		}
	}

	return db, nil
}

func (db DB) Close() error {
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
	if driverName == "sqlite3" || driverName == "sqlite" {
		// SQLite: ALTER TABLE ADD COLUMN (ignores duplicate errors)
		columns := []string{
			"author TEXT NOT NULL DEFAULT ''",
			"author_email TEXT NOT NULL DEFAULT ''",
			"author_url TEXT NOT NULL DEFAULT ''",
			"is_file BOOL NOT NULL DEFAULT 0",
			"file_name TEXT NOT NULL DEFAULT ''",
			"mime_type TEXT NOT NULL DEFAULT ''",
			"is_editable BOOL NOT NULL DEFAULT 0",
			"is_private BOOL NOT NULL DEFAULT 0",
			"is_url BOOL NOT NULL DEFAULT 0",
			"original_url TEXT NOT NULL DEFAULT ''",
		}
		for _, col := range columns {
			db.pool.Exec(`ALTER TABLE pastes ADD COLUMN ` + col)
		}

	} else if driverName == "mysql" || driverName == "mariadb" {
		// MySQL/MariaDB: Use ALTER TABLE ADD COLUMN IF NOT EXISTS (MariaDB 10.0+)
		columns := []string{
			"author TEXT NOT NULL DEFAULT ''",
			"author_email TEXT NOT NULL DEFAULT ''",
			"author_url TEXT NOT NULL DEFAULT ''",
			"is_file BOOLEAN NOT NULL DEFAULT false",
			"file_name TEXT NOT NULL DEFAULT ''",
			"mime_type TEXT NOT NULL DEFAULT ''",
			"is_editable BOOLEAN NOT NULL DEFAULT false",
			"is_private BOOLEAN NOT NULL DEFAULT false",
			"is_url BOOLEAN NOT NULL DEFAULT false",
			"original_url TEXT NOT NULL DEFAULT ''",
		}
		for _, col := range columns {
			db.pool.Exec(`ALTER TABLE pastes ADD COLUMN IF NOT EXISTS ` + col)
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
