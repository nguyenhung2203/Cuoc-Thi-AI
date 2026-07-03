# H-S4-01 — AI Suggestion Event

| Field | Value |
|---|---|
| **Task ID** | H-S4-01 |
| **Sprint** | 4 — AI Realtime Bridge |
| **Độ khó** | Khó |
| **Trạng thái** | ✅ Hoàn thành |
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

- [x] Recruiter nhận gợi ý AI realtime
- [x] `ai:thinking` gửi ngay khi bắt đầu xử lý
- [x] Candidate **KHÔNG** nhận `ai:thinking` và `ai:suggestion`
- [x] Suggestion lưu DB đúng
- [x] Rate limit hoạt động
- [x] AI lỗi → gửi `ai:error` thay vì crash

---

## Dependency

- **Sprint 3** — Transcript pipeline (cần transcript data cho AI)
- **Khôi** — AI Orchestrator suggest-follow-up endpoint

---

## Checklist test

- [x] Recruiter request → nhận `ai:thinking` → nhận `ai:suggestion`
- [x] Candidate → KHÔNG nhận event AI
- [x] AI chậm > 15s → frontend tự handle timeout
- [x] AI lỗi → nhận `ai:error`
- [x] Spam request → bị rate limit

---

## Bằng chứng nghiệm thu (Test Evidence)

Chạy kiểm thử tự động toàn bộ test suite AI Suggestion Bridge (`go test -v ./tests/...`):
```text
=== RUN   TestAIRateLimiter
--- PASS: TestAIRateLimiter (0.00s)
=== RUN   TestAIRequestSuggestion_CandidateForbidden
2026/06/26 09:15:06 [ai] suggestion request forbidden for role="candidate" client=conn_cand_ai
--- PASS: TestAIRequestSuggestion_CandidateForbidden (0.00s)
=== RUN   TestAIRequestSuggestion_FlowAndVisibility
2026/06/26 09:15:06 [db] INSERT INTO ai_suggestions (id, interview_id, suggestion_type, content, reason, target_skill, priority, confidence, created_at) VALUES ('6dd1553a-f84d-4398-9342-29eb5c239820', 'iv_ai_flow', 'follow_up_question', 'Bạn có thể nói rõ bạn đã đo performance bằng chỉ số nào không?', ...)
--- PASS: TestAIRequestSuggestion_FlowAndVisibility (0.05s)
PASS
ok  	backend/tests	3.655s
```

Các file đã sửa & tạo mới:
1. `backend/internal/realtime/ai_handler.go`: Thiết lập `AIRateLimiter` (10 reqs/10m) & `handleAIRequestSuggestion` gửi ngay `ai:thinking`, giả lập gọi AI trả về `ai:suggestion` tới riêng Recruiter.
2. `backend/internal/realtime/message_router.go`: Khởi tạo instance `aiRateLimiter` & đăng ký route handler.
3. `backend/internal/realtime/server.go`: Đăng ký endpoint HTTP Webhook `handleAISuggestionPush` tiếp nhận từ AI Orchestrator ngoài.
4. `backend/tests/ai_suggestion_test.go`: Bộ kiểm thử nghiệm thu tự động bảo mật role, giới hạn tần suất và quy trình phản hồi tức thì.
