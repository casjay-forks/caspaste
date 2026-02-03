
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package web

import (
	"net/http"
)

// handleUserDashboard handles GET /users (user dashboard)
func (data *Data) handleUserDashboard(rw http.ResponseWriter, req *http.Request) error {
	if req.Method != http.MethodGet {
		return ErrMethodNotAllowed
	}

	// Get authenticated user from context
	authUser := GetAuthUser(req.Context())
	if authUser == nil {
		http.Redirect(rw, req, "/login", http.StatusFound)
		return nil
	}

	return data.renderUserDashboard(rw, req, authUser)
}

// handleUserSettings handles GET/POST /users/settings
func (data *Data) handleUserSettings(rw http.ResponseWriter, req *http.Request) error {
	authUser := GetAuthUser(req.Context())
	if authUser == nil {
		http.Redirect(rw, req, "/login", http.StatusFound)
		return nil
	}

	// For now, redirect to existing settings page
	http.Redirect(rw, req, "/settings", http.StatusFound)
	return nil
}

// handleUserSecurity handles GET /users/security
func (data *Data) handleUserSecurity(rw http.ResponseWriter, req *http.Request) error {
	authUser := GetAuthUser(req.Context())
	if authUser == nil {
		http.Redirect(rw, req, "/login", http.StatusFound)
		return nil
	}

	return data.renderUserSecurity(rw, req, authUser)
}

// handleUserTokens handles GET/POST /users/tokens
func (data *Data) handleUserTokens(rw http.ResponseWriter, req *http.Request) error {
	authUser := GetAuthUser(req.Context())
	if authUser == nil {
		http.Redirect(rw, req, "/login", http.StatusFound)
		return nil
	}

	return data.renderUserTokens(rw, req, authUser)
}

// handleUserDomains handles GET /users/domains
func (data *Data) handleUserDomains(rw http.ResponseWriter, req *http.Request) error {
	authUser := GetAuthUser(req.Context())
	if authUser == nil {
		http.Redirect(rw, req, "/login", http.StatusFound)
		return nil
	}

	return data.renderUserDomains(rw, req, authUser)
}

// Render functions - these will use templates

func (data *Data) renderUserDashboard(rw http.ResponseWriter, req *http.Request, user *AuthUser) error {
	// Get locale
	locale := data.Locales.findLocale(req)

	templateData := map[string]interface{}{
		"Title":          "Dashboard",
		"User":           user,
		"Version":        data.Version,
		"FQDN":           data.FQDN,
		"ServerTitle":    data.ServerTitle,
		"LocalesList":    data.LocalesList,
		"ThemesList":     data.ThemesList,
		"UiDefaultTheme": data.UiDefaultTheme,
		"Translate":      locale.translate,
	}

	rw.Header().Set("Content-Type", "text/html; charset=UTF-8")

	// For now, use a simple HTML response until we have the full template
	html := `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Dashboard - ` + data.ServerTitle + `</title>
	<link rel="stylesheet" href="/style.css">
</head>
<body>
	<div class="container">
		<h1>Welcome, ` + user.Username + `!</h1>
		<nav>
			<ul>
				<li><a href="/users/settings">Settings</a></li>
				<li><a href="/users/security">Security</a></li>
				<li><a href="/users/tokens">API Tokens</a></li>
				<li><a href="/users/domains">Custom Domains</a></li>
				<li><a href="/orgs">Organizations</a></li>
				<li><a href="/logout">Logout</a></li>
			</ul>
		</nav>
	</div>
</body>
</html>`

	_, err := rw.Write([]byte(html))
	_ = templateData // Will be used when full template is implemented
	return err
}

func (data *Data) renderUserSecurity(rw http.ResponseWriter, req *http.Request, user *AuthUser) error {
	rw.Header().Set("Content-Type", "text/html; charset=UTF-8")

	html := `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Security Settings - ` + data.ServerTitle + `</title>
	<link rel="stylesheet" href="/style.css">
</head>
<body>
	<div class="container">
		<h1>Security Settings</h1>
		<section>
			<h2>Two-Factor Authentication</h2>
			<p>Status: ` + boolToStr(user.TOTPEnabled, "Enabled", "Disabled") + `</p>
			<form action="/api/v1/users/security/2fa/enable" method="POST">
				<button type="submit">` + boolToStr(user.TOTPEnabled, "Manage 2FA", "Enable 2FA") + `</button>
			</form>
		</section>
		<section>
			<h2>Sessions</h2>
			<p><a href="/api/v1/users/sessions">View Active Sessions</a></p>
		</section>
		<section>
			<h2>Change Password</h2>
			<form action="/api/v1/users/security/password" method="POST">
				<input type="password" name="current_password" placeholder="Current Password" required>
				<input type="password" name="new_password" placeholder="New Password" required>
				<button type="submit">Change Password</button>
			</form>
		</section>
		<p><a href="/users">Back to Dashboard</a></p>
	</div>
</body>
</html>`

	_, err := rw.Write([]byte(html))
	return err
}

func (data *Data) renderUserTokens(rw http.ResponseWriter, req *http.Request, user *AuthUser) error {
	rw.Header().Set("Content-Type", "text/html; charset=UTF-8")

	html := `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>API Tokens - ` + data.ServerTitle + `</title>
	<link rel="stylesheet" href="/style.css">
</head>
<body>
	<div class="container">
		<h1>API Tokens</h1>
		<section>
			<h2>Create New Token</h2>
			<form action="/api/v1/users/tokens" method="POST">
				<input type="text" name="name" placeholder="Token Name" required>
				<select name="scopes">
					<option value="read">Read Only</option>
					<option value="read-write">Read/Write</option>
					<option value="global">Full Access</option>
				</select>
				<button type="submit">Create Token</button>
			</form>
		</section>
		<section>
			<h2>Active Tokens</h2>
			<p>View and manage your API tokens via the API: GET /api/v1/users/tokens</p>
		</section>
		<p><a href="/users">Back to Dashboard</a></p>
	</div>
</body>
</html>`

	_, err := rw.Write([]byte(html))
	return err
}

func (data *Data) renderUserDomains(rw http.ResponseWriter, req *http.Request, user *AuthUser) error {
	rw.Header().Set("Content-Type", "text/html; charset=UTF-8")

	html := `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Custom Domains - ` + data.ServerTitle + `</title>
	<link rel="stylesheet" href="/style.css">
</head>
<body>
	<div class="container">
		<h1>Custom Domains</h1>
		<section>
			<h2>Add New Domain</h2>
			<form action="/api/v1/users/domains" method="POST">
				<input type="text" name="domain" placeholder="yourdomain.com" required>
				<button type="submit">Add Domain</button>
			</form>
		</section>
		<section>
			<h2>Your Domains</h2>
			<p>View and manage your domains via the API: GET /api/v1/users/domains</p>
		</section>
		<p><a href="/users">Back to Dashboard</a></p>
	</div>
</body>
</html>`

	_, err := rw.Write([]byte(html))
	return err
}

// Helper function
func boolToStr(b bool, trueStr, falseStr string) string {
	if b {
		return trueStr
	}
	return falseStr
}

// ErrMethodNotAllowed is returned for unsupported HTTP methods
var ErrMethodNotAllowed = &httpError{code: 405, message: "Method not allowed"}

type httpError struct {
	code    int
	message string
}

func (e *httpError) Error() string {
	return e.message
}
