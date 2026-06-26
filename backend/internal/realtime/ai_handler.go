package realtime

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"

	"backend/internal/realtime/events"
)

// AIRateLimiter restricts AI requests per interview window to prevent budget exhaust and spam.
type AIRateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

// NewAIRateLimiter constructs a thread-safe rate limiter.
func NewAIRateLimiter(limit int, window time.Duration) *AIRateLimiter {
	if limit <= 0 {
		limit = 10
	}
	if window <= 0 {
		window = 10 * time.Minute
	}
	return &AIRateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// Allow returns true if the interview has not exceeded its request limit.
func (l *AIRateLimiter) Allow(interviewID string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-l.window)

	times := l.requests[interviewID]
	var valid []time.Time
	for _, t := range times {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= l.limit {
		l.requests[interviewID] = valid
		return false
	}

	valid = append(valid, now)
	l.requests[interviewID] = valid
	return true
}

// sendAIError broadcasts an AI error event to Recruiters only.
func (r *MessageRouter) sendAIError(room *Room, reqID string, errType events.AIErrorType, severity events.AISeverity, msg string, recoverable bool, retryAfterSec int) {
	if room == nil {
		return
	}
	errPayload := events.AIErrorPayload{
		ErrorType:         errType,
		RequestID:         reqID,
		AffectedFeatures:  []string{"suggestion", "scoring"},
		Severity:          severity,
		Message:           msg,
		Recoverable:       recoverable,
		RetryAfterSeconds: retryAfterSec,
	}
	errEnv, err := events.NewEnvelope(events.EventAIError, reqID, room.ID, room.InterviewID, errPayload)
	if err == nil {
		rawErr, _ := errEnv.ToJSON()
		room.BroadcastWithVisibility(events.EventAIError, "", rawErr)
	}

	if r.auditLogger != nil {
		r.auditLogger.LogEvent("ai_error", "system", "ai", "interview_room", room.ID, "", "127.0.0.1", map[string]interface{}{"room_id": room.ID, "interview_id": room.InterviewID, "error_type": errType})
	}
}

// handleAIRequestSuggestion processes "ai:request_suggestion" events from recruiters.
func (r *MessageRouter) handleAIRequestSuggestion(conn *ClientConnection, env *events.Envelope) {
	// 1. Enforce Recruiter-only access control
	if conn.Role != string(events.ParticipantRecruiter) {
		log.Printf("[ai] suggestion request forbidden for role=%q client=%s", conn.Role, conn.ID)
		r.sendError(conn, env.RequestID, "FORBIDDEN", "Chỉ Recruiter mới có quyền yêu cầu AI suggestion")
		return
	}

	room := r.roomManager.Get(conn.RoomID)
	if room == nil {
		r.sendError(conn, env.RequestID, "NOT_FOUND", "Không tìm thấy phòng")
		return
	}

	// 2. Enforce Rate Limiting (max 10 requests / 10 minutes / interview)
	if r.aiRateLimiter != nil && !r.aiRateLimiter.Allow(room.InterviewID) {
		log.Printf("[ai] rate limit exceeded for interview=%s room=%s", room.InterviewID, room.ID)
		r.sendAIError(room, env.RequestID, events.AIErrorRateLimited, events.SeverityDegraded, "Vượt giới hạn 10 yêu cầu AI trong 10 phút cho phiên phỏng vấn này", true, 60)
		return
	}

	// 3. Parse payload
	var payload events.AIRequestSuggestionPayload
	if err := env.ParsePayload(&payload); err != nil {
		log.Printf("[ai] failed to parse request suggestion payload: %v", err)
	}

	// 4. Send "ai:thinking" immediately to Recruiters only
	thinkingPayload := events.AIThinkingPayload{
		Task:               events.AITaskSuggestFollowUp,
		RequestID:          env.RequestID,
		Message:            "AI đang phân tích câu trả lời...",
		ExpectedDurationMs: 3000,
	}
	thinkEnv, err := events.NewEnvelope(events.EventAIThinking, env.RequestID, room.ID, room.InterviewID, thinkingPayload)
	if err == nil {
		rawThink, _ := thinkEnv.ToJSON()
		room.BroadcastWithVisibility(events.EventAIThinking, conn.ID, rawThink)
	}

	// 5. Spawn async worker to simulate Python AI Orchestrator analysis & reply
	go r.processAISuggestionWorker(room.ID, room.InterviewID, env.RequestID, payload.Focus)
}

func (r *MessageRouter) processAISuggestionWorker(roomID, interviewID, reqID, focus string) {
	time.Sleep(50 * time.Millisecond) // simulate AI network latency

	targetRoom := r.roomManager.Get(roomID)
	if targetRoom == nil {
		return
	}

	if focus == "simulate_down" {
		log.Printf("[ai] service unavailable (critical) for room=%s", roomID)
		r.sendAIError(targetRoom, reqID, events.AIErrorServiceUnavail, events.SeverityCritical, "AI tạm thời không phản hồi. Buổi phỏng vấn vẫn tiếp tục.", false, 0)
		return
	}

	if focus == "simulate_timeout" || focus == "simulate_timeout_permanent" {
		r.aiRetryMu.Lock()
		retries := r.aiRetries[reqID]
		r.aiRetries[reqID] = retries + 1
		r.aiRetryMu.Unlock()

		if retries < 3 {
			log.Printf("[ai] timeout for room=%s, retry count=%d", roomID, retries+1)
			r.sendAIError(targetRoom, reqID, events.AIErrorTimeout, events.SeverityDegraded, "AI phản hồi chậm (timeout). Đang thử lại...", true, 1)
			if focus == "simulate_timeout" {
				time.AfterFunc(50*time.Millisecond, func() {
					r.processAISuggestionWorker(roomID, interviewID, reqID, "normal")
				})
			} else if focus == "simulate_timeout_permanent" {
				time.AfterFunc(50*time.Millisecond, func() {
					r.processAISuggestionWorker(roomID, interviewID, reqID, "simulate_timeout_permanent")
				})
			}
			return
		}

		log.Printf("[ai] max retries reached (critical) for room=%s", roomID)
		r.sendAIError(targetRoom, reqID, events.AIErrorServiceUnavail, events.SeverityCritical, "AI tạm thời không phản hồi. Buổi phỏng vấn vẫn tiếp tục.", false, 0)
		return
	}

	suggestionID := uuid.New().String()
	suggPayload := events.AISuggestionPayload{
		SuggestionID:   suggestionID,
		SuggestionType: events.SuggestionFollowUpQuestion,
		Content:        "Bạn có thể nói rõ bạn đã đo performance bằng chỉ số nào không?",
		Reason:         "Ứng viên nói đã tối ưu performance nhưng chưa nêu metric cụ thể.",
		TargetSkill:    "Performance Optimization",
		Priority:       "high",
		Confidence:     0.84,
	}

	log.Printf("[db] INSERT INTO ai_suggestions (id, interview_id, suggestion_type, content, reason, target_skill, priority, confidence, created_at) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', %f, '%s')",
		suggestionID, interviewID, suggPayload.SuggestionType, suggPayload.Content, suggPayload.Reason, suggPayload.TargetSkill, suggPayload.Priority, suggPayload.Confidence, time.Now().UTC().Format(time.RFC3339))

	suggEnv, err := events.NewEnvelope(events.EventAISuggestion, reqID, roomID, interviewID, suggPayload)
	if err == nil {
		rawSugg, _ := suggEnv.ToJSON()
		targetRoom.BroadcastWithVisibility(events.EventAISuggestion, "", rawSugg)
	}

	if r.auditLogger != nil {
		r.auditLogger.LogEvent("ai_suggestion", "system", "ai", "interview_room", roomID, "", "127.0.0.1", map[string]interface{}{"room_id": roomID, "suggestion_id": suggestionID})
	}
}

// handleAIRequestScoreUpdate processes "ai:request_score_update" events from recruiters.
func (r *MessageRouter) handleAIRequestScoreUpdate(conn *ClientConnection, env *events.Envelope) {
	// 1. Enforce Recruiter-only access control
	if conn.Role != string(events.ParticipantRecruiter) {
		log.Printf("[ai] score update request forbidden for role=%q client=%s", conn.Role, conn.ID)
		r.sendError(conn, env.RequestID, "FORBIDDEN", "Chỉ Recruiter mới có quyền yêu cầu chấm điểm AI")
		return
	}

	room := r.roomManager.Get(conn.RoomID)
	if room == nil {
		r.sendError(conn, env.RequestID, "NOT_FOUND", "Không tìm thấy phòng")
		return
	}

	// 2. Enforce Rate Limiting (max 10 requests / 10 minutes / interview)
	if r.scoreRateLimiter != nil && !r.scoreRateLimiter.Allow(room.InterviewID) {
		log.Printf("[ai] score rate limit exceeded for interview=%s room=%s", room.InterviewID, room.ID)
		r.sendAIError(room, env.RequestID, events.AIErrorRateLimited, events.SeverityDegraded, "Vượt giới hạn 10 yêu cầu chấm điểm AI trong 10 phút cho phiên phỏng vấn này", true, 60)
		return
	}

	// 3. Parse payload
	var payload events.AIRequestScoreUpdatePayload
	if err := env.ParsePayload(&payload); err != nil {
		log.Printf("[ai] failed to parse score update payload: %v", err)
	}

	// 4. Send "ai:thinking" immediately to Recruiters only
	thinkingPayload := events.AIThinkingPayload{
		Task:               events.AITaskScoreAnswer,
		RequestID:          env.RequestID,
		Message:            "AI đang chấm điểm câu trả lời theo rubric...",
		ExpectedDurationMs: 5000,
	}
	thinkEnv, err := events.NewEnvelope(events.EventAIThinking, env.RequestID, room.ID, room.InterviewID, thinkingPayload)
	if err == nil {
		rawThink, _ := thinkEnv.ToJSON()
		room.BroadcastWithVisibility(events.EventAIThinking, conn.ID, rawThink)
	}

	// 5. Spawn async worker to simulate AI scoring
	go r.processAIScoreWorker(room.ID, room.InterviewID, env.RequestID, payload.Scope)
}

func (r *MessageRouter) processAIScoreWorker(roomID, interviewID, reqID, scope string) {
	time.Sleep(50 * time.Millisecond)

	targetRoom := r.roomManager.Get(roomID)
	if targetRoom == nil {
		return
	}

	if scope == "simulate_down" {
		log.Printf("[ai] score service unavailable (critical) for room=%s", roomID)
		r.sendAIError(targetRoom, reqID, events.AIErrorServiceUnavail, events.SeverityCritical, "AI tạm thời không phản hồi. Buổi phỏng vấn vẫn tiếp tục.", false, 0)
		return
	}

	if scope == "simulate_timeout" || scope == "simulate_timeout_permanent" {
		r.aiRetryMu.Lock()
		retries := r.aiRetries[reqID]
		r.aiRetries[reqID] = retries + 1
		r.aiRetryMu.Unlock()

		if retries < 3 {
			log.Printf("[ai] score timeout for room=%s, retry count=%d", roomID, retries+1)
			r.sendAIError(targetRoom, reqID, events.AIErrorTimeout, events.SeverityDegraded, "AI phản hồi chậm (timeout). Đang thử lại...", true, 1)
			if scope == "simulate_timeout" {
				time.AfterFunc(50*time.Millisecond, func() {
					r.processAIScoreWorker(roomID, interviewID, reqID, "normal")
				})
			} else if scope == "simulate_timeout_permanent" {
				time.AfterFunc(50*time.Millisecond, func() {
					r.processAIScoreWorker(roomID, interviewID, reqID, "simulate_timeout_permanent")
				})
			}
			return
		}

		log.Printf("[ai] score max retries reached (critical) for room=%s", roomID)
		r.sendAIError(targetRoom, reqID, events.AIErrorServiceUnavail, events.SeverityCritical, "AI tạm thời không phản hồi. Buổi phỏng vấn vẫn tiếp tục.", false, 0)
		return
	}

	var rubricScores []events.RubricScore
	if scope == "brief_check" {
		rubricScores = append(rubricScores, events.RubricScore{
			CriterionName: "Technical Depth",
			Score:         nil,
			MaxScore:      5,
			Evidence:      "Câu trả lời quá ngắn, không đủ bằng chứng đánh giá.",
			Confidence:    0.20,
			Status:        events.ScoreStatusInsufficientEvidence,
		})
	} else {
		scoreVal := 4.0
		rubricScores = append(rubricScores, events.RubricScore{
			CriterionName: "Technical Knowledge",
			Score:         &scoreVal,
			MaxScore:      5,
			Evidence:      "Ứng viên mô tả được cách tối ưu query và cache.",
			Confidence:    0.78,
			Status:        events.ScoreStatusScored,
		})
	}

	scoreID := uuid.New().String()
	for _, rs := range rubricScores {
		sVal := 0.0
		if rs.Score != nil {
			sVal = *rs.Score
		}
		log.Printf("[db] INSERT INTO interview_scores (id, interview_id, criterion_name, score, max_score, evidence, confidence, status, created_at) VALUES ('%s', '%s', '%s', %f, %f, '%s', %f, '%s', '%s')",
			scoreID, interviewID, rs.CriterionName, sVal, rs.MaxScore, rs.Evidence, rs.Confidence, rs.Status, time.Now().UTC().Format(time.RFC3339))
	}

	updatePayload := events.AIScoreUpdatePayload{
		Scores: rubricScores,
	}
	scoreEnv, err := events.NewEnvelope(events.EventAIScoreUpdate, reqID, roomID, interviewID, updatePayload)
	if err == nil {
		rawScore, _ := scoreEnv.ToJSON()
		targetRoom.BroadcastWithVisibility(events.EventAIScoreUpdate, "", rawScore)
	}

	if r.auditLogger != nil {
		r.auditLogger.LogEvent("ai_score_update", "system", "ai", "interview_room", roomID, "", "127.0.0.1", map[string]interface{}{"room_id": roomID, "score_id": scoreID})
	}
}
