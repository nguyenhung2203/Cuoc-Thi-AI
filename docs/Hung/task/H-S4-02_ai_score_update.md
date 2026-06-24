# H-S4-02 — AI Score Update Event

| Field | Value |
|---|---|
| **Task ID** | H-S4-02 |
| **Sprint** | 4 — AI Realtime Bridge |
| **Độ khó** | Khó |
| **Trạng thái** | ⬜ Chưa bắt đầu |
| **Owner** | Hùng |

---

## Mục tiêu

Nhận AI scoring update và push cho **Recruiter only**:
- Sau mỗi câu trả lời, AI cập nhật score theo rubric
- Score update broadcast cho recruiter trong room
- Candidate **KHÔNG** thấy scoring nội bộ

---

## Yêu cầu chi tiết

### `ai:score_update` (Server → Recruiter)

```json
{
  "event": "ai:score_update",
  "payload": {
    "scores": [
      {
        "criterion_name": "Technical Knowledge",
        "score": 4,
        "max_score": 5,
        "evidence": "Ứng viên mô tả được cách tối ưu query và cache.",
        "confidence": 0.78,
        "status": "scored"
      }
    ]
  }
}
```

### Trigger

- Tự động sau mỗi vài transcript final (AI Orchestrator tự trigger)
- Hoặc recruiter yêu cầu qua `ai:request_score_update`
- Lưu score vào DB `interview_scores`

---

## Definition of Done

- [ ] Recruiter nhận score update realtime
- [ ] Candidate **KHÔNG** nhận score
- [ ] Score lưu DB đúng
- [ ] Score có evidence + confidence
- [ ] `insufficient_evidence` khi thiếu data

---

## Dependency

- **H-S4-01** — AI suggestion bridge pattern
- **Khôi** — AI scoring endpoint

---

## Checklist test

- [ ] Score update → recruiter nhận đúng
- [ ] Candidate → KHÔNG nhận score event
- [ ] Score có evidence text
- [ ] Thiếu data → status = `insufficient_evidence`
