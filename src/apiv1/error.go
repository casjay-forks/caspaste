
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package apiv1

import (
	"encoding/json"
	"errors"
	"github.com/casjay-forks/caspaste/src/netshare"
	"github.com/casjay-forks/caspaste/src/storage"
	"net/http"
	"strconv"
)

type errorType struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func (data *Data) writeError(rw http.ResponseWriter, req *http.Request, e error) (int, error) {
	var resp errorType

	var eTmp429 *netshare.RateLimitError

	if e == netshare.ErrBadRequest {
		resp.Code = 400
		resp.Error = "Bad Request"

	} else if e == netshare.ErrUnauthorized {
		rw.Header().Add("WWW-Authenticate", "Basic")
		resp.Code = 401
		resp.Error = "Unauthorized"

	} else if e == storage.ErrNotFoundID {
		resp.Code = 404
		resp.Error = "Could not find ID"

	} else if e == netshare.ErrNotFound {
		resp.Code = 404
		resp.Error = "Not Found"

	} else if e == netshare.ErrMethodNotAllowed {
		resp.Code = 405
		resp.Error = "Method Not Allowed"

	} else if e == netshare.ErrPayloadTooLarge {
		resp.Code = 413
		resp.Error = "Payload Too Large"

	} else if e == netshare.ErrTooManyRequests || errors.As(e, &eTmp429) {
		resp.Code = 429
		resp.Error = "Too Many Requests"
		if eTmp429 != nil {
			rw.Header().Set("Retry-After", strconv.FormatInt(eTmp429.RetryAfter, 10))
		}

	} else {
		resp.Code = 500
		resp.Error = "Internal Server Error"
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(resp.Code)

	err := writeJSON(rw, resp)
	if err != nil {
		return 500, err
	}

	return resp.Code, nil
}

// writeJSON writes a JSON response with proper formatting per AI.md PART 14
// ALL JSON responses MUST be indented (2 spaces) and end with newline
func writeJSON(w http.ResponseWriter, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	// Single trailing newline per AI.md
	_, err = w.Write([]byte("\n"))
	return err
}
