// Copyright (C) 2021-2023 Leonid Maslakov.

// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE file for details.

package apiv1

import (
	"github.com/casjay-forks/caspaste/src/internal/netshare"
	"net/http"
)

// GET /api/v1/
func (data *Data) MainHand(rw http.ResponseWriter, req *http.Request) {
	data.writeError(rw, req, netshare.ErrNotFound)
}
