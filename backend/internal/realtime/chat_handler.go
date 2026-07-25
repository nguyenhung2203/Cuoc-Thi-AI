package realtime

import (
	"log"
	"os"
	"time"

	"github.com/google/uuid"

	"backend/internal/realtime/events"
)

// handleChatSend handles the "chat:send" event from clients.
func (r *MessageRouter) handleChatSend(conn *ClientConnection, env *events.Envelope) {
	// Block requests targetting a different room than authorized (unless mock dev env)
	if env.RoomID != conn.RoomID && os.Getenv("LIVEKIT_API_SECRET") != "" && os.Getenv("LIVEKIT_API_SECRET") != "devsecret" {
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
	} else {
		// In mock dev mode, broadcast across all rooms so candidate and recruiter tabs hear each other
		if os.Getenv("LIVEKIT_API_SECRET") == "" || os.Getenv("LIVEKIT_API_SECRET") == "devsecret" {
			r.roomManager.BroadcastToAllRooms(raw)
		} else {
			room.BroadcastAll(raw)
		}
	}
}
