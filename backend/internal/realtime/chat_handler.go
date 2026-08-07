package realtime

import (
	"log"
	"time"

	"github.com/google/uuid"

	"backend/internal/realtime/events"
)

// targetRoomAllowed reports whether a message addressed to envRoom may be
// processed on a connection bound to connRoom. An empty envRoom is allowed —
// the FE can send room_id="" before its store hydrates, in which case the
// connection's own room is used.
func targetRoomAllowed(connRoom, envRoom string) bool {
	return envRoom == "" || envRoom == connRoom
}

// handleChatSend handles the "chat:send" event from clients.
func (r *MessageRouter) handleChatSend(conn *ClientConnection, env *events.Envelope) {
	// Room check runs unconditionally — it used to be disabled whenever
	// LIVEKIT_API_SECRET was the dev fallback, which leaked chat across rooms
	// on any default-config deployment.
	if !targetRoomAllowed(conn.RoomID, env.RoomID) {
		log.Printf("[chat] send rejected: room ID mismatch client=%s message=%s", conn.RoomID, env.RoomID)
		r.sendError(conn, env.RequestID, "FORBIDDEN", "Không có quyền truy cập phòng này")
		return
	}

	// 1. Get room
	room := r.roomManager.Get(conn.RoomID)
	if room == nil {
		r.sendError(conn, env.RequestID, "NOT_FOUND", "Không tìm thấy phòng")
		return
	}

	// 2. Parse payload
	var payload events.ChatSendPayload
	if err := env.ParsePayload(&payload); err != nil {
		log.Printf("[router] failed to parse chat:send payload: %v", err)
		r.sendError(conn, env.RequestID, "VALIDATION_ERROR", "Dữ liệu payload không hợp lệ")
		return
	}

	// Default visibility to room if not specified
	visibility := payload.Visibility
	if visibility == "" {
		visibility = events.VisibilityRoom
	}

	// 3. Build chat message payload
	msgID := payload.MessageID
	if msgID == "" {
		msgID = uuid.New().String()
	}
	now := time.Now().UTC()
	msgPayload := events.ChatMessagePayload{
		MessageID:           msgID,
		SenderParticipantID: conn.ID,
		SenderName:          conn.DisplayName,
		SenderType:          events.ParticipantType(conn.Role),
		Message:             payload.Message,
		Visibility:          visibility,
		CreatedAt:           now,
	}

	// 4. Keep an in-memory copy for live reconnect history, then persist to
	// Postgres via the batch saver (source='chat').
	r.roomManager.SaveChatMessage(room.ID, msgPayload)

	if r.transcriptSaver != nil {
		record := TranscriptRecord{
			ID:            msgPayload.MessageID,
			InterviewID:   room.InterviewID,
			ParticipantID: msgPayload.SenderParticipantID,
			SpeakerType:   string(msgPayload.SenderType),
			SpeakerName:   msgPayload.SenderName,
			Content:       msgPayload.Message,
			Language:      "vi",
			Confidence:    1.0,
			Source:        "chat",
			IsFinal:       true,
			CreatedAt:     now,
		}
		r.transcriptSaver.Push(record)
	}

	// 5. Broadcast to participants
	chatEnv, err := events.NewEnvelope(events.EventChatMessage, env.RequestID, room.ID, room.InterviewID, msgPayload)
	if err != nil {
		log.Printf("[chat] failed to create envelope: %v", err)
		return
	}
	raw, err := chatEnv.ToJSON()
	if err != nil {
		log.Printf("[chat] failed to marshal envelope: %v", err)
		return
	}

	// Filter and broadcast based on visibility payload
	if msgPayload.Visibility == events.VisibilityRecruiterOnly {
		room.BroadcastRecruitersOnly(raw)
	} else if devBroadcastAllRooms() {
		// Dev-only multi-tab demo mirror; requires explicit
		// REALTIME_DEV_BROADCAST_ALL=true and never activates in production.
		r.roomManager.BroadcastToAllRooms(raw)
	} else {
		room.BroadcastAll(raw)
	}
}
