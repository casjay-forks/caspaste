// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package web

import (
	"fmt"
	"net/http"
)

// Pattern: /.well-known/security.txt
func (data *Data) securityTxtHand(rw http.ResponseWriter, req *http.Request) error {
	var content string

	// Use override if specified
	if data.SecurityTxt != "" {
		content = data.SecurityTxt
	} else {
		// Auto-generate from config
		content = generateSecurityTxt(data.SecurityContactEmail, data.SecurityContactName, data.FQDN)
	}

	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.Write([]byte(content))
	return nil
}

// generateSecurityTxt creates RFC 9116 compliant security.txt
func generateSecurityTxt(email, name, fqdn string) string {
	canonical := fmt.Sprintf("https://%s/.well-known/security.txt", fqdn)

	return fmt.Sprintf(`Contact: mailto:%s
Preferred-Languages: en
Canonical: %s

# Security Contact: %s
# Please report security vulnerabilities to the contact above.
#
# This file follows RFC 9116: https://www.rfc-editor.org/rfc/rfc9116.html
`, email, canonical, name)
}
