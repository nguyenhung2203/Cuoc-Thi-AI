# H-S1-03 — Start/End Interview Event

| Field | Value |
|---|---|
| **Task ID** | H-S1-03 |
| **Sprint** | 1 — Interview Room MVP |
| **Độ khó** | Trung bình |
| **Trạng thái** | ⬜ Chưa bắt đầu |
| **Owner** | Hùng |

---

## Mục tiêu

Implement start/end interview qua WebSocket:
- Chỉ **Recruiter** được start/end
- Cập nhật interview status trong DB
- Broadcast trạng thái cho tất cả participant
- End interview trigger report generation

---

## Yêu cầu chi tiết

### 1. `interview:start` (Client → Server)

```json
{
  "event": "interview:start",
  "request_id": "req_004",
  "room_id": "uuid",
  "interview_id": "uuid",
  "payload": {
    "consent_recording": true,
    "consent_ai": true
  }
}
```

**Server xử lý:**
1. Kiểm tra sender là recruiter assigned hoặc admin
2. Kiểm tra interview status = `scheduled` hoặc `waiting`
3. Cập nhật interview `status = active`, `started_at = now()`
4. Cập nhật room `status = active`
5. Broadcast `interview:started`
6. Nếu consent_ai = true → khởi động AI pipeline (thông báo cho AI Orchestrator)

### 2. `interview:end` (Client → Server)

```json
{
  "event": "interview:end",
  "request_id": "req_005",
  "room_id": "uuid",
  "interview_id": "uuid",
  "payload": {
    "generate_report": true
  }
}
```

**Server xử lý:**
1. Kiểm tra sender là recruiter
2. Cập nhật interview `status = completed`, `ended_at = now()`
3. Cập nhật room `status = closed`
4. Broadcast `interview:completed`
5. Nếu `generate_report = true` → trigger report generation (gọi API Khôi)
6. Close room sau grace period (30 giây cho client cleanup)

### 3. Idempotency

- `interview:end` gửi nhiều lần chỉ xử lý **1 lần**
- Dùng `request_id` để dedup

---

## Tài liệu tham chiếu

- [REALTIME_EVENTS.md](../../REALTIME_EVENTS.md) — Mục 4.4, 4.5, 5.5, 5.6
- [API_SPEC.md](../../API_SPEC.md) — Mục 6.4, 6.5 (Start/End Interview)
- [DATABASE_DESIGN.md](../../DATABASE_DESIGN.md) — Bảng `interviews` (status, started_at, ended_at)

---

## Definition of Done

- [ ] Recruiter start interview → status = active, broadcast `interview:started`
- [ ] Candidate KHÔNG start được → error event
- [ ] Recruiter end interview → status = completed, broadcast `interview:completed`
- [ ] End interview trigger report generation (nếu `generate_report = true`)
- [ ] `interview:end` idempotent
- [ ] Room close sau grace period
- [ ] DB cập nhật đúng `started_at`, `ended_at`

---

## Dependency

- **H-S1-01** — Join/leave room
- **Khôi** — Interview DB update API, Report generation trigger

---

## Checklist test

- [ ] Recruiter start → tất cả nhận `interview:started`
- [ ] Candidate start → nhận error
- [ ] Recruiter end → tất cả nhận `interview:completed`
- [ ] End 2 lần → chỉ xử lý 1 lần
- [ ] Start khi đã active → error
- [ ] DB `started_at` / `ended_at` đúng
