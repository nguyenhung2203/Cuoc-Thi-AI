package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

type AuditLog struct {
	ID          string         `db:"id" json:"id"`
	CompanyID   sql.NullString `db:"company_id" json:"-"`
	ActorUserID sql.NullString `db:"actor_user_id" json:"-"`
	ActorRole   sql.NullString `db:"actor_role" json:"-"`
	Action      string         `db:"action" json:"action"`
	ResourceType string        `db:"resource_type" json:"resource_type"`
	ResourceID  sql.NullString `db:"resource_id" json:"-"`
	BeforeJSON  JSONB          `db:"before_json" json:"before_json"`
	AfterJSON   JSONB          `db:"after_json" json:"after_json"`
	IPAddress   sql.NullString `db:"ip_address" json:"-"`
	UserAgent   sql.NullString `db:"user_agent" json:"-"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at"`
}

func (a AuditLog) MarshalJSON() ([]byte, error) {
	type Alias AuditLog
	return json.Marshal(&struct {
		Alias
		CompanyID   string `json:"company_id"`
		ActorUserID string `json:"actor_user_id"`
		ActorRole   string `json:"actor_role"`
		ResourceID  string `json:"resource_id"`
		IPAddress   string `json:"ip_address"`
		UserAgent   string `json:"user_agent"`
	}{
		Alias:       (Alias)(a),
		CompanyID:   a.CompanyID.String,
		ActorUserID: a.ActorUserID.String,
		ActorRole:   a.ActorRole.String,
		ResourceID:  a.ResourceID.String,
		IPAddress:   a.IPAddress.String,
		UserAgent:   a.UserAgent.String,
	})
}

