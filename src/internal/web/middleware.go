
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package web

import (
	"net/http"
	"os"
)

// SecurityHeadersMiddleware adds security headers to all responses
func SecurityHeadersMiddleware(cfg SecurityHeadersConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Anti-clickjacking
			if cfg.XFrameOptions != "" {
				w.Header().Set("X-Frame-Options", cfg.XFrameOptions)
			}
			
			// Prevent MIME-sniffing
			if cfg.XContentTypeOptions != "" {
				w.Header().Set("X-Content-Type-Options", cfg.XContentTypeOptions)
			}
			
			// Content Security Policy
			if cfg.ContentSecurityPolicy != "" {
				w.Header().Set("Content-Security-Policy", cfg.ContentSecurityPolicy)
			}
			
			// Referrer policy
			if cfg.ReferrerPolicy != "" {
				w.Header().Set("Referrer-Policy", cfg.ReferrerPolicy)
			}
			
			// Permissions policy
			if cfg.PermissionsPolicy != "" {
				w.Header().Set("Permissions-Policy", cfg.PermissionsPolicy)
			}
			
			// HSTS (only if HTTPS)
			if r.TLS != nil && cfg.StrictTransportSecurity != "" {
				w.Header().Set("Strict-Transport-Security", cfg.StrictTransportSecurity)
			}
			
			next.ServeHTTP(w, r)
		})
	}
}

// CORSMiddleware adds CORS headers to all responses
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins (as requested by user)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// MaintenanceMiddleware checks for maintenance mode file
func MaintenanceMiddleware(dataDir string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		maintenanceFile := dataDir + "/.maintenance"

		// Check if maintenance mode file exists
		if _, err := os.Stat(maintenanceFile); err == nil {
			// Maintenance mode is enabled
			w.Header().Set("Content-Type", "text/html; charset=UTF-8")
			w.Header().Set("Retry-After", "3600") // Retry after 1 hour
			w.WriteHeader(http.StatusServiceUnavailable)

			html := `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Maintenance Mode</title>
	<style>
		body { font-family: sans-serif; text-align: center; padding: 50px; }
		h1 { color: #e74c3c; }
	</style>
</head>
<body>
	<h1>503 - Service Unavailable</h1>
	<p>The server is currently in maintenance mode.</p>
	<p>Please try again later.</p>
</body>
</html>`
			w.Write([]byte(html))
			return
		}

		// Not in maintenance mode, continue normally
		next.ServeHTTP(w, r)
	})
}
