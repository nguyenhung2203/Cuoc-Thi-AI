package models

import (
	"database/sql"
	"time"
)

type MockInterview struct {
	ID          string         `db:"id" json:"id"`
	UserID      string         `db:"user_id" json:"user_id"`
	TargetRole  string         `db:"target_role" json:"target_role"`
	TargetLevel sql.NullString `db:"target_level" json:"target_level"`
	CVFileID    sql.NullString `db:"cv_file_id" json:"cv_file_id"`
	Status      string         `db:"status" json:"status"` // draft, active, completed, cancelled
	StartedAt   sql.NullTime   `db:"started_at" json:"started_at"`
	EndedAt     sql.NullTime   `db:"ended_at" json:"ended_at"`
	FinalScore  sql.NullFloat64 `db:"final_score" json:"final_score"`
	FeedbackJSON JSONB         `db:"feedback_json" json:"feedback_json"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at" json:"updated_at"`
}

type MockInterviewMessage struct {
	ID              string         `db:"id" json:"id"`
	MockInterviewID string         `db:"mock_interview_id" json:"mock_interview_id"`
	SenderType      string         `db:"sender_type" json:"sender_type"` // ai, candidate, system
	Content         string         `db:"content" json:"content"`
	QuestionType    sql.NullString `db:"question_type" json:"question_type"`
	ScoreJSON       JSONB          `db:"score_json" json:"score_json"`
	CreatedAt       time.Time      `db:"created_at" json:"created_at"`
}
