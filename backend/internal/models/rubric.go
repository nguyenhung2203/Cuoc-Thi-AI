package models

import (
	"database/sql"
	"time"
)

type Rubric struct {
	ID          string         `db:"id" json:"id"`
	CompanyID   string         `db:"company_id" json:"company_id"`
	JobID       sql.NullString `db:"job_id" json:"job_id"`
	Name        string         `db:"name" json:"name"`
	Description sql.NullString `db:"description" json:"description"`
	TotalWeight float64        `db:"total_weight" json:"total_weight"`
	CreatedBy   string         `db:"created_by" json:"created_by"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at" json:"updated_at"`
}

type RubricCriteria struct {
	ID           string         `db:"id" json:"id"`
	RubricID     string         `db:"rubric_id" json:"rubric_id"`
	Name         string         `db:"name" json:"name"`
	Description  sql.NullString `db:"description" json:"description"`
	Weight       float64        `db:"weight" json:"weight"`
	MinScore     int            `db:"min_score" json:"min_score"`
	MaxScore     int            `db:"max_score" json:"max_score"`
	ScoringGuide JSONB          `db:"scoring_guide" json:"scoring_guide"`
	OrderIndex   int            `db:"order_index" json:"order_index"`
	CreatedAt    time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at" json:"updated_at"`
}
