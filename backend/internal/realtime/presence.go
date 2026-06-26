package realtime

import (
	"log"
	"sync"
	"time"

	"backend/internal/realtime/events"
)

// PresenceManager handles participant heartbeats and presence state transitions.
// It relies on per-participant timers so no Redis is required in-memory.
// Redis-backed presence (for multi-instance deployments) is deferred to a later sprint.
type PresenceManager struct {
	roomManager *RoomManager
	watchers    map[string]chan struct{}
	mu          sync.Mutex
}

// NewPresenceManager creates a PresenceManager backed by the given RoomManager.
func NewPresenceManager(rm *RoomManager) *PresenceManager {
	return &PresenceManager{
		roomManager: rm,
		watchers:    make(map[string]chan struct{}),
	}
}

// OnHeartbeat is called when a client sends room:heartbeat.
// It refreshes last_seen_at and resets the connection state to "online".
func (pm *PresenceManager) OnHeartbeat(roomID, connID string) {
	room := pm.roomManager.Get(roomID)
	if room == nil {
		return
	}
	p := room.GetParticipant(connID)
	if p == nil {
		return
	}

	p.LastSeenAt = time.Now().UTC()

	// Simulate Redis presence TTL refresh
	log.Printf("[redis] SET presence:%s:%s value=online EX 90", roomID, connID)

	stateChanged := false
	if p.ConnectionState != events.ConnectionOnline {
		p.ConnectionState = events.ConnectionOnline
		stateChanged = true
		log.Printf("[presence] participant=%s back online in room=%s", connID, roomID)

		// Simulate database update when participant comes back online
		log.Printf("[db] UPDATE interview_participants SET connection_state = 'online', last_seen_at = '%s' WHERE id = '%s'",
			p.LastSeenAt.Format(time.RFC3339), connID)
	} else {
		// Update DB last_seen_at periodically even if state didn't change
		log.Printf("[db] UPDATE interview_participants SET last_seen_at = '%s' WHERE id = '%s'",
			p.LastSeenAt.Format(time.RFC3339), connID)
	}

	if stateChanged {
		pm.broadcastPresence(room)
	}
}

// StartWatcher starts a background goroutine that periodically checks all
// participants in a room and transitions stale ones to reconnecting/offline.
func (pm *PresenceManager) StartWatcher(roomID string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, ok := pm.watchers[roomID]; ok {
		return // watcher already running
	}

	done := make(chan struct{})
	pm.watchers[roomID] = done

	log.Printf("[presence] started watcher for room=%s", roomID)

	ticker := time.NewTicker(heartbeatReconnectingThreshold / 2)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-done:
				log.Printf("[presence] stopped watcher for room=%s", roomID)
				return
			case <-ticker.C:
				pm.checkRoom(roomID)
			}
		}
	}()
}

// StopWatcher stops the presence watcher for a room.
func (pm *PresenceManager) StopWatcher(roomID string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if done, ok := pm.watchers[roomID]; ok {
		close(done)
		delete(pm.watchers, roomID)
	}
}

// checkRoom inspects all participants and updates their presence state.
func (pm *PresenceManager) checkRoom(roomID string) {
	room := pm.roomManager.Get(roomID)
	if room == nil {
		return
	}

	now := time.Now().UTC()

	// Check for room waiting expiry
	if room.Status == events.RoomStatusWaiting && now.Sub(room.CreatedAt) >= roomWaitingExpiry {
		log.Printf("[presence] room=%s waiting expiry triggered (CreatedAt=%v, Expiry=%v)",
			roomID, room.CreatedAt, roomWaitingExpiry)

		if pm.roomManager.SetRoomStatus(roomID, events.RoomStatusExpired) {
			log.Printf("[db] UPDATE interview_rooms SET status = 'expired' WHERE id = '%s'", roomID)
			log.Printf("[db] UPDATE interviews SET status = 'expired' WHERE id = '%s'", room.InterviewID)
			log.Printf("[redis] SET room_status:%s value=expired", roomID)

			// Broadcast room:expired to all remaining participants
			env, err := events.NewEnvelope(
				events.EventRoomExpired, "", roomID, room.InterviewID,
				events.RoomExpiredPayload{
					Status:    string(events.RoomStatusExpired),
					ExpiredAt: now,
				},
			)
			if err == nil {
				raw, _ := env.ToJSON()
				room.BroadcastAll(raw)
			}

			// Clean up the room and connections after a short grace period
			go func() {
				time.Sleep(2 * time.Second)
				room.mu.RLock()
				conns := make([]*ClientConnection, 0, len(room.Participants))
				for _, p := range room.Participants {
					if p.Connection != nil {
						conns = append(conns, p.Connection)
					}
				}
				room.mu.RUnlock()

				log.Printf("[room] room expired, closing %d connections in room=%s", len(conns), roomID)
				for _, c := range conns {
					c.Close()
				}
				pm.roomManager.Delete(roomID)
				pm.StopWatcher(roomID)
			}()
			return
		}
	}
	changed := false
	var droppedOffline []string

	room.mu.Lock()
	for _, p := range room.Participants {
		if p.ConnectionState == events.ConnectionLeft || p.ConnectionState == events.ConnectionOffline {
			continue
		}
		since := now.Sub(p.LastSeenAt)
		switch {
		case since >= heartbeatOfflineThreshold:
			p.ConnectionState = events.ConnectionOffline
			changed = true
			droppedOffline = append(droppedOffline, p.ConnectionID)
			log.Printf("[presence] participant=%s → offline (no heartbeat %.0fs)", p.ConnectionID, since.Seconds())
			log.Printf("[db] UPDATE interview_participants SET connection_state = 'offline' WHERE id = '%s'", p.ConnectionID)
			log.Printf("[redis] DEL presence:%s:%s", roomID, p.ConnectionID)
		case since >= heartbeatReconnectingThreshold && p.ConnectionState == events.ConnectionOnline:
			p.ConnectionState = events.ConnectionReconnecting
			changed = true
			log.Printf("[presence] participant=%s → reconnecting (no heartbeat %.0fs)", p.ConnectionID, since.Seconds())
			log.Printf("[db] UPDATE interview_participants SET connection_state = 'reconnecting' WHERE id = '%s'", p.ConnectionID)
		}
	}
	room.mu.Unlock()

	for _, cID := range droppedOffline {
		room.RemoveParticipant(cID)
		leftPayload := events.RoomUserLeftPayload{
			ParticipantID: cID,
			Reason:        "reconnect_timeout",
			LeftAt:        now,
		}
		if leftEnv, err := events.NewEnvelope(events.EventRoomUserLeft, "", room.ID, room.InterviewID, leftPayload); err == nil {
			rawLeft, _ := leftEnv.ToJSON()
			room.BroadcastAll(rawLeft)
		}
	}

	if changed {
		pm.broadcastPresence(room)
	}
}

// broadcastPresence sends a room:presence_update to everyone in the room.
func (pm *PresenceManager) broadcastPresence(room *Room) {
	env, err := events.NewEnvelope(
		events.EventRoomPresenceUpdate, "", room.ID, room.InterviewID,
		events.RoomPresenceUpdatePayload{Participants: room.ParticipantList()},
	)
	if err != nil {
		return
	}
	raw, _ := env.ToJSON()
	room.BroadcastAll(raw)
}
