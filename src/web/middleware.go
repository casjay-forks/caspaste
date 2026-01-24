
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"runtime"

	"github.com/google/uuid"
)

// SecurityHeadersMiddleware adds security headers to all responses per AI.md PART 11
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

			// XSS Protection (deprecated but kept for older browser compatibility per AI.md)
			if cfg.XSSProtection != "" {
				w.Header().Set("X-XSS-Protection", cfg.XSSProtection)
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

// uuidRegex validates UUID v4 format
var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// isValidUUID checks if string is a valid UUID format
func isValidUUID(s string) bool {
	return uuidRegex.MatchString(s)
}

// RequestIDKey is the context key for request ID
type RequestIDKey struct{}

// RequestIDMiddleware adds a unique request ID to each request per AI.md PART 11
// Every request MUST have a Request ID for tracing and debugging.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for existing request ID from client or upstream proxy
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = r.Header.Get("X-Correlation-ID")
		}
		if requestID == "" {
			requestID = r.Header.Get("X-Trace-ID")
		}

		// Generate new ID if none provided or invalid
		if requestID == "" || !isValidUUID(requestID) {
			requestID = uuid.New().String()
		}

		// Add to response headers
		w.Header().Set("X-Request-ID", requestID)

		// Add to request context for logging and downstream calls
		ctx := context.WithValue(r.Context(), RequestIDKey{}, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestID retrieves the request ID from context
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey{}).(string); ok {
		return id
	}
	return ""
}

// PanicRecoveryMiddleware recovers from panics and returns appropriate error response
// Per AI.md PART 6:
// - Production: Graceful recovery, logs error, returns 500
// - Development: Verbose, full stack in response
func PanicRecoveryMiddleware(debug bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					requestID := GetRequestID(r.Context())

					// Log the panic with stack trace
					stack := make([]byte, 4096)
					n := runtime.Stack(stack, false)
					stack = stack[:n]

					// Log format includes request_id per AI.md
					logMsg := "panic recovered"
					if requestID != "" {
						logMsg += ", request_id=" + requestID
					}
					logMsg += ", error=" + fmt.Sprintf("%v", err)

					if debug {
						// Development: verbose response with stack trace
						w.Header().Set("Content-Type", "text/plain; charset=utf-8")
						w.Header().Set("X-Content-Type-Options", "nosniff")
						w.WriteHeader(http.StatusInternalServerError)
						fmt.Fprintf(w, "Internal Server Error\n\n")
						fmt.Fprintf(w, "Panic: %v\n\n", err)
						fmt.Fprintf(w, "Stack Trace:\n%s\n", stack)
						if requestID != "" {
							fmt.Fprintf(w, "\nRequest ID: %s\n", requestID)
						}
					} else {
						// Production: generic error message
						w.Header().Set("Content-Type", "text/plain; charset=utf-8")
						w.Header().Set("X-Content-Type-Options", "nosniff")
						w.WriteHeader(http.StatusInternalServerError)
						fmt.Fprint(w, "An unexpected error occurred")
					}
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
