package models

import (
	"database/sql"
	"time"
)

type UserRole string

const (
	RoleAdmin     UserRole = "admin"
	RoleRecruiter UserRole = "recruiter"
	RoleCandidate UserRole = "candidate"
)

type UserStatus string

const (
	UserStatusPending  UserStatus = "pending"
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusBlocked  UserStatus = "blocked"
)

type User struct {
	ID                 string         `db:"id"                json:"id"`
	Email              string         `db:"email"             json:"email"`
	PasswordHash       string         `db:"password_hash"     json:"-"`
	FullName           string         `db:"full_name"         json:"full_name"`
	AvatarURL          sql.NullString `db:"avatar_url"        json:"avatar_url,omitempty"`
	Role               UserRole       `db:"role"              json:"role"`
	Status             UserStatus     `db:"status"            json:"status"`
	Settings           JSONB          `db:"settings"          json:"settings,omitempty"`
	LastLoginAt        sql.NullTime   `db:"last_login_at"     json:"last_login_at,omitempty"`
	EmailVerifiedAt    sql.NullTime   `db:"email_verified_at" json:"email_verified_at,omitempty"`
	CreatedAt          time.Time      `db:"created_at"           json:"created_at"`
	UpdatedAt          time.Time      `db:"updated_at"           json:"updated_at"`
	DeletedAt          sql.NullTime   `db:"deleted_at"           json:"-"`
	VerificationFileID sql.NullString `db:"verification_file_id" json:"verification_file_id,omitempty"`
}
