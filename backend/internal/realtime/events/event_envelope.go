package events

import (
	"encoding/json"
	"time"
)

// Envelope is the standard wrapper for every WebSocket message
// in both directions (client→server and server→client).
type Envelope struct {
	Event       string          `json:"event"`
	RequestID   string          `json:"request_id,omitempty"`
	RoomID      string          `json:"room_id,omitempty"`
	InterviewID string          `json:"interview_id,omitempty"`
	SentAt      time.Time       `json:"sent_at,omitempty"`
	ServerTime  time.Time       `json:"server_time,omitempty"`
	Payload     json.RawMessage `json:"payload"`
}

// NewEnvelope builds a server-originated envelope with the current timestamp.
func NewEnvelope(event, requestID, roomID, interviewID string, payload any) (*Envelope, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &Envelope{
		Event:       event,
		RequestID:   requestID,
		RoomID:      roomID,
		InterviewID: interviewID,
		ServerTime:  time.Now().UTC(),
		Payload:     raw,
	}, nil
}

// MustEnvelope is like NewEnvelope but panics on marshal error (use only in tests).
func MustEnvelope(event, requestID, roomID, interviewID string, payload any) *Envelope {
	env, err := NewEnvelope(event, requestID, roomID, interviewID, payload)
	if err != nil {
		panic(err)
	}
	return env
}

// ParsePayload decodes the raw Payload field into the given target struct.
func (e *Envelope) ParsePayload(target any) error {
	return json.Unmarshal(e.Payload, target)
}

// ToJSON serialises the envelope to JSON bytes.
func (e *Envelope) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
