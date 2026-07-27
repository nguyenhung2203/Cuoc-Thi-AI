package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"backend/internal/middleware"
	"backend/internal/pkg/response"
)

// --- Helper: extract envelope ---
func decodeEnvelope(t *testing.T, body *bytes.Buffer) response.Envelope {
	t.Helper()
	var env response.Envelope
	if err := json.NewDecoder(body).Decode(&env); err != nil {
		t.Fatalf("decode envelope: %v", err)
	}
	return env
}

// --- AuthHandler tests ---

func TestAuthHandler_Register_InvalidBody(t *testing.T) {
	t.Parallel()
	h := &AuthHandler{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/auth/register", bytes.NewReader([]byte(`{invalid`)))
	r.Header.Set("Content-Type", "application/json")

	h.Register(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

// K-AUTH-04/05 regression: the validator must run before the service, so a
// nil-service handler returns 422 instead of panicking (which is exactly what
// this test silently recovered from before the fix).
func TestAuthHandler_Register_MissingFields(t *testing.T) {
	t.Parallel()
	h := &AuthHandler{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/auth/register", bytes.NewReader([]byte(`{}`)))
	r.Header.Set("Content-Type", "application/json")

	h.Register(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("empty body: got %d, want 422", w.Code)
	}
	env := decodeEnvelope(t, w.Body)
	if env.Success || env.Error == nil || env.Error.Code != "VALIDATION_ERROR" {
		t.Fatalf("envelope = %+v, want success=false VALIDATION_ERROR", env)
	}
	if len(env.Error.Details) < 3 {
		t.Fatalf("details = %v, want >= 3 (email, password, full_name)", env.Error.Details)
	}
}

func TestAuthHandler_Register_InvalidEmail(t *testing.T) {
	t.Parallel()
	h := &AuthHandler{}
	body := `{"email":"abc@","password":"secret123","full_name":"Nguyen Van A"}`
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/auth/register", bytes.NewReader([]byte(body)))
	r.Header.Set("Content-Type", "application/json")

	h.Register(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("bad email: got %d, want 422", w.Code)
	}
	env := decodeEnvelope(t, w.Body)
	if !containsDetail(env.Error.Details, "Email: failed email validation") {
		t.Fatalf("details = %v, want email validation message", env.Error.Details)
	}
}

func TestAuthHandler_Register_ShortPassword(t *testing.T) {
	t.Parallel()
	h := &AuthHandler{}
	body := `{"email":"a@b.com","password":"123","full_name":"Nguyen Van A"}`
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/auth/register", bytes.NewReader([]byte(body)))
	r.Header.Set("Content-Type", "application/json")

	h.Register(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("short password: got %d, want 422", w.Code)
	}
	env := decodeEnvelope(t, w.Body)
	if !containsDetail(env.Error.Details, "Password: failed min validation") {
		t.Fatalf("details = %v, want password min message", env.Error.Details)
	}
}

func containsDetail(details []string, want string) bool {
	for _, d := range details {
		if d == want {
			return true
		}
	}
	return false
}

func TestAuthHandler_Me_Unauthenticated(t *testing.T) {
	t.Parallel()
	h := &AuthHandler{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/auth/me", nil)

	h.Me(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

// --- JobHandler tests ---

func TestJobHandler_List_EmptyCompany(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r != nil {
			t.Logf("recovered panic (nil svc): %v", r)
		}
	}()
	h := &JobHandler{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/companies//jobs", nil)
	h.List(w, r)
}

// --- InterviewHandler tests (nil-safe) ---

func TestInterviewHandler_JoinByToken_Invalid(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r != nil {
			t.Logf("recovered panic (nil svc): %v", r)
		}
	}()

	h := &InterviewHandler{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/interviews/join/invalid-token", nil)
	h.JoinByToken(w, r)
}

// --- ReportHandler tests (nil-safe) ---
func TestReportHandler_MissingParams(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r != nil {
			t.Logf("recovered panic (nil svc): %v", r)
		}
	}()

	h := &ReportHandler{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/companies//interviews//report", nil)
	h.GetReport(w, r)
}

// --- Chi routing invariant (K-COMP-05) ---

// Locks the chi v5 precedence this codebase depends on: the param-node mount
// /companies/{company_id} wins over the /companies mount's catch-all for any
// path carrying an id. If a chi upgrade ever changes this, per-company routes
// silently fall back to the unscoped mount — this test turns that into a
// build-time failure.
func TestChiRouting_ParamRouteWinsOverMountedCatchAll(t *testing.T) {
	t.Parallel()
	r := chi.NewRouter()
	r.Route("/companies", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("list")) })
	})
	r.Route("/companies/{company_id}", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("scoped")) })
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/companies/abc", nil))
	if got := w.Body.String(); got != "scoped" {
		t.Fatalf("/companies/abc routed to %q, want scoped (chi precedence changed!)", got)
	}

	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/companies", nil))
	if got := w.Body.String(); got != "list" {
		t.Fatalf("/companies routed to %q, want list", got)
	}
}

// --- CompanyHandler tests ---

// ScopedRoutes must register GET / and PUT / — a 405 (not 404) for another
// method proves both endpoints exist without touching the nil service.
func TestCompanyHandler_ScopedRoutes_RegistersGetAndPut(t *testing.T) {
	t.Parallel()
	r := chi.NewRouter()
	r.Route("/companies/{company_id}", (&CompanyHandler{}).ScopedRoutes)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("PATCH", "/companies/abc", nil))
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("PATCH: got %d, want 405 (GET/PUT registered)", w.Code)
	}
}

func TestCompanyHandler_Create_MissingName(t *testing.T) {
	t.Parallel()
	h := &CompanyHandler{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/companies", bytes.NewReader([]byte(`{"name":"   "}`)))
	r.Header.Set("Content-Type", "application/json")
	r = r.WithContext(context.WithValue(r.Context(), middleware.CtxUserID, "user-1"))

	h.Create(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("blank name: got %d, want 422", w.Code)
	}
	env := decodeEnvelope(t, w.Body)
	if env.Error == nil || env.Error.Code != "VALIDATION_ERROR" {
		t.Fatalf("envelope = %+v, want VALIDATION_ERROR", env)
	}
}

func TestCompanyHandler_Create_InvalidBody(t *testing.T) {
	t.Parallel()
	h := &CompanyHandler{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/companies", bytes.NewReader([]byte(`{invalid`)))
	r.Header.Set("Content-Type", "application/json")
	r = r.WithContext(context.WithValue(r.Context(), middleware.CtxUserID, "user-1"))

	h.Create(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("broken json: got %d, want 400", w.Code)
	}
}

func TestCompanyHandler_Update_MissingName(t *testing.T) {
	t.Parallel()
	h := &CompanyHandler{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/companies/abc", bytes.NewReader([]byte(`{"name":""}`)))
	r.Header.Set("Content-Type", "application/json")

	h.Update(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("empty name: got %d, want 422", w.Code)
	}
}

// --- Sample test ---
func TestSampleHandler(t *testing.T) {
	t.Parallel()
	t.Log("handler package test setup ok")
}
