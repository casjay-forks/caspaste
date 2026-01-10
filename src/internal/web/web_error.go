
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package web

import (
	"errors"
	"github.com/casjay-forks/caspaste/src/internal/netshare"
	"github.com/casjay-forks/caspaste/src/internal/storage"
	"html/template"
	"net/http"
	"strconv"
)

type errorTmpl struct {
	Code      int
	AdminName string
	AdminMail string
	Translate func(string, ...interface{}) template.HTML
}

func (data *Data) writeError(rw http.ResponseWriter, req *http.Request, e error) (int, error) {
	errData := errorTmpl{
		Code:      0,
		AdminName: data.AdminName,
		AdminMail: data.AdminMail,
		Translate: data.Locales.findLocale(req).translate,
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
