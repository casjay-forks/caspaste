
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package web

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/casjay-forks/caspaste/src/internal/netshare"
	"github.com/casjay-forks/caspaste/src/internal/storage"
)

type errorTmpl struct {
	Code      int
	AdminName string
	AdminMail string
	Language  string       // Language for base template
	Theme     func(string) string // Theme function to get theme values
	Translate func(string, ...interface{}) template.HTML
}

func (data *Data) writeError(rw http.ResponseWriter, req *http.Request, e error) (int, error) {
	// DEBUG: Print error to stdout
	fmt.Printf("DEBUG writeError called with error: %v\n", e)
	
	locale := data.Locales.findLocale(req)
	
	// Get theme name, use default if not set
	themeName := getCookie(req, "theme")
	if themeName == "" {
		themeName = data.UiDefaultTheme
	}
	
	// Get theme map
	themeMap, exists := data.Themes[themeName]
	if !exists {
		// Fallback to default theme if specified theme doesn't exist
		themeMap = data.Themes[data.UiDefaultTheme]
	}
	
	// Create theme lookup function
	themeLookup := func(key string) string {
		return themeMap[key]
	}
	
	errData := errorTmpl{
		Code:      0,
		AdminName: data.AdminName,
		AdminMail: data.AdminMail,
		Language:  getCookie(req, "lang"), // Get language from cookie
		Theme:     themeLookup,             // Theme lookup function
		Translate: locale.translate,
	}

	// Dectect error
	var eTmp429 *netshare.RateLimitError

	if e == netshare.ErrBadRequest {
		errData.Code = 400

	} else if e == netshare.ErrUnauthorized {
		errData.Code = 401

	} else if e == storage.ErrNotFoundID {
		errData.Code = 404

	} else if e == netshare.ErrNotFound {
		errData.Code = 404

	} else if e == netshare.ErrMethodNotAllowed {
		errData.Code = 405

	} else if e == netshare.ErrPayloadTooLarge {
		errData.Code = 413

	} else if errors.As(e, &eTmp429) {
		errData.Code = 429
		rw.Header().Set("Retry-After", strconv.FormatInt(eTmp429.RetryAfter, 10))

	} else {
		errData.Code = 500
	}

	// Write response header
	rw.Header().Set("Content-type", "text/html; charset=utf-8")
	rw.WriteHeader(errData.Code)

	// Render template
	err := data.ErrorPage.Execute(rw, errData)
	if err != nil {
		return 500, err
	}

	return errData.Code, nil
}
