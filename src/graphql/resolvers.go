// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

// GraphQL resolvers per AI.md PART 14
// Resolvers connect GraphQL queries to data sources
package graphql

import (
	"errors"

	"github.com/casjay-forks/caspaste/src/storage"
)

// Resolvers provides data resolution for GraphQL queries
type Resolvers struct {
	db          *storage.DB
	version     string
	title       string
	public      bool
	maxBodyLen  int
	maxTitleLen int
	lexers      []string
}

// ResolversConfig holds configuration for resolvers
type ResolversConfig struct {
	DB          *storage.DB
	Version     string
	Title       string
	Public      bool
	MaxBodyLen  int
	MaxTitleLen int
	Lexers      []string
}

// NewResolvers creates a new resolvers instance
func NewResolvers(cfg *ResolversConfig) *Resolvers {
	return &Resolvers{
		db:          cfg.DB,
		version:     cfg.Version,
		title:       cfg.Title,
		public:      cfg.Public,
		maxBodyLen:  cfg.MaxBodyLen,
		maxTitleLen: cfg.MaxTitleLen,
		lexers:      cfg.Lexers,
	}
}

// HealthResult represents health check result
type HealthResult struct {
	OK      bool   `json:"ok"`
	Status  string `json:"status"`
	Version string `json:"version"`
}

// ServerInfoResult represents server info result
type ServerInfoResult struct {
	Version     string   `json:"version"`
	Title       string   `json:"title"`
	Public      bool     `json:"public"`
	MaxBodyLen  int      `json:"maxBodyLen"`
	MaxTitleLen int      `json:"maxTitleLen"`
	Lexers      []string `json:"lexers"`
}

// PasteResult represents a paste result
type PasteResult struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	Syntax      string `json:"syntax"`
	CreateTime  int64  `json:"createTime"`
	DeleteTime  int64  `json:"deleteTime"`
	OneUse      bool   `json:"oneUse"`
	IsPrivate   bool   `json:"isPrivate"`
	IsFile      bool   `json:"isFile"`
	FileName    string `json:"fileName"`
	MimeType    string `json:"mimeType"`
	IsURL       bool   `json:"isUrl"`
	OriginalURL string `json:"originalUrl"`
}

// PasteSummaryResult represents a paste summary
type PasteSummaryResult struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Syntax     string `json:"syntax"`
	CreateTime int64  `json:"createTime"`
}

// CreatePasteResult represents create paste result
type CreatePasteResult struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// ResolveHealth resolves the healthz query
func (r *Resolvers) ResolveHealth() *HealthResult {
	return &HealthResult{
		OK:      true,
		Status:  "healthy",
		Version: r.version,
	}
}

// ResolveServerInfo resolves the serverInfo query
func (r *Resolvers) ResolveServerInfo() *ServerInfoResult {
	return &ServerInfoResult{
		Version:     r.version,
		Title:       r.title,
		Public:      r.public,
		MaxBodyLen:  r.maxBodyLen,
		MaxTitleLen: r.maxTitleLen,
		Lexers:      r.lexers,
	}
}

// ResolvePaste resolves a single paste by ID
func (r *Resolvers) ResolvePaste(id string) (*PasteResult, error) {
	if r.db == nil {
		return nil, errors.New("database not available")
	}

	paste, err := r.db.PasteGet(id)
	if err != nil {
		return nil, errors.New("paste not found")
	}

	return &PasteResult{
		ID:          paste.ID,
		Title:       paste.Title,
		Body:        paste.Body,
		Syntax:      paste.Syntax,
		CreateTime:  paste.CreateTime,
		DeleteTime:  paste.DeleteTime,
		OneUse:      paste.OneUse,
		IsPrivate:   paste.IsPrivate,
		IsFile:      paste.IsFile,
		FileName:    paste.FileName,
		MimeType:    paste.MimeType,
		IsURL:       paste.IsURL,
		OriginalURL: paste.OriginalURL,
	}, nil
}

// ResolvePastes resolves the list of public pastes
func (r *Resolvers) ResolvePastes() ([]*PasteSummaryResult, error) {
	if r.db == nil {
		return nil, errors.New("database not available")
	}

	// PasteList(limit, offset)
	pastes, err := r.db.PasteList(100, 0)
	if err != nil {
		return nil, err
	}

	results := make([]*PasteSummaryResult, len(pastes))
	for i, p := range pastes {
		results[i] = &PasteSummaryResult{
			ID:         p.ID,
			Title:      p.Title,
			Syntax:     p.Syntax,
			CreateTime: p.CreateTime,
		}
	}

	return results, nil
}

// ResolveCreatePaste creates a new paste
func (r *Resolvers) ResolveCreatePaste(input map[string]interface{}) (*CreatePasteResult, error) {
	if r.db == nil {
		return nil, errors.New("database not available")
	}

	paste := storage.Paste{}

	if title, ok := input["title"].(string); ok {
		paste.Title = title
	}
	if body, ok := input["body"].(string); ok {
		paste.Body = body
	} else {
		return nil, errors.New("body is required")
	}
	if syntax, ok := input["syntax"].(string); ok {
		paste.Syntax = syntax
	}
	if oneUse, ok := input["oneUse"].(bool); ok {
		paste.OneUse = oneUse
	}
	if isPrivate, ok := input["isPrivate"].(bool); ok {
		paste.IsPrivate = isPrivate
	}

	// PasteAdd returns (id, createTime, deleteTime, error)
	id, _, _, err := r.db.PasteAdd(paste)
	if err != nil {
		return nil, err
	}

	return &CreatePasteResult{
		ID:  id,
		URL: "/" + id,
	}, nil
}
