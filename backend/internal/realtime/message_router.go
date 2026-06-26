package realtime

import (
	"log"
	"sync"
	"time"

	"backend/internal/livekit"
	"backend/internal/realtime/events"
	"backend/internal/repository"
)

// HandlerFunc is the signature for event handlers.
type HandlerFunc func(conn *ClientConnection, env *events.Envelope)

// MessageRouter maps incoming event names to their handler functions.
type MessageRouter struct {
	handlers           map[string]HandlerFunc
	connManager        *ConnectionManager
	roomManager        *RoomManager
	presenceManager    *PresenceManager
	audioHook          *livekit.AudioHookService
	transcriptPipeline *TranscriptPipeline
	transcriptSaver    *TranscriptBatchSaver
	aiRateLimiter      *AIRateLimiter
	scoreRateLimiter   *AIRateLimiter
	auditLogger        *AuditLogger
	aiRetryMu          sync.Mutex
	aiRetries          map[string]int
}

// NewMessageRouter constructs a router and registers all known event handlers.
func NewMessageRouter(cm *ConnectionManager, rm *RoomManager) *MessageRouter {
	r := &MessageRouter{
		handlers:         make(map[string]HandlerFunc),
		connManager:      cm,
		roomManager:      rm,
		presenceManager:  NewPresenceManager(rm),
		audioHook:        livekit.NewAudioHookService(),
		transcriptSaver:  NewTranscriptBatchSaver(100, 5*time.Second),
		aiRateLimiter:    NewAIRateLimiter(10, 10*time.Minute),
		scoreRateLimiter: NewAIRateLimiter(10, 10*time.Minute),
		auditLogger:      NewAuditLogger(repository.NewAuditRepository()),
		aiRetries:        make(map[string]int),
	}
	r.registerHandlers()
	return r
}

// SetTranscriptPipeline attaches the async transcript pipeline to the router.
func (r *MessageRouter) SetTranscriptPipeline(tp *TranscriptPipeline) {
	r.transcriptPipeline = tp
}

// GetTranscriptSaver returns the batch saver instance.
func (r *MessageRouter) GetTranscriptSaver() *TranscriptBatchSaver {
	return r.transcriptSaver
}

// registerHandlers wires every client→server event to its handler.
// Handlers are implemented in backend/internal/realtime/handlers/*.go
func (r *MessageRouter) registerHandlers() {
	// Room lifecycle
	r.register(events.EventRoomJoin, r.handleRoomJoin)
	r.register(events.EventRoomLeave, r.handleRoomLeave)
	r.register(events.EventRoomHeartbeat, r.handleRoomHeartbeat)

	// Interview control
	r.register(events.EventInterviewStart, r.handleInterviewStart)
	r.register(events.EventInterviewEnd, r.handleInterviewEnd)
	r.register(events.EventInterviewPause, r.handleInterviewPause)
	r.register(events.EventInterviewResume, r.handleInterviewResume)
	r.register(events.EventInterviewCancel, r.handleInterviewCancel)

	// Media
	r.register(events.EventMediaStatus, r.handleMediaStatus)

	// Chat & notes
	r.register(events.EventChatSend, r.handleChatSend)
	r.register(events.EventNoteCreate, r.stubHandler("note:create"))
	r.register(events.EventQuestionMarkAsked, r.stubHandler("question:mark_asked"))

	// Transcript (from STT pipeline)
	r.register(events.EventTranscriptPartial, r.handleTranscriptPartial)
	r.register(events.EventTranscriptFinal, r.handleTranscriptFinal)

	// AI requests
	r.register(events.EventAIRequestSuggestion, r.handleAIRequestSuggestion)
	r.register(events.EventAIRequestScoreUpdate, r.handleAIRequestScoreUpdate)
}

// register adds a handler for the given event name.
func (r *MessageRouter) register(event string, h HandlerFunc) {
	r.handlers[event] = h
}

// Route dispatches an incoming envelope to the appropriate handler.
// Unknown events are logged and an error is sent back to the client.
func (r *MessageRouter) Route(conn *ClientConnection, env *events.Envelope) {
	h, ok := r.handlers[env.Event]
	if !ok {
		log.Printf("[router] unknown event=%q connID=%s", env.Event, conn.ID)
		return
	}
	h(conn, env)
}

// stubHandler returns a placeholder HandlerFunc that logs the event.
// Sprint 1–4 tasks will replace each stub with a real implementation.
func (r *MessageRouter) stubHandler(name string) HandlerFunc {
	return func(conn *ClientConnection, env *events.Envelope) {
		log.Printf("[router] stub handler for event=%q connID=%s requestID=%s",
			name, conn.ID, env.RequestID)
	}
}
