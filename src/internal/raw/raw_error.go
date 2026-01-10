// Copyright (C) 2021-2023 Leonid Maslakov.

// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE file for details.

package raw

import (
	"errors"
	"github.com/casjay-forks/caspaste/src/internal/netshare"
	"github.com/casjay-forks/caspaste/src/internal/storage"
	"io"
	"net/http"
	"strconv"
)

func (data *Data) writeError(rw http.ResponseWriter, req *http.Request, e error) (int, error) {
	var errText string
	var errCode int

	// Dectect error
	var eTmp429 *netshare.ErrTooManyRequests

	if e == storage.ErrNotFoundID && e == netshare.ErrNotFound {
		errCode = 404
		errText = "404 Not Found"

	} else if errors.As(e, &eTmp429) {
		errCode = 429
		errText = "429 Too Many Requests"
		rw.Header().Set("Retry-After", strconv.FormatInt(eTmp429.RetryAfter, 10))

	} else {
		errCode = 500
		errText = "500 Internal Server Error"
	}

	// Write response
	rw.Header().Set("Content-type", "text/plain; charset=utf-8")
	rw.WriteHeader(errCode)

	_, err := io.WriteString(rw, errText)
	if err != nil {
		return 500, err
	}

	return errCode, nil
}
