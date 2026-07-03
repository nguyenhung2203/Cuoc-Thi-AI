package models

import (
	"database/sql"
	"time"
)

type InterviewReport struct {
	ID                 string          `json:"id" db:"id"`
	InterviewID        string          `json:"interview_id" db:"interview_id"`
	Summary            string          `json:"summary" db:"summary"`
	FinalScore         sql.NullFloat64 `json:"final_score" db:"final_score"`
	Recommendation     string          `json:"recommendation" db:"recommendation"`
	Strengths          sql.NullString  `json:"strengths" db:"strengths"`
	Weaknesses         sql.NullString  `json:"weaknesses" db:"weaknesses"`
	Risks              sql.NullString  `json:"risks" db:"risks"`
	EvidenceJSON       sql.NullString  `json:"evidence_json" db:"evidence_json"`
	AIReasoningSummary sql.NullString  `json:"ai_reasoning_summary" db:"ai_reasoning_summary"`
	RecruiterDecision  sql.NullString  `json:"recruiter_decision" db:"recruiter_decision"`
	RecruiterComment   sql.NullString  `json:"recruiter_comment" db:"recruiter_comment"`
	ReportJSON         string          `json:"report_json" db:"report_json"`
	GeneratedBy        string          `json:"generated_by" db:"generated_by"`
	GeneratedAt        time.Time       `json:"generated_at" db:"generated_at"`
	CreatedAt          time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at" db:"updated_at"`
}
