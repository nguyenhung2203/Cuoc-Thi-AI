package realtime

import (
	"log"
	"time"

	"golang.org/x/time/rate"

	"backend/internal/realtime/events"
)

// sendError helper sends a structured error payload back to the client.
func (r *MessageRouter) sendError(conn *ClientConnection, requestID, code, message string) {
	env, err := events.NewEnvelope(events.EventError, requestID, conn.RoomID, conn.InterviewID,
		events.ErrorPayload{Code: code, Message: message, Recoverable: false})
	if err != nil {
		return
	}
	raw, _ := env.ToJSON()
	if !conn.TrySend(raw) {
		log.Printf("[ws] send buffer full/closed for connID=%s, dropping error event", conn.ID)
	}
}

// handleRoomJoin handles the "room:join" event from clients.
func (r *MessageRouter) handleRoomJoin(conn *ClientConnection, env *events.Envelope) {
	// Block requests targetting a different room than authorized
	if env.RoomID != conn.RoomID {
		log.Printf("[room] join rejected: room ID mismatch client=%s message=%s", conn.RoomID, env.RoomID)
		r.sendError(conn, env.RequestID, "FORBIDDEN", "Không có quyền truy cập phòng này")
		return
	}

	var payload events.RoomJoinPayload
	if err := env.ParsePayload(&payload); err != nil {
		log.Printf("[router] failed to parse room:join payload: %v", err)
		r.sendError(conn, env.RequestID, "VALIDATION_ERROR", "Dữ liệu payload không hợp lệ")
		return
	}

	// 1. Get or create room
	room, created := r.roomManager.GetOrCreate(conn.RoomID, conn.InterviewID)
	if created {
		r.presenceManager.StartWatcher(room.ID)
	}
	r.roomManager.CancelEmptyRoomTimer(room.ID)

	// 2. Determine display name: prioritize token claims, fallback to payload
	displayName := conn.DisplayName
	if displayName == "" {
		displayName = payload.DisplayName
	}
	if displayName == "" {
		displayName = "User"
	}

	if conn.DisplayName == "" {
		conn.DisplayName = displayName
	}

	// 3. Create or restore participant
	oldP := room.FindParticipantByUserID(conn.UserID)
	isReconnect := false
	var missedEvents [][]byte
	var oldLastSeen time.Time
	var p *Participant

	if oldP != nil {
		isReconnect = true
		oldConnID := oldP.ConnectionID
		oldLastSeen = oldP.LastSeenAt

		if oldP.Connection != nil && oldP.ConnectionID != conn.ID {
			oldP.Connection.Close()
		}
		room.RemoveParticipant(oldConnID)

		oldP.ConnectionID = conn.ID
		oldP.Connection = conn
		oldP.ConnectionState = events.ConnectionOnline
		oldP.LastSeenAt = time.Now().UTC()
		if displayName != "" && displayName != "User" {
			oldP.DisplayName = displayName
		}
		p = oldP
		room.AddParticipant(p)
		room.UpdateParticipantIDInTracks(oldConnID, conn.ID)

		missedEvents = room.GetMissedEvents(oldLastSeen, p.ParticipantType)
		log.Printf("[room] participant=%s reconnected to room=%s (recovered %d missed events)", conn.ID, conn.RoomID, len(missedEvents))

		if r.auditLogger != nil {
			r.auditLogger.LogEvent("reconnect", conn.UserID, conn.Role, "interview_room", conn.RoomID, "", conn.IPAddress, map[string]interface{}{"room_id": conn.RoomID, "connection_id": conn.ID})
		}
	} else {
		p = &Participant{
			ConnectionID:     conn.ID,
			UserID:           conn.UserID,
			ParticipantType:  events.ParticipantType(conn.Role),
			DisplayName:      displayName,
			ConnectionState:  events.ConnectionOnline,
			JoinedAt:         time.Now().UTC(),
			LastSeenAt:       time.Now().UTC(),
			Connection:       conn,
			MediaRateLimiter: rate.NewLimiter(rate.Every(2*time.Second), 10),
		}
		room.AddParticipant(p)
		log.Printf("[room] participant=%s joined room=%s as %s", conn.ID, conn.RoomID, p.ParticipantType)

		if r.auditLogger != nil {
			r.auditLogger.LogEvent("room_join", conn.UserID, conn.Role, "interview_room", conn.RoomID, "", conn.IPAddress, map[string]interface{}{"room_id": conn.RoomID, "connection_id": conn.ID})
		}
	}

	log.Printf("[db] INSERT INTO interview_participants (id, interview_id, user_id, participant_type, display_name, joined_at, connection_state) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', 'online')",
		conn.ID, conn.InterviewID, conn.UserID, p.ParticipantType, p.DisplayName, p.JoinedAt.Format(time.RFC3339))
	log.Printf("[redis] SET presence:%s:%s value=online EX 90", conn.RoomID, conn.ID)

	// 4. Send ACK room:joined to the client who joined
	ackPayload := events.RoomJoinedPayload{
		ParticipantID:   conn.ID,
		RoomStatus:      room.Status,
		InterviewStatus: string(room.Status),
		Participants:    room.ParticipantList(),
		StartedAt:       room.StartedAt,
		EndedAt:         room.EndedAt,
	}
	if isReconnect {
		ackPayload.MediaStatus = &p.MediaStatus
		ackPayload.MissedEventsCount = len(missedEvents)
		ackPayload.SyncFromTimestamp = oldLastSeen
	}

	ackEnv, err := events.NewEnvelope(events.EventRoomJoined, env.RequestID, room.ID, room.InterviewID, ackPayload)
	if err == nil {
		rawAck, _ := ackEnv.ToJSON()
		if !conn.TrySend(rawAck) {
			log.Printf("[room] send buffer full/closed for participant=%s ACK", conn.ID)
		}
	}

	// Stream historical chat messages to joining client so chat history persists across browser refresh (F5)
	for _, chatMsg := range r.roomManager.GetChatHistory(room.ID) {
		if chatEnv, err := events.NewEnvelope(events.EventChatMessage, "", room.ID, room.InterviewID, chatMsg); err == nil {
			if rawChat, err := chatEnv.ToJSON(); err == nil {
				conn.TrySend(rawChat)
			}
		}
	}

	if isReconnect {
		for _, rawMissed := range missedEvents {
			conn.TrySend(rawMissed)
		}
	} else {
		joinPayload := events.RoomUserJoinedPayload{
			ParticipantID:   conn.ID,
			DisplayName:     p.DisplayName,
			ParticipantType: p.ParticipantType,
			JoinedAt:        p.JoinedAt,
		}
		if joinEnv, err := events.NewEnvelope(events.EventRoomUserJoined, env.RequestID, room.ID, room.InterviewID, joinPayload); err == nil {
			rawJoin, _ := joinEnv.ToJSON()
			room.BroadcastExcept(conn.ID, rawJoin)
		}
	}

	// 5. Broadcast room:presence_update to all participants in the room
	presencePayload := events.RoomPresenceUpdatePayload{
		Participants: room.ParticipantList(),
	}
	presenceEnv, err := events.NewEnvelope(events.EventRoomPresenceUpdate, env.RequestID, room.ID, room.InterviewID, presencePayload)
	if err == nil {
		rawPresence, _ := presenceEnv.ToJSON()
		room.BroadcastAll(rawPresence)
	}
}

// handleRoomLeave handles the "room:leave" event from clients.
func (r *MessageRouter) handleRoomLeave(conn *ClientConnection, env *events.Envelope) {
	// Block requests targetting a different room than authorized
	if env.RoomID != conn.RoomID {
		log.Printf("[room] leave rejected: room ID mismatch client=%s message=%s", conn.RoomID, env.RoomID)
		r.sendError(conn, env.RequestID, "FORBIDDEN", "Không có quyền truy cập phòng này")
		return
	}

	room := r.roomManager.Get(conn.RoomID)
	if room == nil {
		r.sendError(conn, env.RequestID, "NOT_FOUND", "Không tìm thấy phòng")
		return
	}

	// 1. Remove participant from the room
	room.RemoveParticipant(conn.ID)
	log.Printf("[room] participant=%s left room=%s", conn.ID, conn.RoomID)

	if r.auditLogger != nil {
		r.auditLogger.LogEvent("room_leave", conn.UserID, conn.Role, "interview_room", conn.RoomID, "", conn.IPAddress, map[string]interface{}{"room_id": conn.RoomID, "connection_id": conn.ID})
	}

	// Simulate database update: left_at in interview_participants table
	log.Printf("[db] UPDATE interview_participants SET left_at = '%s', connection_state = 'left' WHERE id = '%s'",
		time.Now().UTC().Format(time.RFC3339), conn.ID)

	// Simulate Redis presence delete
	log.Printf("[redis] DEL presence:%s:%s", conn.RoomID, conn.ID)

	// 2. Broadcast room:user_left to everyone remaining in the room
	leftPayload := events.RoomUserLeftPayload{
		ParticipantID: conn.ID,
		Reason:        "user_clicked_leave",
		LeftAt:        time.Now().UTC(),
	}
	leftEnv, err := events.NewEnvelope(events.EventRoomUserLeft, env.RequestID, room.ID, room.InterviewID, leftPayload)
	if err == nil {
		rawLeft, _ := leftEnv.ToJSON()
		room.BroadcastAll(rawLeft)
	}

	// 3. Broadcast room:presence_update to remaining participants
	presencePayload := events.RoomPresenceUpdatePayload{
		Participants: room.ParticipantList(),
	}
	presenceEnv, err := events.NewEnvelope(events.EventRoomPresenceUpdate, env.RequestID, room.ID, room.InterviewID, presencePayload)
	if err == nil {
		rawPresence, _ := presenceEnv.ToJSON()
		room.BroadcastAll(rawPresence)
	}

	// 4. Start empty room grace timer if no online participants remain
	if !room.HasOnlineParticipants() {
		r.roomManager.StartEmptyRoomTimer(room.ID, r.presenceManager)
	}
}

// handleRoomHeartbeat handles the "room:heartbeat" event from clients.
func (r *MessageRouter) handleRoomHeartbeat(conn *ClientConnection, env *events.Envelope) {
	// Block requests targetting a different room than authorized
	if env.RoomID != conn.RoomID {
		return
	}
	r.presenceManager.OnHeartbeat(conn.RoomID, conn.ID)
}
