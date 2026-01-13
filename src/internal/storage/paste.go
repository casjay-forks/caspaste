
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package storage

import (
	"database/sql"
	"time"
)

type Paste struct {
	ID         string `json:"id"` // Ignored when creating
	Title      string `json:"title"`
	Body       string `json:"body"`
	CreateTime int64  `json:"createTime"` // Ignored when creating
	DeleteTime int64  `json:"deleteTime"`
	OneUse     bool   `json:"oneUse"`
	Syntax     string `json:"syntax"`

	Author      string `json:"author"`
	AuthorEmail string `json:"authorEmail"`
	AuthorURL   string `json:"authorURL"`

	// MicroBin-inspired features
	IsFile      bool   `json:"isFile"`      // True if this is a file upload
	FileName    string `json:"fileName"`    // Original filename for file uploads
	MimeType    string `json:"mimeType"`    // MIME type for file uploads
	IsEditable  bool   `json:"isEditable"`  // Allow paste editing
	IsPrivate   bool   `json:"isPrivate"`   // Private paste (not listed publicly)
	IsURL       bool   `json:"isURL"`       // True if this is a URL shortener entry
	OriginalURL string `json:"originalURL"` // Original URL for shortener
}

func (db DB) PasteAdd(paste Paste) (string, int64, int64, error) {
	var err error

	// Generate ID
	paste.ID, err = genTokenCrypto(8)
	if err != nil {
		return paste.ID, paste.CreateTime, paste.DeleteTime, err
	}

	// Set paste create time
	paste.CreateTime = time.Now().Unix()

	// Check delete time
	if paste.DeleteTime < 0 {
		paste.DeleteTime = 0
	}

	// Add to primary database
	_, err = db.pool.Exec(
		`INSERT INTO pastes (id, title, body, syntax, create_time, delete_time, one_use, author, author_email, author_url, is_file, file_name, mime_type, is_editable, is_private, is_url, original_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`,
		paste.ID, paste.Title, paste.Body, paste.Syntax, paste.CreateTime, paste.DeleteTime, paste.OneUse,
		paste.Author, paste.AuthorEmail, paste.AuthorURL,
		paste.IsFile, paste.FileName, paste.MimeType, paste.IsEditable, paste.IsPrivate, paste.IsURL, paste.OriginalURL,
	)
	if err != nil {
		return paste.ID, paste.CreateTime, paste.DeleteTime, err
	}

	// Also add to SQLite backup/cache if available
	if db.backupPool != nil {
		_, backupErr := db.backupPool.Exec(
			`INSERT OR REPLACE INTO pastes (id, title, body, syntax, create_time, delete_time, one_use, author, author_email, author_url, is_file, file_name, mime_type, is_editable, is_private, is_url, original_url)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			paste.ID, paste.Title, paste.Body, paste.Syntax, paste.CreateTime, paste.DeleteTime, paste.OneUse,
			paste.Author, paste.AuthorEmail, paste.AuthorURL,
			paste.IsFile, paste.FileName, paste.MimeType, paste.IsEditable, paste.IsPrivate, paste.IsURL, paste.OriginalURL,
		)
		// Log backup errors but don't fail primary operation
		if backupErr != nil {
			// TODO: Log this error when logger is available in this context
			_ = backupErr
		}
	}

	return paste.ID, paste.CreateTime, paste.DeleteTime, nil
}

func (db DB) PasteUpdate(paste Paste) error {
	// Update in primary database
	result, err := db.pool.Exec(
		`UPDATE pastes SET title = $2, body = $3, syntax = $4, delete_time = $5, one_use = $6,
		author = $7, author_email = $8, author_url = $9,
		is_file = $10, file_name = $11, mime_type = $12, is_editable = $13, is_private = $14, is_url = $15, original_url = $16
		WHERE id = $1`,
		paste.ID, paste.Title, paste.Body, paste.Syntax, paste.DeleteTime, paste.OneUse,
		paste.Author, paste.AuthorEmail, paste.AuthorURL,
		paste.IsFile, paste.FileName, paste.MimeType, paste.IsEditable, paste.IsPrivate, paste.IsURL, paste.OriginalURL,
	)
	if err != nil {
		return err
	}

	// Check result
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFoundID
	}

	// Also update in SQLite backup/cache if available
	if db.backupPool != nil {
		_, backupErr := db.backupPool.Exec(
			`UPDATE pastes SET title = ?, body = ?, syntax = ?, delete_time = ?, one_use = ?,
			author = ?, author_email = ?, author_url = ?,
			is_file = ?, file_name = ?, mime_type = ?, is_editable = ?, is_private = ?, is_url = ?, original_url = ?
			WHERE id = ?`,
			paste.Title, paste.Body, paste.Syntax, paste.DeleteTime, paste.OneUse,
			paste.Author, paste.AuthorEmail, paste.AuthorURL,
			paste.IsFile, paste.FileName, paste.MimeType, paste.IsEditable, paste.IsPrivate, paste.IsURL, paste.OriginalURL,
			paste.ID,
		)
		if backupErr != nil {
			_ = backupErr
		}
	}

	return nil
}

func (db DB) PasteDelete(id string) error {
	// Delete from primary database
	result, err := db.pool.Exec(
		`DELETE FROM pastes WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	// Check result
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFoundID
	}

	// Also delete from SQLite backup/cache if available
	if db.backupPool != nil {
		_, backupErr := db.backupPool.Exec(`DELETE FROM pastes WHERE id = ?`, id)
		if backupErr != nil {
			_ = backupErr
		}
	}

	return nil
}

func (db DB) PasteGet(id string) (Paste, error) {
	var paste Paste

	// Make query
	row := db.pool.QueryRow(
		`SELECT id, title, body, syntax, create_time, delete_time, one_use, author, author_email, author_url,
		is_file, file_name, mime_type, is_editable, is_private, is_url, original_url
		FROM pastes WHERE id = $1`,
		id,
	)

	// Read query
	err := row.Scan(&paste.ID, &paste.Title, &paste.Body, &paste.Syntax, &paste.CreateTime, &paste.DeleteTime, &paste.OneUse,
		&paste.Author, &paste.AuthorEmail, &paste.AuthorURL,
		&paste.IsFile, &paste.FileName, &paste.MimeType, &paste.IsEditable, &paste.IsPrivate, &paste.IsURL, &paste.OriginalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return paste, ErrNotFoundID
		}

		return paste, err
	}

	// Check paste expiration
	if paste.DeleteTime < time.Now().Unix() && paste.DeleteTime > 0 {
		// Delete expired paste
		_, err = db.pool.Exec(
			`DELETE FROM pastes WHERE id = $1`,
			paste.ID,
		)
		if err != nil {
			return Paste{}, err
		}

		// Return ErrNotFound
		return Paste{}, ErrNotFoundID
	}

	return paste, nil
}

func (db DB) PasteDeleteExpired() (int64, error) {
	// Delete from primary database
	result, err := db.pool.Exec(
		`DELETE FROM pastes WHERE (delete_time < $1) AND (delete_time > 0)`,
		time.Now().Unix(),
	)
	if err != nil {
		return 0, err
	}

	// Check result
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return rowsAffected, err
	}

	// Also delete from SQLite backup/cache if available
	if db.backupPool != nil {
		_, backupErr := db.backupPool.Exec(
			`DELETE FROM pastes WHERE (delete_time < ?) AND (delete_time > 0)`,
			time.Now().Unix(),
		)
		if backupErr != nil {
			_ = backupErr
		}
	}

	return rowsAffected, nil
}

type PasteListItem struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Syntax     string `json:"syntax"`
	CreateTime int64  `json:"createTime"`
	DeleteTime int64  `json:"deleteTime"`
}

func (db DB) PasteList(limit int, offset int) ([]PasteListItem, error) {
	if limit <= 0 || limit > 100 {
		limit = 50 // Default limit
	}
	if offset < 0 {
		offset = 0
	}

	// Query pastes (exclude expired, one-use, and private pastes)
	rows, err := db.pool.Query(
		`SELECT id, title, syntax, create_time, delete_time
		FROM pastes
		WHERE (delete_time > $1 OR delete_time = 0)
		AND is_private = false
		ORDER BY create_time DESC
		LIMIT $2 OFFSET $3`,
		time.Now().Unix(),
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pastes []PasteListItem
	for rows.Next() {
		var paste PasteListItem
		err := rows.Scan(&paste.ID, &paste.Title, &paste.Syntax, &paste.CreateTime, &paste.DeleteTime)
		if err != nil {
			return nil, err
		}
		pastes = append(pastes, paste)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pastes, nil
}
