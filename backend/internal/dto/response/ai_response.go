package response

import "time"

type SuggestFollowUpResponse struct {
	SuggestedQuestion string             `json:"suggested_question"`
	Reason            string             `json:"reason"`
	TargetSkill       string             `json:"target_skill"`
	Priority          string             `json:"priority"`
	Confidence        float64            `json:"confidence"`
	CreatedAt         time.Time          `json:"created_at"`
	ExpiresAt         *time.Time         `json:"expires_at,omitempty"`
	TranscriptWindow  []TranscriptItem   `json:"transcript_window,omitempty"`
}

type TranscriptItem struct {
	ID       string `json:"id"`
	Speaker  string `json:"speaker"`
	Content  string `json:"content"`
	IsFinal  bool   `json:"is_final"`
}
