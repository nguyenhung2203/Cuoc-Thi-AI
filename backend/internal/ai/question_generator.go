package ai

import (
	"context"
	"encoding/json"
	"fmt"
)

type QuestionGenerator struct {
	orchestrator Orchestrator
}

func NewQuestionGenerator(orchestrator Orchestrator) *QuestionGenerator {
	return &QuestionGenerator{orchestrator: orchestrator}
}

func (g *QuestionGenerator) GenerateQuestions(ctx context.Context, job, candidateSummary, rubric, level, count, questionTypes, companyID string) (*QuestionGenerationResult, error) {
	variables := map[string]string{
		"job":                  job,
		"candidate_cv_summary": candidateSummary,
		"rubric":               rubric,
		"level":                level,
		"count":                count,
		"question_types":       questionTypes,
	}

	dataBytes, err := g.orchestrator.CallAI(ctx, "generate_questions", companyID, variables)
	if err != nil {
		return nil, fmt.Errorf("ai orchestrator failed: %w", err)
	}

	var result QuestionGenerationResult
	if err := json.Unmarshal(dataBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal question generation result: %w", err)
	}

	return &result, nil
}
