package models

import "time"

// RefreshToken đại diện cho bảng refresh_tokens trong database
type RefreshToken struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	FamilyID  string    `json:"family_id" db:"family_id"`
	TokenHash string    `json:"token_hash" db:"token_hash"`
	IsRevoked bool      `json:"is_revoked" db:"is_revoked"`
	IPAddress *string   `json:"ip_address" db:"ip_address"`
	UserAgent *string   `json:"user_agent" db:"user_agent"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
