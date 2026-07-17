package realtime

import (
	"fmt"
	"log"
	"strings"
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
	var speakerName, participantID, trackID, identity string
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
		speakerName = payload.SpeakerName
		participantID = payload.ParticipantID
		trackID = payload.TrackID
		identity = payload.Identity
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
		speakerName = payload.SpeakerName
		participantID = payload.ParticipantID
		trackID = payload.TrackID
		identity = payload.Identity
		content = payload.Content
		startTimeMs = payload.StartTimeMs
		endTimeMs = payload.EndTimeMs
		confidence = payload.Confidence
	}

	if speakerType == "" {
		speakerType = events.SpeakerType(conn.Role)
	}
	if speakerName == "" {
		speakerName = conn.DisplayName
	}

	resolvedPID, resolvedType, resolvedName := room.ResolveSpeaker(participantID, trackID, identity, speakerType, speakerName)

	// Deterministic TranscriptID for deduping partial → final match based on room, speaker, and start time.
	var transcriptID string
	if startTimeMs > 0 {
		transcriptID = fmt.Sprintf("%s-%s-%d", conn.RoomID, resolvedType, startTimeMs)
	} else {
		transcriptID = uuid.New().String()
	}

	if confidence < 0.5 {
		log.Printf("[transcript] low confidence (%f) detected for room=%s speaker=%s", confidence, conn.RoomID, resolvedName)
	}

	now := time.Now().UTC()
	updatePayload := events.TranscriptUpdatePayload{
		TranscriptID:  transcriptID,
		ParticipantID: resolvedPID,
		TrackID:       trackID,
		SpeakerType:   resolvedType,
		SpeakerName:   resolvedName,
		Content:       content,
		StartTimeMs:   startTimeMs,
		EndTimeMs:     endTimeMs,
		Confidence:    confidence,
		IsFinal:       isFinal,
		CreatedAt:     now,
	}

	if isFinal && r.transcriptSaver != nil {
		record := TranscriptRecord{
			ID:            transcriptID,
			InterviewID:   room.InterviewID,
			ParticipantID: resolvedPID,
			SpeakerType:   string(resolvedType),
			SpeakerName:   resolvedName,
			Content:       content,
			Language:      "vi",
			StartTimeMs:   startTimeMs,
			EndTimeMs:     endTimeMs,
			Confidence:    confidence,
			Source:        "audio",
			IsFinal:       true,
			CreatedAt:     now,
		}
		r.transcriptSaver.Push(record)
	}

	// Push to async transcript pipeline
	if r.transcriptPipeline != nil {
		r.transcriptPipeline.Push(conn.RoomID, updatePayload)
	}

	// Auto-trigger an AI follow-up suggestion when the CANDIDATE finishes a
	// substantive answer, so the recruiter gets live help without clicking.
	// Guards: candidate-only, long-enough content (skip "vâng"/"ừ"), and the
	// existing per-interview rate limiter (10/10min) to protect the AI quota.
	if isFinal && resolvedType == events.SpeakerCandidate {
		r.maybeAutoSuggest(room, content)
	}
}

// maybeAutoSuggest fires the suggestion worker for a candidate's final answer
// when it is long enough and the rate limiter allows it. Best-effort and
// non-blocking: it never fails the transcript path.
func (r *MessageRouter) maybeAutoSuggest(room *Room, content string) {
	if r.suggestionSvc == nil || r.interviewRepo == nil {
		return // no real AI wired (dev) — skip auto-suggest
	}
	if len([]rune(strings.TrimSpace(content))) < 40 {
		return // too short to be worth an AI call
	}
	if r.aiRateLimiter != nil && !r.aiRateLimiter.Allow(room.InterviewID) {
		return // respect the 10/10min budget; recruiter can still ask manually
	}

	reqID := uuid.New().String()
	thinkingPayload := events.AIThinkingPayload{
		Task:               events.AITaskSuggestFollowUp,
		RequestID:          reqID,
		Message:            "AI đang phân tích câu trả lời...",
		ExpectedDurationMs: 3000,
	}
	if thinkEnv, err := events.NewEnvelope(events.EventAIThinking, reqID, room.ID, room.InterviewID, thinkingPayload); err == nil {
		rawThink, _ := thinkEnv.ToJSON()
		room.BroadcastWithVisibility(events.EventAIThinking, "", rawThink)
	}

	go r.processAISuggestionWorker(room.ID, room.InterviewID, reqID, "", "")
}
