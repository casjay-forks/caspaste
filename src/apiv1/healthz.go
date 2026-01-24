// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package apiv1

import (
	"encoding/json"
	"net/http"
	"time"
)

type healthzResponse struct {
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
	Version   string `json:"version"`
	Database  string `json:"database"`
	Uptime    int64  `json:"uptime"`
}

var startTime = time.Now()

// GET /api/healthz
func (data *Data) handleHealthz(rw http.ResponseWriter, req *http.Request) error {
	if req.Method != "GET" {
		rw.Header().Set("Allow", "GET")
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return nil
	}

	resp := healthzResponse{
		Status:    "healthy",
		Timestamp: time.Now().Unix(),
		Version:   data.Version,
		Database:  "connected",
		Uptime:    int64(time.Since(startTime).Seconds()),
	}

	// Try to ping database
	_, err := data.DB.PasteDeleteExpired()
	if err != nil {
		resp.Status = "degraded"
		resp.Database = "error"
	}

	// Set status code and return response per AI.md PART 14 (indented JSON with newline)
	rw.Header().Set("Content-Type", "application/json")
	if err != nil {
		rw.WriteHeader(http.StatusServiceUnavailable)
	} else {
		rw.WriteHeader(http.StatusOK)
	}
	jsonData, _ := json.MarshalIndent(resp, "", "  ")
	rw.Write(jsonData)
	rw.Write([]byte("\n"))
	return nil
}
