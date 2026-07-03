# H-S3-04 — Save Transcript Final

| Field | Value |
|---|---|
| **Task ID** | H-S3-04 |
| **Sprint** | 3 — Transcript Realtime |
| **Độ khó** | Trung bình |
| **Trạng thái** | ✅ Hoàn thành |
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

- [x] Transcript final lưu DB đúng schema
- [x] Không lưu partial transcript vào DB
- [x] Batch insert hiệu quả
- [x] Data đủ cho AI scoring/report sử dụng
- [x] Index `(interview_id, created_at)` hoạt động

---

## Dependency

- **H-S3-01** đến **H-S3-03** — Pipeline + speaker mapping
- **Khôi** — DB schema `interview_transcripts`

---

## Checklist test

- [x] Final transcript lưu DB đúng field
- [x] Partial transcript KHÔNG lưu DB
- [x] 100 transcript items → batch insert nhanh
- [x] Query by interview_id → kết quả đúng thứ tự

---

## Bằng chứng nghiệm thu (Test Evidence)

Chạy kiểm thử tự động toàn bộ test suite cơ chế Batch Saver (`go test -v ./tests/...`):
```text
=== RUN   TestTranscriptBatchSaver_FinalOnly
--- PASS: TestTranscriptBatchSaver_FinalOnly (0.00s)
=== RUN   TestTranscriptBatchSaver_FlushOnCapacity
--- PASS: TestTranscriptBatchSaver_FlushOnCapacity (0.00s)
=== RUN   TestTranscriptBatchSaver_FlushInterval
--- PASS: TestTranscriptBatchSaver_FlushInterval (0.25s)
=== RUN   TestTranscriptBatchSaver_QueryOrder
--- PASS: TestTranscriptBatchSaver_QueryOrder (0.00s)
PASS
ok  	backend/tests	3.076s
```

Các file đã sửa & tạo mới:
1. `backend/internal/realtime/transcript_saver.go`: Thiết lập struct `TranscriptRecord` & `TranscriptBatchSaver` async queue gom lô mỗi 5s hoặc đủ 100 items.
2. `backend/internal/realtime/message_router.go`: Khởi tạo instance `transcriptSaver` gắn vào bộ định tuyến tin nhắn.
3. `backend/internal/realtime/transcript_handler.go` & `chat_handler.go`: Gửi bản ghi final từ audio/chat vào queue.
4. `backend/internal/realtime/server.go`: Gửi bản ghi final từ AI webhook và đảm bảo `Flush()` an toàn khi Shutdown server.
5. `backend/tests/transcript_saver_test.go`: Unit test tự động nghiệm thu toàn bộ tiêu chí bài toán.
