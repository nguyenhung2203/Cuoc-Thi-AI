# Nghiệp vụ — AI Interview Platform

> Tài liệu đọc nhanh: sản phẩm làm gì, ai dùng, luồng chính, module.  
> Chi tiết kỹ thuật xem `PRODUCT_REQUIREMENTS.md`.

---

## 1. Sản phẩm là gì?

**AI Interview Platform** = nền tảng **phỏng vấn tuyển dụng trực tuyến có AI**.

| | |
|---|---|
| **Không phải** | Cuộc thi coding / contest |
| **Là** | Công cụ hỗ trợ nhà tuyển dụng + ứng viên trong phỏng vấn |
| **AI làm gì** | Transcript, gợi ý câu hỏi, chấm theo rubric, báo cáo, feedback mock |
| **AI không làm** | Tự reject ứng viên, quyết định tuyển thay người |

**Hai mục tiêu:**

1. Recruiter lọc & đánh giá ứng viên **nhanh hơn, nhất quán hơn**, có dữ liệu.
2. Candidate **luyện phỏng vấn** với AI và nhận feedback.

---

## 2. Ai dùng hệ thống?

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  RECRUITER  │     │  CANDIDATE  │     │    ADMIN    │
│  HR / NTĐ   │     │  Ứng viên   │     │  Quản trị   │
└──────┬──────┘     └──────┬──────┘     └──────┬──────┘
       │                   │                   │
       ▼                   ▼                   ▼
  Job, CV, lịch PV    Vào phòng PV        Users, AI
  Phòng + AI panel    Mock interview      Settings, logs
  Report + quyết định  Tìm việc / apply
```

| Vai trò | Được làm | Không được |
|---|---|---|
| **Recruiter** | Quản lý company, job, candidate, lịch PV, vào phòng, xem AI/scoring/report, quyết định cuối | Xem data company khác; để AI tự reject |
| **Candidate** | Vào phòng qua link, mic/camera/chat, mock interview, xem feedback của mình | Xem điểm/note/report nội bộ của recruiter |
| **Admin** | Quản lý users, companies, logs, cấu hình AI | — |
| **AI** | Phân tích JD/CV, gợi ý, chấm, report | Tự loại ứng viên; đánh giá theo tuổi/giới/ngoại hình/giọng |

---

## 3. Luồng nghiệp vụ chính

### 3.1. Recruiter — phỏng vấn thật (flow cốt lõi)

```
Đăng nhập
   → Tạo Company / Workspace
   → Tạo Job + JD
   → AI phân tích JD (rubric / câu hỏi)
   → Thêm Candidate + Upload CV
   → AI parse CV / fit score
   → Tạo lịch phỏng vấn + link mời
   → Vào phòng (Recruiter + Candidate)
   → Start → Transcript realtime + AI gợi ý / chấm
   → End → AI tạo Report
   → Recruiter xem report → Quyết định (pass / consider / reject)
```

### 3.2. Candidate — vào phòng

```
Mở link mời
   → Kiểm tra token
   → Xác nhận thông tin
   → Check mic / camera
   → Đồng ý ghi âm / AI (nếu bật)
   → Waiting room
   → Vào phòng khi recruiter Start
   → Trả lời câu hỏi
   → Màn hình kết thúc
```

### 3.3. Candidate — Mock interview (luyện tập)

```
Đăng nhập → Chọn vị trí / level → CV hoặc profile
   → AI sinh câu hỏi
   → Trả lời (text / audio)
   → AI feedback (điểm mạnh / yếu / gợi ý)
```

> Mock **không** dùng làm quyết định tuyển dụng thật.

---

## 4. Pipeline trạng thái (hiểu nhanh)

### Job
`draft` → `open` → `paused` → `closed`

### Candidate (trong pipeline)
`new` → `screening` → `invited` → `interviewing` → `passed` / `rejected` / `talent_pool`

### Interview
`scheduled` → `waiting` → `in_progress` → `completed` (hoặc `cancelled` / `expired`)

---

## 5. Module nghiệp vụ

| Module | Việc kinh doanh |
|---|---|
| **Auth & Account** | Đăng ký / đăng nhập; role admin, recruiter, candidate |
| **Company / Workspace** | Không gian làm việc; dữ liệu gắn theo company |
| **Job** | Vị trí tuyển + JD; AI phân tích JD |
| **Candidate + CV** | Hồ sơ ứng viên; upload CV; AI parse / tóm tắt |
| **Interview Schedule** | Lịch PV + link mời |
| **Interview Room** | Phòng realtime: A/V, chat, transcript, AI panel |
| **Rubric / Question Bank** | Tiêu chí chấm + ngân hàng câu hỏi |
| **Report** | Báo cáo sau PV + recommendation (recruiter quyết) |
| **Mock Interview** | Luyện PV với AI |
| **Job Board / Careers** | Trang việc làm công khai, apply *(mở rộng)* |
| **Admin / AI Settings** | Cấu hình model, prompt, giám sát hệ thống |
| **Audit & Notification** | Log hành động quan trọng; thông báo |

---

## 6. Nguyên tắc nghiệp vụ quan trọng

1. **AI hỗ trợ, người quyết định** — không auto-reject.
2. **Có evidence** — điểm / kết luận AI phải dựa trên JD, CV, transcript.
3. **Không bias** — không đánh giá theo ngoại hình, giọng, tuổi, giới tính…
4. **Phân quyền rõ** — candidate không thấy scoring/note/report nội bộ.
5. **Multi-tenant** — job / candidate / interview luôn thuộc một company.
6. **Không lưu video** — chỉ transcript / audio theo scope sản phẩm.

---

## 7. Thực thể dữ liệu (business)

| Entity | Ý nghĩa |
|---|---|
| **User** | Tài khoản (admin / recruiter / candidate) |
| **Company** | Workspace tuyển dụng |
| **Job** | Vị trí + JD |
| **Candidate** | Hồ sơ ứng viên thuộc company |
| **Interview** | Buổi PV thật hoặc mock |
| **Room / Transcript** | Phòng realtime + lời nói đã ghi |
| **Score / Rubric** | Điểm theo tiêu chí |
| **Report** | Báo cáo + recommendation sau PV |
| **File / CV** | File đính kèm (metadata, không nhét binary vào DB) |

---

## 8. Stack (tóm tắt)

| Layer | Công nghệ |
|---|---|
| Frontend | Vue 3 + Tailwind + Pinia + LiveKit |
| Backend API | Go |
| Realtime | Go WebSocket |
| AI | Python FastAPI + Gemini |
| DB / Cache | PostgreSQL + Redis |
| Media | LiveKit |
| Deploy | Docker + Caddy |

---

## 9. Team & phạm vi

| Thành viên | Phạm vi |
|---|---|
| **Hùng** | Realtime, Room, WebSocket, Transcript |
| **Khôi** | Backend, DB, AI Engine, Scoring, Report |
| **Lai** | UI/UX, Portal, Business flow |

---

## 10. Đọc tiếp khi cần sâu hơn

| Muốn hiểu… | Đọc file |
|---|---|
| Scope MVP, acceptance | `PRODUCT_REQUIREMENTS.md` |
| API | `API_SPEC.md` |
| Bảng DB | `DATABASE_DESIGN.md` |
| Event phòng realtime | `REALTIME_EVENTS.md` |
| Prompt / chấm điểm AI | `AI_PROMPT_AND_SCORING.md` |
| Demo tuyển dụng | `DEMO_SCRIPT_TUYEN_DUNG.md` |
