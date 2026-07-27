package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	apierrors "backend/internal/pkg/errors"
)

// defaultMultipartMemory mirrors net/http's default in-memory threshold for
// multipart parsing; anything larger spools to a temp file (already capped by
// MaxBytesReader below).
const defaultMultipartMemory = 10 << 20

// parseUploadForm enforces the upload size limit and parses the multipart
// form, classifying failures:
//
//	body over maxBytes      → 413 PAYLOAD_TOO_LARGE
//	body not multipart/bad  → 400 BAD_REQUEST
//
// Missing file parts and bad MIME types are the caller's 422s.
func parseUploadForm(w http.ResponseWriter, r *http.Request, maxBytes int64) *apierrors.AppError {
	limitMB := maxBytes >> 20
	tooLarge := apierrors.NewPayloadTooLarge(
		fmt.Sprintf("file quá lớn: giới hạn tối đa %d MB", limitMB))
	tooLarge.Details = []string{fmt.Sprintf("max_bytes: %d", maxBytes)}

	// Well-behaved browsers always send Content-Length for multipart bodies;
	// rejecting here avoids MaxBytesReader's mid-upload connection reset.
	if r.ContentLength > maxBytes {
		return tooLarge
	}

	// Backstop for chunked or lying clients.
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	if err := r.ParseMultipartForm(defaultMultipartMemory); err != nil {
		// The limit does not reliably surface as *http.MaxBytesError — it often
		// arrives as the plain string "http: request body too large".
		var mbe *http.MaxBytesError
		if errors.As(err, &mbe) || strings.Contains(err.Error(), "request body too large") {
			return tooLarge
		}
		return apierrors.NewBadRequest("failed to parse form data")
	}
	return nil
}
