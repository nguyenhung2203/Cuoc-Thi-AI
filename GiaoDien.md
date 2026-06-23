Thiết kế giao diện web app SaaS cho dự án **AI Interview Platform** — nền tảng phỏng vấn trực tuyến có AI real-time hỗ trợ nhà tuyển dụng và ứng viên.

## Mục tiêu thiết kế

Tạo một giao diện hiện đại, sáng, chuyên nghiệp, dễ dùng, phù hợp với môi trường phỏng vấn, tuyển dụng và đánh giá ứng viên. Giao diện cần tạo cảm giác tin cậy, rõ ràng, công nghệ cao nhưng không quá lạnh lẽo. Ưu tiên UX đơn giản, sạch sẽ, nhiều khoảng trắng, dễ thao tác cho cả recruiter và candidate.

## Phong cách UI

* Modern SaaS dashboard
* Clean, minimal, professional
* Bright theme
* Friendly but serious
* Phù hợp môi trường HR, tuyển dụng, phỏng vấn online
* Không dùng màu quá chói hoặc quá game
* Không dùng giao diện tối làm chủ đạo
* Không dùng quá nhiều gradient màu mè
* Không dùng icon marketing quá lố

## Tông màu đề xuất

* Nền chính: trắng, xám rất nhạt hoặc xanh rất nhạt
* Màu chủ đạo: xanh dương hiện đại
* Màu phụ: cyan hoặc teal nhẹ
* Màu trạng thái:

  * Thành công: xanh lá nhẹ
  * Cảnh báo: vàng/cam nhẹ
  * Lỗi: đỏ nhẹ
  * Trung lập: xám xanh
* Tổng thể phải sáng, sạch, dễ nhìn, phù hợp cho làm việc lâu

Gợi ý palette:

* Primary: #2563EB
* Secondary: #06B6D4
* Background: #F8FAFC
* Surface: #FFFFFF
* Border: #E2E8F0
* Text primary: #0F172A
* Text secondary: #64748B
* Success: #16A34A
* Warning: #F59E0B
* Danger: #DC2626

## Font và style

* Dùng font sans-serif hiện đại, dễ đọc
* Heading rõ ràng, mạnh mẽ
* Body text dễ đọc
* Card bo góc mềm, shadow nhẹ
* Button rõ trạng thái hover/focus
* Input, select, textarea sạch và dễ thao tác
* Icon line-style đơn giản
* Layout responsive cho desktop trước, sau đó mobile/tablet

## Các màn hình cần thiết kế

### 1. Landing Page

Thiết kế trang giới thiệu sản phẩm cho nhà tuyển dụng và ứng viên.

Nội dung chính:

* Hero section:

  * Tiêu đề: “Phỏng vấn thông minh cùng AI real-time”
  * Mô tả: “Tạo phòng phỏng vấn online, AI ghi transcript, gợi ý câu hỏi, đánh giá ứng viên và tạo báo cáo sau phỏng vấn.”
  * CTA chính: “Bắt đầu dùng thử”
  * CTA phụ: “Xem demo”
* Hình minh họa dashboard/phòng phỏng vấn
* Section lợi ích:

  * Lọc ứng viên nhanh hơn
  * Đánh giá nhất quán hơn
  * Có transcript và báo cáo tự động
  * Ứng viên có thể luyện phỏng vấn thử
* Section tính năng:

  * AI Interview Room
  * AI Question Suggestion
  * AI Candidate Scoring
  * Mock Interview
  * Interview Report
* Section dành cho recruiter
* Section dành cho candidate
* Footer đơn giản

### 2. Login / Register Page

Thiết kế giao diện đăng nhập và đăng ký.

Yêu cầu:

* Bố cục 2 cột trên desktop
* Bên trái là form
* Bên phải là hình minh họa hoặc preview dashboard
* Form đơn giản, chuyên nghiệp
* Có chọn loại tài khoản:

  * Recruiter
  * Candidate
* Có đăng nhập Google
* Có quên mật khẩu
* Có thông báo lỗi rõ ràng

### 3. Recruiter Dashboard

Dashboard dành cho nhà tuyển dụng.

Thành phần:

* Sidebar bên trái
* Header có search, notification, avatar
* Cards thống kê:

  * Jobs đang mở
  * Ứng viên mới
  * Lịch phỏng vấn hôm nay
  * Báo cáo chờ xem
* Khu vực lịch phỏng vấn sắp tới
* Danh sách ứng viên mới nhất
* Card AI Insight:

  * “3 ứng viên có điểm phù hợp cao”
  * “2 buổi phỏng vấn cần review”
* Quick actions:

  * Tạo job
  * Thêm ứng viên
  * Tạo lịch phỏng vấn
  * Mở kho câu hỏi

### 4. Job Management Page

Trang quản lý vị trí tuyển dụng.

Thành phần:

* Danh sách job dạng table/card
* Bộ lọc theo trạng thái, phòng ban, level
* Nút tạo job
* Mỗi job hiển thị:

  * Tên vị trí
  * Phòng ban
  * Level
  * Số ứng viên
  * Số buổi phỏng vấn
  * Trạng thái
  * Điểm phù hợp trung bình
* Có panel AI gợi ý cải thiện JD

### 5. Job Detail Page

Trang chi tiết vị trí tuyển dụng.

Thành phần:

* Header job: title, status, department, level
* Tab:

  * Tổng quan
  * Ứng viên
  * Bộ câu hỏi
  * Rubric đánh giá
  * Lịch phỏng vấn
  * Báo cáo
* Khu vực JD
* AI summary của JD
* Rubric scoring
* Danh sách ứng viên theo pipeline:

  * New
  * Screening
  * Interview
  * Passed
  * Rejected

### 6. Candidate Management Page

Trang quản lý ứng viên.

Thành phần:

* Table danh sách ứng viên
* Filter theo job, status, score, source
* Search theo tên/email/kỹ năng
* Mỗi dòng hiển thị:

  * Avatar
  * Tên ứng viên
  * Email
  * Vị trí ứng tuyển
  * Trạng thái
  * AI fit score
  * Lịch phỏng vấn gần nhất
  * Action xem chi tiết

### 7. Candidate Detail Page

Trang chi tiết ứng viên.

Thành phần:

* Profile ứng viên
* CV preview
* AI CV summary
* Skills extracted
* Experience timeline
* Interview history
* Score cards
* Recruiter notes
* AI recommendation
* Action:

  * Mời phỏng vấn
  * Tạo mock test
  * Xem báo cáo
  * Cập nhật trạng thái

### 8. Interview Scheduling Page

Trang tạo lịch phỏng vấn.

Thành phần:

* Chọn job
* Chọn candidate
* Chọn recruiter
* Chọn thời gian
* Chọn interview template
* Chọn hình thức:

  * Online interview
  * Mock interview
* Preview email mời
* Nút gửi lời mời
* Giao diện đơn giản, giống wizard 3 bước

### 9. Interview Room — màn hình quan trọng nhất

Thiết kế phòng phỏng vấn real-time có 2 actor: Recruiter và Candidate, cùng AI assistant.

Layout đề xuất:

* Header:

  * Tên job
  * Tên candidate
  * Timer
  * Trạng thái recording/transcript
  * Nút kết thúc phỏng vấn
* Khu vực chính:

  * Video recruiter
  * Video candidate
  * Control bar: mic, camera, screen share, chat, end
* Sidebar bên phải:

  * Candidate profile mini
  * CV summary
  * AI suggested questions
  * Rubric score realtime
* Bottom panel:

  * Transcript realtime
  * Chat
  * Recruiter notes
* AI assistant panel:

  * Gợi ý câu hỏi tiếp theo
  * Cảnh báo điểm cần hỏi sâu
  * Tóm tắt nhanh câu trả lời
  * Đánh giá confidence

Yêu cầu cảm giác:

* Rõ ràng, tập trung
* Không làm recruiter bị rối
* AI panel phải hữu ích nhưng không chiếm quá nhiều không gian
* Candidate không thấy phần scoring nội bộ
* Có trạng thái “AI đang phân tích…”

### 10. Interview Report Page

Trang báo cáo sau phỏng vấn.

Thành phần:

* Tổng điểm ứng viên
* Recommendation:

  * Strong Hire
  * Hire
  * Consider
  * Reject
* Score theo tiêu chí:

  * Technical Knowledge
  * Problem Solving
  * Communication
  * Experience Relevance
  * Culture Fit
* Điểm mạnh
* Điểm yếu
* Rủi ro
* Evidence từ transcript
* Recruiter final decision
* Transcript đầy đủ
* Nút export PDF
* Nút chia sẻ nội bộ

### 11. Candidate Mock Interview Portal

Giao diện dành cho ứng viên luyện phỏng vấn thử.

Thành phần:

* Dashboard ứng viên
* Nút bắt đầu mock interview
* Chọn vị trí muốn luyện:

  * Frontend Developer
  * Backend Developer
  * Product Manager
  * Designer
  * Sales
  * Marketing
* Chọn level:

  * Intern
  * Junior
  * Middle
  * Senior
* Upload CV
* AI tạo câu hỏi
* Màn hình luyện phỏng vấn:

  * AI interviewer hỏi
  * Candidate trả lời bằng text/audio
  * Timer
  * Feedback sau từng câu hoặc cuối buổi
* Trang kết quả:

  * Điểm tổng
  * Điểm mạnh
  * Điểm cần cải thiện
  * Câu trả lời mẫu tốt hơn
  * Lộ trình luyện tập tiếp theo

### 12. Settings Page

Trang cài đặt.

Thành phần:

* Profile
* Company information
* Team members
* Roles & permissions
* AI settings
* Interview templates
* Data privacy
* Billing

## Navigation đề xuất

Sidebar recruiter:

* Dashboard
* Jobs
* Candidates
* Interviews
* Question Bank
* Reports
* AI Templates
* Settings

Sidebar candidate:

* Dashboard
* Mock Interview
* My Results
* My CV
* Practice History
* Settings

## Component cần có

Thiết kế bộ component đồng bộ:

* Button
* Input
* Select
* Textarea
* Search bar
* Filter chip
* Badge status
* Score badge
* Avatar
* Card
* Table
* Tabs
* Modal
* Drawer
* Toast
* Empty state
* Loading state
* AI thinking state
* Transcript item
* Video tile
* Timeline
* Progress bar
* Rating score
* Report section

## UX yêu cầu

* Dễ hiểu cho người mới dùng
* Tránh quá nhiều thông tin trên một màn hình
* AI insight phải có giải thích ngắn gọn
* Các hành động chính phải nổi bật
* Các trạng thái phải rõ:

  * Scheduled
  * Waiting
  * Active
  * Completed
  * Cancelled
  * Report Ready
* Có empty state đẹp khi chưa có dữ liệu
* Có loading skeleton
* Có confirmation khi kết thúc phỏng vấn
* Có warning khi AI chưa đủ dữ liệu để chấm điểm

## Copywriting mẫu

Sử dụng tiếng Việt có dấu, ngắn gọn, chuyên nghiệp.

Ví dụ:

* “Tạo buổi phỏng vấn”
* “Mời ứng viên”
* “AI đang phân tích câu trả lời…”
* “Gợi ý câu hỏi tiếp theo”
* “Báo cáo đã sẵn sàng”
* “Ứng viên này phù hợp với vị trí ở mức cao”
* “Chưa đủ dữ liệu để đánh giá tiêu chí này”
* “Bắt đầu luyện phỏng vấn thử”
* “Cải thiện câu trả lời của bạn”

## Yêu cầu cuối cùng

Hãy tạo giao diện web app hoàn chỉnh, hiện đại, sáng, chuyên nghiệp cho sản phẩm AI Interview Platform. Ưu tiên thiết kế các màn hình quan trọng nhất trước:

1. Recruiter Dashboard
2. Interview Room
3. Interview Report
4. Candidate Mock Interview
5. Job Detail
6. Candidate Detail

Thiết kế phải có tính hệ thống, thống nhất component, màu sắc, layout và trải nghiệm người dùng. Giao diện cần nhìn giống một sản phẩm SaaS thật, có thể dùng để demo với team và khách hàng.