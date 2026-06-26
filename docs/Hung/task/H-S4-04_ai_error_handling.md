# H-S4-04 — AI Error Event Handling

| Field | Value |
|---|---|
| **Task ID** | H-S4-04 |
| **Sprint** | 4 — AI Realtime Bridge |
| **Độ khó** | Trung bình |
| **Trạng thái** | ✅ Hoàn thành |
| **Owner** | Hùng |

---

## Mục tiêu

Handle AI errors gracefully — AI lỗi KHÔNG crash room:
- Gửi `ai:error` cho recruiter khi AI thực sự lỗi
- Phân biệt "chậm" vs "lỗi thực sự"
- Room tiếp tục hoạt động bình thường khi AI down

---

## Yêu cầu chi tiết

### `ai:error` Event

```json
{
  "event": "ai:error",
  "payload": {
    "error_type": "ai_service_unavailable",
    "request_id": "req_012",
    "affected_features": ["suggestion", "scoring"],
    "severity": "degraded",
    "message": "AI tạm thời không phản hồi. Buổi phỏng vấn vẫn tiếp tục.",
    "recoverable": true,
    "retry_after_seconds": 30
  }
}
```

### Error Types

| error_type | severity | Ý nghĩa |
|---|---|---|
| `ai_throttled` | degraded | AI bị throttle, đang retry |
| `ai_timeout` | degraded | Vượt hard timeout |
| `ai_service_unavailable` | critical | AI down hoàn toàn |
| `ai_rate_limited` | degraded | Vượt quota |
| `ai_invalid_input` | info | Input không hợp lệ |

### Fallback Behavior

- **degraded**: Room + chat + media vẫn hoạt động, AI panel chờ
- **critical**: Ẩn AI panel, room vẫn hoạt động
- **Kết thúc phỏng vấn**: Luôn cho phép dù AI lỗi

### Auto-retry

- `recoverable: true` → tự retry sau `retry_after_seconds`
- Max 3 retries → sau đó chuyển `critical`

---

## Definition of Done

- [x] AI lỗi → gửi `ai:error` cho recruiter
- [x] Room KHÔNG crash khi AI down
- [x] Video/audio/chat vẫn hoạt động khi AI lỗi
- [x] Auto-retry cho recoverable errors
- [x] Candidate KHÔNG nhận `ai:error`
- [x] End interview vẫn hoạt động khi AI lỗi

---

## Dependency

- **H-S4-01**, **H-S4-02** — Cần AI bridge hoạt động trước

---

## Checklist test

- [x] AI timeout → `ai:error` severity=degraded
- [x] AI down → `ai:error` severity=critical
- [x] AI lỗi → room chat/video vẫn hoạt động
- [x] recoverable → auto retry
- [x] End interview khi AI lỗi → vẫn thành công

### Bằng chứng test (Test Evidence)
```
=== RUN   TestAIErrorHandling
=== RUN   TestAIErrorHandling/AI_timeout_->_ai:error_severity=degraded_and_auto_retry
=== RUN   TestAIErrorHandling/AI_down_->_ai:error_severity=critical
=== RUN   TestAIErrorHandling/AI_lỗi_->_room_chat/video_vẫn_hoạt_động
=== RUN   TestAIErrorHandling/End_interview_khi_AI_lỗi_->_vẫn_thành_công
--- PASS: TestAIErrorHandling (0.41s)
    --- PASS: TestAIErrorHandling/AI_timeout_->_ai:error_severity=degraded_and_auto_retry (0.35s)
    --- PASS: TestAIErrorHandling/AI_down_->_ai:error_severity=critical (0.05s)
    --- PASS: TestAIErrorHandling/AI_lỗi_->_room_chat/video_vẫn_hoạt_động (0.00s)
    --- PASS: TestAIErrorHandling/End_interview_khi_AI_lỗi_->_vẫn_thành_công (0.00s)
PASS
```
