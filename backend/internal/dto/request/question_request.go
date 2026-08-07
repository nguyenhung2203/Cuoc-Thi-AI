package request

import "backend/internal/models"

type CreateQuestionRequest struct {
	JobID           string       `json:"job_id"`
	QuestionText    string       `json:"question_text" validate:"required"`
	QuestionType    string       `json:"question_type" validate:"required"`
	SkillTags       models.JSONB `json:"skill_tags"`
	Level           string       `json:"level"`
	ExpectedSignals models.JSONB `json:"expected_signals"`
}

type UpdateQuestionRequest struct {
	QuestionText    string       `json:"question_text"`
	QuestionType    string       `json:"question_type"`
	SkillTags       models.JSONB `json:"skill_tags"`
	Level           string       `json:"level"`
	ExpectedSignals models.JSONB `json:"expected_signals"`
}

type GenerateQuestionsRequest struct {
	JobID         string   `json:"job_id"`
	CandidateID   string   `json:"candidate_id,omitempty"`
	RubricID      string   `json:"rubric_id,omitempty"`
	Count         int      `json:"count" validate:"required,min=1,max=20"`
	Difficulty    string   `json:"difficulty"`
	Level         string   `json:"level"`
	Mode          string   `json:"mode"`
	QuestionTypes []string `json:"question_types"`
}
