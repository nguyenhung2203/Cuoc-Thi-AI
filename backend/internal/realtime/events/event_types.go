package events

// Client → Server events
const (
	// Room lifecycle
	EventRoomJoin      = "room:join"
	EventRoomLeave     = "room:leave"
	EventRoomHeartbeat = "room:heartbeat"

	// Interview control (Recruiter only)
	EventInterviewStart  = "interview:start"
	EventInterviewEnd    = "interview:end"
	EventInterviewPause  = "interview:pause"
	EventInterviewResume = "interview:resume"
	EventInterviewCancel = "interview:cancel"

	// Media
	EventMediaStatus = "media:status"

	// Chat
	EventChatSend = "chat:send"

	// Notes (Recruiter only)
	EventNoteCreate = "note:create"

	// Questions (Recruiter only)
	EventQuestionMarkAsked = "question:mark_asked"

	// Transcript (System/STT)
	EventTranscriptPartial = "transcript:partial"
	EventTranscriptFinal   = "transcript:final"

	// AI (Recruiter only)
	EventAIRequestSuggestion  = "ai:request_suggestion"
	EventAIRequestScoreUpdate = "ai:request_score_update"

	// WebRTC signaling (raw mode, not used with LiveKit)
	EventRTCOffer        = "rtc:offer"
	EventRTCAnswer       = "rtc:answer"
	EventRTCIceCandidate = "rtc:ice_candidate"
)

// Server → Client events
const (
	// Room events
	EventRoomJoined         = "room:joined"
	EventRoomUserJoined     = "room:user_joined"
	EventRoomUserLeft       = "room:user_left"
	EventRoomPresenceUpdate = "room:presence_update"

	// Interview events
	EventInterviewStarted   = "interview:started"
	EventInterviewCompleted = "interview:completed"
	EventInterviewPaused    = "interview:paused"
	EventInterviewResumed   = "interview:resumed"
	EventInterviewCancelled = "interview:cancelled"
	EventRoomExpired        = "room:expired"

	// Media events
	EventMediaStatusChanged = "media:status_changed"

	// Chat events
	EventChatMessage = "chat:message"

	// Note & question events (Recruiter only)
	EventNoteCreated   = "note:created"
	EventQuestionAsked = "question:asked"

	// Transcript events
	EventTranscriptUpdate = "transcript:update"

	// AI events (Recruiter only)
	EventAIThinking    = "ai:thinking"
	EventAISuggestion  = "ai:suggestion"
	EventAIScoreUpdate = "ai:score_update"
	EventAIWarning     = "ai:warning"
	EventAIError       = "ai:error"

	// Report events (Recruiter only)
	EventReportReady = "report:ready"

	// Error
	EventError = "error"
)

// ParticipantType defines who is in the room
type ParticipantType string

const (
	ParticipantRecruiter ParticipantType = "recruiter"
	ParticipantCandidate ParticipantType = "candidate"
	ParticipantAI        ParticipantType = "ai"
	ParticipantGuest     ParticipantType = "guest"
)

// ChatVisibility defines who can see a chat message
type ChatVisibility string

const (
	VisibilityRoom          ChatVisibility = "room"
	VisibilityRecruiterOnly ChatVisibility = "recruiter_only"
)

// ConnectionState tracks participant connection status
type ConnectionState string

const (
	ConnectionOnline       ConnectionState = "online"
	ConnectionReconnecting ConnectionState = "reconnecting"
	ConnectionOffline      ConnectionState = "offline"
	ConnectionLeft         ConnectionState = "left"
)

// RoomStatus tracks interview room lifecycle
type RoomStatus string

const (
	RoomStatusScheduled       RoomStatus = "scheduled"
	RoomStatusWaiting         RoomStatus = "waiting"
	RoomStatusCandidateJoined RoomStatus = "candidate_joined"
	RoomStatusRecruiterJoined RoomStatus = "recruiter_joined"
	RoomStatusActive          RoomStatus = "active"
	RoomStatusPaused          RoomStatus = "paused"
	RoomStatusCompleted       RoomStatus = "completed"
	RoomStatusCancelled       RoomStatus = "cancelled"
	RoomStatusExpired         RoomStatus = "expired"
)

// SpeakerType identifies who said something in a transcript
type SpeakerType string

const (
	SpeakerRecruiter SpeakerType = "recruiter"
	SpeakerCandidate SpeakerType = "candidate"
	SpeakerAI        SpeakerType = "ai"
	SpeakerSystem    SpeakerType = "system"
	SpeakerUnknown   SpeakerType = "unknown"
)

// AISuggestionType categorises suggestions
type AISuggestionType string

const (
	SuggestionFollowUpQuestion AISuggestionType = "follow_up_question"
	SuggestionWarning          AISuggestionType = "warning"
	SuggestionSummary          AISuggestionType = "summary"
	SuggestionQuestion         AISuggestionType = "question"
)

// AIErrorType categorises AI failures
type AIErrorType string

const (
	AIErrorThrottled      AIErrorType = "ai_throttled"
	AIErrorTimeout        AIErrorType = "ai_timeout"
	AIErrorServiceUnavail AIErrorType = "ai_service_unavailable"
	AIErrorRateLimited    AIErrorType = "ai_rate_limited"
	AIErrorInvalidInput   AIErrorType = "ai_invalid_input"
)

// AISeverity defines how bad an AI failure is
type AISeverity string

const (
	SeverityInfo     AISeverity = "info"
	SeverityDegraded AISeverity = "degraded"
	SeverityCritical AISeverity = "critical"
)

// AITaskType names the AI task being executed
type AITaskType string

const (
	AITaskSuggestFollowUp AITaskType = "suggest_follow_up"
	AITaskScoreAnswer     AITaskType = "score_answer"
	AITaskGenerateSummary AITaskType = "generate_summary"
)

// ScoreStatus represents whether a rubric criterion was scored
type ScoreStatus string

const (
	ScoreStatusScored               ScoreStatus = "scored"
	ScoreStatusInsufficientEvidence ScoreStatus = "insufficient_evidence"
	ScoreStatusManualOverride       ScoreStatus = "manual_override"
)
