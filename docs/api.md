# API Documentation

CasPaste provides a RESTful API for programmatic access.

## Base URL

```
https://your-instance.com/api/v1/
```

## Authentication

Most endpoints are public. Private instances require authentication via session cookie or API token.

## Endpoints

### Create Paste

**POST** `/api/v1/new`

Create a new paste or upload a file.

#### Text Paste

```bash
curl -X POST https://paste.example.com/api/v1/new \
  -d "body=Hello World" \
  -d "syntax=plaintext"
```

#### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `body` | string | Yes* | Paste content (*or file) |
| `syntax` | string | No | Syntax highlighting language (default: plaintext) |
| `title` | string | No | Paste title (max 120 chars) |
| `expires` | string | No | Expiration: `never`, `10m`, `1h`, `1d`, `1w`, `1M` |
| `oneUse` | boolean | No | Burn after reading |
| `password` | string | No | Password protection |

#### File Upload

```bash
curl -X POST https://paste.example.com/api/v1/new \
  -F "file=@image.png"
```

#### URL Shortener

```bash
curl -X POST https://paste.example.com/api/v1/new \
  -d "url=true" \
  -d "originalURL=https://example.com/long/url"
```

#### Response

```json
{
  "id": "abc123",
  "url": "https://paste.example.com/abc123",
  "deleteToken": "del_xyz789"
}
```

### Get Paste

**GET** `/api/v1/get/{id}`

Retrieve a paste by ID.

```bash
curl https://paste.example.com/api/v1/get/abc123
```

#### Response

```json
{
  "id": "abc123",
  "title": "My Paste",
  "body": "Hello World",
  "syntax": "plaintext",
  "created": "2024-01-15T10:30:00Z",
  "expires": null,
  "views": 5
}
```

### List Pastes

**GET** `/api/v1/list`

List recent pastes (public mode only).

```bash
curl https://paste.example.com/api/v1/list
```

#### Query Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `limit` | int | 20 | Max results (1-100) |
| `offset` | int | 0 | Pagination offset |

#### Response

```json
{
  "pastes": [
    {
      "id": "abc123",
      "title": "My Paste",
      "syntax": "python",
      "created": "2024-01-15T10:30:00Z",
      "views": 5
    }
  ],
  "total": 100,
  "limit": 20,
  "offset": 0
}
```

### Server Info

**GET** `/api/v1/getServerInfo`

Get server metadata and capabilities.

```bash
curl https://paste.example.com/api/v1/getServerInfo
```

#### Response

```json
{
  "version": "1.0.0",
  "name": "CasPaste",
  "public": true,
  "features": {
    "fileUpload": true,
    "urlShortener": true,
    "burnAfterReading": true,
    "passwordProtection": true
  },
  "limits": {
    "maxPasteSize": 0,
    "maxTitleLength": 120
  }
}
```

### Health Check

**GET** `/api/v1/healthz`

Health check endpoint (always JSON).

```bash
curl https://paste.example.com/api/v1/healthz
```

#### Response

```json
{
  "status": "ok",
  "version": "1.0.0",
  "database": "ok"
}
```

## Frontend Health Check

**GET** `/healthz`

Smart content negotiation:

- **Browser:** HTML response
- **CLI (curl):** Formatted text
- **API client:** JSON (with `Accept: application/json`)

## Error Responses

All errors return JSON:

```json
{
  "error": "paste not found",
  "code": "NOT_FOUND",
  "status": 404
}
```

### Error Codes

| Code | Status | Description |
|------|--------|-------------|
| `NOT_FOUND` | 404 | Paste not found |
| `INVALID_INPUT` | 400 | Invalid request parameters |
| `RATE_LIMITED` | 429 | Rate limit exceeded |
| `UNAUTHORIZED` | 401 | Authentication required |
| `FORBIDDEN` | 403 | Access denied |
| `SERVER_ERROR` | 500 | Internal server error |

## Rate Limiting

Endpoints have configurable rate limits:

| Window | Default Limit |
|--------|---------------|
| 5 minutes | 100 requests |
| 15 minutes | 300 requests |
| 1 hour | 1000 requests |

Rate limit headers are included in responses:

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1705312200
```
