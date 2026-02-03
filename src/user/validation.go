
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package user

import (
	"errors"
	"regexp"
	"strings"
)

// Username validation rules per PART 34
const (
	UsernameMinLength = 3
	UsernameMaxLength = 32
	PasswordMinLength = 8
	PasswordMaxLength = 128
	BioMaxLength      = 500
	EmailMaxLength    = 254
)

// Username regex: lowercase alphanumeric, underscore, hyphen
// Must start with letter, cannot end with _ or -
var usernameRegex = regexp.MustCompile(`^[a-z][a-z0-9_-]*[a-z0-9]$|^[a-z]$`)
var consecutiveChars = regexp.MustCompile(`__|--|_-|-_`)
var emailRegex = regexp.MustCompile(`^[a-z0-9.+_-]+@[a-z0-9][a-z0-9.-]*[a-z0-9]\.[a-z]{2,}$`)

// UsernameBlocklist contains words that cannot be used as usernames per PART 34
var UsernameBlocklist = []string{
	// System & Administrative
	"admin", "administrator", "root", "system", "sysadmin", "superuser",
	"master", "owner", "operator", "manager", "moderator", "mod",
	"staff", "support", "helpdesk", "help", "service", "daemon",

	// Server & Technical
	"server", "host", "node", "cluster", "api", "www", "web", "mail",
	"email", "smtp", "ftp", "ssh", "dns", "proxy", "gateway", "router",
	"firewall", "localhost", "local", "internal", "external", "public",
	"private", "network", "database", "db", "cache", "redis", "mysql",
	"postgres", "mongodb", "elastic", "nginx", "apache", "docker",

	// Application & Service Names
	"app", "application", "bot", "robot", "crawler", "spider", "scraper",
	"webhook", "callback", "cron", "scheduler", "worker", "queue", "job",
	"task", "process", "service", "microservice", "lambda", "function",

	// Authentication & Security
	"auth", "authentication", "login", "logout", "signin", "signout",
	"signup", "register", "password", "passwd", "token", "oauth", "sso",
	"saml", "ldap", "kerberos", "security", "secure", "ssl", "tls",
	"certificate", "cert", "key", "secret", "credential", "session",

	// Roles & Permissions
	"guest", "anonymous", "anon", "user", "users", "member", "members",
	"subscriber", "editor", "author", "contributor", "reviewer", "auditor",
	"analyst", "developer", "dev", "devops", "engineer", "architect",
	"designer", "tester", "qa", "billing", "finance", "legal", "hr",
	"sales", "marketing", "ceo", "cto", "cfo", "coo", "founder", "cofounder",

	// Common Reserved
	"account", "accounts", "profile", "profiles", "settings", "config",
	"configuration", "dashboard", "panel", "console", "portal", "home",
	"index", "main", "default", "null", "nil", "undefined", "void",
	"true", "false", "test", "testing", "debug", "demo", "example",
	"sample", "temp", "temporary", "tmp", "backup", "archive", "log",
	"logs", "audit", "report", "reports", "analytics", "stats", "status",

	// API & Endpoints
	"rest", "graphql", "grpc", "websocket", "ws", "wss", "http",
	"https", "endpoint", "endpoints", "route", "routes", "path", "url",
	"uri", "hook", "hooks", "event", "events", "stream",

	// Content & Media
	"blog", "news", "article", "articles", "post", "posts", "page", "pages",
	"feed", "rss", "atom", "sitemap", "robots", "favicon", "static",
	"assets", "images", "image", "img", "media", "upload", "uploads",
	"download", "downloads", "file", "files", "document", "documents",

	// Communication
	"contact", "message", "messages", "chat", "notification", "notifications",
	"alert", "alerts", "inbox", "outbox", "sent", "draft", "drafts",
	"spam", "abuse", "flag", "block", "mute", "ban",

	// Commerce & Billing
	"shop", "store", "cart", "checkout", "order", "orders", "invoice",
	"invoices", "payment", "payments", "subscription", "subscriptions",
	"plan", "plans", "pricing", "refund", "coupon", "discount",

	// Social Features
	"follow", "follower", "followers", "following", "friend", "friends",
	"like", "likes", "share", "shares", "comment", "comments", "reply",
	"mention", "mentions", "tag", "tags", "group", "groups", "team", "teams",
	"community", "communities", "forum", "forums", "channel", "channels",

	// Brand & Legal
	"official", "verified", "trusted", "partner", "affiliate", "sponsor",
	"brand", "trademark", "copyright", "terms", "privacy",
	"policy", "policies", "tos", "eula", "gdpr", "dmca",

	// Numbers & Special
	"0", "1", "123", "1234", "12345", "000", "111", "666", "911", "420", "69",

	// Common Spam Patterns
	"info", "noreply", "no-reply", "donotreply", "mailer", "postmaster",
	"webmaster", "hostmaster", "junk", "trash",

	// Project-specific
	"caspaste", "casjay-forks", "casjay",
}

// Critical terms that block substrings
var criticalBlockedTerms = []string{
	"admin", "root", "system", "mod", "official", "verified",
}

// ValidateUsername validates a username per PART 34 rules
func ValidateUsername(username string) error {
	username = strings.ToLower(strings.TrimSpace(username))

	// Length checks
	if len(username) < UsernameMinLength {
		return errors.New("username must be at least 3 characters")
	}
	if len(username) > UsernameMaxLength {
		return errors.New("username cannot exceed 32 characters")
	}

	// Format check
	if !usernameRegex.MatchString(username) {
		return errors.New("username can only contain lowercase letters, numbers, underscore, and hyphen")
	}

	// Check for consecutive special chars
	if consecutiveChars.MatchString(username) {
		return errors.New("username cannot have consecutive _, -, _-, or -_")
	}

	// Cannot end with _ or -
	if strings.HasSuffix(username, "_") || strings.HasSuffix(username, "-") {
		return errors.New("username cannot end with _ or -")
	}

	// Check blocklist
	if err := checkUsernameBlocklist(username); err != nil {
		return err
	}

	return nil
}

// checkUsernameBlocklist checks if username contains blocked words
func checkUsernameBlocklist(username string) error {
	username = strings.ToLower(username)

	// Exact match check
	for _, blocked := range UsernameBlocklist {
		if username == blocked {
			return ErrUsernameBlocked
		}
	}

	// Critical terms substring check
	for _, term := range criticalBlockedTerms {
		if strings.Contains(username, term) {
			return ErrUsernameBlocked
		}
	}

	return nil
}

// ValidateEmail validates an email address per PART 34 rules
func ValidateEmail(email string) error {
	email = strings.TrimSpace(strings.ToLower(email))

	// Length checks
	if len(email) > EmailMaxLength {
		return errors.New("email too long (max 254 characters)")
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return errors.New("invalid email format")
	}

	local, domain := parts[0], parts[1]

	// Local part checks
	if len(local) == 0 || len(local) > 64 {
		return errors.New("invalid local part length")
	}
	if strings.HasPrefix(local, ".") || strings.HasSuffix(local, ".") {
		return errors.New("local part cannot start or end with dot")
	}
	if strings.Contains(local, "..") {
		return errors.New("local part cannot have consecutive dots")
	}

	// Domain checks
	if len(domain) == 0 || len(domain) > 255 {
		return errors.New("invalid domain length")
	}
	if !strings.Contains(domain, ".") {
		return errors.New("domain must have valid TLD")
	}

	// Regex validation
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}

	return nil
}

// ValidatePassword validates a password per PART 34 rules
func ValidatePassword(password string) error {
	if len(password) < PasswordMinLength {
		return errors.New("password must be at least 8 characters")
	}
	if len(password) > PasswordMaxLength {
		return errors.New("password too long")
	}
	return nil
}

// ValidatePasswordStrength validates password against configurable requirements
func ValidatePasswordStrength(password string, requireUppercase, requireNumber, requireSpecial bool) error {
	if err := ValidatePassword(password); err != nil {
		return err
	}

	if requireUppercase {
		if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
			return errors.New("password must contain at least one uppercase letter")
		}
	}

	if requireNumber {
		if !regexp.MustCompile(`[0-9]`).MatchString(password) {
			return errors.New("password must contain at least one number")
		}
	}

	if requireSpecial {
		if !regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password) {
			return errors.New("password must contain at least one special character")
		}
	}

	return nil
}

// ValidateBio validates a user bio
func ValidateBio(bio string) error {
	if len(bio) > BioMaxLength {
		return errors.New("bio cannot exceed 500 characters")
	}
	return nil
}

// ValidateWebsite validates a website URL
func ValidateWebsite(website string) error {
	if website == "" {
		return nil
	}

	// Must start with http:// or https://
	if !strings.HasPrefix(website, "http://") && !strings.HasPrefix(website, "https://") {
		return errors.New("website must start with http:// or https://")
	}

	return nil
}

// ValidateVisibility validates a visibility setting
func ValidateVisibility(visibility string) error {
	switch visibility {
	case VisibilityPublic, VisibilityPrivate:
		return nil
	default:
		return errors.New("visibility must be 'public' or 'private'")
	}
}

// ValidateAvatarType validates an avatar type
func ValidateAvatarType(avatarType string) error {
	switch avatarType {
	case AvatarTypeGravatar, AvatarTypeUpload, AvatarTypeURL:
		return nil
	default:
		return errors.New("avatar type must be 'gravatar', 'upload', or 'url'")
	}
}

// NormalizeUsername normalizes a username for storage
func NormalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

// NormalizeEmail normalizes an email for storage
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
