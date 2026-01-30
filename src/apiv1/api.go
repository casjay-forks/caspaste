
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package apiv1

import (
	"net/http"
	"time"

	chromaLexers "github.com/alecthomas/chroma/v2/lexers"
	"github.com/casjay-forks/caspaste/src/config"
	"github.com/casjay-forks/caspaste/src/caspasswd"
	"github.com/casjay-forks/caspaste/src/logger"
	"github.com/casjay-forks/caspaste/src/netshare"
	"github.com/casjay-forks/caspaste/src/storage"
)

type Data struct {
	Log logger.Logger
	DB  storage.DB

	RateLimitNew *netshare.RateLimitSystem
	RateLimitGet *netshare.RateLimitSystem

	Lexers []string

	Version string

	TitleMaxLen int
	BodyMaxLen  int
	MaxLifeTime int64

	ServerAbout      string
	ServerRules      string
	ServerTermsOfUse string

	AdminName string
	AdminMail string

	// true = open/public, false = auth required
	Public        bool
	CasPasswdFile string
	BruteForce    *caspasswd.BruteForceProtection

	UiDefaultLifeTime string
}

func Load(db storage.DB, cfg config.Config) *Data {
	lexers := chromaLexers.Names(false)

	// Initialize brute force protection if authentication is required (server.public=false)
	var bruteForce *caspasswd.BruteForceProtection
	if !cfg.Public && cfg.CasPasswdFile != "" {
		// 5 failed attempts = 15 minute lockout
		bruteForce = caspasswd.NewBruteForceProtection(5, 15*time.Minute)
	}

	return &Data{
		DB:                db,
		Log:               cfg.Log,
		RateLimitNew:      cfg.RateLimitNew,
		RateLimitGet:      cfg.RateLimitGet,
		Lexers:            lexers,
		Version:           cfg.Version,
		TitleMaxLen:       cfg.TitleMaxLen,
		BodyMaxLen:        cfg.BodyMaxLen,
		MaxLifeTime:       cfg.MaxLifeTime,
		ServerAbout:       cfg.ServerAbout,
		ServerRules:       cfg.ServerRules,
		ServerTermsOfUse:  cfg.ServerTermsOfUse,
		AdminName:         cfg.AdminName,
		AdminMail:         cfg.AdminMail,
		Public:            cfg.Public,
		CasPasswdFile:     cfg.CasPasswdFile,
		BruteForce:        bruteForce,
		UiDefaultLifeTime: cfg.UiDefaultLifetime,
	}
}

func (data *Data) Hand(rw http.ResponseWriter, req *http.Request) {
	// Process request
	var err error

	rw.Header().Set("Server", config.Software+"/"+data.Version)

	switch req.URL.Path {
	// Health check per AI.md PART 13
	case "/api/v1/healthz":
		err = data.handleHealthz(rw, req)
	// API v1 endpoints
	case "/api/v1/new":
		err = data.newHand(rw, req)
	case "/api/v1/get":
		err = data.getHand(rw, req)
	case "/api/v1/list":
		err = data.handleList(rw, req)
	case "/api/v1/getServerInfo":
		err = data.getServerInfoHand(rw, req)
	default:
		err = netshare.ErrNotFound
	}

	// Log
	if err == nil {
		data.Log.HttpRequest(req, 200)

	} else {
		// Log the original error before writing HTTP response
		data.Log.HttpError(req, err)

		code, writeErr := data.writeError(rw, req, err)
		if writeErr != nil {
			data.Log.HttpError(req, writeErr)
		}
		data.Log.HttpRequest(req, code)
	}
}
