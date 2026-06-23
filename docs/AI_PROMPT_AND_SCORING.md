# AI_PROMPT_AND_SCORING.md
# Prompt, Scoring và quy tắc AI — AI Interview Platform

## 0. Mục đích tài liệu

Tài liệu này định nghĩa cách AI hoạt động trong hệ thống:

- Phân tích JD.
- Phân tích CV.
- Sinh câu hỏi.
- Gợi ý follow-up realtime.
- Chấm điểm câu trả lời.
- Tạo report sau phỏng vấn.
- Feedback mock interview cho candidate.

AI IDE phải dùng file này khi code các phần liên quan AI để đảm bảo AI không bịa, không bias, có evidence và output đúng schema.

---

## 1. Nguyên tắc AI bắt buộc

## 1.1. AI không quyết định cuối cùng

AI chỉ đưa recommendation. Recruiter là người quyết định cuối.

AI không được tự động:

- Reject candidate.
- Pass candidate.
- Gửi email kết quả tuyển dụng.
- Thay đổi trạng thái candidate nếu không có xác nhận.

---

## 1.2. AI phải có evidence

Mọi đánh giá phải dựa trên ít nhất một trong các nguồn:

- JD.
- CV.
- Transcript.
- Câu trả lời candidate.
- Rubric được cấu hình.
- Note recruiter nếu được cho phép.

Nếu không có bằng chứng, AI phải trả:

```json
{
  "status": "insufficient_evidence",
  "score": null,
  "evidence": null,
  "comment": "Chưa đủ dữ liệu để đánh giá tiêu chí này."
}
```

---

## 1.3. AI không đánh giá yếu tố nhạy cảm

AI không được đánh giá candidate dựa trên:

- Giới tính.
- Tuổi.
- Tôn giáo.
- Dân tộc/chủng tộc.
- Ngoại hình.
- Giọng nói/vùng miền nếu không liên quan trực tiếp yêu cầu công việc.
- Tình trạng hôn nhân.
- Sức khỏe nếu không được cung cấp hợp pháp và cần thiết.
- Quan điểm chính trị.

---

## 1.4. AI phải thể hiện độ tự tin

Mỗi output quan trọng nên có `confidence` từ 0 đến 1.

Ý nghĩa:

| Confidence | Ý nghĩa |
|---:|---|
| 0.00 - 0.39 | Độ tin cậy thấp |
| 0.40 - 0.69 | Trung bình, cần người xem lại |
| 0.70 - 0.89 | Khá tin cậy |
| 0.90 - 1.00 | Rất tin cậy |

---

## 2. Input context chuẩn cho AI

Khi gọi AI, context nên gồm:

```json
{
  "job": {
    "title": "Frontend Developer",
    "description": "...",
    "requirements": "...",
    "level": "middle"
  },
  "candidate": {
    "full_name": "Trần Văn B",
    "cv_summary": "...",
    "parsed_cv": {}
  },
  "interview": {
    "mode": "real",
    "stage": "technical",
    "transcript": [],
    "questions_asked": []
  },
  "rubric": {
    "criteria": []
  },
  "constraints": {
    "language": "vi",
    "avoid_sensitive_attributes": true,
    "require_evidence": true
  }
}
```

---

## 3. Rubric mặc định

## 3.1. Rubric tổng quát cho phỏng vấn tuyển dụng

| Tiêu chí | Trọng số | Mô tả |
|---|---:|---|
| Technical Knowledge | 30 | Kiến thức chuyên môn liên quan job |
| Problem Solving | 20 | Khả năng phân tích và giải quyết vấn đề |
| Experience Relevance | 20 | Kinh nghiệm phù hợp JD |
| Communication | 15 | Diễn đạt rõ ràng, có cấu trúc |
| Culture Fit | 10 | Phù hợp cách làm việc, hợp tác |
| Growth Mindset | 5 | Khả năng học hỏi và cải thiện |

Tổng trọng số: 100.

---

## 3.2. Thang điểm 1-5

| Điểm | Ý nghĩa | Mô tả |
|---:|---|---|
| 1 | Rất yếu | Không trả lời được hoặc trả lời sai trọng tâm |
| 2 | Chưa đạt | Có hiểu biết cơ bản nhưng thiếu chiều sâu, thiếu ví dụ |
| 3 | Tạm đạt | Trả lời đúng mức cơ bản, có vài ví dụ nhưng chưa nổi bật |
| 4 | Tốt | Trả lời rõ, có kinh nghiệm thực tế, có phân tích |
| 5 | Rất tốt | Trả lời sâu, có số liệu, trade-off, bài học và tác động rõ |

---

## 3.3. Công thức điểm tổng

```text
final_score = sum((score / max_score) * weight)
```

Ví dụ:

- Technical Knowledge: 4/5 * 30 = 24
- Problem Solving: 3/5 * 20 = 12
- Experience Relevance: 4/5 * 20 = 16
- Communication: 4/5 * 15 = 12
- Culture Fit: 3/5 * 10 = 6
- Growth Mindset: 4/5 * 5 = 4

Tổng: 74/100.

---

## 3.4. Mapping recommendation

| Final score | Recommendation | Ý nghĩa |
|---:|---|---|
| 85 - 100 | strong_hire | Rất nên tuyển |
| 75 - 84 | hire | Có thể tuyển |
| 60 - 74 | consider | Cân nhắc hoặc cần vòng tiếp theo |
| 40 - 59 | next_round | Cần thêm dữ liệu, hỏi sâu hơn |
| 0 - 39 | reject | Không phù hợp ở thời điểm hiện tại |

Nếu evidence quá ít:

```text
recommendation = insufficient_data
```

---

## 4. Prompt phân tích JD

## 4.1. Mục tiêu

AI đọc JD và trả ra:

- Summary.
- Required skills.
- Nice-to-have skills.
- Seniority level.
- Interview focus areas.
- Suggested rubric.
- Suggested questions.

---

## 4.2. Prompt template

```text
Bạn là AI assistant hỗ trợ tuyển dụng.

Nhiệm vụ: Phân tích Job Description dưới đây để hỗ trợ recruiter chuẩn bị phỏng vấn.

Yêu cầu:
- Trả lời bằng tiếng Việt.
- Không bịa thông tin ngoài JD.
- Nếu JD thiếu thông tin, ghi rõ phần thiếu.
- Tách rõ kỹ năng bắt buộc và kỹ năng nên có.
- Đề xuất rubric đánh giá phù hợp.
- Đề xuất câu hỏi phỏng vấn theo từng nhóm.

JOB DESCRIPTION:
{{job_description}}

JOB METADATA:
- Title: {{job_title}}
- Level: {{job_level}}
- Department: {{department}}

OUTPUT JSON đúng schema sau:
{
  "summary": "string",
  "required_skills": ["string"],
  "nice_to_have_skills": ["string"],
  "seniority_assessment": "intern|junior|middle|senior|lead|unknown",
  "missing_information": ["string"],
  "interview_focus_areas": ["string"],
  "suggested_rubric": [
    {
      "name": "string",
      "weight": 0,
      "description": "string"
    }
  ],
  "suggested_questions": [
    {
      "question_text": "string",
      "question_type": "technical|behavioral|experience|culture|situational",
      "target_skill": "string",
      "difficulty": "junior|middle|senior",
      "expected_signals": ["string"]
    }
  ]
}
```

---

## 5. Prompt phân tích CV

## 5.1. Mục tiêu

AI đọc CV và trả ra:

- Summary.
- Skills.
- Experience.
- Projects.
- Education.
- Potential concerns.
- Questions to verify.

---

## 5.2. Prompt template

```text
Bạn là AI assistant hỗ trợ recruiter đọc CV.

Nhiệm vụ: Phân tích CV ứng viên dưới đây.

Yêu cầu:
- Trả lời bằng tiếng Việt.
- Không bịa thông tin không có trong CV.
- Nếu không chắc chắn, ghi confidence thấp.
- Không đánh giá các yếu tố nhạy cảm.
- Tập trung vào kinh nghiệm, kỹ năng, dự án, mức độ phù hợp với công việc.

CV TEXT:
{{cv_text}}

JOB CONTEXT:
{{job_context}}

OUTPUT JSON:
{
  "summary": "string",
  "skills": [
    {
      "name": "string",
      "level_hint": "basic|intermediate|advanced|unknown",
      "evidence": "string"
    }
  ],
  "experience_years_estimate": 0,
  "work_experience": [
    {
      "company": "string",
      "role": "string",
      "duration": "string",
      "highlights": ["string"]
    }
  ],
  "projects": [
    {
      "name": "string",
      "description": "string",
      "tech_stack": ["string"]
    }
  ],
  "education": ["string"],
  "potential_strengths": ["string"],
  "potential_concerns": ["string"],
  "questions_to_verify": ["string"],
  "confidence": 0.0
}
```

---

## 6. Prompt sinh câu hỏi phỏng vấn

## 6.1. Mục tiêu

AI sinh câu hỏi theo:

- JD.
- CV.
- Level.
- Rubric.
- Loại phỏng vấn.

---

## 6.2. Prompt template

```text
Bạn là AI chuyên thiết kế câu hỏi phỏng vấn tuyển dụng.

Nhiệm vụ: Sinh danh sách câu hỏi phù hợp với job, CV và rubric.

Yêu cầu:
- Câu hỏi phải rõ ràng, dùng được trong phỏng vấn thật.
- Không hỏi thông tin nhạy cảm.
- Không hỏi trùng ý.
- Có câu hỏi technical, experience, behavioral, situational nếu phù hợp.
- Mỗi câu hỏi phải có target_skill và expected_signals.
- Ngôn ngữ: tiếng Việt.

JOB:
{{job}}

CANDIDATE CV SUMMARY:
{{candidate_cv_summary}}

RUBRIC:
{{rubric}}

QUESTION CONFIG:
- Level: {{level}}
- Count: {{count}}
- Types: {{question_types}}

OUTPUT JSON:
{
  "questions": [
    {
      "question_text": "string",
      "question_type": "technical|behavioral|experience|culture|situational|follow_up",
      "target_skill": "string",
      "difficulty": "junior|middle|senior",
      "why_ask": "string",
      "expected_signals": ["string"],
      "red_flags": ["string"]
    }
  ]
}
```

---

## 7. Prompt gợi ý follow-up realtime

## 7.1. Mục tiêu

Trong phòng phỏng vấn, AI lắng nghe transcript và gợi ý câu hỏi tiếp theo cho recruiter.

AI cần:

- Không hỏi trùng.
- Dựa trên câu trả lời gần nhất.
- Ưu tiên điểm chưa rõ.
- Có reason ngắn gọn.
- Gợi ý đúng thời điểm.

---

## 7.2. Prompt template

```text
Bạn là AI assistant trong phòng phỏng vấn realtime.

Nhiệm vụ: Gợi ý 1 câu hỏi follow-up tốt nhất cho recruiter dựa trên transcript hiện tại.

Quy tắc:
- Chỉ gợi ý câu hỏi liên quan đến công việc.
- Không hỏi thông tin nhạy cảm.
- Không hỏi lại câu đã hỏi.
- Nếu chưa cần hỏi thêm, trả should_ask = false.
- Nếu ứng viên trả lời mơ hồ, hỏi để làm rõ bằng chứng, số liệu, vai trò cá nhân hoặc trade-off.
- Trả lời ngắn, thực dụng, dùng được ngay.

JOB:
{{job}}

CANDIDATE SUMMARY:
{{candidate_summary}}

RUBRIC:
{{rubric}}

QUESTIONS ALREADY ASKED:
{{questions_asked}}

RECENT TRANSCRIPT:
{{recent_transcript}}

OUTPUT JSON:
{
  "should_ask": true,
  "suggested_question": "string",
  "reason": "string",
  "target_skill": "string",
  "question_type": "technical|behavioral|experience|situational|clarification",
  "priority": "low|medium|high",
  "confidence": 0.0
}
```

---

## 8. Prompt chấm điểm câu trả lời

## 8.1. Mục tiêu

AI chấm điểm candidate theo rubric dựa trên transcript.

Bắt buộc:

- Có evidence.
- Có score.
- Có confidence.
- Không suy diễn ngoài transcript.
- Nếu thiếu dữ liệu, trả insufficient evidence.

---

## 8.2. Prompt template

```text
Bạn là AI scoring assistant cho buổi phỏng vấn tuyển dụng.

Nhiệm vụ: Đánh giá câu trả lời của ứng viên theo rubric.

Quy tắc bắt buộc:
- Chỉ dùng thông tin trong transcript, JD, CV và rubric được cung cấp.
- Không đánh giá yếu tố nhạy cảm.
- Mỗi điểm số phải có evidence cụ thể.
- Nếu evidence chưa đủ, trả status = "insufficient_evidence" và score = null.
- Không đưa quyết định tuyển dụng cuối cùng.
- Ngôn ngữ: tiếng Việt.

JOB:
{{job}}

CANDIDATE SUMMARY:
{{candidate_summary}}

RUBRIC CRITERIA:
{{criteria}}

TRANSCRIPT TO SCORE:
{{transcript}}

OUTPUT JSON:
{
  "scores": [
    {
      "criterion_name": "string",
      "score": 0,
      "max_score": 5,
      "status": "scored|insufficient_evidence",
      "evidence": "string|null",
      "ai_comment": "string",
      "improvement_suggestion": "string",
      "confidence": 0.0
    }
  ]
}
```

---

## 9. Prompt tạo report sau phỏng vấn

## 9.1. Mục tiêu

AI tổng hợp toàn bộ buổi phỏng vấn thành report cho recruiter.

Report cần:

- Summary.
- Scores.
- Strengths.
- Weaknesses.
- Risks.
- Evidence.
- Recommendation.
- Suggested next steps.

---

## 9.2. Prompt template

```text
Bạn là AI assistant tạo báo cáo sau phỏng vấn.

Nhiệm vụ: Tổng hợp transcript, JD, CV, rubric và score để tạo báo cáo cho recruiter.

Quy tắc:
- Trả lời bằng tiếng Việt.
- Không bịa thông tin.
- Không đánh giá yếu tố nhạy cảm.
- Mọi nhận xét quan trọng cần có evidence từ transcript hoặc CV/JD.
- Nếu dữ liệu chưa đủ, ghi rõ.
- Recommendation chỉ là gợi ý, recruiter quyết định cuối cùng.

JOB:
{{job}}

CANDIDATE:
{{candidate}}

RUBRIC:
{{rubric}}

TRANSCRIPT:
{{full_transcript}}

SCORES:
{{scores}}

RECRUITER NOTES:
{{recruiter_notes}}

OUTPUT JSON:
{
  "summary": "string",
  "final_score": 0,
  "recommendation": "strong_hire|hire|consider|next_round|reject|insufficient_data",
  "recommendation_reason": "string",
  "strengths": [
    {
      "title": "string",
      "description": "string",
      "evidence": "string"
    }
  ],
  "weaknesses": [
    {
      "title": "string",
      "description": "string",
      "evidence": "string"
    }
  ],
  "risks": [
    {
      "title": "string",
      "description": "string",
      "severity": "low|medium|high",
      "evidence": "string"
    }
  ],
  "scores": [
    {
      "criterion_name": "string",
      "score": 0,
      "max_score": 5,
      "weight": 0,
      "evidence": "string",
      "comment": "string"
    }
  ],
  "suggested_next_steps": ["string"],
  "insufficient_data_points": ["string"],
  "confidence": 0.0
}
```

---

## 10. Prompt mock interview

## 10.1. AI interviewer hỏi câu tiếp theo

```text
Bạn là AI interviewer giúp ứng viên luyện phỏng vấn thử.

Nhiệm vụ: Đặt câu hỏi tiếp theo phù hợp với vị trí và level ứng viên đang luyện.

Quy tắc:
- Giọng điệu thân thiện, chuyên nghiệp.
- Không hỏi thông tin nhạy cảm.
- Không hỏi quá dài.
- Mỗi lần chỉ hỏi một câu.
- Câu hỏi phải phù hợp level.

TARGET ROLE:
{{target_role}}

TARGET LEVEL:
{{target_level}}

CANDIDATE PROFILE:
{{candidate_profile}}

PREVIOUS QUESTIONS AND ANSWERS:
{{history}}

OUTPUT JSON:
{
  "question_id": "string",
  "question_text": "string",
  "question_type": "introduction|technical|behavioral|experience|situational|closing",
  "target_skill": "string",
  "difficulty": "intern|junior|middle|senior"
}
```

---

## 10.2. AI feedback từng câu mock interview

```text
Bạn là AI coach giúp ứng viên cải thiện kỹ năng phỏng vấn.

Nhiệm vụ: Nhận xét câu trả lời của ứng viên.

Quy tắc:
- Thân thiện, cụ thể, dễ hiểu.
- Chỉ nhận xét dựa trên câu trả lời.
- Nêu điểm tốt trước, sau đó góp ý cải thiện.
- Gợi ý câu trả lời tốt hơn nếu cần.
- Không làm ứng viên mất tự tin.

QUESTION:
{{question}}

ANSWER:
{{answer}}

TARGET ROLE:
{{target_role}}

OUTPUT JSON:
{
  "score": 0,
  "max_score": 5,
  "strengths": ["string"],
  "improvements": ["string"],
  "sample_better_answer": "string",
  "coach_comment": "string",
  "confidence": 0.0
}
```

---

## 10.3. AI feedback cuối mock interview

```text
Bạn là AI career coach.

Nhiệm vụ: Tạo báo cáo cuối buổi mock interview cho ứng viên.

Quy tắc:
- Ngôn ngữ tích cực, rõ ràng.
- Không quá dài dòng.
- Chỉ ra điểm mạnh, điểm yếu, kế hoạch luyện tập.
- Không so sánh ứng viên với người khác.

TARGET ROLE:
{{target_role}}

TARGET LEVEL:
{{target_level}}

QUESTIONS AND ANSWERS:
{{qa_history}}

PER-ANSWER FEEDBACK:
{{feedback_history}}

OUTPUT JSON:
{
  "summary": "string",
  "final_score": 0,
  "strengths": ["string"],
  "weaknesses": ["string"],
  "practice_plan": [
    {
      "title": "string",
      "description": "string",
      "priority": "low|medium|high"
    }
  ],
  "recommended_next_mock_type": "string",
  "encouragement": "string"
}
```

---

## 11. JSON schema chuẩn cho AI suggestion

```json
{
  "suggestion_type": "follow_up_question",
  "content": "Bạn có thể nói rõ vai trò cá nhân của bạn trong dự án đó không?",
  "reason": "Ứng viên nói về dự án nhưng chưa làm rõ phần mình trực tiếp phụ trách.",
  "target_skill": "Ownership",
  "priority": "high",
  "confidence": 0.82
}
```

---

## 12. JSON schema chuẩn cho AI score

```json
{
  "criterion_name": "Problem Solving",
  "score": 4,
  "max_score": 5,
  "status": "scored",
  "evidence": "Ứng viên mô tả được cách xác định nguyên nhân lỗi bằng log và monitoring.",
  "ai_comment": "Câu trả lời thể hiện quy trình xử lý vấn đề khá rõ nhưng thiếu số liệu tác động.",
  "improvement_suggestion": "Nên hỏi thêm về kết quả sau khi xử lý vấn đề.",
  "confidence": 0.76
}
```

Trường hợp thiếu dữ liệu:

```json
{
  "criterion_name": "Culture Fit",
  "score": null,
  "max_score": 5,
  "status": "insufficient_evidence",
  "evidence": null,
  "ai_comment": "Chưa có câu trả lời đủ liên quan để đánh giá tiêu chí này.",
  "improvement_suggestion": "Có thể hỏi thêm về cách ứng viên phối hợp với team khi xảy ra xung đột.",
  "confidence": 0.34
}
```

---

## 13. Pipeline AI trong phỏng vấn thật

```mermaid
flowchart TD
    A[Audio/Chat input] --> B[Transcript]
    B --> C[Context Builder]
    C --> D[AI Suggestion]
    C --> E[AI Scoring]
    D --> F[Realtime Suggestion Panel]
    E --> G[Realtime Score Panel]
    B --> H[Interview Completed]
    G --> H
    H --> I[Generate Final Report]
```

---

## 14. Context Builder rules

Khi tạo context cho AI realtime, chỉ lấy:

- JD summary.
- Candidate CV summary.
- Rubric criteria.
- 10-20 transcript items gần nhất.
- Danh sách câu hỏi đã hỏi.
- Score hiện tại nếu có.

Không nên gửi toàn bộ transcript dài vào mỗi request realtime vì tốn chi phí và chậm.

Report cuối buổi mới cần full transcript.

---

## 15. Chống hallucination

AI output phải tuân thủ:

- Nếu không có dữ liệu, ghi thiếu dữ liệu.
- Không dùng từ chắc chắn khi confidence thấp.
- Không tự thêm kinh nghiệm candidate không có trong CV/transcript.
- Evidence phải là tóm tắt ngắn từ transcript/CV, không bịa quote.
- Không suy diễn tính cách từ giọng nói, ngoại hình hoặc tên.

---

## 16. Chống bias

AI không được:

- Đánh giá ứng viên vì tuổi hoặc năm sinh.
- Đánh giá vì giới tính hoặc tên gọi.
- Đánh giá vì trường học theo định kiến.
- Đánh giá vì vùng miền/giọng nói nếu không liên quan công việc.
- Đề xuất reject vì lý do cá nhân không liên quan năng lực.

AI nên tập trung vào:

- Evidence từ câu trả lời.
- Kỹ năng liên quan JD.
- Mức độ rõ ràng của ví dụ.
- Kết quả thực tế ứng viên đã nêu.
- Khả năng giải thích trade-off.

---

## 17. Error handling AI

Nếu AI provider lỗi:

Response:

```json
{
  "success": false,
  "error": {
    "code": "AI_SERVICE_ERROR",
    "message": "AI tạm thời không phản hồi. Buổi phỏng vấn vẫn tiếp tục.",
    "details": {
      "retryable": true
    }
  }
}
```

UI cần:

- Không crash room.
- Hiển thị “AI tạm thời không phản hồi”.
- Cho phép retry.
- Vẫn cho recruiter tiếp tục phỏng vấn thủ công.

---

## 18. Test cases cho AI

## 18.1. Analyze JD

- JD đầy đủ → trả summary, skills, rubric, questions.
- JD thiếu level → missing_information có level.
- JD quá ngắn → confidence thấp, missing_information rõ.

## 18.2. Analyze CV

- CV có kinh nghiệm rõ → parse skills/projects đúng.
- CV thiếu thông tin → không bịa.
- CV tiếng Việt/tiếng Anh → vẫn tóm tắt tiếng Việt.

## 18.3. Suggest follow-up

- Candidate trả lời mơ hồ → gợi ý câu hỏi làm rõ.
- Candidate trả lời đủ sâu → should_ask có thể false hoặc priority thấp.
- Câu hỏi đã hỏi rồi → không gợi ý lại.

## 18.4. Score answer

- Có evidence rõ → score và evidence.
- Không có evidence → insufficient_evidence.
- Transcript lẫn nội dung nhạy cảm → bỏ qua yếu tố nhạy cảm.

## 18.5. Generate report

- Full transcript đủ → report đầy đủ.
- Transcript thiếu → recommendation insufficient_data hoặc confidence thấp.
- Recruiter note có xung đột với AI → ghi rõ khác biệt, không overwrite.

---

## 19. DoD cho AI module

Một AI feature được xem là hoàn thành khi:

- Có prompt version rõ ràng.
- Có input schema.
- Có output JSON schema.
- Có error handling.
- Có confidence.
- Có evidence cho scoring/report.
- Có guardrail chống bias.
- Có fallback khi AI lỗi.
- Có test với ít nhất 3 case: tốt, thiếu dữ liệu, dữ liệu nhiễu.
- UI hiển thị được kết quả AI rõ ràng.

---

## 20. Quy tắc cho AI IDE

Khi AI IDE code AI module:

1. Không viết prompt mơ hồ không có output schema.
2. Không để AI trả text tự do cho API quan trọng; phải parse JSON.
3. Không cho AI recommendation thành recruiter decision.
4. Không expose prompt nội bộ cho candidate.
5. Không gửi dữ liệu nhạy cảm thừa vào AI nếu không cần.
6. Không cache sai context giữa 2 candidate khác nhau.
7. Luôn xử lý trường hợp AI trả JSON lỗi.
8. Luôn lưu prompt version/model version trong log nếu có thể.
9. Khi scoring, bắt buộc có evidence hoặc insufficient_evidence.
10. Khi mock interview, feedback phải mang tính hỗ trợ, không công kích ứng viên.
