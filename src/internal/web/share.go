
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package web

import (
	"net/http"
)

func getCookie(req *http.Request, name string) string {
	cookie, err := req.Cookie(name)
	if err != nil {
		return ""
	}

	return cookie.Value
}
