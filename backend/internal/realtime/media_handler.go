package realtime

import (
	"log"

	"backend/internal/realtime/events"
)

// handleMediaStatus handles the "media:status" event from clients.
func (r *MessageRouter) handleMediaStatus(conn *ClientConnection, env *events.Envelope) {
	// Block requests targetting a different room than authorized
	if env.RoomID != conn.RoomID {
		log.Printf("[media] status rejected: room ID mismatch client=%s message=%s", conn.RoomID, env.RoomID)
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
	var payload events.MediaStatusPayload
	if err := env.ParsePayload(&payload); err != nil {
		log.Printf("[router] failed to parse media:status payload: %v", err)
		r.sendError(conn, env.RequestID, "VALIDATION_ERROR", "Dữ liệu payload không hợp lệ")
		return
	}

	// 3. Get participant and check rate limit
	p := room.GetParticipant(conn.ID)
	if p == nil {
		r.sendError(conn, env.RequestID, "NOT_FOUND", "Không tìm thấy người tham gia")
		return
	}

	if p.MediaRateLimiter != nil && !p.MediaRateLimiter.Allow() {
		log.Printf("[media] rate limited for participant=%s", conn.ID)
		r.sendError(conn, env.RequestID, "RATE_LIMITED", "Vui lòng thao tác chậm lại")
		return
	}

	// 4. Update participant state in room
	statusInfo := events.MediaStatusInfo{
		MicEnabled:    payload.MicEnabled,
		CameraEnabled: payload.CameraEnabled,
		ScreenSharing: payload.ScreenSharing,
	}
	room.UpdateMediaStatus(conn.ID, statusInfo)

	// Simulate database update for media_status (as per task description: Cập nhật media_status JSONB trong interview_participants)
	log.Printf("[db] UPDATE interview_participants SET media_status = '{\"mic_enabled\": %v, \"camera_enabled\": %v, \"screen_sharing\": %v}' WHERE id = '%s'",
		payload.MicEnabled, payload.CameraEnabled, payload.ScreenSharing, conn.ID)

	// 5. Broadcast media:status_changed to everyone in the room
	changedPayload := events.MediaStatusChangedPayload{
		ParticipantID: conn.ID,
		MicEnabled:    payload.MicEnabled,
		CameraEnabled: payload.CameraEnabled,
		ScreenSharing: payload.ScreenSharing,
	}

	changedEnv, err := events.NewEnvelope(events.EventMediaStatusChanged, env.RequestID, room.ID, room.InterviewID, changedPayload)
	if err != nil {
		log.Printf("[media] failed to create envelope: %v", err)
		return
	}

	raw, err := changedEnv.ToJSON()
	if err != nil {
		log.Printf("[media] failed to marshal envelope: %v", err)
		return
	}

	room.BroadcastAll(raw)
}
