package models

import (
	"database/sql"
	"time"
)

type InterviewTemplate struct {
	ID               string         `db:"id" json:"id"`
	CompanyID        sql.NullString `db:"company_id" json:"company_id"`
	Name             string         `db:"name" json:"name"`
	Type             string         `db:"type" json:"type"`
	DurationMinutes  int            `db:"duration_minutes" json:"duration_minutes"`
	Description      sql.NullString `db:"description" json:"description"`
	ConfigJSON       JSONB          `db:"config_json" json:"config_json"`
	CreatedBy        sql.NullString `db:"created_by" json:"created_by"`
	CreatedAt        time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time      `db:"updated_at" json:"updated_at"`
}
