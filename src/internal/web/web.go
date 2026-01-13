
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package web

import (
	"embed"
	"os"

	chromaLexers "github.com/alecthomas/chroma/v2/lexers"
	"github.com/casjay-forks/caspaste/src/internal/config"
	"github.com/casjay-forks/caspaste/src/internal/logger"
	"github.com/casjay-forks/caspaste/src/internal/netshare"
	"github.com/casjay-forks/caspaste/src/internal/storage"
	"html/template"
	"net/http"
	"strings"
	textTemplate "text/template"
)

//go:embed data/*
var embFS embed.FS

type Data struct {
	DB  storage.DB
	Log logger.Logger

	RateLimitNew *netshare.RateLimitSystem
	RateLimitGet *netshare.RateLimitSystem

	Lexers      []string
	Locales     Locales
	LocalesList LocalesList
	Themes      Themes
	ThemesList  ThemesList

	StyleCSS       *textTemplate.Template
	ErrorPage      *template.Template
	Main           *template.Template
	MainJS         *[]byte
	HistoryJS      *textTemplate.Template
	CodeJS         *textTemplate.Template
	PastePage      *template.Template
	PasteJS        *textTemplate.Template
	PasteContinue  *template.Template
	Settings       *template.Template
	ListPage       *template.Template
	About          *template.Template
	TermsOfUse     *template.Template
	Authors        *template.Template
	License        *template.Template
	SourceCodePage *template.Template

	Docs           *template.Template
	DocsApiV1      *template.Template
	DocsLibraries  *template.Template
	DocsCustomize  *template.Template

	EmbeddedPage     *template.Template
	EmbeddedHelpPage *template.Template

	Version string

	TitleMaxLen int
	BodyMaxLen  int
	MaxLifeTime int64

	ServerAbout      string
	ServerRules      string
	ServerTermsExist bool
	ServerTermsOfUse string
	SecurityTxt      string

	// Server info
	FQDN        string
	ServerTitle string
	AdminName   string
	AdminMail   string

	// Security contact
	SecurityContactEmail string
	SecurityContactName  string

	// Robots
	SiteRobotsAllow      string
	SiteRobotsDeny       string
	SiteRobotsAgentsDeny []string

	// Branding
	Logo    string
	Favicon string

	CasPasswdFile string

	UiDefaultLifeTime string
	UiDefaultTheme    string
}

// LoadContentWithOverride loads content from embedded FS or overrides from file
// If overridePath is specified and file exists, uses that; otherwise uses embedded
func LoadContentWithOverride(embeddedPath, overridePath string) (string, error) {
	var content []byte
	var err error

	// Try override file first if specified
	if overridePath != "" {
		content, err = os.ReadFile(overridePath)
		if err == nil {
			return string(content), nil
		}
		// File doesn't exist or error, fall back to embedded
	}

	// Use embedded content
	content, err = embFS.ReadFile(embeddedPath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func Load(db storage.DB, cfg config.Config) (*Data, error) {
	var data Data
	var err error

	// Setup base info
	data.DB = db
	data.Log = cfg.Log

	data.RateLimitNew = cfg.RateLimitNew
	data.RateLimitGet = cfg.RateLimitGet

	data.Version = cfg.Version

	data.TitleMaxLen = cfg.TitleMaxLen
	data.BodyMaxLen = cfg.BodyMaxLen
	data.MaxLifeTime = cfg.MaxLifeTime
	data.UiDefaultLifeTime = cfg.UiDefaultLifetime
	data.UiDefaultTheme = cfg.UiDefaultTheme
	data.CasPasswdFile = cfg.CasPasswdFile

	data.ServerAbout = cfg.ServerAbout
	data.ServerRules = cfg.ServerRules
	data.ServerTermsOfUse = cfg.ServerTermsOfUse

	serverTermsExist := false
	if cfg.ServerTermsOfUse != "" {
		serverTermsExist = true
	}
	data.ServerTermsExist = serverTermsExist

	data.AdminName = cfg.AdminName
	data.AdminMail = cfg.AdminMail

	data.FQDN = cfg.FQDN
	data.ServerTitle = cfg.ServerTitle
	data.SecurityContactEmail = cfg.SecurityContactEmail
	data.SecurityContactName = cfg.SecurityContactName
	data.SecurityTxt = cfg.SecurityTxt
	data.SiteRobotsAllow = cfg.SiteRobotsAllow
	data.SiteRobotsDeny = cfg.SiteRobotsDeny
	data.SiteRobotsAgentsDeny = cfg.SiteRobotsAgentsDeny
	data.Logo = cfg.Logo
	data.Favicon = cfg.Favicon

	// Get Chroma lexers
	data.Lexers = chromaLexers.Names(false)

	// Load locales
	data.Locales, data.LocalesList, err = loadLocales(embFS, "data/locale")
	if err != nil {
		return nil, err
	}

	// Load themes
	data.Themes, data.ThemesList, err = loadThemes(cfg.UiThemesDir, data.LocalesList, data.UiDefaultTheme)
	if err != nil {
		return nil, err
	}

	// style.css file
	data.StyleCSS, err = textTemplate.ParseFS(embFS, "data/style.css")
	if err != nil {
		return nil, err
	}

	// main.tmpl
	data.Main, err = template.ParseFS(embFS, "data/base.tmpl", "data/main.tmpl")
	if err != nil {
		return nil, err
	}

	// main.js
	mainJS, err := embFS.ReadFile("data/main.js")
	if err != nil {
		return nil, err
	}
	data.MainJS = &mainJS

	// history.js
	data.HistoryJS, err = textTemplate.ParseFS(embFS, "data/history.js")
	if err != nil {
		return nil, err
	}

	// code.js
	data.CodeJS, err = textTemplate.ParseFS(embFS, "data/code.js")
	if err != nil {
		return nil, err
	}

	// paste.tmpl
	data.PastePage, err = template.ParseFS(embFS, "data/base.tmpl", "data/paste.tmpl")
	if err != nil {
		return nil, err
	}

	// paste.js
	data.PasteJS, err = textTemplate.ParseFS(embFS, "data/paste.js")
	if err != nil {
		return nil, err
	}

	// paste_continue.tmpl
	data.PasteContinue, err = template.ParseFS(embFS, "data/base.tmpl", "data/paste_continue.tmpl")
	if err != nil {
		return nil, err
	}

	// settings.tmpl
	data.Settings, err = template.ParseFS(embFS, "data/base.tmpl", "data/settings.tmpl")
	if err != nil {
		return nil, err
	}

	// list.tmpl
	data.ListPage, err = template.ParseFS(embFS, "data/base.tmpl", "data/list.tmpl")
	if err != nil {
		return nil, err
	}

	// about.tmpl
	data.About, err = template.ParseFS(embFS, "data/base.tmpl", "data/about.tmpl")
	if err != nil {
		return nil, err
	}

	// terms.tmpl
	data.TermsOfUse, err = template.ParseFS(embFS, "data/base.tmpl", "data/terms.tmpl")
	if err != nil {
		return nil, err
	}

	// authors.tmpl
	data.Authors, err = template.ParseFS(embFS, "data/base.tmpl", "data/authors.tmpl")
	if err != nil {
		return nil, err
	}

	// license.tmpl
	data.License, err = template.ParseFS(embFS, "data/base.tmpl", "data/license.tmpl")
	if err != nil {
		return nil, err
	}

	// source_code.tmpl
	data.SourceCodePage, err = template.ParseFS(embFS, "data/base.tmpl", "data/source_code.tmpl")
	if err != nil {
		return nil, err
	}

	// docs.tmpl
	data.Docs, err = template.ParseFS(embFS, "data/base.tmpl", "data/docs.tmpl")
	if err != nil {
		return nil, err
	}

	// docs_apiv1.tmpl
	data.DocsApiV1, err = template.ParseFS(embFS, "data/base.tmpl", "data/docs_apiv1.tmpl")
	if err != nil {
		return nil, err
	}

	// docs_libraries.tmpl
	data.DocsLibraries, err = template.ParseFS(embFS, "data/base.tmpl", "data/docs_libraries.tmpl")
	if err != nil {
		return nil, err
	}

	// docs_customize.tmpl
	data.DocsCustomize, err = template.ParseFS(embFS, "data/base.tmpl", "data/docs_customize.tmpl")
	if err != nil {
		return nil, err
	}

	// error.tmpl
	data.ErrorPage, err = template.ParseFS(embFS, "data/base.tmpl", "data/error.tmpl")
	if err != nil {
		return nil, err
	}

	// emb.tmpl
	data.EmbeddedPage, err = template.ParseFS(embFS, "data/emb.tmpl")
	if err != nil {
		return nil, err
	}

	// emb_help.tmpl
	data.EmbeddedHelpPage, err = template.ParseFS(embFS, "data/base.tmpl", "data/emb_help.tmpl")
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (data *Data) Handler(rw http.ResponseWriter, req *http.Request) {
	// Process request
	var err error

	rw.Header().Set("Server", config.Software+"/"+data.Version)

	switch req.URL.Path {
	// Health checks
	case "/healthz":
		err = data.handleHealthz(rw, req)
	case "/api/healthz":
		err = data.handleAPIHealthz(rw, req)
	// Search engines
	case "/robots.txt":
		err = data.handleRobotsTxt(rw, req)
	case "/sitemap.xml":
		err = data.handleSitemap(rw, req)
	// Security
	case "/.well-known/security.txt":
		err = data.handleSecurityTxt(rw, req)
	// Resources
	case "/style.css":
		err = data.handleStyleCSS(rw, req)
	case "/main.js":
		err = data.handleMainJS(rw, req)
	case "/history.js":
		err = data.handleHistoryJS(rw, req)
	case "/code.js":
		err = data.handleCodeJS(rw, req)
	case "/paste.js":
		err = data.handlePasteJS(rw, req)
	// PWA Support
	case "/manifest.json":
		err = data.handleManifest(rw, req)
	case "/sw.js":
		err = data.handleServiceWorker(rw, req)
	case "/about":
		err = data.handleAbout(rw, req)
	case "/about/authors":
		err = data.handleAuthors(rw, req)
	case "/about/license":
		err = data.handleLicense(rw, req)
	case "/about/source_code":
		err = data.handleSourceCodePage(rw, req)
	case "/docs":
		err = data.handleDocs(rw, req)
	case "/docs/apiv1":
		err = data.handleDocsAPIv1(rw, req)
	case "/docs/libraries":
		err = data.handleDocsLibraries(rw, req)
	case "/docs/api_libs": // Redirect old URL
		http.Redirect(rw, req, "/docs/libraries", http.StatusMovedPermanently)
	case "/docs/customize":
		err = data.handleDocsCustomize(rw, req)
	// Pages
	case "/":
		err = data.handleNewPaste(rw, req)
	case "/list":
		err = data.handleList(rw, req)
	case "/settings":
		err = data.handleSettings(rw, req)
	case "/terms":
		err = data.handleTermsOfUse(rw, req)
	// Else
	default:
		if strings.HasPrefix(req.URL.Path, "/dl/") {
			err = data.handleDownload(rw, req)

		} else if strings.HasPrefix(req.URL.Path, "/emb/") {
			err = data.handleEmbedded(rw, req)

		} else if strings.HasPrefix(req.URL.Path, "/emb_help/") {
			err = data.handleEmbeddedHelp(rw, req)

		} else if strings.HasPrefix(req.URL.Path, "/u/") {
			err = data.handleURLRedirect(rw, req)

		} else if strings.HasPrefix(req.URL.Path, "/qr/") {
			err = data.handleQRCode(rw, req)

		} else if strings.HasPrefix(req.URL.Path, "/edit/") {
			err = data.handleEditPaste(rw, req)

		} else {
			err = data.handleGetPaste(rw, req)
		}
	}

	// Log
	if err == nil {
		data.Log.HttpRequest(req, 200)

	} else {
		code, err := data.writeError(rw, req, err)
		if err != nil {
			data.Log.HttpError(req, err)
		} else {
			data.Log.HttpRequest(req, code)
		}
	}
}
