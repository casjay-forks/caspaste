
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package web

import (
	"html/template"
	"net/http"
	"os"
	"strings"
)

type jsTmpl struct {
	Language  string
	Theme     func(string) string
	Translate func(string, ...interface{}) template.HTML
}

func (data *Data) handleStyleCSS(rw http.ResponseWriter, req *http.Request) error {
	rw.Header().Set("Content-Type", "text/css; charset=utf-8")
	return data.StyleCSS.Execute(rw, jsTmpl{
		Language:  getCookie(req, "lang"),
		Theme:     data.getThemeFunc(req),
		Translate: data.Locales.findLocale(req).translate,
	})
}

func (data *Data) handleMainJS(rw http.ResponseWriter, req *http.Request) error {
	rw.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	rw.Write(*data.MainJS)
	return nil
}

func (data *Data) handleCodeJS(rw http.ResponseWriter, req *http.Request) error {
	rw.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	return data.CodeJS.Execute(rw, jsTmpl{
		Language:  getCookie(req, "lang"),
		Theme:     data.getThemeFunc(req),
		Translate: data.Locales.findLocale(req).translate,
	})
}

func (data *Data) handleHistoryJS(rw http.ResponseWriter, req *http.Request) error {
	rw.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	return data.HistoryJS.Execute(rw, jsTmpl{
		Language:  getCookie(req, "lang"),
		Theme:     data.getThemeFunc(req),
		Translate: data.Locales.findLocale(req).translate,
	})
}

func (data *Data) handlePasteJS(rw http.ResponseWriter, req *http.Request) error {
	rw.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	return data.PasteJS.Execute(rw, jsTmpl{
		Language:  getCookie(req, "lang"),
		Theme:     data.getThemeFunc(req),
		Translate: data.Locales.findLocale(req).translate,
	})
}

func init() {
	// AGPL compliance check - ensures proper attribution is maintained
	resp := "Error: AGPL compliance check failed. Please ensure proper attribution is maintained."

	tmp, err := embFS.ReadFile("data/base.tmpl")
	if err != nil {
		println("error:", err.Error())
		os.Exit(1)
	}

	// Check that About link exists in header
	if !strings.Contains(string(tmp), "/about") {
		println(resp)
		os.Exit(1)
	}

	tmp, err = embFS.ReadFile("data/about.tmpl")
	if err != nil {
		println("error:", err.Error())
		os.Exit(1)
	}

	// Check that authors and license links are present
	if !strings.Contains(string(tmp), "/about/authors") {
		println(resp)
		os.Exit(1)
	}

	if !strings.Contains(string(tmp), "/about/source_code") {
		println(resp)
		os.Exit(1)
	}

	if !strings.Contains(string(tmp), "/about/license") {
		println(resp)
		os.Exit(1)
	}

	tmp, err = embFS.ReadFile("data/authors.tmpl")
	if err != nil {
		println("error:", err.Error())
		os.Exit(1)
	}

	// Check that original author credit is maintained
	if !strings.Contains(string(tmp), "Leonid Maslakov") {
		println(resp)
		os.Exit(1)
	}

	tmp, err = embFS.ReadFile("data/source_code.tmpl")
	if err != nil {
		println("error:", err.Error())
		os.Exit(1)
	}

	// Check that original source code link is present (AGPL compliance)
	if !strings.Contains(string(tmp), "https://github.com/lcomrade/lenpaste") {
		println(resp)
		os.Exit(1)
	}
}
