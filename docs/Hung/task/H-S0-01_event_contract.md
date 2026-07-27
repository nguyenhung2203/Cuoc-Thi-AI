# H-S0-01 — Thiết kế Event Contract Realtime

| Field | Value |
|---|---|
| **Task ID** | H-S0-01 |
| **Sprint** | 0 — Setup Realtime Foundation |
| **Độ khó** | Khó |
| **Trạng thái** | ✅ Hoàn thành |
| **Owner** | Hùng |

---

## Mục tiêu

Thiết kế và định nghĩa toàn bộ event contract cho WebSocket realtime gateway, bao gồm:
- Danh sách event client → server
- Danh sách event server → client
- Payload schema cho từng event
- Role visibility (Recruiter/Candidate)
- Event envelope chuẩn

---

## Yêu cầu chi tiết

### 1. Định nghĩa Event Envelope chuẩn

```json
{
  "event": "domain:action",
  "request_id": "uuid",
  "room_id": "uuid",
  "interview_id": "uuid",
  "sent_at": "ISO8601",
  "payload": {}
}
```

### 2. Client → Server Events cần định nghĩa

| Event | Mô tả | Permission |
|---|---|---|
| `room:join` | User vào phòng | All |
| `room:leave` | User rời phòng | All |
| `room:heartbeat` | Heartbeat 15-30s | All |
| `interview:start` | Bắt đầu phỏng vấn | Recruiter only |
| `interview:end` | Kết thúc phỏng vấn | Recruiter only |
| `media:status` | Cập nhật mic/camera | All |
| `chat:send` | Gửi tin nhắn | All |
| `note:create` | Tạo note nội bộ | Recruiter only |
| `question:mark_asked` | Đánh dấu câu hỏi đã hỏi | Recruiter only |
| `transcript:partial` | Transcript tạm (nếu client stream) | System |
| `transcript:final` | Transcript final | System |
| `ai:request_suggestion` | Yêu cầu AI gợi ý | Recruiter only |
| `ai:request_score_update` | Yêu cầu AI score | Recruiter only |

### 3. Server → Client Events cần định nghĩa

| Event | Visibility | Mô tả |
|---|---|---|
| `room:joined` | Sender | Xác nhận join thành công |
| `room:user_joined` | All | Broadcast có người vào |
| `room:user_left` | All | Broadcast có người rời |
| `room:presence_update` | All | Cập nhật danh sách online |
| `interview:started` | All | Phỏng vấn bắt đầu |
| `interview:completed` | All | Phỏng vấn kết thúc |
| `media:status_changed` | All | Trạng thái media thay đổi |
| `chat:message` | Theo visibility | Tin nhắn mới |
| `transcript:update` | Recruiter + tùy config | Transcript cập nhật |
| `ai:thinking` | Recruiter only | AI đang xử lý |
| `ai:suggestion` | Recruiter only | AI gợi ý câu hỏi |
| `ai:score_update` | Recruiter only | AI cập nhật điểm |
| `ai:warning` | Recruiter only | AI cảnh báo |
| `ai:error` | Recruiter only | AI lỗi |
| `report:ready` | Recruiter only | Report đã sẵn sàng |
| `error` | Sender | Lỗi chung |

### 4. Tạo file định nghĩa TypeScript/Go types

- Tạo file types cho event payload
- Tạo constants cho event names
- Tạo enum cho visibility rules

---

## Tài liệu tham chiếu

- [REALTIME_EVENTS.md](../../REALTIME_EVENTS.md) — Source of truth cho event spec
- [API_SPEC.md](../../API_SPEC.md) — Mục 7 (Room API)
- [DATABASE_DESIGN.md](../../DATABASE_DESIGN.md) — Bảng `interview_rooms`, `interview_participants`

---

## Definition of Done

- [x] Có danh sách đầy đủ event với payload schema
- [x] Có role visibility cho từng event
- [x] Có file type definitions (Go structs)
- [x] Có constants cho event names
- [x] Event naming theo convention `domain:action`
- [x] Payload thống nhất với `REALTIME_EVENTS.md`
- [ ] Đã review với Khôi (AI events) và Lai (frontend events)

---

## Dependency

- **Không phụ thuộc** task khác — đây là task đầu tiên
- Output của task này là input cho tất cả task còn lại

---

## Files dự kiến tạo/sửa

```
backend/
├── internal/realtime/
│   ├── events/
│   │   ├── event_types.go        # Event name constants
│   │   ├── event_envelope.go     # Envelope struct
│   │   ├── room_events.go        # Room event payloads
│   │   ├── interview_events.go   # Interview event payloads
│   │   ├── chat_events.go        # Chat event payloads
│   │   ├── transcript_events.go  # Transcript event payloads
│   │   ├── ai_events.go          # AI event payloads
│   │   └── visibility.go         # Role visibility rules
```
