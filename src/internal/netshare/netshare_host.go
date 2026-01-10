// Copyright (C) 2021-2023 Leonid Maslakov.

// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE file for details.

package netshare

import (
	"net"
	"net/http"
	"strings"
)

func GetHost(req *http.Request) string {
	// Try RFC 7239 Forwarded header first
	forwarded := req.Header.Get("Forwarded")
	if forwarded != "" {
		// Parse "Forwarded: for=192.0.2.60;proto=http;host=example.com"
		parts := strings.Split(forwarded, ";")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if strings.HasPrefix(part, "host=") {
				host := strings.TrimPrefix(part, "host=")
				host = strings.Trim(host, "\"")
				if host != "" {
					return host
				}
			}
		}
	}

	// X-Forwarded-Host (common reverse proxy header)
	xHost := req.Header.Get("X-Forwarded-Host")
	if xHost != "" {
		// Take the first host if multiple are listed
		return strings.Split(xHost, ",")[0]
	}

	// Fallback to request Host header
	return req.Host
}

func GetProtocol(req *http.Request) string {
	// Try RFC 7239 Forwarded header first
	forwarded := req.Header.Get("Forwarded")
	if forwarded != "" {
		// Parse "Forwarded: for=192.0.2.60;proto=http;host=example.com"
		parts := strings.Split(forwarded, ";")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if strings.HasPrefix(part, "proto=") {
				proto := strings.TrimPrefix(part, "proto=")
				proto = strings.Trim(proto, "\"")
				if proto != "" {
					return proto
				}
			}
		}
	}

	// X-Forwarded-Proto (common reverse proxy header)
	xProto := req.Header.Get("X-Forwarded-Proto")
	if xProto != "" {
		// Take the first protocol if multiple are listed
		return strings.Split(xProto, ",")[0]
	}

	// X-Forwarded-Ssl (some proxies use this: on/off)
	xSsl := req.Header.Get("X-Forwarded-Ssl")
	if xSsl == "on" {
		return "https"
	}

	// X-Scheme (alternative header used by some proxies)
	xScheme := req.Header.Get("X-Scheme")
	if xScheme != "" {
		return xScheme
	}

	// Check if request came over TLS
	if req.TLS != nil {
		return "https"
	}

	// Check URL scheme if available
	if req.URL.Scheme != "" {
		return req.URL.Scheme
	}

	// Default to http
	return "http"
}

// GetClientAddrTrusted extracts client IP address from request
// Set trustProxy=true only when behind a trusted reverse proxy
// WARNING: trustProxy=true without a proxy allows IP spoofing for rate limit bypass
func GetClientAddrTrusted(req *http.Request, trustProxy bool) net.IP {
	// If not trusting proxy headers, use direct connection IP only
	if !trustProxy {
		host, _, err := net.SplitHostPort(req.RemoteAddr)
		if err != nil {
			return nil
		}
		return net.ParseIP(host)
	}

	// Try RFC 7239 Forwarded header first
	forwarded := req.Header.Get("Forwarded")
	if forwarded != "" {
		// Parse "Forwarded: for=192.0.2.60;proto=http;host=example.com"
		parts := strings.Split(forwarded, ";")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if strings.HasPrefix(part, "for=") {
				forVal := strings.TrimPrefix(part, "for=")
				forVal = strings.Trim(forVal, "\"")
				// Remove port if present (e.g., "192.0.2.60:47011" or "[2001:db8::1]:47011")
				if strings.Contains(forVal, "]:") {
					// IPv6 with port: [2001:db8::1]:47011
					forVal = strings.TrimPrefix(forVal, "[")
					forVal = strings.Split(forVal, "]:")[0]
				} else if strings.Contains(forVal, ":") && strings.Count(forVal, ":") == 1 {
					// IPv4 with port: 192.0.2.60:47011
					forVal = strings.Split(forVal, ":")[0]
				}
				if ip := net.ParseIP(forVal); ip != nil {
					return ip
				}
			}
		}
	}

	// X-Real-IP (common in nginx)
	xReal := req.Header.Get("X-Real-IP")
	if xReal != "" {
		if ip := net.ParseIP(strings.TrimSpace(xReal)); ip != nil {
			return ip
		}
	}

	// X-Forwarded-For (common in many proxies, takes the first IP)
	xFor := req.Header.Get("X-Forwarded-For")
	if xFor != "" {
		// Take the first IP from the list
		firstIP := strings.TrimSpace(strings.Split(xFor, ",")[0])
		if ip := net.ParseIP(firstIP); ip != nil {
			return ip
		}
	}

	// CF-Connecting-IP (Cloudflare specific)
	cfIP := req.Header.Get("CF-Connecting-IP")
	if cfIP != "" {
		if ip := net.ParseIP(strings.TrimSpace(cfIP)); ip != nil {
			return ip
		}
	}

	// True-Client-IP (Akamai and Cloudflare)
	trueIP := req.Header.Get("True-Client-IP")
	if trueIP != "" {
		if ip := net.ParseIP(strings.TrimSpace(trueIP)); ip != nil {
			return ip
		}
	}

	// Fallback: use real client address from connection
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return nil
	}

	return net.ParseIP(host)
}

// GetClientAddr extracts client IP address using direct connection only
// This is the safe default - use GetClientAddrTrusted when behind a reverse proxy
func GetClientAddr(req *http.Request) net.IP {
	return GetClientAddrTrusted(req, false)
}
