package models

import (
	"database/sql"
	"time"
)

type InterviewRoom struct {
	ID               string         `json:"id" db:"id"`
	InterviewID      string         `json:"interview_id" db:"interview_id"`
	RoomCode         string         `json:"room_code" db:"room_code"`
	Status           string         `json:"status" db:"status"` // waiting, active, closed
	Provider         sql.NullString `json:"provider" db:"provider"`
	ConnectionConfig []byte         `json:"connection_config" db:"connection_config"` // jsonb
	OpenedAt         sql.NullTime   `json:"opened_at" db:"opened_at"`
	ClosedAt         sql.NullTime   `json:"closed_at" db:"closed_at"`
	CreatedAt        time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at" db:"updated_at"`
}
