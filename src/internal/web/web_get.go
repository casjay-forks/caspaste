
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package web

import (
	"github.com/casjay-forks/caspaste/src/internal/lineend"
	"github.com/casjay-forks/caspaste/src/internal/netshare"
	"html/template"
	"net/http"
	"time"
)

type pasteTmpl struct {
	ID         string
	Title      string
	Body       template.HTML
	Syntax     string
	CreateTime int64
	DeleteTime int64
	OneUse     bool

	LineEnd       string
	CreateTimeStr string
	DeleteTimeStr string

	Author      string
	AuthorEmail string
	AuthorURL   string

	Translate func(string, ...interface{}) template.HTML
}

type pasteContinueTmpl struct {
	ID        string
	Translate func(string, ...interface{}) template.HTML
}

func (data *Data) getPasteHand(rw http.ResponseWriter, req *http.Request) error {
	// Check rate limit
	err := data.RateLimitGet.CheckAndUse(netshare.GetClientAddr(req))
	if err != nil {
		return err
	}

	// Get paste ID
	pasteID := string([]rune(req.URL.Path)[1:])

	// Read DB
	paste, err := data.DB.PasteGet(pasteID)
	if err != nil {
		return err
	}

	// If "one use" paste
	if paste.OneUse == true {
		// If continue button not pressed
		req.ParseForm()

		if req.PostForm.Get("oneUseContinue") != "true" {
			tmplData := pasteContinueTmpl{
				ID:        paste.ID,
				Translate: data.Locales.findLocale(req).translate,
			}

			return data.PasteContinue.Execute(rw, tmplData)
		}

		// If continue button pressed delete paste
		err = data.DB.PasteDelete(pasteID)
		if err != nil {
			return err
		}
	}

	// Prepare template data
	createTime := time.Unix(paste.CreateTime, 0).UTC()
	deleteTime := time.Unix(paste.DeleteTime, 0).UTC()

	tmplData := pasteTmpl{
		ID:         paste.ID,
		Title:      paste.Title,
		Body:       data.Themes.findTheme(req, data.UiDefaultTheme).tryHighlight(paste.Body, paste.Syntax),
		Syntax:     paste.Syntax,
		CreateTime: paste.CreateTime,
		DeleteTime: paste.DeleteTime,
		OneUse:     paste.OneUse,

		CreateTimeStr: createTime.Format("Mon, 02 Jan 2006 15:04:05 -0700"),
		DeleteTimeStr: deleteTime.Format("Mon, 02 Jan 2006 15:04:05 -0700"),

		Author:      paste.Author,
		AuthorEmail: paste.AuthorEmail,
		AuthorURL:   paste.AuthorURL,

		Translate: data.Locales.findLocale(req).translate,
	}

	// Get body line end
	switch lineend.GetLineEnd(paste.Body) {
	case "\r\n":
		tmplData.LineEnd = "CRLF"
	case "\r":
		tmplData.LineEnd = "CR"
	default:
		tmplData.LineEnd = "LF"
	}

	// Show paste
	return data.PastePage.Execute(rw, tmplData)
}
