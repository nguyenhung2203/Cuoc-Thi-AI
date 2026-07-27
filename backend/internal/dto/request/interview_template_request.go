package request

import "encoding/json"

type CreateInterviewTemplateRequest struct {
	Name             string          `json:"name" validate:"required"`
	Type             string          `json:"type" validate:"required"`
	DurationMinutes  int             `json:"duration_minutes"`
	Description      string          `json:"description"`
	Config           json.RawMessage `json:"config"`
}

type UpdateInterviewTemplateRequest struct {
	Name             *string         `json:"name"`
	Type             *string         `json:"type"`
	DurationMinutes  *int            `json:"duration_minutes"`
	Description      *string         `json:"description"`
	Config           json.RawMessage `json:"config"`
}
