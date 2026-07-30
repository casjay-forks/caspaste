package backup

import "errors"

// Backup errors
var (
	// ErrInvalidPassword indicates the provided password is incorrect
	ErrInvalidPassword = errors.New("invalid backup password")

	// ErrInvalidBackupFile indicates the backup file is corrupted or invalid
	ErrInvalidBackupFile = errors.New("invalid backup file format")

	// ErrInvalidManifest indicates the backup manifest is invalid
	ErrInvalidManifest = errors.New("invalid backup manifest")

	// ErrVerificationFailed indicates backup verification failed
	ErrVerificationFailed = errors.New("backup verification failed")

	// ErrBackupNotFound indicates the requested backup was not found
	ErrBackupNotFound = errors.New("backup not found")

	// ErrPasswordRequired indicates a password is required for this operation
	ErrPasswordRequired = errors.New("password required for encrypted backup")

	// ErrComplianceRequiresEncryption indicates compliance mode requires encryption
	ErrComplianceRequiresEncryption = errors.New("compliance mode requires encrypted backups")

	// ErrBackupInProgress indicates a backup is already in progress
	ErrBackupInProgress = errors.New("backup already in progress")

	// ErrRestoreInProgress indicates a restore is already in progress
	ErrRestoreInProgress = errors.New("restore already in progress")

	// ErrVersionMismatch indicates a version mismatch between backup and application
	ErrVersionMismatch = errors.New("backup version mismatch")

	// ErrChecksumMismatch indicates the backup checksum does not match
	ErrChecksumMismatch = errors.New("backup checksum mismatch")

	// ErrDatabaseIntegrity indicates database integrity check failed
	ErrDatabaseIntegrity = errors.New("database integrity check failed")

	// ErrExtractFailed indicates archive extraction failed
	ErrExtractFailed = errors.New("failed to extract backup archive")
)

// IsPasswordError checks if the error is related to password issues
func IsPasswordError(err error) bool {
	return errors.Is(err, ErrInvalidPassword) || errors.Is(err, ErrPasswordRequired)
}

// IsVerificationError checks if the error is related to verification
func IsVerificationError(err error) bool {
	return errors.Is(err, ErrVerificationFailed) ||
		errors.Is(err, ErrChecksumMismatch) ||
		errors.Is(err, ErrDatabaseIntegrity) ||
		errors.Is(err, ErrInvalidManifest)
}
