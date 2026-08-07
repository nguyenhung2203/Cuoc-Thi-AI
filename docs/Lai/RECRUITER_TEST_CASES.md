# BỘ TEST CASE CHỨC NĂNG ACTOR NHÀ TUYỂN DỤNG (RECRUITER)

## 1. Thông tin tài liệu

| Thuộc tính | Giá trị |
|---|---|
| Dự án | ViệcLàm AI – Nền tảng tuyển dụng và phỏng vấn AI |
| Actor chính | Recruiter/Nhà tuyển dụng |
| Loại kiểm thử | Functional, API, RBAC, Company Scope, UI/UX, Realtime, AI, Security, Performance |
| Frontend | Vue 3/Vite |
| Backend | Go REST API, Realtime Gateway, AI Service |
| API prefix | `/api/v1` |

> Tài liệu được xây dựng theo route, service frontend và handler backend hiện có. Mã trạng thái HTTP cụ thể có thể là `400/401/403/404/409/422` tùy quy ước lỗi chung; điều kiện bắt buộc là không làm lộ dữ liệu và không thay đổi dữ liệu khi request thất bại.

## 2. Phạm vi kiểm thử

| Phân hệ | Route frontend | API/chức năng chính |
|---|---|---|
| Đăng ký, xác thực | `/register`, `/login` | Register, OTP, login, refresh, logout, hồ sơ tài khoản |
| Dashboard | `/dashboard` | Số liệu tuyển dụng, lịch sắp tới, trạng thái pipeline |
| Công ty | `/settings` | Xem/cập nhật thông tin công ty và tài khoản |
| Việc làm | `/jobs`, `/jobs/:id` | CRUD Job, trạng thái, phân tích JD, tạo câu hỏi AI |
| Ứng viên | `/candidates`, `/candidates/:id` | CRUD, CV, parse CV, tìm kiếm/lọc |
| Pipeline | Chi tiết Job/Ứng viên | Gán, bỏ gán, cập nhật pipeline |
| Lịch phỏng vấn | `/interviews`, `/interviews/new`, `/interviews/:id` | Tạo, xem, ghi chú, nhắc lịch, bắt đầu, kết thúc, hủy |
| Phòng phỏng vấn | `/recruiter-room/:interviewId` | LiveKit, realtime, chat, transcript, reconnect |
| Báo cáo | `/interviews/:id/report` | Báo cáo AI, retry, quyết định Recruiter |
| Ngân hàng câu hỏi | `/question-bank` | CRUD, lọc và tạo câu hỏi AI |
| Rubric | `/rubrics` | CRUD tiêu chí chấm điểm |
| Template | `/templates` | CRUD mẫu phỏng vấn |
| Thông báo | Header/notification UI | Danh sách, đọc một, đọc tất cả |
| Audit log | API company scope | Tra cứu hoạt động trong công ty |

## 3. Dữ liệu kiểm thử

### 3.1 Tài khoản và công ty

| Mã | Vai trò/quyền | Công ty | Trạng thái | Mục đích |
|---|---|---|---|---|
| REC-OWNER-A | Recruiter/Owner | COM-A | Active | Toàn quyền trong COM-A |
| REC-MEMBER-A | Recruiter/Member | COM-A | Active | Kiểm tra permission chi tiết |
| REC-READ-A | Recruiter chỉ đọc | COM-A | Active | Kiểm tra quyền read-only |
| REC-B | Recruiter | COM-B | Active | Kiểm tra cách ly công ty |
| REC-PENDING | Recruiter | Chưa duyệt | Pending | Kiểm tra duyệt tài khoản |
| REC-BLOCKED | Recruiter | COM-A | Blocked | Kiểm tra khóa tài khoản |
| CAN-A | Candidate | — | Active | Kiểm tra RBAC |
| ADM-01 | Admin | — | Active | Kiểm tra phân quyền chéo |

### 3.2 Dữ liệu nghiệp vụ

- `JOB-A-OPEN`: Job đang mở của COM-A.
- `JOB-A-CLOSED`: Job đã đóng của COM-A.
- `JOB-B-OPEN`: Job của COM-B.
- `CAN-A-CV`: Ứng viên COM-A có CV hợp lệ và dữ liệu parse.
- `CAN-A-NOCV`: Ứng viên COM-A chưa có CV.
- `CAN-A-STALE`: Ứng viên trỏ tới metadata/file CV cũ không còn tồn tại.
- `CAN-B-CV`: Ứng viên COM-B có CV.
- `IV-A-SCHEDULED`: Lịch sắp tới của COM-A.
- `IV-A-COMPLETED`: Lịch đã kết thúc và có báo cáo.
- `IV-B-SCHEDULED`: Lịch của COM-B.
- `valid-cv.pdf`: PDF hợp lệ, dưới giới hạn.
- `fake.pdf`: Đuôi PDF nhưng magic bytes không hợp lệ.
- `large-cv.pdf`: Vượt giới hạn upload.

---

## 4. Đăng ký, đăng nhập và quản lý phiên

| ID | Test case | Bước/dữ liệu chính | Kết quả mong đợi | Ưu tiên |
|---|---|---|---|---|
| TC-REC-AUTH-001 | Đăng ký Recruiter hợp lệ | Chọn Recruiter, nhập dữ liệu hợp lệ, xác thực OTP nếu được yêu cầu | Tạo tài khoản đúng role; không lưu mật khẩu plaintext; chuyển đúng bước tiếp theo | Critical |
| TC-REC-AUTH-002 | Email sai định dạng | Nhập `recruiter-example` | FE và BE từ chối | High |
| TC-REC-AUTH-003 | Email đã tồn tại | Dùng email của REC-OWNER-A | Không tạo tài khoản trùng; trả lỗi nghiệp vụ phù hợp | High |
| TC-REC-AUTH-004 | OTP đúng | Nhập OTP còn hiệu lực | Xác thực thành công, OTP không tái sử dụng được | High |
| TC-REC-AUTH-005 | OTP sai/hết hạn | Nhập OTP sai hoặc cũ | Từ chối, không kích hoạt tài khoản | High |
| TC-REC-AUTH-006 | Login thành công | REC-OWNER-A nhập đúng thông tin | Lưu access/refresh token và role; chuyển `/dashboard` | Critical |
| TC-REC-AUTH-007 | Sai mật khẩu | Mật khẩu không đúng | Không cấp token; thông báo không tiết lộ quá mức | Critical |
| TC-REC-AUTH-008 | Tài khoản pending đăng nhập | REC-PENDING | Không vào dashboard; thông báo trạng thái duyệt rõ ràng | Critical |
| TC-REC-AUTH-009 | Tài khoản bị khóa | REC-BLOCKED | Không cấp phiên mới | Critical |
| TC-REC-AUTH-010 | Refresh token hợp lệ | Access token hết hạn, refresh token còn hạn | Cấp access token mới đúng user/role | Critical |
| TC-REC-AUTH-011 | Refresh token sai/đã thu hồi | Sửa token hoặc logout trước đó | HTTP 401; xóa phiên frontend | Critical |
| TC-REC-AUTH-012 | Logout | Bấm đăng xuất | Thu hồi phiên/token phù hợp; quay về login | High |
| TC-REC-AUTH-013 | Logout tất cả thiết bị | Gọi `/auth/logout-all` | Các refresh token của user không dùng lại được | High |
| TC-REC-AUTH-014 | Double click Login | Bấm liên tục khi loading | Chỉ tạo một request logic; nút bị disable | Medium |
| TC-REC-AUTH-015 | Mất mạng khi Login | Ngắt API | Dừng loading, báo lỗi, không lưu token rác | High |

## 5. RBAC, permission và company scope

| ID | Test case | Kết quả mong đợi | Ưu tiên |
|---|---|---|---|
| TC-REC-RBAC-001 | Chưa đăng nhập truy cập `/dashboard` | Chuyển tới login; không tải dữ liệu | Critical |
| TC-REC-RBAC-002 | Candidate truy cập `/jobs` hoặc `/candidates` | Chuyển `/403`; không gọi/không nhận API Recruiter | Critical |
| TC-REC-RBAC-003 | Recruiter truy cập route Admin | Chuyển `/403`; API Admin trả 403 | Critical |
| TC-REC-RBAC-004 | API company không có token | HTTP 401 | Critical |
| TC-REC-RBAC-005 | Token bị chỉnh sửa/hết hạn | HTTP 401; frontend kết thúc phiên | Critical |
| TC-REC-RBAC-006 | REC-B đọc dữ liệu COM-A bằng sửa URL | HTTP 403 hoặc 404; không trả dữ liệu COM-A | Critical |
| TC-REC-RBAC-007 | REC-B sửa/xóa Job COM-A | Bị từ chối; Job không thay đổi | Critical |
| TC-REC-RBAC-008 | REC-B xem CV ứng viên COM-A | Không cấp signed URL; không trả metadata nhạy cảm | Critical |
| TC-REC-RBAC-009 | REC-B truy cập interview/report COM-A | Bị từ chối | Critical |
| TC-REC-RBAC-010 | Thành viên thiếu `job:create` tạo Job | HTTP 403; không tạo bản ghi | Critical |
| TC-REC-RBAC-011 | Thành viên read-only cập nhật interview | HTTP 403 | Critical |
| TC-REC-RBAC-012 | ID UUID sai định dạng | Trả lỗi an toàn; không panic/SQL error lộ ra ngoài | High |
| TC-REC-RBAC-013 | Thay `company_id` trong query signed URL | Chỉ công ty sở hữu/được liên kết Candidate mới truy cập được | Critical |

## 6. Dashboard Recruiter

| ID | Test case | Kết quả mong đợi | Ưu tiên |
|---|---|---|---|
| TC-REC-DASH-001 | Tải Dashboard bình thường | Hiển thị số liệu, lịch gần nhất và dữ liệu COM-A | High |
| TC-REC-DASH-002 | Công ty chưa có dữ liệu | Hiển thị trạng thái rỗng và CTA phù hợp, không lỗi | Medium |
| TC-REC-DASH-003 | Đếm ứng viên theo pipeline | Tổng và từng trạng thái khớp DB, không đếm ứng viên công ty khác | High |
| TC-REC-DASH-004 | Lịch phỏng vấn sắp tới | Chỉ hiển thị lịch hợp lệ, thời gian đúng timezone | High |
| TC-REC-DASH-005 | API từng widget lỗi | Trang không trắng; widget báo lỗi/retry phù hợp | Medium |
| TC-REC-DASH-006 | Refresh nhanh nhiều lần | Không nhân đôi dữ liệu, không hiển thị response cũ đè response mới | Medium |

## 7. Công ty và cài đặt tài khoản

| ID | Test case | Kết quả mong đợi | Ưu tiên |
|---|---|---|---|
| TC-REC-COM-001 | Xem công ty hiện tại | `GET /companies/{id}` trả đúng COM-A | High |
| TC-REC-COM-002 | Cập nhật thông tin hợp lệ | Lưu tên, website, ngành, quy mô đúng; UI đồng bộ | High |
| TC-REC-COM-003 | Thiếu trường bắt buộc | FE/BE từ chối, dữ liệu cũ giữ nguyên | High |
| TC-REC-COM-004 | Website sai định dạng | Báo validation phù hợp | Medium |
| TC-REC-COM-005 | XSS trong tên/mô tả | Hiển thị dạng text, không thực thi script | Critical |
| TC-REC-COM-006 | REC-B cập nhật COM-A | Bị từ chối | Critical |
| TC-REC-COM-007 | Cập nhật hồ sơ cá nhân | `/auth/me` phản ánh dữ liệu mới | High |
| TC-REC-COM-008 | Đổi mật khẩu đúng mật khẩu cũ | Thành công; phiên xử lý theo chính sách | High |
| TC-REC-COM-009 | Đổi mật khẩu với mật khẩu cũ sai | Bị từ chối; mật khẩu không đổi | High |
| TC-REC-COM-010 | Xóa tài khoản | Có xác nhận; dữ liệu/quyền truy cập xử lý theo chính sách | Critical |

## 8. Quản lý việc làm

| ID | Test case | Bước/dữ liệu chính | Kết quả mong đợi | Ưu tiên |
|---|---|---|---|---|
| TC-REC-JOB-001 | Danh sách Job | `GET /companies/COM-A/jobs` | Chỉ trả Job COM-A, phân trang đúng | High |
| TC-REC-JOB-002 | Tìm kiếm/lọc | keyword, status, page, page_size | Kết quả đúng điều kiện; giữ filter khi phân trang | High |
| TC-REC-JOB-003 | Tạo Job hợp lệ | Nhập title, mô tả, yêu cầu, trạng thái | Tạo một Job thuộc COM-A | Critical |
| TC-REC-JOB-004 | Thiếu title | Gửi title rỗng | Từ chối, không tạo Job | High |
| TC-REC-JOB-005 | Dữ liệu biên | Unicode, nội dung dài, salary/date biên | Validation đúng, không lỗi encoding | Medium |
| TC-REC-JOB-006 | Xem chi tiết | Mở `/jobs/:id` | Dữ liệu đúng Job; trạng thái loading/not-found rõ ràng | High |
| TC-REC-JOB-007 | Cập nhật Job | Đổi mô tả/trạng thái | DB và UI cập nhật đồng nhất | Critical |
| TC-REC-JOB-008 | Đóng/mở Job | Chuyển trạng thái hợp lệ | Job public chỉ xuất hiện khi trạng thái cho phép | High |
| TC-REC-JOB-009 | Xóa Job chưa có ràng buộc | Xác nhận xóa | Xóa/soft-delete thành công; biến mất khỏi danh sách | High |
| TC-REC-JOB-010 | Xóa Job có ứng viên/lịch | Thử xóa | Hệ thống chặn hoặc xử lý quan hệ đúng, không tạo dữ liệu mồ côi | Critical |
| TC-REC-JOB-011 | Hủy modal xóa | Bấm Hủy | Không gọi DELETE | Medium |
| TC-REC-JOB-012 | Phân tích JD bằng AI | `POST .../analyze` | Trả/lưu kết quả phân tích đúng Job | High |
| TC-REC-JOB-013 | Force refresh phân tích JD | `force_refresh=true` | Tạo kết quả mới; không dùng cache cũ sai mục đích | Medium |
| TC-REC-JOB-014 | AI service lỗi/timeout | Ngắt AI service | Job vẫn tồn tại; UI báo lỗi và cho retry | High |
| TC-REC-JOB-015 | Tạo câu hỏi AI từ Job | Gọi endpoint generate questions | Câu hỏi gắn đúng Job/công ty, không trùng ngoài ý muốn | High |
| TC-REC-JOB-016 | Double submit Create/Update | Bấm nhanh nhiều lần | Không tạo bản ghi trùng; loading đúng | Medium |

## 9. Quản lý ứng viên, CV và AI parse

| ID | Test case | Kết quả mong đợi | Ưu tiên |
|---|---|---|---|
| TC-REC-CAN-001 | Danh sách ứng viên | Chỉ trả Candidate COM-A; phân trang đúng | Critical |
| TC-REC-CAN-002 | Tìm kiếm/lọc | Tìm theo tên/email/status | Kết quả đúng, không phân biệt encoding sai | High |
| TC-REC-CAN-003 | Tạo ứng viên thủ công | Dữ liệu hợp lệ | Tạo Candidate thuộc COM-A | High |
| TC-REC-CAN-004 | Email sai/trường bắt buộc rỗng | FE và BE từ chối | High |
| TC-REC-CAN-005 | Email trùng trong phạm vi nghiệp vụ | Không tạo trùng hoặc báo xung đột rõ ràng | High |
| TC-REC-CAN-006 | Xem chi tiết ứng viên | Hiển thị hồ sơ, CV, pipeline và AI summary đúng | Critical |
| TC-REC-CAN-007 | Cập nhật ứng viên | Dữ liệu lưu đúng; UI đồng bộ | High |
| TC-REC-CAN-008 | Xóa ứng viên | Có xác nhận; soft-delete/quan hệ xử lý đúng | Critical |
| TC-REC-CAN-009 | Upload CV hợp lệ | Metadata, checksum, storage key và `cv_file_id` được tạo đúng | Critical |
| TC-REC-CAN-010 | File vượt giới hạn | Từ chối trước khi lưu; không để file rác | High |
| TC-REC-CAN-011 | MIME spoofing | `fake.pdf` bị từ chối bằng magic-byte validation | Critical |
| TC-REC-CAN-012 | Tên file Unicode/ký tự đặc biệt | Lưu original name an toàn; storage key không dùng trực tiếp tên nguy hiểm | High |
| TC-REC-CAN-013 | Xem CV qua signed URL | URL hợp lệ, có thời hạn và mở đúng file | Critical |
| TC-REC-CAN-014 | Candidate portal lưu CV, Recruiter xem | Sau khi Candidate bấm lưu, Recruiter COM-A xem được CV gắn vào hồ sơ | Critical |
| TC-REC-CAN-015 | Parse CV hợp lệ | `POST .../parse-cv` cập nhật parsed data/summary | Critical |
| TC-REC-CAN-016 | Candidate chưa có CV | Nút xem/parse bị vô hiệu hoặc API báo lỗi rõ ràng | High |
| TC-REC-CAN-017 | Metadata/file CV cũ bị thiếu | UI báo “CV không khả dụng”, yêu cầu upload lại; không gọi parse vô ích | Critical |
| TC-REC-CAN-018 | Signed URL hết hạn | URL cũ bị từ chối; lấy URL mới hoạt động | High |
| TC-REC-CAN-019 | Recruiter COM-B đoán file ID COM-A | Không nhận signed URL hoặc nội dung file | Critical |
| TC-REC-CAN-020 | File trùng checksum | Dedup không làm sai ownership/company access | Critical |
| TC-REC-CAN-021 | AI parse lỗi | Không mất CV; hiện lỗi và cho retry | High |
| TC-REC-CAN-022 | XSS trong parsed CV | Nội dung kỹ năng/summary không thực thi HTML/script | Critical |

## 10. Gán Job và pipeline ứng viên

| ID | Test case | Kết quả mong đợi | Ưu tiên |
|---|---|---|---|
| TC-REC-PIPE-001 | Danh sách ứng viên theo Job | `GET /jobs/{job}/candidates` chỉ trả đúng company/job | High |
| TC-REC-PIPE-002 | Gán Candidate vào Job | Tạo quan hệ một lần, pipeline mặc định hợp lệ | Critical |
| TC-REC-PIPE-003 | Gán trùng | Không tạo hai quan hệ trùng | High |
| TC-REC-PIPE-004 | Gán Candidate COM-B vào Job COM-A | Bị từ chối | Critical |
| TC-REC-PIPE-005 | Cập nhật pipeline hợp lệ | Trạng thái đổi đúng và phản ánh trên dashboard | Critical |
| TC-REC-PIPE-006 | Pipeline status không hợp lệ | HTTP 400/422; trạng thái cũ giữ nguyên | High |
| TC-REC-PIPE-007 | Hai Recruiter cập nhật đồng thời | Không mất cập nhật âm thầm; kết quả nhất quán theo chính sách | High |
| TC-REC-PIPE-008 | Bỏ gán Candidate | Xóa đúng quan hệ, không xóa hồ sơ Candidate | High |
| TC-REC-PIPE-009 | Fit score sau khi CV thay đổi | Điểm match được tính lại cho các đơn liên quan | High |

## 11. Lịch phỏng vấn

| ID | Test case | Kết quả mong đợi | Ưu tiên |
|---|---|---|---|
| TC-REC-IV-001 | Danh sách lịch | Filter status/date và pagination đúng; chỉ COM-A | Critical |
| TC-REC-IV-002 | Tạo lịch hợp lệ | Job, Candidate, Recruiter, thời gian, duration, mode hợp lệ | Tạo lịch và invite token; xuất hiện ở hai phía | Critical |
| TC-REC-IV-003 | Thiếu Candidate/Job/thời gian | FE/BE từ chối | High |
| TC-REC-IV-004 | Thời gian trong quá khứ | Từ chối theo quy tắc nghiệp vụ | High |
| TC-REC-IV-005 | Duration bằng 0/âm/quá lớn | Validation phù hợp | High |
| TC-REC-IV-006 | Job và Candidate khác công ty | Bị từ chối | Critical |
| TC-REC-IV-007 | Trùng lịch Recruiter | Cảnh báo/chặn theo chính sách; không tạo âm thầm | High |
| TC-REC-IV-008 | Xem chi tiết | Trả đúng lịch, người tham gia, trạng thái | High |
| TC-REC-IV-009 | Lưu ghi chú nội bộ | `PUT .../notes`; Candidate không nhìn thấy ghi chú nội bộ | Critical |
| TC-REC-IV-010 | Gửi nhắc lịch | Gửi một reminder đúng ứng viên/lịch; UI báo kết quả | High |
| TC-REC-IV-011 | Bắt đầu lịch Scheduled | Trạng thái chuyển đúng; kiểm tra consent recording/AI | Critical |
| TC-REC-IV-012 | Bắt đầu lịch đã hủy/hoàn thành | Bị từ chối | Critical |
| TC-REC-IV-013 | Kết thúc lịch đang diễn ra | Chuyển Completed và khởi tạo report nếu yêu cầu | Critical |
| TC-REC-IV-014 | Kết thúc hai lần | Idempotent hoặc từ chối rõ ràng, không tạo report trùng | High |
| TC-REC-IV-015 | Mở modal Hủy lịch | Hiển thị đúng lịch, ô lý do và nút xác nhận/hủy | High |
| TC-REC-IV-016 | Xác nhận Hủy lịch | `POST .../cancel`, trạng thái `cancelled`, lưu lý do | Critical |
| TC-REC-IV-017 | Hủy modal | Không gọi API, trạng thái không đổi | Medium |
| TC-REC-IV-018 | API hủy lỗi | Modal/notification báo lỗi; trạng thái UI không giả thành cancelled | High |
| TC-REC-IV-019 | Hủy lịch Completed/Cancelled | Bị từ chối hoặc action không xuất hiện | High |
| TC-REC-IV-020 | Candidate xem sau khi Recruiter hủy | `/portal/interviews` phản ánh trạng thái cancelled | Critical |
| TC-REC-IV-021 | Timezone | Giờ tạo, danh sách và phòng phỏng vấn thống nhất | High |

## 12. Phòng phỏng vấn, LiveKit và realtime

| ID | Test case | Kết quả mong đợi | Ưu tiên |
|---|---|---|---|
| TC-REC-ROOM-001 | Lấy thông tin room | Chỉ thành viên có quyền của COM-A nhận room | Critical |
| TC-REC-ROOM-002 | Lấy recruiter room token | Token có identity/room/quyền đúng và thời hạn hữu hạn | Critical |
| TC-REC-ROOM-003 | Token cho interview COM-B | REC-A bị từ chối | Critical |
| TC-REC-ROOM-004 | Vào phòng thành công | Kết nối audio/video và hiển thị participants | Critical |
| TC-REC-ROOM-005 | Từ chối camera/microphone | UI báo rõ, vẫn cho retry/chọn thiết bị | High |
| TC-REC-ROOM-006 | Mất kết nối ngắn | Tự reconnect, không nhân đôi participant/message | High |
| TC-REC-ROOM-007 | Token hết hạn giữa phiên | Refresh/reconnect an toàn hoặc yêu cầu vào lại rõ ràng | High |
| TC-REC-ROOM-008 | Candidate dùng invite token hợp lệ | Vào đúng interview, không truy cập interview khác | Critical |
| TC-REC-ROOM-009 | Invite token sai/hết hạn | Bị từ chối, không lộ chi tiết lịch | Critical |
| TC-REC-ROOM-010 | Gửi chat | Message đúng sender, thứ tự và timestamp | High |
| TC-REC-ROOM-011 | XSS trong chat | Nội dung được escape | Critical |
| TC-REC-ROOM-012 | Push transcript | Chỉ quyền `interview:update`; lưu đúng speaker/timestamp | Critical |
| TC-REC-ROOM-013 | Xem transcript | Chỉ quyền `interview:read`; đúng thứ tự | High |
| TC-REC-ROOM-014 | Sửa transcript | Lưu nội dung mới và audit phù hợp | High |
| TC-REC-ROOM-015 | Hai client gửi transcript đồng thời | Không mất/nhân đôi đoạn ngoài ý muốn | High |
| TC-REC-ROOM-016 | Kết thúc phòng | Track được đóng, trạng thái và report trigger đúng | Critical |
| TC-REC-ROOM-017 | Realtime gateway không khả dụng | REST chính vẫn an toàn; UI báo mất realtime và cho reconnect | High |

## 13. Ngân hàng câu hỏi

| ID | Test case | Kết quả mong đợi | Ưu tiên |
|---|---|---|---|
| TC-REC-QB-001 | Danh sách và filter | Chỉ câu hỏi COM-A, filter/pagination đúng | High |
| TC-REC-QB-002 | Tạo câu hỏi hợp lệ | Nội dung, loại, mức độ, kỹ năng được lưu đúng | High |
| TC-REC-QB-003 | Nội dung rỗng | Từ chối | High |
| TC-REC-QB-004 | Cập nhật câu hỏi | UI và DB đồng bộ | High |
| TC-REC-QB-005 | Xóa câu hỏi | Có xác nhận; xử lý tham chiếu template an toàn | High |
| TC-REC-QB-006 | CRUD câu hỏi COM-B | REC-A bị từ chối | Critical |
| TC-REC-QB-007 | Tạo bằng AI từ Job | Lưu đúng company/job; lỗi AI có retry | High |
| TC-REC-QB-008 | XSS/Unicode trong câu hỏi | Hiển thị an toàn và đúng encoding | Critical |

## 14. Rubric chấm điểm

| ID | Test case | Kết quả mong đợi | Ưu tiên |
|---|---|---|---|
| TC-REC-RUB-001 | Danh sách Rubric | Chỉ Rubric COM-A | High |
| TC-REC-RUB-002 | Tạo Rubric hợp lệ | Tên, tiêu chí, trọng số lưu đúng | High |
| TC-REC-RUB-003 | Trọng số không hợp lệ | Từ chối số âm, vượt giới hạn hoặc tổng sai theo quy tắc | High |
| TC-REC-RUB-004 | Sửa Rubric | Không làm sai dữ liệu report đã chốt | High |
| TC-REC-RUB-005 | Xóa Rubric đang được dùng | Chặn hoặc xử lý tham chiếu rõ ràng | Critical |
| TC-REC-RUB-006 | Truy cập Rubric COM-B | Bị từ chối | Critical |
| TC-REC-RUB-007 | Tên/tiêu chí chứa XSS | Không thực thi script | Critical |

## 15. Template phỏng vấn

| ID | Test case | Kết quả mong đợi | Ưu tiên |
|---|---|---|---|
| TC-REC-TPL-001 | Danh sách/filter Template | Chỉ dữ liệu COM-A | High |
| TC-REC-TPL-002 | Tạo Template hợp lệ | Lưu cấu trúc, câu hỏi, Rubric đúng thứ tự | High |
| TC-REC-TPL-003 | Template rỗng/thiếu tên | Từ chối | High |
| TC-REC-TPL-004 | Sửa Template | Dữ liệu mới dùng cho lần sau; lịch cũ không hỏng | High |
| TC-REC-TPL-005 | Xóa Template đang dùng | Chặn hoặc xử lý quan hệ an toàn | Critical |
| TC-REC-TPL-006 | Truy cập Template COM-B | Bị từ chối | Critical |
| TC-REC-TPL-007 | Duplicate submit | Không tạo hai Template giống nhau ngoài ý muốn | Medium |

## 16. AI trong phỏng vấn và báo cáo

| ID | Test case | Kết quả mong đợi | Ưu tiên |
|---|---|---|---|
| TC-REC-AI-001 | Gợi ý câu hỏi tiếp theo | Dựa đúng interview/transcript, không lẫn công ty | High |
| TC-REC-AI-002 | Chấm câu trả lời | Score nằm trong miền hợp lệ và gắn đúng Rubric | Critical |
| TC-REC-AI-003 | Tạo báo cáo khi kết thúc | Trạng thái generating → ready; dữ liệu đúng interview | Critical |
| TC-REC-AI-004 | AI timeout | Không làm mất interview/transcript; hiển thị trạng thái retry | High |
| TC-REC-AI-005 | Xem report hoàn thành | Hiển thị score, summary, điểm mạnh/yếu, recommendation | Critical |
| TC-REC-AI-006 | Report chưa sẵn sàng | Hiển thị loading/pending, không crash | High |
| TC-REC-AI-007 | Retry report lỗi | Giữ report cũ/trạng thái nhất quán; báo lỗi | High |
| TC-REC-AI-008 | Retry report thành công | Chỉ một job hợp lệ hoặc xử lý idempotent | High |
| TC-REC-AI-009 | Recruiter ghi đè quyết định | `PUT .../report/decision`; lưu decision và lý do | Critical |
| TC-REC-AI-010 | Decision không hợp lệ | Từ chối, report không đổi | High |
| TC-REC-AI-011 | Report COM-B | REC-A bị từ chối | Critical |
| TC-REC-AI-012 | Prompt injection trong CV/transcript | AI output chỉ là dữ liệu; không thực thi hành động/quyền ngoài phạm vi | Critical |
| TC-REC-AI-013 | Nội dung report có XSS | FE escape/sanitize trước khi render | Critical |

## 17. Thông báo và Audit log

| ID | Test case | Kết quả mong đợi | Ưu tiên |
|---|---|---|---|
| TC-REC-NOT-001 | Danh sách thông báo | Chỉ thông báo của user hiện tại; phân trang đúng | High |
| TC-REC-NOT-002 | Đánh dấu đã đọc một thông báo | Badge/count cập nhật đúng | Medium |
| TC-REC-NOT-003 | Đọc tất cả | Mọi thông báo thuộc user được cập nhật, không ảnh hưởng user khác | High |
| TC-REC-NOT-004 | Sửa notification ID người khác | Bị từ chối | Critical |
| TC-REC-NOT-005 | Sự kiện lịch/report tạo notification | Nội dung, người nhận và link điều hướng đúng | High |
| TC-REC-AUD-001 | Danh sách audit COM-A | Chỉ log COM-A và filter/pagination đúng | High |
| TC-REC-AUD-002 | Hành động nhạy cảm tạo audit | Create/update/delete/cancel/decision ghi actor, resource, time đúng | Critical |
| TC-REC-AUD-003 | REC-B đọc audit COM-A | Bị từ chối | Critical |
| TC-REC-AUD-004 | Dữ liệu nhạy cảm trong audit | Không log password, token, nội dung file bí mật | Critical |
| TC-REC-AUD-005 | Realtime audit khi tải cao | Không làm block luồng chính; không mất log vượt chính sách queue | High |

## 18. Kiểm thử bảo mật

| ID | Test case | Kết quả mong đợi | Ưu tiên |
|---|---|---|---|
| TC-REC-SEC-001 | SQL injection trong keyword/filter | Không thay đổi câu query; không lộ lỗi DB | Critical |
| TC-REC-SEC-002 | Stored/reflected XSS | Mọi dữ liệu Job, Candidate, note, chat, transcript được escape/sanitize | Critical |
| TC-REC-SEC-003 | Path traversal qua tên file | Không đọc/ghi ngoài storage cho phép | Critical |
| TC-REC-SEC-004 | Upload executable đổi đuôi | Magic-byte/MIME validation từ chối | Critical |
| TC-REC-SEC-005 | IDOR toàn bộ resource | Job/Candidate/CV/Interview/Report/QB/Rubric/Template đều company-scoped | Critical |
| TC-REC-SEC-006 | CORS từ origin không cho phép | Browser request bị chặn theo cấu hình | High |
| TC-REC-SEC-007 | Brute force Login/OTP | Có rate limit/lockout phù hợp; không DoS user hợp lệ | High |
| TC-REC-SEC-008 | Token trong URL/log | Access/refresh token không xuất hiện trong URL hoặc log công khai | Critical |
| TC-REC-SEC-009 | Signed URL bị chia sẻ | URL hết hạn ngắn, không mở rộng quyền sang file khác | Critical |
| TC-REC-SEC-010 | Mass assignment | Không cập nhật `company_id`, role, owner hoặc field hệ thống trái phép | Critical |
| TC-REC-SEC-011 | Payload JSON quá lớn | Từ chối có kiểm soát, không làm cạn tài nguyên | High |
| TC-REC-SEC-012 | CSRF nếu dùng cookie auth | Yêu cầu biện pháp CSRF/SameSite phù hợp | High |
| TC-REC-SEC-013 | Lỗi nội bộ | Response không lộ stack trace, SQL, storage path hay secret | Critical |

## 19. Hiệu năng, tải và độ tin cậy

| ID | Test case | Tiêu chí mong đợi | Ưu tiên |
|---|---|---|---|
| TC-REC-PERF-001 | Danh sách 10.000 Candidate | API dùng pagination, UI không render toàn bộ, thời gian đạt SLA dự án | High |
| TC-REC-PERF-002 | Danh sách Job/Interview lớn | Filter và pagination ổn định, không trùng/thiếu giữa trang | High |
| TC-REC-PERF-003 | 20 Recruiter cập nhật pipeline | Không corrupt dữ liệu; lỗi xung đột rõ ràng | High |
| TC-REC-PERF-004 | Nhiều upload CV đồng thời | Không trùng storage key, checksum/metadata đúng | Critical |
| TC-REC-PERF-005 | Nhiều yêu cầu parse/report AI | Có timeout/queue/retry hợp lý; API chính không treo | High |
| TC-REC-PERF-006 | Ổn định phòng phỏng vấn dài | Không rò memory/track; reconnect hoạt động | High |
| TC-REC-PERF-007 | DB/AI/Redis tạm ngắt | Trả lỗi kiểm soát; phục hồi không tạo thao tác trùng | High |
| TC-REC-PERF-008 | Race cancel/start/end | Chỉ một transition trạng thái hợp lệ được chấp nhận | Critical |
| TC-REC-PERF-009 | Refresh trang khi đang submit | Không tạo Job/interview/decision trùng | High |

## 20. UI/UX và tương thích

| ID | Test case | Kết quả mong đợi | Ưu tiên |
|---|---|---|---|
| TC-REC-UX-001 | Responsive | Các trang chính dùng được trên desktop/tablet/mobile hỗ trợ | Medium |
| TC-REC-UX-002 | Loading state | Nút submit bị disable, skeleton/spinner nhất quán | Medium |
| TC-REC-UX-003 | Empty state | Có thông báo và CTA phù hợp | Medium |
| TC-REC-UX-004 | Error state | Thông báo tiếng Việt rõ ràng; có retry khi phù hợp | High |
| TC-REC-UX-005 | Modal nguy hiểm | Xóa/hủy có xác nhận, focus và Escape hoạt động hợp lý | High |
| TC-REC-UX-006 | Điều hướng bằng bàn phím | Tab order, focus visible và nút hoạt động bằng keyboard | Medium |
| TC-REC-UX-007 | Màu/trạng thái | Không chỉ dựa vào màu; label Scheduled/Cancelled/Completed rõ ràng | Medium |
| TC-REC-UX-008 | Trình duyệt | Chrome, Edge, Firefox bản hỗ trợ hoạt động nhất quán | Medium |

## 21. Ma trận API Recruiter cần kiểm thử

| Nhóm | Method | Endpoint chính |
|---|---|---|
| Auth | POST | `/auth/register`, `/auth/register/send-otp`, `/auth/login`, `/auth/refresh`, `/auth/logout`, `/auth/logout-all` |
| Account | GET/PUT | `/auth/me`, `/auth/me/settings`, `/auth/me/password` |
| Company | GET/POST/PUT | `/companies`, `/companies/{company_id}` |
| Job | GET/POST | `/companies/{company_id}/jobs` |
| Job | GET/PUT/DELETE | `/companies/{company_id}/jobs/{job_id}` |
| Job AI | POST | `/companies/{company_id}/jobs/{job_id}/analyze`, `.../ai/generate-questions` |
| Candidate | GET/POST | `/companies/{company_id}/candidates` |
| Candidate | GET/PUT/DELETE | `/companies/{company_id}/candidates/{candidate_id}` |
| Candidate CV | POST | `/companies/{company_id}/candidates/{candidate_id}/cv`, `.../parse-cv` |
| File | GET | `/files/{file_id}/signed-url?company_id={company_id}` |
| Pipeline | GET | `/companies/{company_id}/jobs/{job_id}/candidates` |
| Pipeline | POST/PUT/DELETE | `.../candidates/{candidate_id}/assign`, `.../pipeline`, `.../unassign` |
| Interview | GET/POST | `/companies/{company_id}/interviews` |
| Interview | GET | `/companies/{company_id}/interviews/{interview_id}` |
| Interview lifecycle | POST | `.../start`, `.../end`, `.../cancel`, `.../send-reminder` |
| Notes | PUT | `.../interviews/{interview_id}/notes` |
| Room | GET/POST | `.../room`, `.../room/access-token`, `.../room/token` |
| Transcript | GET/POST/PUT | `.../interviews/{interview_id}/transcripts[/ {transcript_id}]` |
| AI interview | POST | `.../ai/suggest-follow-up`, `.../ai/score-answer`, `.../ai/generate-report` |
| Report | GET/PUT/POST | `.../report`, `.../report/decision`, `.../report/retry` |
| Question bank | GET/POST/PUT/DELETE | `/companies/{company_id}/question-bank[/ {question_id}]` |
| Rubric | GET/POST/PUT/DELETE | `/companies/{company_id}/rubrics[/ {rubric_id}]` |
| Template | GET/POST/PUT/DELETE | `/companies/{company_id}/templates[/ {template_id}]` |
| Notification | GET/PUT | `/notifications`, `/notifications/read-all`, `/notifications/{id}/read` |
| Audit | GET | `/companies/{company_id}/audit-logs` |

> Khi thi hành automation, cần xác nhận lại dấu gạch chéo cuối endpoint và tên prefix được mount trong `cmd/api/main.go`; bảng trên mô tả contract mà frontend hiện đang gọi.

## 22. Chuỗi E2E ưu tiên cao

### E2E-REC-01 – Tuyển dụng từ Job đến quyết định

1. Recruiter đăng nhập và tạo Job.
2. Candidate ứng tuyển Job với CV.
3. Recruiter thấy Candidate trong đúng Job.
4. Recruiter mở và parse CV.
5. Cập nhật pipeline.
6. Tạo lịch phỏng vấn.
7. Candidate thấy lịch và tham gia bằng token.
8. Recruiter bắt đầu, ghi transcript/notes và kết thúc.
9. AI tạo report.
10. Recruiter xem report và lưu quyết định.

**Kết quả:** Mọi resource thuộc cùng công ty, trạng thái nối tiếp nhất quán, không cần nhập lại dữ liệu và có audit phù hợp.

### E2E-REC-02 – Candidate cập nhật CV từ Portal

1. Candidate có đơn ứng tuyển tại COM-A.
2. Candidate upload và bấm lưu CV mới.
3. Recruiter COM-A tải lại trang chi tiết.
4. Bấm **Xem CV** và **Phân tích CV**.
5. Recruiter COM-B thử truy cập cùng `file_id`.

**Kết quả:** Recruiter COM-A xem/parse được file thật; COM-B bị từ chối; không còn lỗi `CV file not found` với dữ liệu mới hợp lệ.

### E2E-REC-03 – Hủy lịch phỏng vấn

1. Tạo lịch Scheduled.
2. Mở danh sách và bấm **Hủy lịch**.
3. Nhập lý do, xác nhận.
4. Kiểm tra Candidate Portal và thử vào phòng.

**Kết quả:** Trạng thái cancelled ở cả hai phía; không thể start/join trái phép; lý do/audit/notification đúng chính sách.

### E2E-REC-04 – Cách ly hai công ty

1. Tạo Job, Candidate, CV, Interview, Report tại COM-A và COM-B.
2. Dùng REC-OWNER-A thay ID/URL sang toàn bộ resource COM-B.
3. Thử cả GET và thao tác ghi/xóa.

**Kết quả:** Tất cả bị từ chối; không lộ sự tồn tại, metadata hoặc signed URL của COM-B.

## 23. Tiêu chí chấp nhận

Bản Recruiter đạt yêu cầu khi:

1. 100% test Critical về đăng nhập, RBAC, company scope, CV và lifecycle interview đạt.
2. Recruiter không thể đọc hoặc sửa resource của công ty khác bằng thay URL/ID.
3. Candidate lưu CV mới thì Recruiter thuộc công ty ứng tuyển xem và parse được.
4. Job, Candidate, pipeline, interview và report duy trì trạng thái nhất quán qua cả FE và BE.
5. Lịch đã hủy không thể bắt đầu; lịch đã hoàn thành không bị kết thúc/hủy lần nữa trái quy tắc.
6. File upload chống MIME spoofing/path traversal và signed URL có phạm vi, thời hạn phù hợp.
7. Realtime mất kết nối không làm mất/nhân đôi dữ liệu nghiệp vụ quan trọng.
8. AI lỗi hoặc timeout không làm mất CV, transcript, Job hoặc interview.
9. Không có lỗi 5xx do input thông thường; không lộ stack trace, SQL, secret hay storage path.
10. Các thao tác nguy hiểm có xác nhận, thông báo lỗi và loading state rõ ràng.

## 24. Lưu ý/điểm cần theo dõi từ code hiện tại

1. **CV cũ:** Các Candidate được tạo trước khi sửa luồng portal có thể vẫn trỏ tới metadata giả/file thiếu; cần upload lại hoặc chạy migration sửa dữ liệu trước khi kỳ vọng test xem/parse CV đạt.
2. **Quyền xem Portal CV:** Backend cho phép file owner-scoped khi file đó đang được gắn với Candidate thuộc công ty; phải giữ test IDOR để tránh mở quyền sang file bất kỳ của Candidate.
3. **Hủy lịch:** Frontend đã có modal và gọi `POST .../cancel`; cần kiểm tra trạng thái UI chỉ đổi sau khi API thành công.
4. **AI service Docker:** Cấu hình port container từng có dấu hiệu không khớp giữa Docker mapping và Uvicorn; cần xác nhận trước khi chạy test AI E2E.
5. **Realtime health:** Không nên kết luận gateway hỏng chỉ từ `/health` trả 404; cần kiểm tra đúng WebSocket route và flow lấy room token.
6. **Route/API contract:** Một số service AI/template/rubric dùng endpoint theo API spec; automation cần đối chiếu route thực tế được mount trong `backend/cmd/api/main.go` trước khi khóa contract.
7. **Reschedule:** Service Interview hiện có create/start/end/cancel/notes/reminder nhưng chưa thấy method reschedule riêng; nếu UI yêu cầu đổi lịch, cần xác định đang dùng update endpoint nào hoặc ghi nhận là chức năng chưa triển khai.
8. **Company membership UI:** Cần xác nhận giao diện thêm/xóa thành viên và phân vai có route/API hoàn chỉnh trước khi đưa vào regression bắt buộc; RBAC backend vẫn phải được test độc lập.

---

## 25. Test case mở rộng – Lifecycle và trạng thái nghiệp vụ

| ID | Kịch bản mở rộng | Dữ liệu/bước chính | Kết quả mong đợi | Ưu tiên |
|---|---|---|---|---|
| TC-REC-EXT-LIFE-001 | Job Draft → Open → Closed | Chuyển lần lượt trạng thái và kiểm tra Public Job Board | Chỉ Job Open hợp lệ được công khai; dữ liệu Candidate cũ không mất | Critical |
| TC-REC-EXT-LIFE-002 | Candidate pipeline đầy đủ | New → Screening → Interviewing → Offered → Hired/Rejected | Chỉ transition hợp lệ; dashboard/report cập nhật đúng | Critical |
| TC-REC-EXT-LIFE-003 | Candidate rút đơn giữa pipeline | Candidate withdraw khi đang screening/interviewing | Recruiter thấy trạng thái mới và không tiếp tục thao tác trái quy tắc | High |
| TC-REC-EXT-LIFE-004 | Interview no-show | Recruiter đánh dấu/ghi nhận vắng mặt nếu được hỗ trợ | Trạng thái, audit và notification đúng; không giả thành completed | High |
| TC-REC-EXT-LIFE-005 | Đổi lịch phỏng vấn | Thay scheduled_at/duration khi chức năng được triển khai | Hai phía nhận giờ mới; invite/link cũ xử lý đúng; không tạo lịch trùng | Critical |
| TC-REC-EXT-LIFE-006 | Đổi lịch chưa có endpoint | Kiểm tra UI hiện tại | Không hiển thị action giả hoặc gọi endpoint không tồn tại; ghi nhận gap rõ | High |
| TC-REC-EXT-LIFE-007 | Reminder sau cancel/completed | Gửi reminder cho lịch không còn Scheduled | Bị từ chối, không gửi email/notification sai | High |
| TC-REC-EXT-LIFE-008 | Decision sau report | Chọn hire/reject/hold và nhập lý do | Decision đúng whitelist, audit before/after, không sửa nhầm interview | Critical |
| TC-REC-EXT-LIFE-009 | Sửa decision đồng thời | Hai Recruiter gửi quyết định khác nhau | Có conflict/version hoặc last-write policy rõ, không mất dấu audit | Critical |

## 26. Test case mở rộng – Realtime, transcript và báo cáo

| ID | Kịch bản | Kết quả mong đợi | Ưu tiên |
|---|---|---|---|
| TC-REC-EXT-RT-001 | Refresh room khi đang phỏng vấn | Lấy token mới và reconnect đúng identity/room | Critical |
| TC-REC-EXT-RT-002 | Mở cùng interview ở hai tab | Không nhân đôi track/transcript/command ngoài ý muốn | High |
| TC-REC-EXT-RT-003 | Candidate reconnect nhiều lần | Participant state được dọn đúng, không có ghost user | High |
| TC-REC-EXT-RT-004 | Message đến sai thứ tự | UI sắp theo sequence/timestamp và không mất message | High |
| TC-REC-EXT-RT-005 | Transcript duplicate event | Backend/UI dedup theo ID/sequence | Critical |
| TC-REC-EXT-RT-006 | Edit transcript sau report generating | Áp dụng chính sách khóa hoặc đánh dấu report cần regenerate | High |
| TC-REC-EXT-RT-007 | End interview khi gateway mất kết nối | REST end vẫn nhất quán; room đóng khi gateway phục hồi | Critical |
| TC-REC-EXT-REP-001 | Generate report double-click | Chỉ một job/report active hoặc response idempotent | Critical |
| TC-REC-EXT-REP-002 | Worker restart lúc generating | Job retry/khôi phục, không kẹt vô hạn | High |
| TC-REC-EXT-REP-003 | Report dựa trên transcript đã sửa | Version nguồn được ghi nhận; kết quả không dùng dữ liệu stale âm thầm | High |
| TC-REC-EXT-REP-004 | Score ngoài miền | Backend từ chối/chuẩn hóa; UI không vỡ biểu đồ | Critical |
| TC-REC-EXT-REP-005 | Report partial/null field | UI dùng fallback, không hiện `undefined` hoặc crash | Medium |
| TC-REC-EXT-REP-006 | AI trả nội dung nguy hiểm | Escape/sanitize HTML, không thực thi link/script | Critical |

## 27. Test case mở rộng – Multi-company và quyền chi tiết

| ID | Kịch bản | Kết quả mong đợi | Ưu tiên |
|---|---|---|---|
| TC-REC-EXT-SCOPE-001 | Một Recruiter thuộc hai công ty | Chuyển workspace và gọi API liên tiếp | company_id hiện tại đúng; cache/store công ty cũ không rò sang công ty mới | Critical |
| TC-REC-EXT-SCOPE-002 | Mở hai tab hai công ty | Tab COM-A và COM-B thao tác song song | Mỗi request dùng company rõ ràng, không ghi nhầm tenant | Critical |
| TC-REC-EXT-SCOPE-003 | Bị xóa membership giữa phiên | Admin/Owner thu hồi quyền khi Recruiter đang mở trang | Request tiếp theo 403; token không tiếp tục truy cập resource công ty | Critical |
| TC-REC-EXT-SCOPE-004 | File CV owner-scoped chưa gắn Candidate | Đoán file ID Candidate Portal | Không cấp signed URL chỉ vì biết file ID | Critical |
| TC-REC-EXT-SCOPE-005 | Portal CV gắn Candidate đã soft-delete | Yêu cầu signed URL sau soft delete | Không truy cập qua quan hệ Candidate đã xóa | Critical |
| TC-REC-EXT-SCOPE-006 | Cross-company Question/Rubric/Template | Dùng ID COM-B trong endpoint COM-A | 403/404 cho GET/PUT/DELETE, không lộ metadata | Critical |
| TC-REC-EXT-SCOPE-007 | Cross-company AI endpoint | Gửi Job/Interview ID COM-B dưới path COM-A | Bị từ chối trước khi gọi AI, không phát sinh chi phí/job | Critical |
| TC-REC-EXT-SCOPE-008 | Permission thay đổi giữa phiên | Owner hạ quyền update xuống read-only | Request ghi tiếp theo bị từ chối; UI ẩn/disable action sau refresh | High |

## 28. Test case mở rộng – API robustness, quan sát và phục hồi

| ID | Kịch bản | Kết quả mong đợi | Ưu tiên |
|---|---|---|---|
| TC-REC-EXT-API-001 | JSON/MIME sai | Body lỗi hoặc Content-Type không hỗ trợ | 400/415, không ghi dữ liệu nửa chừng | High |
| TC-REC-EXT-API-002 | page_size cực lớn | Danh sách Job/Candidate/Interview | Clamp/từ chối, không cạn tài nguyên | High |
| TC-REC-EXT-API-003 | UUID lỗi | Truy cập detail với ID sai format | 400/404 an toàn, không panic/SQL leak | High |
| TC-REC-EXT-API-004 | Request timeout và retry | Client retry POST create | Dùng idempotency/chống double-submit phù hợp | High |
| TC-REC-EXT-OBS-001 | Correlation ID E2E | Theo dõi FE → API → queue → AI | Có request/job ID để truy vết mà không chứa secret | Medium |
| TC-REC-EXT-OBS-002 | Audit hành động realtime | Start/end/edit transcript/decision | Actor, company, resource, timestamp chính xác | Critical |
| TC-REC-EXT-OBS-003 | Metric queue/report | Tạo tải AI có kiểm soát | Quan sát backlog/failure/retry; không lộ payload nhạy cảm | Medium |
| TC-REC-EXT-REC-001 | Redis ngắt | Create Job/Candidate và end interview | Luồng cốt lõi hoạt động hoặc fail rõ; job async không mất âm thầm | High |
| TC-REC-EXT-REC-002 | AI service ngắt | Parse CV/generate report | Giữ dữ liệu gốc, timeout hợp lý và cho retry | Critical |
| TC-REC-EXT-REC-003 | Storage ngắt | Upload/view CV | Không tạo liên kết CV thành công giả; giữ file cũ | Critical |
| TC-REC-EXT-REC-004 | Database fail giữa transaction | Assign/pipeline/create interview | Rollback đầy đủ, không có relation mồ côi | Critical |

## 29. Test case mở rộng – Accessibility và UI đa thiết bị

| ID | Kịch bản | Kết quả mong đợi | Ưu tiên |
|---|---|---|---|
| TC-REC-EXT-A11Y-001 | Quản lý bảng bằng keyboard | Search/filter/pagination/action menu dùng được không cần chuột | High |
| TC-REC-EXT-A11Y-002 | Modal tạo/xóa/hủy | Focus trap, Escape, label và thông báo lỗi đúng | High |
| TC-REC-EXT-A11Y-003 | Room controls với screen reader | Mic/camera/end call có accessible name và state | High |
| TC-REC-EXT-A11Y-004 | Báo cáo không chỉ dùng màu | Score/recommendation có text/icon hỗ trợ | Medium |
| TC-REC-EXT-UX-001 | Tablet/mobile Recruiter | Danh sách, detail, modal và room không che CTA quan trọng | Medium |
| TC-REC-EXT-UX-002 | Zoom 200% | Không mất chức năng, text không chồng nghiêm trọng | Medium |
| TC-REC-EXT-UX-003 | Back/forward với form | Không gửi lại POST hoặc mất dữ liệu chưa lưu ngoài cảnh báo | Medium |

## 30. Checklist regression Recruiter tối thiểu

- [ ] Login/logout/token expiry, route guard và permission backend.
- [ ] Chuyển workspace và IDOR giữa hai công ty.
- [ ] Job CRUD/status/analyze/generate questions.
- [ ] Candidate CRUD, CV upload/view/parse và pipeline.
- [ ] Candidate Portal lưu CV thì Recruiter đúng công ty xem được.
- [ ] Interview create/list/detail/notes/reminder/start/end/cancel.
- [ ] Room token, media, reconnect, chat và transcript consistency.
- [ ] Report generate/retry/decision và AI failure recovery.
- [ ] Question Bank, Rubric, Template CRUD và company scope.
- [ ] Notification/audit không lộ secret.
- [ ] Các test Critical về concurrency, transaction rollback và storage đạt.

