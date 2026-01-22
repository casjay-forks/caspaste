
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package web

import (
	"encoding/base64"
	"html/template"
	"net/http"
	"time"

	"github.com/casjay-forks/caspaste/src/internal/lineend"
	"github.com/casjay-forks/caspaste/src/internal/netshare"
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

	// File upload fields
	IsFile   bool
	FileName string
	MimeType string
	FileSize int

	Language  string
	Theme     func(string) string
	Translate func(string, ...interface{}) template.HTML
}

type pasteContinueTmpl struct {
	ID        string
	Language  string
	Theme     func(string) string
	Translate func(string, ...interface{}) template.HTML
}

func (data *Data) handleGetPaste(rw http.ResponseWriter, req *http.Request) error {
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
	if paste.OneUse {
		// If continue button not pressed
		req.ParseForm()

		if req.PostForm.Get("oneUseContinue") != "true" {
			tmplData := pasteContinueTmpl{
				ID:        paste.ID,
				Language:  getCookie(req, "lang"),
				Theme:     data.getThemeFunc(req),
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

	// Determine body content based on whether this is a file upload
	var bodyContent string
	var fileSize int
	if paste.IsFile {
		// File upload: try to decode base64, fall back to raw for legacy data
		fileData, err := base64.StdEncoding.DecodeString(paste.Body)
		if err != nil {
			// Legacy data stored without base64 encoding - use as-is
			bodyContent = paste.Body
			fileSize = len(paste.Body)
		} else {
			bodyContent = string(fileData)
			fileSize = len(fileData)
		}
	} else {
		bodyContent = paste.Body
	}

	tmplData := pasteTmpl{
		ID:         paste.ID,
		Title:      paste.Title,
		Body:       data.Themes.findTheme(req, data.UiDefaultTheme).tryHighlight(bodyContent, paste.Syntax),
		Syntax:     paste.Syntax,
		CreateTime: paste.CreateTime,
		DeleteTime: paste.DeleteTime,
		OneUse:     paste.OneUse,

		CreateTimeStr: createTime.Format("Mon, 02 Jan 2006 15:04:05 -0700"),
		DeleteTimeStr: deleteTime.Format("Mon, 02 Jan 2006 15:04:05 -0700"),

		Author:      paste.Author,
		AuthorEmail: paste.AuthorEmail,
		AuthorURL:   paste.AuthorURL,

		IsFile:   paste.IsFile,
		FileName: paste.FileName,
		MimeType: paste.MimeType,
		FileSize: fileSize,

		Language:  getCookie(req, "lang"),
		Theme:     data.getThemeFunc(req),
		Translate: data.Locales.findLocale(req).translate,
	}

	// Get body line end
	switch lineend.GetLineEnd(bodyContent) {
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
