package backup

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/argon2"
)

const (
	// Argon2id parameters for key derivation
	argon2Time = 1
	// 64MB
	argon2Memory  = 64 * 1024
	argon2Threads = 4
	// AES-256
	argon2KeyLen = 32

	// Salt size
	saltSize = 16

	// Nonce size for AES-GCM
	nonceSize = 12
)

// EncryptionService handles backup encryption/decryption
type EncryptionService struct{}

// NewEncryptionService creates a new encryption service
func NewEncryptionService() *EncryptionService {
	return &EncryptionService{}
}

// EncryptFile encrypts a file using AES-256-GCM with Argon2id key derivation
func (e *EncryptionService) EncryptFile(inputPath, password string) (string, error) {
	// Read input file
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Generate salt
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive key using Argon2id
	key := argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)

	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate nonce
	nonce := make([]byte, nonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	// Prepare output: salt + ciphertext (nonce is prepended by Seal)
	output := make([]byte, len(salt)+len(ciphertext))
	copy(output[:saltSize], salt)
	copy(output[saltSize:], ciphertext)

	// Write encrypted file
	outputPath := inputPath + ".enc"
	if err := os.WriteFile(outputPath, output, 0600); err != nil {
		return "", fmt.Errorf("failed to write encrypted file: %w", err)
	}

	return outputPath, nil
}

// DecryptFile decrypts an encrypted backup file
func (e *EncryptionService) DecryptFile(inputPath, password, outputDir string) (string, error) {
	// Read encrypted file
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	if len(data) < saltSize+nonceSize {
		return "", ErrInvalidBackupFile
	}

	// Extract salt
	salt := data[:saltSize]
	ciphertext := data[saltSize:]

	// Derive key using Argon2id
	key := argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)

	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Extract nonce
	if len(ciphertext) < nonceSize {
		return "", ErrInvalidBackupFile
	}
	nonce := ciphertext[:nonceSize]
	ciphertext = ciphertext[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", ErrInvalidPassword
	}

	// Write decrypted file
	basename := filepath.Base(inputPath)
	// Remove .enc extension
	if len(basename) > 4 && basename[len(basename)-4:] == ".enc" {
		basename = basename[:len(basename)-4]
	}
	outputPath := filepath.Join(outputDir, basename)

	if err := os.WriteFile(outputPath, plaintext, 0600); err != nil {
		return "", fmt.Errorf("failed to write decrypted file: %w", err)
	}

	return outputPath, nil
}

// EncryptData encrypts data in memory using AES-256-GCM
func (e *EncryptionService) EncryptData(plaintext []byte, password string) ([]byte, error) {
	// Generate salt
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive key using Argon2id
	key := argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)

	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate nonce
	nonce := make([]byte, nonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	// Prepare output: salt + ciphertext
	output := make([]byte, len(salt)+len(ciphertext))
	copy(output[:saltSize], salt)
	copy(output[saltSize:], ciphertext)

	return output, nil
}

// DecryptData decrypts data in memory using AES-256-GCM
func (e *EncryptionService) DecryptData(data []byte, password string) ([]byte, error) {
	if len(data) < saltSize+nonceSize {
		return nil, ErrInvalidBackupFile
	}

	// Extract salt
	salt := data[:saltSize]
	ciphertext := data[saltSize:]

	// Derive key using Argon2id
	key := argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)

	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Extract nonce
	if len(ciphertext) < nonceSize {
		return nil, ErrInvalidBackupFile
	}
	nonce := ciphertext[:nonceSize]
	ciphertext = ciphertext[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, ErrInvalidPassword
	}

	return plaintext, nil
}

// HashPassword hashes a password using SHA-256 (for hint storage, NOT for encryption)
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}

// GenerateRandomBytes generates cryptographically secure random bytes
func GenerateRandomBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return nil, err
	}
	return bytes, nil
}
