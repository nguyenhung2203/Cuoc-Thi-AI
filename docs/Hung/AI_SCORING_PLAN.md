# Kế hoạch Hùng — AI chấm điểm & đánh giá phỏng vấn

> **Owner:** Hùng (AI giám khảo / scoring)
> **Phối hợp:** Khôi (AI service, schema, API, database), Lai (UI report/profile)
> **Mục tiêu:** Chấm câu trả lời dựa trên năng lực và bằng chứng, không chỉ đúng/sai kiến thức; đánh giá đúng theo cấp bậc tuyển dụng.
> **Tài liệu tham chiếu:** `docs/AI_PROMPT_AND_SCORING.md`, `docs/Hung/MODULE_HUNG_REALTIME_AI.md`, `docs/Hung/task/H-S4-02_ai_score_update.md`, `docs/Khoi/MODULE_KHOI_CORE_AI_BACKEND.md`

---

## 1. Phạm vi và nguyên tắc bắt buộc

### 1.1. AI là trợ lý đánh giá

- AI chỉ đưa ra điểm, bằng chứng và khuyến nghị.
- Recruiter là người quyết định cuối cùng.
- Không tự động reject/pass ứng viên.
- Mọi điểm phải trích dẫn từ câu trả lời, transcript, JD, CV hoặc rubric.
- Nếu thiếu dữ liệu: trả `insufficient_evidence`, không đoán.

### 1.2. Không đánh giá yếu tố nhạy cảm

Không chấm theo giới tính, tuổi, vùng miền, ngoại hình, giọng địa phương, tôn giáo, dân tộc, tình trạng hôn nhân hoặc các yếu tố không liên quan trực tiếp đến công việc.

### 1.3. "Tính cách" phải được giới hạn đúng

Không suy đoán tính cách tâm lý. Chỉ đánh giá **hành vi thể hiện trong câu trả lời có liên quan công việc**, ví dụ:

- ownership/trách nhiệm;
- tinh thần hợp tác;
- khả năng nhận phản hồi;
- tư duy học hỏi;
- tính trung thực khi mô tả kinh nghiệm.

---

## 2. Chuẩn hóa cấp bậc và độ kỳ vọng

### 2.1. Enum cấp bậc

Dùng thống nhất: `intern`, `fresher`, `junior`, `middle`, `senior`, `lead`, `unknown`.

Nếu hệ thống hiện dùng `mid`, chuẩn hóa ở API boundary thành `middle` hoặc hỗ trợ alias `mid -> middle`, không để prompt nhận nhiều tên không thống nhất.

### 2.2. Kỳ vọng theo cấp bậc

| Cấp bậc | Kỳ vọng khi chấm |
|---|---|
| Intern/Fresher | Nắm khái niệm nền tảng, tư duy học hỏi, giải thích được lựa chọn cơ bản; không trừ điểm chỉ vì thiếu kinh nghiệm production |
| Junior | Làm được task có hướng dẫn, hiểu quy trình, debug cơ bản, đưa được ví dụ dự án/thực hành |
| Middle | Tự chủ giải quyết vấn đề, phân tích trade-off, thiết kế ở mức vừa, có kinh nghiệm vận hành/thực tế |
| Senior | Ra quyết định có hệ thống, đánh giá rủi ro, scale, observability, mentoring, ownership và tác động kinh doanh |
| Lead | Định hướng kỹ thuật/đội nhóm, chiến lược, phân bổ nguồn lực, quản trị rủi ro và kết quả liên phòng ban |

### 2.3. Quy tắc chống chấm sai level

- Luôn truyền `job.level` vào scoring context.
- Rubric phải ghi rõ `expected_signals` theo level.
- Không yêu cầu senior-level signal ở fresher/junior.
- Nếu level là `unknown`, AI phải trả cảnh báo và dùng rubric tổng quát thận trọng.
- Độ khó câu hỏi do Khôi sinh phải khớp level; Hùng không tự dùng câu hỏi senior để kết luận ứng viên junior yếu.

---

## 3. Rubric scoring phiên bản GK1/GK3/GK4

### 3.1. Tiêu chí và trọng số mặc định

| Tiêu chí | Trọng số | Nội dung |
|---|---:|---|
| Technical Knowledge | 25 | Kiến thức đúng và liên quan vị trí |
| Problem Solving | 20 | Phân tích, cách tiếp cận, kiểm chứng và xử lý edge case |
| Experience Relevance | 15 | Mức độ liên quan của ví dụ với JD và level |
| Communication | 15 | Rõ ràng, có cấu trúc, trả lời đúng trọng tâm |
| Collaboration & Work Behavior | 10 | Hợp tác, ownership, nhận feedback qua bằng chứng |
| Growth Mindset | 5 | Nhận diện thiếu sót và cách cải thiện |
| Role/Level Fit | 10 | Đáp ứng kỳ vọng của level, không so với level khác |

Tổng trọng số: 100. Cho phép rubric theo job override trọng số, nhưng phải kiểm tra tổng bằng 100.

### 3.2. Thang điểm từng tiêu chí

- `1`: Không có bằng chứng hoặc sai nghiêm trọng.
- `2`: Có hiểu biết sơ bộ nhưng thiếu phần cốt lõi.
- `3`: Đạt kỳ vọng tối thiểu của level.
- `4`: Vượt kỳ vọng, có ví dụ và lập luận tốt.
- `5`: Rất tốt, có chiều sâu, trade-off, tác động và bài học rõ.

Điểm `3` phải được diễn giải theo level. Ví dụ, câu trả lời đạt 3 ở fresher không cần có scale production như senior.

### 3.3. Điểm tổng

```text
weighted_score = sum(score / 5 * weight)
final_score = round(weighted_score, 0..100)
```

Nếu thiếu bằng chứng ở một tiêu chí, không tự chấm 0; đánh dấu `insufficient_evidence` và hạ confidence/đề xuất câu hỏi bổ sung.

---

## 4. Output contract đề xuất

```json
{
  "status": "scored|insufficient_evidence|degraded",
  "level": "junior",
  "question_id": "uuid",
  "overall_score": 76,
  "confidence": 0.82,
  "criteria": [
    {
      "key": "communication",
      "score": 4,
      "max_score": 5,
      "weight": 15,
      "weighted_score": 12,
      "evidence": ["...trích ý từ câu trả lời..."],
      "strengths": ["..."],
      "gaps": ["..."],
      "status": "scored"
    }
  ],
  "strengths": ["..."],
  "weaknesses": ["..."],
  "feedback": {
    "missing": ["Thiếu số liệu tác động và cách kiểm chứng."],
    "emphasize": ["Nêu rõ vai trò cá nhân, trade-off và kết quả."],
    "example_better_answer": "..."
  },
  "follow_up_questions": ["..."],
  "recommendation": "consider"
}
```

### Acceptance criteria

- Có tối thiểu technical, communication, work behavior/growth, strengths, weaknesses và feedback.
- Có `evidence` cho từng tiêu chí đã chấm.
- Có `level` trong output.
- Có confidence và status.
- Không trả điểm nếu thiếu context nghiêm trọng.
- JSON parse được và validate đúng schema.

---

## 5. Checklist triển khai

### H-SC-01 — Audit hiện trạng scoring

- [ ] Đọc prompt/model/parser hiện tại của AI service.
- [ ] Đối chiếu schema score realtime của Hùng với schema Python/Go.
- [ ] Liệt kê field đang mất khi đi qua WebSocket/webhook/database.
- [ ] Chốt enum level và criterion key với Khôi.

### H-SC-02 — Thiết kế rubric theo level

- [ ] Tạo bảng expected signals cho từng level.
- [ ] Tách tiêu chí kiến thức khỏi giao tiếp/hành vi.
- [ ] Định nghĩa rule không phạt thiếu production experience ở fresher.
- [ ] Thêm test case cùng một câu trả lời nhưng chấm ở fresher/junior/senior phải khác kỳ vọng.

### H-SC-03 — Nâng cấp prompt giám khảo

- [ ] Bắt AI đóng vai evaluator hỗ trợ recruiter, không phải người quyết định.
- [ ] Bắt buộc trích evidence.
- [ ] Yêu cầu đánh giá clarity, structure, concision và collaboration behavior.
- [ ] Cấm suy đoán personality/giọng nói/vùng miền.
- [ ] Bắt buộc trả strengths, weaknesses, missing, emphasize, follow-up.

### H-SC-04 — Tích hợp realtime

- [ ] Mở rộng payload `ai:score_update` nhưng giữ backward compatibility.
- [ ] Chỉ recruiter nhận chi tiết nhạy cảm của scoring trong live interview.
- [ ] Candidate chỉ nhận feedback được phép ở mock interview hoặc sau khi recruiter publish.
- [ ] Handle AI degraded/error mà không làm gián đoạn phòng phỏng vấn.

### H-SC-05 — Lưu report và tổng điểm

- [ ] Lưu snapshot criterion name, weight, level, score, evidence.
- [ ] Tính overall score ở backend từ criterion scores, không tin tổng điểm client.
- [ ] Lưu strengths/weaknesses/feedback trong report JSONB.
- [ ] Expose API report/profile cho Lai dùng.
- [ ] Có version prompt/rubric để audit kết quả cũ.

### H-SC-06 — Kiểm thử

- [ ] Unit test schema validation.
- [ ] Test thiếu transcript/evidence.
- [ ] Test level mapping fresher/junior/middle/senior.
- [ ] Test câu trả lời đúng kiến thức nhưng giao tiếp kém.
- [ ] Test giao tiếp tốt nhưng kiến thức thiếu.
- [ ] Test không xuất hiện sensitive scoring.
- [ ] Test retry/degraded AI và không duplicate score.

---

## 6. Nâng cấp sau MVP

- [ ] Feedback cụ thể theo STAR/CARE: Situation, Task, Action, Result.
- [ ] So sánh score theo nhiều câu hỏi nhưng tránh double-count.
- [ ] Dashboard điểm tổng và confidence.
- [ ] Recruiter chỉnh rubric/weight theo job.
- [ ] Human override có audit log.
