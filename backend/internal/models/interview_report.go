package models

import (
	"database/sql"
	"time"
)

type InterviewReport struct {



	ID                   string         `db:"id" json:"id"`
	InterviewID          string         `db:"interview_id" json:"interview_id"`
	Summary              string         `db:"summary" json:"summary"`
	FinalScore           sql.NullFloat64 `db:"final_score" json:"final_score"`
	Recommendation       string         `db:"recommendation" json:"recommendation"`
	Strengths            JSONB          `db:"strengths" json:"strengths"`
	Weaknesses           JSONB          `db:"weaknesses" json:"weaknesses"`
	Risks                JSONB          `db:"risks" json:"risks"`
	EvidenceJSON         JSONB          `db:"evidence_json" json:"evidence_json"`
	AIReasoningSummary   sql.NullString `db:"ai_reasoning_summary" json:"ai_reasoning_summary"`
	RecruiterDecision    sql.NullString `db:"recruiter_decision" json:"recruiter_decision"`
	RecruiterComment     sql.NullString `db:"recruiter_comment" json:"recruiter_comment"`
	ReportJSON           JSONB          `db:"report_json" json:"report_json"`
	GeneratedBy          string         `db:"generated_by" json:"generated_by"`
	GeneratedAt          time.Time      `db:"generated_at" json:"generated_at"`
	CreatedAt            time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time      `db:"updated_at" json:"updated_at"`

}
