# Câu hỏi giảng viên có thể hỏi — AI Interview Platform

> Gợi ý ôn trước khi bảo vệ / demo. Trả lời ngắn, rõ; có thể mở rộng khi được hỏi sâu.

---

## A. Sản phẩm & nghiệp vụ

### 1. Dự án của nhóm giải quyết bài toán gì?

**Trả lời:**  
Giúp nhà tuyển dụng phỏng vấn online có **dữ liệu và nhất quán hơn**, đồng thời giúp ứng viên **luyện phỏng vấn** với AI. Hệ thống hỗ trợ: tạo job/JD, quản lý candidate & CV, phòng phỏng vấn realtime, transcript, AI gợi ý câu hỏi, chấm theo rubric, và báo cáo sau buổi PV. AI **không thay** quyết định tuyển dụng.

---

### 2. Khác gì so với Zoom / Google Meet + chấm tay?

**Trả lời:**  
Meet chỉ là kênh gọi video. Hệ thống của nhóm gắn **ngữ cảnh tuyển dụng**: JD, CV, rubric, transcript lưu lại, AI gợi ý follow-up, scoring có evidence, report sau buổi. Recruiter có pipeline candidate và quyết định trên dữ liệu, không chỉ “gọi xong rồi quên”.

---

### 3. Ai là người dùng chính? Mỗi role làm gì?

**Trả lời:**  
- **Recruiter:** quản lý company, job, candidate, lịch PV, vào phòng, xem AI/report, quyết định.  
- **Candidate:** vào phòng qua link, trả lời; luyện mock; (mở rộng) tìm việc/apply.  
- **Admin:** users, logs, cấu hình AI.  
Candidate **không** xem điểm/note/report nội bộ của recruiter.

---

### 4. Mô tả luồng nghiệp vụ chính từ đầu đến cuối.

**Trả lời:**  
Recruiter tạo company → tạo job + JD → AI phân tích JD → thêm candidate + upload CV → AI parse CV → tạo lịch + link mời → hai bên vào phòng → start → transcript + AI gợi ý/chấm → end → AI tạo report → recruiter quyết định pass/consider/reject.

---

### 5. Mock interview khác phỏng vấn thật thế nào?

**Trả lời:**  
Mock là **luyện tập**: candidate chọn vị trí/level, AI hỏi và cho feedback. Kết quả mock **không** dùng để tuyển dụng thật. Phỏng vấn thật có recruiter, room realtime, scoring/report nội bộ cho phía tuyển dụng.

---

### 6. Tại sao AI không được tự động loại ứng viên?

**Trả lời:**  
Tuyển dụng ảnh hưởng quyền lợi con người; AI có thể sai hoặc thiên lệch. Nhóm thiết kế AI chỉ **hỗ trợ** (gợi ý, chấm có evidence, report). **Người** luôn xác nhận quyết định cuối — tránh rủi ro pháp lý và bias.

---

## B. AI & đạo đức

### 7. AI trong hệ thống làm những việc gì cụ thể?

**Trả lời:**  
Phân tích JD, parse/tóm tắt CV, sinh/gợi ý câu hỏi (kể cả follow-up realtime), chấm theo rubric dựa trên transcript, tạo interview report, feedback mock interview. Không đánh giá qua khuôn mặt, cảm xúc camera, hay giọng vùng miền.

---

### 8. Làm sao giảm thiên lệch (bias) của AI?

**Trả lời:**  
- Không dùng tín hiệu nhạy cảm (tuổi, giới, ngoại hình, giọng…).  
- Chấm theo **rubric** gắn JD/CV/transcript, yêu cầu **evidence**.  
- Không auto-reject.  
- Recruiter thấy và có thể không đồng ý với AI.  
- Audit log các hành động quan trọng.

---

### 9. “Evidence” trong scoring nghĩa là gì?

**Trả lời:**  
Mỗi điểm/kết luận AI phải gắn với **bằng chứng** lấy từ câu trả lời trong transcript (hoặc nội dung CV/JD), không bịa. Giảng viên/recruiter có thể đối chiếu lại.

---

### 10. Hệ thống có lưu video không? Vì sao?

**Trả lời:**  
**Không lưu video** trong scope hiện tại. Chỉ dùng realtime A/V qua LiveKit và lưu **transcript** (và metadata cần thiết). Lý do: chi phí lưu trữ, rủi ro PII/pháp lý, và scope MVP tập trung vào nội dung trả lời chứ không phân tích hình ảnh.

---

## C. Kiến trúc & kỹ thuật

### 11. Vì sao tách Frontend / Backend / AI Service / Realtime?

**Trả lời:**  
- **Frontend (Vue):** UI cho recruiter/candidate/admin.  
- **Backend API (Go):** nghiệp vụ, auth, CRUD, quyền.  
- **Realtime (Go WebSocket):** phòng PV, event join/leave/chat/transcript.  
- **AI Service (Python FastAPI):** gọi LLM (Gemini), prompt, scoring — Python phù hợp ML/LLM; Go phù hợp API/concurrency.  
Tách service để scale độc lập, lỗi AI không làm sập toàn bộ API.

---

### 12. LiveKit dùng để làm gì? WebSocket dùng để làm gì?

**Trả lời:**  
- **LiveKit (SFU):** truyền audio/video realtime giữa recruiter và candidate.  
- **WebSocket:** sự kiện nghiệp vụ trong phòng (trạng thái, chat, transcript, AI suggestion) — không thay thế media pipeline.

---

### 13. Multi-tenant nghĩa là gì trong dự án?

**Trả lời:**  
Mỗi company là một workspace. Job, candidate, interview **luôn thuộc một company**. Recruiter không xem dữ liệu company khác — bảo vệ dữ liệu tuyển dụng giữa các tổ chức.

---

### 14. Auth / bảo mật cơ bản nhóm xử lý thế nào?

**Trả lời:**  
JWT access token + refresh token (HttpOnly cookie, token family). API bảo vệ theo role. Candidate vào phòng bằng **invite token** có kiểm tra hết hạn. Phân quyền: candidate không đọc scoring/report nội bộ. Có audit log.

---

### 15. Database dùng gì? Entity chính?

**Trả lời:**  
PostgreSQL. Entity chính: User, Company, Job, Candidate, Interview, Room/Transcript, Rubric/Score, Report, File/CV, AI logs. Redis dùng cache/queue hỗ trợ realtime/ops.

---

## D. Team, scope, hạn chế

### 16. Nhóm phân công thế nào?

**Trả lời:**  
- **Hùng:** Realtime, Room, WebSocket, Transcript.  
- **Khôi:** Backend core, DB, AI orchestration, Scoring, Report.  
- **Lai:** UI/UX, portal, business flow trên frontend.

---

### 17. MVP gồm gì? Cố ý chưa làm gì?

**Trả lời:**  
**Có:** auth, company, job, candidate/CV, lịch PV, phòng realtime, AI gợi ý & scoring cơ bản, report, mock, audit.  
**Chưa / ngoài scope:** ATS/LinkedIn sâu, calendar Google/Microsoft đầy đủ, AI avatar 3D, coding test IDE, SSO enterprise, lưu video, đánh giá cảm xúc khuôn mặt, AI tự gửi mail từ chối.

---

### 18. Điểm mạnh / điểm yếu của đồ án hiện tại?

**Trả lời gợi ý:**  
**Mạnh:** luồng tuyển dụng end-to-end rõ; AI gắn nghiệp vụ (không chỉ chatbot); realtime room; nguyên tắc AI có kiểm soát.  
**Yếu / hạn chế:** phụ thuộc API Gemini (chi phí, latency); chấm AI chưa hoàn hảo; chưa tích hợp ATS sâu; cần thêm dữ liệu thực tế để đo accuracy scoring; vận hành production còn cần giám sát/ổn định.

---

### 19. Làm sao đo hệ thống “thành công”?

**Trả lời:**  
Ví dụ metric trong PRD: giảm thời gian sàng lọc; tỷ lệ hoàn thành buổi PV; độ trễ transcript; độ tương đồng AI vs đánh giá recruiter; thời gian tạo report; tỷ lệ hoàn thành mock. Cần chạy pilot để chỉnh số liệu.

---

### 20. Hướng phát triển tiếp theo?

**Trả lời:**  
Ổn định production (deploy, monitoring, AI settings/ops); cải thiện độ chính xác scoring & tiếng Việt; tích hợp calendar/ATS nếu có nhu cầu; làm sâu candidate portal/job board một cách có kiểm soát; mở rộng multi-round hiring khi nghiệp vụ sẵn sàng.

---

## E. Câu hỏi “đào” khi demo

### 21. Nếu transcript sai / AI chấm lệch thì sao?

**Trả lời:**  
Recruiter vẫn nghe trực tiếp và ghi note; AI chỉ tham khảo. Report có evidence để đối chiếu. Có thể cải thiện bằng prompt, rubric rõ hơn, và (tương lai) cho phép recruiter chỉnh/ghi nhận sai lệch — không lấy AI làm nguồn sự thật duy nhất.

---

### 22. Candidate có biết mình bị AI chấm không?

**Trả lời:**  
Trước khi vào phòng có bước **consent** (ghi âm/AI) nếu hệ thống bật. Candidate biết có hỗ trợ AI trong buổi, nhưng **không** xem điểm/report nội bộ trừ khi công ty cho phép.

---

### 23. Nếu nhiều ứng viên cùng lúc / nhiều phòng?

**Trả lời:**  
Kiến trúc tách API / Realtime / LiveKit / AI giúp scale từng lớp. Mỗi interview có room và invite riêng. Thực tế production cần giới hạn concurrent, queue gọi AI, và giám sát tài nguyên — đây là hướng vận hành tiếp theo.

---

### 24. Dữ liệu CV/JD nhạy cảm lưu ở đâu?

**Trả lời:**  
Metadata và nội dung phân tích lưu PostgreSQL; file CV qua storage (S3-compatible / uploads), không nhét binary vào DB. Truy cập theo company và quyền. Production dùng HTTPS, auth, và kiểm soát AI settings/API key ở tầng admin.

---

### 25. Tại sao chọn Vue + Go + Python, không all-in một stack?

**Trả lời:**  
Mỗi tầng đúng sở trường: Vue cho SPA nhanh; Go cho API/realtime hiệu năng; Python cho gọi LLM và xử lý AI. Docker Compose gom chạy local/deploy. Trade-off: nhiều service hơn nhưng rõ trách nhiệm và dễ chia việc trong nhóm 3 người.

---

## F. Câu trả lời “an toàn” khi bí

| Tình huống | Cách nói |
|---|---|
| Ngoài scope | “Phần này cố ý để phase sau / ngoài MVP vì …; hiện nhóm ưu tiên luồng PV cốt lõi.” |
| Chưa đo được số | “Đây là mục tiêu thiết kế trong PRD; cần pilot thực tế để xác nhận.” |
| AI sai | “AI hỗ trợ có evidence; quyết định cuối thuộc recruiter.” |
| So với sản phẩm thương mại | “Nhóm làm MVP học thuật/demo đủ end-to-end; chưa cạnh tranh feature enterprise.” |

---

## Gợi ý ôn nhanh 1 phút

1. Bài toán + AI hỗ trợ không thay người.  
2. 3 role + luồng recruiter 1 vòng.  
3. Vue / Go API / Go Realtime / Python AI / Postgres / LiveKit.  
4. Không bias, không auto-reject, không lưu video.  
5. Phân công Hùng–Khôi–Lai + hướng phát triển tiếp.
