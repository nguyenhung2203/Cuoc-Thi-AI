package realtime

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"backend/internal/realtime/events"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	// Allow all origins in development; restrict in production via env config.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// handleUpgrade upgrades an HTTP request to a WebSocket connection,
// authenticates the room_access_token, registers the connection, and
// starts the read/write pump goroutines.
func (s *Server) handleUpgrade(w http.ResponseWriter, r *http.Request) {
	// ── 1. Authenticate ──────────────────────────────────────────────────────
	claims, err := extractAndValidateToken(r)
	if err != nil {
		log.Printf("[ws] auth failed: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// ── 1.1 Room Validation ──────────────────────────────────────────────────
	// If room_id is explicitly requested in query, it must match the token's room ID
	if reqRoomID := r.URL.Query().Get("room_id"); reqRoomID != "" && reqRoomID != claims.RoomID {
		log.Printf("[ws] auth failed: room mismatch, request=%s token=%s", reqRoomID, claims.RoomID)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// ── 2. Upgrade HTTP → WebSocket ──────────────────────────────────────────
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ws] upgrade failed: %v", err)
		return
	}

	// ── 3. Build ClientConnection ────────────────────────────────────────────
	connID := uuid.New().String()
	conn := newClientConnection(
		connID,
		claims.UserID,
		claims.Role,
		claims.DisplayName,
		claims.RoomID,
		claims.InterviewID,
	)
	conn.IPAddress = r.RemoteAddr
	conn.Conn = wsConn

	s.connManager.Add(conn)
	log.Printf("[ws] connected connID=%s userID=%s role=%s roomID=%s", connID, claims.UserID, claims.Role, claims.RoomID)

	// ── 4. Start pumps ───────────────────────────────────────────────────────
	go s.writePump(conn)
	go s.readPump(conn)
}

// readPump reads inbound messages from the client and routes them.
// It runs in its own goroutine and cleans up when the connection closes.
func (s *Server) readPump(conn *ClientConnection) {
	defer func() {
		s.onDisconnect(conn)
	}()

	conn.Conn.SetReadLimit(maxMessageSize)
	conn.Conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.Conn.SetPongHandler(func(string) error {
		conn.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, raw, err := conn.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[ws] read error connID=%s: %v", conn.ID, err)
			}
			return
		}

		var env events.Envelope
		if err := json.Unmarshal(raw, &env); err != nil {
			log.Printf("[ws] invalid JSON from connID=%s: %v", conn.ID, err)
			// Send error back to client but keep connection open
			s.sendError(conn, "", "INVALID_JSON", "Dữ liệu không hợp lệ", false)
			continue
		}

		s.router.Route(conn, &env)
	}
}

// writePump drains the connection's Send channel and writes messages to the WebSocket.
// It also sends periodic pings to keep the connection alive.
func (s *Server) writePump(conn *ClientConnection) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		conn.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-conn.Send:
			conn.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Channel closed → send close frame
				conn.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := conn.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("[ws] write error connID=%s: %v", conn.ID, err)
				return
			}

		case <-ticker.C:
			conn.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// onDisconnect handles cleanup when a connection drops.
func (s *Server) onDisconnect(conn *ClientConnection) {
	log.Printf("[ws] disconnected connID=%s userID=%s", conn.ID, conn.UserID)

	if s.router != nil && s.router.auditLogger != nil {
		s.router.auditLogger.LogEvent(
			"disconnect",
			conn.UserID,
			conn.Role,
			"interview_room",
			conn.RoomID,
			"",
			conn.IPAddress,
			map[string]interface{}{"room_id": conn.RoomID, "connection_id": conn.ID},
		)
	}

	s.connManager.Remove(conn.ID)
	conn.Close()

	if conn.RoomID != "" {
		room := s.roomManager.Get(conn.RoomID)
		if room != nil {
			p := room.GetParticipant(conn.ID)
			if p != nil {
				room.mu.Lock()
				p.ConnectionState = events.ConnectionReconnecting
				p.LastSeenAt = time.Now().UTC()
				p.Connection = nil
				room.mu.Unlock()

				log.Printf("[db] UPDATE interview_participants SET connection_state = 'reconnecting', last_seen_at = '%s' WHERE id = '%s'",
					p.LastSeenAt.Format(time.RFC3339), conn.ID)
				log.Printf("[redis] SET presence:%s:%s value=reconnecting EX 120", conn.RoomID, conn.ID)

				presencePayload := events.RoomPresenceUpdatePayload{
					Participants: room.ParticipantList(),
				}
				presenceEnv, err := events.NewEnvelope(events.EventRoomPresenceUpdate, "", room.ID, room.InterviewID, presencePayload)
				if err == nil {
					rawPresence, _ := presenceEnv.ToJSON()
					room.BroadcastAll(rawPresence)
				}
			}

			if !room.HasOnlineParticipants() {
				s.roomManager.StartEmptyRoomTimer(room.ID, s.router.presenceManager)
			}
		}
	}
}

// sendError enqueues a generic error event to a single connection.
func (s *Server) sendError(conn *ClientConnection, requestID, code, message string, recoverable bool) {
	env, err := events.NewEnvelope(events.EventError, requestID, conn.RoomID, conn.InterviewID,
		events.ErrorPayload{Code: code, Message: message, Recoverable: recoverable})
	if err != nil {
		return
	}
	raw, _ := env.ToJSON()
	if !conn.TrySend(raw) {
		log.Printf("[ws] send buffer full/closed for connID=%s, dropping error event", conn.ID)
	}
}
