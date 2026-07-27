package events

import "time"

// ── Client → Server ──────────────────────────────────────────────────────────

// ChatSendPayload is sent by any participant to post a message.
type ChatSendPayload struct {
	MessageID  string         `json:"message_id,omitempty"`
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

// NoteCreatedPayload is broadcast (recruiters only) when a note is saved.
type NoteCreatedPayload struct {
	NoteID      string    `json:"note_id"`
	AuthorID    string    `json:"author_id"`
	AuthorName  string    `json:"author_name"`
	Content     string    `json:"content"`
	Tags        []string  `json:"tags,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// QuestionAskedPayload is broadcast (recruiters only) when a prepared question
// is marked as asked.
type QuestionAskedPayload struct {
	QuestionID string    `json:"question_id"`
	MarkedByID string    `json:"marked_by_id"`
	AskedAtMs  int64     `json:"asked_at_ms"`
	MarkedAt   time.Time `json:"marked_at"`
}
