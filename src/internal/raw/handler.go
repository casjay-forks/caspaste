
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package raw

import (
	"github.com/casjay-forks/caspaste/src/internal/netshare"
	"io"
	"net/http"
)

// Pattern: /raw/
func (data *Data) rawHand(rw http.ResponseWriter, req *http.Request) error {
	// Check rate limit
	err := data.RateLimitGet.CheckAndUse(netshare.GetClientAddr(req))
	if err != nil {
		return err
	}

	// Read DB
	pasteID := string([]rune(req.URL.Path)[5:])

	paste, err := data.DB.PasteGet(pasteID)
	if err != nil {
		return err
	}

	// If "one use" paste
	if paste.OneUse {
		// Delete paste
		err = data.DB.PasteDelete(pasteID)
		if err != nil {
			return err
		}
	}

	// Write result
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")

	_, err = io.WriteString(rw, paste.Body)
	if err != nil {
		return err
	}

	return nil
}
