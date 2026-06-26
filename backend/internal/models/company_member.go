package models

import (
	"database/sql"
	"time"
)

type CompanyMemberRole string

const (
	CompanyMemberRoleOwner  CompanyMemberRole = "owner"
	CompanyMemberRoleAdmin  CompanyMemberRole = "admin"
	CompanyMemberRoleMember CompanyMemberRole = "member"
	CompanyMemberRoleViewer CompanyMemberRole = "viewer"
)

type CompanyMemberStatus string

const (
	CompanyMemberStatusInvited  CompanyMemberStatus = "invited"
	CompanyMemberStatusActive   CompanyMemberStatus = "active"
	CompanyMemberStatusRemoved  CompanyMemberStatus = "removed"
)

type CompanyMember struct {
	ID         string              `db:"id"          json:"id"`
	CompanyID  string              `db:"company_id"  json:"company_id"`
	UserID     string              `db:"user_id"     json:"user_id"`
	Role       CompanyMemberRole   `db:"role"        json:"role"`
	Status     CompanyMemberStatus `db:"status"      json:"status"`
	InvitedBy  sql.NullString      `db:"invited_by"  json:"invited_by,omitempty"`
	JoinedAt   sql.NullTime        `db:"joined_at"   json:"joined_at,omitempty"`
	CreatedAt  time.Time           `db:"created_at"  json:"created_at"`
	UpdatedAt  time.Time           `db:"updated_at"  json:"updated_at"`
}
