package models

import (
	"database/sql"
	"time"
)

type Notification struct {
	ID        string         `json:"id" db:"id"`
	UserID    string         `json:"user_id" db:"user_id"`
	Type      string         `json:"type" db:"type"`
	Title     string         `json:"title" db:"title"`
	Content   sql.NullString `json:"content" db:"content"`
	DataJSON  sql.NullString `json:"data_json" db:"data_json"`
	ReadAt    sql.NullTime   `json:"read_at" db:"read_at"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
}
