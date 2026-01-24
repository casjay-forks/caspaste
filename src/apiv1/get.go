
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package apiv1

import (
	"net/http"

	"github.com/casjay-forks/caspaste/src/netshare"
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

	// If "one use" (burn after reading) paste - delete it after returning content
	if paste.OneUse {
		// Delete paste immediately - burn after reading just works
		err = data.DB.PasteDelete(pasteID)
		if err != nil {
			return err
		}
	}

	// Return response per AI.md PART 14 (indented JSON with newline)
	return writeJSON(rw, paste)
}
