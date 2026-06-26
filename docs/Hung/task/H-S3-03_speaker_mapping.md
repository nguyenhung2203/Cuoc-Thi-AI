# H-S3-03 — Speaker Mapping

| Field | Value |
|---|---|
| **Task ID** | H-S3-03 |
| **Sprint** | 3 — Transcript Realtime |
| **Độ khó** | Khó |
| **Trạng thái** | ✅ Hoàn thành |
| **Owner** | Hùng |

---

## Mục tiêu

Map audio transcript với đúng speaker (recruiter/candidate):
- Biết ai đang nói dựa trên audio track source
- Gắn `speaker_type` và `speaker_name` vào transcript
- Phân biệt rõ giữa recruiter, candidate, system

---

## Yêu cầu chi tiết

### Speaker Types

| Type | Mô tả |
|---|---|
| `recruiter` | Người phỏng vấn |
| `candidate` | Ứng viên |
| `ai` | AI system message |
| `system` | Thông báo hệ thống |

### Mapping Strategy

1. **LiveKit track identity** → mỗi participant có identity riêng khi publish audio
2. **Audio track → participant_id** → map qua LiveKit room participant list
3. **participant_id → speaker_type + speaker_name** → map qua room participant data

### Fallback

- Nếu không map được speaker → `speaker_type: "unknown"`, log warning
- AI Orchestrator có thể dùng speaker diarization nếu cần

---

## Definition of Done

- [x] Transcript có `speaker_type` đúng (recruiter/candidate)
- [x] Transcript có `speaker_name` đúng
- [x] Không gán sai speaker giữa recruiter và candidate
- [x] Unknown speaker được handle gracefully

---

## Dependency

- **H-S3-01** — Transcript pipeline
- **H-S2-02** — LiveKit integration (cần track identity)

---

## Checklist test

- [x] Recruiter nói → transcript speaker_type = "recruiter"
- [x] Candidate nói → transcript speaker_type = "candidate"
- [x] 2 người nói xen kẽ → speaker đúng
- [x] Không nhận diện được → "unknown" + log

---

## Bằng chứng nghiệm thu (Test Evidence)

Chạy kiểm thử tự động toàn bộ module Realtime Gateway và Speaker Mapping (`go test -v ./tests/...`):
```text
=== RUN   TestResolveSpeaker
2026/06/26 09:03:21 [room-mgr] created room=room_123 interview=iv_123 defaulting to waiting
--- PASS: TestResolveSpeaker (0.00s)
=== RUN   TestResolveSpeakerFallback
2026/06/26 09:03:21 [room-mgr] created room=room_empty interview=iv_empty defaulting to waiting
--- PASS: TestResolveSpeakerFallback (0.00s)
=== RUN   TestTranscriptRoleVisibility
--- PASS: TestTranscriptRoleVisibility (0.00s)
PASS
ok  	backend/tests	3.627s
```

Các file đã sửa & tạo mới:
1. `backend/internal/realtime/events/transcript_events.go`: Thêm `participant_id`, `track_id`, `identity`, `speaker_name` vào struct payload.
2. `backend/internal/realtime/room.go`: Thêm `TrackParticipantMap`, `RegisterTrack`, `ResolveSpeaker`.
3. `backend/internal/realtime/transcript_handler.go`: Phân giải speaker trước khi broadcast và insert DB.
4. `backend/internal/realtime/server.go`: Phân giải speaker cho Webhook push `handleTranscriptPush`.
5. `backend/tests/transcript_test.go`: Unit test xác minh kiểm thử tự động.
