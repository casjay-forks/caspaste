
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package netshare

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/casjay-forks/caspaste/src/internal/lineend"
	"github.com/casjay-forks/caspaste/src/internal/storage"
)

func PasteAddFromForm(req *http.Request, db storage.DB, rateSys *RateLimitSystem, titleMaxLen int, bodyMaxLen int, maxLifeTime int64, lexerNames []string) (string, int64, int64, error) {
	// Check HTTP method
	if req.Method != "POST" {
		return "", 0, 0, ErrMethodNotAllowed
	}

	// Check rate limit
	err := rateSys.CheckAndUse(GetClientAddr(req))
	if err != nil {
		return "", 0, 0, err
	}

	// Read form
	req.ParseMultipartForm(52428800) // 50MB max

	paste := storage.Paste{
		Title:       req.PostFormValue("title"),
		Body:        req.PostFormValue("body"),
		Syntax:      req.PostFormValue("syntax"),
		DeleteTime:  0,
		OneUse:      false,
		Author:      req.PostFormValue("author"),
		AuthorEmail: req.PostFormValue("authorEmail"),
		AuthorURL:   req.PostFormValue("authorURL"),
		IsEditable:  req.PostFormValue("editable") == "true",
		IsPrivate:   req.PostFormValue("private") == "true",
		IsURL:       req.PostFormValue("url") == "true",
		OriginalURL: req.PostFormValue("originalURL"),
	}

	// Handle file upload
	file, handler, err := req.FormFile("file")
	if err == nil {
		defer file.Close()

		// Read file contents
		fileData, err := io.ReadAll(file)
		if err != nil {
			return "", 0, 0, err
		}

		// Set file fields
		paste.IsFile = true
		paste.FileName = handler.Filename
		paste.MimeType = handler.Header.Get("Content-Type")
		if paste.MimeType == "" {
			paste.MimeType = "application/octet-stream"
		}

		// Store file data as base64 in Body field
		paste.Body = string(fileData)

		// Default syntax for files
		if paste.Syntax == "" {
			paste.Syntax = "file"
		}
	}

	// Remove new line from title
	paste.Title = strings.Replace(paste.Title, "\n", "", -1)
	paste.Title = strings.Replace(paste.Title, "\r", "", -1)
	paste.Title = strings.Replace(paste.Title, "\t", " ", -1)

	// Check title
	if utf8.RuneCountInString(paste.Title) > titleMaxLen && titleMaxLen >= 0 {
		return "", 0, 0, ErrPayloadTooLarge
	}

	// Check paste body
	if paste.Body == "" {
		return "", 0, 0, ErrBadRequest
	}

	if utf8.RuneCountInString(paste.Body) > bodyMaxLen && bodyMaxLen > 0 {
		return "", 0, 0, ErrPayloadTooLarge
	}

	// Change paste body lines end
	switch req.PostForm.Get("lineEnd") {
	case "", "LF", "lf":
		paste.Body = lineend.UnknownToUnix(paste.Body)

	case "CRLF", "crlf":
		paste.Body = lineend.UnknownToDos(paste.Body)

	case "CR", "cr":
		paste.Body = lineend.UnknownToOldMac(paste.Body)

	default:
		return "", 0, 0, ErrBadRequest
	}

	// Check syntax
	if paste.Syntax == "" {
		paste.Syntax = "plaintext"
	}

	// Validate syntax (allow "autodetect" as special value)
	syntaxOk := false
	if paste.Syntax == "autodetect" {
		syntaxOk = true
		// Leave as "autodetect" - will be detected during display
	} else {
		for _, name := range lexerNames {
			if name == paste.Syntax {
				syntaxOk = true
				break
			}
		}
	}

	if syntaxOk == false {
		return "", 0, 0, ErrBadRequest
	}

	// Get delete time
	expirStr := req.PostForm.Get("expiration")
	if expirStr != "" {
		// Convert string to int
		expir, err := strconv.ParseInt(expirStr, 10, 64)
		if err != nil {
			return "", 0, 0, ErrBadRequest
		}

		// Check limits
		if maxLifeTime > 0 {
			if expir > maxLifeTime || expir <= 0 {
				return "", 0, 0, ErrBadRequest
			}
		}

		// Save if ok
		if expir > 0 {
			paste.DeleteTime = time.Now().Unix() + expir
		}
	}

	// Get "one use" parameter
	if req.PostForm.Get("oneUse") == "true" {
		paste.OneUse = true
	}

	// Check author name, email and URL length.
	if utf8.RuneCountInString(paste.Author) > MaxLengthAuthorAll {
		return "", 0, 0, ErrPayloadTooLarge
	}

	if utf8.RuneCountInString(paste.AuthorEmail) > MaxLengthAuthorAll {
		return "", 0, 0, ErrPayloadTooLarge
	}

	if utf8.RuneCountInString(paste.AuthorURL) > MaxLengthAuthorAll {
		return "", 0, 0, ErrPayloadTooLarge
	}

	// Validate Author URL scheme to prevent XSS via javascript: or data: URLs
	if paste.AuthorURL != "" {
		// Convert to lowercase for comparison
		urlLower := strings.ToLower(strings.TrimSpace(paste.AuthorURL))

		// Only allow http:// and https:// schemes
		if !strings.HasPrefix(urlLower, "http://") && !strings.HasPrefix(urlLower, "https://") {
			return "", 0, 0, ErrBadRequest
		}

		// Prevent data:, javascript:, vbscript:, file:, etc.
		if strings.Contains(urlLower, "javascript:") ||
		   strings.Contains(urlLower, "data:") ||
		   strings.Contains(urlLower, "vbscript:") ||
		   strings.Contains(urlLower, "file:") {
			return "", 0, 0, ErrBadRequest
		}
	}

	// Create paste
	pasteID, createTime, deleteTime, err := db.PasteAdd(paste)
	if err != nil {
		return pasteID, createTime, deleteTime, err
	}

	return pasteID, createTime, deleteTime, nil
}
