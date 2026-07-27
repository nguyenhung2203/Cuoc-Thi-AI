package service

import "testing"

func TestSanitizeExt(t *testing.T) {
	cases := map[string]string{
		".pdf":    ".pdf",
		".PDF":    ".pdf",
		".docx":   ".docx",
		".p df":   "", // client filename garbage — the K-FILE key-poisoning class
		".php":    ".php", // matches the pattern; safety comes from mimeExt only emitting whitelisted ext
		"":        "",
		"pdf":     "", // no leading dot
		".toolong1": "",
		`..\x`:    "",
	}
	for in, want := range cases {
		if got := sanitizeExt(in); got != want {
			t.Errorf("sanitizeExt(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestIsAllowedMIME(t *testing.T) {
	allowed := []string{
		"application/pdf",
		"image/jpeg",
		"image/png",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	}
	for _, m := range allowed {
		if !isAllowedMIME(m) {
			t.Errorf("isAllowedMIME(%q) = false, want true", m)
		}
	}
	// text/html and image/svg+xml render in the browser → same-origin XSS if
	// ever stored and served; must stay rejected.
	denied := []string{"text/html", "image/svg+xml", "application/x-msdownload", "text/plain", ""}
	for _, m := range denied {
		if isAllowedMIME(m) {
			t.Errorf("isAllowedMIME(%q) = true, want false", m)
		}
	}
}

// Every MIME type in the upload whitelist must map to a sanitizable extension —
// otherwise uploads of that type silently lose their extension and download
// as octet-stream attachments.
func TestMimeExtCoversWhitelist(t *testing.T) {
	for mime := range allowedMIMEs {
		ext, ok := mimeExt[mime]
		if !ok {
			t.Errorf("mimeExt missing entry for whitelisted %q", mime)
			continue
		}
		if sanitizeExt(ext) != ext {
			t.Errorf("mimeExt[%q] = %q does not survive sanitizeExt", mime, ext)
		}
	}
}
