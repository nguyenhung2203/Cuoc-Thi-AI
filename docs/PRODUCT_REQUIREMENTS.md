# PRODUCT_REQUIREMENTS.md
# Yêu cầu sản phẩm — AI Interview Platform

## 0. Mục đích tài liệu

Tài liệu này là nguồn sự thật chính về sản phẩm để 3 người trong team và AI IDE cùng hiểu đúng phạm vi cần làm.

Tài liệu dùng cho:

- Hùng: làm realtime interview room, transcript realtime, AI realtime trong phòng.
- Khôi: làm backend core, database, API, AI orchestration, report.
- Lai: làm UI/UX, dashboard, candidate portal, business flow, tích hợp API.

AI IDE khi bắt đầu một task phải đọc file này trước để hiểu:

- Sản phẩm là gì.
- Actor nào được làm gì.
- MVP gồm những gì.
- Không được tự ý bịa thêm tính năng ngoài scope.
- Luồng nghiệp vụ chuẩn của recruiter và candidate.

---

## 1. Tầm nhìn sản phẩm

**AI Interview Platform** là nền tảng phỏng vấn trực tuyến có AI real-time hỗ trợ nhà tuyển dụng và ứng viên.

Hệ thống phục vụ 2 mục tiêu chính:

1. **Hỗ trợ nhà tuyển dụng lọc ứng viên nhanh, có dữ liệu và nhất quán hơn.**
2. **Hỗ trợ ứng viên luyện phỏng vấn thử với AI và nhận feedback chi tiết.**

AI trong sản phẩm không thay thế quyết định tuyển dụng của con người. AI chỉ hỗ trợ:

- Ghi transcript.
- Gợi ý câu hỏi.
- Chấm điểm theo rubric.
- Trích dẫn bằng chứng từ câu trả lời.
- Tạo báo cáo sau phỏng vấn.
- Feedback cho mock interview.

---

## 2. Actor chính

## 2.1. Recruiter

Recruiter là người đại diện phía tuyển dụng.

Recruiter có thể:

- Đăng nhập hệ thống.
- Quản lý công ty/workspace.
- Tạo job.
- Thêm candidate.
- Upload CV candidate.
- Tạo lịch phỏng vấn.
- Mời candidate tham gia phòng phỏng vấn.
- Vào phòng phỏng vấn.
- Bắt đầu/kết thúc buổi phỏng vấn.
- Xem AI suggestion.
- Xem AI scoring.
- Ghi chú nội bộ.
- Xem report sau phỏng vấn.
- Ra quyết định cuối cùng.

Recruiter không nên:

- Xem dữ liệu candidate thuộc công ty khác.
- Sửa transcript gốc mà không ghi audit log.
- Cho AI tự động reject candidate mà không có xác nhận người dùng.

---

## 2.2. Candidate

Candidate là người tham gia phỏng vấn hoặc luyện phỏng vấn thử.

Candidate có thể:

- Nhận link mời phỏng vấn.
- Xác nhận thông tin trước khi vào phòng.
- Upload hoặc cập nhật CV nếu được cho phép.
- Vào phòng phỏng vấn.
- Bật/tắt mic, camera của bản thân.
- Chat trong phòng.
- Trả lời câu hỏi.
- Tham gia mock interview với AI.
- Xem feedback mock interview của chính mình.

Candidate không được:

- Xem AI scoring nội bộ của recruiter.
- Xem note riêng của recruiter.
- Xem report tuyển dụng nội bộ nếu công ty không cho phép.
- Vào phòng khi link hết hạn hoặc không đúng token.

---

## 2.3. AI Assistant

AI là hệ thống hỗ trợ, không phải actor quyết định nghiệp vụ cuối cùng.

AI có thể:

- Phân tích JD.
- Phân tích CV.
- Sinh câu hỏi.
- Gợi ý câu hỏi follow-up realtime.
- Tóm tắt transcript.
- Chấm điểm theo rubric.
- Tạo report.
- Đưa feedback mock interview.

AI không được:

- Tự động loại ứng viên.
- Đánh giá dựa trên tuổi, giới tính, tôn giáo, ngoại hình, vùng miền, tình trạng hôn nhân hoặc đặc điểm nhạy cảm.
- Tự bịa thông tin không có trong JD, CV hoặc transcript.
- Đưa kết luận chắc chắn khi dữ liệu chưa đủ.

---

## 3. MVP Scope

## 3.1. MVP bắt buộc có

| Nhóm | Chức năng | Mô tả | Owner chính |
|---|---|---|---|
| Auth | Đăng nhập/đăng ký | Recruiter và candidate đăng nhập được | Khôi |
| Workspace | Company cơ bản | Recruiter quản lý công ty/workspace | Khôi |
| Job | Job CRUD | Tạo/sửa/xem/xóa job | Khôi + Lai |
| Candidate | Candidate CRUD | Thêm/xem/sửa candidate, upload CV | Khôi + Lai |
| Interview | Scheduling | Tạo lịch phỏng vấn, sinh room link | Khôi + Lai |
| Room | Interview Room | Recruiter/candidate vào phòng realtime | Hùng + Lai |
| Realtime | WebSocket events | Join/leave, status, chat, transcript, AI suggestion | Hùng |
| AI | Question suggestion | AI gợi ý câu hỏi theo JD/CV/transcript | Hùng + Khôi |
| AI | Scoring cơ bản | AI chấm điểm theo rubric | Khôi |
| Report | Interview report | Tạo report sau phỏng vấn | Khôi + Lai |
| Mock | Mock interview cơ bản | Candidate luyện phỏng vấn với AI | Lai + Khôi |
| Audit | Audit log cơ bản | Ghi log hành động quan trọng | Khôi |

---

## 3.1.1. Success Metrics — Tiêu chí đo lường thành công

| Metric | Mục tiêu | Đo bằng |
|---|---|---|
| Giảm thời gian sàng lọc ứng viên | Giảm ≥ 30% thời gian so với quy trình thủ công | Thời gian từ tạo job → quyết định |
| Tỷ lệ hoàn thành phỏng vấn | ≥ 90% buổi phỏng vấn scheduled được hoàn thành | `completed / scheduled` |
| Độ trễ transcript realtime | < 1.5 giây audio-to-text | P95 latency STT |
| AI scoring vs đánh giá recruiter | Tương đồng ≥ 75% (cùng kết luận pass/consider/reject) | So sánh sau khi recruiter xác nhận |
| System uptime phòng phỏng vấn | ≥ 99.5% trong giờ hành chính | Monitoring uptime |
| Mock interview completion rate | ≥ 70% phiên luyện được hoàn thành đến cuối | `completed / started` |
| AI report generation time | < 60 giây sau khi phỏng vấn kết thúc | P95 thời gian tạo report |

> Các con số trên là mục tiêu ban đầu — có thể điều chỉnh sau khi có dữ liệu thực tế từ pilot.

---

## 3.2. Ngoài MVP, chưa làm ngay

Các phần sau không bắt buộc ở MVP:

- ATS integration.
- Calendar integration sâu với Google/Microsoft.
- AI avatar 3D.
- Video recording — hệ thống không lưu video, chỉ lưu audio transcript.
- Coding test realtime.
- Multi-round hiring pipeline nâng cao.
- Billing/subscription nâng cao.
- Enterprise SSO.
- Multi-language phức tạp.
- Mobile app native.

AI IDE không được tự ý triển khai các phần này nếu task không yêu cầu.

---

## 3.3. Non-Scope — Hệ thống KHÔNG làm

Các tính năng này **không thuộc bất kỳ sprint nào** trừ khi có quyết định thay đổi scope rõ ràng:

| Tính năng | Lý do loại |
|---|---|
| Đánh giá qua nét mặt / cảm xúc qua camera | Chống bias — vi phạm nguyên tắc AI không đánh giá theo ngoại hình |
| Lưu trữ video recording | Không thuộc scope, tránh chi phí storage và rủi ro pháp lý PII |
| Coding test IDE tích hợp | Feature phức tạp — thuộc phase 2+ |
| Personality test / Psychometric test | Không trong scope AI Interview |
| AI tự động gửi email từ chối / chấp nhận | AI không được tự hành động thay recruiter |
| Tích hợp LinkedIn / ATS | Phase 2+ |
| Đánh giá giọng điệu, tốc độ nói | Tránh bias về chất giọng vùng miền |
| Public job board (candidate tự ứng tuyển) | Scope khác — đây là platform B2B cho recruiter |

---

## 4. Luồng nghiệp vụ chính

## 4.1. Luồng recruiter phỏng vấn thật

1. Recruiter đăng nhập.
2. Recruiter tạo workspace/company nếu chưa có.
3. Recruiter tạo job.
4. Recruiter nhập JD.
5. AI phân tích JD và đề xuất rubric/câu hỏi.
6. Recruiter thêm candidate.
7. Recruiter upload CV candidate.
8. AI parse/tóm tắt CV.
9. Recruiter tạo lịch phỏng vấn.
10. Hệ thống tạo interview room link.
11. Candidate nhận link.
12. Recruiter và candidate vào phòng.
13. Recruiter bắt đầu buổi phỏng vấn.
14. Hệ thống chạy transcript realtime.
15. AI gợi ý câu hỏi trong phòng.
16. Recruiter có thể ghi note.
17. Recruiter kết thúc buổi phỏng vấn.
18. AI tạo report.
19. Recruiter xem report và ra quyết định.

---

## 4.2. Luồng candidate vào phòng phỏng vấn

1. Candidate mở link mời.
2. Hệ thống kiểm tra token/link.
3. Candidate xác nhận họ tên/email.
4. Candidate kiểm tra mic/camera.
5. Candidate đồng ý điều khoản ghi âm/AI nếu hệ thống bật.
6. Candidate vào waiting room.
7. Khi recruiter bắt đầu, candidate vào phòng chính.
8. Candidate trả lời câu hỏi.
9. Candidate kết thúc khi recruiter kết thúc room.
10. Candidate thấy màn hình cảm ơn hoặc hướng dẫn tiếp theo.

---

## 4.3. Luồng candidate mock interview

1. Candidate đăng nhập.
2. Candidate chọn vị trí muốn luyện.
3. Candidate chọn level.
4. Candidate upload CV hoặc nhập kinh nghiệm.
5. AI tạo kịch bản câu hỏi.
6. Candidate bắt đầu mock interview.
7. AI hỏi từng câu.
8. Candidate trả lời bằng text/audio tùy phase.
9. AI lưu câu trả lời.
10. AI tạo feedback cuối buổi.
11. Candidate xem điểm, điểm mạnh, điểm yếu, gợi ý cải thiện.

---

## 5. Module sản phẩm

## 5.1. Auth & Account

Chức năng:

- Register.
- Login.
- Logout.
- Refresh token.
- Forgot password.
- Get current user.
- Role cơ bản: admin, recruiter, candidate.

Yêu cầu:

- Email là duy nhất.
- Password phải được hash.
- API cần bảo vệ bằng access token.
- Candidate có thể vào phòng bằng invite token mà chưa cần full account ở MVP nếu cần.

---

## 5.2. Company / Workspace

Chức năng:

- Tạo company.
- Cập nhật company.
- Mời thành viên recruiter.
- Phân quyền owner/member cơ bản.

Business rules:

- Mỗi recruiter phải thuộc ít nhất một company để tạo job.
- Candidate không được truy cập workspace recruiter.
- Dữ liệu job/candidate/interview luôn thuộc company.

---

## 5.3. Job Management

Chức năng:

- Tạo job.
- Cập nhật job.
- Đóng/mở job.
- Xem danh sách job.
- Xem chi tiết job.
- AI phân tích JD.
- AI sinh câu hỏi theo JD.

Job status:

| Status | Ý nghĩa |
|---|---|
| draft | Đang soạn |
| open | Đang tuyển |
| paused | Tạm dừng |
| closed | Đã đóng |

Yêu cầu UI:

- Có danh sách job.
- Có filter theo status, level, department.
- Có job detail gồm tab overview, candidates, questions, rubric, interviews, reports.

---

## 5.4. Candidate Management

Chức năng:

- Thêm candidate.
- Cập nhật candidate.
- Upload CV.
- AI parse CV.
- Gán candidate vào job.
- Xem candidate detail.
- Cập nhật trạng thái pipeline.

Candidate status:

| Status | Ý nghĩa |
|---|---|
| new | Mới thêm |
| screening | Đang sàng lọc |
| invited | Đã mời phỏng vấn |
| interviewing | Đang phỏng vấn |
| completed | Đã phỏng vấn |
| passed | Đạt |
| rejected | Loại |
| talent_pool | Lưu hồ sơ tiềm năng |

---

## 5.5. Interview Scheduling

Chức năng:

- Tạo buổi phỏng vấn.
- Chọn job.
- Chọn candidate.
- Chọn recruiter.
- Chọn thời gian.
- Chọn template/rubric.
- Tạo invite link.
- Gửi email mời nếu tích hợp email.

Interview status:

| Status | Ý nghĩa |
|---|---|
| scheduled | Đã lên lịch |
| waiting | Đang chờ |
| active | Đang diễn ra |
| paused | Tạm dừng |
| completed | Hoàn thành |
| cancelled | Đã hủy |
| expired | Hết hạn |

---

## 5.6. Interview Room

Chức năng:

- Join room.
- Leave room.
- Start interview.
- End interview.
- Chat.
- Transcript realtime.
- AI suggestion.
- Recruiter note.
- Candidate profile sidebar.
- Rubric scoring panel cho recruiter.

Business rules:

- Candidate chỉ thấy các thông tin được phép.
- Recruiter thấy AI suggestion và scoring.
- Candidate không thấy note nội bộ.
- Khi interview completed, room không cho start lại nếu không có quyền admin/recruiter.
- Nếu mất mạng, user được reconnect trong thời gian cho phép.

---

## 5.7. AI Assistant

Chức năng:

- Analyze JD.
- Analyze CV.
- Generate questions.
- Suggest follow-up.
- Score answer.
- Generate report.
- Mock interview feedback.

Business rules:

- Mọi AI score phải có evidence.
- AI phải có confidence.
- Nếu thiếu dữ liệu, trả về `insufficient_evidence`.
- Recommendation của AI không phải quyết định cuối cùng.

---

## 5.8. Report

Chức năng:

- Tạo report sau phỏng vấn.
- Xem report.
- Export PDF ở phase sau.
- Recruiter cập nhật final decision.

Report gồm:

- Summary.
- Final score.
- Score theo rubric.
- Strengths.
- Weaknesses.
- Risks.
- Evidence.
- AI recommendation.
- Recruiter final decision.
- Transcript link.

---

## 5.9. Mock Interview

Chức năng:

- Candidate chọn vị trí.
- Candidate chọn level.
- Candidate upload CV hoặc nhập profile.
- AI sinh câu hỏi.
- Candidate trả lời.
- AI feedback.
- Lưu lịch sử luyện tập.

Business rules:

- Mock interview không dùng làm quyết định tuyển dụng thật.
- Feedback chỉ hiển thị cho candidate.
- Candidate chỉ thấy lịch sử mock của chính mình.

---

## 6. Quy tắc phân quyền MVP

| Chức năng | Admin | Recruiter | Candidate |
|---|---:|---:|---:|
| Quản lý toàn hệ thống | Có | Không | Không |
| Tạo company | Có | Có | Không |
| Tạo job | Có | Có | Không |
| Thêm candidate | Có | Có | Không |
| Tạo lịch phỏng vấn | Có | Có | Không |
| Vào room phỏng vấn | Có | Có | Có nếu được mời |
| Start/end interview | Có | Có | Không |
| Xem AI scoring | Có | Có | Không |
| Xem recruiter note | Có | Có | Không |
| Xem report tuyển dụng | Có | Có | Không |
| Mock interview | Không bắt buộc | Không bắt buộc | Có |
| Xem feedback mock | Không | Không | Có |

---

## 7. Quy tắc dữ liệu và bảo mật

- Mọi dữ liệu job/candidate/interview phải gắn với `company_id`.
- API phải check quyền theo company.
- File CV, recording, transcript cần bảo vệ bằng signed URL hoặc permission check.
- Phải có audit log cho hành động quan trọng.
- Cần consent nếu ghi âm/ghi hình hoặc dùng AI realtime.
- Không dùng dữ liệu candidate để train model nếu chưa có consent.
- Candidate không được xem dữ liệu nội bộ recruiter.

---

## 8. Yêu cầu UX chính

Giao diện cần:

- Sáng, hiện đại, dễ nhìn.
- Phù hợp môi trường tuyển dụng/phỏng vấn.
- Không màu mè quá mức.
- Không gây rối trong interview room.
- Có loading state, empty state, error state.
- Có warning rõ khi AI chưa đủ dữ liệu.
- Các hành động nguy hiểm như end interview/cancel phải có confirm.

---

## 9. Acceptance Criteria tổng quát cho MVP

MVP đạt khi:

- Recruiter đăng nhập được.
- Recruiter tạo job được.
- Recruiter thêm candidate được.
- Recruiter tạo lịch phỏng vấn được.
- Candidate vào phòng bằng invite link được.
- Recruiter và candidate thấy trạng thái room realtime.
- Chat trong room hoạt động.
- Transcript realtime hoặc transcript text hoạt động ở mức cơ bản.
- AI gợi ý câu hỏi theo JD/CV/transcript.
- Recruiter kết thúc interview được.
- AI tạo report sau interview.
- Recruiter xem report được.
- Candidate chạy mock interview cơ bản được.
- AI feedback mock interview được.
- Phân quyền recruiter/candidate không bị lẫn.

---

## 10. Tech stack đã chốt

| Thành phần | Đã chọn |
|---|---|
| Frontend | Vue 3 (Composition API) + Pinia + TailwindCSS + VeeValidate |
| Backend | Golang (REST API) |
| Realtime | WebSocket (Golang) |
| AI Service | Python (LLM orchestration) |
| LLM | Gemini |
| STT | Whisper |
| Database | PostgreSQL |
| Cache/Queue | Redis |
| File Storage | S3-compatible |
| Auth | JWT + Refresh Token |

AI IDE không được tự ý đổi sang framework/ngôn ngữ khác.

---

## 11. Quy tắc cho AI IDE

Khi AI IDE code dự án này:

1. Luôn đọc `PRODUCT_REQUIREMENTS.md` trước khi làm task.
2. Nếu task liên quan API, đọc thêm `API_SPEC.md`.
3. Nếu task liên quan DB/model/migration, đọc thêm `DATABASE_DESIGN.md`.
4. Nếu task liên quan room/websocket, đọc thêm `REALTIME_EVENTS.md`.
5. Nếu task liên quan AI, prompt, score, report, đọc thêm `AI_PROMPT_AND_SCORING.md`.
6. Nếu task liên quan UI/giao diện, đọc thêm `GiaoDien.md`.
7. Không tự đổi tên endpoint, field, enum nếu chưa cập nhật tài liệu.
8. Không tự thêm tính năng ngoài MVP nếu không được giao.
9. Không tự đổi tech stack (framework, ngôn ngữ, database).
10. Mỗi task xong phải ghi rõ file đã sửa, API đã tạo, test đã chạy.
11. Nếu phát hiện spec mâu thuẫn, dừng lại và báo rõ mâu thuẫn.
12. Luôn giữ quyền riêng tư candidate và phân quyền recruiter/candidate.

