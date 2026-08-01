# Security

This page documents CasPaste's actual security-relevant endpoints and
behavior as shipped — not a generic feature list. Anything not mentioned
here (SAML, WebFinger, OpenID Provider Metadata, App Links, Apple
association, MTA-STS) is not implemented in this codebase.

## Authentication & Identity

- Local, session-based authentication lives under `/server/auth/*`
  (`/server/auth/login`, `/server/auth/logout`, and related password
  routes). This is the only working authentication backend.
- The admin panel has configuration *pages* at
  `/server/{admin_path}/config/security/auth/oidc` and
  `/server/{admin_path}/config/security/auth/ldap`, but as shipped these
  render static "Not configured" placeholders
  (`src/admin/handlers.go`: `handleSecurityOIDC`, `handleSecurityLDAP`).
  There is no OIDC or LDAP client library in `go.mod`, and no discovery,
  token-exchange, or bind logic anywhere in the codebase. Treat these pages
  as reserved UI for a future release, not a working integration.
- There is no SAML support (metadata endpoint, SLO, or otherwise).

## Public Security Endpoints

These paths are reachable without authentication even when the instance
requires login (see `IsPublicPath` in `src/web/auth.go`):

| Path | Purpose |
|------|---------|
| `/.well-known/security.txt` | RFC 9116 security contact info (see below) |
| `/.well-known/change-password` | [RFC 8615](https://www.rfc-editor.org/rfc/rfc8615) well-known redirect to the password-change flow |
| `/healthz`, `/server/healthz` | Content-negotiated health check (HTML/JSON/text) |
| `/api/v1/server/healthz` | JSON health check for scripts/monitoring |
| `/server/about/security` | Human-readable security overview, linked as the `Policy` field in security.txt |
| `/openapi`, `/openapi.json` | OpenAPI/Swagger UI and spec |

There is no `/.well-known/pgp-key.asc` and no `/server/security` route
distinct from `/server/about/security`.

## Security Reporting

`/.well-known/security.txt` is served by `src/web/security.go` and is
either an operator-supplied override (`SecurityTxt` config value) or
auto-generated from `SecurityContactEmail`, `SecurityContactName`, and the
instance's FQDN, per RFC 9116. The generated file includes:

- `Contact: mailto:{security contact email}`
- `Expires:` one year from generation time
- `Acknowledgments:` linking to `/server/about/authors`
- `Canonical:` the file's own URL
- `Policy:` linking to `/server/about/security`
- `Preferred-Languages: en`

Security researchers should email the configured contact address, or read
`/server/about/security` for the human-readable policy page. There is no
dedicated `/server/contact` endpoint or security-report web form.

## Well-Known Namespace

Only two entries under `/.well-known/` are implemented:

- `/.well-known/security.txt`
- `/.well-known/change-password`

`/.well-known/acme-challenge/{token}` also exists (`src/ssl/acme.go`) but
is an ACME/TLS provisioning mechanism, not a security-disclosure or
identity endpoint. Any other `/.well-known/*` path returns `404`. WebFinger,
OpenID Provider Metadata, Android App Links
(`/.well-known/assetlinks.json`), Apple association
(`/.well-known/apple-app-site-association`), and MTA-STS are not present
in this codebase.
