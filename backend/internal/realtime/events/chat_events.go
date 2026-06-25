package events

import "time"

// ── Client → Server ──────────────────────────────────────────────────────────

// ChatSendPayload is sent by any participant to post a message.
type ChatSendPayload struct {
	Message    string         `json:"message"`
	Visibility ChatVisibility `json:"visibility"` // "room" | "recruiter_only"
}

// NoteCreatePayload is used by the recruiter to save an internal note.
type NoteCreatePayload struct {
	Content string   `json:"content"`
	Tags    []string `json:"tags,omitempty"`
}

// QuestionMarkAskedPayload marks a prepared question as asked.
type QuestionMarkAskedPayload struct {
	QuestionID string `json:"question_id"`
	AskedAtMs  int64  `json:"asked_at_ms"`
}

// ── Server → Client ──────────────────────────────────────────────────────────

// ChatMessagePayload is broadcast when a chat message is delivered.
type ChatMessagePayload struct {
	MessageID           string          `json:"message_id"`
	SenderParticipantID string          `json:"sender_participant_id"`
	SenderName          string          `json:"sender_name"`
	SenderType          ParticipantType `json:"sender_type"`
	Message             string          `json:"message"`
	Visibility          ChatVisibility  `json:"visibility"`
	CreatedAt           time.Time       `json:"created_at"`
}
