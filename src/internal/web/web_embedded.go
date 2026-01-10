// Copyright (C) 2021-2023 Leonid Maslakov.

// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE file for details.

package web

import (
	"github.com/casjay-forks/caspaste/src/internal/netshare"
	"github.com/casjay-forks/caspaste/src/internal/storage"
	"html/template"
	"net/http"
	"time"
)

type embTmpl struct {
	ID            string
	CreateTimeStr string
	DeleteTime    int64
	OneUse        bool
	Title         string
	Body          template.HTML

	ErrorNotFound bool
	Translate     func(string, ...interface{}) template.HTML
}

// Pattern: /emb/
func (data *Data) embeddedHand(rw http.ResponseWriter, req *http.Request) error {
	errorNotFound := false

	// Check rate limit
	err := data.RateLimitGet.CheckAndUse(netshare.GetClientAddr(req))
	if err != nil {
		return err
	}

	// Get paste ID
	pasteID := string([]rune(req.URL.Path)[5:])

	// Read DB
	paste, err := data.DB.PasteGet(pasteID)
	if err != nil {
		if err == storage.ErrNotFoundID {
			errorNotFound = true

		} else {
			return err
		}
	}

	// Prepare template data
	createTime := time.Unix(paste.CreateTime, 0).UTC()

	tmplData := embTmpl{
		ID:            paste.ID,
		CreateTimeStr: createTime.Format("1 Jan, 2006"),
		DeleteTime:    paste.DeleteTime,
		OneUse:        paste.OneUse,
		Title:         paste.Title,
		Body:          tryHighlight(paste.Body, paste.Syntax, "monokai"),

		ErrorNotFound: errorNotFound,
		Translate:     data.Locales.findLocale(req).translate,
	}

	// Show paste
	return data.EmbeddedPage.Execute(rw, tmplData)
}
