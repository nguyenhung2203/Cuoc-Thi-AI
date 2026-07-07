package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireCandidate_AllowsCandidate(t *testing.T) {
	handler := RequireCandidate()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(r.Context(), CtxUserRole, "candidate")
	r = r.WithContext(ctx)

	handler.ServeHTTP(w, r)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for candidate role, got %d", w.Result().StatusCode)
	}
}

func TestRequireCandidate_BlocksRecruiter(t *testing.T) {
	handler := RequireCandidate()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(r.Context(), CtxUserRole, "recruiter")
	r = r.WithContext(ctx)

	handler.ServeHTTP(w, r)

	if w.Result().StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 for recruiter, got %d", w.Result().StatusCode)
	}
}

func TestRequireCandidate_BlocksNoRole(t *testing.T) {
	handler := RequireCandidate()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)

	handler.ServeHTTP(w, r)

	// Without CtxUserRole set, the value will be empty string, not "candidate" → 403
	if w.Result().StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 for unauthenticated, got %d", w.Result().StatusCode)
	}
}

func TestRequireRole_AllowsMatchingRole(t *testing.T) {
	mw := RequireRole("admin", "recruiter")
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(r.Context(), CtxUserRole, "recruiter")
	ctx = context.WithValue(ctx, CtxRequestID, "req-1")
	r = r.WithContext(ctx)

	handler.ServeHTTP(w, r)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for recruiter, got %d", w.Result().StatusCode)
	}
}

func TestRequireRole_BlocksNonMatchingRole(t *testing.T) {
	mw := RequireRole("admin", "recruiter")
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(r.Context(), CtxUserRole, "candidate")
	ctx = context.WithValue(ctx, CtxRequestID, "req-2")
	r = r.WithContext(ctx)

	handler.ServeHTTP(w, r)

	if w.Result().StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 for candidate, got %d", w.Result().StatusCode)
	}
}

func TestRequireCompanyRole_AdminBypass(t *testing.T) {
	mw := RequireCompanyRole("owner", "admin")
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(r.Context(), CtxUserRole, "admin") // system admin bypasses
	ctx = context.WithValue(ctx, CtxCompanyRole, "viewer")
	ctx = context.WithValue(ctx, CtxRequestID, "req-3")
	r = r.WithContext(ctx)

	handler.ServeHTTP(w, r)

	// System admin should always be allowed
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for system admin, got %d", w.Result().StatusCode)
	}
}

func TestContextKeys(t *testing.T) {
	// Verify context keys are consistently typed
	keys := []contextKey{CtxUserID, CtxUserRole, CtxCompanyID, CtxRequestID, CtxCompanyRole}
	for _, k := range keys {
		val := string(k)
		if val == "" {
			t.Error("context key has empty string value")
		}
	}
}

func TestGetRequestIDFromCtx(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	if id := getRequestIDFromCtx(r); id != "" {
		t.Fatalf("expected empty request id from empty context, got %s", id)
	}

	ctx := context.WithValue(r.Context(), CtxRequestID, "req-42")
	r = r.WithContext(ctx)
	if id := getRequestIDFromCtx(r); id != "req-42" {
		t.Fatalf("expected req-42, got %s", id)
	}
}
