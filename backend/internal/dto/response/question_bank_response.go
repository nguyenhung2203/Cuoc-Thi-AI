package response

import "time"

type QuestionBankItem struct {
	ID              string    `json:"id"`
	CompanyID       *string   `json:"company_id,omitempty"`
	JobID           *string   `json:"job_id,omitempty"`
	CreatedBy       *string   `json:"created_by,omitempty"`
	QuestionText    string    `json:"question_text"`
	QuestionType    string    `json:"question_type"`
	SkillTags       []string  `json:"skill_tags,omitempty"`
	Level           *string   `json:"level,omitempty"`
	ExpectedSignals []string  `json:"expected_signals,omitempty"`
	IsAIGenerated   bool      `json:"is_ai_generated"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
