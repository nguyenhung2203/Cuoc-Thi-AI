package models

import (
	"database/sql"
	"time"
)

// JobCandidate maps the job_candidates table (job ↔ candidate join).
type JobCandidate struct {
	ID             string          `db:"id"              json:"id"`
	CompanyID      string          `db:"company_id"      json:"company_id"`
	JobID          string          `db:"job_id"          json:"job_id"`
	CandidateID    string          `db:"candidate_id"    json:"candidate_id"`
	PipelineStatus string          `db:"pipeline_status" json:"pipeline_status"`
	FitScore       sql.NullFloat64 `db:"fit_score"       json:"fit_score,omitempty"`
	AIMatchJSON    JSONB           `db:"ai_match_json"   json:"ai_match_json,omitempty"`
	AppliedAt      sql.NullTime    `db:"applied_at"      json:"applied_at,omitempty"`
	CreatedBy      sql.NullString  `db:"created_by"      json:"created_by,omitempty"`
	CreatedAt      time.Time       `db:"created_at"      json:"created_at"`
	UpdatedAt      time.Time       `db:"updated_at"      json:"updated_at"`
}
