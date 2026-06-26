package realtime

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ClientConnection represents a single active WebSocket connection.
type ClientConnection struct {
	ID          string
	UserID      string
	Role        string // "recruiter" | "candidate" | "ai" | "guest"
	DisplayName string
	RoomID      string
	InterviewID string

	Conn        *websocket.Conn
	Send        chan []byte // outbound message queue
	ConnectedAt time.Time

	mu     sync.Mutex
	closed bool
}

// WriteMessage safely writes a JSON-encoded byte slice to the connection.
// It is safe for concurrent use.
func (c *ClientConnection) WriteMessage(data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed || c.Conn == nil {
		return nil
	}
	c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.Conn.WriteMessage(websocket.TextMessage, data)
}

// TrySend safely attempts to enqueue a message on Send channel without blocking or panicking.
func (c *ClientConnection) TrySend(data []byte) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return false
	}
	select {
	case c.Send <- data:
		return true
	default:
		return false
	}
}

// Close marks the connection as closed and closes the underlying WebSocket.
// Safe to call multiple times and safe when Conn is nil (e.g. in unit tests).
func (c *ClientConnection) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return
	}
	c.closed = true
	// Drain and close the Send channel so writePump goroutine exits cleanly.
	close(c.Send)
	if c.Conn != nil {
		c.Conn.Close()
	}
}
