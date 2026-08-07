# Kế hoạch Khôi — Model sinh câu hỏi phỏng vấn

> **Owner:** Khôi (AI Orchestrator / Question Generation)
> **Phối hợp:** Hùng (realtime follow-up/scoring), Lai (question UI), Âu (job data)
> **Mục tiêu:** Sinh câu hỏi đúng vị trí, đúng cấp bậc và đúng mode; hỗ trợ cả luyện tập AI và recruiter phỏng vấn thật.
> **Tài liệu tham chiếu:** `docs/AI_PROMPT_AND_SCORING.md`, `docs/Khoi/MODULE_KHOI_CORE_AI_BACKEND.md`, `docs/Khoi/CHECKLIST_KHOI.md`, `ai-service/app/prompts/question_generation_v1.py`

---

## 1. Chốt vai trò của AI

Hệ thống có hai mode, phải truyền rõ trong input và không trộn prompt:

### `real` — AI hỗ trợ recruiter thật

- Recruiter là người hỏi và quyết định.
- AI sinh question set trước buổi phỏng vấn theo JD/rubric.
- Trong buổi phỏng vấn, AI có thể đề xuất follow-up dựa trên transcript gần nhất.
- AI suggestion chỉ hiển thị cho recruiter; không tự gửi câu hỏi cho candidate.

### `mock` — AI là interviewer luyện tập

- AI đóng vai nhà tuyển dụng mô phỏng.
- AI hỏi từng câu, chờ candidate trả lời, rồi sinh follow-up hoặc feedback.
- Phải ghi rõ đây là mô phỏng, không phải quyết định tuyển dụng thật.

Input bắt buộc có `interview.mode`, `job.level`, `job.title`, `requirements`, `rubric`, `language`.

---

## 2. Chuẩn hóa level và ma trận độ khó

### 2.1. Enum

`intern | fresher | junior | middle | senior | lead | unknown`.

Hỗ trợ alias API `mid -> middle`, không gửi alias vào prompt. Nếu không xác định được level, trả `unknown` và sinh bộ câu hỏi nền tảng có cảnh báo.

### 2.2. Ma trận câu hỏi

| Level | Tỷ lệ focus đề xuất | Câu hỏi đạt yêu cầu |
|---|---|---|
| Intern/Fresher | 50% fundamentals, 30% project/learning, 20% behavior | Khái niệm nền, ví dụ học tập/dự án, cách học và nhận feedback |
| Junior | 35% fundamentals, 35% implementation/debugging, 30% behavior | Task thực tế có hướng dẫn, debug, test, giao tiếp trong team |
| Middle | 25% fundamentals, 40% design/trade-off, 35% delivery | Tự chủ, lựa chọn kiến trúc, edge cases, vận hành và tác động |
| Senior | 15% fundamentals, 45% architecture/scale, 40% leadership | Scale, reliability, risk, mentoring, ownership, business impact |
| Lead | 10% fundamentals, 40% strategy/architecture, 50% leadership | Chiến lược, team, ưu tiên, stakeholder, đo lường kết quả |

### 2.3. Guardrails

- Không hỏi kiến thức/scale cao hơn nhiều so với level rồi kết luận ứng viên yếu.
- Câu hỏi phải liên quan skill trong JD hoặc rubric.
- Không bịa công nghệ không có trong JD nếu không ghi là câu hỏi mở rộng.
- Mỗi câu ghi `difficulty`, `level`, `skill_tags`, `expected_signals`.
- Tránh câu hỏi yêu cầu thông tin nhạy cảm hoặc không liên quan công việc.

---

## 3. Output contract

```json
{
  "mode": "real",
  "level": "junior",
  "questions": [
    {
      "id": "q-1",
      "question_text": "...",
      "category": "technical|behavioral|problem_solving|communication",
      "difficulty": "basic|intermediate|advanced",
      "skill_tags": ["..."],
      "expected_signals": ["..."],
      "follow_up_prompts": ["..."],
      "timebox_minutes": 5,
      "evidence_required": true
    }
  ],
  "coverage": ["skill-a", "skill-b"],
  "warnings": [],
  "confidence": 0.85
}
```

### Acceptance criteria

- Có level/mode trong output.
- Không có câu hỏi ngoài level nếu không ghi rõ lý do.
- Bao phủ skill bắt buộc của JD.
- Có expected signals để Hùng chấm thống nhất.
- Có follow-up prompt có điều kiện, không chỉ danh sách câu hỏi cố định.
- JSON validate được; nếu thiếu JD/level trả warning thay vì bịa.

---

## 4. Checklist triển khai

### K-Q-01 — Audit và contract

- [ ] Đọc route/service/model hiện tại trong `ai-service`.
- [ ] Đối chiếu response Python với `question_bank` và API Go.
- [ ] Chốt enum level/difficulty/category với Hùng và Lai.
- [ ] Bổ sung `mode`, `level`, `skill_tags`, `expected_signals` nếu đang thiếu.

### K-Q-02 — Prompt sinh question set

- [ ] Viết prompt version mới, không sửa ngầm version cũ.
- [ ] Tách system instruction cho `real` và `mock`.
- [ ] Truyền level + ma trận tỷ lệ focus vào prompt.
- [ ] Bắt buộc output JSON schema.
- [ ] Thêm validation và fallback nếu model trả câu hỏi sai level.

### K-Q-03 — Adaptive follow-up

- [ ] Nhận transcript/câu trả lời gần nhất.
- [ ] Phân loại thiếu evidence, mâu thuẫn, trả lời quá chung hoặc đã đủ evidence.
- [ ] Sinh tối đa 1–2 follow-up liên quan, có lý do.
- [ ] Rate-limit và timeout; AI lỗi không làm hỏng interview room.
- [ ] Hùng nhận event contract rõ ràng cho realtime.

### K-Q-04 — Question bank và report

- [ ] Lưu question snapshot vào `question_bank`/interview snapshot.
- [ ] Lưu prompt version/model/metadata để audit.
- [ ] Không thay đổi câu hỏi đã dùng trong interview khi prompt được nâng version.
- [ ] Cung cấp API cho Lai load question set và expected signals theo quyền.

### K-Q-05 — Kiểm thử

- [ ] Cùng JD nhưng đổi level phải tạo độ khó/expected signals khác.
- [ ] Test fresher không bị hỏi scale/leadership như senior.
- [ ] Test senior có trade-off/scale/risk.
- [ ] Test mode mock sinh câu hỏi trực tiếp cho candidate.
- [ ] Test mode real chỉ sinh suggestion cho recruiter.
- [ ] Test thiếu level/JD, malformed JSON, timeout, retry và fallback.
- [ ] Test không sinh câu hỏi nhạy cảm/ngoài job scope.

---

## 5. Nâng cấp sau MVP

- [ ] Sinh câu hỏi theo transcript realtime cho recruiter.
- [ ] Cho candidate luyện tập trước phỏng vấn thật.
- [ ] Điều chỉnh độ khó theo câu trả lời nhưng giữ giới hạn level.
- [ ] Cho recruiter duyệt/sửa question set trước khi dùng.
- [ ] Đo coverage skill và tránh lặp câu hỏi.
