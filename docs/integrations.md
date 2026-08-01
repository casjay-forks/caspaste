# Integrations

This page documents the integration surfaces CasPaste actually ships.
Features described in generic templates but not present in this codebase
(SAML, WebFinger, OpenID Provider Metadata, Android App Links, Apple
Universal Links, MTA-STS, federation) are intentionally omitted.

## External Identity

There is no working external identity integration. The admin panel ships
configuration pages for OIDC and LDAP
(`/server/{admin_path}/config/security/auth/oidc` and `.../ldap`), but
they only render "Not configured" placeholders — no OAuth2/OIDC client or
LDAP client library is wired into the codebase (`go.mod` has no such
dependency). Do not configure external identity providers expecting them
to work; only local, session-based auth under `/server/auth/*` is
functional. See `docs/security.md` for details.

## Discovery & Protocol Endpoints

CasPaste exposes machine-readable API documentation directly, without a
`/server/docs/*` prefix:

| Endpoint | Purpose |
|----------|---------|
| `/openapi` | Swagger UI for the REST API |
| `/openapi.json` | Raw OpenAPI spec |
| `/graphql` | GraphiQL UI on `GET`, GraphQL queries on `POST` |

There are no `/api/swagger`, `/api/graphql`, or `/api/healthz` aliases —
use `/openapi`, `/graphql`, and `/api/v1/server/healthz` respectively.

## External Paste-Service Compatibility

CasPaste's most significant integration surface is a set of drop-in
compatibility shims (`src/compat/`) that let existing clients built for
other paste services talk to CasPaste unmodified. Mode selection order
(first match wins, per `src/compat/compat.go`):

1. `CASPASTE_API_MODE` environment variable, set once at startup
2. `Host` header pattern matching, evaluated per request
3. Native CasPaste API (default, no interception)

By default the `Host` header takes precedence over `CASPASTE_API_MODE`;
set `CASPASTE_FORCE_HOST=no` to reverse that.

| Service | Mode detection |
|---------|----------------|
| sprunge.us | Always active — `POST /sprunge` |
| ix.io | Always active — `POST /ix` |
| termbin.com | Always active — `POST /termbin` |
| hastebin | `Host: haste.*` or `CASPASTE_API_MODE=hastebin` |
| pastebin.com | `Host: pb.*` or `CASPASTE_API_MODE=pastebin` |
| Stikked | `Host: sk.*` or `CASPASTE_API_MODE=stikked` |
| Microbin | `Host: mb.*` or `CASPASTE_API_MODE=microbin` |
| Lenpaste | `Host: lp.*` or `CASPASTE_API_MODE=lenpaste` |

All compat endpoints validate input and are rate-limited the same as the
native paste-creation path. See `src/compat/{lenpaste,stikked,microbin,
hastebin,pastebin,termbin}.go` for the per-service handler implementations.

## Operator Notes

- The compatibility shims require no configuration to use the
  always-active modes (sprunge, ix, termbin); host-based or
  `CASPASTE_API_MODE`-based modes require either DNS/reverse-proxy routing
  for the relevant hostname pattern or setting the environment variable.
- OIDC/LDAP admin pages are visible but non-functional; do not document
  them to end users as working login options.
- `/openapi` and `/graphql` are public endpoints (see `docs/security.md`)
  and require no authentication to view documentation or introspect the
  schema.
