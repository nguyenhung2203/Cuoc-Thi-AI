# H-S4-01 — AI Suggestion Event

| Field | Value |
|---|---|
| **Task ID** | H-S4-01 |
| **Sprint** | 4 — AI Realtime Bridge |
| **Độ khó** | Khó |
| **Trạng thái** | ⬜ Chưa bắt đầu |
| **Owner** | Hùng |

---

## Mục tiêu

Nhận AI suggestion từ AI Orchestrator và push cho **Recruiter only** qua WebSocket:
- Recruiter yêu cầu `ai:request_suggestion`
- Server gửi `ai:thinking` ngay
- AI Orchestrator trả suggestion → server gửi `ai:suggestion`
- Candidate **KHÔNG** nhận event này

---

## Yêu cầu chi tiết

### 1. `ai:request_suggestion` (Client → Server)

```json
{
  "event": "ai:request_suggestion",
  "payload": {
    "focus": "technical_depth",
    "last_transcript_id": "uuid"
  }
}
```

### 2. `ai:thinking` (Server → Recruiter)

```json
{
  "event": "ai:thinking",
  "payload": {
    "task": "suggest_follow_up",
    "request_id": "req_012",
    "message": "AI đang phân tích câu trả lời...",
    "expected_duration_ms": 3000
  }
}
```

### 3. `ai:suggestion` (Server → Recruiter)

```json
{
  "event": "ai:suggestion",
  "payload": {
    "suggestion_id": "uuid",
    "suggestion_type": "follow_up_question",
    "content": "Bạn có thể nói rõ bạn đã đo performance bằng chỉ số nào không?",
    "reason": "Ứng viên nói đã tối ưu performance nhưng chưa nêu metric cụ thể.",
    "target_skill": "Performance Optimization",
    "priority": "high",
    "confidence": 0.84
  }
}
```

### 4. Flow

```
Recruiter → ai:request_suggestion → Server → ai:thinking → Recruiter
                                         ↓
                                    AI Orchestrator
                                         ↓
                                    Server → ai:suggestion → Recruiter
```

### 5. Lưu DB

- Lưu suggestion vào bảng `ai_suggestions`
- Rate limit: max 10 request / 10 phút / interview

---

## Definition of Done

- [ ] Recruiter nhận gợi ý AI realtime
- [ ] `ai:thinking` gửi ngay khi bắt đầu xử lý
- [ ] Candidate **KHÔNG** nhận `ai:thinking` và `ai:suggestion`
- [ ] Suggestion lưu DB đúng
- [ ] Rate limit hoạt động
- [ ] AI lỗi → gửi `ai:error` thay vì crash

---

## Dependency

- **Sprint 3** — Transcript pipeline (cần transcript data cho AI)
- **Khôi** — AI Orchestrator suggest-follow-up endpoint

---

## Checklist test

- [ ] Recruiter request → nhận `ai:thinking` → nhận `ai:suggestion`
- [ ] Candidate → KHÔNG nhận event AI
- [ ] AI chậm > 15s → frontend tự handle timeout
- [ ] AI lỗi → nhận `ai:error`
- [ ] Spam request → bị rate limit
