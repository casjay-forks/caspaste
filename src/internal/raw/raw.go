// Copyright (C) 2021-2023 Leonid Maslakov.

// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE file for details.

package raw

import (
	"github.com/casjay-forks/caspaste/src/internal/config"
	"github.com/casjay-forks/caspaste/src/internal/logger"
	"github.com/casjay-forks/caspaste/src/internal/netshare"
	"github.com/casjay-forks/caspaste/src/internal/storage"
	"net/http"
)

type Data struct {
	DB  storage.DB
	Log logger.Logger

	RateLimitGet *netshare.RateLimitSystem

	Version string
}

func Load(db storage.DB, cfg config.Config) *Data {
	return &Data{
		DB:           db,
		Log:          cfg.Log,
		RateLimitGet: cfg.RateLimitGet,
		Version:      cfg.Version,
	}
}

func (data *Data) Hand(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Server", config.Software+"/"+data.Version)

	err := data.rawHand(rw, req)

	if err == nil {
		data.Log.HttpRequest(req, 200)

	} else {
		code, err := data.writeError(rw, req, err)
		if err != nil {
			data.Log.HttpError(req, err)
		} else {
			data.Log.HttpRequest(req, code)
		}
	}
}
