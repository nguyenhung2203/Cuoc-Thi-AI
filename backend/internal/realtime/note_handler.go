package realtime

import (
	"log"
	"time"

	"github.com/google/uuid"

	"backend/internal/realtime/events"
)

// handleNoteCreate handles the "note:create" event. Recruiter-only. The note is
// persisted (as a recruiter-only transcript record) and echoed back to
// recruiters in the room so multiple interviewers stay in sync.
func (r *MessageRouter) handleNoteCreate(conn *ClientConnection, env *events.Envelope) {
	if conn.Role != string(events.ParticipantRecruiter) {
		log.Printf("[note] create forbidden for role=%q client=%s", conn.Role, conn.ID)
		r.sendError(conn, env.RequestID, "FORBIDDEN", "Chỉ Recruiter mới có quyền tạo ghi chú")
		return
	}

	room := r.roomManager.Get(conn.RoomID)
	if room == nil {
		r.sendError(conn, env.RequestID, "NOT_FOUND", "Không tìm thấy phòng")
		return
	}

	var payload events.NoteCreatePayload
	if err := env.ParsePayload(&payload); err != nil {
		r.sendError(conn, env.RequestID, "VALIDATION_ERROR", "Dữ liệu ghi chú không hợp lệ")
		return
	}
	if payload.Content == "" {
		r.sendError(conn, env.RequestID, "VALIDATION_ERROR", "Nội dung ghi chú không được để trống")
		return
	}

	now := time.Now().UTC()
	noteID := uuid.New().String()

	// Persist as a recruiter-only note record via the batch saver.
	if r.transcriptSaver != nil {
		r.transcriptSaver.Push(TranscriptRecord{
			ID:            noteID,
			InterviewID:   room.InterviewID,
			ParticipantID: conn.ID,
			SpeakerType:   string(events.ParticipantRecruiter),
			SpeakerName:   conn.DisplayName,
			Content:       payload.Content,
			Language:      "vi",
			Confidence:    1.0,
			Source:        "note",
			IsFinal:       true,
			CreatedAt:     now,
		})
	}

	if r.auditLogger != nil {
		r.auditLogger.LogEvent("note_create", conn.UserID, conn.Role, "interview_room", room.ID, "", conn.IPAddress,
			map[string]interface{}{"room_id": room.ID, "interview_id": room.InterviewID, "note_id": noteID})
	}

	out := events.NoteCreatedPayload{
		NoteID:     noteID,
		AuthorID:   conn.ID,
		AuthorName: conn.DisplayName,
		Content:    payload.Content,
		Tags:       payload.Tags,
		CreatedAt:  now,
	}
	if noteEnv, err := events.NewEnvelope(events.EventNoteCreated, env.RequestID, room.ID, room.InterviewID, out); err == nil {
		if raw, mErr := noteEnv.ToJSON(); mErr == nil {
			room.BroadcastRecruitersOnly(raw)
		}
	}
}

// handleQuestionMarkAsked handles "question:mark_asked". Recruiter-only. Records
// that a prepared question was asked and notifies recruiters in the room.
func (r *MessageRouter) handleQuestionMarkAsked(conn *ClientConnection, env *events.Envelope) {
	if conn.Role != string(events.ParticipantRecruiter) {
		log.Printf("[question] mark_asked forbidden for role=%q client=%s", conn.Role, conn.ID)
		r.sendError(conn, env.RequestID, "FORBIDDEN", "Chỉ Recruiter mới có quyền đánh dấu câu hỏi")
		return
	}

	room := r.roomManager.Get(conn.RoomID)
	if room == nil {
		r.sendError(conn, env.RequestID, "NOT_FOUND", "Không tìm thấy phòng")
		return
	}

	var payload events.QuestionMarkAskedPayload
	if err := env.ParsePayload(&payload); err != nil {
		r.sendError(conn, env.RequestID, "VALIDATION_ERROR", "Dữ liệu câu hỏi không hợp lệ")
		return
	}
	if payload.QuestionID == "" {
		r.sendError(conn, env.RequestID, "VALIDATION_ERROR", "Thiếu question_id")
		return
	}

	now := time.Now().UTC()
	if r.auditLogger != nil {
		r.auditLogger.LogEvent("question_mark_asked", conn.UserID, conn.Role, "interview_room", room.ID, "", conn.IPAddress,
			map[string]interface{}{"room_id": room.ID, "interview_id": room.InterviewID, "question_id": payload.QuestionID})
	}

	out := events.QuestionAskedPayload{
		QuestionID: payload.QuestionID,
		MarkedByID: conn.ID,
		AskedAtMs:  payload.AskedAtMs,
		MarkedAt:   now,
	}
	if qEnv, err := events.NewEnvelope(events.EventQuestionAsked, env.RequestID, room.ID, room.InterviewID, out); err == nil {
		if raw, mErr := qEnv.ToJSON(); mErr == nil {
			room.BroadcastRecruitersOnly(raw)
		}
	}
}
