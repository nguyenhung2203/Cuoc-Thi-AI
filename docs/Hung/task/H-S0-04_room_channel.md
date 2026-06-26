# H-S0-04 — Room Channel Architecture

| Field | Value |
|---|---|
| **Task ID** | H-S0-04 |
| **Sprint** | 0 — Setup Realtime Foundation |
| **Độ khó** | Khó |
| **Trạng thái** | ✅ Hoàn thành |
| **Owner** | Hùng |

---

## Mục tiêu

Thiết kế kiến trúc room channel để:
- Mỗi interview room là một channel riêng
- Broadcast event chỉ trong cùng room
- Gửi event riêng theo role (recruiter-only events)
- Map connection vào room khi join
- Remove connection khỏi room khi leave

---

## Yêu cầu chi tiết

### 1. Room Manager

```go
type RoomManager struct {
    rooms map[string]*Room  // room_id → Room
    mu    sync.RWMutex
}

type Room struct {
    ID           string
    InterviewID  string
    Status       string  // waiting/active/paused/completed
    Participants map[string]*Participant
    CreatedAt    time.Time
}

type Participant struct {
    ConnectionID   string
    UserID         string
    ParticipantType string  // recruiter/candidate/ai/guest
    DisplayName    string
    Connection     *ClientConnection
    JoinedAt       time.Time
}
```

### 2. Broadcast Functions

```go
// Gửi cho TẤT CẢ participant trong room
func (r *Room) BroadcastAll(event Event)

// Gửi cho participant CỤ THỂ
func (r *Room) SendTo(participantID string, event Event)

// Gửi cho TẤT CẢ recruiter trong room (AI events)
func (r *Room) BroadcastRecruitersOnly(event Event)

// Gửi cho TẤT CẢ trừ sender
func (r *Room) BroadcastExcept(senderID string, event Event)
```

### 3. Visibility Rules Engine

```go
func ShouldSendToParticipant(event string, participant *Participant) bool {
    // Dựa vào event type và participant role
    // Ví dụ: ai:suggestion → chỉ recruiter
    // Ví dụ: chat:message room → tất cả
    // Ví dụ: note:create → chỉ recruiter
}
```

### 4. Room Lifecycle

```
CreateRoom → AddParticipant → BroadcastEvents → RemoveParticipant → CloseRoom
```

---

## Tài liệu tham chiếu

- [REALTIME_EVENTS.md](../../REALTIME_EVENTS.md) — Mục 6 (Visibility rules)
- [DATABASE_DESIGN.md](../../DATABASE_DESIGN.md) — Bảng `interview_rooms`, `interview_participants`

---

## Definition of Done

- [x] Room được tạo khi có participant đầu tiên join
- [x] Participant được map đúng room
- [x] BroadcastAll gửi đúng cho tất cả trong room
- [x] BroadcastRecruitersOnly chỉ gửi cho recruiter
- [x] Event không bị gửi nhầm sang room khác
- [x] Remove participant khi leave/disconnect
- [x] Room cleanup khi không còn participant

---

## Dependency

- **H-S0-01** — Event types + visibility rules
- **H-S0-02** — Connection manager

---

## Files dự kiến tạo/sửa

```
backend/
├── internal/realtime/
│   ├── room.go                    # Room struct + methods
│   ├── room_manager.go            # Room pool management
│   ├── participant.go             # Participant struct
│   ├── broadcast.go               # Broadcast functions
│   └── visibility.go              # Cập nhật: visibility engine
```

---

## Checklist test

- [x] 2 client join cùng room → broadcast đúng cả 2
- [x] 2 client ở 2 room khác nhau → không nhận event lẫn nhau
- [x] Recruiter-only event → candidate không nhận
- [x] Client leave → bị remove khỏi room
- [x] Room trống → room được cleanup
- [x] Nhiều room đồng thời → không race condition
