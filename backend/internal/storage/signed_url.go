package storage

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Signed download URLs replace the old anonymous /uploads/* FileServer:
// the HMAC signature carried in the query string IS the credential, so the
// download handler can stay outside AuthMiddleware and browser-native
// consumers (window.open, <iframe>, <object>, <a href>) keep working without
// an Authorization header.

var (
	// ErrSignatureInvalid is returned for missing, malformed, or forged
	// signatures — including a valid signature lifted onto another key.
	ErrSignatureInvalid = errors.New("signed url: invalid signature")
	// ErrSignatureExpired is returned when the signature is authentic but the
	// exp timestamp (plus leeway) has passed.
	ErrSignatureExpired = errors.New("signed url: expired")
)

// clockLeeway absorbs clock drift between the signing and verifying process
// (single process today, but survives multi-replica NTP skew).
const clockLeeway = 60 * time.Second

// signScheme versions the canonical string so the format can evolve without
// old signatures remaining valid under a new scheme.
const signScheme = "v1"

// Signer mints and verifies HMAC-SHA256 signed download URLs.
type Signer struct {
	key    []byte
	ttl    time.Duration
	leeway time.Duration
	now    func() time.Time // injectable for tests
}

// NewSigner builds a Signer. secret must be non-empty; ttl <= 0 falls back to
// 15 minutes.
func NewSigner(secret []byte, ttl time.Duration) *Signer {
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}
	return &Signer{key: secret, ttl: ttl, leeway: clockLeeway, now: time.Now}
}

// canonical builds the exact byte string covered by the MAC. storageKey is the
// DECODED key; '\n' cannot appear in a storage key so the fields cannot bleed
// into each other.
func canonical(storageKey string, exp int64) string {
	return signScheme + "\n" + storageKey + "\n" + strconv.FormatInt(exp, 10)
}

func (s *Signer) mac(storageKey string, exp int64) string {
	h := hmac.New(sha256.New, s.key)
	h.Write([]byte(canonical(storageKey, exp)))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

// Sign returns the query string ("exp=...&sig=...") authorizing storageKey
// until now+ttl, plus the expiry time for the API response.
func (s *Signer) Sign(storageKey string) (query string, expiresAt time.Time) {
	expiresAt = s.now().Add(s.ttl)
	exp := expiresAt.Unix()
	return "exp=" + strconv.FormatInt(exp, 10) + "&sig=" + s.mac(storageKey, exp), expiresAt
}

// SignURL builds the full browser-reachable URL for storageKey.
func (s *Signer) SignURL(baseURL, storageKey string) (string, time.Time) {
	query, expiresAt := s.Sign(storageKey)
	return strings.TrimRight(baseURL, "/") + "/uploads/" + url.PathEscape(storageKey) + "?" + query, expiresAt
}

// Verify checks the exp/sig query parameters against storageKey (the DECODED
// key, e.g. from chi.URLParam). The signature is checked before expiry so a
// forged exp can never influence the outcome.
func (s *Signer) Verify(storageKey string, q url.Values) error {
	expStr := q.Get("exp")
	sig := q.Get("sig")
	if expStr == "" || sig == "" {
		return ErrSignatureInvalid
	}
	exp, err := strconv.ParseInt(expStr, 10, 64)
	if err != nil {
		return ErrSignatureInvalid
	}
	if !hmac.Equal([]byte(s.mac(storageKey, exp)), []byte(sig)) {
		return ErrSignatureInvalid
	}
	if s.now().After(time.Unix(exp, 0).Add(s.leeway)) {
		return ErrSignatureExpired
	}
	return nil
}
