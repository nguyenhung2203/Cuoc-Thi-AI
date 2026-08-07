package storage

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

func testSigner(secret string, ttl time.Duration, now time.Time) *Signer {
	s := NewSigner([]byte(secret), ttl)
	s.now = func() time.Time { return now }
	return s
}

func queryValues(t *testing.T, query string) url.Values {
	t.Helper()
	q, err := url.ParseQuery(query)
	if err != nil {
		t.Fatalf("parse query %q: %v", query, err)
	}
	return q
}

func TestSigner_SignVerify_RoundTrip(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	s := testSigner("secret", 15*time.Minute, now)

	query, expiresAt := s.Sign("cv-abc.pdf")
	if got, want := expiresAt.Unix(), now.Add(15*time.Minute).Unix(); got != want {
		t.Fatalf("expiresAt = %d, want %d", got, want)
	}
	if err := s.Verify("cv-abc.pdf", queryValues(t, query)); err != nil {
		t.Fatalf("Verify: %v", err)
	}
}

// The core property: a signature minted for one key must never authorize
// another key.
func TestSigner_Verify_RejectsTamperedKey(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	s := testSigner("secret", 15*time.Minute, now)

	query, _ := s.Sign("file-a.pdf")
	if err := s.Verify("file-b.pdf", queryValues(t, query)); !errors.Is(err, ErrSignatureInvalid) {
		t.Fatalf("want ErrSignatureInvalid for swapped key, got %v", err)
	}
}

// exp is inside the MAC — bumping it must invalidate the signature.
func TestSigner_Verify_RejectsTamperedExp(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	s := testSigner("secret", 15*time.Minute, now)

	query, expiresAt := s.Sign("cv.pdf")
	q := queryValues(t, query)
	q.Set("exp", strconv.FormatInt(expiresAt.Unix()+3600, 10))
	if err := s.Verify("cv.pdf", q); !errors.Is(err, ErrSignatureInvalid) {
		t.Fatalf("want ErrSignatureInvalid for bumped exp, got %v", err)
	}
}

func TestSigner_Verify_RejectsExpired(t *testing.T) {
	signedAt := time.Unix(1_700_000_000, 0)
	s := testSigner("secret", 15*time.Minute, signedAt)
	query, expiresAt := s.Sign("cv.pdf")

	// Move the clock past exp + leeway.
	s.now = func() time.Time { return expiresAt.Add(clockLeeway + time.Second) }
	if err := s.Verify("cv.pdf", queryValues(t, query)); !errors.Is(err, ErrSignatureExpired) {
		t.Fatalf("want ErrSignatureExpired, got %v", err)
	}
}

func TestSigner_Verify_AllowsWithinLeeway(t *testing.T) {
	signedAt := time.Unix(1_700_000_000, 0)
	s := testSigner("secret", 15*time.Minute, signedAt)
	query, expiresAt := s.Sign("cv.pdf")

	s.now = func() time.Time { return expiresAt.Add(clockLeeway / 2) }
	if err := s.Verify("cv.pdf", queryValues(t, query)); err != nil {
		t.Fatalf("within leeway should verify, got %v", err)
	}
}

func TestSigner_Verify_MalformedParams(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	s := testSigner("secret", 15*time.Minute, now)
	_, expiresAt := s.Sign("cv.pdf")

	cases := map[string]url.Values{
		"empty":       {},
		"missing sig": {"exp": {strconv.FormatInt(expiresAt.Unix(), 10)}},
		"missing exp": {"sig": {"deadbeef"}},
		"exp text":    {"exp": {"abc"}, "sig": {"deadbeef"}},
		"exp empty":   {"exp": {""}, "sig": {"deadbeef"}},
		"exp huge":    {"exp": {"99999999999999999999999999"}, "sig": {"deadbeef"}},
	}
	for name, q := range cases {
		if err := s.Verify("cv.pdf", q); !errors.Is(err, ErrSignatureInvalid) {
			t.Errorf("%s: want ErrSignatureInvalid, got %v", name, err)
		}
	}
}

// Rotating the secret revokes every outstanding URL.
func TestSigner_Verify_DifferentSecret(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	a := testSigner("secret-a", 15*time.Minute, now)
	b := testSigner("secret-b", 15*time.Minute, now)

	query, _ := a.Sign("cv.pdf")
	if err := b.Verify("cv.pdf", queryValues(t, query)); !errors.Is(err, ErrSignatureInvalid) {
		t.Fatalf("want ErrSignatureInvalid across secrets, got %v", err)
	}
}

// Legacy keys may contain spaces/odd chars (ext came from client filenames);
// SignURL must escape them into a parseable URL that round-trips.
func TestSigner_SignURL_EscapesKey(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	s := testSigner("secret", 15*time.Minute, now)

	raw, _ := s.SignURL("http://localhost:18080", "id.p df")
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("signed URL unparseable: %v", err)
	}
	key := strings.TrimPrefix(u.Path, "/uploads/")
	unescaped, err := url.PathUnescape(key)
	if err != nil || unescaped != "id.p df" {
		t.Fatalf("key round-trip: got %q (err %v), want %q", unescaped, err, "id.p df")
	}
	if err := s.Verify(unescaped, u.Query()); err != nil {
		t.Fatalf("Verify after round-trip: %v", err)
	}
}

func TestSigner_SignURL_TrailingSlashBase(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	s := testSigner("secret", 15*time.Minute, now)

	withSlash, _ := s.SignURL("http://x/", "k.pdf")
	withoutSlash, _ := s.SignURL("http://x", "k.pdf")
	if withSlash != withoutSlash {
		t.Fatalf("base slash handling differs:\n%s\n%s", withSlash, withoutSlash)
	}
}
