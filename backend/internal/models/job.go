package models

import (
	"database/sql"
	"time"
)

type JobStatus string

const (
	JobStatusDraft  JobStatus = "draft"
	JobStatusOpen   JobStatus = "open"
	JobStatusPaused JobStatus = "paused"
	JobStatusClosed JobStatus = "closed"
)

type Job struct {
	ID             string          `db:"id"                json:"id"`
	CompanyID      string          `db:"company_id"        json:"company_id"`
	Title          string          `db:"title"             json:"title"`
	Department     sql.NullString  `db:"department"        json:"department,omitempty"`
	Level          sql.NullString  `db:"level"             json:"level,omitempty"`
	Location       sql.NullString  `db:"location"          json:"location,omitempty"`
	EmploymentType sql.NullString  `db:"employment_type"   json:"employment_type,omitempty"`
	SalaryMin      sql.NullFloat64 `db:"salary_min"        json:"salary_min,omitempty"`
	SalaryMax      sql.NullFloat64 `db:"salary_max"        json:"salary_max,omitempty"`
	Currency       sql.NullString  `db:"currency"          json:"currency,omitempty"`
	Description    string          `db:"description"       json:"description"`
	Requirements   sql.NullString  `db:"requirements"      json:"requirements,omitempty"`
	Benefits       sql.NullString  `db:"benefits"          json:"benefits,omitempty"`
	Status         JobStatus       `db:"status"            json:"status"`
	AISummary      sql.NullString  `db:"ai_summary"        json:"ai_summary,omitempty"`
	AIAnalysisJSON JSONB           `db:"ai_analysis_json"  json:"ai_analysis_json,omitempty"`
	CreatedBy      string          `db:"created_by"        json:"created_by"`
	CreatedAt      time.Time       `db:"created_at"        json:"created_at"`
	UpdatedAt      time.Time       `db:"updated_at"        json:"updated_at"`
	DeletedAt      sql.NullTime    `db:"deleted_at"        json:"-"`
	CompanyName    sql.NullString  `db:"company_name"      json:"company_name,omitempty"`
	CompanyLogoURL sql.NullString  `db:"company_logo_url"  json:"company_logo_url,omitempty"`
}
