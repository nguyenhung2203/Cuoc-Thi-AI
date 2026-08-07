package handler

import (
	"errors"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"

	apierrors "backend/internal/pkg/errors"
	pkgresponse "backend/internal/pkg/response"
	"backend/internal/storage"
)

// DownloadHandler serves /uploads/{key} guarded by an HMAC signature instead
// of the old anonymous http.FileServer. The signature in the query string is
// the credential, so this handler intentionally lives outside AuthMiddleware —
// browser-native consumers (window.open, <iframe>, <object>) cannot send an
// Authorization header.
type DownloadHandler struct {
	store  *storage.LocalStore
	signer *storage.Signer
}

func NewDownloadHandler(store *storage.LocalStore, signer *storage.Signer) *DownloadHandler {
	return &DownloadHandler{store: store, signer: signer}
}

// inlineTypes maps allowed extensions to a served Content-Type. Anything not
// listed is delivered as an opaque attachment — never text/html or image/svg,
// which would be same-origin XSS on the API host.
var inlineTypes = map[string]string{
	".pdf":  "application/pdf",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
}

var attachmentTypes = map[string]string{
	".doc":  "application/msword",
	".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
}

func (h *DownloadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// When the request path contains percent-escapes, chi routes on
	// r.URL.RawPath and the wildcard arrives STILL ESCAPED (mux.go routeHTTP).
	// The signer signed the decoded key, so decode before verifying — legacy
	// keys with spaces/unicode would otherwise 403. A '%2F' that decodes to
	// '/' is still caught by the store's key validation below.
	key := chi.URLParam(r, "*")
	if dec, err := url.PathUnescape(key); err == nil {
		key = dec
	}

	// Signature check first: a forged request learns nothing about existence.
	if err := h.signer.Verify(key, r.URL.Query()); err != nil {
		pkgresponse.Error(w, apierrors.NewForbidden("invalid or expired download link"), "")
		return
	}

	f, err := h.store.Open(key)
	if err != nil {
		if errors.Is(err, storage.ErrInvalidKey) {
			pkgresponse.Error(w, apierrors.NewForbidden("invalid or expired download link"), "")
			return
		}
		// Signature already proved authorization at issuance time, so a plain
		// 404 for a since-deleted file is not an information leak.
		pkgresponse.Error(w, apierrors.NewNotFound("file not found"), "")
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil || fi.IsDir() {
		pkgresponse.Error(w, apierrors.NewNotFound("file not found"), "")
		return
	}

	ext := strings.ToLower(filepath.Ext(key))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "private, no-store")
	w.Header().Set("Cross-Origin-Resource-Policy", "cross-origin")
	if ct, ok := inlineTypes[ext]; ok {
		w.Header().Set("Content-Type", ct)
		w.Header().Set("Content-Disposition", "inline")
	} else if ct, ok := attachmentTypes[ext]; ok {
		w.Header().Set("Content-Type", ct)
		w.Header().Set("Content-Disposition", `attachment; filename="`+key+`"`)
	} else {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", `attachment; filename="`+key+`"`)
	}

	// Empty name disables ServeContent's extension sniffing; Range/HEAD and
	// If-Modified-Since come for free (PDF viewers issue Range requests).
	http.ServeContent(w, r, "", fi.ModTime(), f)
}
