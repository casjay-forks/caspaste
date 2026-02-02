
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package apiv1

import (
	"net/http"
	"strconv"

	"github.com/casjay-forks/caspaste/src/caspasswd"
	"github.com/casjay-forks/caspaste/src/netshare"
)

type newPasteAnswer struct {
	ID         string `json:"id"`
	URL        string `json:"url"`
	CreateTime int64  `json:"createTime"`
	DeleteTime int64  `json:"deleteTime"`
}

// POST /api/v1/new
func (data *Data) newHand(rw http.ResponseWriter, req *http.Request) error {
	var err error

	// Check auth (required when server.public=false)
	if !data.Public && data.CasPasswdFile != "" {
		clientIP := netshare.GetClientAddr(req)

		// Check if IP is blocked due to too many failed attempts
		if data.BruteForce != nil && data.BruteForce.CheckBlocked(clientIP) {
			// Return 429 Too Many Requests with retry-after header
			remaining := data.BruteForce.GetRemainingLockout(clientIP)
			rw.Header().Set("Retry-After", strconv.Itoa(int(remaining.Seconds())))
			return netshare.ErrTooManyRequests
		}

		isAuthenticated := false

		user, pass, authProvided := req.BasicAuth()
		if authProvided {
			isAuthenticated, err = caspasswd.LoadAndCheck(data.CasPasswdFile, user, pass)
			if err != nil {
				return err
			}
		}

		if !isAuthenticated {
			// Record failed attempt
			if data.BruteForce != nil {
				data.BruteForce.RecordFailure(clientIP)
			}
			return netshare.ErrUnauthorized
		}

		// Record successful login
		if data.BruteForce != nil {
			data.BruteForce.RecordSuccess(clientIP)
		}
	}

	// Check method
	if req.Method != "POST" {
		return netshare.ErrMethodNotAllowed
	}

	// Get form data and create paste
	pasteID, createTime, deleteTime, err := netshare.PasteAddFromForm(req, data.DB, data.RateLimitNew, data.TitleMaxLen, data.BodyMaxLen, data.MaxLifeTime, data.Lexers)
	if err != nil {
		return err
	}

	// Construct full URL for paste
	url := netshare.BuildPasteURL(req, pasteID)

	// Return response per AI.md PART 14 (indented JSON with newline)
	return writeJSON(rw, newPasteAnswer{
		ID:         pasteID,
		URL:        url,
		CreateTime: createTime,
		DeleteTime: deleteTime,
	})
}
