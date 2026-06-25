# H-S3-04 — Save Transcript Final

| Field | Value |
|---|---|
| **Task ID** | H-S3-04 |
| **Sprint** | 3 — Transcript Realtime |
| **Độ khó** | Trung bình |
| **Trạng thái** | ⬜ Chưa bắt đầu |
| **Owner** | Hùng |

---

## Mục tiêu

Lưu transcript final vào database `interview_transcripts`:
- Chỉ lưu `is_final = true`
- Lưu đầy đủ: speaker, content, timing, confidence, source
- Data dùng cho report generation và AI scoring sau phỏng vấn

---

## Yêu cầu chi tiết

### Bảng `interview_transcripts`

| Field | Value example |
|---|---|
| `id` | UUID |
| `interview_id` | UUID |
| `participant_id` | UUID |
| `speaker_type` | "candidate" |
| `speaker_name` | "Trần Văn B" |
| `content` | "Tôi từng làm dự án..." |
| `language` | "vi" |
| `start_time_ms` | 120000 |
| `end_time_ms` | 128000 |
| `confidence` | 0.91 |
| `source` | "audio" |
| `is_final` | true |
| `created_at` | timestamp |

### Batch Insert

- Nếu transcript đến nhanh → batch insert mỗi 5-10 giây
- Không block WebSocket pipeline khi insert DB

---

## Definition of Done

- [ ] Transcript final lưu DB đúng schema
- [ ] Không lưu partial transcript vào DB
- [ ] Batch insert hiệu quả
- [ ] Data đủ cho AI scoring/report sử dụng
- [ ] Index `(interview_id, created_at)` hoạt động

---

## Dependency

- **H-S3-01** đến **H-S3-03** — Pipeline + speaker mapping
- **Khôi** — DB schema `interview_transcripts`

---

## Checklist test

- [ ] Final transcript lưu DB đúng field
- [ ] Partial transcript KHÔNG lưu DB
- [ ] 100 transcript items → batch insert nhanh
- [ ] Query by interview_id → kết quả đúng thứ tự
