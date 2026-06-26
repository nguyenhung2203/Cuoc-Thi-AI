package events

import "time"

// ── Client → Server ──────────────────────────────────────────────────────────

// InterviewStartPayload is sent by the recruiter to begin the session.
type InterviewStartPayload struct {
	ConsentRecording bool `json:"consent_recording"`
	ConsentAI        bool `json:"consent_ai"`
}

// InterviewEndPayload is sent by the recruiter to close the session.
type InterviewEndPayload struct {
	GenerateReport bool `json:"generate_report"`
}

// ── Server → Client ──────────────────────────────────────────────────────────

// InterviewStartedPayload is broadcast to all participants once the interview is active.
type InterviewStartedPayload struct {
	Status            string    `json:"status"`
	StartedAt         time.Time `json:"started_at"`
	TranscriptEnabled bool      `json:"transcript_enabled"`
	AIEnabled         bool      `json:"ai_enabled"`
}

// InterviewCompletedPayload is broadcast when the interview ends.
type InterviewCompletedPayload struct {
	Status       string    `json:"status"`
	EndedAt      time.Time `json:"ended_at"`
	ReportStatus string    `json:"report_status"` // "generating" | "ready" | "failed"
}

// InterviewPausedPayload is broadcast when the interview is paused.
type InterviewPausedPayload struct {
	Status   string    `json:"status"`
	PausedAt time.Time `json:"paused_at"`
}

// InterviewResumedPayload is broadcast when the interview is resumed.
type InterviewResumedPayload struct {
	Status    string    `json:"status"`
	ResumedAt time.Time `json:"resumed_at"`
}

// InterviewCancelledPayload is broadcast when the interview is cancelled.
type InterviewCancelledPayload struct {
	Status      string    `json:"status"`
	CancelledAt time.Time `json:"cancelled_at"`
	Reason      string    `json:"reason,omitempty"`
}

// RoomExpiredPayload is broadcast when the room waiting timer expires.
type RoomExpiredPayload struct {
	Status    string    `json:"status"`
	ExpiredAt time.Time `json:"expired_at"`
}

