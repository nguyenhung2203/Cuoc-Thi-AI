package models

import (
	"database/sql"
	"time"
)

// AIPromptTemplate represents the ai_prompt_templates table.
type AIPromptTemplate struct {
	ID              string          `db:"id" json:"id"`
	CompanyID       sql.NullString  `db:"company_id" json:"company_id"`
	Name            string          `db:"name" json:"name"`
	Version         int             `db:"version" json:"version"`
	Content         string          `db:"content" json:"content"`
	VariablesSchema JSONB `db:"variables_schema" json:"variables_schema"`
	Model           string          `db:"model" json:"model"`
	Params          JSONB `db:"params" json:"params"`
	IsActive        bool            `db:"is_active" json:"is_active"`
	CreatedBy       sql.NullString  `db:"created_by" json:"created_by"`
	CreatedAt       time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time       `db:"updated_at" json:"updated_at"`
	DeletedAt       sql.NullTime    `db:"deleted_at" json:"deleted_at"`
}

// AIRequestLog represents the ai_request_logs table.
type AIRequestLog struct {
	ID              string          `db:"id" json:"id"`
	CompanyID       string          `db:"company_id" json:"company_id"`
	TemplateID      sql.NullString  `db:"template_id" json:"template_id"`
	TemplateVersion sql.NullInt32   `db:"template_version" json:"template_version"`
	InterviewID     sql.NullString  `db:"interview_id" json:"interview_id"`
	JobID           sql.NullString  `db:"job_id" json:"job_id"`
	CandidateID     sql.NullString  `db:"candidate_id" json:"candidate_id"`
	InputJSON       JSONB `db:"input_json" json:"input_json"`
	OutputJSON      JSONB `db:"output_json" json:"output_json"`
	LatencyMs       sql.NullInt32   `db:"latency_ms" json:"latency_ms"`
	TokensIn        sql.NullInt32   `db:"tokens_in" json:"tokens_in"`
	TokensOut       sql.NullInt32   `db:"tokens_out" json:"tokens_out"`
	Cost            sql.NullFloat64 `db:"cost" json:"cost"`
	Status          string          `db:"status" json:"status"`
	Error           sql.NullString  `db:"error" json:"error"`
	CreatedAt       time.Time       `db:"created_at" json:"created_at"`
}
