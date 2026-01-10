// Copyright (C) 2021-2023 Leonid Maslakov.

// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE file for details.

package web

import (
	"html/template"
	"net/http"
)

type aboutTmpl struct {
	Version     string
	TitleMaxLen int
	BodyMaxLen  int
	MaxLifeTime int64

	ServerAbout      string
	ServerRules      string
	ServerTermsExist bool

	AdminName string
	AdminMail string

	Highlight func(string, string) template.HTML
	Translate func(string, ...interface{}) template.HTML
}

type aboutMinTmp struct {
	Translate func(string, ...interface{}) template.HTML
}

// Pattern: /about
func (data *Data) aboutHand(rw http.ResponseWriter, req *http.Request) error {
	dataTmpl := aboutTmpl{
		Version:          data.Version,
		TitleMaxLen:      data.TitleMaxLen,
		BodyMaxLen:       data.BodyMaxLen,
		MaxLifeTime:      data.MaxLifeTime,
		ServerAbout:      data.ServerAbout,
		ServerRules:      data.ServerRules,
		ServerTermsExist: data.ServerTermsExist,
		AdminName:        data.AdminName,
		AdminMail:        data.AdminMail,
		Highlight:        data.Themes.findTheme(req, data.UiDefaultTheme).tryHighlight,
		Translate:        data.Locales.findLocale(req).translate,
	}

	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	return data.About.Execute(rw, dataTmpl)
}

// Pattern: /about/authors
func (data *Data) authorsHand(rw http.ResponseWriter, req *http.Request) error {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	return data.Authors.Execute(rw, aboutMinTmp{Translate: data.Locales.findLocale(req).translate})
}

// Pattern: /about/license
func (data *Data) licenseHand(rw http.ResponseWriter, req *http.Request) error {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	return data.License.Execute(rw, aboutMinTmp{Translate: data.Locales.findLocale(req).translate})
}

// Pattern: /about/source_code
func (data *Data) sourceCodePageHand(rw http.ResponseWriter, req *http.Request) error {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	return data.SourceCodePage.Execute(rw, aboutMinTmp{Translate: data.Locales.findLocale(req).translate})
}
