
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package netshare

import (
	"errors"
)

const (
	// Max length for paste author name, email and URL
	MaxLengthAuthorAll = 100
)

var (
	// HTTP 400
	ErrBadRequest = errors.New("Bad Request")
	// HTTP 401
	ErrUnauthorized = errors.New("Unauthorized")
	// HTTP 404
	ErrNotFound = errors.New("Not Found")
	// HTTP 405
	ErrMethodNotAllowed = errors.New("Method Not Allowed")
	// HTTP 413
	ErrPayloadTooLarge = errors.New("Payload Too Large")
	// HTTP 429
	ErrTooManyRequests = errors.New("Too Many Requests")
	// HTTP 500
	ErrInternal = errors.New("Internal Server Error")
)

type RateLimitError struct {
	s          string
	RetryAfter int64
}

func (e *RateLimitError) Error() string {
	return e.s
}

func ErrTooManyRequestsNew(retryAfter int64) *RateLimitError {
	return &RateLimitError{
		s:          "Too Many Requests",
		RetryAfter: retryAfter,
	}
}
