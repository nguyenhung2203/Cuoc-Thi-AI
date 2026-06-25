# CHECKLIST CHI TIẾT CÔNG VIỆC - KHÔI

> **Vai trò:** Owner Backend Core, Database, AI Engine, Scoring, Report & Audit Log
> **Tham chiếu:**
> - [MODULE_KHOI_CORE_AI_BACKEND.md](MODULE_KHOI_CORE_AI_BACKEND.md) — phân công module
> - [DATABASE_DESIGN.md](../DATABASE_DESIGN.md) — Source of Truth cho schema/tên bảng/field
> - [API_SPEC.md](../API_SPEC.md) — Source of Truth cho endpoint/response/permission
> - [SYSTEM_ANALYSIS.md](../SYSTEM_ANALYSIS.md) — phân tích hệ thống
> **Stack:** Go backend + Python AI service + PostgreSQL
> **Convention API:** Base URL `/api/v1`, đa phần endpoint nghiệp vụ scope theo company → `/api/v1/companies/:company_id/...`

---

## SPRINT 0 — BACKEND FOUNDATION

### K-S0-01: Thiết kế Database Schema lõi `[Rất khó]`

> Bám đúng tên bảng và field theo [DATABASE_DESIGN.md](../DATABASE_DESIGN.md). Không tự đổi tên — nếu cần đổi, cập nhật tài liệu trước.

**Bảng định danh & tổ chức:**
- [x] `users` — tài khoản đăng nhập (admin/recruiter/candidate, có `email_verified_at`, `last_login_at`, soft delete)
- [x] `companies` — workspace tuyển dụng (slug unique, `created_by`, `settings` jsonb)
- [x] `company_members` — thành viên company (role owner/admin/member/viewer, unique `(company_id, user_id)`)

**Bảng nghiệp vụ tuyển dụng:**
- [x] `jobs` — vị trí tuyển dụng (có `ai_summary`, `ai_analysis_json`, `status` draft/open/paused/closed)
- [x] `candidates` — hồ sơ CRM của company (`user_id` nullable, `cv_file_id`, `parsed_cv_json`)
- [x] `job_candidates` — bảng nối job ↔ candidate **(KHÔNG phải `candidate_jobs`)** — có `pipeline_status`, `fit_score`, `ai_match_json`, unique `(job_id, candidate_id)`

**Bảng phỏng vấn:**
- [x] `interview_templates` — mẫu phỏng vấn theo company
- [x] `question_bank` — kho câu hỏi **(KHÔNG phải `interview_questions`)** — có `question_type`, `skill_tags`, `level`, `expected_signals`, `is_ai_generated`
- [x] `rubrics` — bộ tiêu chí đánh giá theo job/company
- [x] `rubric_criteria` — tiêu chí trong rubric (`weight`, `min_score`, `max_score`, `scoring_guide`, `order_index`)
- [x] `interviews` — phiên phỏng vấn (mode real/mock, `consent_recording`, `consent_ai`, `invite_token_hash`, `invite_expires_at`)
- [x] `interview_rooms` — metadata phòng realtime **(BẢNG MỚI BỔ SUNG)** — có `room_code` unique, `provider`, `connection_config`, `opened_at`/`closed_at`, unique `interview_id`
- [x] `interview_participants` — người trong phòng (`participant_type`, `connection_state`, `media_status`)
- [x] `interview_transcripts` — transcript từng đoạn (`speaker_type`, `is_final`, `confidence`, `edited_content`)
- [x] `ai_suggestions` — gợi ý AI realtime **(BẢNG MỚI BỔ SUNG)** — có `suggestion_type`, `priority`, `confidence`, `accepted_by`/`dismissed_by`
- [x] `interview_scores` — điểm theo tiêu chí (snapshot `criterion_name`/`weight`, `weighted_score`, `evidence`, `status`)
- [x] `interview_reports` — report cuối (1-1 với `interview_id`, có `recommendation`, `recruiter_decision`, `report_json`)

**Bảng mock interview:**
- [x] `mock_interviews` — phiên luyện cá nhân (dùng `user_id`, **không** dùng `candidate_id`)
- [x] `mock_interview_messages` — tin nhắn/câu trả lời mock **(KHÔNG phải `mock_interview_answers`)** — có `sender_type` ai/candidate/system, `score_json`

**Bảng hệ thống:**
- [x] `files` — metadata file CV/recording/avatar **(BẢNG MỚI BỔ SUNG)** — có `storage_key`, `mime_type`, `size_bytes`, `file_type`, `checksum`
- [x] `audit_logs` — log hành động (`actor_user_id`, `action`, `resource_type`, `before_json`/`after_json`, `ip_address`)
- [x] `notifications` — thông báo trong app **(BẢNG MỚI BỔ SUNG)** — có `type`, `title`, `data_json`, `read_at`. Cần cho K-S5-04 (notify khi report ready/failed)

**Cần đồng bộ với DATABASE_DESIGN.md trước khi code:**
- [x] **CHƯA CÓ trong DATABASE_DESIGN.md:** `ai_prompt_templates` và `ai_request_logs` (dùng cho K-S3-02, K-S3-03). **Action:** soạn schema 2 bảng này, tạo PR cập nhật `DATABASE_DESIGN.md` mục 19.5/19.6 trước khi viết migration. Đề xuất field tối thiểu:
  - `ai_prompt_templates`: `id`, `name`, `version`, `content`, `variables_schema jsonb`, `model`, `params jsonb`, `is_active`, `created_by`, `created_at`, `updated_at` — unique `(name, version)`
  - `ai_request_logs`: `id`, `template_id`, `template_version`, `interview_id?`, `job_id?`, `candidate_id?`, `input_json`, `output_json`, `latency_ms`, `tokens_in`, `tokens_out`, `cost`, `status`, `error`, `created_at`

**Quy ước chung:**
- [x] PostgreSQL + extensions `uuid-ossp`, `pgcrypto`
- [x] PK `UUID`, tên bảng số nhiều snake_case, field snake_case
- [x] Mọi entity nghiệp vụ có `created_at`, `updated_at`, soft delete `deleted_at` khi cần
- [x] Mọi dữ liệu tuyển dụng có `company_id` (multi-tenant isolation)
- [x] Dữ liệu AI linh hoạt → `jsonb`, file lớn → chỉ lưu `storage_key`
- [x] Enum đồng bộ với DATABASE_DESIGN mục 1.3 (user role/status, job status, candidate status, interview status, mode, recommendation, speaker type)
- [x] Migration theo thứ tự DATABASE_DESIGN mục 26 (users → companies → company_members → files → jobs → candidates → job_candidates → templates → question_bank → rubrics → rubric_criteria → interviews → interview_rooms → participants → transcripts → ai_suggestions → scores → reports → mock_interviews → mock_interview_messages → audit_logs → notifications)
- [x] Index cho query phổ biến (FK, `(company_id, status)`, `(interview_id, created_at)`, ...)
- [x] Rollback migration cho mỗi step

**DoD:**
- ERD đủ 22 bảng (20 bảng trong DATABASE_DESIGN + 2 bảng AI mới đã được cập nhật vào tài liệu)
- Migration chạy `up` từ đầu sạch + `down` rollback được
- Không lệch tên bảng/field so với DATABASE_DESIGN.md
- API response không leak `password_hash`, `invite_token_hash`

---

### K-S0-02: API Response/Error Standard `[Trung bình]`

> Bám [API_SPEC.md](../API_SPEC.md) mục 1.

- [x] Base URL `/api/v1`
- [x] Response success envelope: `{ success: true, data, meta, request_id }`
- [x] Response error envelope: `{ success: false, error: { code, message, details }, request_id }`
- [x] Sinh `request_id` cho mỗi request và log kèm
- [x] Error code chuẩn theo API_SPEC mục 1.5: `UNAUTHORIZED` (401), `FORBIDDEN` (403), `NOT_FOUND` (404), `VALIDATION_ERROR` (422), `CONFLICT` (409), `RATE_LIMITED` (429), `AI_SERVICE_ERROR` (502), `REALTIME_ERROR` (502), `INTERNAL_ERROR` (500)
- [x] Pagination chuẩn: query `?page=1&page_size=20&sort=created_at:desc`, meta `{ page, page_size, total, total_pages }`
- [x] Helper/middleware bọc response thống nhất
- [x] Validation request (body/query/params) — trả `VALIDATION_ERROR` với mảng `details`
- [x] Logging request/response (không log token/password)
- [x] Config environment (.env, secrets management)
- [x] Tài liệu API (OpenAPI hoặc bám API_SPEC.md)
- [x] **DoD:** Response/error format khớp API_SPEC.md, FE Lai parse được không cần đoán

---

### K-S0-03: Auth API & Middleware `[Khó]`

> Endpoint theo [API_SPEC.md](../API_SPEC.md) mục 2.

- [x] `POST /auth/register` — body `{ email, password, full_name, role }`, trả `{ user, access_token, refresh_token }`
- [x] `POST /auth/login` — body `{ email, password }`, trả token giống register
- [x] `GET /auth/me` **(BỔ SUNG)** — trả `{ id, email, full_name, role, companies: [{ id, name, role }] }`. Lai cần để hiển thị user info sau login
- [x] `POST /auth/refresh` — đổi access token
- [x] `POST /auth/logout` — revoke token
- [x] Hash password (bcrypt/argon2)
- [x] Middleware xác thực Bearer token, attach `user` + danh sách `company_members` vào request context
- [x] Xử lý invite token cho candidate join room: route public `/interviews/join/:invite_token`
- [x] Token hết hạn → trả `UNAUTHORIZED` 401 chuẩn
- [x] **DoD:** API protected chặn user chưa đăng nhập, `GET /auth/me` trả đúng companies kèm role

---

### K-S0-04: RBAC Company Scope `[Rất khó]`

- [x] Phân role hệ thống: `admin`, `recruiter`, `candidate` (theo `users.role`)
- [x] Phân role trong company: `owner`, `admin`, `member`, `viewer` (theo `company_members.role`)
- [x] Middleware permission lấy `company_id` từ URL → kiểm tra user là member của company
- [x] Helper `requirePermission(action, resource)` cho các action quan trọng (job:create, job:update, interview:create...)
- [x] Scope dữ liệu mọi query theo `company_id` (theo DATABASE_DESIGN mục 27)
- [ ] Candidate chỉ xem dữ liệu của chính mình (mock interview, report được share)
- [x] Test cross-company không leak (tạo 2 company, recruiter A không xem được job của B)
- [x] **DoD:** Không leak dữ liệu giữa company, API chặn sai quyền với code `FORBIDDEN`

---

### K-S0-05: Company API `[Trung bình]` **(TASK BỔ SUNG)**

> Endpoint theo [API_SPEC.md](../API_SPEC.md) mục 3. Khôi phụ trách RBAC company scope nên cần luôn CRUD company.

- [x] `GET /companies` — list company của user hiện tại (kèm `role` của user trong từng company)
- [x] `POST /companies` — tạo company mới (`name`, `website`, `industry`, `size`); user tạo tự động thành `owner` trong `company_members`
- [x] `GET /companies/:company_id` — chi tiết company (chỉ member xem được)
- [x] `PUT /companies/:company_id` — update company (chỉ owner/admin)
- [x] Sinh `slug` unique từ `name`
- [ ] **DoD:** User tạo được company và tự động trở thành owner; member khác không sửa được info company

---

## SPRINT 1 — JOB + CANDIDATE BACKEND

> **Lưu ý đường dẫn:** Mọi endpoint job/candidate scope theo company. Dùng `/companies/:company_id/jobs/...` chứ KHÔNG phải `/jobs/...`.

### K-S1-01: Job CRUD API `[Trung bình]` ✅ DONE

- [x] `GET /companies/:company_id/jobs?status=&keyword=&page=&page_size=` — list + filter + pagination, response item có `candidate_count`, `interview_count`, `avg_fit_score`
- [x] `POST /companies/:company_id/jobs` — body theo API_SPEC mục 4.2, validate `title`, `description`, `status`
- [x] `GET /companies/:company_id/jobs/:job_id` — chi tiết kèm `ai_summary`, `ai_analysis_json`, `stats`
- [x] `PUT /companies/:company_id/jobs/:job_id`
- [x] `DELETE /companies/:company_id/jobs/:job_id` — soft delete (set `deleted_at`)
- [x] Đóng/mở job qua `status` (draft/open/paused/closed)
- [x] RBAC: `job:create`, `job:update`, `job:delete` theo company role
- [x] **DoD:** CRUD đầy đủ, validation chặt, scope company, soft delete

---

### K-S1-02: Candidate CRUD API `[Trung bình]` ✅ DONE

- [x] `GET /companies/:company_id/candidates?job_id=&status=&keyword=&page=&page_size=`
- [x] `POST /companies/:company_id/candidates` — body `{ full_name, email, phone, source, job_id? }`; nếu có `job_id` tạo luôn record `job_candidates`
- [x] `GET /companies/:company_id/candidates/:candidate_id` — chi tiết kèm `cv_file` (signed url), `ai_cv_summary`, `skills`, `interviews`, `reports`
- [x] `PUT /companies/:company_id/candidates/:candidate_id`
- [x] `DELETE /companies/:company_id/candidates/:candidate_id` — soft delete
- [x] Cập nhật `status` (new/screening/invited/interviewing/...)
- [x] Lưu `source` (linkedin/referral/import/manual)
- [x] `tags` lưu jsonb array
- [x] Phân biệt rõ `candidates` (CRM record của company) ≠ `users` (account đăng nhập). `candidates.user_id` nullable
- [x] **DoD:** CRUD đúng, không leak cross-company, mỗi candidate gắn đúng `company_id`

---

### K-S1-03: Candidate Assign to Job `[Khó]` ✅ DONE

- [x] Endpoint gán: dùng `POST /companies/:company_id/candidates` với `job_id` HOẶC tạo endpoint riêng `POST /companies/:company_id/jobs/:job_id/candidates/:candidate_id/assign` (chốt với Lai 1 trong 2)
- [x] Lưu vào `job_candidates` với `pipeline_status`, `fit_score`, `applied_at`
- [x] Validate candidate và job cùng `company_id`
- [x] Unique `(job_id, candidate_id)` — không assign trùng
- [x] API list candidate theo job: `GET /companies/:company_id/jobs/:job_id/candidates`
- [x] API update pipeline status: `PUT /companies/:company_id/jobs/:job_id/candidates/:candidate_id/pipeline`
- [ ] API unassign nếu cần
- [x] **DoD:** Gán đúng, không trùng, có lịch sử pipeline status

---

### K-S1-04: CV Upload + Files API `[Trung bình]` ⚠️ PARTIAL

> Theo [API_SPEC.md](../API_SPEC.md) mục 5.5, 5.6, 12. Dùng bảng `files`.

- [x] `POST /companies/:company_id/candidates/:candidate_id/cv` — multipart, lưu vào object storage, ghi metadata vào `files` (`file_type=cv`), set `candidates.cv_file_id`
- [ ] `POST /companies/:company_id/candidates/:candidate_id/parse-cv` — trigger AI parse, lưu `parsed_cv_json`, `ai_cv_summary`
- [x] `POST /files` (general upload) — trả `{ file_id, signed_url }`
- [x] `GET /files/:file_id/signed-url` — sinh signed URL có expire (stub 15 phút, chưa kết nối object storage thật)
- [x] Validate mime type, size
- [x] Tính `checksum` (sha256) để dedup
- [x] **DoD:** CV upload lưu đúng vào `files`, không lưu blob trong DB, có signed URL

---

## SPRINT 2 — INTERVIEW BACKEND

### K-S2-01: Interview Scheduling API `[Khó]`

> Theo [API_SPEC.md](../API_SPEC.md) mục 6.

- [x] `GET /companies/:company_id/interviews?status=&date=&page=&page_size=`
- [x] `POST /companies/:company_id/interviews` — body `{ job_id, candidate_id, recruiter_id, template_id?, rubric_id?, scheduled_at, duration_minutes, mode, send_invite }`; tạo luôn record `interview_rooms` và sinh `invite_token`/`invite_url`
- [x] `GET /companies/:company_id/interviews/:interview_id` — chi tiết kèm `room`, `rubric`, `job`, `candidate`, `recruiter`
- [x] `PUT /companies/:company_id/interviews/:interview_id` — reschedule (chỉ khi `status=scheduled`)
- [x] `POST /companies/:company_id/interviews/:interview_id/cancel`
- [x] Validate trùng lịch recruiter/candidate (cảnh báo, không cứng)
- [x] Lưu `consent_recording`, `consent_ai`
- [x] **DoD:** Tạo lịch đúng, có room + invite_url, RBAC company scope

---

### K-S2-02: Room Link/Token Backend `[Khó]`

> Theo [API_SPEC.md](../API_SPEC.md) mục 7. Dùng bảng `interview_rooms`.

- [x] Sinh `room_code` unique, `invite_token_hash` lưu trong `interviews`, `invite_expires_at`
- [x] `GET /companies/:company_id/interviews/:interview_id/room` — trả room detail (chỉ member)
- [x] `POST /companies/:company_id/interviews/:interview_id/room/access-token` — sinh token realtime cho recruiter (có expire ngắn)
- [x] `GET /interviews/join/:invite_token` (public) — validate invite token → trả room info cho candidate
- [x] Revoke token khi end/cancel interview
- [x] **DoD:** Hùng dùng được API validate token, không leak `invite_token_hash` ra response

---

### K-S2-03: Start/End Interview API `[Trung bình]`

- [x] `POST /companies/:company_id/interviews/:interview_id/start` — set `status=active`, `started_at`, mở `interview_rooms.status=active`, emit event cho Hùng
- [x] `POST /companies/:company_id/interviews/:interview_id/end` — set `status=completed`, `ended_at`, đóng room, **trigger báo cáo tự động** (K-S5-03)
- [x] `POST /companies/:company_id/interviews/:interview_id/cancel` — set `status=cancelled`, revoke token
- [x] Kiểm tra quyền: chỉ recruiter của interview hoặc owner/admin company
- [x] State machine: scheduled → waiting → active → completed (hoặc cancelled/expired)
- [x] **DoD:** State chuyển đúng, end interview trigger report generation

---

### K-S2-04: Transcript Storage API `[Khó]`

> Theo [API_SPEC.md](../API_SPEC.md) mục 8. Đây là endpoint Hùng push lên.

- [x] `GET /companies/:company_id/interviews/:interview_id/transcripts?since=&limit=` — list transcript
- [x] `POST /companies/:company_id/interviews/:interview_id/transcripts` — Hùng push (`speaker_type`, `content`, `language`, `start_time_ms`, `end_time_ms`, `confidence`, `is_final`, `source`)
- [x] `PUT /companies/:company_id/interviews/:interview_id/transcripts/:transcript_id` — edit transcript (`edited_content`, `edited_by`, `edited_at`)
- [x] Append-only cho item `is_final=true`, partial có thể update
- [x] Index `(interview_id, created_at)` để load nhanh
- [x] **DoD:** Hùng push transcript được realtime, Lai đọc được, hỗ trợ edit thủ công

---

## SPRINT 3 — AI FOUNDATION

### K-S3-01: AI Orchestrator Service `[Rất khó]`

- [ ] Service trung tâm trong Go gọi sang Python `ai-service`
- [ ] Function `analyzeJD(job) -> JDAnalysis`
- [ ] Function `analyzeCV(cv_file, job?) -> CVAnalysis`
- [ ] Function `generateQuestions(job, candidate?, rubric?, level) -> Question[]`
- [ ] Function `suggestFollowUp(transcript_window, context) -> Suggestion`
- [ ] Function `scoreAnswer(question, answer_segments, rubric_criterion) -> Score`
- [ ] Function `generateReport(interview, scores, transcript) -> Report`
- [ ] Function `mockInterviewQuestion(...)` + `mockFeedback(...)`
- [ ] Chuẩn input chung (job/candidate/transcript/rubric/...)
- [ ] Chuẩn output JSON có `evidence`, `confidence`, `insufficient_data`
- [ ] **DoD:** Một entry point gọi AI thống nhất, output có schema validate được

---

### K-S3-02: Prompt Template Manager `[Khó]`

> Cần đảm bảo bảng `ai_prompt_templates` đã được thêm vào DATABASE_DESIGN.md (xem K-S0-01).

- [ ] Migration bảng `ai_prompt_templates`
- [ ] Hàm `loadTemplate(name, version)` (default version = latest active)
- [ ] Variable substitution `{{var}}` an toàn (escape nếu render vào prompt LLM)
- [ ] Versioning: không sửa version cũ, tạo version mới
- [ ] Seed prompt cho: analyze_jd, analyze_cv, generate_questions, suggest_follow_up, score_answer, generate_report, mock_question, mock_feedback
- [ ] **DoD:** Đổi prompt không cần redeploy, có version để rollback

---

### K-S3-03: AI Request Log `[Trung bình]`

> Cần đảm bảo bảng `ai_request_logs` đã được thêm vào DATABASE_DESIGN.md (xem K-S0-01).

- [ ] Migration bảng `ai_request_logs`
- [ ] Log mỗi AI call: `template_id`, `template_version`, input/output, latency, status, error, tokens, cost
- [ ] Truy vấn được theo `interview_id` / `job_id` / `candidate_id`
- [ ] Không log PII thừa (mask password, token nếu có)
- [ ] API admin xem log: `GET /admin/ai-logs?interview_id=...`
- [ ] **DoD:** Có log đủ để debug AI và audit cost

---

### K-S3-04: AI Error/Retry Handling `[Khó]`

- [ ] Retry policy (exponential backoff, max N lần)
- [ ] Timeout cho AI call (config được)
- [ ] Circuit breaker khi AI provider lỗi liên tục
- [ ] Fallback graceful: trả `insufficient_data: true` thay vì 500
- [ ] Map lỗi AI → error code `AI_SERVICE_ERROR` (502)
- [ ] Không làm crash backend khi AI provider lỗi
- [ ] **DoD:** AI lỗi được handle gọn, không sập backend

---

## SPRINT 4 — JD / CV / QUESTION AI

### K-S4-01: Analyze JD `[Khó]`

> Endpoint theo [API_SPEC.md](../API_SPEC.md) mục 4.6 và 9.

- [ ] `POST /companies/:company_id/jobs/:job_id/analyze` — body `{ force_refresh }`
- [ ] Output schema: `summary`, `required_skills[]`, `nice_to_have_skills[]`, `suggested_rubric[]`, `suggested_questions[]`
- [ ] Cache vào `jobs.ai_analysis_json` + `jobs.ai_summary`
- [ ] Validate output schema trước khi lưu
- [ ] **DoD:** Trả đúng schema, lưu cache, có `force_refresh`

---

### K-S4-02: Analyze CV `[Khó]`

- [ ] `POST /companies/:company_id/candidates/:candidate_id/parse-cv`
- [ ] Output: `summary`, `skills[]`, `experience_years`, `education[]`, `projects[]`
- [ ] Lưu vào `candidates.parsed_cv_json` + `ai_cv_summary`
- [ ] Nếu có `job_id`, sinh `ai_match_json` cho `job_candidates` (matched/missing requirements, focus points)
- [ ] **DoD:** CV parse đúng, có match score nếu có job

---

### K-S4-03: Generate Questions `[Khó]`

> Theo [API_SPEC.md](../API_SPEC.md) mục 9.1.

- [ ] `POST /ai/generate-questions` (hoặc scope theo company tùy chốt với Lai) — body `{ job_id, candidate_id?, rubric_id?, count, difficulty?, types[] }`
- [ ] Loại câu hỏi: introduction, technical, behavioral, situational, experience, culture, follow_up, closing
- [ ] Mỗi câu hỏi: `question_text`, `question_type`, `skill_tags`, `level`, `expected_signals`, `reason`, `is_ai_generated=true`
- [ ] Tránh trùng: check với `question_bank` của job
- [ ] **DoD:** Câu hỏi đa dạng, có metadata, không trùng

---

### K-S4-04: Save Question Bank `[Trung bình]`

> Lưu ý: bảng tên là `question_bank` (KHÔNG phải `interview_questions`).

- [ ] Insert câu hỏi AI sinh vào `question_bank`
- [ ] API `GET /companies/:company_id/question-bank?job_id=&type=&level=`
- [ ] API `POST /companies/:company_id/question-bank` — recruiter thêm câu hỏi thủ công
- [ ] API `PUT/DELETE` từng câu hỏi
- [ ] Gán question vào `interview_templates.config_json` hoặc bảng nối nếu cần
- [ ] **DoD:** Question bank tái sử dụng giữa các interview của cùng job/company

---

## SPRINT 5 — SCORING + REPORT

### K-S5-01: Rubric Scoring Engine `[Rất khó]`

> Theo [API_SPEC.md](../API_SPEC.md) mục 9.3. Dùng bảng `rubrics`, `rubric_criteria`, `interview_scores`.

- [ ] CRUD rubric: `GET/POST/PUT/DELETE /companies/:company_id/rubrics`, kèm `rubric_criteria`
- [ ] `POST /ai/score-answer` — body `{ interview_id, criterion_id, transcript_segments }`
- [ ] Output: `criterion_name`, `score`, `max_score`, `weight`, `weighted_score`, `evidence`, `ai_comment`, `confidence`, `status`, `scored_by`
- [ ] Snapshot `criterion_name` và `weight` vào `interview_scores` (rubric đổi sau không ảnh hưởng score cũ)
- [ ] Tính `final_score` tổng theo weight
- [ ] Unique `(interview_id, criterion_name)` cho score cuối
- [ ] Emit event để Hùng đẩy realtime cho recruiter
- [ ] **DoD:** Score có weight đúng, lưu snapshot, có endpoint Hùng dùng

---

### K-S5-02: Score Evidence Extraction `[Rất khó]`

- [ ] Evidence trích từ `interview_transcripts` thật, không bịa (verify substring)
- [ ] Đánh dấu `status=insufficient_evidence` nếu không tìm được evidence
- [ ] Đánh dấu `status=manual_override` khi recruiter sửa
- [ ] Không score khi thiếu dữ liệu nghiêm trọng
- [ ] Có `confidence` (0-1)
- [ ] **DoD:** Evidence luôn có nguồn từ transcript, recruiter sửa được

---

### K-S5-03: Report Generator `[Rất khó]`

> Theo [API_SPEC.md](../API_SPEC.md) mục 9.4 và 10. Dùng bảng `interview_reports` (1-1 với interview).

- [ ] Trigger sau khi end interview (K-S2-03)
- [ ] `POST /ai/generate-report` (internal, hoặc gọi qua orchestrator)
- [ ] `GET /companies/:company_id/interviews/:interview_id/report`
- [ ] Output `report_json`: `summary`, `final_score`, `recommendation` (strong_hire/hire/consider/next_round/reject/insufficient_data), `strengths[]`, `weaknesses[]`, `risks[]`, `evidence_json`, `ai_reasoning_summary`
- [ ] Field tách rời cho recruiter: `recruiter_decision`, `recruiter_comment`
- [ ] Lưu unique `interview_id` (1 interview = 1 report)
- [ ] **DoD:** Report sinh tự động sau end, recruiter xem và override decision được

---

### K-S5-04: Report Retry/Status + Notification `[Khó]`

> Dùng bảng `notifications`.

- [ ] Field `report_status` trong `interview_reports` hoặc dùng `generated_by` + state riêng: `pending` / `generating` / `ready` / `failed`
- [ ] Retry tự động khi failed (max N lần, exponential backoff)
- [ ] `POST /companies/:company_id/interviews/:interview_id/report/retry` — retry thủ công
- [ ] `PUT /companies/:company_id/interviews/:interview_id/report/decision` — recruiter cập nhật `recruiter_decision`, `recruiter_comment`
- [ ] Insert notification cho recruiter khi report `ready` hoặc `failed`
- [ ] `GET /notifications` — list notification của user, đánh dấu `read_at`
- [ ] **DoD:** Recruiter nhận thông báo khi report sẵn sàng, có thể retry khi failed

---

## SPRINT 6 — MOCK INTERVIEW BACKEND

> Lưu ý: dùng `user_id` (KHÔNG phải `candidate_id`). Mock interview là tính năng cá nhân.

### K-S6-01: Mock Interview Session API `[Khó]`

- [ ] `POST /mock-interviews` — body `{ target_role, target_level, cv_file_id? }`
- [ ] `GET /mock-interviews/me` — list mock của user hiện tại
- [ ] `GET /mock-interviews/:id` — chi tiết
- [ ] `POST /mock-interviews/:id/start` — set `status=active`, `started_at`
- [ ] `POST /mock-interviews/:id/end` — set `status=completed`, `ended_at`, sinh feedback tổng
- [ ] **DoD:** Candidate (user) tạo và chạy session mock

---

### K-S6-02: Mock AI Question Flow `[Khó]`

- [ ] `POST /mock-interviews/:id/messages` — gửi câu trả lời, AI sinh câu hỏi tiếp
- [ ] Lưu vào `mock_interview_messages` (KHÔNG phải `mock_interview_answers`) với `sender_type` ai/candidate/system
- [ ] AI hỏi follow-up theo context
- [ ] Giới hạn số câu hỏi mỗi session
- [ ] **DoD:** AI hỏi liền mạch, có context

---

### K-S6-03: Mock Answer Scoring `[Khó]`

- [ ] Score từng câu trả lời (tái dùng scoring engine K-S5-01 nhưng dùng rubric mặc định cho target_role)
- [ ] Lưu `score_json` trong `mock_interview_messages`
- [ ] Có evidence + confidence
- [ ] **DoD:** Candidate nhận feedback ngay sau câu trả lời

---

### K-S6-04: Mock Report API `[Trung bình]`

- [ ] `GET /mock-interviews/:id/report` — tổng hợp `feedback_json` + `final_score` từ `mock_interviews`
- [ ] Điểm mạnh, điểm cần cải thiện, gợi ý
- [ ] **DoD:** Candidate xem được report mock

---

## XUYÊN SUỐT — AUDIT LOG

> Dùng bảng `audit_logs`. Theo [API_SPEC.md](../API_SPEC.md) mục 13.

- [ ] Log: login/logout
- [ ] Log: tạo/sửa/xóa job (`before_json`/`after_json`)
- [ ] Log: thêm/sửa/xóa candidate + upload CV
- [ ] Log: tạo/cancel/reschedule lịch phỏng vấn
- [ ] Log: start/end interview
- [ ] Log: xem/tải report
- [ ] Log: thay đổi `recruiter_decision`
- [ ] Log: invite/remove company member
- [ ] Mỗi log có: `actor_user_id`, `actor_role`, `action`, `resource_type`, `resource_id`, `before_json`, `after_json`, `ip_address`, `user_agent`, `created_at`
- [ ] Không log password/token
- [ ] `GET /companies/:company_id/audit-logs?resource_type=&actor_user_id=&page=&page_size=` (admin/owner)

---

## DEPENDENCY — CUNG CẤP CHO HÙNG

- [ ] Interview API (CRUD + start/end/cancel)
- [ ] Room access token endpoint (mục 7)
- [ ] Transcript storage API (`POST /companies/:company_id/interviews/:id/transcripts`)
- [ ] AI suggestion endpoint (mục 9.2 — `POST /ai/suggest-follow-up`)
- [ ] AI scoring endpoint (mục 9.3 — `POST /ai/score-answer`)
- [ ] Report generation trigger
- [ ] Event contract: room.opened, room.closed, transcript.created, score.created

## DEPENDENCY — NHẬN TỪ HÙNG

- [ ] Event payload Hùng cần emit (room status, participant join/leave, transcript chunk)
- [ ] Format `connection_state` thật của participant
- [ ] Bộ test fixture transcript để Khôi test scoring

## DEPENDENCY — CUNG CẤP CHO LAI

- [ ] API Company (CRUD + list)
- [ ] API Auth (`/auth/me` đặc biệt cho hiển thị user)
- [ ] API Job (CRUD + analyze)
- [ ] API Candidate (CRUD + upload CV + parse)
- [ ] API Interview (CRUD + start/end + transcript)
- [ ] API Rubric (CRUD)
- [ ] API Question Bank (CRUD)
- [ ] API Report (read + update decision)
- [ ] API Mock Interview (session + messages + report)
- [ ] API File (upload + signed URL)
- [ ] API Notification (list + mark read)
- [ ] OpenAPI spec hoặc Postman collection

## DEPENDENCY — NHẬN TỪ LAI

- [ ] UI yêu cầu field nào trong response (đặc biệt list view: count, summary)
- [ ] Form validation rules cần backend đảm bảo
- [ ] Dashboard thống kê cần endpoint nào
- [ ] Report UI cần data shape nào
- [ ] Filter/sort UI cần parameter nào

---

## CHECKLIST TEST CUỐI CỦA KHÔI

- [ ] Test: Auth register/login/logout/refresh
- [ ] Test: `GET /auth/me` trả đúng companies + role
- [ ] Test: Token hết hạn → 401 chuẩn
- [ ] Test: Recruiter không xem được company khác (cross-tenant leak)
- [ ] Test: Candidate không xem được report nội bộ company
- [ ] Test: Company CRUD + member tự động owner khi tạo
- [ ] Test: Job CRUD + scope company + soft delete
- [ ] Test: Candidate CRUD + assign job (`job_candidates` unique)
- [ ] Test: CV upload → file metadata vào `files`, không lưu blob
- [ ] Test: Tạo interview → có `interview_rooms` + `invite_token`
- [ ] Test: Start/end interview đúng quyền + state machine
- [ ] Test: End interview trigger report generation
- [ ] Test: Transcript push từ Hùng lưu đúng schema
- [ ] Test: Analyze JD đúng schema, lưu vào `ai_analysis_json`
- [ ] Test: Parse CV đúng schema, lưu vào `parsed_cv_json`
- [ ] Test: Generate questions không trùng, lưu vào `question_bank`
- [ ] Test: Score có evidence từ transcript thật (không bịa)
- [ ] Test: Score `insufficient_evidence` khi thiếu data
- [ ] Test: Report sinh sau end interview, unique 1-1
- [ ] Test: Report retry khi failed, notification gửi đúng user
- [ ] Test: Recruiter override `recruiter_decision`
- [ ] Test: AI provider lỗi → `AI_SERVICE_ERROR` 502, không sập backend
- [ ] Test: Mock interview chạy với `user_id`, lưu vào `mock_interview_messages`
- [ ] Test: Audit log lưu đúng cho mọi action quan trọng
- [ ] Test: Response không leak `password_hash`, `invite_token_hash`

---

## BÀN GIAO CUỐI

- [ ] Backend core ổn định, response/error envelope chuẩn
- [ ] DATABASE_DESIGN.md được cập nhật bổ sung `ai_prompt_templates` + `ai_request_logs`
- [ ] Migration đầy đủ 22 bảng, chạy up/down sạch
- [ ] Auth/RBAC/Company scope hoạt động đúng
- [ ] Company/Job/Candidate/Interview API khớp API_SPEC.md
- [ ] AI Orchestrator + Prompt Manager + Request Log
- [ ] JD/CV analysis có cache + force_refresh
- [ ] Question generation lưu vào `question_bank`
- [ ] Scoring engine có evidence + confidence + snapshot weight
- [ ] Report generator + retry + notification
- [ ] Mock interview backend
- [ ] Audit log đầy đủ
- [ ] OpenAPI/Postman collection cho Hùng + Lai
- [ ] Test cơ bản pass cho mọi luồng quan trọng
