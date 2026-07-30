package backup

import (
	"encoding/json"
	"os"
	"time"
)

// Manifest represents a backup manifest
type Manifest struct {
	Version          string    `json:"version"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedBy        string    `json:"created_by"`
	AppVersion       string    `json:"app_version"`
	Contents         []string  `json:"contents"`
	Encrypted        bool      `json:"encrypted"`
	EncryptionMethod string    `json:"encryption_method,omitempty"`
	Checksum         string    `json:"checksum"`
}

// NewManifest creates a new manifest with default values
func NewManifest() *Manifest {
	hostname, _ := os.Hostname()
	return &Manifest{
		Version:   "1.0.0",
		CreatedAt: time.Now().UTC(),
		CreatedBy: hostname,
		Contents:  make([]string, 0),
	}
}

// Save saves the manifest to a file
func (m *Manifest) Save(path string) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// LoadManifest loads a manifest from a file
func LoadManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}

// Validate validates the manifest
func (m *Manifest) Validate() error {
	if m.Version == "" {
		return ErrInvalidManifest
	}
	if m.CreatedAt.IsZero() {
		return ErrInvalidManifest
	}
	if len(m.Contents) == 0 {
		return ErrInvalidManifest
	}
	return nil
}
