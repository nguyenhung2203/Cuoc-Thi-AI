package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"backend/internal/pkg/utils"
)

// Orchestrator interface to avoid circular dependency
type Orchestrator interface {
	CallAI(ctx context.Context, templateName, companyID string, variables map[string]string) ([]byte, error)
}

type JDAnalyzer struct {
	orchestrator Orchestrator
}

func NewJDAnalyzer(orchestrator Orchestrator) *JDAnalyzer {
	return &JDAnalyzer{orchestrator: orchestrator}
}

func (a *JDAnalyzer) AnalyzeJD(ctx context.Context, jobDescription, jobTitle, jobLevel, department, companyID string) (*JDAnalysisResult, error) {
	variables := map[string]string{
		"job_description": jobDescription,
		"job_title":       jobTitle,
		"job_level":       jobLevel,
		"department":      department,
	}

	dataBytes, err := a.orchestrator.CallAI(ctx, "analyze_jd", companyID, variables)
	if err != nil {
		return nil, fmt.Errorf("ai orchestrator failed: %w", err)
	}

	var result JDAnalysisResult
	cleanJSON := utils.CleanJSON(string(dataBytes))
	if err := json.Unmarshal([]byte(cleanJSON), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JD analysis result: %w", err)
	}

	return &result, nil
}
