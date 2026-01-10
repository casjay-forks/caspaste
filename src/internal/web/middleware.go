// Copyright (C) 2021-2023 Leonid Maslakov.

// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE file for details.

package web

import (
	"net/http"
	"os"
)

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
