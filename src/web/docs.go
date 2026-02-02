
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package web

import (
	"github.com/casjay-forks/caspaste/src/netshare"
	"html/template"
	"net/http"
)

type docsTmpl struct {
	Language  string
	Theme     func(string) string
	Highlight func(string, string) template.HTML
	Translate func(string, ...interface{}) template.HTML
}

type docsApiV1Tmpl struct {
	MaxLenAuthorAll int

	Language  string
	Theme     func(string) string
	Highlight func(string, string) template.HTML
	Translate func(string, ...interface{}) template.HTML
}

// Pattern: /docs
func (data *Data) handleDocs(rw http.ResponseWriter, req *http.Request) error {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	return data.Docs.Execute(rw, docsTmpl{
		Language:  getCookie(req, "lang"),
		Theme:     data.getThemeFunc(req),
		Translate: data.Locales.findLocale(req).translate,
	})
}

// Pattern: /docs/apiv1
func (data *Data) handleDocsAPIv1(rw http.ResponseWriter, req *http.Request) error {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	return data.DocsApiV1.Execute(rw, docsApiV1Tmpl{
		MaxLenAuthorAll: netshare.MaxLengthAuthorAll,
		Language:        getCookie(req, "lang"),
		Theme:           data.getThemeFunc(req),
		Translate:       data.Locales.findLocale(req).translate,
		Highlight:       data.Themes.findTheme(req, data.UiDefaultTheme).tryHighlight,
	})
}

// Pattern: /docs/libraries
func (data *Data) handleDocsLibraries(rw http.ResponseWriter, req *http.Request) error {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	return data.DocsLibraries.Execute(rw, docsTmpl{
		Language:  getCookie(req, "lang"),
		Theme:     data.getThemeFunc(req),
		Translate: data.Locales.findLocale(req).translate,
	})
}

// Pattern: /docs/customize
func (data *Data) handleDocsCustomize(rw http.ResponseWriter, req *http.Request) error {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	return data.DocsCustomize.Execute(rw, docsTmpl{
		Language:  getCookie(req, "lang"),
		Theme:     data.getThemeFunc(req),
		Translate: data.Locales.findLocale(req).translate,
		Highlight: data.Themes.findTheme(req, data.UiDefaultTheme).tryHighlight,
	})
}
