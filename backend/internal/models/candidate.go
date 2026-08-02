package models

import (
	"database/sql"
	"time"
)

type CandidateStatus string

const (
	CandidateStatusNew          CandidateStatus = "new"
	CandidateStatusScreening    CandidateStatus = "screening"
	CandidateStatusInvited      CandidateStatus = "invited"
	CandidateStatusInterviewing CandidateStatus = "interviewing"
	CandidateStatusCompleted    CandidateStatus = "completed"
	CandidateStatusPassed       CandidateStatus = "passed"
	CandidateStatusRejected     CandidateStatus = "rejected"
	CandidateStatusTalentPool   CandidateStatus = "talent_pool"
)

type Candidate struct {
	ID             string          `db:"id"              json:"id"`
	CompanyID      string          `db:"company_id"      json:"company_id"`
	UserID         sql.NullString  `db:"user_id"         json:"user_id,omitempty"`
	FullName       string          `db:"full_name"       json:"full_name"`
	Email          string          `db:"email"           json:"email"`
	Phone          sql.NullString  `db:"phone"           json:"phone,omitempty"`
	AvatarURL      sql.NullString  `db:"avatar_url"      json:"avatar_url,omitempty"`
	CVFileID       sql.NullString  `db:"cv_file_id"      json:"cv_file_id,omitempty"`
	ParsedCVJSON   JSONB           `db:"parsed_cv_json"  json:"parsed_cv_json,omitempty"`
	AICVSummary    sql.NullString  `db:"ai_cv_summary"   json:"ai_cv_summary,omitempty"`
	CVAIStatus     string          `db:"cv_ai_status"     json:"cv_ai_status"`
	CVAIError      sql.NullString  `db:"cv_ai_error"      json:"cv_ai_error,omitempty"`
	CVAIAttempts   int             `db:"cv_ai_attempts"   json:"cv_ai_attempts"`
	CVAIUpdatedAt  sql.NullTime    `db:"cv_ai_updated_at" json:"cv_ai_updated_at,omitempty"`
	Source         sql.NullString  `db:"source"          json:"source,omitempty"`
	Status         CandidateStatus `db:"status"          json:"status"`
	Tags           JSONB           `db:"tags"            json:"tags,omitempty"`
	CreatedBy      sql.NullString  `db:"created_by"      json:"created_by,omitempty"`
	CreatedAt      time.Time       `db:"created_at"      json:"created_at"`
	UpdatedAt      time.Time       `db:"updated_at"      json:"updated_at"`
	DeletedAt      sql.NullTime    `db:"deleted_at"      json:"-"`
	LatestJobID    sql.NullString  `db:"latest_job_id"   json:"-"`
	LatestJobTitle sql.NullString  `db:"latest_job_title" json:"-"`
	CVOriginalName sql.NullString  `db:"cv_original_name" json:"-"`
}
