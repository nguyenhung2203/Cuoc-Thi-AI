# Kịch bản Demo dự án AI Interview Platform
**Dành cho luồng: Tuyển dụng (Recruiter) & Admin**

Tài liệu này hướng dẫn bạn từng bước để demo dự án "AI Interview Platform" cho thầy cô, tập trung vào các tính năng nổi bật của AI trong quá trình tuyển dụng.

---

## 🚀 Trước khi Demo (Chuẩn bị)

1. **Khởi động hệ thống:** Đảm bảo tất cả các service (Frontend, Backend API, Backend Realtime, AI Service) đều đang chạy. Có thể dùng lệnh `docker-compose up -d` hoặc chạy tay từng service.
2. **Dữ liệu mẫu (Seed Data):** Nên chuẩn bị sẵn:
   - 1 tài khoản Recruiter.
   - 1 Job Description (JD) mẫu (ví dụ: Tuyển Frontend Developer VueJS).
   - 1-2 CV mẫu (file PDF) để upload.
3. **Trình duyệt:** 
   - Mở sẵn 1 cửa sổ chính (đăng nhập tài khoản Recruiter).
   - Mở sẵn 1 cửa sổ ẩn danh (Incognito) để giả lập làm Ứng viên (Candidate) khi vào phòng phỏng vấn.

---

## 🎬 Kịch bản Demo chi tiết

### Bước 1: Tổng quan (Dashboard)
**Mục tiêu:** Cho thầy cô thấy cái nhìn toàn cảnh của hệ thống quản lý tuyển dụng.
- **Hành động:** Đăng nhập vào tài khoản Recruiter. Mở trang Dashboard.
- **Thuyết minh:** 
  > "Chào thầy, đây là hệ thống AI Interview Platform dành cho Nhà tuyển dụng (Recruiter). Ở màn hình Dashboard, hệ thống cung cấp các chỉ số tổng quan như số lượng Job đang mở, số ứng viên mới, lịch phỏng vấn hôm nay. Điểm đặc biệt là mục **AI Insights**, nơi AI tự động đưa ra các gợi ý nhắc nhở recruiter như 'Ứng viên X có mức độ phù hợp cao' hoặc 'Cần review báo cáo phỏng vấn Y'."

### Bước 2: Quản lý Công việc & Phân tích JD bằng AI
**Mục tiêu:** Demo khả năng AI đọc hiểu và xử lý Job Description.
- **Hành động:** Vào mục **Việc làm (Jobs)** -> Chọn **Tạo Job mới**. Nhập thông tin và paste một đoạn JD mẫu.
- **Thuyết minh:**
  > "Khi nhà tuyển dụng tạo một vị trí mới, AI của hệ thống sẽ tự động phân tích Job Description (JD). Nó không chỉ tóm tắt yêu cầu mà còn **tự động sinh ra bộ câu hỏi phỏng vấn dự kiến** và **Rubric đánh giá (tiêu chí chấm điểm)** phù hợp nhất cho vị trí này. Điều này giúp tiết kiệm rất nhiều thời gian chuẩn bị cho Recruiter."

### Bước 3: Quản lý Ứng viên & AI Parse CV
**Mục tiêu:** Demo khả năng AI đọc hiểu CV của ứng viên.
- **Hành động:** Vào mục **Ứng viên (Candidates)** -> Chọn **Thêm ứng viên** -> Upload một file CV PDF.
- **Thuyết minh:**
  > "Sau khi có Job, chúng ta sẽ thêm ứng viên. Khi upload CV, AI Service (sử dụng LLM) sẽ tự động **Parse (trích xuất) thông tin** từ CV, đối chiếu với JD của vị trí ứng tuyển để đưa ra **Điểm phù hợp (Fit Score)**. AI cũng sẽ chỉ ra điểm mạnh, điểm yếu của ứng viên ngay trên hệ thống trước khi phỏng vấn."

### Bước 4: Lên lịch phỏng vấn
**Mục tiêu:** Cho thấy luồng kết nối giữa Recruiter và Candidate.
- **Hành động:** Chọn một ứng viên -> Tạo lịch phỏng vấn. Sinh ra link tham gia phòng phỏng vấn.
- **Thuyết minh:**
  > "Hệ thống cho phép lên lịch phỏng vấn và sinh ra link phòng họp trực tuyến (giống Google Meet nhưng chuyên dụng cho tuyển dụng). Em sẽ copy link này để giả lập ứng viên tham gia."

### Bước 5: Phòng Phỏng vấn Real-time (Tính năng Core/Ăn tiền nhất)
**Mục tiêu:** Demo khả năng Real-time AI Assistant trong phòng phỏng vấn.
- **Hành động:**
  1. Ở cửa sổ chính: Recruiter bấm **Join Room**.
  2. Ở cửa sổ ẩn danh: Paste link, nhập tên Candidate và **Join Room**.
  3. Bật mic và nói thử 1-2 câu tiếng Việt.
- **Thuyết minh:**
  > "Đây là không gian phòng phỏng vấn trực tuyến. Hệ thống của bọn em tích hợp WebSocket và LiveKit để stream audio/video. 
  > 
  > Điểm đột phá ở đây là bảng điều khiển bên phải dành riêng cho Recruiter. 
  > 1. **Live Transcript:** Mọi lời nói của ứng viên và nhà tuyển dụng đều được AI (Whisper) chuyển thành văn bản theo thời gian thực.
  > 2. **AI Suggestion (Gợi ý câu hỏi):** Dựa vào câu trả lời vừa rồi của ứng viên, AI sẽ phân tích xem câu trả lời có đủ ý chưa, có bằng chứng chưa. Nếu chưa, nó sẽ **gợi ý ngay các câu hỏi follow-up (hỏi xoáy)** trên màn hình để Recruiter bấm hỏi tiếp.
  > 3. **AI Scoring (Chấm điểm trực tiếp):** AI sẽ đối chiếu câu trả lời với Rubric để tự động đề xuất điểm số cho từng tiêu chí."

### Bước 6: Báo cáo sau phỏng vấn (Interview Report)
**Mục tiêu:** Khả năng tổng hợp và tự động hóa báo cáo.
- **Hành động:** Bấm **Kết thúc phỏng vấn (End Interview)**. Mở trang Báo cáo phỏng vấn của ứng viên đó.
- **Thuyết minh:**
  > "Ngay sau khi kết thúc buổi phỏng vấn, thay vì Recruiter phải ngồi nhớ lại và gõ báo cáo, AI sẽ dựa vào toàn bộ Transcript để **tạo ra một Báo cáo phỏng vấn hoàn chỉnh (Interview Report)** trong vòng chưa đầy 60 giây.
  > 
  > Báo cáo này bao gồm: Điểm số tổng quan, Đề xuất (Pass/Fail), Điểm mạnh, Rủi ro, và đặc biệt là **trích dẫn bằng chứng (evidence)** rõ ràng từ lời nói của ứng viên trong transcript để chứng minh tại sao AI lại chấm điểm như vậy. Cuối cùng, quyền quyết định vẫn thuộc về Recruiter, AI chỉ đóng vai trò trợ lý đắc lực."

---

## 💡 Mẹo để buổi Demo ấn tượng hơn
- **Nhấn mạnh vào "Sự phân quyền":** Nói rõ với thầy rằng Candidate sẽ KHÔNG thể nhìn thấy điểm số AI hay transcript nội bộ, màn hình của Candidate rất sạch sẽ và thân thiện, chỉ có Recruiter mới thấy bảng AI Insights.
- **Nói về giới hạn của AI:** Thể hiện sự hiểu biết sâu bằng cách nói: *"Hệ thống của em thiết kế AI chỉ là Trợ lý (Assistant), không tự động loại ứng viên. Mọi đánh giá của AI đều yêu cầu trích dẫn bằng chứng từ Transcript để chống Halucination (ảo giác) và chống bias (thiên kiến)."*
- **Sử dụng Micro tốt:** Phần Realtime Transcript rất nhạy cảm với chất lượng âm thanh, hãy đảm bảo mic của bạn thu âm rõ ràng khi nói thử trong phòng phỏng vấn.

Chúc bạn có một buổi Demo thật thành công!
