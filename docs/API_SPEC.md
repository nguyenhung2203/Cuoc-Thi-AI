# API_SPEC.md
# Đặc tả API — AI Interview Platform

## 0. Mục đích tài liệu

Tài liệu này là hợp đồng giao tiếp giữa frontend, backend, realtime và AI service.

AI IDE của Hùng, Khôi, Lai phải bám theo file này để tránh lệch endpoint, field, response và permission.

- Hùng dùng cho realtime room, transcript, AI realtime.
- Khôi dùng cho backend API, DB, AI orchestration, report.
- Lai dùng cho UI gọi API và render đúng dữ liệu.

---

## 1. API Convention

## 1.1. Base URL

```text
/api/v1
```

Ví dụ:

```text
GET /api/v1/jobs
POST /api/v1/interviews/:id/start
```

---

## 1.2. Authentication

Dùng Bearer token:

```http
Authorization: Bearer <access_token>
```

Một số endpoint candidate join room có thể dùng invite token:

```text
/interviews/join/:invite_token
```

---

## 1.3. Response envelope chuẩn

Success:

```json
{
  "success": true,
  "data": {},
  "meta": {},
  "request_id": "req_123"
}
```

Error:

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Dữ liệu không hợp lệ",
    "details": []
  },
  "request_id": "req_123"
}
```

---

## 1.4. Pagination chuẩn

Request:

```text
?page=1&page_size=20&sort=created_at:desc
```

Response meta:

```json
{
  "page": 1,
  "page_size": 20,
  "total": 120,
  "total_pages": 6
}
```

---

## 1.5. Error code chuẩn

| Code | HTTP | Mô tả |
|---|---:|---|
| UNAUTHORIZED | 401 | Chưa đăng nhập hoặc token sai |
| FORBIDDEN | 403 | Không có quyền |
| NOT_FOUND | 404 | Không tìm thấy dữ liệu |
| VALIDATION_ERROR | 422 | Dữ liệu không hợp lệ |
| CONFLICT | 409 | Trùng hoặc xung đột trạng thái |
| RATE_LIMITED | 429 | Quá nhiều request |
| AI_SERVICE_ERROR | 502 | Lỗi AI provider |
| REALTIME_ERROR | 502 | Lỗi realtime service |
| INTERNAL_ERROR | 500 | Lỗi hệ thống |

---

## 2. Auth API

## 2.1. Register

```http
POST /auth/register
```

Request:

```json
{
  "email": "recruiter@example.com",
  "password": "Password123!",
  "full_name": "Nguyễn Văn A",
  "role": "recruiter"
}
```

Response:

```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "email": "recruiter@example.com",
      "full_name": "Nguyễn Văn A",
      "role": "recruiter",
      "status": "active"
    },
    "access_token": "string"
  }
}

*Lưu ý: `refresh_token` sẽ được trả về ngầm qua header `Set-Cookie`:*

```http
Set-Cookie: refresh_token=<TOKEN>; HttpOnly; Secure; SameSite=Lax; Path=/api/v1/auth; Max-Age=604800
```

*Giải thích các thuộc tính:*
- `HttpOnly` — Javascript không thể đọc Cookie, chống XSS.
- `Secure` — Chỉ gửi qua HTTPS.
- `SameSite=Lax` — Cho phép Cookie khi user click link từ email/Slack vào trang web (Strict sẽ block trường hợp này).
- `Path=/api/v1/auth` — Chỉ gửi Cookie khi gọi Auth endpoints, giảm attack surface.
- `Max-Age=604800` — 7 ngày, đồng bộ với `expires_at` trong DB.

*Nếu FE và BE khác domain (ví dụ FE ở `app.example.com`, BE ở `api.example.com`), cần thêm `Domain=.example.com`.*
```

Permission: public.

---

## 2.2. Login

```http
POST /auth/login
```

Request:

```json
{
  "email": "recruiter@example.com",
  "password": "Password123!"
}
```

Response: giống register.

---

## 2.3. Me

```http
GET /auth/me
```

Response:

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "email": "recruiter@example.com",
    "full_name": "Nguyễn Văn A",
    "role": "recruiter",
    "companies": [
      {
        "id": "uuid",
        "name": "ABC Company",
        "role": "owner"
      }
    ]
  }
}
```

Permission: authenticated.

---

## 2.4. Refresh token

```http
POST /auth/refresh
```

Request: Không cần body. Backend sẽ tự đọc `refresh_token` từ HttpOnly Cookie.

Response: Trả về `access_token` mới trong JSON body và tự động cập nhật `refresh_token` mới qua `Set-Cookie` (Token Rotation).

Lưu ý: Nếu phát hiện gửi lại Refresh Token cũ (đã sử dụng), hệ thống sẽ từ chối và **thu hồi toàn bộ phiên đăng nhập (Token Family Revocation)**.

---

## 2.5. Logout

```http
POST /auth/logout
```

Permission: authenticated.

*Lưu ý: Endpoint này sẽ vô hiệu hóa (revoke) Token Family hiện tại trong Database và trả về header `Set-Cookie` với Max-Age=-1 để xóa Cookie ở Frontend.*

---

## 2.6. Logout All Sessions

```http
POST /auth/logout-all
```

Permission: authenticated.

*Lưu ý: Endpoint này sẽ vô hiệu hóa (revoke) **tất cả** Token Family của user trong Database (tất cả thiết bị) và xóa Cookie ở thiết bị hiện tại. Dùng khi user nghi ngờ tài khoản bị xâm phạm.*

Response:

```json
{
  "success": true,
  "data": {
    "message": "All sessions revoked. Please login again.",
    "revoked_sessions": 3
  }
}
```

---

## 3. Company API

## 3.1. List companies

```http
GET /companies
```

Permission: authenticated recruiter/admin.

Response:

```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "name": "ABC Company",
      "slug": "abc-company",
      "role": "owner",
      "logo_url": null
    }
  ]
}
```

---

## 3.2. Create company

```http
POST /companies
```

Request:

```json
{
  "name": "ABC Company",
  "website": "https://example.com",
  "industry": "Software",
  "size": "11-50"
}
```

Permission: recruiter/admin.

---

## 3.3. Get company detail

```http
GET /companies/:company_id
```

Permission: company member.

---

## 3.4. Update company

```http
PUT /companies/:company_id
```

Permission: company owner/admin.

---

## 4. Job API

## 4.1. List jobs

```http
GET /companies/:company_id/jobs?status=open&keyword=frontend&page=1&page_size=20
```

Permission: company member.

Response item:

```json
{
  "id": "uuid",
  "title": "Frontend Developer",
  "department": "Engineering",
  "level": "middle",
  "status": "open",
  "candidate_count": 12,
  "interview_count": 4,
  "avg_fit_score": 78.5,
  "created_at": "2026-06-23T09:00:00Z"
}
```

---

## 4.2. Create job

```http
POST /companies/:company_id/jobs
```

Request:

```json
{
  "title": "Frontend Developer",
  "department": "Engineering",
  "level": "middle",
  "location": "Đà Nẵng",
  "employment_type": "full_time",
  "description": "Mô tả công việc...",
  "requirements": "Yêu cầu...",
  "benefits": "Phúc lợi...",
  "status": "draft"
}
```

Permission: recruiter/company member with job:create.

---

## 4.3. Get job detail

```http
GET /companies/:company_id/jobs/:job_id
```

Response:

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "title": "Frontend Developer",
    "description": "...",
    "requirements": "...",
    "status": "open",
    "ai_summary": "Vị trí cần ứng viên có kinh nghiệm Vue/React...",
    "ai_analysis_json": {},
    "stats": {
      "candidate_count": 12,
      "interview_count": 4,
      "report_ready_count": 3
    }
  }
}
```

---

## 4.4. Update job

```http
PUT /companies/:company_id/jobs/:job_id
```

Permission: company member with job:update.

---

## 4.5. Delete job

```http
DELETE /companies/:company_id/jobs/:job_id
```

Permission: company owner/admin or job owner.

Recommended: soft delete.

---

## 4.6. Analyze JD

```http
POST /companies/:company_id/jobs/:job_id/analyze
```

Request:

```json
{
  "force_refresh": false
}
```

Response:

```json
{
  "success": true,
  "data": {
    "summary": "Vị trí Frontend Developer tập trung vào Vue, UI performance và API integration.",
    "required_skills": ["Vue 3", "TypeScript", "REST API"],
    "nice_to_have_skills": ["Testing", "Design System"],
    "suggested_rubric": [
      {
        "name": "Technical Knowledge",
        "weight": 30
      }
    ],
    "suggested_questions": []
  }
}
```

Permission: company member with job:update.

---

## 5. Candidate API

## 5.1. List candidates

```http
GET /companies/:company_id/candidates?job_id=uuid&status=screening&keyword=nguyen&page=1&page_size=20
```

Response item:

```json
{
  "id": "uuid",
  "full_name": "Trần Văn B",
  "email": "candidate@example.com",
  "phone": "0900000000",
  "status": "screening",
  "source": "manual",
  "latest_job": {
    "id": "uuid",
    "title": "Frontend Developer"
  },
  "fit_score": 82.5,
  "latest_interview_at": null
}
```

Permission: company member.

---

## 5.2. Create candidate

```http
POST /companies/:company_id/candidates
```

Request:

```json
{
  "full_name": "Trần Văn B",
  "email": "candidate@example.com",
  "phone": "0900000000",
  "source": "manual",
  "job_id": "uuid"
}
```

---

## 5.3. Get candidate detail

```http
GET /companies/:company_id/candidates/:candidate_id
```

Response:

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "full_name": "Trần Văn B",
    "email": "candidate@example.com",
    "phone": "0900000000",
    "status": "screening",
    "cv_file": {
      "id": "uuid",
      "original_name": "cv.pdf",
      "download_url": "signed-url"
    },
    "ai_cv_summary": "Ứng viên có 3 năm kinh nghiệm frontend...",
    "skills": ["Vue", "React", "TypeScript"],
    "interviews": [],
    "reports": []
  }
}
```

---

## 5.4. Update candidate

```http
PUT /companies/:company_id/candidates/:candidate_id
```

---

## 5.5. Upload CV

```http
POST /companies/:company_id/candidates/:candidate_id/cv
Content-Type: multipart/form-data
```

Form data:

```text
file=<cv.pdf>
```

Response:

```json
{
  "success": true,
  "data": {
    "file_id": "uuid",
    "original_name": "cv.pdf",
    "size_bytes": 123456
  }
}
```

---

## 5.6. Parse CV

```http
POST /companies/:company_id/candidates/:candidate_id/parse-cv
```

Response:

```json
{
  "success": true,
  "data": {
    "summary": "Ứng viên có kinh nghiệm...",
    "skills": ["Vue", "Node.js"],
    "experience_years": 3,
    "education": [],
    "projects": []
  }
}
```

---

## 6. Interview API

## 6.1. List interviews

```http
GET /companies/:company_id/interviews?status=scheduled&date=2026-06-23&page=1&page_size=20
```

Response item:

```json
{
  "id": "uuid",
  "title": "Frontend Developer - Trần Văn B",
  "job": {
    "id": "uuid",
    "title": "Frontend Developer"
  },
  "candidate": {
    "id": "uuid",
    "full_name": "Trần Văn B"
  },
  "recruiter": {
    "id": "uuid",
    "full_name": "Nguyễn Văn A"
  },
  "scheduled_at": "2026-06-25T03:00:00Z",
  "status": "scheduled",
  "mode": "real"
}
```

---

## 6.2. Create interview

```http
POST /companies/:company_id/interviews
```

Request:

```json
{
  "job_id": "uuid",
  "candidate_id": "uuid",
  "recruiter_id": "uuid",
  "template_id": "uuid",
  "rubric_id": "uuid",
  "scheduled_at": "2026-06-25T03:00:00Z",
  "duration_minutes": 60,
  "mode": "real",
  "send_invite": true
}
```

Response:

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "status": "scheduled",
    "room_id": "uuid",
    "invite_url": "https://app.example.com/interview/join/token"
  }
}
```

---

## 6.3. Get interview detail

```http
GET /companies/:company_id/interviews/:interview_id
```

Response:

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "status": "scheduled",
    "mode": "real",
    "job": {},
    "candidate": {},
    "recruiter": {},
    "room": {},
    "rubric": {},
    "created_at": "2026-06-23T09:00:00Z"
  }
}
```

---

## 6.4. Start interview

```http
POST /companies/:company_id/interviews/:interview_id/start
```

Request:

```json
{
  "consent_recording": true,
  "consent_ai": true
}
```

Response:

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "status": "active",
    "started_at": "2026-06-25T03:00:00Z"
  }
}
```

Permission: recruiter assigned to interview or company admin.

---

## 6.5. End interview

```http
POST /companies/:company_id/interviews/:interview_id/end
```

Request:

```json
{
  "generate_report": true
}
```

Response:

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "status": "completed",
    "ended_at": "2026-06-25T04:00:00Z",
    "report_status": "generating"
  }
}
```

---

## 6.6. Cancel interview

```http
POST /companies/:company_id/interviews/:interview_id/cancel
```

Request:

```json
{
  "reason": "Ứng viên đổi lịch"
}
```

---

## 6.7. Join interview by invite token

```http
GET /interviews/join/:invite_token
```

> **One-Time Token (OTT) — Cơ chế xác thực link mời:**
> - `invite_token` trong URL là token dùng một lần, được hash trong DB (`invite_token_hash`).
> - Token có thời hạn (`invite_expires_at`). Sau khi hết hạn trả về `403 INVITE_EXPIRED`.
> - Khi gọi endpoint này thành công, server trả về `room_access_token` — một JWT phạm vi hẹp:
>   - Chỉ có quyền vào đúng `room_id` đó.
>   - TTL: 4 giờ kể từ lúc phát sinh.
>   - Không dùng được cho bất kỳ API nào khác.
> - `invite_token` gốc **không bị invalidate** sau lần đầu dùng (để candidate có thể reload trang). Tuy nhiên một `invite_token` chỉ được cấp cho một `candidate_id` — nếu dùng sai người trả về `403 INVITE_UNAUTHORIZED`.

Response:

```json
{
  "success": true,
  "data": {
    "interview_id": "uuid",
    "room_id": "uuid",
    "candidate_name": "Trần Văn B",
    "job_title": "Frontend Developer",
    "scheduled_at": "2026-06-25T03:00:00Z",
    "requires_consent_ai": true,
    "requires_consent_recording": false,
    "room_access_token": "jwt-scoped-to-room",
    "room_access_token_expires_at": "2026-06-25T08:00:00Z"
  }
}
```

Lỗi:

| HTTP | Error code | Khi nào |
|---|---|---|
| 404 | `INVITE_NOT_FOUND` | Token không tồn tại |
| 403 | `INVITE_EXPIRED` | Token hết hạn |
| 403 | `INTERVIEW_CANCELLED` | Buổi phỏng vấn đã hủy |

Permission: valid invite token (không cần Bearer auth — đây là entry point cho candidate chưa đăng nhập).

---

## 7. Room API

REST API cho room metadata. Realtime event xem `REALTIME_EVENTS.md`.

## 7.1. Get room detail

```http
GET /companies/:company_id/interviews/:interview_id/room
```

Response:

```json
{
  "success": true,
  "data": {
    "room_id": "uuid",
    "room_code": "room_abc",
    "status": "waiting",
    "participants": [],
    "realtime_url": "wss://api.example.com/ws/interview-room",
    "ice_servers": []
  }
}
```

---

## 7.2. Create room access token

```http
POST /companies/:company_id/interviews/:interview_id/room/token
```

> Dùng cho **recruiter** xin token vào phòng. Candidate dùng endpoint `GET /interviews/join/:invite_token` (mục 6.7).
> Token trả về là LiveKit Access Token — dùng để khởi tạo kết nối LiveKit SDK ở frontend.

Request:

```json
{
  "participant_type": "recruiter"
}
```

Response:

```json
{
  "success": true,
  "data": {
    "room_access_token": "livekit-jwt-token",
    "livekit_url": "wss://livekit.example.com",
    "expires_at": "2026-06-25T04:00:00Z"
  }
}
```

---

## 8. Transcript API

## 8.1. List transcripts

```http
GET /companies/:company_id/interviews/:interview_id/transcripts
```

Response item:

```json
{
  "id": "uuid",
  "speaker_type": "candidate",
  "speaker_name": "Trần Văn B",
  "content": "Tôi có 3 năm kinh nghiệm Vue...",
  "start_time_ms": 120000,
  "end_time_ms": 135000,
  "confidence": 0.91,
  "source": "audio",
  "is_final": true,
  "created_at": "2026-06-25T03:10:00Z"
}
```

Permission:

- Recruiter/company member: có.
- Candidate: chỉ nếu policy cho phép.

---

## 8.2. Create manual transcript/note-like transcript

```http
POST /companies/:company_id/interviews/:interview_id/transcripts
```

Request:

```json
{
  "speaker_type": "candidate",
  "speaker_name": "Trần Văn B",
  "content": "Nội dung transcript chỉnh tay",
  "source": "manual"
}
```

---

## 8.3. Edit transcript

```http
PUT /companies/:company_id/interviews/:interview_id/transcripts/:transcript_id
```

Request:

```json
{
  "edited_content": "Nội dung đã chỉnh sửa"
}
```

Rule:

- Không sửa mất bản gốc `content`.
- Lưu `edited_content`, `edited_by`, `edited_at`.
- Ghi audit log.

---

## 9. AI API

> **Lưu ý kiến trúc:** Các endpoint trong mục 9 là **REST API của Backend Golang** — frontend gọi vào đây.
> Golang Backend sau đó gọi tiếp đến **Python AI Orchestrator** qua **internal REST API** (không expose ra ngoài).
> Frontend không bao giờ gọi trực tiếp Python AI service.
>
> Với các task nặng bất đồng bộ (generate report, xử lý audio batch), Golang đẩy job vào **Redis Queue** thay vì gọi AI synchronously. Response trả về ngay với `status: "processing"` và client dùng WebSocket event `report:ready` để biết khi nào xong.

## 9.1. Generate questions

```http
POST /companies/:company_id/jobs/:job_id/ai/generate-questions
```

Request:

```json
{
  "question_types": ["technical", "behavioral", "experience"],
  "level": "middle",
  "count": 10
}
```

Response:

```json
{
  "success": true,
  "data": {
    "questions": [
      {
        "question_text": "Bạn hãy mô tả một lần bạn tối ưu performance frontend?",
        "question_type": "experience",
        "target_skill": "Frontend Performance",
        "difficulty": "middle",
        "expected_signals": ["biết đo metric", "có ví dụ thực tế"]
      }
    ]
  }
}
```

---

## 9.2. Suggest follow-up question

```http
POST /companies/:company_id/interviews/:interview_id/ai/suggest-follow-up
```

Request:

```json
{
  "last_transcript_id": "uuid",
  "focus": "technical_depth"
}
```

Response:

```json
{
  "success": true,
  "data": {
    "suggested_question": "Bạn có thể nói rõ bạn đo performance bằng chỉ số nào không?",
    "reason": "Ứng viên có nhắc tối ưu performance nhưng chưa nêu metric cụ thể.",
    "target_skill": "Performance Optimization",
    "priority": "high",
    "confidence": 0.84
  }
}
```

---

## 9.3. Score answer

```http
POST /companies/:company_id/interviews/:interview_id/ai/score-answer
```

Request:

```json
{
  "transcript_ids": ["uuid1", "uuid2"],
  "criterion_ids": ["uuid3"]
}
```

Response:

```json
{
  "success": true,
  "data": {
    "scores": [
      {
        "criterion_name": "Technical Knowledge",
        "score": 4,
        "max_score": 5,
        "evidence": "Ứng viên mô tả được cách dùng index và cache.",
        "ai_comment": "Câu trả lời có ví dụ thực tế, nhưng thiếu số liệu đo lường.",
        "confidence": 0.78,
        "status": "scored"
      }
    ]
  }
}
```

---

## 9.4. Generate report

```http
POST /companies/:company_id/interviews/:interview_id/ai/generate-report
```

Request:

```json
{
  "force_refresh": false
}
```

Response:

```json
{
  "success": true,
  "data": {
    "report_id": "uuid",
    "status": "ready",
    "final_score": 78.5,
    "recommendation": "consider"
  }
}
```

Note:

- Nếu report tạo lâu, API có thể trả `status: generating` và frontend poll report endpoint.

---

## 10. Report API

## 10.1. Get interview report

```http
GET /companies/:company_id/interviews/:interview_id/report
```

Response:

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "interview_id": "uuid",
    "summary": "Ứng viên có nền tảng frontend khá tốt...",
    "final_score": 78.5,
    "recommendation": "consider",
    "strengths": ["Giao tiếp rõ", "Có kinh nghiệm dự án thực tế"],
    "weaknesses": ["Thiếu số liệu đo performance"],
    "risks": ["Chưa chứng minh kinh nghiệm scale lớn"],
    "scores": [
      {
        "criterion_name": "Technical Knowledge",
        "score": 4,
        "weight": 30,
        "evidence": "..."
      }
    ],
    "recruiter_decision": null,
    "generated_at": "2026-06-25T04:05:00Z"
  }
}
```

---

## 10.2. Update recruiter decision

```http
PUT /companies/:company_id/interviews/:interview_id/report/decision
```

Request:

```json
{
  "decision": "next_round",
  "comment": "Ứng viên phù hợp, mời vào vòng technical sâu hơn."
}
```

Rule:

- Recruiter decision là quyết định cuối.
- AI recommendation không tự overwrite recruiter decision.

---

## 11. Mock Interview API

## 11.1. Create mock interview

```http
POST /mock-interviews
```

Request:

```json
{
  "target_role": "Frontend Developer",
  "target_level": "junior",
  "cv_file_id": "uuid"
}
```

Permission: candidate.

Response:

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "status": "draft",
    "target_role": "Frontend Developer",
    "target_level": "junior"
  }
}
```

---

## 11.2. Start mock interview

```http
POST /mock-interviews/:mock_interview_id/start
```

Response:

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "status": "active",
    "first_question": {
      "id": "q_1",
      "content": "Bạn hãy giới thiệu ngắn gọn về bản thân."
    }
  }
}
```

---

## 11.3. Submit mock answer

```http
POST /mock-interviews/:mock_interview_id/answer
```

Request:

```json
{
  "question_id": "q_1",
  "answer_text": "Em là..."
}
```

Response:

```json
{
  "success": true,
  "data": {
    "feedback": {
      "score": 3.5,
      "strengths": ["Trả lời rõ ràng"],
      "improvements": ["Nên thêm ví dụ cụ thể"],
      "sample_better_answer": "Bạn có thể trả lời theo cấu trúc..."
    },
    "next_question": {
      "id": "q_2",
      "content": "Bạn đã từng làm dự án frontend nào đáng chú ý?"
    }
  }
}
```

---

## 11.4. End mock interview

```http
POST /mock-interviews/:mock_interview_id/end
```

Response:

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "status": "completed",
    "final_score": 72.0,
    "feedback": {
      "summary": "Bạn có nền tảng tốt nhưng cần trả lời có cấu trúc hơn.",
      "strengths": [],
      "weaknesses": [],
      "practice_plan": []
    }
  }
}
```

---

## 11.5. List my mock interviews

```http
GET /mock-interviews/my?page=1&page_size=20
```

Permission: candidate.

---

## 12. File API

## 12.1. Upload file

```http
POST /files
Content-Type: multipart/form-data
```

Form data:

```text
file=<file>
file_type=cv
company_id=uuid optional
```

Response:

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "original_name": "cv.pdf",
    "mime_type": "application/pdf",
    "size_bytes": 123456,
    "file_type": "cv"
  }
}
```

---

## 12.2. Get signed download URL

```http
GET /files/:file_id/download-url
```

Response:

```json
{
  "success": true,
  "data": {
    "download_url": "signed-url",
    "expires_at": "2026-06-23T10:00:00Z"
  }
}
```

---

## 13. Audit API

## 13.1. List audit logs

```http
GET /companies/:company_id/audit-logs?page=1&page_size=20&resource_type=interview
```

Permission: company owner/admin.

---

## 14. Permission summary

| API group | Admin | Recruiter | Candidate |
|---|---:|---:|---:|
| Auth | Có | Có | Có |
| Company | Có | Có nếu member | Không |
| Job | Có | Có nếu member | Không |
| Candidate | Có | Có nếu member | Không |
| Interview scheduling | Có | Có nếu member | Không |
| Interview join | Có | Có nếu assigned | Có nếu invite hợp lệ |
| Room metadata | Có | Có nếu assigned/member | Có nếu invite hợp lệ |
| AI scoring | Có | Có | Không |
| Report tuyển dụng | Có | Có | Không |
| Mock interview | Không bắt buộc | Không bắt buộc | Có |

---

## 15. Integration contract cho 3 người

## 15.1. Hùng

Hùng không tự tạo endpoint mới cho realtime nếu chưa cập nhật file này và `REALTIME_EVENTS.md`.

Hùng tập trung:

- Room token.
- WebSocket events.
- Transcript event.
- AI suggestion event.
- Room state.

## 15.2. Khôi

Khôi đảm bảo backend trả response đúng schema để Lai tích hợp.

Khôi tập trung:

- Auth.
- Company.
- Job.
- Candidate.
- Interview.
- AI orchestration.
- Report.
- DB.

## 15.3. Lai

Lai không hardcode field ngoài API spec.

Lai tập trung:

- UI gọi API.
- Dashboard.
- Job detail.
- Candidate detail.
- Interview room UI.
- Report UI.
- Candidate mock portal.

---

## 16. DoD cho API

Một API được xem là hoàn thành khi:

- Đúng endpoint trong spec.
- Đúng method.
- Đúng request/response schema.
- Có validation.
- Có permission check.
- Có tenant check bằng `company_id` nếu cần.
- Có error response chuẩn.
- Có test case cơ bản.
- Có audit log nếu là hành động nhạy cảm.
- Frontend có thể gọi thử thành công.
