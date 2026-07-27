package models

import (
	"database/sql"
	"time"
)

type QuestionBank struct {
	ID              string         `db:"id" json:"id"`
	CompanyID       sql.NullString `db:"company_id" json:"company_id"`
	JobID           sql.NullString `db:"job_id" json:"job_id"`
	CreatedBy       sql.NullString `db:"created_by" json:"created_by"`
	QuestionText    string         `db:"question_text" json:"question_text"`
	QuestionType    string         `db:"question_type" json:"question_type"`
	SkillTags       JSONB          `db:"skill_tags" json:"skill_tags"`
	Level           sql.NullString `db:"level" json:"level"`
	ExpectedSignals JSONB          `db:"expected_signals" json:"expected_signals"`
	IsAIGenerated   bool           `db:"is_ai_generated" json:"is_ai_generated"`
	CreatedAt       time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time      `db:"updated_at" json:"updated_at"`
}
