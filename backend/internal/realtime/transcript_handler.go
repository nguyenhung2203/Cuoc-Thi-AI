package realtime

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"backend/internal/realtime/events"
)

// handleTranscriptPartial handles incoming partial transcript events.
func (r *MessageRouter) handleTranscriptPartial(conn *ClientConnection, env *events.Envelope) {
	r.processTranscriptWebSocketEvent(conn, env, false)
}

// handleTranscriptFinal handles incoming final transcript events.
func (r *MessageRouter) handleTranscriptFinal(conn *ClientConnection, env *events.Envelope) {
	r.processTranscriptWebSocketEvent(conn, env, true)
}

func (r *MessageRouter) processTranscriptWebSocketEvent(conn *ClientConnection, env *events.Envelope, isFinal bool) {
	if env.RoomID != conn.RoomID {
		r.sendError(conn, env.RequestID, "FORBIDDEN", "Không có quyền truy cập phòng này")
		return
	}

	room := r.roomManager.Get(conn.RoomID)
	if room == nil {
		r.sendError(conn, env.RequestID, "NOT_FOUND", "Không tìm thấy phòng")
		return
	}

	var speakerType events.SpeakerType
	var content string
	var startTimeMs, endTimeMs int64
	var confidence float64

	if isFinal {
		var payload events.TranscriptFinalPayload
		if err := env.ParsePayload(&payload); err != nil {
			r.sendError(conn, env.RequestID, "VALIDATION_ERROR", "Payload final không hợp lệ")
			return
		}
		speakerType = payload.SpeakerType
		content = payload.Content
		startTimeMs = payload.StartTimeMs
		endTimeMs = payload.EndTimeMs
		confidence = payload.Confidence
	} else {
		var payload events.TranscriptPartialPayload
		if err := env.ParsePayload(&payload); err != nil {
			r.sendError(conn, env.RequestID, "VALIDATION_ERROR", "Payload partial không hợp lệ")
			return
		}
		speakerType = payload.SpeakerType
		content = payload.Content
		startTimeMs = payload.StartTimeMs
		endTimeMs = payload.EndTimeMs
		confidence = payload.Confidence
	}

	if speakerType == "" {
		speakerType = events.SpeakerType(conn.Role)
	}

	// Deterministic TranscriptID for deduping partial → final match based on room, speaker, and start time.
	var transcriptID string
	if startTimeMs > 0 {
		transcriptID = fmt.Sprintf("%s-%s-%d", conn.RoomID, speakerType, startTimeMs)
	} else {
		transcriptID = uuid.New().String()
	}

	if confidence < 0.5 {
		log.Printf("[transcript] low confidence (%f) detected for room=%s speaker=%s", confidence, conn.RoomID, conn.DisplayName)
	}

	now := time.Now().UTC()
	updatePayload := events.TranscriptUpdatePayload{
		TranscriptID: transcriptID,
		SpeakerType:  speakerType,
		SpeakerName:  conn.DisplayName,
		Content:      content,
		StartTimeMs:  startTimeMs,
		EndTimeMs:    endTimeMs,
		Confidence:   confidence,
		IsFinal:      isFinal,
		CreatedAt:    now,
	}

	if isFinal {
		log.Printf("[db] INSERT INTO interview_transcripts (id, interview_id, speaker_type, speaker_name, content, is_final, confidence, created_at) VALUES ('%s', '%s', '%s', '%s', '%s', true, %f, '%s')",
			transcriptID, room.InterviewID, speakerType, conn.DisplayName, content, confidence, now.Format(time.RFC3339))
	}

	// Push to async transcript pipeline
	if r.transcriptPipeline != nil {
		r.transcriptPipeline.Push(conn.RoomID, updatePayload)
	}
}
