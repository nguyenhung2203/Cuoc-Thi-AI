# MODULE_KHOI_CORE_AI_BACKEND.md
# Phân công Module cho Khôi

## 1. Vai trò của Khôi

Khôi phụ trách nhóm module **khó và phức tạp nhất về backend core, database, AI engine, scoring, report và nghiệp vụ dữ liệu lõi**.

Khôi là owner chính của các module:

1. Backend architecture.
2. Auth/RBAC core.
3. Database schema/migration.
4. Job backend.
5. Candidate backend.
6. Interview backend.
7. AI Orchestrator.
8. JD/CV Analysis.
9. Question Generation.
10. AI Scoring Engine.
11. Report Generator.
12. Audit Log.
13. API contract chuẩn cho Hùng và Lai dùng.

---

## 2. Mục tiêu module của Khôi

Mục tiêu chính:

- Có backend core ổn định cho toàn hệ thống.
- Có database thiết kế đúng để mở rộng.
- Có API chuẩn cho Job, Candidate, Interview, Report.
- Có AI engine xử lý JD/CV/question/scoring/report.
- Có phân quyền rõ Recruiter/Candidate.
- Có audit log cho hành động quan trọng.
- Có output AI có cấu trúc, có evidence, có confidence.

---

## 3. Phạm vi công việc

## 3.1. Backend Foundation

### Chức năng

- Thiết kế cấu trúc service/module.
- Chuẩn hóa response API.
- Chuẩn hóa error handling.
- Middleware auth.
- Middleware permission.
- Validation request.
- Logging.
- Config environment.

### DoD

- API response thống nhất.
- Error code rõ ràng.
- Middleware auth hoạt động.
- Có base service/repository pattern nếu phù hợp.
- Có tài liệu API cơ bản.

---

## 3.2. Auth & RBAC Core

### Chức năng

- Đăng ký/đăng nhập.
- JWT/refresh token.
- Phân biệt role:
  - Admin
  - Recruiter Owner
  - Recruiter Member
  - Candidate
- Permission theo company/workspace.
- Candidate chỉ xem được dữ liệu của mình.
- Recruiter chỉ xem dữ liệu trong company của mình.

### DoD

- Login/logout hoạt động.
- Token hết hạn xử lý đúng.
- API protected chặn user chưa đăng nhập.
- API chặn sai quyền.
- Không leak dữ liệu giữa company.

---

## 3.3. Database Schema

### Bảng lõi Khôi cần thiết kế

- users
- companies
- company_members
- jobs
- candidates
- candidate_jobs
- interview_templates
- interview_questions
- interviews
- interview_participants
- interview_transcripts
- interview_scores
- interview_reports
- mock_interviews
- mock_interview_answers
- audit_logs
- ai_prompt_templates
- ai_request_logs

### Nguyên tắc

- Dùng UUID.
- Dùng foreign key rõ ràng.
- Entity chính có created_at/updated_at.
- Dữ liệu AI phức tạp lưu JSONB nhưng field quan trọng tách riêng.
- File CV/recording lưu storage key/url, không lưu blob.

### DoD

- Migration chạy được từ đầu.
- Có rollback nếu dự án yêu cầu.
- Quan hệ dữ liệu rõ.
- Không thiếu index cho query chính.

---

## 3.4. Job Backend

### Chức năng

- CRUD job.
- Đóng/mở job.
- Lưu JD.
- AI analyze JD.
- AI đề xuất rubric.
- AI đề xuất question set.

### API đề xuất

```text
GET    /jobs
POST   /jobs
GET    /jobs/:id
PUT    /jobs/:id
DELETE /jobs/:id
POST   /jobs/:id/analyze
POST   /jobs/:id/generate-questions
```

### DoD

- Recruiter tạo/sửa/xóa job trong company của mình.
- Candidate không truy cập job nội bộ nếu không được phép.
- AI analyze JD trả output có cấu trúc.
- Có validation cho title/description/status.

---

## 3.5. Candidate Backend

### Chức năng

- CRUD candidate.
- Upload CV metadata.
- Parse CV bằng AI.
- Gán candidate vào job.
- Cập nhật trạng thái candidate.
- Lưu nguồn ứng viên.

### API đề xuất

```text
GET    /candidates
POST   /candidates
GET    /candidates/:id
PUT    /candidates/:id
DELETE /candidates/:id
POST   /candidates/:id/upload-cv
POST   /candidates/:id/parse-cv
POST   /jobs/:job_id/candidates/:candidate_id/assign
```

### DoD

- Candidate được tạo và gán job.
- CV metadata lưu đúng.
- Parsed CV lưu JSONB.
- Không leak candidate giữa company.

---

## 3.6. Interview Backend

### Chức năng

- Tạo lịch phỏng vấn.
- Gán job, candidate, recruiter.
- Tạo room token/link.
- Start/end interview.
- Lưu status.
- Lưu transcript.
- Trigger report generation.

### API đề xuất

```text
GET    /interviews
POST   /interviews
GET    /interviews/:id
PUT    /interviews/:id
POST   /interviews/:id/start
POST   /interviews/:id/end
POST   /interviews/:id/cancel
GET    /interviews/:id/transcripts
GET    /interviews/:id/report
```

### DoD

- Tạo lịch phỏng vấn được.
- Link room hợp lệ.
- Start/end đúng quyền.
- End interview trigger report generation.
- Transcript lưu được.

---

## 3.7. AI Orchestrator

### Chức năng

AI Orchestrator là lớp trung tâm điều phối các AI service:

- Analyze JD.
- Analyze CV.
- Generate questions.
- Suggest follow-up.
- Score answer.
- Generate report.
- Mock interview question.
- Mock feedback.

### Input chung

- JD.
- CV.
- Transcript.
- Rubric.
- Questions asked.
- Recruiter note nếu có quyền.
- Interview metadata.

### Output chung

- JSON có cấu trúc.
- Có evidence nếu là scoring/report.
- Có confidence nếu phù hợp.
- Có insufficient_data nếu thiếu dữ liệu.

### DoD

- Có service gọi AI tập trung.
- Có prompt template version.
- Có log request/response cần thiết.
- Có retry/error handling.
- Không làm crash backend khi AI provider lỗi.

---

## 3.8. JD Analysis

### Mục tiêu

AI đọc JD và trả về:

- Tóm tắt vị trí.
- Skill chính.
- Skill phụ.
- Level kỳ vọng.
- Bộ tiêu chí đánh giá.
- Gợi ý câu hỏi.

### Output đề xuất

```json
{
  "job_summary": "string",
  "required_skills": ["Vue", "TypeScript"],
  "nice_to_have_skills": ["Testing", "CI/CD"],
  "level": "middle",
  "rubric": [
    {
      "criterion": "Technical Knowledge",
      "weight": 0.3,
      "description": "Kiến thức chuyên môn phù hợp vị trí"
    }
  ],
  "suggested_questions": []
}
```

---

## 3.9. CV Analysis

### Mục tiêu

AI đọc CV và trả về:

- Thông tin tổng quan.
- Kỹ năng.
- Kinh nghiệm.
- Điểm phù hợp với JD.
- Điểm cần xác minh khi phỏng vấn.

### Output đề xuất

```json
{
  "candidate_summary": "string",
  "skills": ["Vue", "React", "Node.js"],
  "years_of_experience": 2,
  "matched_requirements": [],
  "missing_requirements": [],
  "interview_focus_points": []
}
```

---

## 3.10. Question Generation

### Mục tiêu

Sinh câu hỏi phù hợp:

- Theo JD.
- Theo CV.
- Theo level.
- Theo loại phỏng vấn.
- Không hỏi trùng.

### Loại câu hỏi

- introduction
- technical
- behavioral
- situational
- experience
- culture
- follow_up
- closing

### DoD

- Câu hỏi lưu được vào question bank.
- Có difficulty.
- Có target skill.
- Có reason vì sao hỏi câu đó.

---

## 3.11. AI Scoring Engine

### Chức năng

- Chấm từng câu trả lời.
- Chấm theo rubric.
- Tính điểm tổng theo weight.
- Có evidence từ transcript.
- Có confidence.
- Ghi rõ thiếu dữ liệu.

### Output đề xuất

```json
{
  "criterion": "Problem Solving",
  "score": 4,
  "max_score": 5,
  "weight": 0.2,
  "evidence": "Ứng viên mô tả được cách chia nhỏ vấn đề và kiểm tra giả thuyết.",
  "comment": "Tư duy giải quyết vấn đề tốt.",
  "confidence": 0.81,
  "insufficient_data": false
}
```

### DoD

- Score lưu vào DB.
- Không bịa evidence.
- Không score khi thiếu dữ liệu nghiêm trọng.
- Có API hoặc event cho Hùng đẩy realtime.
- Có dữ liệu cho Lai hiển thị UI.

---

## 3.12. Report Generator

### Chức năng

Tạo báo cáo sau phỏng vấn:

- Thông tin chung.
- Tóm tắt buổi phỏng vấn.
- Điểm tổng.
- Điểm theo rubric.
- Điểm mạnh.
- Điểm yếu.
- Rủi ro.
- Câu trả lời nổi bật.
- Câu trả lời cần cải thiện.
- AI recommendation.
- Recruiter final decision.
- Transcript.

### Output đề xuất

```json
{
  "summary": "string",
  "final_score": 4.1,
  "strengths": ["string"],
  "weaknesses": ["string"],
  "risks": ["string"],
  "recommendation": "hire",
  "reasoning": "string",
  "scores": []
}
```

### DoD

- Report sinh sau khi end interview.
- Lưu report vào DB.
- Recruiter xem được report.
- Có trạng thái report: pending/generating/ready/failed.
- Có retry nếu report failed.

---

## 3.13. Audit Log

### Hành động cần log

- Login.
- Tạo/sửa/xóa job.
- Thêm/sửa/xóa candidate.
- Upload CV.
- Tạo lịch phỏng vấn.
- Start/end/cancel interview.
- Xem/tải report.
- Thay đổi final decision.

### DoD

- Audit log lưu actor/action/resource.
- Có metadata đủ dùng.
- Không lưu dữ liệu quá nhạy cảm vào log.

---

## 4. Sprint chia việc cho Khôi

## Sprint 0: Backend foundation

| Task | Mô tả | Độ khó | DoD |
|---|---|---:|---|
| K-S0-01 | Thiết kế database schema lõi | Rất khó | ERD + migration ban đầu |
| K-S0-02 | Setup API response/error standard | Trung bình | Response thống nhất |
| K-S0-03 | Auth middleware | Khó | API protected hoạt động |
| K-S0-04 | RBAC company scope | Rất khó | Không leak dữ liệu cross-company |

## Sprint 1: Job + Candidate backend

| Task | Mô tả | Độ khó | DoD |
|---|---|---:|---|
| K-S1-01 | Job CRUD API | Trung bình | CRUD + validation |
| K-S1-02 | Candidate CRUD API | Trung bình | CRUD + company scope |
| K-S1-03 | Candidate assign to job | Khó | Gán ứng viên vào job |
| K-S1-04 | CV upload metadata API | Trung bình | Lưu file metadata |

## Sprint 2: Interview backend

| Task | Mô tả | Độ khó | DoD |
|---|---|---:|---|
| K-S2-01 | Interview scheduling API | Khó | Tạo lịch phỏng vấn |
| K-S2-02 | Room link/token backend | Khó | Link hợp lệ, có expire |
| K-S2-03 | Start/end interview API | Trung bình | Đúng quyền |
| K-S2-04 | Transcript storage API | Khó | Hùng lưu transcript được |

## Sprint 3: AI foundation

| Task | Mô tả | Độ khó | DoD |
|---|---|---:|---|
| K-S3-01 | AI Orchestrator service | Rất khó | Gọi AI tập trung |
| K-S3-02 | Prompt template manager | Khó | Prompt có version |
| K-S3-03 | AI request log | Trung bình | Log input/output cần thiết |
| K-S3-04 | AI error/retry handling | Khó | AI lỗi không crash |

## Sprint 4: JD/CV/Question AI

| Task | Mô tả | Độ khó | DoD |
|---|---|---:|---|
| K-S4-01 | Analyze JD | Khó | Trả rubric/skills/questions |
| K-S4-02 | Analyze CV | Khó | Trả summary/skills/focus points |
| K-S4-03 | Generate questions | Khó | Câu hỏi theo JD/CV/level |
| K-S4-04 | Save question bank | Trung bình | Lưu và gán question |

## Sprint 5: Scoring + Report

| Task | Mô tả | Độ khó | DoD |
|---|---|---:|---|
| K-S5-01 | Rubric scoring engine | Rất khó | Score theo tiêu chí |
| K-S5-02 | Score evidence extraction | Rất khó | Có evidence từ transcript |
| K-S5-03 | Report generator | Rất khó | Sinh report sau interview |
| K-S5-04 | Report retry/status | Khó | pending/generating/ready/failed |

## Sprint 6: Mock Interview backend

| Task | Mô tả | Độ khó | DoD |
|---|---|---:|---|
| K-S6-01 | Mock interview session API | Khó | Candidate tạo session |
| K-S6-02 | Mock AI question flow | Khó | AI hỏi theo vị trí |
| K-S6-03 | Mock answer scoring | Khó | Feedback từng câu |
| K-S6-04 | Mock report API | Trung bình | Candidate xem feedback |

---

## 5. Dependency với Hùng

Khôi cần cung cấp cho Hùng:

- Interview API.
- Room token validation.
- Transcript storage API/service.
- AI suggestion endpoint/service.
- AI scoring endpoint/service.
- Report generation trigger.

Khôi cần nhận từ Hùng:

- Event contract thực tế.
- Room status update.
- Transcript realtime data.
- Start/end room event.

---

## 6. Dependency với Lai

Khôi cần cung cấp cho Lai:

- API Job.
- API Candidate.
- API Interview.
- API Report.
- API Mock Interview.
- API AI analyze JD/CV.
- API question bank.

Khôi cần nhận từ Lai:

- UI yêu cầu field nào.
- Form validation cần gì.
- Dashboard cần thống kê gì.
- Report UI cần dữ liệu dạng nào.

---

## 7. Checklist test riêng của Khôi

- Auth login/logout.
- Token hết hạn.
- Recruiter không xem được company khác.
- Candidate không xem được report nội bộ.
- CRUD job đúng.
- CRUD candidate đúng.
- Assign candidate vào job đúng.
- Tạo interview đúng.
- Start/end đúng quyền.
- Transcript lưu đúng.
- AI analyze JD đúng schema.
- AI analyze CV đúng schema.
- AI scoring có evidence.
- Report sinh được.
- AI lỗi có retry hoặc failed status.
- Audit log tạo đúng.

---

## 8. Prompt mẫu cho Khôi dùng với AI coding assistant

```text
Bạn đang làm module của Khôi trong dự án phỏng vấn cùng AI real-time.
Phạm vi của Khôi là Backend Core, Database, Auth/RBAC, Job, Candidate, Interview API, AI Orchestrator, Scoring, Report và Audit Log.

Yêu cầu task hiện tại:
[Điền task]

Trước khi code hãy:
1. Đọc schema/model/API hiện có.
2. Xác định có cần migration không.
3. Không phá API contract mà Hùng/Lai đang dùng.
4. Output AI phải có JSON schema rõ ràng.
5. Scoring phải có evidence, confidence, insufficient_data.
6. Recruiter là người quyết định cuối cùng, AI chỉ recommendation.
7. Có validation, permission, error handling.
8. Sau khi code, liệt kê file đã sửa và checklist test.
```

---

## 9. Kết quả cuối cùng Khôi cần bàn giao

Khôi hoàn thành khi hệ thống có:

- Backend core ổn định.
- Database schema đầy đủ.
- Auth/RBAC hoạt động.
- Job/Candidate/Interview API hoạt động.
- AI Orchestrator hoạt động.
- JD/CV analysis hoạt động.
- Question generation hoạt động.
- Scoring engine có evidence.
- Report generator hoạt động.
- Mock interview backend hoạt động.
- Audit log cơ bản hoạt động.
- API contract đủ cho Hùng và Lai tích hợp.
