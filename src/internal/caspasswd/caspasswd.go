// Copyright (C) 2021-2023 Leonid Maslakov.

// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE file for details.

package caspasswd

import (
	"bytes"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
)

type Data map[string]string

func LoadFile(path string) (Data, error) {
	// Open file
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.New("caspasswd: " + err.Error())
	}
	defer file.Close()

	// Read file
	fileByte, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.New("caspasswd: " + err.Error())
	}

	// Convert []byte to string
	fileTxt := bytes.NewBuffer(fileByte).String()

	// Parse file
	data := make(Data)
	for i, line := range strings.Split(fileTxt, "\n") {
		if line == "" {
			continue
		}

		lineSplit := strings.Split(line, ":")
		if len(lineSplit) != 2 {
			return nil, errors.New("caspasswd: error in line " + strconv.Itoa(i))
		}

		user := lineSplit[0]
		pass := lineSplit[1]

		_, exist := data[user]
		if exist == true {
			return nil, errors.New("caspasswd: overriding user " + user + " in line " + strconv.Itoa(i))
		}

		data[user] = pass
	}

	return data, nil
}

// Argon2id parameters (recommended by OWASP)
const (
	argon2Time    = 3
	argon2Memory  = 64 * 1024 // 64 MB
	argon2Threads = 4
	argon2KeyLen  = 32
	argon2SaltLen = 16
)

func (data Data) Check(user string, pass string) bool {
	storedPass, exist := data[user]
	if exist == false {
		return false
	}

	// Check if the stored password is an argon2id hash
	// Format: $argon2id$v=19$m=65536,t=3,p=4$salt$hash
	if strings.HasPrefix(storedPass, "$argon2id$") {
		return verifyArgon2Hash(storedPass, pass)
	}

	// Check if the stored password is a bcrypt hash (legacy support)
	// Bcrypt hashes start with $2a$, $2b$, or $2y$
	if strings.HasPrefix(storedPass, "$2a$") ||
	   strings.HasPrefix(storedPass, "$2b$") ||
	   strings.HasPrefix(storedPass, "$2y$") {
		// Use bcrypt comparison
		err := bcrypt.CompareHashAndPassword([]byte(storedPass), []byte(pass))
		return err == nil
	}

	// Legacy plain text password (INSECURE - deprecated)
	// This is kept for backward compatibility only
	// TODO: Remove in future versions and require argon2id hashes
	if pass != storedPass {
		return false
	}

	return true
}

// HashPassword generates an argon2id hash from a plain text password
// Use this to create password hashes for the caspasswd file
// Returns hash in format: $argon2id$v=19$m=65536,t=3,p=4$salt$hash
func HashPassword(password string) (string, error) {
	// Generate random salt
	salt := make([]byte, argon2SaltLen)
	_, err := cryptoRandRead(salt)
	if err != nil {
		return "", err
	}

	// Generate the hash
	hash := argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)

	// Encode to base64
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Format: $argon2id$v=19$m=65536,t=3,p=4$salt$hash
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, argon2Memory, argon2Time, argon2Threads, b64Salt, b64Hash), nil
}

// verifyArgon2Hash verifies an argon2id hash
func verifyArgon2Hash(encodedHash, password string) bool {
	// Parse the hash: $argon2id$v=19$m=65536,t=3,p=4$salt$hash
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return false
	}

	var memory, time uint32
	var threads uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads); err != nil {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}

	// Generate hash with provided password
	hash := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(expectedHash)))

	// Use constant-time comparison to prevent timing attacks
	return subtle.ConstantTimeCompare(hash, expectedHash) == 1
}

// cryptoRandRead is a helper to read random bytes using crypto/rand
func cryptoRandRead(b []byte) (int, error) {
	return rand.Read(b)
}

func LoadAndCheck(path string, user string, pass string) (bool, error) {
	data, err := LoadFile(path)
	if err != nil {
		return false, err
	}

	return data.Check(user, pass), nil
}
