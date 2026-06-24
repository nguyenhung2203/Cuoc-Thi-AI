# H-S3-02 — Partial/Final Transcript Handling

| Field | Value |
|---|---|
| **Task ID** | H-S3-02 |
| **Sprint** | 3 — Transcript Realtime |
| **Độ khó** | Khó |
| **Trạng thái** | ⬜ Chưa bắt đầu |
| **Owner** | Hùng |

---

## Mục tiêu

Xử lý transcript partial (đang nói) và final (đã hoàn thành):
- Partial transcript hiển thị realtime khi đang nói (có thể thay đổi)
- Final transcript là kết quả cuối cùng, lưu DB
- Không duplicate final text trên UI

---

## Yêu cầu chi tiết

### Partial Transcript

- `is_final: false` — nội dung đang thay đổi
- Client hiển thị ở dạng "đang gõ..." với text tạm
- Partial mới thay thế partial cũ của cùng speaker

### Final Transcript

- `is_final: true` — nội dung cuối cùng
- Client thêm vào danh sách transcript cố định
- Lưu vào DB (task H-S3-04)

### Dedup Logic

- Dùng `start_time_ms` + `speaker_type` để match partial → final
- Khi nhận final → xóa partial tương ứng trên UI
- Không hiển thị cùng lúc partial và final cho cùng đoạn

---

## Definition of Done

- [ ] Partial transcript hiển thị realtime
- [ ] Final transcript thay thế partial tương ứng
- [ ] Không duplicate final text
- [ ] Confidence thấp → có indicator trên UI

---

## Dependency

- **H-S3-01** — Transcript pipeline

---

## Checklist test

- [ ] Nhận partial → hiển thị "đang nói..."
- [ ] Nhận final → thay thế partial, hiển thị cố định
- [ ] 2 partial liên tiếp → text cập nhật, không nhân đôi
- [ ] confidence < 0.5 → hiển thị cảnh báo
