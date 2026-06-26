package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Company struct {
	ID        string          `db:"id"         json:"id"`
	Name      string          `db:"name"       json:"name"`
	Slug      string          `db:"slug"       json:"slug"`
	LogoURL   sql.NullString  `db:"logo_url"   json:"logo_url,omitempty"`
	Website   string          `db:"website"    json:"website"`
	Industry  string          `db:"industry"   json:"industry"`
	Size      string          `db:"size"       json:"size"`
	CreatedBy string          `db:"created_by" json:"created_by"`
	Settings  json.RawMessage `db:"settings"   json:"settings,omitempty"`
	CreatedAt time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt time.Time       `db:"updated_at" json:"updated_at"`
	DeletedAt sql.NullTime    `db:"deleted_at" json:"-"`
}
