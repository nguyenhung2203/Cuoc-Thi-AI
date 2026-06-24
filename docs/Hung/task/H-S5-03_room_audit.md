# H-S5-03 — Room Event Audit Log

| Field | Value |
|---|---|
| **Task ID** | H-S5-03 |
| **Sprint** | 5 — Hardening |
| **Độ khó** | Trung bình |
| **Trạng thái** | ⬜ Chưa bắt đầu |
| **Owner** | Hùng |

---

## Mục tiêu

Ghi audit log cho các sự kiện quan trọng trong room:
- Ai join/leave room
- Khi nào start/end interview
- AI events quan trọng
- Connection/reconnect events

---

## Yêu cầu chi tiết

### Events cần ghi audit log

| Event | Log content |
|---|---|
| `room:join` | User X joined room Y |
| `room:leave` | User X left room Y |
| `interview:start` | Interview Z started by recruiter X |
| `interview:end` | Interview Z ended by recruiter X |
| `ai:suggestion` | AI suggestion generated for interview Z |
| `ai:score_update` | AI score updated for interview Z |
| `ai:error` | AI error in interview Z |
| Connection lost | User X disconnected from room Y |
| Reconnect | User X reconnected to room Y |

### Audit Log Format

```json
{
  "action": "room_join",
  "actor_id": "uuid",
  "actor_type": "recruiter",
  "resource_type": "interview_room",
  "resource_id": "uuid",
  "company_id": "uuid",
  "metadata": { "room_id": "uuid" },
  "ip_address": "1.2.3.4",
  "created_at": "timestamp"
}
```

### Lưu trữ

- Lưu vào bảng `audit_logs` (Khôi cung cấp schema)
- Async write — không block room events

---

## Definition of Done

- [ ] Ghi log đúng cho tất cả event quan trọng
- [ ] Audit log không block realtime events
- [ ] Log format đúng schema `audit_logs`
- [ ] Có actor_id, resource_id, timestamp
- [ ] Async write (goroutine hoặc channel buffer)

---

## Dependency

- **Khôi** — DB schema `audit_logs`
- **Sprint 1-4** — Các event handlers

---

## Checklist test

- [ ] Join room → audit log có record
- [ ] Start interview → audit log có record
- [ ] End interview → audit log có record
- [ ] Audit log không làm chậm room events
- [ ] Query audit log by interview_id → kết quả đúng
