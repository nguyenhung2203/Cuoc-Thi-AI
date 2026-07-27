# BỘ TEST CASE CHỨC NĂNG ACTOR USER (CANDIDATE)

## 1. Thông tin tài liệu

| Thuộc tính | Giá trị |
|---|---|
| Dự án | ViệcLàm AI – Nền tảng tuyển dụng và phỏng vấn AI |
| Actor chính | User – Candidate/Ứng viên |
| Loại kiểm thử | Functional, API, RBAC, UI/UX, Security, Performance |
| Frontend | Vue 3/Vite |
| Backend | Go REST API, Realtime Gateway, AI Service |
| API prefix | `/api/v1` |

## 2. Phạm vi kiểm thử

| Phân hệ | Route frontend | API/chức năng chính |
|---|---|---|
| Đăng ký | `/register` | Đăng ký Candidate |
| Đăng nhập | `/login` | Xác thực và tạo phiên |
| Quên mật khẩu | `/forgot-password` | Gửi yêu cầu khôi phục |
| Dashboard | `/home` | `GET /portal/dashboard`, `GET /portal/interviews` |
| Hồ sơ | `/profile` | `GET/PUT /portal/profile` |
| CV | `/profile` | `POST /portal/cv` |
| Việc làm | `/job-board` | Danh sách/tìm kiếm công việc |
| Chi tiết & ứng tuyển | `/careers/:company_id/jobs/:job_id` | Apply, AI CV-job match |
| Đơn ứng tuyển | `/my-applications` | `GET/DELETE /portal/applications` |
| Việc đã lưu | `/saved-jobs` | Local storage `candidate_saved_jobs` |
| Lịch phỏng vấn | `/my-interviews` | `GET /portal/interviews` |
| Phòng phỏng vấn | `/candidate-room` | Token phòng, realtime, media |
| Phỏng vấn thử | `/mock-setup`, `/mock-room` | Tạo/bắt đầu/kết thúc mock interview |
| Lịch sử luyện tập | `/mock-results` | Danh sách mock interview |
| Kết quả chi tiết | `/mock-results/:id` | Tin nhắn và kết quả mock |
| Bảo mật tài khoản | `/candidate-settings` | Redirect `/profile?tab=security` |

## 3. Dữ liệu kiểm thử

| Mã | Vai trò | Trạng thái | CV | Mục đích |
|---|---|---|---|---|
| USR-CAN-01 | Candidate | Active | Có CV hợp lệ | Luồng đầy đủ |
| USR-CAN-02 | Candidate | Active | Chưa có CV | Kiểm tra CTA upload |
| USR-CAN-03 | Candidate | Blocked | Có CV | Kiểm tra khóa tài khoản |
| USR-CAN-04 | Candidate | Pending/Inactive | Không CV | Kiểm tra trạng thái chưa hoạt động |
| USR-REC-01 | Recruiter | Active | Không áp dụng | Kiểm tra RBAC |
| USR-ADM-01 | Admin | Active | Không áp dụng | Kiểm tra RBAC |

Dữ liệu việc làm:

- JOB-01: Đang mở, còn hạn, thuộc công ty hoạt động.
- JOB-02: Đã đóng.
- JOB-03: Đã hết hạn.
- JOB-04: Đang mở nhưng Candidate đã ứng tuyển.
- JOB-05: Có yêu cầu kỹ năng khớp cao với CV.
- JOB-06: Có yêu cầu kỹ năng khớp thấp với CV.

Dữ liệu CV:

- `valid-cv.pdf`: PDF hợp lệ dưới giới hạn.
- `valid-cv.docx`: DOCX hợp lệ.
- `large-cv.pdf`: Vượt giới hạn upload.
- `fake.pdf`: Đuôi PDF nhưng nội dung executable/text không hợp lệ.
- `image.jpg`: Ảnh hợp lệ nếu hệ thống hỗ trợ.
- `script.html`: Loại file bị cấm.

---

## 4. Test case đăng ký

### TC-USR-REG-001 – Hiển thị trang đăng ký

- **Ưu tiên:** High
- **Bước thực hiện:** Truy cập `/register` khi chưa đăng nhập.
- **Kết quả mong đợi:** Form hiển thị họ tên, email, mật khẩu, xác nhận mật khẩu và lựa chọn vai trò.

### TC-USR-REG-002 – Đăng ký Candidate thành công

- **Ưu tiên:** Critical
- **Bước thực hiện:**
  1. Nhập đầy đủ dữ liệu hợp lệ.
  2. Chọn vai trò Ứng viên.
  3. Bấm đăng ký.
- **Kết quả mong đợi:**
  - API tạo tài khoản thành công.
  - Không lưu mật khẩu dạng plaintext.
  - Chuyển tới `/login` và hiển thị thông báo thành công.

### TC-USR-REG-003 – Thiếu họ tên

- **Kết quả mong đợi:** Không gửi API; hiển thị validation bắt buộc.

### TC-USR-REG-004 – Thiếu email

- **Kết quả mong đợi:** Không gửi API; yêu cầu nhập email.

### TC-USR-REG-005 – Email sai định dạng

- **Dữ liệu:** `candidate-example`
- **Kết quả mong đợi:** Email bị từ chối phía frontend và backend.

### TC-USR-REG-006 – Email đã tồn tại

- **Ưu tiên:** High
- **Kết quả mong đợi:** HTTP 409 hoặc lỗi nghiệp vụ phù hợp; không tạo user trùng.

### TC-USR-REG-007 – Mật khẩu và xác nhận không khớp

- **Kết quả mong đợi:** Không gọi API; hiển thị “Mật khẩu xác nhận không khớp”.

### TC-USR-REG-008 – Mật khẩu quá yếu

- **Dữ liệu:** `123`
- **Ưu tiên:** High
- **Kết quả mong đợi đề xuất:** Từ chối và hiển thị yêu cầu độ dài/độ mạnh.

### TC-USR-REG-009 – Họ tên có Unicode

- **Dữ liệu:** `Đinh Lưu Lại`
- **Kết quả mong đợi:** Đăng ký và hiển thị đúng, không lỗi encoding.

### TC-USR-REG-010 – XSS trong họ tên

- **Dữ liệu:** `<script>alert(1)</script>`
- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Không thực thi script; dữ liệu được validate/escape.

### TC-USR-REG-011 – Bấm đăng ký nhiều lần

- **Kết quả mong đợi:** Nút disable khi xử lý; không tạo tài khoản trùng.

### TC-USR-REG-012 – API đăng ký lỗi

- **Kết quả mong đợi:** Loading kết thúc; hiển thị lỗi; dữ liệu form không bị mất không cần thiết.

---

## 5. Test case đăng nhập và phiên

### TC-USR-AUTH-001 – Candidate đăng nhập thành công

- **Ưu tiên:** Critical
- **Kết quả mong đợi:**
  - Backend trả access/refresh token và user role Candidate.
  - Frontend lưu phiên.
  - Chuyển tới `/home`.

### TC-USR-AUTH-002 – Bỏ trống email/mật khẩu

- **Kết quả mong đợi:** Không gọi API; hiển thị yêu cầu nhập đầy đủ.

### TC-USR-AUTH-003 – Sai mật khẩu

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Không cấp token; hiển thị lỗi xác thực chung.

### TC-USR-AUTH-004 – Email không tồn tại

- **Kết quả mong đợi:** Không tiết lộ tài khoản tồn tại; không tạo phiên.

### TC-USR-AUTH-005 – Candidate bị khóa

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Backend từ chối; không cấp token; thông báo trạng thái tài khoản.

### TC-USR-AUTH-006 – Token hết hạn

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Thử refresh hợp lệ hoặc logout; không tiếp tục dùng API bảo vệ bằng token hết hạn.

### TC-USR-AUTH-007 – Refresh token không hợp lệ

- **Kết quả mong đợi:** Xóa phiên và chuyển login.

### TC-USR-AUTH-008 – Truy cập trang bảo vệ khi chưa đăng nhập

- **Route:** `/profile`, `/my-applications`, `/saved-jobs`, `/my-interviews`, `/mock-setup`.
- **Kết quả mong đợi:** Chuyển `/login` kèm thông báo yêu cầu đăng nhập.

### TC-USR-AUTH-009 – Candidate truy cập trang Recruiter

- **Route:** `/dashboard`, `/jobs`, `/candidates`, `/interviews`, `/settings`.
- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Chuyển `/403?reason=recruiter_required`; backend API cũng trả 403.

### TC-USR-AUTH-010 – Candidate truy cập Admin

- **Kết quả mong đợi:** Chuyển `/403?reason=admin_required`; API Admin trả 403.

### TC-USR-AUTH-011 – Chỉnh role localStorage thành Recruiter/Admin

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Backend vẫn xác định role từ token/database và từ chối quyền sai.

### TC-USR-AUTH-012 – Đăng xuất

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Xóa/revoke token; xóa user state; chuyển về trang công khai/login.

### TC-USR-AUTH-013 – Back sau đăng xuất

- **Kết quả mong đợi:** Không xem lại dữ liệu riêng tư từ cache.

### TC-USR-AUTH-014 – Đăng nhập đồng thời nhiều tab

- **Kết quả mong đợi:** Trạng thái phiên đồng bộ; logout một tab xử lý theo chính sách xác định.

---

## 6. Test case quên mật khẩu

### TC-USR-FORGOT-001 – Hiển thị trang

- **Kết quả mong đợi:** Có email và nút gửi yêu cầu.

### TC-USR-FORGOT-002 – Gửi email hợp lệ

- **Kết quả mong đợi:** Hiển thị thông báo chung rằng hướng dẫn đã được gửi, không tiết lộ email tồn tại.

### TC-USR-FORGOT-003 – Email sai định dạng

- **Kết quả mong đợi:** Không gửi request hoặc backend trả 400.

### TC-USR-FORGOT-004 – Email không tồn tại

- **Ưu tiên:** High
- **Kết quả mong đợi:** Response bên ngoài giống email tồn tại để chống account enumeration.

### TC-USR-FORGOT-005 – Gửi liên tục

- **Ưu tiên:** High
- **Kết quả mong đợi:** Có rate limit/cooldown.

### TC-USR-FORGOT-006 – Token reset hết hạn/đã dùng

- **Kết quả mong đợi:** Từ chối và yêu cầu gửi lại liên kết.

---

## 7. Test case Dashboard Candidate

### TC-USR-DASH-001 – Tải Dashboard

- **API:** `GET /portal/dashboard`, `GET /portal/interviews`.
- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Hiển thị tên Candidate, thống kê, lịch sắp tới và việc gợi ý đúng dữ liệu.

### TC-USR-DASH-002 – Thống kê chính xác

- **Kết quả mong đợi:** Số lịch sắp tới, số lần luyện tập, điểm trung bình và mức hoàn thiện hồ sơ khớp database.

### TC-USR-DASH-003 – Candidate mới chưa có dữ liệu

- **Kết quả mong đợi:** Hiển thị 0/empty state và CTA tìm việc, hoàn thiện hồ sơ, luyện tập.

### TC-USR-DASH-004 – Lịch phỏng vấn sắp tới

- **Kết quả mong đợi:** Thời gian, công ty, vị trí và trạng thái chính xác; chỉ lịch của Candidate hiện tại.

### TC-USR-DASH-005 – Tham gia lịch thực

- **Kết quả mong đợi:** Chỉ hiện nút tham gia với lịch hợp lệ; điều hướng join link đúng.

### TC-USR-DASH-006 – Việc làm gợi ý

- **Kết quả mong đợi:** Bấm card/nút xem việc chuyển đúng `/careers/{company_id}/jobs/{job_id}`.

### TC-USR-DASH-007 – CTA luyện tập

- **Kết quả mong đợi:** Chuyển `/mock-setup`.

### TC-USR-DASH-008 – API Dashboard lỗi

- **Kết quả mong đợi:** Không trắng trang; loading kết thúc; hiển thị lỗi/thử lại.

---

## 8. Test case hồ sơ cá nhân

### TC-USR-PROFILE-001 – Tải hồ sơ

- **API:** `GET /portal/profile`
- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Họ tên, email, avatar, CV và dữ liệu phân tích hiển thị đúng user hiện tại.

### TC-USR-PROFILE-002 – Cập nhật họ tên

- **API:** `PUT /portal/profile`
- **Ưu tiên:** High
- **Kết quả mong đợi:** Lưu thành công; header/store cập nhật; reload vẫn giữ giá trị.

### TC-USR-PROFILE-003 – Họ tên rỗng

- **Kết quả mong đợi:** Frontend/backend từ chối; không lưu chuỗi rỗng nếu bắt buộc.

### TC-USR-PROFILE-004 – Họ tên quá dài

- **Kết quả mong đợi:** Validation theo giới hạn; không gây lỗi layout/database.

### TC-USR-PROFILE-005 – Cập nhật avatar URL hợp lệ

- **Kết quả mong đợi:** Avatar hiển thị đúng sau khi lưu/reload.

### TC-USR-PROFILE-006 – Avatar URL không hợp lệ

- **Kết quả mong đợi:** Dùng fallback an toàn; không thực thi URL nguy hiểm.

### TC-USR-PROFILE-007 – Email không thể sửa trái phép

- **Kết quả mong đợi:** Trường email read-only/disabled; backend bỏ qua field email không hỗ trợ.

### TC-USR-PROFILE-008 – User A truy cập hồ sơ User B

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Không có user ID tùy ý trong portal API; dữ liệu lấy từ token, không xảy ra IDOR.

### TC-USR-PROFILE-009 – API cập nhật lỗi

- **Kết quả mong đợi:** Không báo thành công; UI giữ hoặc rollback dữ liệu phù hợp.

### TC-USR-PROFILE-010 – Mở tab bảo mật

- **Route:** `/candidate-settings`.
- **Kết quả mong đợi:** Redirect `/profile?tab=security` và mở đúng tab.

---

## 9. Test case upload và phân tích CV

### TC-USR-CV-001 – Upload PDF hợp lệ

- **API:** `POST /portal/cv`
- **Field:** `file`
- **Ưu tiên:** Critical
- **Kết quả mong đợi:**
  - File được lưu vật lý và có metadata hợp lệ.
  - Candidate `cv_file_id` trỏ đúng file record.
  - Response có `file_name`, `cv_url`, `cv_file_id`, `parsed_data`.
  - Recruiter của công ty Candidate ứng tuyển có thể xem CV theo quyền.

### TC-USR-CV-002 – Upload DOCX/ảnh hợp lệ

- **Kết quả mong đợi:** Chấp nhận nếu MIME nằm trong whitelist; parse AI xử lý/fallback đúng khả năng định dạng.

### TC-USR-CV-003 – Không chọn file

- **Kết quả mong đợi:** HTTP 400 `file is required`.

### TC-USR-CV-004 – File vượt giới hạn

- **Ưu tiên:** High
- **Kết quả mong đợi:** HTTP 400/413; không tạo file rác hoặc metadata mồ côi.

### TC-USR-CV-005 – File type bị cấm

- **Dữ liệu:** `.html`, `.exe`, `.js`.
- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Từ chối theo magic bytes/MIME, không chỉ phần mở rộng.

### TC-USR-CV-006 – File giả mạo đuôi PDF

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Từ chối do nội dung không phải PDF hợp lệ.

### TC-USR-CV-007 – Upload lại CV

- **Kết quả mong đợi:** Candidate trỏ CV mới; profile và đơn ứng tuyển liên quan dùng đúng CV mới theo nghiệp vụ; file cũ được xử lý theo chính sách.

### TC-USR-CV-008 – Upload file trùng checksum

- **Kết quả mong đợi:** Tái sử dụng storage an toàn hoặc tạo metadata đúng owner; không cấp quyền chéo giữa user/công ty.

### TC-USR-CV-009 – AI parse thành công

- **Kết quả mong đợi:** Skills, experience, education và summary được lưu/hiển thị chính xác.

### TC-USR-CV-010 – AI parse lỗi

- **Ưu tiên:** High
- **Kết quả mong đợi:** Upload vẫn thành công theo thiết kế best-effort; UI không báo đã phân tích nếu chưa có parsed data; có thể thử lại.

### TC-USR-CV-011 – File metadata tồn tại nhưng file vật lý mất

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Không hiển thị sai “Đã lưu”; yêu cầu upload lại; recruiter không gặp lỗi mơ hồ.

### TC-USR-CV-012 – Truy cập CV người khác

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** HTTP 403/404; Candidate chỉ xem file thuộc owner; recruiter chỉ xem CV gắn với Candidate của công ty.

### TC-USR-CV-013 – Tên file Unicode/ký tự đặc biệt

- **Kết quả mong đợi:** Lưu và hiển thị an toàn; không path traversal.

### TC-USR-CV-014 – Tên file path traversal

- **Dữ liệu:** `../../secret.pdf`
- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Storage key do server tạo; không ghi file ngoài thư mục upload.

---

## 10. Test case bảng việc làm

### TC-USR-JOB-001 – Tải danh sách việc làm

- **Route:** `/job-board`
- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Chỉ hiển thị job public/open hợp lệ.

### TC-USR-JOB-002 – Tìm kiếm theo tiêu đề

- **Kết quả mong đợi:** Lọc đúng, không phân biệt hoa/thường.

### TC-USR-JOB-003 – Tìm kiếm theo công ty/kỹ năng

- **Kết quả mong đợi:** Kết quả chính xác theo điều kiện hỗ trợ.

### TC-USR-JOB-004 – Lọc địa điểm, loại việc, cấp độ

- **Kết quả mong đợi:** Bộ lọc kết hợp chính xác; có thể reset.

### TC-USR-JOB-005 – Không có kết quả

- **Kết quả mong đợi:** Empty state và CTA xóa bộ lọc.

### TC-USR-JOB-006 – Job đóng/hết hạn

- **Kết quả mong đợi:** Không hiển thị hoặc hiển thị trạng thái không thể ứng tuyển.

### TC-USR-JOB-007 – Xem chi tiết

- **Kết quả mong đợi:** Chuyển đúng career company/job; dữ liệu đúng ID.

### TC-USR-JOB-008 – Job ID không tồn tại

- **Kết quả mong đợi:** 404/empty state; không crash.

### TC-USR-JOB-009 – Phân trang/infinite scroll

- **Kết quả mong đợi:** Không trùng/thiếu job; thứ tự ổn định.

### TC-USR-JOB-010 – API danh sách lỗi

- **Kết quả mong đợi:** Loading kết thúc; có lỗi/thử lại.

---

## 11. Test case lưu việc làm

### TC-USR-SAVED-001 – Lưu một việc làm

- **Kết quả mong đợi:** Job được thêm vào `candidate_saved_jobs`; icon/trạng thái cập nhật.

### TC-USR-SAVED-002 – Lưu trùng

- **Kết quả mong đợi:** Không tạo bản ghi trùng theo job ID.

### TC-USR-SAVED-003 – Bỏ lưu

- **Kết quả mong đợi:** Job bị xóa khỏi danh sách đã lưu và trạng thái card cập nhật.

### TC-USR-SAVED-004 – Tải trang việc đã lưu

- **Route:** `/saved-jobs`
- **Kết quả mong đợi:** Hiển thị đúng local storage của Candidate hiện tại.

### TC-USR-SAVED-005 – Không có việc đã lưu

- **Kết quả mong đợi:** Empty state và CTA `/job-board`.

### TC-USR-SAVED-006 – Local storage JSON hỏng

- **Ưu tiên:** High
- **Kết quả mong đợi:** Không crash; reset/fallback danh sách rỗng.

### TC-USR-SAVED-007 – Job đã lưu sau đó bị đóng/xóa

- **Kết quả mong đợi:** Hiển thị trạng thái không khả dụng hoặc loại bỏ an toàn.

### TC-USR-SAVED-008 – Bấm job đã lưu

- **Kết quả mong đợi:** Điều hướng đúng chi tiết.

---

## 12. Test case chi tiết và ứng tuyển việc làm

### TC-USR-APPLY-001 – Xem chi tiết job mở

- **Kết quả mong đợi:** Hiển thị công ty, mô tả, yêu cầu, địa điểm, hình thức, hạn nộp.

### TC-USR-APPLY-002 – AI match khi có CV

- **API:** `GET /portal/jobs/{jobID}/match`
- **Ưu tiên:** High
- **Kết quả mong đợi:** `has_cv=true`; điểm 0–100; kỹ năng khớp/thiếu, summary, recommendation đúng định dạng.

### TC-USR-APPLY-003 – AI match khi chưa có CV

- **Kết quả mong đợi:** HTTP 200 với `has_cv=false`; UI hiển thị CTA upload CV, không tạo điểm giả.

### TC-USR-APPLY-004 – Ứng tuyển bằng CV hiện tại

- **API:** `POST /portal/jobs/{jobID}/apply`
- **Ưu tiên:** Critical
- **Kết quả mong đợi:**
  - Tạo/update Candidate thuộc đúng công ty.
  - Tạo job_candidate duy nhất.
  - Đơn xuất hiện trong `/my-applications`.
  - Recruiter nhận thông báo ứng viên mới.
  - Candidate nhận thông báo thành công.

### TC-USR-APPLY-005 – Ứng tuyển với CV mới

- **Field:** `cv_file`
- **Kết quả mong đợi:** Upload file đúng owner, gắn Candidate/đơn ứng tuyển và recruiter xem được.

### TC-USR-APPLY-006 – Job đã đóng

- **Kết quả mong đợi:** HTTP 400/409; không tạo đơn.

### TC-USR-APPLY-007 – Job không tồn tại

- **Kết quả mong đợi:** HTTP 404.

### TC-USR-APPLY-008 – Ứng tuyển trùng

- **Kết quả mong đợi:** Không tạo `job_candidates` trùng; trả trạng thái rõ ràng hoặc idempotent.

### TC-USR-APPLY-009 – Ứng tuyển chưa đăng nhập

- **Kết quả mong đợi:** Chuyển login/HTTP 401; không tạo đơn ẩn danh.

### TC-USR-APPLY-010 – CV upload lỗi trong apply

- **Kết quả mong đợi:** Không tạo đơn nửa chừng; không để metadata/file mồ côi không cần thiết.

### TC-USR-APPLY-011 – AI tính fit score lỗi

- **Kết quả mong đợi:** Đơn vẫn tạo thành công theo best-effort; fit score có thể null; không rollback ứng tuyển.

### TC-USR-APPLY-012 – Hai request apply đồng thời

- **Kết quả mong đợi:** Chỉ một đơn cho cùng Candidate/job.

---

## 13. Test case đơn ứng tuyển

### TC-USR-APP-001 – Tải đơn của tôi

- **API:** `GET /portal/applications`
- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Chỉ trả đơn thuộc Candidate đang đăng nhập; sắp xếp mới nhất trước.

### TC-USR-APP-002 – Hiển thị trạng thái

- **Trạng thái:** New, screening, interviewing, offered, rejected/withdrawn.
- **Kết quả mong đợi:** Nhãn và màu đúng; không hiện raw value khó hiểu.

### TC-USR-APP-003 – Không có đơn

- **Kết quả mong đợi:** Empty state và CTA tìm việc.

### TC-USR-APP-004 – Xem job từ đơn

- **Kết quả mong đợi:** Chuyển đúng company/job detail.

### TC-USR-APP-005 – Hủy/rút đơn

- **API:** `DELETE /portal/applications/{id}`
- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Có xác nhận; chỉ hủy đơn của chính user; UI cập nhật sau thành công.

### TC-USR-APP-006 – Hủy xác nhận

- **Kết quả mong đợi:** Không gọi API; đơn giữ nguyên.

### TC-USR-APP-007 – Hủy đơn người khác bằng ID

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** HTTP 403/404; dữ liệu không đổi.

### TC-USR-APP-008 – Hủy đơn đã vào giai đoạn không cho phép

- **Kết quả mong đợi:** Backend áp dụng business rule và trả 409/400 rõ ràng.

### TC-USR-APP-009 – API danh sách/hủy lỗi

- **Kết quả mong đợi:** Không xóa card trước khi backend xác nhận; hiển thị lỗi.

---

## 14. Test case lịch phỏng vấn

### TC-USR-IV-001 – Tải lịch của Candidate

- **API:** `GET /portal/interviews`
- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Chỉ lịch thuộc user; đúng công ty, vị trí, mode, trạng thái và thời gian.

### TC-USR-IV-002 – Sắp xếp mới nhất trước

- **Kết quả mong đợi:** Danh sách giảm dần theo `scheduled_at`; không mutate dữ liệu nguồn ngoài ý muốn.

### TC-USR-IV-003 – Không có lịch

- **Kết quả mong đợi:** Empty state và CTA luyện tập AI.

### TC-USR-IV-004 – Lịch scheduled thực

- **Kết quả mong đợi:** Hiển thị nút tham gia nếu join link hợp lệ.

### TC-USR-IV-005 – Lịch completed

- **Kết quả mong đợi:** Có badge hoàn thành; không hiện nút tham gia.

### TC-USR-IV-006 – Lịch cancelled

- **Kết quả mong đợi:** Badge đã hủy; không hiện nút tham gia.

### TC-USR-IV-007 – Join link thiếu

- **Kết quả mong đợi:** Nút disable/ẩn hoặc thông báo; không click no-op khó hiểu.

### TC-USR-IV-008 – Lịch ở timezone khác

- **Kết quả mong đợi:** Hiển thị đúng giờ địa phương, định dạng Việt Nam.

### TC-USR-IV-009 – Candidate A truy cập lịch Candidate B

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Backend scope theo `user_id` trong token, không cho IDOR.

### TC-USR-IV-010 – Recruiter hủy lịch

- **Kết quả mong đợi:** Candidate thấy trạng thái cancelled sau refresh/realtime và không vào được phòng.

---

## 15. Test case consent và phòng phỏng vấn thật

### TC-USR-ROOM-001 – Mở trang consent

- **Route:** `/interview-consent`
- **Kết quả mong đợi:** Hiển thị nội dung ghi âm, transcript, AI và quyền riêng tư.

### TC-USR-ROOM-002 – Chưa đồng ý consent

- **Kết quả mong đợi:** Nút tham gia bị disable.

### TC-USR-ROOM-003 – Đồng ý và tham gia

- **Kết quả mong đợi:** Điều hướng phòng hợp lệ; consent được ghi nhận theo phiên/interview nếu yêu cầu.

### TC-USR-ROOM-004 – Từ chối consent

- **Kết quả mong đợi:** Không vào phòng; chuyển trang phù hợp.

### TC-USR-ROOM-005 – Truy cập room không có token mời

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Từ chối/redirect; không chỉ dựa token đăng nhập.

### TC-USR-ROOM-006 – Token phòng hết hạn

- **Kết quả mong đợi:** Chuyển `/interview-expired` hoặc hiển thị thông báo hết hạn.

### TC-USR-ROOM-007 – Token của cuộc phỏng vấn khác

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Từ chối tham gia.

### TC-USR-ROOM-008 – Bật/tắt microphone

- **Kết quả mong đợi:** Track media thực tế thay đổi và icon đồng bộ.

### TC-USR-ROOM-009 – Bật/tắt camera

- **Kết quả mong đợi:** Track video thực tế thay đổi và preview đồng bộ.

### TC-USR-ROOM-010 – Từ chối quyền camera/microphone

- **Kết quả mong đợi:** Thông báo rõ ràng, hướng dẫn cấp quyền; không crash.

### TC-USR-ROOM-011 – Mất kết nối và reconnect

- **Ưu tiên:** High
- **Kết quả mong đợi:** Hiển thị trạng thái; tự kết nối lại an toàn; không tạo participant trùng.

### TC-USR-ROOM-012 – Recruiter hủy khi Candidate đang trong phòng

- **Kết quả mong đợi:** Candidate nhận sự kiện, rời phòng và thấy thông báo đã hủy.

### TC-USR-ROOM-013 – Cuộc phỏng vấn kết thúc

- **Kết quả mong đợi:** Candidate nhận sự kiện, chuyển home và thấy thông báo cảm ơn.

### TC-USR-ROOM-014 – Rời phòng chủ động

- **Kết quả mong đợi:** Có modal xác nhận; cancel giữ phòng; confirm đóng media/socket và chuyển home.

### TC-USR-ROOM-015 – Chat

- **Kết quả mong đợi:** Chỉ gửi vào đúng room; chống XSS; không gửi rỗng; thứ tự và timestamp đúng.

### TC-USR-ROOM-016 – Ghi chú Candidate

- **Kết quả mong đợi:** Thực hiện theo chính sách lưu/không lưu rõ ràng; không lộ cho người khác nếu là ghi chú riêng.

---

## 16. Test case thiết lập phỏng vấn thử

### TC-USR-MOCK-SET-001 – Tải trang thiết lập

- **Kết quả mong đợi:** Có vai trò mục tiêu, cấp độ, loại phỏng vấn, phong cách AI và lựa chọn CV.

### TC-USR-MOCK-SET-002 – Tạo phiên hợp lệ

- **Kết quả mong đợi:**
  - Gọi API tạo mock interview.
  - Bắt đầu phiên thành công.
  - Chuyển `/mock-room` kèm mock ID.

### TC-USR-MOCK-SET-003 – Thay đổi từng lựa chọn

- **Kết quả mong đợi:** Payload gửi đúng role, level, type, style.

### TC-USR-MOCK-SET-004 – Dùng CV hiện tại

- **Kết quả mong đợi:** Gửi đúng `cv_file_id` thuộc user.

### TC-USR-MOCK-SET-005 – Chọn dùng CV nhưng chưa có CV

- **Kết quả mong đợi:** Hiển thị yêu cầu upload hoặc tự tắt tùy chọn; không gửi ID rỗng sai.

### TC-USR-MOCK-SET-006 – API tạo/start lỗi

- **Kết quả mong đợi:** Không chuyển room; loading kết thúc; có thể thử lại.

### TC-USR-MOCK-SET-007 – Bấm bắt đầu nhiều lần

- **Kết quả mong đợi:** Chỉ tạo một phiên.

---

## 17. Test case phòng phỏng vấn thử

### TC-USR-MOCK-ROOM-001 – Vào phiên hợp lệ

- **Kết quả mong đợi:** Mock ID hợp lệ, trạng thái running, câu hỏi AI hiển thị.

### TC-USR-MOCK-ROOM-002 – Truy cập không có mock ID

- **Kết quả mong đợi:** Redirect setup/lịch sử hoặc báo phiên không hợp lệ.

### TC-USR-MOCK-ROOM-003 – Mock ID của user khác

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** HTTP 403/404.

### TC-USR-MOCK-ROOM-004 – Gửi câu trả lời text

- **Kết quả mong đợi:** Không chấp nhận rỗng/whitespace; lưu turn đúng; AI phản hồi và đánh giá.

### TC-USR-MOCK-ROOM-005 – Gửi câu trả lời voice

- **Kết quả mong đợi:** Xin quyền mic, thu âm/transcript, gửi dữ liệu đúng và hiển thị phản hồi.

### TC-USR-MOCK-ROOM-006 – Đang phân tích

- **Kết quả mong đợi:** Disable input/action gây request trùng; hiển thị loading.

### TC-USR-MOCK-ROOM-007 – Chuyển câu hỏi

- **Kết quả mong đợi:** Chỉ cho chuyển sau khi trả lời/đánh giá; reset input đúng.

### TC-USR-MOCK-ROOM-008 – Kết thúc đủ câu hỏi

- **Kết quả mong đợi:** Lưu transcript, gọi end API, tạo report/kết quả và chuyển lịch sử/chi tiết.

### TC-USR-MOCK-ROOM-009 – Kết thúc sớm

- **Kết quả mong đợi:** Có confirm; cancel giữ phiên; confirm kết thúc theo trạng thái xác định.

### TC-USR-MOCK-ROOM-010 – AI Service lỗi/timeout

- **Kết quả mong đợi:** Có retry/fallback; không mất câu trả lời đã gửi; không treo loading vô hạn.

### TC-USR-MOCK-ROOM-011 – Refresh khi đang phỏng vấn

- **Kết quả mong đợi:** Khôi phục phiên hoặc cảnh báo; không tạo phiên mới ngoài ý muốn.

### TC-USR-MOCK-ROOM-012 – Mất mạng

- **Kết quả mong đợi:** Hiển thị offline/reconnect; đồng bộ turns chưa lưu an toàn.

### TC-USR-MOCK-ROOM-013 – Nội dung XSS

- **Kết quả mong đợi:** Tin nhắn user/AI được escape; không thực thi script.

---

## 18. Test case lịch sử và kết quả luyện tập

### TC-USR-HISTORY-001 – Tải lịch sử

- **Route:** `/mock-results`
- **Kết quả mong đợi:** Chỉ phiên của user hiện tại; mới nhất trước; summary cards chính xác.

### TC-USR-HISTORY-002 – Không có lịch sử

- **Kết quả mong đợi:** Empty state và CTA luyện tập.

### TC-USR-HISTORY-003 – Hiển thị điểm

- **Kết quả mong đợi:** Điểm và badge màu đúng ngưỡng; không vượt thang điểm.

### TC-USR-HISTORY-004 – Xem chi tiết

- **Kết quả mong đợi:** Chuyển `/mock-results/{id}`, dữ liệu đúng phiên.

### TC-USR-HISTORY-005 – ID không tồn tại

- **Kết quả mong đợi:** 404/empty state; không hiển thị kết quả hardcoded khác.

### TC-USR-HISTORY-006 – ID của user khác

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** HTTP 403/404.

### TC-USR-HISTORY-007 – Nội dung chi tiết

- **Kết quả mong đợi:** Câu hỏi, câu trả lời, score, feedback, điểm tốt và cần cải thiện chính xác.

### TC-USR-HISTORY-008 – API lịch sử/kết quả lỗi

- **Kết quả mong đợi:** Loading kết thúc; có lỗi/thử lại.

---

## 19. Test case thông báo

### TC-USR-NOTIF-001 – Nhận thông báo ứng tuyển thành công

- **Kết quả mong đợi:** Có title/message/link đúng job.

### TC-USR-NOTIF-002 – Nhận thông báo lịch phỏng vấn

- **Kết quả mong đợi:** Đúng công ty, thời gian, interview link.

### TC-USR-NOTIF-003 – Nhận thông báo hủy lịch

- **Kết quả mong đợi:** Candidate biết lịch bị hủy và không thể tham gia.

### TC-USR-NOTIF-004 – Đánh dấu đã đọc

- **Kết quả mong đợi:** Badge unread giảm và trạng thái lưu bền vững.

### TC-USR-NOTIF-005 – Candidate chỉ xem thông báo của mình

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** API scope theo user ID trong token.

### TC-USR-NOTIF-006 – Link thông báo

- **Kết quả mong đợi:** Chuyển đúng route/resource và áp dụng RBAC.

---

## 20. Test case bảo mật

### TC-USR-SEC-001 – SQL Injection ở tìm kiếm

- **Dữ liệu:** `' OR 1=1 --`
- **Kết quả mong đợi:** Không lỗi SQL hoặc lộ dữ liệu.

### TC-USR-SEC-002 – Stored/reflected XSS

- **Vị trí:** Họ tên, CV filename, chat, ghi chú, câu trả lời mock.
- **Kết quả mong đợi:** Không thực thi JavaScript.

### TC-USR-SEC-003 – IDOR hồ sơ/đơn/lịch/mock result/file

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** User không thể truy cập/chỉnh dữ liệu user khác bằng thay ID.

### TC-USR-SEC-004 – Upload malware và MIME spoofing

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Kiểm tra magic bytes, MIME whitelist, kích thước; không thực thi file upload.

### TC-USR-SEC-005 – Path traversal filename

- **Kết quả mong đợi:** Server tạo storage key; không ghi ngoài thư mục upload.

### TC-USR-SEC-006 – Rate limit auth/apply/upload/AI

- **Kết quả mong đợi:** Chống brute-force, spam đơn, upload và lạm dụng AI.

### TC-USR-SEC-007 – Token trong URL/log

- **Kết quả mong đợi:** Không log access/refresh token hoặc token phòng đầy đủ; URL nhạy cảm có thời hạn.

### TC-USR-SEC-008 – CORS

- **Kết quả mong đợi:** Chỉ origin tin cậy; không wildcard với credential.

### TC-USR-SEC-009 – Response không chứa secret

- **Kết quả mong đợi:** Không password hash, API key, internal storage key hoặc dữ liệu user khác.

### TC-USR-SEC-010 – Quyền realtime room

- **Kết quả mong đợi:** Chỉ Candidate của interview được join/send chat/heartbeat.

---

## 24. Test case mở rộng – Luồng E2E Candidate

| ID | Kịch bản mở rộng | Dữ liệu/bước chính | Kết quả mong đợi | Ưu tiên |
|---|---|---|---|---|
| TC-USR-EXT-E2E-001 | Từ đăng ký đến ứng tuyển | Register → login → profile → upload CV → tìm Job → apply | Dữ liệu xuyên suốt nhất quán; đơn xuất hiện ở Candidate và Recruiter đúng công ty | Critical |
| TC-USR-EXT-E2E-002 | Cập nhật CV sau khi đã ứng tuyển | Upload CV mới từ Profile, Recruiter mở Candidate và parse | Candidate trỏ file thật; Recruiter được phép xem/parse; match được tính lại | Critical |
| TC-USR-EXT-E2E-003 | Từ lịch đến phòng phỏng vấn | Nhận lịch → consent → kiểm tra media → join → kết thúc | Vào đúng room, trạng thái và thời gian thống nhất, không lộ dữ liệu nội bộ | Critical |
| TC-USR-EXT-E2E-004 | Recruiter hủy lịch | Candidate đang xem lịch khi Recruiter hủy | UI cập nhật cancelled/notification; link không cho join trái phép | Critical |
| TC-USR-EXT-E2E-005 | Mock interview hoàn chỉnh | Setup → start → trả lời → end → history → result | Session và report gắn đúng user, kết quả hiển thị ổn định | Critical |
| TC-USR-EXT-E2E-006 | Rút đơn ứng tuyển | Apply → Recruiter thấy Candidate → Candidate rút đơn | Đơn biến mất/chuyển trạng thái theo contract ở cả hai phía, có thông báo phù hợp | High |

## 25. Test case mở rộng – Invite, waiting room và media

| ID | Kịch bản | Kết quả mong đợi | Ưu tiên |
|---|---|---|---|
| TC-USR-EXT-INV-001 | Invite token hợp lệ | Trả đúng metadata tối thiểu và room token có thời hạn | Critical |
| TC-USR-EXT-INV-002 | Invite token hết hạn | Chuyển `/interview-expired`, không cấp room token | Critical |
| TC-USR-EXT-INV-003 | Invite token bị sửa | Từ chối mà không tiết lộ Candidate/công ty/lịch | Critical |
| TC-USR-EXT-INV-004 | Token dùng trên hai thiết bị | Áp dụng đúng chính sách session; không tạo identity trùng gây chiếm phòng | High |
| TC-USR-EXT-INV-005 | Reload waiting room | Giữ dữ liệu/consent hợp lệ hoặc tải lại an toàn | High |
| TC-USR-EXT-MEDIA-001 | Từ chối microphone | Hiện hướng dẫn cấp quyền và cho thử lại, không crash | High |
| TC-USR-EXT-MEDIA-002 | Không có camera | Cho audio-only nếu mode hỗ trợ; báo rõ trạng thái | High |
| TC-USR-EXT-MEDIA-003 | Đổi thiết bị giữa phiên | Track cũ đóng, track mới hoạt động, không nhân đôi | Medium |
| TC-USR-EXT-MEDIA-004 | Mạng chuyển Wi-Fi/4G | Reconnect và khôi phục participant/chat phù hợp | High |
| TC-USR-EXT-MEDIA-005 | Recruiter chưa vào phòng | Candidate thấy waiting state, không tự bắt đầu interview | Medium |
| TC-USR-EXT-MEDIA-006 | Interview đã kết thúc khi đang reconnect | Không quay lại room; chuyển trạng thái hoàn thành | High |

## 26. Test case mở rộng – Mock interview, AI và phục hồi

| ID | Kịch bản | Kết quả mong đợi | Ưu tiên |
|---|---|---|---|
| TC-USR-EXT-MOCK-001 | Refresh khi đang mock | Khôi phục session hiện tại hoặc hướng dẫn resume an toàn | High |
| TC-USR-EXT-MOCK-002 | Gửi câu trả lời double-click | Chỉ lưu một message/answer logic | High |
| TC-USR-EXT-MOCK-003 | AI timeout giữa câu hỏi | Giữ câu trả lời, dừng loading và cho retry | High |
| TC-USR-EXT-MOCK-004 | Kết thúc khi còn request đang chạy | Chỉ một transition `completed`; không tạo report trùng | Critical |
| TC-USR-EXT-MOCK-005 | Report đang generating | History hiển thị pending; tự cập nhật hoặc cho reload | Medium |
| TC-USR-EXT-MOCK-006 | Report generation thất bại | Session vẫn tồn tại; lỗi rõ và có retry theo chính sách | High |
| TC-USR-EXT-MOCK-007 | Truy cập mock result user khác | HTTP 403/404, không lộ score/messages | Critical |
| TC-USR-EXT-MOCK-008 | Nội dung prompt injection | AI không thực thi chỉ dẫn trái quyền từ câu trả lời/CV | Critical |
| TC-USR-EXT-MOCK-009 | Transcript Unicode dài | Không lỗi encoding, có giới hạn payload hợp lý | Medium |

## 27. Test case mở rộng – Quyền riêng tư và dữ liệu cá nhân

| ID | Kịch bản | Kết quả mong đợi | Ưu tiên |
|---|---|---|---|
| TC-USR-EXT-PRIV-001 | Candidate đọc internal notes | Không có endpoint/UI nào trả ghi chú nội bộ Recruiter | Critical |
| TC-USR-EXT-PRIV-002 | Candidate đọc raw scoring/report nội bộ | Chỉ dữ liệu được phép công bố mới xuất hiện | Critical |
| TC-USR-EXT-PRIV-003 | Download CV bằng signed URL hết hạn | URL cũ bị từ chối; owner lấy URL mới được | High |
| TC-USR-EXT-PRIV-004 | Xóa tài khoản có đơn/lịch/mock | Có xác nhận và xử lý dữ liệu theo retention policy | Critical |
| TC-USR-EXT-PRIV-005 | Export dữ liệu cá nhân nếu hỗ trợ | Chỉ owner nhận dữ liệu, file có bảo vệ và thời hạn | High |
| TC-USR-EXT-PRIV-006 | Log frontend/backend | Không chứa access token, room token, CV content hoặc password | Critical |
| TC-USR-EXT-PRIV-007 | Candidate đổi email khi có phiên khác | Phiên/token và ownership được cập nhật an toàn | High |

## 28. Test case mở rộng – Accessibility, mobile và offline

| ID | Kịch bản | Kết quả mong đợi | Ưu tiên |
|---|---|---|---|
| TC-USR-EXT-A11Y-001 | Job board/Profile bằng bàn phím | Focus rõ, form và modal thao tác không cần chuột | High |
| TC-USR-EXT-A11Y-002 | Screen reader upload CV | Label, progress, thành công/lỗi được thông báo | Medium |
| TC-USR-EXT-A11Y-003 | Caption/transcript | Transcript có speaker/time rõ, không chỉ phân biệt bằng màu | Medium |
| TC-USR-EXT-A11Y-004 | Zoom 200% và mobile | CTA apply/join không bị che, không cuộn ngang nghiêm trọng | High |
| TC-USR-EXT-OFF-001 | Offline khi sửa Profile | Không mất dữ liệu form; báo offline và cho retry | High |
| TC-USR-EXT-OFF-002 | Mất mạng khi upload CV | Không tạo metadata nửa chừng; giữ CV cũ | Critical |
| TC-USR-EXT-OFF-003 | Back/forward sau apply | Không gửi lại POST hoặc tạo đơn trùng | High |
| TC-USR-EXT-OFF-004 | Storage local bị xóa | Saved Jobs reset an toàn; dữ liệu server không bị ảnh hưởng | Medium |

## 29. Checklist regression Candidate tối thiểu

- [ ] Register/login/logout/forgot password và token expiry.
- [ ] Profile và CV upload/parse với file hợp lệ, file giả và file quá lớn.
- [ ] Job Board, Job detail, AI match, Save Job và Apply.
- [ ] Danh sách đơn và rút đúng đơn của chính user.
- [ ] Lịch mới nhất, lịch hủy, invite hết hạn và consent.
- [ ] Room permissions, reconnect, chat/transcript và end state.
- [ ] Mock setup/start/messages/end/history/report.
- [ ] Candidate không xem dữ liệu nội bộ Recruiter hoặc dữ liệu user khác.
- [ ] Các case Critical về offline upload, IDOR và duplicate request đạt.

---

## 21. Test case hiệu năng và ổn định

### TC-USR-PERF-001 – Job board dữ liệu lớn

- **Kết quả mong đợi:** Phân trang/lazy load; thời gian phản hồi hợp lý; không treo UI.

### TC-USR-PERF-002 – Nhiều đơn ứng tuyển

- **Kết quả mong đợi:** Phân trang và thứ tự ổn định.

### TC-USR-PERF-003 – Nhiều lịch/mock history

- **Kết quả mong đợi:** Load nhanh, không duplicate, scroll ổn định.

### TC-USR-PERF-004 – Upload CV mạng chậm

- **Kết quả mong đợi:** Có progress/loading; timeout/retry hợp lý; không upload trùng.

### TC-USR-PERF-005 – Hai request apply đồng thời

- **Kết quả mong đợi:** Database chỉ có một application.

### TC-USR-PERF-006 – Realtime reconnect

- **Kết quả mong đợi:** Không participant/message trùng và không mất trạng thái chính.

### TC-USR-PERF-007 – AI timeout

- **Kết quả mong đợi:** UI không loading vô hạn; có retry; dữ liệu nghiệp vụ chính không bị rollback không cần thiết.

---

## 22. Ma trận API Candidate Portal

| STT | Method | Endpoint | Chức năng | Thành công |
|---:|---|---|---|---:|
| 1 | GET | `/portal/dashboard` | Dashboard Candidate | 200 |
| 2 | GET | `/portal/interviews` | Lịch phỏng vấn | 200 |
| 3 | GET | `/portal/profile` | Hồ sơ | 200 |
| 4 | PUT | `/portal/profile` | Cập nhật hồ sơ | 200 |
| 5 | POST | `/portal/cv` | Upload CV | 200 |
| 6 | POST | `/portal/jobs/{jobID}/apply` | Ứng tuyển | 200 |
| 7 | GET | `/portal/jobs/{jobID}/match` | AI CV-job match | 200 |
| 8 | GET | `/portal/applications` | Đơn của tôi | 200 |
| 9 | DELETE | `/portal/applications/{id}` | Rút đơn | 200 |

Quy tắc response chung:

| Trường hợp | HTTP mong đợi |
|---|---:|
| Không có token | 401 |
| Token sai/hết hạn | 401 |
| Sai role/quyền | 403 |
| Payload sai | 400 |
| Resource không tồn tại | 404 |
| Trùng/xung đột trạng thái | 409 |
| File quá lớn | 400/413 |
| Lỗi ngoài dự kiến | 500, không lộ stack trace |

## 23. Tiêu chí nghiệm thu

Actor User/Candidate đạt yêu cầu khi:

1. Đăng ký, đăng nhập, quên mật khẩu và đăng xuất hoạt động đúng.
2. Candidate không truy cập được Recruiter/Admin route hoặc API.
3. Dashboard và hồ sơ chỉ hiển thị dữ liệu của user hiện tại.
4. Candidate cập nhật hồ sơ và upload CV hợp lệ được.
5. CV có metadata/file vật lý đúng và recruiter có quyền mới xem được.
6. Candidate tìm kiếm, lọc, lưu và xem việc làm được.
7. Candidate ứng tuyển job mở và không thể tạo đơn trùng.
8. Candidate xem/rút đúng đơn của mình, không IDOR.
9. Lịch phỏng vấn mới nhất hiển thị trước; lịch hủy/completed không thể join.
10. Consent, room token, camera, microphone và realtime được kiểm soát đúng.
11. Mock interview tạo phiên, lưu transcript, kết thúc và hiển thị kết quả đúng.
12. Không có XSS, SQL injection, path traversal, upload malware hoặc lộ secret.
13. Token hết hạn/giả mạo và role localStorage bị sửa vẫn không vượt backend RBAC.
14. UI xử lý đúng API error, timeout, mất mạng và dữ liệu rỗng.
15. Dữ liệu lớn được phân trang và có hiệu năng chấp nhận được.

## 24. Điểm cần chú ý từ code hiện tại

1. `/home` và `/job-board` đang nằm trong danh sách public route; cần xác nhận đây có đúng yêu cầu bảo mật hay không.
2. Candidate portal API dùng user ID từ auth context, đây là thiết kế đúng để hạn chế IDOR và phải được kiểm thử.
3. Saved jobs hiện lưu trong localStorage nên có thể mất khi đổi thiết bị/trình duyệt và cần xử lý JSON hỏng.
4. Upload CV phải tạo một file record thật và gắn trực tiếp `cv_file_id`, không tạo metadata giả.
5. CV của Candidate upload ở portal không có company trực tiếp; recruiter chỉ được xem khi file gắn với Candidate thuộc công ty của họ.
6. AI parse CV là best-effort; upload có thể thành công dù AI lỗi, vì vậy UI phải phân biệt hai trạng thái.
7. Apply có thể thành công dù AI fit-score lỗi; fit score nên để null thay vì tạo điểm giả.
8. Route `/candidate-room` cần kiểm tra cả invite token/interview ownership, không chỉ access token.
9. Lịch Candidate phải loại bỏ nút tham gia khi `completed` hoặc `cancelled`.
10. GET match trả `has_cv=false` khi chưa có parsed CV; UI phải hiển thị CTA upload thay vì điểm ngẫu nhiên.
