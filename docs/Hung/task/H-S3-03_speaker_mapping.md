# H-S3-03 — Speaker Mapping

| Field | Value |
|---|---|
| **Task ID** | H-S3-03 |
| **Sprint** | 3 — Transcript Realtime |
| **Độ khó** | Khó |
| **Trạng thái** | ⬜ Chưa bắt đầu |
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

- [ ] Transcript có `speaker_type` đúng (recruiter/candidate)
- [ ] Transcript có `speaker_name` đúng
- [ ] Không gán sai speaker giữa recruiter và candidate
- [ ] Unknown speaker được handle gracefully

---

## Dependency

- **H-S3-01** — Transcript pipeline
- **H-S2-02** — LiveKit integration (cần track identity)

---

## Checklist test

- [ ] Recruiter nói → transcript speaker_type = "recruiter"
- [ ] Candidate nói → transcript speaker_type = "candidate"
- [ ] 2 người nói xen kẽ → speaker đúng
- [ ] Không nhận diện được → "unknown" + log
