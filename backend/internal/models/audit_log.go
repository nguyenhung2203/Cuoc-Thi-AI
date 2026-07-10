package models

import (
	"time"
)

type AuditLog struct {
	ID          string    `db:"id" json:"id"`
	CompanyID   *string   `db:"company_id" json:"company_id"`
	ActorUserID *string   `db:"actor_user_id" json:"actor_user_id"`
	ActorRole   *string   `db:"actor_role" json:"actor_role"`
	Action      string    `db:"action" json:"action"`
	ResourceType string   `db:"resource_type" json:"resource_type"`
	ResourceID  *string   `db:"resource_id" json:"resource_id"`
	BeforeJSON  JSONB     `db:"before_json" json:"before_json"`
	AfterJSON   JSONB     `db:"after_json" json:"after_json"`
	IPAddress   *string   `db:"ip_address" json:"ip_address"`
	UserAgent   *string   `db:"user_agent" json:"user_agent"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at"`
}

