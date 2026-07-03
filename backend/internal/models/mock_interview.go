package models

import (
	"database/sql"
	"time"
)

type MockInterview struct {
	ID           string          `json:"id" db:"id"`
	UserID       string          `json:"user_id" db:"user_id"`
	TargetRole   string          `json:"target_role" db:"target_role"`
	TargetLevel  sql.NullString  `json:"target_level" db:"target_level"`
	CVFileID     sql.NullString  `json:"cv_file_id" db:"cv_file_id"`
	Status       string          `json:"status" db:"status"`
	StartedAt    sql.NullTime    `json:"started_at" db:"started_at"`
	EndedAt      sql.NullTime    `json:"ended_at" db:"ended_at"`
	FinalScore   sql.NullFloat64 `json:"final_score" db:"final_score"`
	FeedbackJSON sql.NullString  `json:"feedback_json" db:"feedback_json"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at" db:"updated_at"`
}

type MockInterviewMessage struct {
	ID              string         `json:"id" db:"id"`
	MockInterviewID string         `json:"mock_interview_id" db:"mock_interview_id"`
	SenderType      string         `json:"sender_type" db:"sender_type"` // ai, candidate, system
	Content         string         `json:"content" db:"content"`
	QuestionType    sql.NullString `json:"question_type" db:"question_type"`
	ScoreJSON       sql.NullString `json:"score_json" db:"score_json"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
}
