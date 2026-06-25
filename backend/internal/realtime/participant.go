package realtime

import (
	"time"

	"backend/internal/realtime/events"
)

// Participant represents a user inside an interview room.
// It links the in-memory ClientConnection to the room's presence state.
type Participant struct {
	ConnectionID    string
	UserID          string
	ParticipantType events.ParticipantType
	DisplayName     string
	ConnectionState events.ConnectionState
	MediaStatus     events.MediaStatusInfo
	JoinedAt        time.Time
	LastSeenAt      time.Time
	Connection      *ClientConnection // nil when the participant is offline/reconnecting
}

// ToInfo converts a Participant to the wire-format ParticipantInfo used in events.
func (p *Participant) ToInfo() events.ParticipantInfo {
	info := events.ParticipantInfo{
		ParticipantID:   p.ConnectionID,
		DisplayName:     p.DisplayName,
		ParticipantType: p.ParticipantType,
		ConnectionState: p.ConnectionState,
		JoinedAt:        p.JoinedAt,
		LastSeenAt:      p.LastSeenAt,
	}
	ms := p.MediaStatus
	info.MediaStatus = &ms
	return info
}
