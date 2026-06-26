package models

import "time"

type AuditLog struct {
	ID           uint64                 `json:"id"`
	Action       string                 `json:"action"`
	ActorID      string                 `json:"actor_id"`
	ActorType    string                 `json:"actor_type"`
	ResourceType string                 `json:"resource_type"`
	ResourceID   string                 `json:"resource_id"`
	CompanyID    string                 `json:"company_id"`
	Metadata     map[string]interface{} `json:"metadata"`
	IPAddress    string                 `json:"ip_address"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}
