# H-S3-01 — Transcript Event Pipeline

| Field | Value |
|---|---|
| **Task ID** | H-S3-01 |
| **Sprint** | 3 — Transcript Realtime |
| **Độ khó** | Rất khó |
| **Trạng thái** | ⬜ Chưa bắt đầu |
| **Owner** | Hùng |

---

## Mục tiêu

Xây dựng pipeline nhận transcript từ AI Orchestrator (STT) và broadcast vào room qua WebSocket.

---

## Yêu cầu chi tiết

### Pipeline

```
AI Orchestrator (Whisper STT) → Internal API/callback → WebSocket Gateway → broadcast transcript:update
```

### Event `transcript:update` (Server → Client)

```json
{
  "event": "transcript:update",
  "room_id": "uuid",
  "interview_id": "uuid",
  "payload": {
    "transcript_id": "uuid",
    "speaker_type": "candidate",
    "speaker_name": "Trần Văn B",
    "content": "Tôi từng làm dự án thương mại điện tử bằng Vue 3.",
    "start_time_ms": 120000,
    "end_time_ms": 128000,
    "confidence": 0.91,
    "is_final": true,
    "created_at": "2026-06-23T09:08:00Z"
  }
}
```

### Visibility

- Recruiter: luôn thấy transcript
- Candidate: tùy cấu hình room (`transcript_enabled` cho candidate)

### Internal API để nhận transcript

```
POST /internal/rooms/:room_id/transcript
```
- AI Orchestrator gọi endpoint này khi có transcript mới
- Gateway broadcast vào room

---

## Definition of Done

- [ ] Nhận transcript từ AI Orchestrator thành công
- [ ] Broadcast `transcript:update` vào đúng room
- [ ] Recruiter luôn nhận transcript
- [ ] Candidate nhận/không nhận tùy config
- [ ] Pipeline không block khi AI chậm
- [ ] Có buffer/queue nếu transcript đến quá nhanh

---

## Dependency

- **H-S2-03** — Audio stream hook
- **Khôi** — AI Orchestrator STT endpoint

---

## Checklist test

- [ ] AI gửi transcript → room nhận `transcript:update`
- [ ] Transcript gửi đúng room (không lẫn room khác)
- [ ] Candidate config off → không nhận transcript
- [ ] AI chậm → room vẫn hoạt động
