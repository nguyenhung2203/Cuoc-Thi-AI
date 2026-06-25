# H-S1-02 — Presence Realtime

| Field | Value |
|---|---|
| **Task ID** | H-S1-02 |
| **Sprint** | 1 — Interview Room MVP |
| **Độ khó** | Khó |
| **Trạng thái** | ⬜ Chưa bắt đầu |
| **Owner** | Hùng |

---

## Mục tiêu

Hiển thị trạng thái online/offline/reconnecting của từng participant trong room:
- Heartbeat mechanism để detect online/offline
- Cập nhật `last_seen_at` trong DB
- Broadcast presence update khi trạng thái thay đổi
- Phân biệt `online` / `offline` / `reconnecting` / `left`

---

## Yêu cầu chi tiết

### 1. Heartbeat

- Client gửi `room:heartbeat` mỗi **15-30 giây**
- Server cập nhật `last_seen_at` cho participant
- Nếu không nhận heartbeat > **60 giây** → chuyển sang `reconnecting`
- Nếu không nhận heartbeat > **2 phút** → chuyển sang `offline`

### 2. Connection States

| State | Ý nghĩa |
|---|---|
| `online` | Đang kết nối bình thường |
| `reconnecting` | Mất heartbeat 60s, có thể đang reconnect |
| `offline` | Mất heartbeat > 2 phút |
| `left` | User chủ động rời phòng |

### 3. Presence Update Event

```json
{
  "event": "room:presence_update",
  "room_id": "uuid",
  "payload": {
    "participants": [
      {
        "participant_id": "uuid",
        "display_name": "Nguyễn Văn A",
        "participant_type": "recruiter",
        "connection_state": "online",
        "last_seen_at": "2026-06-23T09:05:00Z"
      }
    ]
  }
}
```

### 4. Redis Presence Store

- Dùng Redis để lưu presence nhanh (TTL key)
- Key: `presence:{room_id}:{participant_id}`
- TTL: 90 giây (tự expire nếu không heartbeat)

---

## Tài liệu tham chiếu

- [REALTIME_EVENTS.md](../../REALTIME_EVENTS.md) — Mục 4.3, 5.4
- [DATABASE_DESIGN.md](../../DATABASE_DESIGN.md) — Bảng `interview_participants` (field `last_seen_at`, `connection_state`)

---

## Definition of Done

- [ ] Heartbeat cập nhật `last_seen_at` đúng
- [ ] Mất heartbeat → trạng thái chuyển `reconnecting` → `offline`
- [ ] Presence update broadcast đúng cho room
- [ ] Không nhân đôi participant khi reconnect
- [ ] Redis presence store hoạt động đúng TTL
- [ ] DB cập nhật `connection_state` đúng

---

## Dependency

- **H-S1-01** — Cần join/leave room hoạt động
- **Khôi** — Redis setup

---

## Checklist test

- [ ] Heartbeat → `last_seen_at` update
- [ ] Dừng heartbeat 60s → `reconnecting`
- [ ] Dừng heartbeat 2 phút → `offline`
- [ ] Gửi lại heartbeat → `online`
- [ ] 2 participant → presence list đúng 2 người
- [ ] Leave → state = `left`
