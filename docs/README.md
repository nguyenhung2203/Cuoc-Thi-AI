# AI Interview Platform — Bộ tài liệu kỹ thuật cho AI IDE

## Tech Stack

- **Frontend**: Vue 3 + Pinia + TailwindCSS
- **Backend**: Golang (Core & Realtime)
- **AI Service**: Python (LLM orchestration, scoring, report)
- **Database**: PostgreSQL
- **Cache/Queue**: Redis
- **Realtime**: WebSocket
- **Auth**: JWT Access Token + HttpOnly Refresh Token (Token Family)

---

## Danh sách tài liệu

### Tài liệu chính (Source of Truth)

| # | File | Mô tả |
|---|---|---|
| 1 | `PRODUCT_REQUIREMENTS.md` | Hiểu sản phẩm, actor, MVP, flow, acceptance criteria |
| 2 | `API_SPEC.md` | Hợp đồng API — endpoint, request/response, permission |
| 3 | `DATABASE_DESIGN.md` | Thiết kế bảng, quan hệ, enum, migration, index |
| 4 | `REALTIME_EVENTS.md` | WebSocket events, room lifecycle, transcript, AI realtime |
| 5 | `AI_PROMPT_AND_SCORING.md` | Prompt template, rubric, scoring, report, guardrail AI |
| 6 | `GiaoDien.md` | Thiết kế giao diện, tông màu, layout, component |

### Tài liệu quy tắc

| # | File | Mô tả |
|---|---|---|
| 7 | `AI_INTERVIEW_PROJECT_RULES.md` | Quy tắc phát triển, branch, commit, test, DoD |
| 8 | `AI_INTERVIEW_PROJECT_SKILL.md` | Bộ skill hướng dẫn AI IDE khi code dự án |

### Tài liệu tham khảo

| # | File | Mô tả |
|---|---|---|
| 9 | `SYSTEM_ANALYSIS.md` | Phân tích hệ thống tổng hợp ban đầu (tham khảo, không phải source of truth) |

### Tài liệu module riêng

| # | File | Owner |
|---|---|---|
| 10 | `Hung/MODULE_HUNG_REALTIME_AI.md` | Hùng — Realtime, Interview Room, WebSocket, Transcript |
| 11 | `Khoi/MODULE_KHOI_CORE_AI_BACKEND.md` | Khôi — Backend core, DB, AI Engine, Scoring, Report |
| 12 | `Lai/MODULE_LAI_PRODUCT_UI_BUSINESS.md` | Lai — UI/UX, Dashboard, Portal, Business flow |

---

## Thứ tự AI IDE nên đọc

1. `PRODUCT_REQUIREMENTS.md` — luôn đọc trước.
2. `AI_INTERVIEW_PROJECT_RULES.md` — hiểu quy tắc chung.
3. File liên quan task hiện tại:
   - API/UI: `API_SPEC.md` + `GiaoDien.md`
   - DB/backend: `DATABASE_DESIGN.md`
   - Realtime room: `REALTIME_EVENTS.md`
   - AI/report/scoring: `AI_PROMPT_AND_SCORING.md`
4. File module của người đang làm:
   - Hùng: `Hung/MODULE_HUNG_REALTIME_AI.md`
   - Khôi: `Khoi/MODULE_KHOI_CORE_AI_BACKEND.md`
   - Lai: `Lai/MODULE_LAI_PRODUCT_UI_BUSINESS.md`

---

## Phân vai

- **Hùng**: đọc kỹ `REALTIME_EVENTS.md` + `AI_PROMPT_AND_SCORING.md` + `Hung/MODULE_HUNG_REALTIME_AI.md`.
- **Khôi**: đọc kỹ `DATABASE_DESIGN.md` + `API_SPEC.md` + `AI_PROMPT_AND_SCORING.md` + `Khoi/MODULE_KHOI_CORE_AI_BACKEND.md`.
- **Lai**: đọc kỹ `PRODUCT_REQUIREMENTS.md` + `API_SPEC.md` + `GiaoDien.md` + `Lai/MODULE_LAI_PRODUCT_UI_BUSINESS.md`.

---

## Quy tắc quan trọng

- Khi có mâu thuẫn giữa `SYSTEM_ANALYSIS.md` và các file chuyên biệt, **ưu tiên file chuyên biệt**.
- Không tự đổi tên endpoint, field, enum nếu chưa cập nhật tài liệu tương ứng.
- Không tự thêm tính năng ngoài MVP nếu không được giao.
