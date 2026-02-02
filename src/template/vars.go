// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package template

import (
	"strings"
)

// Variables holds all available variables for content replacement
type Variables struct {
	FQDN                 string
	Version              string
	Protocol             string
	Port                 string
	ServerTitle          string
	ServerAdminName      string
	ServerAdminEmail     string
	ServerAdminFrom      string
	SecurityContactEmail string
	SecurityContactName  string
}

// ReplaceVariables performs simple string replacement of variables in content
// Variables are in {key} format, e.g., {fqdn}, {server.title}
func ReplaceVariables(content string, vars Variables) string {
	replacements := map[string]string{
		"{fqdn}":                       vars.FQDN,
		"{version}":                    vars.Version,
		"{protocol}":                   vars.Protocol,
		"{port}":                       vars.Port,
		"{server.title}":               vars.ServerTitle,
		"{server.administrator.name}":  vars.ServerAdminName,
		"{server.administrator.email}": vars.ServerAdminEmail,
		"{server.administrator.from}":  vars.ServerAdminFrom,
		"{web.security.contact.email}": vars.SecurityContactEmail,
		"{web.security.contact.name}":  vars.SecurityContactName,
	}

	result := content
	for placeholder, value := range replacements {
		result = strings.ReplaceAll(result, placeholder, value)
	}

	return result
}
