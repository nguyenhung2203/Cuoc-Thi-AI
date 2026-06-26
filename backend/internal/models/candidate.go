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
	ID            string          `db:"id"              json:"id"`
	CompanyID     string          `db:"company_id"      json:"company_id"`
	UserID        sql.NullString  `db:"user_id"         json:"user_id,omitempty"`
	FullName      string          `db:"full_name"       json:"full_name"`
	Email         string          `db:"email"           json:"email"`
	Phone         sql.NullString  `db:"phone"           json:"phone,omitempty"`
	AvatarURL     sql.NullString  `db:"avatar_url"      json:"avatar_url,omitempty"`
	CVFileID      sql.NullString  `db:"cv_file_id"      json:"cv_file_id,omitempty"`
	ParsedCVJSON  JSONB           `db:"parsed_cv_json"  json:"parsed_cv_json,omitempty"`
	AICVSummary   sql.NullString  `db:"ai_cv_summary"   json:"ai_cv_summary,omitempty"`
	Source        sql.NullString  `db:"source"          json:"source,omitempty"`
	Status        CandidateStatus `db:"status"          json:"status"`
	Tags          JSONB           `db:"tags"            json:"tags,omitempty"`
	CreatedBy     sql.NullString  `db:"created_by"      json:"created_by,omitempty"`
	CreatedAt     time.Time       `db:"created_at"      json:"created_at"`
	UpdatedAt     time.Time       `db:"updated_at"      json:"updated_at"`
	DeletedAt     sql.NullTime    `db:"deleted_at"      json:"-"`
}
