package realtime

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"backend/internal/livekit"
	"backend/internal/realtime/events"
)

// Server is the WebSocket realtime gateway.
// It wires together the HTTP upgrader, connection manager, room manager, and message router.
type Server struct {
	connManager        *ConnectionManager
	roomManager        *RoomManager
	router             *MessageRouter
	httpServer         *http.Server
	transcriptPipeline *TranscriptPipeline
	transcriptSaver    *TranscriptBatchSaver
}

// NewServer creates a Server with all dependencies wired up.
func NewServer(addr string) *Server {
	cm := NewConnectionManager()
	rm := NewRoomManager()
	router := NewMessageRouter(cm, rm)

	tp := NewTranscriptPipeline(rm, 1000)
	router.SetTranscriptPipeline(tp)

	mux := http.NewServeMux()
	s := &Server{
		connManager:        cm,
		roomManager:        rm,
		router:             router,
		transcriptPipeline: tp,
		transcriptSaver:    router.GetTranscriptSaver(),
		httpServer: &http.Server{
			Addr:         addr,
			Handler:      mux,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}

	// WebSocket upgrade endpoint
	mux.HandleFunc("/ws/interview-room", s.handleUpgrade)

	// REST endpoint for chat history
	mux.HandleFunc("/api/v1/rooms/", s.handleGetChatHistory)

	// Internal callback from AI Orchestrator (transcript push, AI results)
	mux.HandleFunc("/internal/rooms/", s.handleInternal)

	// Recruiter room token endpoints
	mux.HandleFunc("POST /api/v1/companies/{company_id}/interviews/{interview_id}/room/token", s.handleGetRecruiterRoomToken)
	mux.HandleFunc("POST /companies/{company_id}/interviews/{interview_id}/room/token", s.handleGetRecruiterRoomToken)

	// Candidate invitation join endpoints
	mux.HandleFunc("GET /api/v1/interviews/join/{invite_token}", s.handleCandidateJoinByInviteToken)
	mux.HandleFunc("GET /interviews/join/{invite_token}", s.handleCandidateJoinByInviteToken)

	return s
}

// GetHandler returns the http.Handler for testing purposes.
func (s *Server) GetHandler() http.Handler {
	return s.httpServer.Handler
}

// GetRoomManager returns the RoomManager for testing purposes.
func (s *Server) GetRoomManager() *RoomManager {
	return s.roomManager
}

// ForceCheckRoom forces the presence watcher inspection for testing purposes.
func (s *Server) ForceCheckRoom(roomID string) {
	s.router.presenceManager.checkRoom(roomID)
}

// Start begins listening for connections. Blocks until the context is cancelled.
func (s *Server) Start(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		log.Printf("[realtime] listening on %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return s.Shutdown()
	}
}

// Shutdown gracefully drains all connections and stops the HTTP server.
func (s *Server) Shutdown() error {
	log.Println("[realtime] shutting down...")
	if s.transcriptSaver != nil {
		s.transcriptSaver.Close()
	}
	s.connManager.CloseAll()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}

// handleInternal routes internal callbacks from the AI Orchestrator.
func (s *Server) handleInternal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Example Route: /internal/rooms/{room_id}/transcript
	pathParts := strings.Split(r.URL.Path, "/")
	// Expected parts: ["", "internal", "rooms", "{room_id}", "transcript"]
	if len(pathParts) >= 5 && pathParts[4] == "transcript" {
		s.handleTranscriptPush(w, r, pathParts[3])
		return
	}
	if len(pathParts) >= 6 && pathParts[4] == "ai" && pathParts[5] == "suggestion" {
		s.handleAISuggestionPush(w, r, pathParts[3])
		return
	}
	if len(pathParts) >= 6 && pathParts[4] == "ai" && pathParts[5] == "score_update" {
		s.handleAIScoreUpdatePush(w, r, pathParts[3])
		return
	}
	if len(pathParts) >= 6 && pathParts[4] == "ai" && pathParts[5] == "error" {
		s.handleAIErrorPush(w, r, pathParts[3])
		return
	}

	http.NotFound(w, r)
}

// handleAIErrorPush processes the AI error webhook from Khôi's AI Orchestrator.
func (s *Server) handleAIErrorPush(w http.ResponseWriter, r *http.Request, roomID string) {
	room := s.roomManager.Get(roomID)
	if room == nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	var payload events.AIErrorPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	env, err := events.NewEnvelope(events.EventAIError, payload.RequestID, room.ID, room.InterviewID, payload)
	if err != nil {
		http.Error(w, "Failed to build envelope", http.StatusInternalServerError)
		return
	}

	raw, _ := env.ToJSON()
	room.BroadcastWithVisibility(events.EventAIError, "", raw)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// handleAISuggestionPush processes the AI suggestion webhook from Khôi's AI Orchestrator.
func (s *Server) handleAISuggestionPush(w http.ResponseWriter, r *http.Request, roomID string) {
	room := s.roomManager.Get(roomID)
	if room == nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	var payload events.AISuggestionPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	env, err := events.NewEnvelope(events.EventAISuggestion, "", room.ID, room.InterviewID, payload)
	if err != nil {
		http.Error(w, "Failed to build envelope", http.StatusInternalServerError)
		return
	}

	raw, _ := env.ToJSON()
	room.BroadcastWithVisibility(events.EventAISuggestion, "", raw)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// handleAIScoreUpdatePush processes the AI scoring update webhook from Khôi's AI Orchestrator.
func (s *Server) handleAIScoreUpdatePush(w http.ResponseWriter, r *http.Request, roomID string) {
	room := s.roomManager.Get(roomID)
	if room == nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	var payload events.AIScoreUpdatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	env, err := events.NewEnvelope(events.EventAIScoreUpdate, "", room.ID, room.InterviewID, payload)
	if err != nil {
		http.Error(w, "Failed to build envelope", http.StatusInternalServerError)
		return
	}

	raw, _ := env.ToJSON()
	room.BroadcastWithVisibility(events.EventAIScoreUpdate, "", raw)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// handleTranscriptPush processes the transcript webhook from the AI Orchestrator.
func (s *Server) handleTranscriptPush(w http.ResponseWriter, r *http.Request, roomID string) {
	room := s.roomManager.Get(roomID)
	if room == nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	var payload events.TranscriptUpdatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("[server] failed to parse transcript payload: %v", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	resolvedPID, resolvedType, resolvedName := room.ResolveSpeaker(payload.ParticipantID, payload.TrackID, "", payload.SpeakerType, payload.SpeakerName)
	payload.ParticipantID = resolvedPID
	payload.SpeakerType = resolvedType
	payload.SpeakerName = resolvedName

	// Decouple HTTP request from WebSocket broadcast via async buffer queue
	if !s.transcriptPipeline.Push(roomID, payload) {
		log.Printf("[server] transcript buffer full for room=%s", roomID)
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}

	if payload.IsFinal && s.transcriptSaver != nil {
		record := TranscriptRecord{
			ID:            payload.TranscriptID,
			InterviewID:   room.InterviewID,
			ParticipantID: resolvedPID,
			SpeakerType:   string(resolvedType),
			SpeakerName:   resolvedName,
			Content:       payload.Content,
			Language:      "vi",
			StartTimeMs:   payload.StartTimeMs,
			EndTimeMs:     payload.EndTimeMs,
			Confidence:    payload.Confidence,
			Source:        "audio",
			IsFinal:       true,
			CreatedAt:     payload.CreatedAt,
		}
		s.transcriptSaver.Push(record)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// handleGetChatHistory processes GET /api/v1/rooms/{room_id}/chat to load chat history.
func (s *Server) handleGetChatHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Authenticate using Bearer / token query param
	claims, err := extractAndValidateToken(r)
	if err != nil {
		log.Printf("[api] chat history auth failed: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 2. Extract room ID from path
	// URL format: /api/v1/rooms/{room_id}/chat
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 6 || pathParts[5] != "chat" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	roomID := pathParts[4]

	// 3. Verify user has access to this room
	if roomID != claims.RoomID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// 4. Retrieve simulated chat history
	history := s.roomManager.GetChatHistory(roomID)

	// 5. Filter history based on participant role
	filtered := make([]events.ChatMessagePayload, 0)
	for _, msg := range history {
		if msg.Visibility == events.VisibilityRecruiterOnly {
			if claims.Role == string(events.ParticipantRecruiter) {
				filtered = append(filtered, msg)
			}
		} else {
			filtered = append(filtered, msg)
		}
	}

	// Simulate database SELECT query logging
	log.Printf("[db] SELECT * FROM interview_transcripts WHERE interview_id = '%s' AND source = 'chat'", claims.InterviewID)

	// 6. Write JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    filtered,
	})
}

// handleGetRecruiterRoomToken handles generating a LiveKit token for the recruiter.
func (s *Server) handleGetRecruiterRoomToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "METHOD_NOT_ALLOWED",
			"message": "Only POST is allowed",
		})
		return
	}

	// 1. Authenticate recruiter Bearer token
	claims, err := extractAndValidateRecruiterToken(r)
	if err != nil {
		log.Printf("[api] recruiter auth failed: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "UNAUTHORIZED",
			"message": err.Error(),
		})
		return
	}

	// 2. Extract company_id and interview_id
	var companyID, interviewID string
	parts := strings.Split(r.URL.Path, "/")
	for i, part := range parts {
		if part == "companies" && i+1 < len(parts) {
			companyID = parts[i+1]
		}
		if part == "interviews" && i+1 < len(parts) {
			interviewID = parts[i+1]
		}
	}
	if companyID == "" {
		companyID = r.PathValue("company_id")
	}
	if interviewID == "" {
		interviewID = r.PathValue("interview_id")
	}

	if companyID == "" || interviewID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "BAD_REQUEST",
			"message": "Missing company_id or interview_id in path",
		})
		return
	}

	// For recruiters, let's use a deterministic room ID based on interviewID.
	// In mock environment, roomID is "room-" + interviewID.
	roomID := "room-" + interviewID

	// Log simulated DB SELECT query for verification
	log.Printf("[db] SELECT * FROM interviews WHERE id = '%s' AND company_id = '%s'", interviewID, companyID)

	// 3. Generate LiveKit token
	livekitSecret := os.Getenv("LIVEKIT_API_SECRET")
	if livekitSecret == "" {
		livekitSecret = "devsecret"
	}
	livekitKey := os.Getenv("LIVEKIT_API_KEY")
	if livekitKey == "" {
		livekitKey = "devkey"
	}

	displayName := claims.DisplayName
	if displayName == "" {
		displayName = "Recruiter"
	}

	tokenString, err := livekit.GenerateToken(
		livekitKey,
		livekitSecret,
		roomID,
		claims.UserID,
		displayName,
		"recruiter",
		interviewID,
	)
	if err != nil {
		log.Printf("[api] token generation failed: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "INTERNAL_ERROR",
			"message": "Failed to generate token",
		})
		return
	}

	livekitURL := os.Getenv("LIVEKIT_URL")
	if livekitURL == "" {
		livekitURL = "wss://livekit.example.com"
	}

	expiresAt := time.Now().Add(4 * time.Hour).UTC().Format(time.RFC3339)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"room_access_token": tokenString,
			"livekit_url":       livekitURL,
			"expires_at":        expiresAt,
		},
	})
}

// handleCandidateJoinByInviteToken handles candidate token exchange from invitation token.
func (s *Server) handleCandidateJoinByInviteToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "METHOD_NOT_ALLOWED",
			"message": "Only GET is allowed",
		})
		return
	}

	// 1. Extract invite_token
	var inviteToken string
	parts := strings.Split(r.URL.Path, "/")
	for i, part := range parts {
		if part == "join" && i+1 < len(parts) {
			inviteToken = parts[i+1]
		}
	}
	if inviteToken == "" {
		inviteToken = r.PathValue("invite_token")
	}

	if inviteToken == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "BAD_REQUEST",
			"message": "Missing invite_token in path",
		})
		return
	}

	// Simulate looking up token in DB
	log.Printf("[db] SELECT * FROM interviews WHERE invite_token_hash = SHA256('%s')", inviteToken)

	// 2. Mock checking
	if inviteToken == "invalid-token" || inviteToken == "notfound-token" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "INVITE_NOT_FOUND",
			"message": "Token không tồn tại",
		})
		return
	}

	if inviteToken == "expired-token" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "INVITE_EXPIRED",
			"message": "Token hết hạn",
		})
		return
	}

	if inviteToken == "cancelled-token" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "INTERVIEW_CANCELLED",
			"message": "Buổi phỏng vấn đã hủy",
		})
		return
	}

	// Default case - valid token
	interviewID := "interview-123"
	roomID := "room-" + interviewID
	candidateName := "Trần Văn B"
	jobTitle := "Frontend Developer"

	// 3. Generate candidate LiveKit room token
	livekitSecret := os.Getenv("LIVEKIT_API_SECRET")
	if livekitSecret == "" {
		livekitSecret = "devsecret"
	}
	livekitKey := os.Getenv("LIVEKIT_API_KEY")
	if livekitKey == "" {
		livekitKey = "devkey"
	}

	tokenString, err := livekit.GenerateToken(
		livekitKey,
		livekitSecret,
		roomID,
		"candidate-123", // candidate user ID
		candidateName,
		"candidate",
		interviewID,
	)
	if err != nil {
		log.Printf("[api] candidate token generation failed: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "INTERNAL_ERROR",
			"message": "Failed to generate candidate token",
		})
		return
	}

	now := time.Now().UTC()
	tokenExpiresAt := now.Add(4 * time.Hour)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"interview_id":                 interviewID,
			"room_id":                      roomID,
			"candidate_name":               candidateName,
			"job_title":                    jobTitle,
			"scheduled_at":                 now.Format(time.RFC3339),
			"requires_consent_ai":          true,
			"requires_consent_recording":   false,
			"room_access_token":            tokenString,
			"room_access_token_expires_at": tokenExpiresAt.Format(time.RFC3339),
		},
	})
}
