# H-S2-01 — Chat Realtime trong Room

| Field | Value |
|---|---|
| **Task ID** | H-S2-01 |
| **Sprint** | 2 — Chat + LiveKit WebRTC |
| **Độ khó** | Trung bình |
| **Trạng thái** | ✅ Hoàn thành |
| **Owner** | Hùng |

---

## Mục tiêu

Implement chat realtime trong phòng phỏng vấn:
- Recruiter và Candidate gửi tin nhắn qua WebSocket
- Tin nhắn broadcast realtime cho tất cả participant (hoặc chỉ recruiter)
- Lưu chat message vào DB
- Reload trang xem lại lịch sử chat

---

## Yêu cầu chi tiết

### 1. `chat:send` (Client → Server)

```json
{
  "event": "chat:send",
  "request_id": "req_007",
  "room_id": "uuid",
  "interview_id": "uuid",
  "payload": {
    "message": "Bạn có thể chia sẻ màn hình không?",
    "visibility": "room"
  }
}
```

**Visibility options:**
- `room` — tất cả participant thấy
- `recruiter_only` — chỉ recruiter thấy (nội bộ)

### 2. `chat:message` (Server → Client)

```json
{
  "event": "chat:message",
  "room_id": "uuid",
  "interview_id": "uuid",
  "payload": {
    "message_id": "uuid",
    "sender_participant_id": "uuid",
    "sender_name": "Nguyễn Văn A",
    "sender_type": "recruiter",
    "message": "Bạn có thể chia sẻ màn hình không?",
    "created_at": "2026-06-23T09:05:00Z"
  }
}
```

### 3. Lưu DB

- Lưu vào bảng riêng hoặc dùng `interview_transcripts` với `source = "chat"`
- Fields: interview_id, sender, content, visibility, created_at

### 4. Chat History API

- REST endpoint để load history khi reload trang (Khôi cung cấp hoặc Hùng tạo)

---

## Definition of Done

- [x] Gửi/nhận chat realtime
- [x] Chat `visibility: room` → tất cả nhận
- [x] Chat `visibility: recruiter_only` → candidate không nhận
- [x] Lưu chat vào DB
- [x] Reload vẫn xem được lịch sử
- [x] Không cho người ngoài room gửi chat

---

## Dependency

- **Sprint 0 + Sprint 1** hoàn thành

---

## Checklist test

- [x] Recruiter gửi chat → candidate nhận
- [x] Candidate gửi chat → recruiter nhận
- [x] `recruiter_only` chat → candidate không nhận
- [x] Message lưu DB đúng
- [x] Reload → lấy lại history
