# H-S1-01 — Join/Leave Room

| Field | Value |
|---|---|
| **Task ID** | H-S1-01 |
| **Sprint** | 1 — Interview Room MVP |
| **Độ khó** | Khó |
| **Trạng thái** | ⬜ Chưa bắt đầu |
| **Owner** | Hùng |

---

## Mục tiêu

Implement join/leave room flow hoàn chỉnh:
- Recruiter và Candidate gửi `room:join` → server validate và thêm vào room
- Server broadcast `room:user_joined` cho tất cả participant
- User gửi `room:leave` hoặc disconnect → server broadcast `room:user_left`
- Server trả `room:joined` với danh sách participant hiện tại

---

## Yêu cầu chi tiết

### 1. Client gửi `room:join`

```json
{
  "event": "room:join",
  "request_id": "req_001",
  "room_id": "uuid",
  "interview_id": "uuid",
  "payload": {
    "participant_type": "recruiter",
    "display_name": "Nguyễn Văn A",
    "device_info": { "browser": "Chrome", "os": "Windows" }
  }
}
```

### 2. Server xử lý Join

1. Validate `room_access_token` đã xác thực ở connection
2. Kiểm tra room tồn tại (tạo nếu chưa có)
3. Kiểm tra participant permission
4. Tạo/cập nhật participant trong room
5. Lưu `interview_participants` vào DB
6. Gửi `room:joined` cho sender (kèm danh sách participants hiện tại)
7. Broadcast `room:user_joined` cho tất cả participant khác

### 3. Server gửi `room:joined` (cho sender)

```json
{
  "event": "room:joined",
  "request_id": "req_001",
  "payload": {
    "participant_id": "uuid",
    "room_status": "waiting",
    "participants": [...]
  }
}
```

### 4. Server broadcast `room:user_joined`

```json
{
  "event": "room:user_joined",
  "payload": {
    "participant_id": "uuid",
    "display_name": "Nguyễn Văn A",
    "participant_type": "recruiter",
    "joined_at": "2026-06-23T09:00:00Z"
  }
}
```

### 5. Leave Room

- Client gửi `room:leave` hoặc WebSocket disconnect
- Server cập nhật `left_at` trong DB
- Server broadcast `room:user_left`

---

## Tài liệu tham chiếu

- [REALTIME_EVENTS.md](../../REALTIME_EVENTS.md) — Mục 4.1, 4.2, 5.1, 5.2, 5.3
- [DATABASE_DESIGN.md](../../DATABASE_DESIGN.md) — Bảng `interview_participants` (mục 15)

---

## Definition of Done

- [ ] Recruiter join room thành công
- [ ] Candidate join room thành công
- [ ] User không có quyền bị chặn
- [ ] Participant list trả về đúng khi join
- [ ] Broadcast `room:user_joined` cho participant khác
- [ ] Leave room cập nhật DB và broadcast
- [ ] Disconnect cũng trigger leave logic
- [ ] Không cho user lạ vào room

---

## Dependency

- **H-S0-01** đến **H-S0-04** — Toàn bộ Sprint 0
- **Khôi** — DB schema `interview_participants`, API interview

---

## Checklist test

- [ ] Recruiter join → nhận `room:joined` với participant list
- [ ] Candidate join → recruiter nhận `room:user_joined`
- [ ] Token sai room → bị reject
- [ ] Leave → tất cả nhận `room:user_left`
- [ ] Disconnect (close tab) → tự trigger leave
- [ ] Join lại sau leave → thành công
