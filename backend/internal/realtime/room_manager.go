package realtime

import (
	"log"
	"sync"
	"time"

	"backend/internal/realtime/events"
)

// RoomManager maintains a registry of all active rooms.
// It is safe for concurrent use.
type RoomManager struct {
	rooms             map[string]*Room // roomID → Room
	simulatedStatuses map[string]events.RoomStatus // roomID → RoomStatus (simulates DB/Redis)
	simulatedChat     map[string][]events.ChatMessagePayload // roomID → chat history
	mu                sync.RWMutex
}

// NewRoomManager creates a ready-to-use RoomManager.
func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms:             make(map[string]*Room),
		simulatedStatuses: make(map[string]events.RoomStatus),
		simulatedChat:     make(map[string][]events.ChatMessagePayload),
	}
}

// GetOrCreate returns the room for roomID, creating it if it doesn't exist.
// If the room status is cached in Redis/DB, it restores it. Otherwise, it defaults to "waiting".
func (rm *RoomManager) GetOrCreate(roomID, interviewID string) (*Room, bool) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if r, ok := rm.rooms[roomID]; ok {
		return r, false
	}
	r := newRoom(roomID, interviewID)

	// Simulate checking Redis / DB
	log.Printf("[redis] GET room_status:%s", roomID)
	log.Printf("[db] SELECT status FROM interview_rooms WHERE id = '%s'", roomID)

	if status, ok := rm.simulatedStatuses[roomID]; ok {
		r.Status = status
		log.Printf("[room-mgr] restored room=%s status=%s from simulated DB/Redis", roomID, status)
	} else {
		r.Status = events.RoomStatusWaiting
		rm.simulatedStatuses[roomID] = events.RoomStatusWaiting
		log.Printf("[room-mgr] created room=%s interview=%s defaulting to waiting", roomID, interviewID)
	}

	rm.rooms[roomID] = r
	return r, true
}

// Get returns the room for roomID (nil if not found).
func (rm *RoomManager) Get(roomID string) *Room {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.rooms[roomID]
}

// Delete removes a room from the registry and logs it.
func (rm *RoomManager) Delete(roomID string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	delete(rm.rooms, roomID)
	log.Printf("[room-mgr] deleted room=%s", roomID)
}

// Count returns the number of active rooms.
func (rm *RoomManager) Count() int {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return len(rm.rooms)
}

// OnParticipantDisconnect is called when a WebSocket connection drops.
// It marks the participant as "reconnecting" and starts a grace-period timer.
// If the participant doesn't reconnect within gracePeriodShort, they are
// marked "offline" and room:user_left is broadcast.
// Full reconnect/recovery logic is implemented in H-S5-01.
func (rm *RoomManager) OnParticipantDisconnect(roomID, connID string, cm *ConnectionManager) {
	room := rm.Get(roomID)
	if room == nil {
		return
	}

	p := room.GetParticipant(connID)
	if p == nil {
		return
	}

	// Mark as reconnecting immediately
	p.ConnectionState = events.ConnectionReconnecting
	p.Connection = nil
	log.Printf("[room-mgr] participant=%s marked reconnecting in room=%s", connID, roomID)

	// Start grace-period goroutine
	go func() {
		timer := time.NewTimer(gracePeriodShort)
		defer timer.Stop()
		<-timer.C

		// Re-check: if still reconnecting after grace period → offline
		currentP := room.GetParticipant(connID)
		if currentP == nil || currentP.ConnectionState != events.ConnectionReconnecting {
			return // already reconnected or removed
		}

		currentP.ConnectionState = events.ConnectionOffline
		log.Printf("[room-mgr] participant=%s marked offline in room=%s", connID, roomID)

		// Broadcast room:user_left
		env, err := events.NewEnvelope(
			events.EventRoomUserLeft, "", roomID, room.InterviewID,
			events.RoomUserLeftPayload{
				ParticipantID: connID,
				Reason:        "connection_timeout",
				LeftAt:        time.Now().UTC(),
			},
		)
		if err != nil {
			return
		}
		raw, _ := env.ToJSON()
		room.BroadcastExcept(connID, raw)
	}()
}

// SetRoomStatus updates the room's status. Returns false if the transition is invalid.
func (rm *RoomManager) SetRoomStatus(roomID string, newStatus events.RoomStatus) bool {
	room := rm.Get(roomID)
	if room == nil {
		return false
	}
	room.mu.Lock()
	defer room.mu.Unlock()
	if !isValidTransition(room.Status, newStatus) {
		log.Printf("[room-mgr] invalid transition room=%s %s→%s", roomID, room.Status, newStatus)
		return false
	}
	room.Status = newStatus

	// Sync with simulated DB/Redis
	rm.mu.Lock()
	if rm.simulatedStatuses == nil {
		rm.simulatedStatuses = make(map[string]events.RoomStatus)
	}
	rm.simulatedStatuses[roomID] = newStatus
	rm.mu.Unlock()

	return true
}

// isValidTransition enforces the room state machine defined in REALTIME_EVENTS.md §3.
func isValidTransition(from, to events.RoomStatus) bool {
	allowed := map[events.RoomStatus][]events.RoomStatus{
		events.RoomStatusScheduled: {events.RoomStatusWaiting, events.RoomStatusCancelled, events.RoomStatusActive},
		events.RoomStatusWaiting:   {events.RoomStatusActive, events.RoomStatusExpired, events.RoomStatusCancelled},
		events.RoomStatusActive:    {events.RoomStatusPaused, events.RoomStatusCompleted},
		events.RoomStatusPaused:    {events.RoomStatusActive, events.RoomStatusCompleted},
	}
	for _, valid := range allowed[from] {
		if valid == to {
			return true
		}
	}
	return false
}

// SaveChatMessage stores a chat message in the simulated chat database map.
func (rm *RoomManager) SaveChatMessage(roomID string, msg events.ChatMessagePayload) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if rm.simulatedChat == nil {
		rm.simulatedChat = make(map[string][]events.ChatMessagePayload)
	}
	rm.simulatedChat[roomID] = append(rm.simulatedChat[roomID], msg)
}

// GetChatHistory retrieves historical chat messages for a room.
func (rm *RoomManager) GetChatHistory(roomID string) []events.ChatMessagePayload {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	if rm.simulatedChat == nil {
		return nil
	}
	// Return a copy to avoid race conditions when clients read it
	history := rm.simulatedChat[roomID]
	out := make([]events.ChatMessagePayload, len(history))
	copy(out, history)
	return out
}
