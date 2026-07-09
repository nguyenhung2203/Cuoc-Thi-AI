package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"backend/internal/pkg/utils"
)

// SuggestFollowUpResult matches the output JSON of suggest_follow_up template.
type SuggestFollowUpResult struct {
	SuggestedQuestion string  `json:"suggested_question"`
	Reason            string  `json:"reason"`
	TargetSkill       string  `json:"target_skill"`
	Priority          string  `json:"priority"`
	Confidence        float64 `json:"confidence"`
	InsufficientData  bool    `json:"insufficient_data,omitempty"`
}

// SuggestFollowUpAnalyzer calls the AI to generate a follow-up question
// based on recent transcript context.
type SuggestFollowUpAnalyzer struct {
	orchestrator Orchestrator
}

func NewSuggestFollowUpAnalyzer(orchestrator Orchestrator) *SuggestFollowUpAnalyzer {
	return &SuggestFollowUpAnalyzer{orchestrator: orchestrator}
}

func (a *SuggestFollowUpAnalyzer) SuggestFollowUp(
	ctx context.Context,
	transcriptWindow string,
	jobContext string,
	focus string,
	companyID string,
) (*SuggestFollowUpResult, error) {
	variables := map[string]string{
		"transcript_window": transcriptWindow,
		"job_context":       jobContext,
		"focus":             focus,
	}

	dataBytes, err := a.orchestrator.CallAI(ctx, "suggest_follow_up", companyID, variables)
	if err != nil {
		return nil, fmt.Errorf("ai orchestrator failed: %w", err)
	}

	var result SuggestFollowUpResult
	cleanJSON := utils.CleanJSON(string(dataBytes))
	if err := json.Unmarshal([]byte(cleanJSON), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal suggest follow-up result: %w", err)
	}

	return &result, nil
}
