# H-S0-02 — Tạo WebSocket Gateway Skeleton

| Field | Value |
|---|---|
| **Task ID** | H-S0-02 |
| **Sprint** | 0 — Setup Realtime Foundation |
| **Độ khó** | Khó |
| **Trạng thái** | ⬜ Chưa bắt đầu |
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

- [ ] Client connect WebSocket thành công
- [ ] Client disconnect được handle clean
- [ ] Server parse được event envelope JSON
- [ ] Server log connection/disconnect events
- [ ] Có connection pool quản lý đúng
- [ ] Có read/write goroutine per connection
- [ ] Không memory leak khi disconnect
- [ ] Có unit test cho connection manager

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

- [ ] Connect WebSocket bằng wscat hoặc Postman
- [ ] Gửi JSON message → server log nhận được
- [ ] Gửi invalid JSON → server không crash
- [ ] Disconnect → connection bị remove khỏi pool
- [ ] Nhiều connection cùng lúc → không race condition
- [ ] Server shutdown → tất cả connection đóng clean
