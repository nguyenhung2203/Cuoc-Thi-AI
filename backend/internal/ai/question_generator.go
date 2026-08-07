package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"backend/internal/pkg/utils"
)

type StructuredQuestionRequest struct {
	JobTitle          string         `json:"job_title"`
	JobDescription    string         `json:"job_description"`
	Requirements      []string       `json:"requirements,omitempty"`
	Level             string         `json:"level"`
	Mode              string         `json:"mode"`
	Language          string         `json:"language"`
	SkillTags         []string       `json:"skill_tags,omitempty"`
	Rubric            map[string]any `json:"rubric,omitempty"`
	QuestionCount     int            `json:"question_count"`
	PreviousQuestions []string       `json:"previous_questions,omitempty"`
	RecentAnswer      string         `json:"recent_answer,omitempty"`
}

type QuestionGenerator struct{ orchestrator Orchestrator }

type StructuredQuestionClient interface {
	GenerateStructuredQuestions(context.Context, string, StructuredQuestionRequest) (*QuestionGenerationResult, error)
	GenerateFollowUp(context.Context, string, StructuredQuestionRequest) (*QuestionGenerationResult, error)
}

func (g *QuestionGenerator) GenerateStructuredQuestions(ctx context.Context, companyID string, req StructuredQuestionRequest) (*QuestionGenerationResult, error) {
	client, ok := g.orchestrator.(StructuredQuestionClient)
	if !ok {
		return nil, fmt.Errorf("structured question client unavailable")
	}
	return client.GenerateStructuredQuestions(ctx, companyID, req)
}

func (g *QuestionGenerator) GenerateFollowUp(ctx context.Context, companyID string, req StructuredQuestionRequest) (*QuestionGenerationResult, error) {
	client, ok := g.orchestrator.(StructuredQuestionClient)
	if !ok {
		return nil, fmt.Errorf("structured question client unavailable")
	}
	return client.GenerateFollowUp(ctx, companyID, req)
}

func NewQuestionGenerator(orchestrator Orchestrator) *QuestionGenerator {
	return &QuestionGenerator{orchestrator: orchestrator}
}

func (g *QuestionGenerator) GenerateQuestions(ctx context.Context, job, candidateSummary, rubric, level, count, questionTypes, companyID string) (*QuestionGenerationResult, error) {
	variables := map[string]string{"job": job, "candidate_cv_summary": candidateSummary, "rubric": rubric, "level": level, "count": count, "question_types": questionTypes}
	dataBytes, err := g.orchestrator.CallAI(ctx, "generate_questions", companyID, variables)
	if err != nil {
		return nil, fmt.Errorf("ai orchestrator failed: %w", err)
	}
	var result QuestionGenerationResult
	if err := json.Unmarshal([]byte(utils.CleanJSON(string(dataBytes))), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal question generation result: %w", err)
	}
	return &result, nil
}
