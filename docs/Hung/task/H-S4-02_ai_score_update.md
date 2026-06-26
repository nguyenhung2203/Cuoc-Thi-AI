# H-S4-02 — AI Score Update Event

| Field | Value |
|---|---|
| **Task ID** | H-S4-02 |
| **Sprint** | 4 — AI Realtime Bridge |
| **Độ khó** | Khó |
| **Trạng thái** | ✅ Hoàn thành |
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

- [x] Recruiter nhận score update realtime
- [x] Candidate **KHÔNG** nhận score
- [x] Score lưu DB đúng
- [x] Score có evidence + confidence
- [x] `insufficient_evidence` khi thiếu data

---

## Dependency

- **H-S4-01** — AI suggestion bridge pattern
- **Khôi** — AI scoring endpoint

---

## Checklist test

- [x] Score update → recruiter nhận đúng
- [x] Candidate → KHÔNG nhận score event
- [x] Score có evidence text
- [x] Thiếu data → status = `insufficient_evidence`

---

## Bằng chứng nghiệm thu (Test Evidence)

Chạy kiểm thử tự động toàn bộ test suite AI Score Update Bridge (`go test -v ./tests/...`):
```text
=== RUN   TestAIScoreUpdate_FlowAndEvidence
2026/06/26 09:19:19 [db] INSERT INTO interview_scores (id, interview_id, criterion_name, score, max_score, evidence, confidence, status, created_at) VALUES ('d9b55ab6-873d-4a8e-ac53-ae083df2e38f', 'iv_score_flow', 'Technical Knowledge', 4.000000, 5.000000, 'Ứng viên mô tả được cách tối ưu query và cache.', 0.780000, 'scored', ...)
--- PASS: TestAIScoreUpdate_FlowAndEvidence (0.05s)
=== RUN   TestAIScoreUpdate_InsufficientEvidence
2026/06/26 09:19:19 [db] INSERT INTO interview_scores (id, interview_id, criterion_name, score, max_score, evidence, confidence, status, created_at) VALUES ('9fe087d5-85c3-4182-b938-768839feb604', 'iv_score_insuff', 'Technical Depth', 0.000000, 5.000000, 'Câu trả lời quá ngắn, không đủ bằng chứng đánh giá.', 0.200000, 'insufficient_evidence', ...)
--- PASS: TestAIScoreUpdate_InsufficientEvidence (0.05s)
PASS
ok  	backend/tests	3.627s
```

Các file đã sửa & tạo mới:
1. `backend/internal/realtime/ai_handler.go`: Bổ sung `handleAIRequestScoreUpdate` chấm điểm Rubric kèm bằng chứng, tự động trả `insufficient_evidence` khi thiếu thông tin.
2. `backend/internal/realtime/message_router.go`: Khởi tạo `scoreRateLimiter` (10 reqs/10m) & định tuyến websocket event.
3. `backend/internal/realtime/server.go`: Thiết lập HTTP Webhook `handleAIScoreUpdatePush` tiếp nhận điểm chấm tự động từ Python AI Orchestrator.
4. `backend/tests/ai_score_update_test.go`: Unit test tự động nghiệm thu bảo mật role ứng viên, kiểm chứng evidence text và xử lý ngoại lệ thiếu dữ liệu.
