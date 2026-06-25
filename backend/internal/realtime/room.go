package realtime

import (
	"log"
	"sync"
	"time"

	"backend/internal/realtime/events"
)

// Room represents a single interview room in memory.
// It holds all participants and the current room status.
type Room struct {
	ID                            string
	InterviewID                   string
	Status                        events.RoomStatus
	TranscriptEnabledForCandidate bool
	IsMockInterview               bool
	Participants                  map[string]*Participant // participantID (connectionID) → Participant
	ProcessedRequests             map[string]bool         // requestID → processed (for idempotency)
	CreatedAt                     time.Time
	mu                            sync.RWMutex
}

// newRoom initialises an empty room in "waiting" state.
func newRoom(roomID, interviewID string) *Room {
	return &Room{
		ID:                roomID,
		InterviewID:       interviewID,
		Status:            events.RoomStatusWaiting,
		Participants:      make(map[string]*Participant),
		ProcessedRequests: make(map[string]bool),
		CreatedAt:         time.Now().UTC(),
	}
}

// RecordRequest checks if a request ID has already been processed.
// If it has, it returns true. Otherwise, it records the request ID and returns false.
func (r *Room) RecordRequest(requestID string) bool {
	if requestID == "" {
		return false
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.ProcessedRequests == nil {
		r.ProcessedRequests = make(map[string]bool)
	}
	if r.ProcessedRequests[requestID] {
		return true
	}
	r.ProcessedRequests[requestID] = true
	return false
}

// AddParticipant adds or replaces a participant in the room.
func (r *Room) AddParticipant(p *Participant) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Participants[p.ConnectionID] = p
}

// RemoveParticipant removes a participant by connectionID.
func (r *Room) RemoveParticipant(connID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.Participants, connID)
}

// GetParticipant returns a participant by connectionID (nil if not found).
func (r *Room) GetParticipant(connID string) *Participant {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.Participants[connID]
}

// ParticipantList returns a snapshot of all participants as ParticipantInfo slices.
func (r *Room) ParticipantList() []events.ParticipantInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]events.ParticipantInfo, 0, len(r.Participants))
	for _, p := range r.Participants {
		out = append(out, p.ToInfo())
	}
	return out
}

// ── Broadcast helpers ─────────────────────────────────────────────────────────

// BroadcastAll sends an envelope to every participant in the room.
func (r *Room) BroadcastAll(raw []byte) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.Participants {
		r.send(p, raw)
	}
}

// BroadcastExcept sends to everyone except the given connectionID.
func (r *Room) BroadcastExcept(excludeConnID string, raw []byte) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.Participants {
		if p.ConnectionID == excludeConnID {
			continue
		}
		r.send(p, raw)
	}
}

// BroadcastRecruitersOnly sends only to participants with role "recruiter".
func (r *Room) BroadcastRecruitersOnly(raw []byte) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.Participants {
		if p.ParticipantType == events.ParticipantRecruiter {
			r.send(p, raw)
		}
	}
}

// SendTo sends an envelope to a single participant by connectionID.
func (r *Room) SendTo(connID string, raw []byte) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if p, ok := r.Participants[connID]; ok {
		r.send(p, raw)
	}
}

// BroadcastWithVisibility uses the centralised visibility engine to decide
// which participants receive a given event.
func (r *Room) BroadcastWithVisibility(event string, senderConnID string, raw []byte) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.Participants {
		isSender := p.ConnectionID == senderConnID
		if events.ShouldSendToParticipant(
			event,
			p.ParticipantType,
			isSender,
			r.TranscriptEnabledForCandidate,
			r.IsMockInterview,
		) {
			r.send(p, raw)
		}
	}
}

// send enqueues raw bytes on the participant's outbound channel.
// It is safe to call with the room lock held because it never blocks.
func (r *Room) send(p *Participant, raw []byte) {
	if p.Connection == nil {
		return
	}
	select {
	case p.Connection.Send <- raw:
	default:
		log.Printf("[room] send buffer full for participant=%s", p.ConnectionID)
	}
}
