package models

import (
	"database/sql"
	"time"
)

type InterviewTranscript struct {
	ID             string         `json:"id" db:"id"`
	InterviewID    string         `json:"interview_id" db:"interview_id"`
	SpeakerID      string         `json:"speaker_id" db:"speaker_id"` // user_id or candidate_id
	SpeakerRole    string         `json:"speaker_role" db:"speaker_role"` // recruiter, candidate, ai
	StartTime      float64        `json:"start_time" db:"start_time"`
	EndTime        float64        `json:"end_time" db:"end_time"`
	Content        string         `json:"content" db:"content"`
	IsFinal        bool           `json:"is_final" db:"is_final"`
	Language       string         `json:"language" db:"language"`
	EditedContent  sql.NullString `json:"edited_content" db:"edited_content"`
	EditedBy       sql.NullString `json:"edited_by" db:"edited_by"`
	EditedAt       sql.NullTime   `json:"edited_at" db:"edited_at"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at" db:"updated_at"`
}
