# AI_INTERVIEW_PROJECT_SKILL.md
# Bộ Skill riêng cho dự án: Phỏng vấn cùng AI Real-time

## 1. Vai trò của bộ skill

Bộ skill này dùng để hướng dẫn AI/Developer/Team khi làm việc trên dự án **nền tảng phỏng vấn cùng AI real-time**.

Mục tiêu của bộ skill:

- Giúp AI hiểu đúng sản phẩm trước khi code.
- Giúp team thống nhất cách phân tích, thiết kế, triển khai.
- Giúp chia module rõ ràng cho 3 người: **Hùng, Khôi, Lai**.
- Giúp tránh code lệch hướng, chồng chéo module hoặc phá kiến trúc.
- Giúp AI khi nhận prompt có thể tự kiểm tra requirement, dependency, security, realtime và AI flow.

Dự án gồm 2 actor chính:

| Actor | Vai trò |
|---|---|
| Tuyển dụng / Recruiter | Tạo job, mời ứng viên, phỏng vấn, xem AI report, ra quyết định tuyển dụng |
| Ứng viên / Candidate | Tham gia phỏng vấn thật hoặc phỏng vấn thử với AI |

AI trong hệ thống là **trợ lý hỗ trợ**, không phải người quyết định cuối cùng.

---

## 2. Tư duy sản phẩm bắt buộc

Khi phát triển bất kỳ tính năng nào, phải luôn đặt câu hỏi:

1. Tính năng này phục vụ **Recruiter**, **Candidate**, hay cả hai?
2. Dữ liệu này có nhạy cảm không?
3. AI có được phép xem dữ liệu này không?
4. Kết quả AI là gợi ý hay quyết định cuối cùng?
5. Có cần lưu audit log không?
6. Có ảnh hưởng đến realtime room không?
7. Có cần xử lý khi mất mạng, reconnect, AI lỗi, STT lỗi không?
8. Có thể mở rộng cho nhiều công ty/workspace sau này không?

Nguyên tắc quan trọng:

> AI chỉ được hỗ trợ phân tích, gợi ý, tóm tắt và chấm điểm tham khảo. Recruiter vẫn là người ra quyết định cuối cùng.

---

## 3. Ngôn ngữ nghiệp vụ chính

| Thuật ngữ | Ý nghĩa |
|---|---|
| Job | Vị trí tuyển dụng |
| Candidate | Ứng viên |
| Recruiter | Người tuyển dụng / người phỏng vấn |
| Interview | Buổi phỏng vấn |
| Interview Room | Phòng phỏng vấn online real-time |
| Mock Interview | Phỏng vấn thử với AI |
| Transcript | Nội dung hội thoại được chuyển thành văn bản |
| Rubric | Bộ tiêu chí chấm điểm |
| AI Suggestion | Gợi ý câu hỏi hoặc hướng phỏng vấn từ AI |
| AI Score | Điểm AI đánh giá theo tiêu chí |
| Report | Báo cáo sau phỏng vấn |
| Evidence | Bằng chứng trích từ câu trả lời của ứng viên |
| Workspace / Company | Không gian làm việc của công ty tuyển dụng |

---

## 4. Kiến trúc tư duy tổng quát

Hệ thống nên được chia thành các khối lớn:

```text
Frontend Web App (Vue 3 + Pinia + TailwindCSS)
    |
Backend API (Golang + REST)
    |
Realtime Gateway / WebSocket (Golang)
    |
Interview Room Service (Golang)
    |
AI Orchestrator (Python)
    |-- Speech-to-Text Service (Whisper)
    |-- LLM Service (Gemini)
    |-- Scoring Service
    |-- Report Generator
    |-- Prompt Template Manager
    |
PostgreSQL / S3 / Redis
```

### 4.1. Frontend

Frontend cần phục vụ 2 luồng:

- Recruiter portal.
- Candidate portal.

Các màn hình quan trọng:

- Login/Register.
- Recruiter Dashboard.
- Job List / Job Detail.
- Candidate List / Candidate Detail.
- Interview Schedule.
- Interview Room.
- Interview Report.
- Mock Interview.
- Candidate Feedback.

### 4.2. Backend

Backend chịu trách nhiệm:

- Auth.
- Company/Workspace.
- Job.
- Candidate.
- Interview.
- Question Bank.
- Rubric.
- Report.
- AI request orchestration.
- Permission/RBAC.
- Audit log.

### 4.3. Realtime

Realtime chịu trách nhiệm:

- Join/leave room.
- Presence.
- Audio/video signaling.
- Chat realtime.
- Transcript realtime.
- AI suggestion realtime.
- AI scoring update realtime.
- Reconnect.
- Room state.

### 4.4. AI

AI chịu trách nhiệm:

- Phân tích JD.
- Phân tích CV.
- Sinh câu hỏi.
- Gợi ý follow-up.
- Chấm điểm câu trả lời.
- Tóm tắt transcript.
- Tạo report.
- Tạo mock interview.
- Feedback cho ứng viên.

---

## 5. Phân công theo độ khó

| Người | Vai trò chính | Mức độ module |
|---|---|---|
| Hùng | Realtime + Interview Room + Transcript + WebSocket/WebRTC | Khó nhất |
| Khôi | AI Engine + Core Backend + Scoring + Report + Data model phức tạp | Khó nhất |
| Lai | Product UI + Business module + Candidate/Recruiter portal + các module CRUD còn lại | Trung bình / nhiều việc / cần ổn định |

Nguyên tắc chia việc:

- **Hùng** xử lý các phần cần realtime, đồng bộ trạng thái, room, transcript live.
- **Khôi** xử lý các phần lõi backend, AI orchestration, dữ liệu, scoring, report.
- **Lai** xử lý phần giao diện, nghiệp vụ tuyển dụng, CRUD, dashboard, mock interview portal cơ bản.

---

## 6. Skill dành cho AI khi hỗ trợ dự án

Khi AI được yêu cầu code hoặc phân tích, AI phải làm theo quy trình sau:

### Bước 1: Đọc bối cảnh

AI cần xác định:

- Đang làm module của ai?
- Module thuộc Sprint nào?
- Có dependency với module khác không?
- Có ảnh hưởng database không?
- Có ảnh hưởng API contract không?
- Có ảnh hưởng realtime event không?
- Có ảnh hưởng bảo mật dữ liệu không?

### Bước 2: Phân tích trước khi code

Trước khi sửa code, AI phải trả lời được:

- File nào cần sửa?
- File nào không được đụng?
- API nào được thêm?
- Model/schema nào cần thêm?
- Có cần migration không?
- Có cần test không?
- Có breaking change không?

### Bước 3: Code theo module ownership

AI không được tự ý sửa module người khác nếu không cần thiết.

Ví dụ:

- Đang làm module của Lai thì không tự ý sửa AI scoring của Khôi.
- Đang làm module của Hùng thì không tự ý đổi report schema của Khôi.
- Đang làm backend của Khôi thì không tự ý redesign UI của Lai.

### Bước 4: Kiểm tra DoD

Mỗi task hoàn thành phải có:

- Code chạy được.
- Không lỗi build.
- Không lỗi lint nghiêm trọng.
- API có response chuẩn.
- UI không vỡ layout.
- Có xử lý loading/error/empty state.
- Có kiểm tra quyền.
- Có log/audit nếu là hành động quan trọng.
- Có test hoặc checklist test thủ công.

---

## 7. Skill phân tích AI feature

Khi làm feature liên quan AI, luôn phải tách rõ:

| Thành phần | Câu hỏi cần trả lời |
|---|---|
| Input | AI nhận dữ liệu gì? JD, CV, transcript, rubric hay note? |
| Permission | AI có được phép đọc dữ liệu đó không? |
| Prompt | Prompt template nào được dùng? |
| Output | Output dạng text hay JSON? |
| Evidence | Có trích dẫn bằng chứng không? |
| Confidence | Có độ tin cậy không? |
| Storage | Có lưu kết quả AI không? |
| Review | Recruiter có chỉnh sửa/override được không? |
| Cost | Có giới hạn token/cost không? |

AI output nên ưu tiên JSON có cấu trúc để backend/frontend dễ dùng.

Ví dụ output chấm điểm:

```json
{
  "criterion": "Technical Knowledge",
  "score": 4,
  "max_score": 5,
  "evidence": "Ứng viên mô tả được cách tối ưu query và index.",
  "comment": "Câu trả lời có kinh nghiệm thực tế.",
  "confidence": 0.82,
  "insufficient_data": false
}
```

---

## 8. Skill thiết kế realtime room

Khi làm realtime room, luôn phải có:

- Room state.
- User presence.
- Role trong room.
- Event join/leave.
- Event reconnect.
- Event start/end interview.
- Event chat.
- Event transcript.
- Event AI suggestion.
- Event AI score update.
- Event error.

Room không được phụ thuộc hoàn toàn vào AI. Nếu AI lỗi, phỏng vấn vẫn phải tiếp tục.

---

## 9. Skill thiết kế database

Khi thêm bảng hoặc field mới:

- Tên bảng dùng snake_case.
- Primary key ưu tiên UUID.
- Có `created_at`, `updated_at` nếu là entity chính.
- Có `deleted_at` nếu cần soft delete.
- Có foreign key rõ ràng.
- Dữ liệu AI phức tạp có thể lưu JSONB nhưng field quan trọng phải tách cột riêng.
- Không lưu token/key bí mật dạng plain text.
- File CV/recording chỉ lưu URL hoặc storage key, không lưu blob trong DB.

---

## 10. Skill thiết kế UI

UI phải ưu tiên:

- Rõ ràng.
- Dễ dùng.
- Không màu mè quá mức.
- Có loading state.
- Có empty state.
- Có error state.
- Có confirm khi thao tác nguy hiểm.
- Không hiển thị dữ liệu nội bộ cho Candidate.
- Candidate UI phải đơn giản, ít gây áp lực.
- Recruiter UI phải đủ thông tin để ra quyết định.

---

## 11. Skill kiểm thử

Mỗi module cần có checklist test:

- Happy path.
- Missing data.
- Permission denied.
- Unauthorized.
- Not found.
- Validation error.
- Network error.
- AI error.
- Realtime reconnect.
- Empty state.
- Large data.

---

## 12. Prompt mẫu cho AI code trong dự án

```text
Bạn đang làm dự án nền tảng phỏng vấn cùng AI real-time.
Hãy đọc kỹ source hiện tại trước khi code.
Không tự ý phá kiến trúc cũ.
Không sửa module ngoài phạm vi nếu không cần thiết.

Module đang làm: [Tên module]
Owner: [Hùng/Khôi/Lai]
Mục tiêu: [Mô tả mục tiêu]

Yêu cầu:
1. Phân tích file liên quan trước.
2. Liệt kê kế hoạch sửa.
3. Code theo kiến trúc hiện tại.
4. Đảm bảo không breaking change.
5. Có loading/error/empty state nếu là UI.
6. Có validation và permission nếu là API.
7. Có test/checklist sau khi hoàn thành.
8. Báo rõ file đã sửa và cách test.
```

---

## 13. Kết luận

Bộ skill này là tài liệu nền để AI và team hiểu dự án. Khi bắt đầu task mới, luôn đọc:

1. `AI_INTERVIEW_PROJECT_SKILL.md`
2. `AI_INTERVIEW_PROJECT_RULES.md`
3. File module của người đang làm:
   - `MODULE_HUNG_REALTIME_AI.md`
   - `MODULE_KHOI_CORE_AI_BACKEND.md`
   - `MODULE_LAI_PRODUCT_UI_BUSINESS.md`
