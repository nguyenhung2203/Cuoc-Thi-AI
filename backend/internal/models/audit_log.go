package models

import (
	"database/sql"
	"time"
)

type AuditLog struct {
	ID          string         `db:"id" json:"id"`
	CompanyID   sql.NullString `db:"company_id" json:"company_id"`
	ActorUserID sql.NullString `db:"actor_user_id" json:"actor_user_id"`
	ActorRole   sql.NullString `db:"actor_role" json:"actor_role"`
	Action      string         `db:"action" json:"action"`
	ResourceType string        `db:"resource_type" json:"resource_type"`
	ResourceID  sql.NullString `db:"resource_id" json:"resource_id"`
	BeforeJSON  JSONB          `db:"before_json" json:"before_json"`
	AfterJSON   JSONB          `db:"after_json" json:"after_json"`
	IPAddress   sql.NullString `db:"ip_address" json:"ip_address"`
	UserAgent   sql.NullString `db:"user_agent" json:"user_agent"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at"`
}

