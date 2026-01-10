
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package web

import (
	"github.com/casjay-forks/caspaste/src/internal/netshare"
	"html/template"
	"net/http"
)

type docsTmpl struct {
	Highlight func(string, string) template.HTML
	Translate func(string, ...interface{}) template.HTML
}

type docsApiV1Tmpl struct {
	MaxLenAuthorAll int

	Highlight func(string, string) template.HTML
	Translate func(string, ...interface{}) template.HTML
}

// Pattern: /docs
func (data *Data) docsHand(rw http.ResponseWriter, req *http.Request) error {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	return data.Docs.Execute(rw, docsTmpl{Translate: data.Locales.findLocale(req).translate})
}

// Pattern: /docs/apiv1
func (data *Data) docsApiV1Hand(rw http.ResponseWriter, req *http.Request) error {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	return data.DocsApiV1.Execute(rw, docsApiV1Tmpl{
		MaxLenAuthorAll: netshare.MaxLengthAuthorAll,
		Translate:       data.Locales.findLocale(req).translate,
		Highlight:       data.Themes.findTheme(req, data.UiDefaultTheme).tryHighlight,
	})
}

// Pattern: /docs/api_libs
func (data *Data) docsApiLibsHand(rw http.ResponseWriter, req *http.Request) error {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	return data.DocsApiLibs.Execute(rw, docsTmpl{Translate: data.Locales.findLocale(req).translate})
}
