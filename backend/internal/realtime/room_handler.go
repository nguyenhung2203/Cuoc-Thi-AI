package realtime

import (
	"log"
	"time"

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
	select {
	case conn.Send <- raw:
	default:
		log.Printf("[ws] send buffer full for connID=%s, dropping error event", conn.ID)
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

	// 2. Determine display name: prioritize token claims, fallback to payload
	displayName := conn.DisplayName
	if displayName == "" {
		displayName = payload.DisplayName
	}
	if displayName == "" {
		displayName = "User"
	}

	// Update connection display name if it was empty
	if conn.DisplayName == "" {
		conn.DisplayName = displayName
	}

	// 3. Create participant
	p := &Participant{
		ConnectionID:    conn.ID,
		UserID:          conn.UserID,
		ParticipantType: events.ParticipantType(conn.Role),
		DisplayName:     displayName,
		ConnectionState: events.ConnectionOnline,
		JoinedAt:        time.Now().UTC(),
		LastSeenAt:      time.Now().UTC(),
		Connection:      conn,
	}

	// 4. Add to room
	room.AddParticipant(p)
	log.Printf("[room] participant=%s joined room=%s as %s", conn.ID, conn.RoomID, p.ParticipantType)

	// Simulate database save: interview_participants table
	log.Printf("[db] INSERT INTO interview_participants (id, interview_id, user_id, participant_type, display_name, joined_at, connection_state) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', 'online')",
		conn.ID, conn.InterviewID, conn.UserID, p.ParticipantType, p.DisplayName, p.JoinedAt.Format(time.RFC3339))

	// Simulate Redis presence write
	log.Printf("[redis] SET presence:%s:%s value=online EX 90", conn.RoomID, conn.ID)

	// 5. Send ACK room:joined to the client who joined
	ackPayload := events.RoomJoinedPayload{
		ParticipantID:   conn.ID,
		RoomStatus:      room.Status,
		InterviewStatus: string(room.Status),
		Participants:    room.ParticipantList(),
	}
	ackEnv, err := events.NewEnvelope(events.EventRoomJoined, env.RequestID, room.ID, room.InterviewID, ackPayload)
	if err == nil {
		rawAck, _ := ackEnv.ToJSON()
		select {
		case conn.Send <- rawAck:
		default:
			log.Printf("[room] send buffer full for participant=%s ACK", conn.ID)
		}
	}

	// 6. Broadcast room:user_joined to everyone else in the room
	joinPayload := events.RoomUserJoinedPayload{
		ParticipantID:   conn.ID,
		DisplayName:     p.DisplayName,
		ParticipantType: p.ParticipantType,
		JoinedAt:        p.JoinedAt,
	}
	joinEnv, err := events.NewEnvelope(events.EventRoomUserJoined, env.RequestID, room.ID, room.InterviewID, joinPayload)
	if err == nil {
		rawJoin, _ := joinEnv.ToJSON()
		room.BroadcastExcept(conn.ID, rawJoin)
	}

	// 7. Broadcast room:presence_update to all participants in the room
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

	// 4. Clean up room if empty
	room.mu.RLock()
	participantsCount := len(room.Participants)
	room.mu.RUnlock()
	if participantsCount == 0 {
		r.roomManager.Delete(room.ID)
		r.presenceManager.StopWatcher(room.ID)
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
