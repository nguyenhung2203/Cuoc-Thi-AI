package models

import (
	"database/sql"
	"time"
)

type Interview struct {
	ID               string         `json:"id" db:"id"`
	CompanyID        sql.NullString `json:"company_id" db:"company_id"`
	JobID            sql.NullString `json:"job_id" db:"job_id"`
	CandidateID      string         `json:"candidate_id" db:"candidate_id"`
	RecruiterID      sql.NullString `json:"recruiter_id" db:"recruiter_id"`
	TemplateID       sql.NullString `json:"template_id" db:"template_id"`
	RubricID         sql.NullString `json:"rubric_id" db:"rubric_id"`
	Mode             string         `json:"mode" db:"mode"` // real, mock
	Title            string         `json:"title" db:"title"`
	ScheduledAt      sql.NullTime   `json:"scheduled_at" db:"scheduled_at"`
	StartedAt        sql.NullTime   `json:"started_at" db:"started_at"`
	EndedAt          sql.NullTime   `json:"ended_at" db:"ended_at"`
	Status           string         `json:"status" db:"status"` // scheduled, waiting, active, paused, completed, cancelled, expired
	RoomID           sql.NullString `db:"room_id" json:"room_id"`
	InviteTokenHash  sql.NullString `db:"invite_token_hash" json:"-"`
	InviteExpiresAt  sql.NullTime   `db:"invite_expires_at" json:"invite_expires_at"`
	ConsentRecording bool           `db:"consent_recording" json:"consent_recording"`
	ConsentAI        bool           `db:"consent_ai" json:"consent_ai"`
	ReportStatus     string         `db:"report_status" json:"report_status"`
	CreatedBy        sql.NullString `db:"created_by" json:"created_by"`
	CreatedAt        time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time      `db:"updated_at" json:"updated_at"`
}
