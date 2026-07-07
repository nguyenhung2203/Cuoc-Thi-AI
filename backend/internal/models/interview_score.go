package models

import (
	"database/sql"
	"time"
)

type InterviewScore struct {
	ID                string         `db:"id" json:"id"`
	InterviewID       string         `db:"interview_id" json:"interview_id"`
	RubricCriterionID sql.NullString `db:"rubric_criterion_id" json:"rubric_criterion_id"`
	CriterionName     string         `db:"criterion_name" json:"criterion_name"`
	Score             sql.NullFloat64 `db:"score" json:"score"`
	MaxScore          float64        `db:"max_score" json:"max_score"`
	Weight            float64        `db:"weight" json:"weight"`
	WeightedScore     sql.NullFloat64 `db:"weighted_score" json:"weighted_score"`
	Evidence          sql.NullString `db:"evidence" json:"evidence"`
	AIComment         sql.NullString `db:"ai_comment" json:"ai_comment"`
	Confidence        sql.NullFloat64 `db:"confidence" json:"confidence"`
	Status            string         `db:"status" json:"status"` // scored/insufficient_evidence/manual_override
	ScoredBy          string         `db:"scored_by" json:"scored_by"` // ai/recruiter/system
	CreatedAt         time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time      `db:"updated_at" json:"updated_at"`
}
