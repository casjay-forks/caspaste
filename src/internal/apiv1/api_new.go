// Copyright (C) 2021-2023 Leonid Maslakov.

// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE file for details.

package apiv1

import (
	"encoding/json"
	"strconv"

	"github.com/casjay-forks/caspaste/src/internal/caspasswd"
	"github.com/casjay-forks/caspaste/src/internal/netshare"
	"net/http"
)

type newPasteAnswer struct {
	ID         string `json:"id"`
	CreateTime int64  `json:"createTime"`
	DeleteTime int64  `json:"deleteTime"`
}

// POST /api/v1/new
func (data *Data) newHand(rw http.ResponseWriter, req *http.Request) error {
	var err error

	// Check auth
	if data.CasPasswdFile != "" {
		clientIP := netshare.GetClientAddr(req)

		// Check if IP is blocked due to too many failed attempts
		if data.BruteForce != nil && data.BruteForce.CheckBlocked(clientIP) {
			// Return 429 Too Many Requests with retry-after header
			remaining := data.BruteForce.GetRemainingLockout(clientIP)
			rw.Header().Set("Retry-After", strconv.Itoa(int(remaining.Seconds())))
			return netshare.ErrTooManyRequests
		}

		authOk := false

		user, pass, authExist := req.BasicAuth()
		if authExist == true {
			authOk, err = caspasswd.LoadAndCheck(data.CasPasswdFile, user, pass)
			if err != nil {
				return err
			}
		}

		if authOk == false {
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

	// Return response
	rw.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(newPasteAnswer{ID: pasteID, CreateTime: createTime, DeleteTime: deleteTime})
}
