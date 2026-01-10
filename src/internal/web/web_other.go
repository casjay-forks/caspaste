// Copyright (C) 2021-2023 Leonid Maslakov.

// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE file for details.

package web

import (
	"html/template"
	"net/http"
	"os"
	"strings"
)

type jsTmpl struct {
	Translate func(string, ...interface{}) template.HTML
	Theme     func(string) string
}

func (data *Data) styleCSSHand(rw http.ResponseWriter, req *http.Request) error {
	rw.Header().Set("Content-Type", "text/css; charset=utf-8")
	return data.StyleCSS.Execute(rw, jsTmpl{
		Translate: data.Locales.findLocale(req).translate,
		Theme:     data.Themes.findTheme(req, data.UiDefaultTheme).theme,
	})
}

func (data *Data) mainJSHand(rw http.ResponseWriter, req *http.Request) error {
	rw.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	rw.Write(*data.MainJS)
	return nil
}

func (data *Data) codeJSHand(rw http.ResponseWriter, req *http.Request) error {
	rw.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	return data.CodeJS.Execute(rw, jsTmpl{Translate: data.Locales.findLocale(req).translate})
}

func (data *Data) historyJSHand(rw http.ResponseWriter, req *http.Request) error {
	rw.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	return data.HistoryJS.Execute(rw, jsTmpl{
		Translate: data.Locales.findLocale(req).translate,
		Theme:     data.Themes.findTheme(req, data.UiDefaultTheme).theme,
	})
}

func (data *Data) pasteJSHand(rw http.ResponseWriter, req *http.Request) error {
	rw.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	return data.PasteJS.Execute(rw, jsTmpl{Translate: data.Locales.findLocale(req).translate})
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
	if strings.Contains(string(tmp), "/about") == false {
		println(resp)
		os.Exit(1)
	}

	tmp, err = embFS.ReadFile("data/about.tmpl")
	if err != nil {
		println("error:", err.Error())
		os.Exit(1)
	}

	// Check that authors and license links are present
	if strings.Contains(string(tmp), "/about/authors") == false {
		println(resp)
		os.Exit(1)
	}

	if strings.Contains(string(tmp), "/about/source_code") == false {
		println(resp)
		os.Exit(1)
	}

	if strings.Contains(string(tmp), "/about/license") == false {
		println(resp)
		os.Exit(1)
	}

	tmp, err = embFS.ReadFile("data/authors.tmpl")
	if err != nil {
		println("error:", err.Error())
		os.Exit(1)
	}

	// Check that original author credit is maintained
	if strings.Contains(string(tmp), "Leonid Maslakov") == false {
		println(resp)
		os.Exit(1)
	}

	tmp, err = embFS.ReadFile("data/source_code.tmpl")
	if err != nil {
		println("error:", err.Error())
		os.Exit(1)
	}

	// Check that original source code link is present (AGPL compliance)
	if strings.Contains(string(tmp), "https://github.com/lcomrade/lenpaste") == false {
		println(resp)
		os.Exit(1)
	}
}
