package models

import "time"

type InterviewTranscript struct {
	ID            string     `json:"id" db:"id"`
	InterviewID   string     `json:"interview_id" db:"interview_id"`
	ParticipantID *string    `json:"participant_id,omitempty" db:"participant_id"`
	SpeakerType   string     `json:"speaker_type" db:"speaker_type"`
	SpeakerName   *string    `json:"speaker_name,omitempty" db:"speaker_name"`
	Content       string     `json:"content" db:"content"`
	Language      *string    `json:"language,omitempty" db:"language"`
	StartTimeMs   *int64     `json:"start_time_ms,omitempty" db:"start_time_ms"`
	EndTimeMs     *int64     `json:"end_time_ms,omitempty" db:"end_time_ms"`
	Confidence    *float64   `json:"confidence,omitempty" db:"confidence"`
	Source        string     `json:"source" db:"source"`
	IsFinal       bool       `json:"is_final" db:"is_final"`
	EditedContent *string    `json:"edited_content,omitempty" db:"edited_content"`
	EditedBy      *string    `json:"edited_by,omitempty" db:"edited_by"`
	EditedAt      *time.Time `json:"edited_at,omitempty" db:"edited_at"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
}
