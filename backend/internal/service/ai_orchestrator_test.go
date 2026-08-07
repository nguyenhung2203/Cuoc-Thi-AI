package service

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	apierrors "backend/internal/pkg/errors"
)

// Existing AppErrors (insufficient-data 422, circuit-breaker 502) must pass
// through unchanged — not get re-wrapped into a generic 502.
func TestAIServiceError_PassesThroughAppError(t *testing.T) {
	orig := apierrors.NewValidation("insufficient data", []string{"need more context"})
	got := aiServiceError(orig, "AI call failed")

	appErr, ok := apierrors.IsAppError(got)
	if !ok {
		t.Fatal("expected an AppError")
	}
	if appErr.Code != apierrors.VALIDATION_ERROR || appErr.HTTPStatus != http.StatusUnprocessableEntity {
		t.Fatalf("got %s/%d, want VALIDATION_ERROR/422 (pass-through)", appErr.Code, appErr.HTTPStatus)
	}
}

func TestAIServiceError_WrapsRawError(t *testing.T) {
	got := aiServiceError(errors.New("dial tcp: i/o timeout"), "AI service call failed")

	appErr, ok := apierrors.IsAppError(got)
	if !ok {
		t.Fatal("raw error should become an AppError")
	}
	if appErr.Code != apierrors.AI_SERVICE_ERROR || appErr.HTTPStatus != http.StatusBadGateway {
		t.Fatalf("got %s/%d, want AI_SERVICE_ERROR/502", appErr.Code, appErr.HTTPStatus)
	}
}

// fmt.Errorf("%w") chains around an AppError must unwrap to the original —
// no double wrapping, no downgrade to 500.
func TestAIServiceError_UnwrapsNestedAppError(t *testing.T) {
	inner := apierrors.NewAIServiceError("upstream down")
	wrapped := fmt.Errorf("ai orchestrator failed: %w", inner)

	got := aiServiceError(wrapped, "outer message")
	appErr, ok := apierrors.IsAppError(got)
	if !ok {
		t.Fatal("expected an AppError")
	}
	if appErr != inner {
		t.Fatalf("expected the inner AppError to pass through, got %+v", appErr)
	}
}

func TestAIServiceError_NilIsNil(t *testing.T) {
	if got := aiServiceError(nil, "x"); got != nil {
		t.Fatalf("nil in must be nil out, got %v", got)
	}
}
