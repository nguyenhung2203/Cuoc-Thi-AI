package realtime

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"backend/internal/realtime/events"
)

// EventRecord stores historical broadcast envelopes for reconnect backlog sync.
type EventRecord struct {
	Event          string
	SenderConnID   string
	Raw            []byte
	Timestamp      time.Time
	RecruitersOnly bool
}

// Room represents a single interview room in memory.
// It holds all participants and the current room status.
type Room struct {
	ID                            string
	InterviewID                   string
	Status                        events.RoomStatus
	StartedAt                     *time.Time
	EndedAt                       *time.Time
	TranscriptEnabledForCandidate bool
	IsMockInterview               bool
	Participants                  map[string]*Participant // participantID (connectionID) → Participant
	TrackParticipantMap           map[string]string       // trackID → participantID
	ProcessedRequests             map[string]bool         // requestID → processed (for idempotency)
	EventHistory                  []EventRecord           // backlog buffer
	CreatedAt                     time.Time
	mu                            sync.RWMutex
}

// newRoom initialises an empty room in "waiting" state.
func newRoom(roomID, interviewID string) *Room {
	return &Room{
		ID:                  roomID,
		InterviewID:         interviewID,
		Status:              events.RoomStatusWaiting,
		Participants:        make(map[string]*Participant),
		TrackParticipantMap: make(map[string]string),
		ProcessedRequests:   make(map[string]bool),
		EventHistory:        make([]EventRecord, 0, 500),
		CreatedAt:           time.Now().UTC(),
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

// FindParticipantByUserID returns a participant matching userID (nil if not found).
func (r *Room) FindParticipantByUserID(userID string) *Participant {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.Participants {
		if p.UserID == userID {
			return p
		}
	}
	return nil
}

// UpdateParticipantLastSeen updates LastSeenAt for a participant by userID.
func (r *Room) UpdateParticipantLastSeen(userID string, ts time.Time) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, p := range r.Participants {
		if p.UserID == userID {
			p.LastSeenAt = ts
		}
	}
}

// UpdateParticipantIDInTracks maps old participant connection ID to new connection ID in audio track map.
func (r *Room) UpdateParticipantIDInTracks(oldID, newID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for tID, pID := range r.TrackParticipantMap {
		if pID == oldID {
			r.TrackParticipantMap[tID] = newID
		}
	}
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

// HasOnlineParticipants returns true if at least one participant has connection_state = online.
func (r *Room) HasOnlineParticipants() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.Participants {
		if p.ConnectionState == events.ConnectionOnline {
			return true
		}
	}
	return false
}

// recordEvent records historical broadcasts in the room buffer.
func (r *Room) recordEvent(event, senderConnID string, raw []byte, recruitersOnly bool) {
	if event == "" || event == events.EventRoomPresenceUpdate ||
		event == events.EventRoomUserJoined || event == events.EventRoomUserLeft ||
		event == events.EventRoomHeartbeat {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.EventHistory == nil {
		r.EventHistory = make([]EventRecord, 0, 500)
	}
	r.EventHistory = append(r.EventHistory, EventRecord{
		Event:          event,
		SenderConnID:   senderConnID,
		Raw:            raw,
		Timestamp:      time.Now().UTC(),
		RecruitersOnly: recruitersOnly,
	})
	if len(r.EventHistory) > 500 {
		r.EventHistory = r.EventHistory[len(r.EventHistory)-500:]
	}
}

// GetMissedEvents retrieves all historical events sent after the given timestamp for the target role.
func (r *Room) GetMissedEvents(since time.Time, role events.ParticipantType) [][]byte {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([][]byte, 0)
	for _, rec := range r.EventHistory {
		if rec.Timestamp.After(since) || rec.Timestamp.Equal(since) {
			if rec.RecruitersOnly && role != events.ParticipantRecruiter {
				continue
			}
			if !events.ShouldSendToParticipant(
				rec.Event,
				role,
				false,
				r.TranscriptEnabledForCandidate,
				r.IsMockInterview,
			) {
				continue
			}
			out = append(out, rec.Raw)
		}
	}
	return out
}

// ── Broadcast helpers ─────────────────────────────────────────────────────────

// BroadcastAll sends an envelope to every participant in the room.
func (r *Room) BroadcastAll(raw []byte) {
	var head struct {
		Event string `json:"event"`
	}
	json.Unmarshal(raw, &head)
	r.recordEvent(head.Event, "", raw, false)

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
	var head struct {
		Event string `json:"event"`
	}
	json.Unmarshal(raw, &head)
	r.recordEvent(head.Event, "", raw, true)

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
	r.recordEvent(event, senderConnID, raw, false)
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
	if p == nil || p.Connection == nil {
		return
	}
	if !p.Connection.TrySend(raw) {
		log.Printf("[room] send buffer full/closed for participant=%s", p.ConnectionID)
	}
}

// UpdateMediaStatus updates the media status of a participant safely.
func (r *Room) UpdateMediaStatus(connID string, status events.MediaStatusInfo) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if p, ok := r.Participants[connID]; ok {
		p.MediaStatus = status
	}
}

// RegisterTrack maps an audio track ID to a participant connection ID.
func (r *Room) RegisterTrack(trackID, participantID string) {
	if trackID == "" || participantID == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.TrackParticipantMap == nil {
		r.TrackParticipantMap = make(map[string]string)
	}
	r.TrackParticipantMap[trackID] = participantID
}

// ResolveSpeaker maps track ID, participant ID, or identity to the correct speaker type and name.
func (r *Room) ResolveSpeaker(participantID, trackID, identity string, fallbackType events.SpeakerType, fallbackName string) (string, events.SpeakerType, string) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 1. Resolve participantID from trackID if needed
	if participantID == "" && trackID != "" && r.TrackParticipantMap != nil {
		participantID = r.TrackParticipantMap[trackID]
	}

	// 2. Lookup participant in room registry
	for _, p := range r.Participants {
		if (participantID != "" && p.ConnectionID == participantID) ||
			(identity != "" && (p.ConnectionID == identity || p.UserID == identity || p.DisplayName == identity)) {
			return p.ConnectionID, events.SpeakerType(p.ParticipantType), p.DisplayName
		}
	}

	// 3. Fallback when participant dropped or unknown track
	if fallbackType != "" && fallbackType != events.SpeakerUnknown {
		if fallbackName == "" {
			fallbackName = string(fallbackType)
		}
		return participantID, fallbackType, fallbackName
	}

	if fallbackName == "" {
		fallbackName = "Unknown Speaker"
	}
	return participantID, events.SpeakerUnknown, fallbackName
}
