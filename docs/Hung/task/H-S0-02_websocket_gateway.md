# H-S0-02 — Tạo WebSocket Gateway Skeleton

| Field | Value |
|---|---|
| **Task ID** | H-S0-02 |
| **Sprint** | 0 — Setup Realtime Foundation |
| **Độ khó** | Khó |
| **Trạng thái** | ✅ Hoàn thành |
| **Owner** | Hùng |

---

## Mục tiêu

Tạo WebSocket gateway cơ bản bằng Golang, có thể:
- Chấp nhận WebSocket connection
- Quản lý connection pool
- Gửi/nhận message theo event envelope chuẩn
- Handle disconnect gracefully
- Logging connection events

---

## Yêu cầu chi tiết

### 1. WebSocket Server

- URL: `ws://localhost:8080/ws/interview-room`
- Production: `wss://api.example.com/ws/interview-room`
- Dùng thư viện: `gorilla/websocket` hoặc `nhooyr.io/websocket`
- Upgrade HTTP → WebSocket

### 2. Connection Manager

```go
type ConnectionManager struct {
    connections map[string]*ClientConnection  // connection_id → connection
    mu          sync.RWMutex
}

type ClientConnection struct {
    ID          string
    UserID      string
    Conn        *websocket.Conn
    Send        chan []byte
    RoomID      string
    Role        string  // recruiter/candidate
    ConnectedAt time.Time
}
```

### 3. Message Handler Pipeline

```
Client → WebSocket → Parse Envelope → Route by Event → Handler → Response/Broadcast
```

- Parse JSON envelope
- Validate event name
- Route đến handler tương ứng
- Handler xử lý logic
- Gửi response hoặc broadcast

### 4. Graceful Shutdown

- Close tất cả connections khi server shutdown
- Gửi close frame cho client
- Timeout cho pending messages

---

## Tài liệu tham chiếu

- [REALTIME_EVENTS.md](../../REALTIME_EVENTS.md) — Mục 2 (Connection)
- Task H-S0-01 output (event types)

---

## Definition of Done

- [x] Client connect WebSocket thành công
- [x] Client disconnect được handle clean
- [x] Server parse được event envelope JSON
- [x] Server log connection/disconnect events
- [x] Có connection pool quản lý đúng
- [x] Có read/write goroutine per connection
- [x] Không memory leak khi disconnect
- [x] Có unit test cho connection manager (6 tests PASS)

---

## Dependency

- **H-S0-01** — Cần event type definitions

---

## Files dự kiến tạo/sửa

```
backend/
├── cmd/realtime/
│   └── main.go                    # Entry point cho realtime server
├── internal/realtime/
│   ├── server.go                  # WebSocket server setup
│   ├── handler.go                 # HTTP upgrade handler
│   ├── connection.go              # ClientConnection struct
│   ├── connection_manager.go      # Connection pool
│   ├── message_router.go          # Event routing
│   └── config.go                  # Realtime config
```

---

## Checklist test

- [x] Connect WebSocket bằng wscat hoặc Postman (TestWS_Connect PASS)
- [x] Gửi JSON message → server log nhận được (TestWS_SendValidJSON PASS)
- [x] Gửi invalid JSON → server không crash (TestWS_InvalidJSON PASS)
- [x] Disconnect → connection bị remove khỏi pool (TestWS_DisconnectRemovesFromPool PASS)
- [x] Nhiều connection cùng lúc → không race condition (TestWS_MultipleConcurrentConnections PASS)
- [x] Server shutdown → tất cả connection đóng clean (TestWS_ServerShutdown PASS)
