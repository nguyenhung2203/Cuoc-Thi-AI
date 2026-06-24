# H-S4-04 — AI Error Event Handling

| Field | Value |
|---|---|
| **Task ID** | H-S4-04 |
| **Sprint** | 4 — AI Realtime Bridge |
| **Độ khó** | Trung bình |
| **Trạng thái** | ⬜ Chưa bắt đầu |
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

- [ ] AI lỗi → gửi `ai:error` cho recruiter
- [ ] Room KHÔNG crash khi AI down
- [ ] Video/audio/chat vẫn hoạt động khi AI lỗi
- [ ] Auto-retry cho recoverable errors
- [ ] Candidate KHÔNG nhận `ai:error`
- [ ] End interview vẫn hoạt động khi AI lỗi

---

## Dependency

- **H-S4-01**, **H-S4-02** — Cần AI bridge hoạt động trước

---

## Checklist test

- [ ] AI timeout → `ai:error` severity=degraded
- [ ] AI down → `ai:error` severity=critical
- [ ] AI lỗi → room chat/video vẫn hoạt động
- [ ] recoverable → auto retry
- [ ] End interview khi AI lỗi → vẫn thành công
