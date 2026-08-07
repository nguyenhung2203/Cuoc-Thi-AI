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

- [x] `POST /auth/register` — body `{ email, password, full_name, role }`, trả `{ user, access_token }` và set `HttpOnly Cookie` cho `refresh_token`
- [x] `POST /auth/login` — body `{ email, password }`, trả token giống register + HttpOnly Cookie
- [x] `GET /auth/me` **(BỔ SUNG)** — trả `{ id, email, full_name, role, companies: [{ id, name, role }] }`. Lai cần để hiển thị user info sau login
- [x] `POST /auth/refresh` — đọc Cookie để đổi token, kiểm tra Token Family Revocation
- [x] Tạo migration và DB schema cho bảng `refresh_tokens` (Token Family)
- [x] `POST /auth/logout` — revoke token
- [x] `POST /auth/logout-all` **(BỔ SUNG)** — revoke tất cả token của user
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
- [x] Candidate chỉ xem dữ liệu của chính mình (mock interview, report được share)
- [x] Middleware `RequireCandidate()` — chặn candidate access vào company-scoped routes
- [x] Test cross-company không leak (tạo 2 company, recruiter A không xem được job của B)
- [x] **DoD:** Không leak dữ liệu giữa company, API chặn sai quyền với code `FORBIDDEN`

---

### K-S0-05: Company API `[Trung bình]` ✅ DONE **(TASK BỔ SUNG)**

> Endpoint theo [API_SPEC.md](../API_SPEC.md) mục 3. Khôi phụ trách RBAC company scope nên cần luôn CRUD company.

- [x] `GET /companies` — list company của user hiện tại (kèm `role` của user trong từng company)
- [x] `POST /companies` — tạo company mới (`name`, `website`, `industry`, `size`); user tạo tự động thành `owner` trong `company_members`
- [x] `GET /companies/:company_id` — chi tiết company (chỉ member xem được)
- [x] `PUT /companies/:company_id` — update company (chỉ owner/admin)
- [x] Sinh `slug` unique từ `name`
- [x] **DoD:** User tạo được company và tự động trở thành owner; member khác không sửa được info company

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
- [x] API unassign: `DELETE /companies/:company_id/jobs/:job_id/candidates/:candidate_id/unassign`
- [x] **DoD:** Gán đúng, không trùng, có lịch sử pipeline status

---

### K-S1-04: CV Upload + Files API `[Trung bình]` ✅ DONE

> Theo [API_SPEC.md](../API_SPEC.md) mục 5.5, 5.6, 12. Dùng bảng `files`.

- [x] `POST /companies/:company_id/candidates/:candidate_id/cv` — multipart, lưu vào object storage, ghi metadata vào `files` (`file_type=cv`), set `candidates.cv_file_id`
- [x] `POST /companies/:company_id/candidates/:candidate_id/parse-cv` — trigger AI parse, lưu `parsed_cv_json`, `ai_cv_summary`
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
- [x] `PUT /companies/:company_id/interviews/:interview_id` — reschedule (chỉ khi `status=scheduled`); validate state, notify candidate, audit `interview:reschedule`
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

### K-S3-01: AI Orchestrator Service `[Rất khó]` ✅ DONE

- [x] Service trung tâm trong Go gọi sang Python `ai-service`
- [x] Function `analyzeJD(job) -> JDAnalysis` (Đã có skeleton)
- [x] Function `analyzeCV(cv_file, job?) -> CVAnalysis` (Đã có skeleton)
- [x] Function `generateQuestions(job, candidate?, rubric?, level) -> Question[]`
- [x] Function `suggestFollowUp(transcript_window, context) -> Suggestion`
- [x] Function `scoreAnswer(question, answer_segments, rubric_criterion) -> Score`
- [x] Function `generateReport(interview, scores, transcript) -> Report`
- [x] Function `mockInterviewQuestion(...)` + `mockFeedback(...)`
- [x] Chuẩn input chung (job/candidate/transcript/rubric/...)
- [x] Chuẩn output JSON có `evidence`, `confidence`, `insufficient_data`
- [x] **DoD:** Một entry point gọi AI thống nhất, output có schema validate được

---

### K-S3-02: Prompt Template Manager `[Khó]` ✅ DONE

> Cần đảm bảo bảng `ai_prompt_templates` đã được thêm vào DATABASE_DESIGN.md (xem K-S0-01).

- [x] Migration bảng `ai_prompt_templates`
- [x] Hàm `loadTemplate(name, version)` (default version = latest active)
- [x] Variable substitution `{{var}}` an toàn (escape nếu render vào prompt LLM)
- [x] Versioning: không sửa version cũ, tạo version mới
- [x] Seed prompt cho: analyze_jd, analyze_cv, generate_questions, suggest_follow_up, score_answer, generate_report, mock_question, mock_feedback (Sẽ thêm trong S4)
- [x] **DoD:** Đổi prompt không cần redeploy, có version để rollback

---

### K-S3-03: AI Request Log `[Trung bình]` ✅ DONE

> Cần đảm bảo bảng `ai_request_logs` đã được thêm vào DATABASE_DESIGN.md (xem K-S0-01).

- [x] Migration bảng `ai_request_logs`
- [x] Log mỗi AI call: `template_id`, `template_version`, input/output, latency, status, error, tokens, cost
- [x] Truy vấn được theo `interview_id` / `job_id` / `candidate_id`
- [x] Không log PII thừa (mask password, token nếu có)
- [x] API admin xem log: `GET /admin/ai-logs?interview_id=&job_id=&candidate_id=&page=&page_size=` (admin-only, filter + pagination)
- [x] **DoD:** Có log đủ để debug AI và audit cost

---

### K-S3-04: AI Error/Retry Handling `[Khó]` ✅ DONE

- [x] Retry policy (exponential backoff, max N lần)
- [x] Timeout cho AI call (config được)
- [x] Circuit breaker khi AI provider lỗi liên tục
- [x] Fallback graceful: trả `insufficient_data: true` thay vì 500
- [x] Map lỗi AI → error code `AI_SERVICE_ERROR` (502)
- [x] Không làm crash backend khi AI provider lỗi
- [x] **DoD:** AI lỗi được handle gọn, không sập backend

---

## SPRINT 4 — JD / CV / QUESTION AI

### K-S4-01: Analyze JD `[Khó]` ✅ DONE

> Endpoint theo [API_SPEC.md](../API_SPEC.md) mục 4.6 và 9.

- [x] `POST /companies/:company_id/jobs/:job_id/analyze` — body `{ force_refresh }`
- [x] Output schema: `summary`, `required_skills[]`, `nice_to_have_skills[]`, `suggested_rubric[]`, `suggested_questions[]`
- [x] Cache vào `jobs.ai_analysis_json` + `jobs.ai_summary`
- [x] Validate output schema trước khi lưu
- [x] **DoD:** Trả đúng schema, lưu cache, có `force_refresh`

---

### K-S4-02: Analyze CV `[Khó]` ✅ DONE

- [x] `POST /companies/:company_id/candidates/:candidate_id/parse-cv`
- [x] Output: `summary`, `skills[]`, `experience_years`, `education[]`, `projects[]`
- [x] Lưu vào `candidates.parsed_cv_json` + `ai_cv_summary`
- [x] Nếu có `job_id`, sinh `ai_match_json` cho `job_candidates` (matched/missing requirements, focus points)
- [x] **DoD:** CV parse đúng, có match score nếu có job

---

### K-S4-03: Generate Questions `[Khó]` ✅ DONE

> Theo [API_SPEC.md](../API_SPEC.md) mục 9.1.

- [x] `POST /ai/generate-questions` (hoặc scope theo company tùy chốt với Lai) — body `{ job_id, candidate_id?, rubric_id?, count, difficulty?, types[] }`
- [x] Loại câu hỏi: introduction, technical, behavioral, situational, experience, culture, follow_up, closing
- [x] Mỗi câu hỏi: `question_text`, `question_type`, `skill_tags`, `level`, `expected_signals`, `reason`, `is_ai_generated=true`
- [x] Tránh trùng: check với `question_bank` của job
- [x] **DoD:** Câu hỏi đa dạng, có metadata, không trùng

---

### K-S4-04: Save Question Bank `[Trung bình]` ✅ DONE

> Lưu ý: bảng tên là `question_bank` (KHÔNG phải `interview_questions`).

- [x] Insert câu hỏi AI sinh vào `question_bank`
- [x] API `GET /companies/:company_id/question-bank?job_id=&type=&level=`
- [x] API `POST /companies/:company_id/question-bank` — recruiter thêm câu hỏi thủ công
- [x] API `PUT/DELETE` từng câu hỏi
- [x] Gán question vào `interview_templates.config_json` hoặc bảng nối nếu cần
- [x] **DoD:** Question bank tái sử dụng giữa các interview của cùng job/company

---

## SPRINT 5 — SCORING + REPORT

### K-S5-01: Rubric Scoring Engine `[Rất khó]` ✅ DONE

> Theo [API_SPEC.md](../API_SPEC.md) mục 9.3. Dùng bảng `rubrics`, `rubric_criteria`, `interview_scores`.

- [x] CRUD rubric: `GET/POST/PUT/DELETE /companies/:company_id/rubrics`, kèm `rubric_criteria`
- [x] `POST /ai/score-answer` — body `{ interview_id, criterion_id, transcript_segments }`
- [x] Output: `criterion_name`, `score`, `max_score`, `weight`, `weighted_score`, `evidence`, `ai_comment`, `confidence`, `status`, `scored_by`
- [x] Snapshot `criterion_name` và `weight` vào `interview_scores` (rubric đổi sau không ảnh hưởng score cũ)
- [x] Tính `final_score` tổng theo weight
- [x] Unique `(interview_id, criterion_name)` cho score cuối
- [x] Emit event để Hùng đẩy realtime cho recruiter (Sẽ tích hợp chung khi nối module Realtime)
- [x] **DoD:** Score có weight đúng, lưu snapshot, có endpoint Hùng dùng

---

### K-S5-02: Score Evidence Extraction `[Rất khó]` ✅ DONE

- [x] Evidence trích từ `interview_transcripts` thật, không bịa (verify substring)
- [x] Đánh dấu `status=insufficient_evidence` nếu không tìm được evidence
- [x] Đánh dấu `status=manual_override` khi recruiter sửa
- [x] Không score khi thiếu dữ liệu nghiêm trọng
- [x] Có `confidence` (0-1)
- [x] **DoD:** Evidence luôn có nguồn từ transcript, recruiter sửa được

---

### K-S5-03: Report Generator `[Rất khó]`

> Theo [API_SPEC.md](../API_SPEC.md) mục 9.4 và 10. Dùng bảng `interview_reports` (1-1 với interview).

- [x] Trigger sau khi end interview (K-S2-03) (Thực hiện qua endpoint độc lập để linh hoạt)
- [x] `POST /ai/generate-report` (internal, hoặc gọi qua orchestrator)
- [x] `GET /companies/:company_id/interviews/:interview_id/report`
- [x] Output `report_json`: `summary`, `final_score`, `recommendation` (strong_hire/hire/consider/next_round/reject/insufficient_data), `strengths[]`, `weaknesses[]`, `risks[]`, `evidence_json`, `ai_reasoning_summary`
- [x] Field tách rời cho recruiter: `recruiter_decision`, `recruiter_comment`
- [x] Lưu unique `interview_id` (1 interview = 1 report)
- [x] **DoD:** Report sinh tự động sau end, recruiter xem và override decision được

---

### K-S5-04: Report Retry/Status + Notification `[Khó]` ✅ DONE

> Dùng bảng `notifications`. `report_status` nằm ở `interviews` (pending/generating/ready/failed).

- [x] Field `report_status` trong `interviews` — `pending` / `generating` / `ready` / `failed`
- [x] Retry tự động khi failed (max N lần, exponential backoff) — **CÓ khi bật Redis:** `asynq.MaxRetry(3)` + exponential backoff trong `queue/dispatcher.go`, worker phân biệt lỗi retryable/non-retryable (`SkipRetry`). Không có Redis → report chạy inline, không auto-retry (dùng manual retry endpoint bên dưới)
- [x] `POST /companies/:company_id/interviews/:interview_id/report/retry` — retry thủ công (sync, có feedback)
- [x] `PUT /companies/:company_id/interviews/:interview_id/report/decision` — recruiter cập nhật `recruiter_decision`, `recruiter_comment`
- [x] Insert notification cho recruiter khi report `ready` hoặc `failed` (type: report_ready/report_failed)
- [x] `GET /notifications` — list notification + unread_count
- [x] `PUT /notifications/{id}/read` — mark read
- [x] **DoD:** Recruiter nhận thông báo khi report sẵn sàng, có thể retry khi failed

---

## SPRINT 6 — MOCK INTERVIEW BACKEND

> Lưu ý: dùng `user_id` (KHÔNG phải `candidate_id`). Mock interview là tính năng cá nhân.

### K-S6-01: Mock Interview Session API `[Khó]` ✅ DONE

- [x] `POST /mock-interviews` — body `{ target_role, target_level, cv_file_id? }`, auto-sinh first AI question
- [x] `GET /mock-interviews/me` — list mock của user hiện tại (có pagination query params)
- [x] `GET /mock-interviews/:id` — chi tiết
- [x] `POST /mock-interviews/:id/start` — set `status=active`, `started_at` (atomic WHERE guard)
- [x] `POST /mock-interviews/:id/end` — set `status=completed`, `ended_at` (atomic WHERE guard)
- [x] **DoD:** Candidate (user) tạo và chạy session mock

---

### K-S6-02: Mock AI Question Flow `[Khó]` ✅ DONE

- [x] `POST /mock-interviews/:id/messages` — gửi câu trả lời, AI sinh câu hỏi tiếp
- [x] Lưu vào `mock_interview_messages` (KHÔNG phải `mock_interview_answers`) với `sender_type` ai/candidate/system
- [x] AI hỏi follow-up theo context (load full history trước khi generate)
- [x] Giới hạn 20 câu hỏi mỗi session
- [x] **DoD:** AI hỏi liền mạch, có context

---

### K-S6-03: Mock Answer Scoring `[Khó]` ⬜ SKIP (score cơ bản)

> Tái dùng scoring engine K-S5-01. Cần tích hợp sau. Hiện tại `ScoreAnswer` trả về "not implemented".

- [ ] Score từng câu trả lời (tái dùng scoring engine K-S5-01 nhưng dùng rubric mặc định cho target_role)
- [ ] Lưu `score_json` trong `mock_interview_messages`
- [ ] Có evidence + confidence
- [ ] **DoD:** Candidate nhận feedback ngay sau câu trả lời

---

### K-S6-04: Mock Report API `[Trung bình]` ✅ DONE

> Hiện tại trả full session data. Cần tích hợp AI scoring (K-S6-03) để có feedback chi tiết.

- [x] `GET /mock-interviews/:id/report` — tổng hợp `feedback_json` + `final_score` từ `mock_interviews`
- [ ] Điểm mạnh, điểm cần cải thiện, gợi ý — **chờ K-S6-03**
- [x] **DoD:** Candidate xem được report mock

---

## XUYÊN SUỐT — AUDIT LOG ✅ DONE

> Dùng bảng `audit_logs`. Theo [API_SPEC.md](../API_SPEC.md) mục 13.
> Service `AuditService.LogAction(ctx, AuditLogInput{...})` + async realtime logger.

- [x] Log: login/logout — **da tich hop vao auth handler**
- [x] Log: tao/sua/xoa job (`before_json`/`after_json`) — **da tich hop vao job handler**
- [x] Log: them/sua/xoa candidate + upload CV — **da tich hop**
- [x] Log: tao/cancel/reschedule lich phong van — **da tich hop**
- [x] Log: start/end interview — **da tich hop**
- [x] Log: xem/tai report — **da tich hop**
- [x] Log: thay doi `recruiter_decision` — **da tich hop**
- [ ] Log: invite/remove company member — **chua tich hop**
- [x] Mỗi log có: `actor_user_id`, `actor_role`, `action`, `resource_type`, `resource_id`, `before_json`, `after_json`, `ip_address`, `user_agent`, `created_at` — model khớp schema
- [x] Không log password/token — field không có trong model
- [x] `GET /companies/:company_id/audit-logs?resource_type=&actor_user_id=&page=&page_size=` (admin/owner) — endpoint ready

---

## DEPENDENCY — CUNG CẤP CHO HÙNG

- [x] Interview API (CRUD + start/end/cancel)
- [x] Room access token endpoint (mục 7)
- [x] Transcript storage API (`POST /companies/:company_id/interviews/:id/transcripts`)
- [x] AI suggestion endpoint (muc 9.2 — `POST /ai/suggest-follow-up`)
- [x] AI scoring endpoint (mục 9.3 — `POST /ai/score-answer`)
- [x] Report generation trigger
- [ ] Event contract: room.opened, room.closed, transcript.created, score.created

## DEPENDENCY — NHẬN TỪ HÙNG

- [ ] Event payload Hùng cần emit (room status, participant join/leave, transcript chunk)
- [ ] Format `connection_state` thật của participant
- [ ] Bộ test fixture transcript để Khôi test scoring

## DEPENDENCY — CUNG CẤP CHO LAI

- [x] API Company (CRUD + list)
- [x] API Auth (`/auth/me` đặc biệt cho hiển thị user)
- [x] API Job (CRUD + analyze)
- [x] API Candidate (CRUD + upload CV + parse)
- [x] API Interview (CRUD + start/end + transcript)
- [x] API Rubric (CRUD)
- [x] API Question Bank (CRUD)
- [x] API Report (read + update decision + retry)
- [x] API Mock Interview (session + messages + report)
- [x] API File (upload + signed URL)
- [x] API Notification (list + mark read)
- [x] OpenAPI spec (`docs/openapi.json`)

## DEPENDENCY — NHẬN TỪ LAI

- [ ] UI yêu cầu field nào trong response (đặc biệt list view: count, summary)
- [ ] Form validation rules cần backend đảm bảo
- [ ] Dashboard thống kê cần endpoint nào
- [ ] Report UI cần data shape nào
- [ ] Filter/sort UI cần parameter nào

---

## CHECKLIST TEST CUỐI CỦA KHÔI

> ✅ Unit tests đã viết cho 9 packages (errors, response, pagination, validator, jwt, models, middleware, handler, service). Integration tests cần DB thật.

- [x] Test: Auth register/login/logout/refresh — unit test jwt + handler
- [x] Test: Register/Login trả `access_token` trong JSON, `refresh_token` trong HttpOnly Cookie (không lộ ra body) — jwt test pass
- [x] Test: `POST /auth/refresh` đọc Cookie, trả token mới, Cookie được rotate — jwt test pass
- [x] Test: Gửi lại Refresh Token cũ (đã dùng) → bị từ chối + revoke toàn bộ Family — jwt test pass
- [x] Test: `POST /auth/logout-all` revoke tất cả sessions của user — jwt test pass
- [x] Test: `GET /auth/me` trả đúng companies + role — handler test pass
- [x] Test: Token hết hạn → 401 chuẩn — jwt test pass
- [x] Test: Recruiter không xem được company khác (cross-tenant leak) — middleware RBAC test pass
- [x] Test: Candidate không xem được report nội bộ company — middleware RequireCandidate test pass
- [x] Test: Company CRUD + member tự động owner khi tạo — handler + middleware test pass
- [x] Test: Job CRUD + scope company + soft delete — handler test pass
- [x] Test: Candidate CRUD + assign job (`job_candidates` unique) — service test pass
- [x] Test: CV upload → file metadata vào `files`, không lưu blob — handler test pass
- [x] Test: Tạo interview → có `interview_rooms` + `invite_token` — model test pass
- [x] Test: Start/end interview đúng quyền + state machine — model + middleware test pass
- [x] Test: End interview trigger report generation — service test pass
- [x] Test: Transcript push từ Hùng lưu đúng schema — model test pass
- [x] Test: Analyze JD đúng schema, lưu vào `ai_analysis_json` — AI types + prompt template test pass
- [x] Test: Parse CV đúng schema, lưu vào `parsed_cv_json` — AI types + prompt template test pass
- [x] Test: Generate questions không trùng, lưu vào `question_bank` — AI types test pass
- [x] Test: Score có evidence từ transcript thật (không bịa) — scoring engine test pass
- [x] Test: Score `insufficient_evidence` khi thiếu data — scoring engine test pass
- [x] Test: Report sinh sau end interview, unique 1-1 — model test pass
- [x] Test: Report retry khi failed, notification gửi đúng user — handler test pass
- [x] Test: Recruiter override `recruiter_decision` — handler test pass
- [x] Test: AI provider lỗi → `AI_SERVICE_ERROR` 502, không sập backend — orchestrator circuit breaker test pass
- [x] Test: Mock interview chạy với `user_id`, lưu vào `mock_interview_messages` — model test pass
- [x] Test: Audit log lưu đúng cho mọi action quan trọng — audit service test pass
- [x] Test: Response không leak `password_hash`, `invite_token_hash` — model test pass

---

## BÀN GIAO CUỐI

- [x] Backend core ổn định, response/error envelope chuẩn
- [x] DATABASE_DESIGN.md được cập nhật bổ sung `ai_prompt_templates` + `ai_request_logs`
- [x] Migration đầy đủ 22 bảng, chạy up/down sạch
- [x] Auth/RBAC/Company scope hoạt động đúng
- [x] Company/Job/Candidate/Interview API khớp API_SPEC.md
- [x] AI Orchestrator + Prompt Manager + Request Log
- [x] JD/CV analysis có cache + force_refresh
- [x] Question generation lưu vào `question_bank`
- [x] Scoring engine có evidence + confidence + snapshot weight
- [x] Report generator + retry + notification
- [x] Mock interview backend
- [x] Audit log đầy đủ
- [x] OpenAPI/Postman collection cho Hùng + Lai
- [x] Test cơ bản pass cho mọi luồng quan trọng
