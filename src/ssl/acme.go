// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

// ACME/Let's Encrypt support per AI.md PART 15
// Provides automatic certificate issuance via HTTP-01 challenge
package ssl

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

// ACMEConfig holds ACME/Let's Encrypt configuration
type ACMEConfig struct {
	Enabled   bool
	Email     string
	CacheDir  string
	Staging   bool
	Domains   []string
	Challenge string
}

// ACMEManager handles automatic certificate management
type ACMEManager struct {
	config      *ACMEConfig
	autocert    *autocert.Manager
	challenges  map[string]string
	challengeMu sync.RWMutex
	enabled     bool
}

// NewACMEManager creates a new ACME manager
func NewACMEManager(cfg *ACMEConfig) (*ACMEManager, error) {
	if cfg == nil || !cfg.Enabled {
		return &ACMEManager{enabled: false}, nil
	}

	// Ensure cache directory exists
	cacheDir := cfg.CacheDir
	if cacheDir == "" {
		cacheDir = filepath.Join(os.TempDir(), "acme-certs")
	}
	if err := os.MkdirAll(cacheDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create ACME cache directory: %w", err)
	}

	// Determine ACME directory URL
	var directoryURL string
	if cfg.Staging {
		directoryURL = "https://acme-staging-v02.api.letsencrypt.org/directory"
	} else {
		directoryURL = autocert.DefaultACMEDirectory
	}

	// Create autocert manager
	m := &autocert.Manager{
		Prompt:      autocert.AcceptTOS,
		Cache:       autocert.DirCache(cacheDir),
		Email:       cfg.Email,
		HostPolicy:  hostWhitelist(cfg.Domains),
		RenewBefore: 30 * 24 * time.Hour,
	}

	// Set staging directory if configured
	if cfg.Staging {
		m.Client = &acme.Client{
			DirectoryURL: directoryURL,
		}
	}

	return &ACMEManager{
		config:     cfg,
		autocert:   m,
		challenges: make(map[string]string),
		enabled:    true,
	}, nil
}

// hostWhitelist creates a host policy that allows specified domains
func hostWhitelist(domains []string) autocert.HostPolicy {
	allowed := make(map[string]bool)
	for _, d := range domains {
		allowed[d] = true
	}
	return func(ctx context.Context, host string) error {
		if len(allowed) == 0 {
			// No whitelist, allow all
			return nil
		}
		if allowed[host] {
			return nil
		}
		return fmt.Errorf("host %q not allowed", host)
	}
}

// IsEnabled returns true if ACME is enabled
func (m *ACMEManager) IsEnabled() bool {
	return m.enabled
}

// GetTLSConfig returns a TLS config with automatic certificate management
func (m *ACMEManager) GetTLSConfig() *tls.Config {
	if !m.enabled || m.autocert == nil {
		return nil
	}

	return &tls.Config{
		GetCertificate: m.autocert.GetCertificate,
		NextProtos:     []string{"h2", "http/1.1", acme.ALPNProto},
		MinVersion:     tls.VersionTLS12,
	}
}

// HTTPHandler returns the HTTP-01 challenge handler
// Mount this at /.well-known/acme-challenge/
func (m *ACMEManager) HTTPHandler(fallback http.Handler) http.Handler {
	if !m.enabled || m.autocert == nil {
		return fallback
	}
	return m.autocert.HTTPHandler(fallback)
}

// SetChallenge stores a challenge token for HTTP-01
func (m *ACMEManager) SetChallenge(domain, token, response string) {
	m.challengeMu.Lock()
	defer m.challengeMu.Unlock()
	key := domain + "/" + token
	m.challenges[key] = response
}

// GetChallenge retrieves a challenge response
func (m *ACMEManager) GetChallenge(domain, token string) (string, bool) {
	m.challengeMu.RLock()
	defer m.challengeMu.RUnlock()
	key := domain + "/" + token
	response, ok := m.challenges[key]
	return response, ok
}

// ClearChallenge removes a challenge token
func (m *ACMEManager) ClearChallenge(domain, token string) {
	m.challengeMu.Lock()
	defer m.challengeMu.Unlock()
	key := domain + "/" + token
	delete(m.challenges, key)
}

// HandleChallenge is an HTTP handler for ACME HTTP-01 challenges
// Route: GET /.well-known/acme-challenge/{token}
func (m *ACMEManager) HandleChallenge(w http.ResponseWriter, r *http.Request) {
	// Extract token from path
	token := filepath.Base(r.URL.Path)
	host := r.Host

	// Strip port if present
	if idx := len(host) - 1; idx > 0 {
		for i := idx; i >= 0; i-- {
			if host[i] == ':' {
				host = host[:i]
				break
			}
		}
	}

	response, ok := m.GetChallenge(host, token)
	if !ok {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(response))
}

// CertificateStatus represents the status of a certificate
type CertificateStatus struct {
	Domain    string
	Valid     bool
	ExpiresAt time.Time
	Issuer    string
	Error     string
}

// GetCertificateStatus checks the status of a certificate for a domain
func (m *ACMEManager) GetCertificateStatus(domain string) *CertificateStatus {
	status := &CertificateStatus{
		Domain: domain,
	}

	if !m.enabled || m.autocert == nil {
		status.Error = "ACME not enabled"
		return status
	}

	// Try to get certificate from cache
	cert, err := m.autocert.GetCertificate(&tls.ClientHelloInfo{
		ServerName: domain,
	})
	if err != nil {
		status.Error = err.Error()
		return status
	}

	if cert != nil && len(cert.Certificate) > 0 {
		leaf, err := x509.ParseCertificate(cert.Certificate[0])
		if err == nil {
			status.Valid = time.Now().Before(leaf.NotAfter)
			status.ExpiresAt = leaf.NotAfter
			if len(leaf.Issuer.Organization) > 0 {
				status.Issuer = leaf.Issuer.Organization[0]
			}
		}
	}

	return status
}
