# AI_INTERVIEW_PROJECT_RULES.md
# Bộ Quy Tắc Dự Án: Phỏng vấn cùng AI Real-time

## 1. Mục đích

File này định nghĩa các quy tắc bắt buộc khi phát triển dự án **phỏng vấn cùng AI real-time**.

Tất cả thành viên và AI coding assistant phải tuân thủ để tránh:

- Chồng chéo module.
- Phá luồng realtime.
- Sai phân quyền dữ liệu.
- AI đánh giá thiếu kiểm soát.
- Database/API không nhất quán.
- UI thiếu trạng thái lỗi/loading.
- Code xong nhưng không test được.

---

## 2. Quy tắc phân công

| Người | Phạm vi chính | Không tự ý làm |
|---|---|---|
| Hùng | Interview Room, Realtime, WebSocket/WebRTC, Transcript live, Room State | Không tự ý đổi AI scoring/report schema nếu chưa thống nhất với Khôi |
| Khôi | Backend core, Database, AI Engine, Scoring, Report, API chuẩn | Không tự ý redesign UI/phân luồng màn hình nếu chưa thống nhất với Lai |
| Lai | UI/UX, Recruiter Portal, Candidate Portal, Job/Candidate/Mock CRUD | Không tự ý sửa realtime gateway hoặc AI engine lõi |

Nếu cần sửa module của người khác:

1. Ghi rõ lý do.
2. Chỉ sửa phần tối thiểu.
3. Báo lại file đã sửa.
4. Không thay đổi contract nếu chưa được đồng ý.

---

## 3. Quy tắc branch và commit

### 3.1. Đặt tên branch

```text
feature/<owner>/<module-name>
fix/<owner>/<bug-name>
refactor/<owner>/<scope>
```

Ví dụ:

```text
feature/hung/interview-room-realtime
feature/khoi/ai-scoring-engine
feature/lai/candidate-portal
fix/hung/reconnect-room-state
fix/khoi/report-generation-error
fix/lai/job-list-empty-state
```

### 3.2. Commit message

Format:

```text
<type>(<module>): <short description>
```

Type:

| Type | Ý nghĩa |
|---|---|
| feat | Thêm tính năng |
| fix | Sửa lỗi |
| refactor | Tái cấu trúc code |
| docs | Tài liệu |
| test | Test |
| chore | Việc phụ trợ |

Ví dụ:

```text
feat(interview-room): add room join and presence events
feat(ai-scoring): add rubric based scoring API
fix(candidate-portal): handle empty mock interview history
```

---

## 4. Quy tắc code chung

- Không hardcode secret, API key, token.
- Không để console log rác trong production code.
- Không copy code lặp lại quá nhiều.
- Không tạo file quá lớn nếu có thể tách component/service.
- Không sửa tên field/API đã dùng nơi khác nếu chưa kiểm tra toàn bộ reference.
- Không xóa code người khác nếu chưa hiểu rõ.
- Không bỏ qua xử lý lỗi.
- Không để UI treo loading vô hạn.
- Không để Candidate xem dữ liệu nội bộ của Recruiter.
- Không để AI tự động quyết định pass/fail cuối cùng.

---

## 5. Quy tắc API

### 5.1. Response chuẩn

API nên thống nhất response:

```json
{
  "success": true,
  "data": {},
  "message": "OK",
  "error": null
}
```

Khi lỗi:

```json
{
  "success": false,
  "data": null,
  "message": "Validation failed",
  "error": {
    "code": "VALIDATION_ERROR",
    "details": []
  }
}
```

### 5.2. Status code

| Code | Dùng khi |
|---:|---|
| 200 | Thành công |
| 201 | Tạo mới thành công |
| 400 | Request sai |
| 401 | Chưa đăng nhập |
| 403 | Không có quyền |
| 404 | Không tìm thấy |
| 409 | Xung đột dữ liệu |
| 422 | Validation lỗi |
| 500 | Lỗi server |

### 5.3. API naming

Dùng REST rõ ràng:

```text
GET    /jobs
POST   /jobs
GET    /jobs/:id
PUT    /jobs/:id
DELETE /jobs/:id

GET    /candidates
POST   /candidates
GET    /candidates/:id
PUT    /candidates/:id

POST   /interviews/:id/start
POST   /interviews/:id/end
GET    /interviews/:id/report
```

---

## 6. Quy tắc database

- Tên bảng: snake_case, số nhiều.
- Tên cột: snake_case.
- ID chính: UUID.
- Entity chính phải có `created_at`, `updated_at`.
- Không xóa cứng dữ liệu quan trọng nếu chưa cần thiết.
- Dữ liệu nhạy cảm phải hạn chế quyền truy cập.
- Migration phải có thứ tự rõ ràng.
- Không sửa migration cũ đã chạy trên môi trường chung, hãy tạo migration mới.

### 6.1. Bảng lõi nên có

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
- audit_logs

---

## 7. Quy tắc phân quyền

### 7.1. Recruiter được phép

- Quản lý job trong company của mình.
- Quản lý candidate trong company của mình.
- Tạo lịch phỏng vấn.
- Vào room với vai trò recruiter.
- Xem AI suggestion.
- Xem AI score.
- Xem report tuyển dụng.
- Ghi note nội bộ.
- Ra quyết định cuối.

### 7.2. Candidate được phép

- Xem thông tin buổi phỏng vấn của mình.
- Vào room bằng link hợp lệ.
- Cập nhật profile cá nhân nếu được cho phép.
- Upload CV nếu được cho phép.
- Tham gia mock interview.
- Xem feedback mock interview của chính mình.

### 7.3. Candidate không được phép

- Xem note nội bộ của recruiter.
- Xem AI score tuyển dụng nếu company không bật.
- Xem report tuyển dụng nội bộ.
- Xem thông tin candidate khác.
- Xem job/candidate của company khác.

---

## 8. Quy tắc AI

### 8.1. AI không được làm

- Không tự quyết định tuyển/loại ứng viên.
- Không đánh giá dựa trên tuổi, giới tính, tôn giáo, chủng tộc, ngoại hình.
- Không suy diễn khi không đủ dữ liệu.
- Không bịa bằng chứng.
- Không dùng dữ liệu ứng viên để training nếu chưa có consent.
- Không hiển thị output nhạy cảm cho sai actor.

### 8.2. AI phải làm

- Có evidence khi chấm điểm.
- Có confidence nếu có thể.
- Ghi rõ `insufficient_data` khi thiếu dữ liệu.
- Tách rõ recommendation của AI và final decision của recruiter.
- Output nên có cấu trúc JSON để dễ lưu và hiển thị.

### 8.3. AI scoring format

```json
{
  "criterion": "Communication",
  "score": 4,
  "max_score": 5,
  "weight": 0.15,
  "evidence": "Ứng viên trả lời rõ ràng, có ví dụ cụ thể.",
  "comment": "Khả năng giao tiếp tốt.",
  "confidence": 0.78,
  "insufficient_data": false
}
```

---

## 9. Quy tắc realtime

Realtime module phải đảm bảo:

- User join/leave chính xác.
- Có presence state.
- Có room status.
- Có reconnect.
- Không mất toàn bộ state khi reload.
- Không crash khi AI service lỗi.
- Không gửi event nhạy cảm cho sai role.
- Event phải có format thống nhất.

### 9.1. Event format

```json
{
  "event": "transcript:update",
  "room_id": "uuid",
  "sender_id": "uuid",
  "sender_role": "candidate",
  "payload": {},
  "timestamp": "2026-06-23T10:00:00Z"
}
```

### 9.2. Event quan trọng

- room:join
- room:leave
- room:reconnect
- interview:start
- interview:end
- chat:send
- transcript:update
- ai:suggestion
- ai:score_update
- ai:error
- report:ready

---

## 10. Quy tắc UI/UX

- Tất cả màn hình phải có loading state.
- Tất cả danh sách phải có empty state.
- Tất cả API lỗi phải có error state.
- Form phải validate trước khi submit.
- Nút nguy hiểm phải có confirm.
- Không dùng quá nhiều màu gây rối.
- Ưu tiên layout rõ ràng, dễ đọc.
- Candidate UI phải đơn giản, tránh tạo áp lực.
- Recruiter UI phải ưu tiên thông tin ra quyết định.

---

## 11. Quy tắc bảo mật

- Bắt buộc HTTPS ở production.
- JWT/Session phải có thời hạn.
- Refresh token phải an toàn.
- File CV/recording dùng signed URL hoặc quyền truy cập kiểm soát.
- Không expose storage key trực tiếp nếu không cần.
- Audit log cho hành động quan trọng.
- Rate limit cho API nhạy cảm.
- Validate input cả frontend và backend.
- Không tin dữ liệu từ client.

---

## 12. Quy tắc audit log

Phải ghi audit log khi:

- Tạo/sửa/xóa job.
- Thêm/sửa/xóa candidate.
- Upload/xóa CV.
- Tạo lịch phỏng vấn.
- Bắt đầu/kết thúc phỏng vấn.
- Xem/tải report.
- Xem/tải transcript.
- Thay đổi final decision.
- Thay đổi quyền thành viên.

Audit log nên có:

```json
{
  "actor_id": "uuid",
  "actor_role": "recruiter",
  "action": "interview.end",
  "resource_type": "interview",
  "resource_id": "uuid",
  "metadata": {},
  "created_at": "timestamp"
}
```

---

## 13. Quy tắc test

Mỗi task phải có ít nhất checklist test thủ công.

### 13.1. Checklist chung

- Đăng nhập/chưa đăng nhập.
- Đúng quyền/sai quyền.
- Dữ liệu hợp lệ.
- Dữ liệu thiếu.
- Dữ liệu không tồn tại.
- API lỗi.
- Network lỗi.
- Loading state.
- Empty state.
- Reload trang.
- Responsive cơ bản.

### 13.2. Checklist realtime

- Join room thành công.
- Leave room thành công.
- Reconnect sau mất mạng.
- Start/end interview đúng role.
- Chat realtime hoạt động.
- Transcript update không bị nhân đôi.
- AI lỗi nhưng room vẫn hoạt động.

### 13.3. Checklist AI

- AI nhận đúng input.
- AI output đúng schema.
- AI không bịa evidence.
- AI trả insufficient_data khi thiếu dữ liệu.
- AI report sinh được sau phỏng vấn.
- Lỗi AI được hiển thị/thử lại.

---

## 14. Definition of Done chung

Một task được xem là xong khi:

- Đúng yêu cầu nghiệp vụ.
- Không phá module khác.
- Code chạy được.
- Không lỗi build.
- Không lỗi lint nghiêm trọng.
- Có validate dữ liệu.
- Có xử lý quyền.
- Có loading/error/empty state nếu là UI.
- Có response lỗi rõ ràng nếu là API.
- Có audit log nếu là hành động quan trọng.
- Có checklist test.
- Có ghi lại file đã thay đổi.

---

## 15. Quy tắc review

Khi review task, cần kiểm tra:

1. Có đúng module owner không?
2. Có đúng requirement không?
3. Có thiếu edge case không?
4. Có phá API contract không?
5. Có lộ dữ liệu nhạy cảm không?
6. Có xử lý lỗi không?
7. Có test được không?
8. Có dễ bảo trì không?
9. Có phù hợp roadmap không?

---

## 16. Quy tắc ưu tiên Sprint

Thứ tự ưu tiên:

1. Auth + phân quyền cơ bản.
2. Job + Candidate.
3. Interview scheduling.
4. Interview room realtime cơ bản.
5. Transcript + chat.
6. AI generate question.
7. AI scoring.
8. AI report.
9. Mock interview.
10. Dashboard + analytics.

Không làm tính năng nâng cao khi nền tảng chưa ổn.
