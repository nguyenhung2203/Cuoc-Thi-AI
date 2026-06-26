package ai

import (
	"context"
	"encoding/json"
	"fmt"
)

type CVAnalyzer struct {
	orchestrator Orchestrator
}

func NewCVAnalyzer(orchestrator Orchestrator) *CVAnalyzer {
	return &CVAnalyzer{orchestrator: orchestrator}
}

func (a *CVAnalyzer) AnalyzeCV(ctx context.Context, cvText, jobContext, companyID string) (*CVAnalysisResult, error) {
	variables := map[string]string{
		"cv_text":     cvText,
		"job_context": jobContext,
	}

	dataBytes, err := a.orchestrator.CallAI(ctx, "analyze_cv", companyID, variables)
	if err != nil {
		return nil, fmt.Errorf("ai orchestrator failed: %w", err)
	}

	var result CVAnalysisResult
	if err := json.Unmarshal(dataBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal CV analysis result: %w", err)
	}

	return &result, nil
}
