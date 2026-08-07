package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/time/rate"
)

func TestAuthRateLimitMiddleware_DefaultLimit(t *testing.T) {
	mu.Lock()
	visitors = make(map[string]*rate.Limiter)
	mu.Unlock()

	handler := AuthRateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	for i := 0; i < 11; i++ {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/auth", nil)
		r.RemoteAddr = "default-client:1234"
		handler.ServeHTTP(w, r)
		if i < 10 && w.Code != http.StatusNoContent {
			t.Fatalf("request %d: expected 204, got %d", i+1, w.Code)
		}
		if i == 10 && w.Code != http.StatusTooManyRequests {
			t.Fatalf("request 11: expected 429, got %d", w.Code)
		}
	}
}

func TestLoginRateLimitMiddleware_AllowsThirty(t *testing.T) {
	mu.Lock()
	visitors = make(map[string]*rate.Limiter)
	mu.Unlock()

	handler := LoginRateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	for i := 0; i < 31; i++ {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
		r.RemoteAddr = "login-client:1234"
		handler.ServeHTTP(w, r)
		if i < 30 && w.Code != http.StatusNoContent {
			t.Fatalf("request %d: expected 204, got %d", i+1, w.Code)
		}
		if i == 30 && w.Code != http.StatusTooManyRequests {
			t.Fatalf("request 31: expected 429, got %d", w.Code)
		}
	}
}

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

// requirePermissionStatus runs RequirePermission(action) against a request
// whose context carries the given company role (systemRole optional).
func requirePermissionStatus(t *testing.T, action, companyRole, systemRole string) int {
	t.Helper()
	mw := RequirePermission(action)
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(r.Context(), CtxRequestID, "req-perm")
	if companyRole != "" {
		ctx = context.WithValue(ctx, CtxCompanyRole, companyRole)
	}
	if systemRole != "" {
		ctx = context.WithValue(ctx, CtxUserRole, systemRole)
	}
	h.ServeHTTP(w, r.WithContext(ctx))
	return w.Result().StatusCode
}

// K-RBAC-03 regression locks: the job/candidate actions wired onto the routes
// must keep their role mapping.
func TestRequirePermission_RoleMapping(t *testing.T) {
	cases := []struct {
		action, companyRole string
		want                int
	}{
		{"job:create", "member", http.StatusOK},
		{"job:create", "viewer", http.StatusForbidden},
		{"job:delete", "viewer", http.StatusForbidden},
		{"job:read", "viewer", http.StatusOK},
		{"candidate:read", "viewer", http.StatusOK},
		{"candidate:delete", "viewer", http.StatusForbidden},
		{"candidate:update", "member", http.StatusOK},
		{"interview:update", "owner", http.StatusOK},
	}
	for _, c := range cases {
		if got := requirePermissionStatus(t, c.action, c.companyRole, "recruiter"); got != c.want {
			t.Errorf("%s as %s: got %d, want %d", c.action, c.companyRole, got, c.want)
		}
	}
}

// Documents the trap in rbac.go: an action with no mapping silently falls
// back to owner/admin only — a typo'd action locks members out, not open.
func TestRequirePermission_UnknownAction_DefaultsToOwnerAdmin(t *testing.T) {
	if got := requirePermissionStatus(t, "totally:unknown", "member", "recruiter"); got != http.StatusForbidden {
		t.Fatalf("member on unknown action: got %d, want 403", got)
	}
	if got := requirePermissionStatus(t, "totally:unknown", "owner", "recruiter"); got != http.StatusOK {
		t.Fatalf("owner on unknown action: got %d, want 200", got)
	}
}

// Documents the second trap: guards attached outside CompanyScopeMiddleware
// see an empty company role and reject everyone but system admins.
func TestRequireCompanyRole_MissingCompanyRole_Blocks(t *testing.T) {
	mw := RequireCompanyRole("owner", "admin")
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(r.Context(), CtxUserRole, "recruiter") // no CtxCompanyRole
	h.ServeHTTP(w, r.WithContext(ctx))

	if w.Result().StatusCode != http.StatusForbidden {
		t.Fatalf("missing company role: got %d, want 403", w.Result().StatusCode)
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
