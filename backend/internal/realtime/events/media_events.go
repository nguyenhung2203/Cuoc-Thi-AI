package events

// ── Client → Server ──────────────────────────────────────────────────────────

// MediaStatusPayload is sent whenever the client toggles mic/camera/screen.
type MediaStatusPayload struct {
	MicEnabled    bool `json:"mic_enabled"`
	CameraEnabled bool `json:"camera_enabled"`
	ScreenSharing bool `json:"screen_sharing"`
}

// ── Server → Client ──────────────────────────────────────────────────────────

// MediaStatusChangedPayload is broadcast to all participants in the room.
type MediaStatusChangedPayload struct {
	ParticipantID string `json:"participant_id"`
	MicEnabled    bool   `json:"mic_enabled"`
	CameraEnabled bool   `json:"camera_enabled"`
	ScreenSharing bool   `json:"screen_sharing"`
}
