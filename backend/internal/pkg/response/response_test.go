package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/internal/pkg/errors"
)

func TestJSON_Success(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"key": "value"}
	JSON(w, http.StatusOK, data, nil, "req_123")

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Fatalf("expected json content-type, got %s", ct)
	}

	var env Envelope
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		t.Fatal(err)
	}
	if !env.Success {
		t.Fatal("expected success=true")
	}
	if env.RequestID != "req_123" {
		t.Fatalf("expected req_123, got %s", env.RequestID)
	}
	if env.Error != nil {
		t.Fatal("expected no error")
	}
}

func TestJSON_WithMeta(t *testing.T) {
	w := httptest.NewRecorder()
	meta := &Meta{Page: 1, PageSize: 20, Total: 100, TotalPages: 5}
	JSON(w, http.StatusOK, []string{"a"}, meta, "")

	var env Envelope
	json.NewDecoder(w.Body).Decode(&env)
	if env.Meta == nil {
		t.Fatal("expected meta")
	}
	if env.Meta.Page != 1 {
		t.Fatalf("expected page 1, got %d", env.Meta.Page)
	}
}

func TestError_Envelope(t *testing.T) {
	w := httptest.NewRecorder()
	appErr := errors.NewValidation("invalid", []string{"name: required"})
	Error(w, appErr, "req_456")

	resp := w.Result()
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", resp.StatusCode)
	}

	var env Envelope
	json.NewDecoder(resp.Body).Decode(&env)
	if env.Success {
		t.Fatal("expected success=false")
	}
	if env.Error == nil {
		t.Fatal("expected error body")
	}
	if env.Error.Code != "VALIDATION_ERROR" {
		t.Fatalf("expected VALIDATION_ERROR, got %s", env.Error.Code)
	}
	if len(env.Error.Details) != 1 {
		t.Fatalf("expected 1 detail, got %d", len(env.Error.Details))
	}
}

func TestError_NoDetails(t *testing.T) {
	w := httptest.NewRecorder()
	appErr := errors.NewNotFound("resource")
	Error(w, appErr, "")

	var env Envelope
	json.NewDecoder(w.Body).Decode(&env)
	if env.Error == nil {
		t.Fatal("expected error body")
	}
	if env.Error.Details != nil {
		t.Fatal("expected nil details for not-found error")
	}
}
