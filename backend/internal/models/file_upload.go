package models

import (
	"database/sql"
	"time"
)

// File holds metadata for uploaded files (table: files).
// NOTE: storage_key must never be returned directly in API responses;
// generate a signed URL instead.
type File struct {
	ID           string         `db:"id"            json:"id"`
	CompanyID    sql.NullString `db:"company_id"    json:"company_id,omitempty"`
	OwnerUserID  sql.NullString `db:"owner_user_id" json:"owner_user_id,omitempty"`
	OriginalName string         `db:"original_name" json:"original_name"`
	StorageKey   string         `db:"storage_key"   json:"-"`
	MimeType     string         `db:"mime_type"     json:"mime_type"`
	SizeBytes    int64          `db:"size_bytes"    json:"size_bytes"`
	FileType     string         `db:"file_type"     json:"file_type"` // cv | audio_recording | avatar | attachment
	Checksum     sql.NullString `db:"checksum"      json:"checksum,omitempty"`
	ParsedJSON   JSONB          `db:"parsed_json"   json:"parsed_json,omitempty"`
	AISummary    sql.NullString `db:"ai_summary"    json:"ai_summary,omitempty"`
	ParseStatus  sql.NullString `db:"parse_status"  json:"parse_status,omitempty"`
	ParseError   sql.NullString `db:"parse_error"   json:"parse_error,omitempty"`
	ParsedAt     sql.NullTime   `db:"parsed_at"     json:"parsed_at,omitempty"`
	CreatedAt    time.Time      `db:"created_at"    json:"created_at"`
}
