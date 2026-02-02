// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package validation

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

// ValidateFQDN validates that a string is a valid Fully Qualified Domain Name
// Returns error if invalid
func ValidateFQDN(fqdn string) error {
	if fqdn == "" {
		return fmt.Errorf("FQDN is empty")
	}

	// Must have at least one dot (domain.tld at minimum)
	if !strings.Contains(fqdn, ".") {
		return fmt.Errorf("not a valid FQDN: %s (must contain domain)", fqdn)
	}

	// Remove port if present
	if strings.Contains(fqdn, ":") {
		parts := strings.Split(fqdn, ":")
		fqdn = parts[0]
	}

	// Basic DNS validation (alphanumeric, dots, hyphens)
	// Pattern: label.label.tld where each label is alphanumeric/hyphen, max 63 chars
	dnsPattern := `^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`
	if !regexp.MustCompile(dnsPattern).MatchString(fqdn) {
		return fmt.Errorf("invalid FQDN format: %s", fqdn)
	}

	// No trailing or leading dots
	if strings.HasPrefix(fqdn, ".") || strings.HasSuffix(fqdn, ".") {
		return fmt.Errorf("FQDN cannot start or end with dot: %s", fqdn)
	}

	// Not localhost or variations
	lower := strings.ToLower(fqdn)
	forbiddenHosts := []string{
		"localhost",
		"localhost.localdomain",
		"localdomain",
		"local",
		"localhost.local",
		"localhost6",
		"localhost6.localdomain6",
		"localdomain6",
	}
	for _, forbidden := range forbiddenHosts {
		if lower == forbidden || strings.HasSuffix(lower, "."+forbidden) {
			return fmt.Errorf("localhost not allowed as FQDN: %s", fqdn)
		}
	}

	// Not IP address (basic check - IPs have no letters except for IPv6)
	if regexp.MustCompile(`^[0-9.:]+$`).MatchString(fqdn) {
		return fmt.Errorf("IP address not allowed as FQDN: %s", fqdn)
	}

	return nil
}

// DetermineFQDN determines the server's FQDN or IP address
// Priority: Config (override if set) > Reverse proxy header > OS hostname > Global IP
// NEVER returns localhost - uses global IP as fallback
func DetermineFQDN(fromProxy, fromConfig string) (string, error) {
	// Priority 1: Config server.fqdn (OVERRIDE - use if explicitly set)
	if fromConfig != "" {
		// Remove port if present
		if strings.Contains(fromConfig, ":") {
			parts := strings.Split(fromConfig, ":")
			fromConfig = parts[0]
		}

		if strings.Contains(fromConfig, ".") {
			if err := ValidateFQDN(fromConfig); err == nil {
				return fromConfig, nil
			}
		}
	}

	// Priority 2: Reverse proxy header (from X-Forwarded-Host, Forwarded, etc.)
	if fromProxy != "" {
		// Remove port if present
		if strings.Contains(fromProxy, ":") {
			parts := strings.Split(fromProxy, ":")
			fromProxy = parts[0]
		}

		if strings.Contains(fromProxy, ".") {
			if err := ValidateFQDN(fromProxy); err == nil {
				return fromProxy, nil
			}
		}
	}

	// Priority 3: Get hostname from OS
	hostname, err := os.Hostname()
	if err == nil && strings.Contains(hostname, ".") {
		if err := ValidateFQDN(hostname); err == nil {
			return hostname, nil
		}
	}

	// Priority 4: NEVER use localhost - use global IP as fallback
	return GetGlobalIP()
}

// GetGlobalIP returns the first non-loopback global IP address
// Prefers IPv4, falls back to IPv6
func GetGlobalIP() (string, error) {
	// Use UDP dial to determine the IP associated with the default route
	// This doesn't actually connect, just determines which interface would be used
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", fmt.Errorf("failed to determine default route IP: %w", err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

// ValidateEmail validates an email address format
// Checks basic RFC format and that domain is a valid FQDN
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is empty")
	}

	// Basic format: local@domain
	if !strings.Contains(email, "@") {
		return fmt.Errorf("invalid email: %s (missing @)", email)
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return fmt.Errorf("invalid email format: %s (too many @ symbols)", email)
	}

	local := parts[0]
	domain := parts[1]

	if local == "" {
		return fmt.Errorf("invalid email: %s (empty local part)", email)
	}

	if domain == "" {
		return fmt.Errorf("invalid email: %s (empty domain part)", email)
	}

	// Validate domain part is FQDN
	return ValidateFQDN(domain)
}
