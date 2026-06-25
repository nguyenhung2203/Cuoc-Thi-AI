# H-S5-01 — Reconnect Recovery

| Field | Value |
|---|---|
| **Task ID** | H-S5-01 |
| **Sprint** | 5 — Hardening |
| **Độ khó** | Rất khó |
| **Trạng thái** | ⬜ Chưa bắt đầu |
| **Owner** | Hùng |

---

## Mục tiêu

Khi user mất kết nối ngắn rồi reconnect, khôi phục state đầy đủ:
- Giữ participant state trong room (không bị xóa ngay)
- Khi reconnect thành công → gửi current state (room status, participants, recent events)
- Client sync lại transcript/chat bằng REST hoặc event backlog

---

## Yêu cầu chi tiết

### 1. Mất mạng ngắn (< 2 phút)

- Server giữ participant `connection_state = "reconnecting"`
- KHÔNG broadcast `room:user_left` ngay
- Chờ reconnect trong grace period (2 phút)
- Khi reconnect → gửi `room:joined` lại với current state

### 2. Mất mạng lâu (> 2 phút)

- Chuyển `connection_state = "offline"`
- Broadcast `room:user_left` cho room
- User vẫn có thể join lại nếu interview chưa completed

### 3. Reconnect State Recovery

Khi user reconnect, server gửi:

```json
{
  "event": "room:joined",
  "payload": {
    "participant_id": "uuid",
    "room_status": "active",
    "participants": [...],
    "interview_status": "active",
    "media_status": {...},
    "missed_events_count": 5,
    "sync_from_timestamp": "2026-06-23T09:05:00Z"
  }
}
```

### 4. Event Backlog (Optional)

- Lưu N event gần nhất trong Redis (buffer 5 phút)
- Khi reconnect → gửi missed events cho client
- Hoặc client tự sync bằng REST (load transcript, chat history)

---

## Definition of Done

- [ ] Mất mạng < 2 phút → reconnect không mất state
- [ ] Mất mạng > 2 phút → user_left broadcast, vẫn join lại được
- [ ] Reconnect → nhận đúng room state hiện tại
- [ ] Không tạo duplicate participant khi reconnect
- [ ] Transcript không bị mất khi reconnect

---

## Dependency

- **H-S1-02** — Presence realtime
- **H-S1-04** — Room status sync

---

## Checklist test

- [ ] Disconnect 30s → reconnect → vẫn trong room
- [ ] Disconnect 30s → người kia KHÔNG nhận user_left
- [ ] Disconnect 3 phút → user_left broadcast → join lại thành công
- [ ] Reconnect → nhận room_status, participants đúng
- [ ] Reconnect → transcript panel sync lại
