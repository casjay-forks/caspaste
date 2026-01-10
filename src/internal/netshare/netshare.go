// Copyright (C) 2021-2023 Leonid Maslakov.

// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE file for details.

package netshare

import (
	"errors"
)

const (
	MaxLengthAuthorAll = 100 // Max length or paste author name, email and URL.
)

var (
	ErrBadRequest       = errors.New("Bad Request")        // 400
	ErrUnauthorized     = errors.New("Unauthorized")       // 401
	ErrNotFound         = errors.New("Not Found")          // 404
	ErrMethodNotAllowed = errors.New("Method Not Allowed") // 405
	ErrPayloadTooLarge  = errors.New("Payload Too Large")  // 413
	ErrTooManyRequests  = errors.New("Too Many Requests")  // 429
	ErrInternal         = errors.New("Internal Server Error") // 500
)

type ErrTooManyRequests struct {
	s          string
	RetryAfter int64
}

func (e *ErrTooManyRequests) Error() string {
	return e.s
}

func ErrTooManyRequestsNew(retryAfter int64) *ErrTooManyRequests {
	return &ErrTooManyRequests{
		s:          "Too Many Requests",
		RetryAfter: retryAfter,
	}
}
