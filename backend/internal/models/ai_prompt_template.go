package models

import (
	"database/sql"
	"time"
)

type AiPromptTemplate struct {
	ID              string         `db:"id"              json:"id"`
	Name            string         `db:"name"            json:"name"`
	Version         int            `db:"version"         json:"version"`
	Content         string         `db:"content"         json:"content"`
	VariablesSchema JSONB          `db:"variables_schema" json:"variables_schema,omitempty"`
	Model           string         `db:"model"           json:"model"`
	Params          JSONB          `db:"params"          json:"params,omitempty"`
	IsActive        bool           `db:"is_active"       json:"is_active"`
	CompanyID       sql.NullString `db:"company_id"      json:"company_id,omitempty"`
	CreatedBy       sql.NullString `db:"created_by"      json:"created_by,omitempty"`
	CreatedAt       time.Time      `db:"created_at"      json:"created_at"`
	UpdatedAt       time.Time      `db:"updated_at"      json:"updated_at"`
	DeletedAt       sql.NullTime   `db:"deleted_at"      json:"deleted_at,omitempty"`
}
