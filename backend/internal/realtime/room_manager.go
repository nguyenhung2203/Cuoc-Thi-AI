package realtime

import (
	"log"
	"os"
	"sync"
	"time"

	"backend/internal/realtime/events"
)

// RoomManager maintains a registry of all active rooms.
// It is safe for concurrent use.
type RoomManager struct {
	rooms             map[string]*Room                       // roomID → Room
	simulatedStatuses map[string]events.RoomStatus           // roomID → RoomStatus (simulates DB/Redis)
	simulatedChat     map[string][]events.ChatMessagePayload // roomID → chat history
	emptyTimers       map[string]*time.Timer                 // roomID → grace timer
	timerMu           sync.Mutex
	mu                sync.RWMutex
}

// NewRoomManager creates a ready-to-use RoomManager.
func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms:             make(map[string]*Room),
		simulatedStatuses: make(map[string]events.RoomStatus),
		simulatedChat:     make(map[string][]events.ChatMessagePayload),
		emptyTimers:       make(map[string]*time.Timer),
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
	rm.CancelEmptyRoomTimer(roomID)
	rm.mu.Lock()
	defer rm.mu.Unlock()
	delete(rm.rooms, roomID)
	log.Printf("[room-mgr] deleted room=%s", roomID)
}

// StartEmptyRoomTimer starts a grace-period timer for an empty room.
// When gracePeriodRoom expires, the room is cleaned up from memory and DB.
func (rm *RoomManager) StartEmptyRoomTimer(roomID string, pm *PresenceManager) {
	rm.timerMu.Lock()
	defer rm.timerMu.Unlock()
	if rm.emptyTimers == nil {
		rm.emptyTimers = make(map[string]*time.Timer)
	}
	if _, exists := rm.emptyTimers[roomID]; exists {
		return // timer already running
	}
	log.Printf("[room-mgr] all participants disconnected in room=%s, starting room grace timer (%v)", roomID, gracePeriodRoom)

	rm.emptyTimers[roomID] = time.AfterFunc(gracePeriodRoom, func() {
		rm.CleanupExpiredRoom(roomID, pm)
	})
}

// CancelEmptyRoomTimer cancels any pending empty room grace timer.
func (rm *RoomManager) CancelEmptyRoomTimer(roomID string) {
	rm.timerMu.Lock()
	defer rm.timerMu.Unlock()
	if rm.emptyTimers == nil {
		return
	}
	if t, exists := rm.emptyTimers[roomID]; exists {
		t.Stop()
		delete(rm.emptyTimers, roomID)
		log.Printf("[room-mgr] cancelled room grace timer for room=%s", roomID)
	}
}

// CleanupExpiredRoom performs cleanup after grace period expires.
func (rm *RoomManager) CleanupExpiredRoom(roomID string, pm *PresenceManager) {
	rm.timerMu.Lock()
	delete(rm.emptyTimers, roomID)
	rm.timerMu.Unlock()

	room := rm.Get(roomID)
	if room == nil {
		return
	}
	nowStr := time.Now().UTC().Format(time.RFC3339)
	log.Printf("[room-mgr] empty room grace period expired for room=%s, triggering cleanup", roomID)

	// Cập nhật DB room status
	log.Printf("[db] UPDATE interview_rooms SET status = 'closed', updated_at = '%s' WHERE id = '%s'", nowStr, roomID)

	// Cập nhật participant left_at
	log.Printf("[db] UPDATE interview_participants SET left_at = '%s', connection_state = 'offline' WHERE interview_id = '%s' AND left_at IS NULL", nowStr, room.InterviewID)

	// Release Redis keys
	log.Printf("[redis] DEL room_status:%s", roomID)
	log.Printf("[redis] DEL presence:%s:*", roomID)

	// Remove room khỏi memory
	rm.Delete(roomID)
	if pm != nil {
		pm.StopWatcher(roomID)
	}
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
	// In dev mode, replicate chat history across all active rooms so reconnecting Candidate or Recruiter tab sees everything
	if os.Getenv("LIVEKIT_API_SECRET") == "" || os.Getenv("LIVEKIT_API_SECRET") == "devsecret" {
		for id := range rm.rooms {
			if id != roomID {
				rm.simulatedChat[id] = append(rm.simulatedChat[id], msg)
			}
		}
	}
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

// BroadcastToAllRooms broadcasts raw JSON to every active room in the registry (used in mock dev environment).
func (rm *RoomManager) BroadcastToAllRooms(data []byte) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	for _, r := range rm.rooms {
		r.BroadcastAll(data)
	}
}
