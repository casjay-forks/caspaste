
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package recovery

import (
	"database/sql"
	"errors"
	"time"

	"github.com/casjay-forks/caspaste/src/totp"
)

// Default number of recovery keys to generate
const DefaultKeyCount = 10

// Common errors
var (
	ErrKeyNotFound      = errors.New("recovery key not found")
	ErrKeyAlreadyUsed   = errors.New("recovery key has already been used")
	ErrNoKeysRemaining  = errors.New("no recovery keys remaining")
	ErrInvalidKeyFormat = errors.New("invalid recovery key format")
)

// RecoveryKey represents a recovery key
type RecoveryKey struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	KeyHash   string `json:"-"`
	UsedAt    *int64 `json:"used_at,omitempty"`
	CreatedAt int64  `json:"created_at"`
}

// Service provides recovery key operations
type Service struct {
	db *sql.DB
}

// NewService creates a new recovery service
func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

// GenerateKeys generates new recovery keys for a user
// This invalidates all existing keys
func (s *Service) GenerateKeys(userID int64) ([]string, error) {
	// Delete existing keys
	_, err := s.db.Exec("DELETE FROM recovery_keys WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}

	// Generate new keys
	keys, err := totp.GenerateRecoveryKeys(DefaultKeyCount)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()

	// Store hashed keys
	for _, key := range keys {
		keyHash := totp.HashRecoveryKey(key)
		_, err = s.db.Exec(`
			INSERT INTO recovery_keys (user_id, key_hash, created_at)
			VALUES (?, ?, ?)
		`, userID, keyHash, now)
		if err != nil {
			return nil, err
		}
	}

	return keys, nil
}

// VerifyAndConsumeKey verifies a recovery key and marks it as used
func (s *Service) VerifyAndConsumeKey(userID int64, key string) error {
	// Validate format
	if !totp.VerifyRecoveryKeyFormat(key) {
		return ErrInvalidKeyFormat
	}

	// Hash the key
	keyHash := totp.HashRecoveryKey(totp.NormalizeRecoveryKey(key))

	// Find the key
	var keyID int64
	var usedAt sql.NullInt64

	err := s.db.QueryRow(`
		SELECT id, used_at FROM recovery_keys
		WHERE user_id = ? AND key_hash = ?
	`, userID, keyHash).Scan(&keyID, &usedAt)

	if err == sql.ErrNoRows {
		return ErrKeyNotFound
	}
	if err != nil {
		return err
	}

	// Check if already used
	if usedAt.Valid {
		return ErrKeyAlreadyUsed
	}

	// Mark as used
	_, err = s.db.Exec("UPDATE recovery_keys SET used_at = ? WHERE id = ?",
		time.Now().Unix(), keyID)
	if err != nil {
		return err
	}

	return nil
}

// CountRemainingKeys returns the number of unused recovery keys
func (s *Service) CountRemainingKeys(userID int64) (int, error) {
	var count int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM recovery_keys
		WHERE user_id = ? AND used_at IS NULL
	`, userID).Scan(&count)
	return count, err
}

// HasKeys checks if a user has any recovery keys (used or unused)
func (s *Service) HasKeys(userID int64) (bool, error) {
	var count int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM recovery_keys WHERE user_id = ?
	`, userID).Scan(&count)
	return count > 0, err
}

// DeleteAllKeys deletes all recovery keys for a user
func (s *Service) DeleteAllKeys(userID int64) error {
	_, err := s.db.Exec("DELETE FROM recovery_keys WHERE user_id = ?", userID)
	return err
}

// ListKeys returns all recovery keys for a user (showing which are used)
// Note: This does NOT return the actual key values (only hashed)
func (s *Service) ListKeys(userID int64) ([]RecoveryKey, error) {
	rows, err := s.db.Query(`
		SELECT id, user_id, used_at, created_at
		FROM recovery_keys WHERE user_id = ?
		ORDER BY id
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []RecoveryKey
	for rows.Next() {
		var k RecoveryKey
		var usedAt sql.NullInt64

		err := rows.Scan(&k.ID, &k.UserID, &usedAt, &k.CreatedAt)
		if err != nil {
			return nil, err
		}

		if usedAt.Valid {
			k.UsedAt = &usedAt.Int64
		}

		keys = append(keys, k)
	}

	return keys, nil
}

// KeysStatus represents the status of a user's recovery keys
type KeysStatus struct {
	Total     int  `json:"total"`
	Used      int  `json:"used"`
	Remaining int  `json:"remaining"`
	HasKeys   bool `json:"has_keys"`
}

// GetKeysStatus returns the status of a user's recovery keys
func (s *Service) GetKeysStatus(userID int64) (*KeysStatus, error) {
	var total, used int

	err := s.db.QueryRow(`
		SELECT COUNT(*), SUM(CASE WHEN used_at IS NOT NULL THEN 1 ELSE 0 END)
		FROM recovery_keys WHERE user_id = ?
	`, userID).Scan(&total, &used)
	if err != nil {
		return nil, err
	}

	return &KeysStatus{
		Total:     total,
		Used:      used,
		Remaining: total - used,
		HasKeys:   total > 0,
	}, nil
}
