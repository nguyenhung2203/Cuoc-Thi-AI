package ai

import (
	"context"
	"testing"

	apierrors "backend/internal/pkg/errors"
)

// fakeOrchestrator implements the Orchestrator interface with canned output.
type fakeOrchestrator struct {
	data []byte
	err  error
}

func (f *fakeOrchestrator) CallAI(ctx context.Context, templateName, companyID string, variables map[string]string) ([]byte, error) {
	return f.data, f.err
}

// K-JOB-11 regression: when the orchestrator fails, AnalyzeJD must return
// (nil, err) with the AppError surviving the %w wrap — the caller depends on
// this to map a 502 instead of panicking on a nil result.
func TestAnalyzeJD_OrchestratorError_ReturnsNilResult(t *testing.T) {
	upstream := apierrors.NewAIServiceError("upstream down")
	a := NewJDAnalyzer(&fakeOrchestrator{err: upstream})

	result, err := a.AnalyzeJD(context.Background(), "jd text", "title", "senior", "eng", "company-1")
	if result != nil {
		t.Fatalf("result should be nil on error, got %+v", result)
	}
	if err == nil {
		t.Fatal("expected an error")
	}
	appErr, ok := apierrors.IsAppError(err)
	if !ok {
		t.Fatal("AppError must survive the %w wrap (errors.As)")
	}
	if appErr.Code != apierrors.AI_SERVICE_ERROR {
		t.Fatalf("got code %s, want AI_SERVICE_ERROR", appErr.Code)
	}
}

func TestAnalyzeJD_InvalidJSON_ReturnsError(t *testing.T) {
	a := NewJDAnalyzer(&fakeOrchestrator{data: []byte("{not json")})

	result, err := a.AnalyzeJD(context.Background(), "jd", "t", "l", "d", "c")
	if err == nil {
		t.Fatal("expected an unmarshal error")
	}
	if result != nil {
		t.Fatalf("result should be nil, got %+v", result)
	}
}

func TestAnalyzeJD_Success_ParsesRubric(t *testing.T) {
	payload := `{
		"summary": "solid role",
		"suggested_rubric": [
			{"name": "Go", "weight": 60, "description": "lang"},
			{"name": "SQL", "weight": 40, "description": "db"}
		]
	}`
	a := NewJDAnalyzer(&fakeOrchestrator{data: []byte(payload)})

	result, err := a.AnalyzeJD(context.Background(), "jd", "t", "l", "d", "c")
	if err != nil {
		t.Fatalf("AnalyzeJD: %v", err)
	}
	if result.Summary != "solid role" {
		t.Errorf("Summary = %q", result.Summary)
	}
	if len(result.SuggestedRubric) != 2 {
		t.Fatalf("rubric len = %d, want 2", len(result.SuggestedRubric))
	}
}
