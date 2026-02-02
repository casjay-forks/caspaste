
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package config

import (
	"github.com/casjay-forks/caspaste/src/logger"
	"github.com/casjay-forks/caspaste/src/netshare"
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

	// Content
	ServerAbout      string
	ServerRules      string
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

	// Authentication
	// true = open/public (no auth), false = auth required
	Public        bool
	CasPasswdFile string

	// Trusted proxy configuration (for X-Forwarded-* headers)
	TrustedProxies []string

	UiDefaultLifetime string
	UiDefaultTheme    string
	UiThemesDir       string
}
