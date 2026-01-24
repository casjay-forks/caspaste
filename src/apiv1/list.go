
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package apiv1

import (
	"net/http"
	"strconv"

	"github.com/casjay-forks/caspaste/src/netshare"
)

// GET /api/v1/list
func (data *Data) handleList(rw http.ResponseWriter, req *http.Request) error {
	// Check method
	if req.Method != "GET" {
		return netshare.ErrMethodNotAllowed
	}

	// Check rate limit
	err := data.RateLimitGet.CheckAndUse(netshare.GetClientAddr(req))
	if err != nil {
		return err
	}

	// Parse query parameters
	query := req.URL.Query()

	limit := 50
	if limitStr := query.Get("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 || parsedLimit > 100 {
			return netshare.ErrBadRequest
		}
		limit = parsedLimit
	}

	offset := 0
	if offsetStr := query.Get("offset"); offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil || parsedOffset < 0 {
			return netshare.ErrBadRequest
		}
		offset = parsedOffset
	}

	// Get paste list from database
	pastes, err := data.DB.PasteList(limit, offset)
	if err != nil {
		return err
	}

	// Return response per AI.md PART 14 (indented JSON with newline)
	return writeJSON(rw, pastes)
}
