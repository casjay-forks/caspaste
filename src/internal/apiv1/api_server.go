
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package apiv1

import (
	"encoding/json"
	"github.com/casjay-forks/caspaste/src/internal/netshare"
	"net/http"
)

type serverInfoType struct {
	Software          string   `json:"software"`
	Version           string   `json:"version"`
	TitleMaxLen       int      `json:"titleMaxlength"`
	BodyMaxLen        int      `json:"bodyMaxlength"`
	MaxLifeTime       int64    `json:"maxLifeTime"`
	ServerAbout       string   `json:"serverAbout"`
	ServerRules       string   `json:"serverRules"`
	ServerTermsOfUse  string   `json:"serverTermsOfUse"`
	AdminName         string   `json:"adminName"`
	AdminMail         string   `json:"adminMail"`
	Syntaxes          []string `json:"syntaxes"`
	UiDefaultLifeTime string   `json:"uiDefaultLifeTime"`
	AuthRequired      bool     `json:"authRequired"`
}

// GET /api/v1/getServerInfo
func (data *Data) getServerInfoHand(rw http.ResponseWriter, req *http.Request) error {
	// Check method
	if req.Method != "GET" {
		return netshare.ErrMethodNotAllowed
	}

	// Prepare data
	serverInfo := serverInfoType{
		Software:          "CasPaste",
		Version:           data.Version,
		TitleMaxLen:       data.TitleMaxLen,
		BodyMaxLen:        data.BodyMaxLen,
		MaxLifeTime:       data.MaxLifeTime,
		ServerAbout:       data.ServerAbout,
		ServerRules:       data.ServerRules,
		ServerTermsOfUse:  data.ServerTermsOfUse,
		AdminName:         data.AdminName,
		AdminMail:         data.AdminMail,
		Syntaxes:          data.Lexers,
		UiDefaultLifeTime: data.UiDefaultLifeTime,
		AuthRequired:      data.CasPasswdFile != "",
	}

	// Return response
	rw.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(serverInfo)
}
