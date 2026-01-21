
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package web

import (
	"html/template"
	"net/http"
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
