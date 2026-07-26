package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"backend/internal/storage"
)

// newDownloadFixture builds a jailed store in a temp dir with one PDF file,
// a signer, and a chi router mounted like cmd/api/main.go does.
func newDownloadFixture(t *testing.T, ttl time.Duration) (*storage.Signer, chi.Router, string) {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "cv-test.pdf"), []byte("%PDF-1.4 test"), 0o644); err != nil {
		t.Fatal(err)
	}
	store, err := storage.NewLocalStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	// Windows: the os.Root dir handle must be released before TempDir cleanup.
	t.Cleanup(func() { store.Close() })
	signer := storage.NewSigner([]byte("test-secret"), ttl)
	dl := NewDownloadHandler(store, signer)

	r := chi.NewRouter()
	r.Method(http.MethodGet, "/uploads/*", dl)
	r.Method(http.MethodHead, "/uploads/*", dl)
	return signer, r, dir
}

func signedPath(signer *storage.Signer, key string) string {
	query, _ := signer.Sign(key)
	return "/uploads/" + key + "?" + query
}

func TestDownload_ValidSignature_ServesFile(t *testing.T) {
	signer, r, _ := newDownloadFixture(t, 15*time.Minute)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, signedPath(signer, "cv-test.pdf"), nil))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body: %s)", w.Code, w.Body.String())
	}
	if got := w.Header().Get("Content-Type"); got != "application/pdf" {
		t.Errorf("Content-Type = %q, want application/pdf", got)
	}
	if got := w.Header().Get("X-Content-Type-Options"); got != "nosniff" {
		t.Errorf("X-Content-Type-Options = %q, want nosniff", got)
	}
	if got := w.Header().Get("Content-Disposition"); got != "inline" {
		t.Errorf("Content-Disposition = %q, want inline (iframe/object viewers)", got)
	}
	if !strings.Contains(w.Body.String(), "%PDF-1.4") {
		t.Error("body does not contain the file content")
	}
}

// Regression lock for K-FILE-01: the old FileServer served everything with no
// credential — a bare /uploads/{key} must now be rejected.
func TestDownload_NoSignature_403(t *testing.T) {
	_, r, _ := newDownloadFixture(t, 15*time.Minute)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/uploads/cv-test.pdf", nil))
	if w.Code != http.StatusForbidden {
		t.Fatalf("unsigned request: status = %d, want 403", w.Code)
	}
	if strings.Contains(w.Body.String(), "%PDF") {
		t.Fatal("unsigned request leaked file content")
	}
}

func TestDownload_InvalidSignature_403(t *testing.T) {
	_, r, _ := newDownloadFixture(t, 15*time.Minute)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/uploads/cv-test.pdf?exp=9999999999&sig=deadbeef", nil))
	if w.Code != http.StatusForbidden {
		t.Fatalf("forged sig: status = %d, want 403", w.Code)
	}
}

// Crafts a correctly-signed URL whose exp is far in the past. Recomputing the
// MAC here also locks the wire format ("v1\n"+key+"\n"+exp, HMAC-SHA256,
// base64url) — a silent format change breaks this test.
func TestDownload_Expired_403(t *testing.T) {
	_, r, _ := newDownloadFixture(t, 15*time.Minute)

	exp := time.Now().Add(-time.Hour).Unix() // long past exp+leeway
	mac := hmac.New(sha256.New, []byte("test-secret"))
	mac.Write([]byte("v1\ncv-test.pdf\n" + strconv.FormatInt(exp, 10)))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	path := "/uploads/cv-test.pdf?exp=" + strconv.FormatInt(exp, 10) + "&sig=" + sig
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
	if w.Code != http.StatusForbidden {
		t.Fatalf("expired sig: status = %d, want 403", w.Code)
	}
	if strings.Contains(w.Body.String(), "%PDF") {
		t.Fatal("expired request leaked file content")
	}
}

// The old http.FileServer rendered a directory listing at /uploads/.
func TestDownload_DirectoryListing_NotServed(t *testing.T) {
	signer, r, _ := newDownloadFixture(t, 15*time.Minute)

	for _, path := range []string{"/uploads/", signedPath(signer, "")} {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
		if w.Code == http.StatusOK {
			t.Fatalf("%s: directory listing served (200)", path)
		}
		if strings.Contains(w.Body.String(), "cv-test.pdf") {
			t.Fatalf("%s: response leaked directory contents", path)
		}
	}
}

// Traversal defense must hold even when the attacker holds a VALID signature
// for the malicious key (i.e. it is independent of the signature layer).
func TestDownload_PathTraversal(t *testing.T) {
	signer, r, dir := newDownloadFixture(t, 15*time.Minute)

	// Plant a file outside the jail that traversal would reach.
	outside := filepath.Join(filepath.Dir(dir), "secret.txt")
	if err := os.WriteFile(outside, []byte("outside-secret"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(outside) })

	keys := []string{
		"../secret.txt",
		"..%2Fsecret.txt",
		"%2e%2e/secret.txt",
		"sub/../../secret.txt",
		`..\secret.txt`,
		"..",
		".hidden",
	}
	for _, key := range keys {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, signedPath(signer, key), nil))
		if w.Code == http.StatusOK {
			t.Errorf("key %q: served with 200", key)
		}
		if strings.Contains(w.Body.String(), "outside-secret") {
			t.Errorf("key %q: escaped the uploads jail", key)
		}
	}
}

// Valid signature but the file has since been deleted → 404, distinct from 403.
func TestDownload_MissingFile_404(t *testing.T) {
	signer, r, _ := newDownloadFixture(t, 15*time.Minute)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, signedPath(signer, "gone.pdf"), nil))
	if w.Code != http.StatusNotFound {
		t.Fatalf("missing file: status = %d, want 404", w.Code)
	}
}

func TestDownload_IsDirectory_404(t *testing.T) {
	signer, r, dir := newDownloadFixture(t, 15*time.Minute)
	if err := os.Mkdir(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, signedPath(signer, "sub"), nil))
	if w.Code != http.StatusNotFound {
		t.Fatalf("directory key: status = %d, want 404", w.Code)
	}
}

func TestDownload_HEAD(t *testing.T) {
	signer, r, _ := newDownloadFixture(t, 15*time.Minute)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodHead, signedPath(signer, "cv-test.pdf"), nil))
	if w.Code != http.StatusOK {
		t.Fatalf("HEAD: status = %d, want 200", w.Code)
	}
	if w.Header().Get("Content-Length") == "" {
		t.Error("HEAD: missing Content-Length")
	}
	if w.Body.Len() != 0 {
		t.Errorf("HEAD: body should be empty, got %d bytes", w.Body.Len())
	}
}

// Unknown extensions must never render in the browser (same-origin XSS).
func TestDownload_UnknownType_Attachment(t *testing.T) {
	signer, r, dir := newDownloadFixture(t, 15*time.Minute)
	if err := os.WriteFile(filepath.Join(dir, "evil.html"), []byte("<script>alert(1)</script>"), 0o644); err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, signedPath(signer, "evil.html"), nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	if got := w.Header().Get("Content-Type"); got != "application/octet-stream" {
		t.Errorf("Content-Type = %q, want application/octet-stream", got)
	}
	if got := w.Header().Get("Content-Disposition"); !strings.HasPrefix(got, "attachment") {
		t.Errorf("Content-Disposition = %q, want attachment", got)
	}
}

func TestDownload_MethodNotAllowed(t *testing.T) {
	signer, r, _ := newDownloadFixture(t, 15*time.Minute)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, signedPath(signer, "cv-test.pdf"), nil))
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST: status = %d, want 405", w.Code)
	}
}
