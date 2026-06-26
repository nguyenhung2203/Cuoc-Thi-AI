package events

import "time"

// ── Client → Server ──────────────────────────────────────────────────────────

// RoomJoinPayload is sent by a client when entering a room.
type RoomJoinPayload struct {
	ParticipantType ParticipantType `json:"participant_type"`
	DisplayName     string          `json:"display_name"`
	DeviceInfo      *DeviceInfo     `json:"device_info,omitempty"`
}

// DeviceInfo carries optional browser/OS metadata from the client.
type DeviceInfo struct {
	Browser string `json:"browser,omitempty"`
	OS      string `json:"os,omitempty"`
}

// RoomLeavePayload is sent by a client when leaving a room.
type RoomLeavePayload struct {
	Reason string `json:"reason,omitempty"`
}

// RoomHeartbeatPayload is sent periodically to signal the connection is alive.
type RoomHeartbeatPayload struct {
	ConnectionState ConnectionState `json:"connection_state"`
}

// ── Server → Client ──────────────────────────────────────────────────────────

// ParticipantInfo is a snapshot of a single participant's state, included in
// room:joined and room:presence_update events.
type ParticipantInfo struct {
	ParticipantID   string           `json:"participant_id"`
	DisplayName     string           `json:"display_name"`
	ParticipantType ParticipantType  `json:"participant_type"`
	ConnectionState ConnectionState  `json:"connection_state"`
	MediaStatus     *MediaStatusInfo `json:"media_status,omitempty"`
	JoinedAt        time.Time        `json:"joined_at,omitempty"`
	LastSeenAt      time.Time        `json:"last_seen_at,omitempty"`
}

// MediaStatusInfo is the mic/camera/screen state for a participant.
type MediaStatusInfo struct {
	MicEnabled    bool `json:"mic_enabled"`
	CameraEnabled bool `json:"camera_enabled"`
	ScreenSharing bool `json:"screen_sharing"`
}

// RoomJoinedPayload is the ACK sent back only to the participant who just joined.
type RoomJoinedPayload struct {
	ParticipantID   string            `json:"participant_id"`
	RoomStatus      RoomStatus        `json:"room_status"`
	InterviewStatus string            `json:"interview_status"`
	Participants    []ParticipantInfo `json:"participants"`
}

// RoomUserJoinedPayload is broadcast to everyone else when a new participant joins.
type RoomUserJoinedPayload struct {
	ParticipantID   string          `json:"participant_id"`
	DisplayName     string          `json:"display_name"`
	ParticipantType ParticipantType `json:"participant_type"`
	JoinedAt        time.Time       `json:"joined_at"`
}

// RoomUserLeftPayload is broadcast when a participant leaves or disconnects.
type RoomUserLeftPayload struct {
	ParticipantID string    `json:"participant_id"`
	Reason        string    `json:"reason,omitempty"`
	LeftAt        time.Time `json:"left_at"`
}

// RoomPresenceUpdatePayload is broadcast whenever the presence list changes.
type RoomPresenceUpdatePayload struct {
	Participants []ParticipantInfo `json:"participants"`
}
