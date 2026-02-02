
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package web

import (
	"encoding/base64"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/casjay-forks/caspaste/src/lineend"
	"github.com/casjay-forks/caspaste/src/netshare"
)

// File type detection helpers
func isImageMimeType(mimeType string) bool {
	return strings.HasPrefix(mimeType, "image/")
}

func isVideoMimeType(mimeType string) bool {
	return strings.HasPrefix(mimeType, "video/")
}

func isAudioMimeType(mimeType string) bool {
	return strings.HasPrefix(mimeType, "audio/")
}

func isPDFMimeType(mimeType string) bool {
	return mimeType == "application/pdf"
}

func isTextMimeType(mimeType string) bool {
	if strings.HasPrefix(mimeType, "text/") {
		return true
	}
	// Common text-based MIME types
	textTypes := []string{
		"application/json",
		"application/xml",
		"application/javascript",
		"application/x-javascript",
		"application/ecmascript",
		"application/x-sh",
		"application/x-csh",
	}
	for _, t := range textTypes {
		if mimeType == t {
			return true
		}
	}
	return false
}

type pasteTmpl struct {
	ID         string
	Title      string
	Body       template.HTML
	Syntax     string
	CreateTime int64
	DeleteTime int64
	OneUse     bool

	LineEnd       string
	CreateTimeStr string
	DeleteTimeStr string

	Author      string
	AuthorEmail string
	AuthorURL   string

	// File upload fields
	IsFile   bool
	FileName string
	MimeType string
	FileSize int

	// File type flags for template rendering
	IsImage bool
	IsVideo bool
	IsAudio bool
	IsPDF   bool
	IsText  bool

	// Data URL for embedding media (images, video, audio)
	// Using template.URL to mark as safe for embedding
	MediaDataURL template.URL

	Language  string
	Theme     func(string) string
	Translate func(string, ...interface{}) template.HTML
}

type pasteContinueTmpl struct {
	ID        string
	Language  string
	Theme     func(string) string
	Translate func(string, ...interface{}) template.HTML
	// CSRF token for form protection per AI.md PART 11
	CSRFToken string
}

func (data *Data) handleGetPaste(rw http.ResponseWriter, req *http.Request) error {
	// Check rate limit
	err := data.RateLimitGet.CheckAndUse(netshare.GetClientAddr(req))
	if err != nil {
		return err
	}

	// Get paste ID
	pasteID := string([]rune(req.URL.Path)[1:])

	// Read DB
	paste, err := data.DB.PasteGet(pasteID)
	if err != nil {
		return err
	}

	// If "one use" paste
	if paste.OneUse {
		// If continue button not pressed
		req.ParseForm()

		if req.PostForm.Get("oneUseContinue") != "true" {
			tmplData := pasteContinueTmpl{
				ID:        paste.ID,
				Language:  getCookie(req, "lang"),
				Theme:     data.getThemeFunc(req),
				Translate: data.Locales.findLocale(req).translate,
				CSRFToken: GetCSRFToken(req, 32),
			}

			return data.PasteContinue.Execute(rw, tmplData)
		}

		// If continue button pressed delete paste
		err = data.DB.PasteDelete(pasteID)
		if err != nil {
			return err
		}
	}

	// Prepare template data
	createTime := time.Unix(paste.CreateTime, 0).UTC()
	deleteTime := time.Unix(paste.DeleteTime, 0).UTC()

	// Determine body content based on whether this is a file upload
	var bodyContent string
	var fileSize int
	var bodyHTML template.HTML
	var mediaDataURL template.URL
	var isImage, isVideo, isAudio, isPDF, isText bool

	if paste.IsFile {
		// File upload: try to decode base64, fall back to raw for legacy data
		var base64Data string
		fileData, err := base64.StdEncoding.DecodeString(paste.Body)
		if err != nil {
			// Legacy data stored without base64 encoding - use as-is
			bodyContent = paste.Body
			fileSize = len(paste.Body)
			base64Data = base64.StdEncoding.EncodeToString([]byte(paste.Body))
		} else {
			bodyContent = string(fileData)
			fileSize = len(fileData)
			base64Data = paste.Body
		}

		// Detect file type from MIME type
		mimeType := paste.MimeType
		isImage = isImageMimeType(mimeType)
		isVideo = isVideoMimeType(mimeType)
		isAudio = isAudioMimeType(mimeType)
		isPDF = isPDFMimeType(mimeType)
		isText = isTextMimeType(mimeType)

		// For media files, create data URL for embedding
		if isImage || isVideo || isAudio || isPDF {
			mediaDataURL = template.URL("data:" + mimeType + ";base64," + base64Data)
			// Don't syntax highlight media - body will be empty
			bodyHTML = ""
		} else if isText {
			// Text files can be syntax highlighted
			bodyHTML = data.Themes.findTheme(req, data.UiDefaultTheme).tryHighlight(bodyContent, paste.Syntax)
		} else {
			// Binary files - show file info, don't try to display content
			bodyHTML = ""
		}
	} else {
		bodyContent = paste.Body
		bodyHTML = data.Themes.findTheme(req, data.UiDefaultTheme).tryHighlight(bodyContent, paste.Syntax)
	}

	tmplData := pasteTmpl{
		ID:         paste.ID,
		Title:      paste.Title,
		Body:       bodyHTML,
		Syntax:     paste.Syntax,
		CreateTime: paste.CreateTime,
		DeleteTime: paste.DeleteTime,
		OneUse:     paste.OneUse,

		CreateTimeStr: createTime.Format("Mon, 02 Jan 2006 15:04:05 -0700"),
		DeleteTimeStr: deleteTime.Format("Mon, 02 Jan 2006 15:04:05 -0700"),

		Author:      paste.Author,
		AuthorEmail: paste.AuthorEmail,
		AuthorURL:   paste.AuthorURL,

		IsFile:       paste.IsFile,
		FileName:     paste.FileName,
		MimeType:     paste.MimeType,
		FileSize:     fileSize,
		IsImage:      isImage,
		IsVideo:      isVideo,
		IsAudio:      isAudio,
		IsPDF:        isPDF,
		IsText:       isText,
		MediaDataURL: mediaDataURL,

		Language:  getCookie(req, "lang"),
		Theme:     data.getThemeFunc(req),
		Translate: data.Locales.findLocale(req).translate,
	}

	// Get body line end (only for text content)
	if !paste.IsFile || isText {
		switch lineend.GetLineEnd(bodyContent) {
		case "\r\n":
			tmplData.LineEnd = "CRLF"
		case "\r":
			tmplData.LineEnd = "CR"
		default:
			tmplData.LineEnd = "LF"
		}
	}

	// Show paste
	return data.PastePage.Execute(rw, tmplData)
}
