package events

import "time"

// ── Client → Server ──────────────────────────────────────────────────────────

// TranscriptPartialPayload is sent by the STT gateway with a live (unstable) chunk.
type TranscriptPartialPayload struct {
	SpeakerType SpeakerType `json:"speaker_type"`
	Content     string      `json:"content"`
	StartTimeMs int64       `json:"start_time_ms"`
	EndTimeMs   int64       `json:"end_time_ms"`
	Confidence  float64     `json:"confidence"`
}

// TranscriptFinalPayload is sent when a speech segment is fully recognised.
type TranscriptFinalPayload struct {
	SpeakerType SpeakerType `json:"speaker_type"`
	Content     string      `json:"content"`
	StartTimeMs int64       `json:"start_time_ms"`
	EndTimeMs   int64       `json:"end_time_ms"`
	Confidence  float64     `json:"confidence"`
}

// ── Server → Client ──────────────────────────────────────────────────────────

// TranscriptUpdatePayload is broadcast to participants whenever a new transcript
// chunk (partial or final) is available.
type TranscriptUpdatePayload struct {
	TranscriptID string      `json:"transcript_id"`
	SpeakerType  SpeakerType `json:"speaker_type"`
	SpeakerName  string      `json:"speaker_name"`
	Content      string      `json:"content"`
	StartTimeMs  int64       `json:"start_time_ms"`
	EndTimeMs    int64       `json:"end_time_ms"`
	Confidence   float64     `json:"confidence"`
	IsFinal      bool        `json:"is_final"`
	CreatedAt    time.Time   `json:"created_at"`
}
