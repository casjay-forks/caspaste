
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package web

import (
	"github.com/casjay-forks/caspaste/src/internal/caspasswd"
	"github.com/casjay-forks/caspaste/src/internal/netshare"
	"html/template"
	"net/http"
)

type createTmpl struct {
	Language          string
	Theme             func(string) string
	TitleMaxLen       int
	BodyMaxLen        int
	AuthorAllMaxLen   int
	MaxLifeTime       int64
	UiDefaultLifeTime string
	Lexers            []string
	ServerTermsExist  bool

	AuthorDefault      string
	AuthorEmailDefault string
	AuthorURLDefault   string

	AuthOk bool

	Translate func(string, ...interface{}) template.HTML
}

func (data *Data) handleNewPaste(rw http.ResponseWriter, req *http.Request) error {
	var err error

	// Check auth
	isAuthenticated := true

	if data.CasPasswdFile != "" {
		isAuthenticated = false

		user, pass, authProvided := req.BasicAuth()
		if authProvided {
			isAuthenticated, err = caspasswd.LoadAndCheck(data.CasPasswdFile, user, pass)
			if err != nil {
				return err
			}
		}

		if !isAuthenticated {
			rw.Header().Add("WWW-Authenticate", "Basic")
			rw.WriteHeader(401)
		}
	}

	// Create paste if need
	if req.Method == "POST" {
		pasteID, _, _, err := netshare.PasteAddFromForm(req, data.DB, data.RateLimitNew, data.TitleMaxLen, data.BodyMaxLen, data.MaxLifeTime, data.Lexers)
		if err != nil {
			return err
		}

		// Redirect to paste
		writeRedirect(rw, req, "/"+pasteID, 302)
		return nil
	}

	// Get theme
	themeName := getCookie(req, "theme")
	if themeName == "" {
		themeName = data.UiDefaultTheme
	}
	themeMap, exists := data.Themes[themeName]
	if !exists {
		themeMap = data.Themes[data.UiDefaultTheme]
	}
	themeLookup := func(key string) string {
		return themeMap[key]
	}

	// Else show create page
	tmplData := createTmpl{
		Language:           getCookie(req, "lang"),
		Theme:              themeLookup,
		TitleMaxLen:        data.TitleMaxLen,
		BodyMaxLen:         data.BodyMaxLen,
		AuthorAllMaxLen:    netshare.MaxLengthAuthorAll,
		MaxLifeTime:        data.MaxLifeTime,
		UiDefaultLifeTime:  data.UiDefaultLifeTime,
		Lexers:             data.Lexers,
		ServerTermsExist:   data.ServerTermsExist,
		AuthorDefault:      getCookie(req, "author"),
		AuthorEmailDefault: getCookie(req, "authorEmail"),
		AuthorURLDefault:   getCookie(req, "authorURL"),
		AuthOk:             isAuthenticated,
		Translate:          data.Locales.findLocale(req).translate,
	}

	rw.Header().Set("Content-Type", "text/html; charset=utf-8")

	return data.Main.Execute(rw, tmplData)
}
