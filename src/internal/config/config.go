// Copyright (C) 2021-2023 Leonid Maslakov.

// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE file for details.

package config

import (
	"github.com/casjay-forks/caspaste/src/internal/logger"
	"github.com/casjay-forks/caspaste/src/internal/netshare"
)

const Software = "CasPaste"

type Config struct {
	Log logger.Logger

	RateLimitNew *netshare.RateLimitSystem
	RateLimitGet *netshare.RateLimitSystem

	Version string

	TitleMaxLen int
	BodyMaxLen  int
	MaxLifeTime int64

	ServerAbout      string
	ServerRules      string
	ServerTermsOfUse string

	AdminName string
	AdminMail string

	RobotsDisallow bool

	CasPasswdFile string

	// TrustReverseProxy controls whether to trust X-Forwarded-* and similar headers
	// Set to true only when behind a trusted reverse proxy (nginx, caddy, etc.)
	// WARNING: Setting to true when not behind a proxy allows IP spoofing
	TrustReverseProxy bool

	UiDefaultLifetime string
	UiDefaultTheme    string
	UiThemesDir       string
}
