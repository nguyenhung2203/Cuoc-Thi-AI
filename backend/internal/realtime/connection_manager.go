package realtime

import (
	"sync"
	"time"
)

// ConnectionManager maintains a registry of all live WebSocket connections.
// It is safe for concurrent use.
type ConnectionManager struct {
	connections map[string]*ClientConnection // connection_id → connection
	mu          sync.RWMutex
}

// NewConnectionManager creates a ready-to-use ConnectionManager.
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*ClientConnection),
	}
}

// Add registers a new connection.
func (cm *ConnectionManager) Add(conn *ClientConnection) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.connections[conn.ID] = conn
}

// Remove unregisters a connection. It does NOT close it.
func (cm *ConnectionManager) Remove(connID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.connections, connID)
}

// Get returns a connection by ID (nil if not found).
func (cm *ConnectionManager) Get(connID string) *ClientConnection {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.connections[connID]
}

// Count returns the total number of live connections.
func (cm *ConnectionManager) Count() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return len(cm.connections)
}

// CloseAll gracefully closes every active connection (used on server shutdown).
func (cm *ConnectionManager) CloseAll() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	for id, conn := range cm.connections {
		conn.Close()
		delete(cm.connections, id)
	}
}

// ByRoom returns all connections currently assigned to the given room.
func (cm *ConnectionManager) ByRoom(roomID string) []*ClientConnection {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	var out []*ClientConnection
	for _, c := range cm.connections {
		if c.RoomID == roomID {
			out = append(out, c)
		}
	}
	return out
}

// newClientConnection allocates a ClientConnection with a buffered Send channel.
// The caller must assign conn.Conn after calling this function.
func newClientConnection(id, userID, role, displayName, roomID, interviewID string) *ClientConnection {
	return &ClientConnection{
		ID:          id,
		UserID:      userID,
		Role:        role,
		DisplayName: displayName,
		RoomID:      roomID,
		InterviewID: interviewID,
		Send:        make(chan []byte, sendBufferSize),
		ConnectedAt: time.Now().UTC(),
	}
}
