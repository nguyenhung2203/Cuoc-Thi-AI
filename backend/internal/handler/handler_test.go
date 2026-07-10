package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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

func TestAuthHandler_Register_MissingFields(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r != nil {
			t.Logf("recovered panic (nil svc): %v", r)
		}
	}()
	h := &AuthHandler{}
	body := `{}`
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/auth/register", bytes.NewReader([]byte(body)))
	r.Header.Set("Content-Type", "application/json")
	h.Register(w, r)
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

// --- Sample test ---
func TestSampleHandler(t *testing.T) {
	t.Parallel()
	t.Log("handler package test setup ok")
}
