package realtime

import (
	"context"
	"log"
	"time"

	"backend/internal/realtime/events"
)

// persistTimeout bounds the synchronous DB writes performed on interview
// lifecycle events; a slow DB must not stall the broadcast for long.
const persistTimeout = 5 * time.Second

// handleInterviewStart handles the "interview:start" event from clients.
func (r *MessageRouter) handleInterviewStart(conn *ClientConnection, env *events.Envelope) {
	// 1. Only recruiter can start
	if conn.Role != string(events.ParticipantRecruiter) {
		r.sendError(conn, env.RequestID, "FORBIDDEN", "Chỉ Recruiter mới có quyền bắt đầu buổi phỏng vấn")
		return
	}

	// 2. Parse payload
	var payload events.InterviewStartPayload
	if err := env.ParsePayload(&payload); err != nil {
		log.Printf("[router] failed to parse interview:start payload: %v", err)
		r.sendError(conn, env.RequestID, "VALIDATION_ERROR", "Dữ liệu payload không hợp lệ")
		return
	}

	// 3. Get room
	room := r.roomManager.Get(conn.RoomID)
	if room == nil {
		r.sendError(conn, env.RequestID, "NOT_FOUND", "Không tìm thấy phòng")
		return
	}

	// 4. Validate current status / prevent double start
	if room.Status == events.RoomStatusActive {
		r.sendError(conn, env.RequestID, "ALREADY_ACTIVE", "Buổi phỏng vấn đang diễn ra")
		return
	}
	if room.Status == events.RoomStatusCompleted {
		r.sendError(conn, env.RequestID, "COMPLETED", "Buổi phỏng vấn đã kết thúc")
		return
	}

	// 5. Update room status
	if !r.roomManager.SetRoomStatus(room.ID, events.RoomStatusActive) {
		r.sendError(conn, env.RequestID, "INVALID_STATE", "Không thể bắt đầu phỏng vấn từ trạng thái hiện tại")
		return
	}

	// Set room configuration from consent flags
	now := time.Now().UTC()
	room.mu.Lock()
	room.TranscriptEnabledForCandidate = payload.ConsentRecording
	room.StartedAt = &now
	room.mu.Unlock()

	// 6. Persist status + started_at. This is the path the real FE flow takes
	// (interview:start over WS, not the REST endpoint), so the write happens
	// here; MarkStarted is idempotent (COALESCE) so a REST start cannot
	// double-stamp. Broadcast continues even if the DB write fails — a DB
	// hiccup must not break a live interview.
	if r.interviewRepo != nil {
		dbCtx, cancel := context.WithTimeout(context.Background(), persistTimeout)
		if _, err := r.interviewRepo.MarkStarted(dbCtx, room.InterviewID, now); err != nil {
			log.Printf("[room] failed to persist interview start interview=%s: %v", room.InterviewID, err)
		}
		if err := r.interviewRepo.MarkRoomOpened(dbCtx, room.InterviewID, now); err != nil {
			log.Printf("[room] failed to persist room open interview=%s: %v", room.InterviewID, err)
		}
		cancel()
	}

	if r.auditLogger != nil {
		r.auditLogger.LogEvent("interview_start", conn.UserID, conn.Role, "interview_room", room.ID, "", conn.IPAddress, map[string]interface{}{"room_id": room.ID, "interview_id": room.InterviewID, "consent_ai": payload.ConsentAI})
	}

	// 7. Simulate AI Orchestrator activation if consent_ai is true
	if payload.ConsentAI {
		log.Printf("[ai-orchestrator] starting AI pipeline for room=%s (interview=%s, consent_recording=%t)",
			room.ID, room.InterviewID, payload.ConsentRecording)
		
		// Bắt đầu capture audio stream để gửi cho AI STT
		r.audioHook.StartHook(room.ID, room.InterviewID)
	}

	// 8. Broadcast interview:started to all participants
	startedPayload := events.InterviewStartedPayload{
		Status:            string(events.RoomStatusActive),
		StartedAt:         now,
		TranscriptEnabled: payload.ConsentRecording,
		AIEnabled:         payload.ConsentAI,
	}
	startedEnv, err := events.NewEnvelope(events.EventInterviewStarted, env.RequestID, room.ID, room.InterviewID, startedPayload)
	if err == nil {
		raw, _ := startedEnv.ToJSON()
		room.BroadcastAll(raw)
	}
}

// handleInterviewEnd handles the "interview:end" event from clients.
func (r *MessageRouter) handleInterviewEnd(conn *ClientConnection, env *events.Envelope) {
	// 1. Only recruiter can end
	if conn.Role != string(events.ParticipantRecruiter) {
		r.sendError(conn, env.RequestID, "FORBIDDEN", "Chỉ Recruiter mới có quyền kết thúc buổi phỏng vấn")
		return
	}

	// 2. Parse payload
	var payload events.InterviewEndPayload
	if err := env.ParsePayload(&payload); err != nil {
		log.Printf("[router] failed to parse interview:end payload: %v", err)
		r.sendError(conn, env.RequestID, "VALIDATION_ERROR", "Dữ liệu payload không hợp lệ")
		return
	}

	// 3. Get room
	room := r.roomManager.Get(conn.RoomID)
	if room == nil {
		r.sendError(conn, env.RequestID, "NOT_FOUND", "Không tìm thấy phòng")
		return
	}

	// 4. Idempotency check: if request ID is already processed, skip.
	// Also if room is already completed, do not re-process.
	if room.RecordRequest(env.RequestID) || room.Status == events.RoomStatusCompleted {
		log.Printf("[room] interview:end requestID=%s already processed or room is completed", env.RequestID)
		return
	}

	// 5. Update room status to completed
	if !r.roomManager.SetRoomStatus(room.ID, events.RoomStatusCompleted) {
		r.sendError(conn, env.RequestID, "INVALID_STATE", "Không thể kết thúc phỏng vấn từ trạng thái hiện tại")
		return
	}

	now := time.Now().UTC()
	room.mu.Lock()
	room.EndedAt = &now
	room.mu.Unlock()

	// 6. Persist status + ended_at (real write — see handleInterviewStart).
	if r.interviewRepo != nil {
		dbCtx, cancel := context.WithTimeout(context.Background(), persistTimeout)
		if _, err := r.interviewRepo.MarkEnded(dbCtx, room.InterviewID, now); err != nil {
			log.Printf("[room] failed to persist interview end interview=%s: %v", room.InterviewID, err)
		}
		if err := r.interviewRepo.MarkRoomClosed(dbCtx, room.InterviewID, now); err != nil {
			log.Printf("[room] failed to persist room close interview=%s: %v", room.InterviewID, err)
		}
		cancel()
	}

	if r.auditLogger != nil {
		r.auditLogger.LogEvent("interview_end", conn.UserID, conn.Role, "interview_room", room.ID, "", conn.IPAddress, map[string]interface{}{"room_id": room.ID, "interview_id": room.InterviewID, "generate_report": payload.GenerateReport})
	}

	reportStatus := "failed"
	// 7. Simulate triggering report generation if requested
	if payload.GenerateReport {
		reportStatus = "generating"
		log.Printf("[api-khoi] triggering report generation for interview=%s", room.InterviewID)
	}

	// 8. Broadcast interview:completed to all participants
	completedPayload := events.InterviewCompletedPayload{
		Status:       string(events.RoomStatusCompleted),
		EndedAt:      now,
		ReportStatus: reportStatus,
	}
	completedEnv, err := events.NewEnvelope(events.EventInterviewCompleted, env.RequestID, room.ID, room.InterviewID, completedPayload)
	if err == nil {
		raw, _ := completedEnv.ToJSON()
		room.BroadcastAll(raw)
	}

	// 9. Close room and all connections after grace period
	go func() {
		// Ngắt hook lấy audio stream
		r.audioHook.StopHook(room.ID)

		// Use configurable end interview grace period (default 30s)
		grace := gracePeriodEndInterview
		time.Sleep(grace)

		// Get all connections in this room and close them
		conns := r.connManager.ByRoom(room.ID)
		log.Printf("[room] end interview grace period expired, closing %d connections in room=%s", len(conns), room.ID)
		for _, c := range conns {
			c.Close()
		}

		// Delete the room from RoomManager
		r.roomManager.Delete(room.ID)

		// Stop presence watcher
		r.presenceManager.StopWatcher(room.ID)
	}()
}

// handleInterviewPause handles the "interview:pause" event from clients.
func (r *MessageRouter) handleInterviewPause(conn *ClientConnection, env *events.Envelope) {
	// 1. Only recruiter can pause
	if conn.Role != string(events.ParticipantRecruiter) {
		r.sendError(conn, env.RequestID, "FORBIDDEN", "Chỉ Recruiter mới có quyền tạm dừng buổi phỏng vấn")
		return
	}

	// 2. Get room
	room := r.roomManager.Get(conn.RoomID)
	if room == nil {
		r.sendError(conn, env.RequestID, "NOT_FOUND", "Không tìm thấy phòng")
		return
	}

	// 3. Update room status to paused
	if !r.roomManager.SetRoomStatus(room.ID, events.RoomStatusPaused) {
		r.sendError(conn, env.RequestID, "INVALID_STATE", "Không thể tạm dừng phỏng vấn từ trạng thái hiện tại")
		return
	}

	now := time.Now().UTC()

	// 4. Simulate DB & Redis updates
	log.Printf("[db] UPDATE interview_rooms SET status = 'paused' WHERE id = '%s'", room.ID)
	log.Printf("[db] UPDATE interviews SET status = 'paused' WHERE id = '%s'", room.InterviewID)
	log.Printf("[redis] SET room_status:%s value=paused", room.ID)

	// 5. Broadcast interview:paused to all participants
	pausedPayload := events.InterviewPausedPayload{
		Status:   string(events.RoomStatusPaused),
		PausedAt: now,
	}
	pausedEnv, err := events.NewEnvelope(events.EventInterviewPaused, env.RequestID, room.ID, room.InterviewID, pausedPayload)
	if err == nil {
		raw, _ := pausedEnv.ToJSON()
		room.BroadcastAll(raw)
	}
}

// handleInterviewResume handles the "interview:resume" event from clients.
func (r *MessageRouter) handleInterviewResume(conn *ClientConnection, env *events.Envelope) {
	// 1. Only recruiter can resume
	if conn.Role != string(events.ParticipantRecruiter) {
		r.sendError(conn, env.RequestID, "FORBIDDEN", "Chỉ Recruiter mới có quyền tiếp tục buổi phỏng vấn")
		return
	}

	// 2. Get room
	room := r.roomManager.Get(conn.RoomID)
	if room == nil {
		r.sendError(conn, env.RequestID, "NOT_FOUND", "Không tìm thấy phòng")
		return
	}

	// 3. Update room status to active
	if !r.roomManager.SetRoomStatus(room.ID, events.RoomStatusActive) {
		r.sendError(conn, env.RequestID, "INVALID_STATE", "Không thể tiếp tục phỏng vấn từ trạng thái hiện tại")
		return
	}

	now := time.Now().UTC()

	// 4. Simulate DB & Redis updates
	log.Printf("[db] UPDATE interview_rooms SET status = 'active' WHERE id = '%s'", room.ID)
	log.Printf("[db] UPDATE interviews SET status = 'active' WHERE id = '%s'", room.InterviewID)
	log.Printf("[redis] SET room_status:%s value=active", room.ID)

	// 5. Broadcast interview:resumed to all participants
	resumedPayload := events.InterviewResumedPayload{
		Status:    string(events.RoomStatusActive),
		ResumedAt: now,
	}
	resumedEnv, err := events.NewEnvelope(events.EventInterviewResumed, env.RequestID, room.ID, room.InterviewID, resumedPayload)
	if err == nil {
		raw, _ := resumedEnv.ToJSON()
		room.BroadcastAll(raw)
	}
}

// handleInterviewCancel handles the "interview:cancel" event from clients.
func (r *MessageRouter) handleInterviewCancel(conn *ClientConnection, env *events.Envelope) {
	// 1. Only recruiter can cancel
	if conn.Role != string(events.ParticipantRecruiter) {
		r.sendError(conn, env.RequestID, "FORBIDDEN", "Chỉ Recruiter mới có quyền hủy buổi phỏng vấn")
		return
	}

	// 2. Get room
	room := r.roomManager.Get(conn.RoomID)
	if room == nil {
		r.sendError(conn, env.RequestID, "NOT_FOUND", "Không tìm thấy phòng")
		return
	}

	// 3. Update room status to cancelled
	if !r.roomManager.SetRoomStatus(room.ID, events.RoomStatusCancelled) {
		r.sendError(conn, env.RequestID, "INVALID_STATE", "Không thể hủy phỏng vấn từ trạng thái hiện tại")
		return
	}

	var payload struct {
		Reason string `json:"reason"`
	}
	_ = env.ParsePayload(&payload)

	now := time.Now().UTC()

	// 4. Simulate DB & Redis updates
	log.Printf("[db] UPDATE interview_rooms SET status = 'cancelled' WHERE id = '%s'", room.ID)
	log.Printf("[db] UPDATE interviews SET status = 'cancelled' WHERE id = '%s'", room.InterviewID)
	log.Printf("[redis] SET room_status:%s value=cancelled", room.ID)

	// 5. Broadcast interview:cancelled to all participants
	cancelledPayload := events.InterviewCancelledPayload{
		Status:      string(events.RoomStatusCancelled),
		CancelledAt: now,
		Reason:      payload.Reason,
	}
	cancelledEnv, err := events.NewEnvelope(events.EventInterviewCancelled, env.RequestID, room.ID, room.InterviewID, cancelledPayload)
	if err == nil {
		raw, _ := cancelledEnv.ToJSON()
		room.BroadcastAll(raw)
	}

	// 6. Close room and connections after a brief delay
	go func() {
		// Ngắt hook lấy audio stream
		r.audioHook.StopHook(room.ID)

		time.Sleep(2 * time.Second)
		conns := r.connManager.ByRoom(room.ID)
		log.Printf("[room] interview:cancel closing %d connections in room=%s", len(conns), room.ID)
		for _, c := range conns {
			c.Close()
		}
		r.roomManager.Delete(room.ID)
		r.presenceManager.StopWatcher(room.ID)
	}()
}
