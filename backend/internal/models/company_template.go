package models

import (
	"encoding/json"
	"time"
)

type CompanyTemplate struct {
	ID          string    `db:"id" json:"id"`
	CompanyID   string    `db:"company_id" json:"company_id"`
	Title       string    `db:"title" json:"title"`
	Type        string    `db:"type" json:"type"`
	Description string    `db:"description" json:"description"`
	Tags        string    `db:"tags" json:"-"` // JSON string from DB
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

func (t *CompanyTemplate) GetTags() []string {
	var tags []string
	if t.Tags != "" && t.Tags != "null" {
		_ = json.Unmarshal([]byte(t.Tags), &tags)
	}
	if tags == nil {
		tags = []string{}
	}
	return tags
}

func (t *CompanyTemplate) SetTags(tags []string) {
	b, _ := json.Marshal(tags)
	t.Tags = string(b)
}
