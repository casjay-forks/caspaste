
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package apiv1

import (
	"encoding/json"
	"github.com/casjay-forks/caspaste/src/internal/netshare"
	"github.com/casjay-forks/caspaste/src/internal/storage"
	"net/http"
)

// GET /api/v1/get
func (data *Data) getHand(rw http.ResponseWriter, req *http.Request) error {
	// Check rate limit
	err := data.RateLimitGet.CheckAndUse(netshare.GetClientAddr(req))
	if err != nil {
		return err
	}

	// Check method
	if req.Method != "GET" {
		return netshare.ErrMethodNotAllowed
	}

	// Get paste ID
	req.ParseForm()

	pasteID := req.Form.Get("id")

	// Check paste id
	if pasteID == "" {
		return netshare.ErrBadRequest
	}

	// Get paste
	paste, err := data.DB.PasteGet(pasteID)
	if err != nil {
		return err
	}

	// If "one use" paste
	if paste.OneUse {
		if req.Form.Get("openOneUse") == "true" {
			// Delete paste
			err = data.DB.PasteDelete(pasteID)
			if err != nil {
				return err
			}

		} else {
			// Remove secret data
			paste = storage.Paste{
				ID:     paste.ID,
				OneUse: true,
			}
		}
	}

	// Return response
	rw.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(paste)
}
