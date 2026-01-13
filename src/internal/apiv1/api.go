
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package apiv1

import (
	"net/http"
	"time"

	chromaLexers "github.com/alecthomas/chroma/v2/lexers"
	"github.com/casjay-forks/caspaste/src/internal/config"
	"github.com/casjay-forks/caspaste/src/internal/caspasswd"
	"github.com/casjay-forks/caspaste/src/internal/logger"
	"github.com/casjay-forks/caspaste/src/internal/netshare"
	"github.com/casjay-forks/caspaste/src/internal/storage"
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

	CasPasswdFile string
	BruteForce    *caspasswd.BruteForceProtection

	UiDefaultLifeTime string
}

func Load(db storage.DB, cfg config.Config) *Data {
	lexers := chromaLexers.Names(false)

	// Initialize brute force protection if authentication is enabled
	var bruteForce *caspasswd.BruteForceProtection
	if cfg.CasPasswdFile != "" {
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
	// Health check
	case "/api/healthz":
		err = data.healthzHand(rw, req)
	// API v1 endpoints
	case "/api/v1/new":
		err = data.newHand(rw, req)
	case "/api/v1/get":
		err = data.getHand(rw, req)
	case "/api/v1/list":
		err = data.listHand(rw, req)
	case "/api/v1/getServerInfo":
		err = data.getServerInfoHand(rw, req)
	default:
		err = netshare.ErrNotFound
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
