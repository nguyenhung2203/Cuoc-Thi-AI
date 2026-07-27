package errors

import (
	"net/http"
	"testing"
)

func TestAppError_Error(t *testing.T) {
	err := NewInternal("test error")
	if err.Error() != "test error" {
		t.Fatalf("expected 'test error', got '%s'", err.Error())
	}
	if err.HTTPStatus != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", err.HTTPStatus)
	}
}

func TestIsAppError(t *testing.T) {
	appErr := NewNotFound("not found")
	got, ok := IsAppError(appErr)
	if !ok {
		t.Fatal("expected true")
	}
	if got.Code != NOT_FOUND {
		t.Fatalf("expected NOT_FOUND, got %s", got.Code)
	}

	// Non-AppError
	_, ok = IsAppError(http.ErrAbortHandler)
	if ok {
		t.Fatal("expected false for non-AppError")
	}
}

func TestErrorConstructors(t *testing.T) {
	tests := []struct {
		name     string
		constructor func(string) *AppError
		wantCode string
		wantHTTP int
	}{
		{"NewUnauthorized", NewUnauthorized, UNAUTHORIZED, http.StatusUnauthorized},
		{"NewForbidden", NewForbidden, FORBIDDEN, http.StatusForbidden},
		{"NewNotFound", NewNotFound, NOT_FOUND, http.StatusNotFound},
		{"NewConflict", NewConflict, CONFLICT, http.StatusConflict},
		{"NewBadRequest", NewBadRequest, BAD_REQUEST, http.StatusBadRequest},
		{"NewInternal", NewInternal, INTERNAL_ERROR, http.StatusInternalServerError},
		{"NewRateLimited", NewRateLimited, RATE_LIMITED, http.StatusTooManyRequests},
		{"NewAccountBlocked", NewAccountBlocked, ACCOUNT_BLOCKED, http.StatusForbidden},
		{"NewPayloadTooLarge", NewPayloadTooLarge, PAYLOAD_TOO_LARGE, http.StatusRequestEntityTooLarge},
		{"NewAIServiceError", NewAIServiceError, AI_SERVICE_ERROR, http.StatusBadGateway},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.constructor(tt.name)
			if err.Code != tt.wantCode {
				t.Errorf("code: got %s, want %s", err.Code, tt.wantCode)
			}
			if err.HTTPStatus != tt.wantHTTP {
				t.Errorf("http: got %d, want %d", err.HTTPStatus, tt.wantHTTP)
			}
			if err.Message != tt.name {
				t.Errorf("msg: got %s, want %s", err.Message, tt.name)
			}
		})
	}
}

// TestAccountBlockedCode_FrontendContract locks the literal error code the
// frontend switches on (api.service.js redirects ACCOUNT_BLOCKED to
// /403?reason=account_blocked). Renaming the constant must fail this test.
func TestAccountBlockedCode_FrontendContract(t *testing.T) {
	if ACCOUNT_BLOCKED != "ACCOUNT_BLOCKED" {
		t.Fatalf("ACCOUNT_BLOCKED constant changed to %q — update frontend/src/services/api.service.js in lockstep", ACCOUNT_BLOCKED)
	}
}

func TestNewValidation(t *testing.T) {
	details := []string{"field1: required", "field2: email"}
	err := NewValidation("validation error", details)
	if err.Code != VALIDATION_ERROR {
		t.Fatalf("expected VALIDATION_ERROR, got %s", err.Code)
	}
	if len(err.Details) != 2 {
		t.Fatalf("expected 2 details, got %d", len(err.Details))
	}
	if err.HTTPStatus != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", err.HTTPStatus)
	}
}
