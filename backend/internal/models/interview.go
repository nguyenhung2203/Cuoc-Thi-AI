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
	RoomID           sql.NullString `json:"room_id" db:"room_id"`
	InviteTokenHash  sql.NullString `json:"-" db:"invite_token_hash"` // Do not leak in JSON
	InviteExpiresAt  sql.NullTime   `json:"invite_expires_at" db:"invite_expires_at"`
	ConsentRecording bool           `json:"consent_recording" db:"consent_recording"`
	ConsentAI        bool           `json:"consent_ai" db:"consent_ai"`
	CreatedBy        sql.NullString `json:"created_by" db:"created_by"`
	CreatedAt        time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at" db:"updated_at"`
}
