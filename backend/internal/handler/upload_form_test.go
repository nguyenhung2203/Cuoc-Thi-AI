package handler

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	apierrors "backend/internal/pkg/errors"
)

// multipartBody builds a multipart body with a single "file" part of n bytes.
func multipartBody(t *testing.T, fieldName string, n int) (*bytes.Buffer, string) {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile(fieldName, "cv.pdf")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fw.Write(bytes.Repeat([]byte("a"), n)); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return &buf, w.FormDataContentType()
}

func TestParseUploadForm_UnderLimit(t *testing.T) {
	body, ct := multipartBody(t, "file", 1024)
	r := httptest.NewRequest(http.MethodPost, "/upload", body)
	r.Header.Set("Content-Type", ct)

	if appErr := parseUploadForm(httptest.NewRecorder(), r, 1<<20); appErr != nil {
		t.Fatalf("under limit: got %v", appErr)
	}
	if _, _, err := r.FormFile("file"); err != nil {
		t.Fatalf("FormFile after parse: %v", err)
	}
}

func TestParseUploadForm_ExactlyAtLimit(t *testing.T) {
	body, ct := multipartBody(t, "file", 1024)
	limit := int64(body.Len()) // the whole body is exactly the limit
	r := httptest.NewRequest(http.MethodPost, "/upload", body)
	r.Header.Set("Content-Type", ct)

	if appErr := parseUploadForm(httptest.NewRecorder(), r, limit); appErr != nil {
		t.Fatalf("exactly at limit should pass, got %v", appErr)
	}
}

func TestParseUploadForm_OverLimit_413(t *testing.T) {
	body, ct := multipartBody(t, "file", 4096)
	r := httptest.NewRequest(http.MethodPost, "/upload", body)
	r.Header.Set("Content-Type", ct)

	appErr := parseUploadForm(httptest.NewRecorder(), r, 1024)
	if appErr == nil {
		t.Fatal("over limit: want error, got nil")
	}
	if appErr.HTTPStatus != http.StatusRequestEntityTooLarge || appErr.Code != apierrors.PAYLOAD_TOO_LARGE {
		t.Fatalf("got %d/%s, want 413/PAYLOAD_TOO_LARGE", appErr.HTTPStatus, appErr.Code)
	}
}

// Chunked/lying clients send no Content-Length — the MaxBytesReader backstop
// must still produce a 413, not a 400.
func TestParseUploadForm_OverLimit_NoContentLength_413(t *testing.T) {
	body, ct := multipartBody(t, "file", 4096)
	r := httptest.NewRequest(http.MethodPost, "/upload", body)
	r.Header.Set("Content-Type", ct)
	r.ContentLength = -1

	appErr := parseUploadForm(httptest.NewRecorder(), r, 1024)
	if appErr == nil || appErr.HTTPStatus != http.StatusRequestEntityTooLarge {
		t.Fatalf("chunked over limit: got %v, want 413", appErr)
	}
}

// The core classification assert: a broken multipart body under the limit is
// the client's fault (400), never a 413.
func TestParseUploadForm_MalformedBody_400(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/upload", bytes.NewBufferString("this is not multipart"))
	r.Header.Set("Content-Type", "multipart/form-data; boundary=xyz")

	appErr := parseUploadForm(httptest.NewRecorder(), r, 1<<20)
	if appErr == nil {
		t.Fatal("malformed body: want error, got nil")
	}
	if appErr.HTTPStatus != http.StatusBadRequest {
		t.Fatalf("malformed body: got %d/%s, want 400", appErr.HTTPStatus, appErr.Code)
	}
}

func TestParseUploadForm_NotMultipart_400(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/upload", bytes.NewBufferString(`{"a":1}`))
	r.Header.Set("Content-Type", "application/json")

	appErr := parseUploadForm(httptest.NewRecorder(), r, 1<<20)
	if appErr == nil || appErr.HTTPStatus != http.StatusBadRequest {
		t.Fatalf("json body: got %v, want 400", appErr)
	}
}

// A valid form without the expected part parses fine — the missing-part 422
// belongs to the caller (FormFile), not the size/parse helper.
func TestParseUploadForm_MissingPart_IsCallerConcern(t *testing.T) {
	body, ct := multipartBody(t, "other_field", 128)
	r := httptest.NewRequest(http.MethodPost, "/upload", body)
	r.Header.Set("Content-Type", ct)

	if appErr := parseUploadForm(httptest.NewRecorder(), r, 1<<20); appErr != nil {
		t.Fatalf("valid form: got %v", appErr)
	}
	if _, _, err := r.FormFile("file"); err == nil {
		t.Fatal("expected FormFile to fail for the missing part")
	}
}
