package response

import (
	"encoding/json"
	"time"
)

type InterviewTemplateItem struct {
	ID              string          `json:"id"`
	CompanyID       *string         `json:"company_id,omitempty"`
	Name            string          `json:"name"`
	Type            string          `json:"type"`
	DurationMinutes int             `json:"duration_minutes"`
	Description     *string         `json:"description,omitempty"`
	Config          json.RawMessage `json:"config,omitempty"`
	CreatedBy       *string         `json:"created_by,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}
