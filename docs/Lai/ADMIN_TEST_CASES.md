# BỘ TEST CASE CHỨC NĂNG ACTOR ADMIN

## 1. Thông tin tài liệu

| Thuộc tính | Giá trị |
|---|---|
| Dự án | ViệcLàm AI – Nền tảng tuyển dụng và phỏng vấn AI |
| Actor | System Admin |
| Loại kiểm thử | Functional, API, RBAC, Security, UI/UX, Performance |
| Frontend Admin | `/admin/*` |
| API prefix | `/api/v1` |

## 2. Phạm vi chức năng

| Phân hệ | Frontend | API chính |
|---|---|---|
| Đăng nhập Admin | `/admin/login` | API đăng nhập dùng chung |
| Dashboard | `/admin/dashboard` | `GET /admin/dashboard-stats`, `GET /admin/reports`, `GET /admin/logs` |
| Người dùng | `/admin/users` | Danh sách, duyệt, khóa/mở khóa |
| Công ty | `/admin/companies` | `GET /admin/companies` |
| Báo cáo | `/admin/reports` | `GET /admin/reports` |
| Nhật ký | `/admin/logs` | `GET /admin/logs` |
| Cài đặt | `/admin/settings` | `GET/PUT /admin/settings` |
| Prompt AI | Tab trong Settings | `GET/POST /admin/ai-prompts` |

## 3. Dữ liệu kiểm thử

| Mã | Vai trò | Trạng thái | Mục đích |
|---|---|---|---|
| ACC-ADMIN-01 | Admin | Active | Tài khoản quản trị hợp lệ |
| ACC-ADMIN-02 | Admin | Blocked | Kiểm tra tài khoản Admin bị khóa |
| ACC-REC-01 | Recruiter | Active | Kiểm tra phân quyền |
| ACC-REC-02 | Recruiter | Pending | Kiểm tra duyệt tài khoản |
| ACC-REC-03 | Recruiter | Blocked | Kiểm tra mở khóa |
| ACC-CAN-01 | Candidate | Active | Kiểm tra phân quyền |
| ACC-CAN-02 | Candidate | Pending | Kiểm tra danh sách pending |

Tài khoản Admin mặc định được giao diện gợi ý:

```text
Email: admin@wemake.vn
Password: admin123
```

Dữ liệu công ty:

- COM-01: Công ty có đầy đủ tên, ngành nghề, website, quy mô.
- COM-02: Công ty không có website.
- COM-03: Công ty thiếu ngành nghề hoặc quy mô.
- COM-04: Website không chứa protocol.
- COM-05: Tên công ty có Unicode/Tiếng Việt.

---

## 4. Test case đăng nhập Admin

### TC-ADM-AUTH-001 – Hiển thị trang đăng nhập

- **Ưu tiên:** Critical
- **Tiền điều kiện:** Chưa đăng nhập.
- **Bước thực hiện:** Truy cập `/admin/login`.
- **Kết quả mong đợi:**
  - Hiển thị ô Email, Mật khẩu và nút đăng nhập.
  - Có nút quay lại trang đăng nhập người dùng.
  - Không hiển thị nội dung Dashboard.

### TC-ADM-AUTH-002 – Đăng nhập Admin thành công

- **Ưu tiên:** Critical
- **Tiền điều kiện:** ACC-ADMIN-01 tồn tại.
- **Bước thực hiện:**
  1. Nhập đúng email và mật khẩu.
  2. Bấm đăng nhập.
- **Kết quả mong đợi:**
  - API trả thành công.
  - Lưu `access_token` và `user_role=admin`.
  - Chuyển tới `/admin/dashboard`.

### TC-ADM-AUTH-003 – Bỏ trống thông tin

- **Ưu tiên:** High
- **Bước thực hiện:** Để trống email/mật khẩu và bấm đăng nhập.
- **Kết quả mong đợi:** Không gọi API; hiển thị yêu cầu nhập đầy đủ thông tin.

### TC-ADM-AUTH-004 – Email sai định dạng

- **Ưu tiên:** Medium
- **Dữ liệu:** `admin-wemake`
- **Kết quả mong đợi:** Trình duyệt hoặc ứng dụng báo email không hợp lệ; không đăng nhập.

### TC-ADM-AUTH-005 – Sai mật khẩu

- **Ưu tiên:** Critical
- **Kết quả mong đợi:**
  - Không tạo phiên đăng nhập.
  - Hiển thị “Email hoặc mật khẩu không chính xác”.
  - Không lưu token và không chuyển trang.

### TC-ADM-AUTH-006 – Email không tồn tại

- **Ưu tiên:** High
- **Kết quả mong đợi:** Đăng nhập thất bại và không tiết lộ email có tồn tại hay không.

### TC-ADM-AUTH-007 – Recruiter đăng nhập tại cổng Admin

- **Ưu tiên:** Critical
- **Kết quả mong đợi:**
  - Frontend logout sau khi phát hiện role không phải Admin.
  - Hiển thị thông báo không có quyền truy cập Admin Portal.
  - Không giữ phiên Recruiter tại cổng Admin.

### TC-ADM-AUTH-008 – Candidate đăng nhập tại cổng Admin

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Candidate bị từ chối, phiên bị xóa và không vào Dashboard.

### TC-ADM-AUTH-009 – Admin đã đăng nhập truy cập lại trang login

- **Ưu tiên:** Medium
- **Kết quả mong đợi:** Tự động chuyển tới `/admin/dashboard`.

### TC-ADM-AUTH-010 – Admin bị khóa đăng nhập

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Backend không cấp token; hiển thị thông báo tài khoản bị khóa/không hoạt động.

### TC-ADM-AUTH-011 – Bấm đăng nhập liên tục

- **Ưu tiên:** Medium
- **Kết quả mong đợi:** Nút bị disable khi loading và chỉ gửi một request hợp lệ.

### TC-ADM-AUTH-012 – API đăng nhập lỗi hoặc mất mạng

- **Ưu tiên:** High
- **Kết quả mong đợi:** Loading kết thúc; có thông báo lỗi; không lưu token.

---

## 5. Test case xác thực và phân quyền

### TC-ADM-RBAC-001 – Truy cập Dashboard khi chưa đăng nhập

- **Ưu tiên:** Critical
- **Bước thực hiện:** Xóa token rồi truy cập `/admin/dashboard`.
- **Kết quả mong đợi:** Chuyển tới `/admin/login` kèm thông báo yêu cầu đăng nhập Admin.

### TC-ADM-RBAC-002 – Recruiter truy cập route Admin

- **Ưu tiên:** Critical
- **Bước thực hiện:** Đăng nhập Recruiter rồi truy cập `/admin/users`.
- **Kết quả mong đợi:** Chuyển tới `/403?reason=admin_required`; không tải dữ liệu Admin.

### TC-ADM-RBAC-003 – Candidate truy cập route Admin

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Chuyển tới `/403`; không hiển thị dữ liệu quản trị.

### TC-ADM-RBAC-004 – API Admin không có token

- **API:** `GET /admin/users`
- **Ưu tiên:** Critical
- **Kết quả mong đợi:** HTTP 401; không trả dữ liệu người dùng.

### TC-ADM-RBAC-005 – API Admin bằng token Candidate

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** HTTP 403.

### TC-ADM-RBAC-006 – API Admin bằng token Recruiter

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** HTTP 403.

### TC-ADM-RBAC-007 – Token hết hạn

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** API trả 401; frontend xóa phiên và chuyển về login.

### TC-ADM-RBAC-008 – Token bị chỉnh sửa

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Backend từ chối chữ ký; không tin riêng `user_role` trong localStorage.

### TC-ADM-RBAC-009 – Sửa `user_role=admin` nhưng dùng token Candidate

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Backend vẫn trả 403; không lộ dữ liệu Admin.

### TC-ADM-RBAC-010 – Admin truy cập đầy đủ route

- **Route:** `/admin/dashboard`, `/admin/users`, `/admin/companies`, `/admin/reports`, `/admin/logs`, `/admin/settings`.
- **Kết quả mong đợi:** Tất cả tải đúng component, không bị chuyển `/403`.

---

## 6. Test case Dashboard

### TC-ADM-DASH-001 – Tải Dashboard thành công

- **API:** `GET /admin/dashboard-stats`, `/admin/reports`, `/admin/logs`, `/admin/users/pending`.
- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Tất cả request kèm token; Dashboard hiển thị không có lỗi JavaScript.

### TC-ADM-DASH-002 – Số liệu tổng quan chính xác

- **Ưu tiên:** High
- **Kết quả mong đợi:** Tổng user, công ty, phỏng vấn và pending khớp database; không có `NaN`, `undefined`, số âm.

### TC-ADM-DASH-003 – Không có dữ liệu thống kê

- **Kết quả mong đợi:** Hiển thị 0; bố cục không lỗi; không dùng số giả.

### TC-ADM-DASH-004 – Một API Dashboard lỗi

- **Ưu tiên:** High
- **Kết quả mong đợi:** Trang không trắng; loading kết thúc; các phần còn lại vẫn hoạt động.

### TC-ADM-DASH-005 – Danh sách pending chính xác

- **Kết quả mong đợi:** Chỉ hiển thị user pending, không có active/blocked.

### TC-ADM-DASH-006 – Duyệt user từ Dashboard

- **Ưu tiên:** Critical
- **Kết quả mong đợi:**
  - Gọi đúng API/user ID.
  - User thành active và biến mất khỏi pending.
  - Số pending giảm.
  - Có audit log action `APPROVE`.

### TC-ADM-DASH-007 – Hủy xác nhận duyệt

- **Kết quả mong đợi:** Không gọi API; trạng thái không đổi.

### TC-ADM-DASH-008 – Điều hướng nhanh

- **Kết quả mong đợi:** Các nút chuyển đúng `/admin/users`, `/admin/companies`, `/admin/logs`.

---

## 7. Test case quản lý người dùng

### TC-ADM-USER-001 – Tải danh sách

- **API:** `GET /admin/users`, `GET /admin/users/pending`.
- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Hai danh sách và số đếm tab chính xác.

### TC-ADM-USER-002 – Chuyển tab

- **Kết quả mong đợi:** Nội dung đổi đúng; bộ lọc role chỉ có ở tab tất cả.

### TC-ADM-USER-003 – Tìm theo tên

- **Kết quả mong đợi:** Lọc đúng, không phân biệt hoa/thường.

### TC-ADM-USER-004 – Tìm theo email

- **Kết quả mong đợi:** Lọc đúng toàn bộ hoặc một phần email.

### TC-ADM-USER-005 – Không có kết quả

- **Kết quả mong đợi:** Hiển thị empty state; xóa từ khóa khôi phục danh sách.

### TC-ADM-USER-006 – Lọc theo vai trò

- **Dữ liệu:** Recruiter, Candidate, Admin.
- **Kết quả mong đợi:** Chỉ hiển thị đúng role được chọn.

### TC-ADM-USER-007 – Kết hợp tìm kiếm và role

- **Kết quả mong đợi:** Kết quả đồng thời khớp cả hai điều kiện.

### TC-ADM-USER-008 – Làm mới danh sách

- **Kết quả mong đợi:** Gọi lại API; icon quay; nút disable khi loading.

### TC-ADM-USER-009 – Duyệt tài khoản pending

- **API:** `PUT /admin/users/{user_id}/approve`
- **Ưu tiên:** Critical
- **Kết quả mong đợi:**
  - Có xác nhận.
  - API 200; user chuyển active.
  - Danh sách tải lại.
  - Audit log đúng actor, action `APPROVE`, resource ID.

### TC-ADM-USER-010 – Hủy duyệt

- **Kết quả mong đợi:** Không gọi API; user vẫn pending.

### TC-ADM-USER-011 – Duyệt ID không tồn tại

- **Kết quả mong đợi:** 404 hoặc lỗi nghiệp vụ phù hợp; frontend không báo thành công.

### TC-ADM-USER-012 – Duyệt lại user active

- **Kết quả mong đợi:** Xử lý idempotent hoặc trả lỗi trạng thái rõ ràng.

### TC-ADM-USER-013 – Khóa user active

- **API:** `PUT /admin/users/{id}/status`
- **Payload:** `{"status":"blocked"}`
- **Ưu tiên:** Critical
- **Kết quả mong đợi:** User thành blocked; UI đổi thành “Mở khóa”; log action `BLOCK`.

### TC-ADM-USER-014 – Mở khóa user

- **Payload:** `{"status":"active"}`
- **Kết quả mong đợi:** User thành active; log action `ACTIVATE`; có thể đăng nhập lại.

### TC-ADM-USER-015 – Không cho khóa Admin từ UI

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Không hiển thị nút khóa/mở khóa với `role=admin`.

### TC-ADM-USER-016 – Tự khóa Admin bằng API

- **Ưu tiên:** Critical
- **Kết quả mong đợi đề xuất:** Backend từ chối tự khóa hoặc khóa Admin active cuối cùng.

### TC-ADM-USER-017 – Status không hợp lệ

- **Payload:** `{"status":"deleted_by_test"}`
- **Kết quả mong đợi:** HTTP 400; database không đổi.

### TC-ADM-USER-018 – Payload status rỗng

- **Kết quả mong đợi:** HTTP 400 với thông báo thiếu status.

### TC-ADM-USER-019 – Xem giấy tờ xác thực

- **Tiền điều kiện:** User pending có `verification_file_id`.
- **Kết quả mong đợi:** File đúng mở ở tab mới; chỉ Admin được phép truy cập.

### TC-ADM-USER-020 – Không có giấy tờ

- **Kết quả mong đợi:** Hiển thị “Chưa tải file”, không có link hỏng.

### TC-ADM-USER-021 – API danh sách lỗi

- **Kết quả mong đợi:** Loading kết thúc; trang không treo; có lỗi/thử lại.

---

## 8. Test case quản lý công ty

### TC-ADM-COMP-001 – Tải danh sách công ty

- **API:** `GET /admin/companies`
- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Dữ liệu và tổng số công ty chính xác.

### TC-ADM-COMP-002 – Hiển thị thông tin

- **Kết quả mong đợi:** Tên, ngành, website, quy mô và trạng thái khớp API; field thiếu có fallback hợp lý.

### TC-ADM-COMP-003 – Tìm theo tên/ngành/website

- **Kết quả mong đợi:** Lọc đúng, không phân biệt hoa/thường.

### TC-ADM-COMP-004 – Công ty không có website

- **Kết quả mong đợi:** Hiển thị “Chưa có website”; không tạo link rỗng.

### TC-ADM-COMP-005 – Website không có protocol

- **Dữ liệu:** `example.com`
- **Kết quả mong đợi:** Mở `https://example.com`.

### TC-ADM-COMP-006 – Website đã có protocol

- **Kết quả mong đợi:** Không thêm `https://` lần hai.

### TC-ADM-COMP-007 – Không có kết quả tìm kiếm

- **Kết quả mong đợi:** Empty state phù hợp; tổng số 0 theo kết quả lọc.

### TC-ADM-COMP-008 – Không có công ty trong hệ thống

- **Kết quả mong đợi:** Hiển thị “Chưa có công ty nào”.

### TC-ADM-COMP-009 – Làm mới danh sách

- **Kết quả mong đợi:** Gọi lại API; loading đúng; chống request lặp.

### TC-ADM-COMP-010 – Nút “Xem hồ sơ”

- **Ưu tiên:** High
- **Kết quả mong đợi:** Điều hướng đến chi tiết nếu hỗ trợ; nếu chưa có phải disable/ẩn hoặc báo đang phát triển, không để nút không phản hồi.

### TC-ADM-COMP-011 – API công ty lỗi

- **Kết quả mong đợi:** Loading kết thúc; có thông báo và tùy chọn thử lại.

---

## 9. Test case báo cáo hệ thống

### TC-ADM-REPORT-001 – Tải báo cáo

- **API:** `GET /admin/reports`
- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Dữ liệu hiển thị đúng, không có `NaN`/`undefined`.

### TC-ADM-REPORT-002 – Số liệu khớp database

- **Kết quả mong đợi:** Tổng phỏng vấn, hoàn thành, điểm trung bình và thống kê theo kỳ chính xác.

### TC-ADM-REPORT-003 – Không có báo cáo

- **Kết quả mong đợi:** Hiển thị 0/empty state, biểu đồ không vỡ.

### TC-ADM-REPORT-004 – Phần trăm và số thập phân

- **Kết quả mong đợi:** Làm tròn nhất quán; phần trăm 0–100; điểm không vượt giới hạn.

### TC-ADM-REPORT-005 – API báo cáo lỗi

- **Kết quả mong đợi:** Trang không trắng; có trạng thái lỗi và thử lại.

### TC-ADM-REPORT-006 – Non-admin gọi API

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** HTTP 403; không trả số liệu toàn hệ thống.

---

## 10. Test case nhật ký hệ thống

### TC-ADM-LOG-001 – Tải danh sách log

- **API:** `GET /admin/logs`
- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Log mới nhất trước; có thời gian, actor, action, resource.

### TC-ADM-LOG-002 – Log duyệt user

- **Kết quả mong đợi:** Có `APPROVE`, đúng Admin actor và user resource.

### TC-ADM-LOG-003 – Log khóa/mở khóa

- **Kết quả mong đợi:** Có `BLOCK`/`ACTIVATE`, IP/User-Agent nếu hỗ trợ.

### TC-ADM-LOG-004 – Tìm kiếm và lọc log

- **Kết quả mong đợi:** Lọc đúng theo từ khóa, action, actor hoặc resource mà UI hỗ trợ.

### TC-ADM-LOG-005 – Phân trang

- **Kết quả mong đợi:** Không trùng/thiếu log, thứ tự ổn định, trang vượt giới hạn trả rỗng.

### TC-ADM-LOG-006 – Không có log

- **Kết quả mong đợi:** Hiển thị empty state.

### TC-ADM-LOG-007 – Dữ liệu nhạy cảm

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Log không có password, token, API key, CV đầy đủ hoặc secret.

### TC-ADM-LOG-008 – Non-admin truy cập log

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** HTTP 403.

### TC-ADM-LOG-009 – API log lỗi

- **Kết quả mong đợi:** Trang không treo; có trạng thái lỗi phù hợp.

---

## 11. Test case cài đặt hệ thống

### TC-ADM-SET-001 – Tải cài đặt

- **API:** `GET /admin/settings`
- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Form nhận đúng dữ liệu backend; không bị giá trị mặc định ghi đè.

### TC-ADM-SET-002 – Cập nhật thương hiệu

- **Dữ liệu:** `system_name`, `brand_name`, `brand_badge`, `brand_slogan`, `brand_logo_url`.
- **Kết quả mong đợi:** PUT thành công; store và layout cập nhật; reload vẫn giữ dữ liệu.

### TC-ADM-SET-003 – Email hỗ trợ

- **Kết quả mong đợi:** Email hợp lệ lưu được; email sai bị từ chối.

### TC-ADM-SET-004 – Bật/tắt maintenance mode

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** User thường bị giới hạn; Admin vẫn truy cập để tắt; tắt xong hệ thống trở lại bình thường.

### TC-ADM-SET-005 – Kích thước upload

- **Kết quả mong đợi:** Số dương hợp lệ được lưu và backend thực thi; 0/âm bị 400.

### TC-ADM-SET-006 – Điểm đạt mặc định

- **Kết quả mong đợi:** Chấp nhận 0–100; từ chối -1 và 101.

### TC-ADM-SET-007 – AI model mặc định

- **Kết quả mong đợi:** Chỉ model hỗ trợ được lưu và sử dụng cho request mới.

### TC-ADM-SET-008 – JWT expiry

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Chỉ số dương trong giới hạn an toàn; token mới dùng expiry mới.

### TC-ADM-SET-009 – Bắt buộc 2FA Admin

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Nếu bật, login tiếp theo yêu cầu 2FA; không cho bật nếu 2FA chưa được triển khai gây khóa Admin.

### TC-ADM-SET-010 – TTL và giới hạn thông báo

- **Kết quả mong đợi:** Chỉ nhận số nguyên dương trong giới hạn backend.

### TC-ADM-SET-011 – Cấu hình email notification

- **Kết quả mong đợi:** Cờ tổng và từng sự kiện được lưu độc lập; luồng gửi email tuân theo cấu hình.

### TC-ADM-SET-012 – Lỗi khi lưu

- **Ưu tiên:** High
- **Kết quả mong đợi:** Không báo thành công; loading kết thúc; platform store phải rollback/đồng bộ lại dữ liệu backend.

### TC-ADM-SET-013 – Field không hỗ trợ

- **Kết quả mong đợi:** Backend bỏ qua hoặc từ chối field ngoài whitelist.

### TC-ADM-SET-014 – Non-admin cập nhật settings

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** HTTP 403; cấu hình không đổi.

---

## 12. Test case Prompt AI

### TC-ADM-PROMPT-001 – Tải danh sách Prompt

- **API:** `GET /admin/ai-prompts`
- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Hiển thị đúng tên, model, nội dung/metadata.

### TC-ADM-PROMPT-002 – Mở modal tạo mới

- **Kết quả mong đợi:** Có trường tên, model, nội dung và giá trị mặc định.

### TC-ADM-PROMPT-003 – Tạo Prompt hợp lệ

- **API:** `POST /admin/ai-prompts`
- **Payload:**

```json
{
  "name": "RUBRIC_EVALUATION_V2",
  "model": "gemini-2.5-flash",
  "content": "Bạn là chuyên gia nhân sự AI..."
}
```

- **Kết quả mong đợi:** Tạo thành công; modal đóng; danh sách tải lại; có thông báo.

### TC-ADM-PROMPT-004 – Thiếu tên/nội dung

- **Kết quả mong đợi:** Không gửi API; hiển thị validation.

### TC-ADM-PROMPT-005 – Model không hợp lệ

- **Kết quả mong đợi:** Backend từ chối hoặc chỉ chấp nhận model trong whitelist.

### TC-ADM-PROMPT-006 – Nội dung rất dài

- **Kết quả mong đợi:** Áp dụng giới hạn; không treo frontend/backend; lỗi rõ nếu vượt.

### TC-ADM-PROMPT-007 – Unicode và xuống dòng

- **Kết quả mong đợi:** Lưu và hiển thị nguyên vẹn.

### TC-ADM-PROMPT-008 – Stored XSS

- **Dữ liệu:** `<script>alert(1)</script>`
- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Không thực thi script; dữ liệu được escape/sanitize.

### TC-ADM-PROMPT-009 – Chỉnh sửa Prompt

- **Kết quả mong đợi:** Modal có dữ liệu cũ; lưu thành phiên bản/cập nhật đúng thiết kế.

### TC-ADM-PROMPT-010 – API lưu lỗi

- **Kết quả mong đợi:** Modal không tự đóng; hiển thị lỗi; không báo thành công.

### TC-ADM-PROMPT-011 – Non-admin tạo Prompt

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** HTTP 403; không tạo dữ liệu.

---

## 13. Test case điều hướng và đăng xuất

### TC-ADM-NAV-001 – Sidebar đầy đủ

- **Kết quả mong đợi:** Có Dashboard, Người dùng, Công ty, Báo cáo, Nhật ký, Cài đặt.

### TC-ADM-NAV-002 – Active menu

- **Kết quả mong đợi:** Chỉ menu của route hiện tại được highlight.

### TC-ADM-NAV-003 – Truy cập `/admin`

- **Kết quả mong đợi:** Redirect `/admin/dashboard`.

### TC-ADM-NAV-004 – Refresh trang con

- **Ưu tiên:** Critical
- **Bước thực hiện:** F5 tại `/admin/users`.
- **Kết quả mong đợi:** Phiên hợp lệ được giữ và đúng trang được tải lại.

### TC-ADM-NAV-005 – Đăng xuất

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Xóa/revoke token, role và dữ liệu phiên; chuyển về Admin login.

### TC-ADM-NAV-006 – Back sau logout

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Không xem lại dữ liệu Admin từ cache; chuyển về login.

### TC-ADM-NAV-007 – Route Admin không tồn tại

- **Kết quả mong đợi:** Hiển thị 404 hoặc redirect Dashboard; không chuyển nhầm Candidate portal.

### TC-ADM-NAV-008 – Responsive

- **Thiết bị:** Desktop, tablet, mobile.
- **Kết quả mong đợi:** Sidebar/header, bảng và modal sử dụng được; không mất action.

---

## 14. Test case bảo mật

### TC-ADM-SEC-001 – SQL Injection

- **Dữ liệu:** `' OR 1=1 --`
- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Không lỗi SQL/không lộ dữ liệu; backend dùng parameterized query.

### TC-ADM-SEC-002 – XSS tên user/công ty

- **Dữ liệu:** `<img src=x onerror=alert(1)>`
- **Kết quả mong đợi:** Hiển thị text; không thực thi JavaScript.

### TC-ADM-SEC-003 – IDOR cập nhật user

- **Bước thực hiện:** Dùng token non-admin và thay `user_id` trong API.
- **Kết quả mong đợi:** HTTP 403; dữ liệu không đổi.

### TC-ADM-SEC-004 – UUID không hợp lệ

- **API:** `/admin/users/not-a-uuid/status`
- **Kết quả mong đợi:** HTTP 400; không lộ lỗi database/stack trace.

### TC-ADM-SEC-005 – Rate limit đăng nhập

- **Ưu tiên:** Critical
- **Kết quả mong đợi:** Brute-force bị hạn chế bằng rate limit/lockout phù hợp.

### TC-ADM-SEC-006 – CORS

- **Kết quả mong đợi:** Chỉ origin được cấu hình; không wildcard với credential.

### TC-ADM-SEC-007 – Response không lộ secret

- **Kết quả mong đợi:** Không trả password hash, refresh token, API key hoặc storage key nội bộ.

### TC-ADM-SEC-008 – Audit hành động quản trị

- **Hành động:** Duyệt, khóa, mở khóa, sửa settings, tạo Prompt.
- **Kết quả mong đợi:** Log có actor, action, resource, thời gian, IP/User-Agent.

### TC-ADM-SEC-009 – Không sửa/xóa audit log

- **Kết quả mong đợi:** Không có API thông thường cho Admin sửa/xóa audit log.

### TC-ADM-SEC-010 – Bảo vệ file xác thực

- **Kết quả mong đợi:** Chỉ Admin được xem; URL yêu cầu token hoặc có thời hạn; không đoán ID để tải file khác.

---

## 15. Test case hiệu năng và đồng thời

### TC-ADM-PERF-001 – 10.000 người dùng

- **Kết quả mong đợi:** Backend phân trang; UI không tải toàn bộ; thời gian phản hồi đạt yêu cầu.

### TC-ADM-PERF-002 – Nhiều công ty

- **Kết quả mong đợi:** Không treo trình duyệt; có pagination/lazy load nếu dữ liệu lớn.

### TC-ADM-PERF-003 – Nhiều audit log

- **Kết quả mong đợi:** Phân trang ổn định; query dùng index; không tải toàn bộ.

### TC-ADM-PERF-004 – Hai Admin cùng duyệt một user

- **Kết quả mong đợi:** Trạng thái cuối active; không dữ liệu trùng; xử lý idempotent/conflict rõ ràng.

### TC-ADM-PERF-005 – Hai Admin cập nhật settings

- **Kết quả mong đợi:** Quy tắc đồng thời rõ ràng; không tạo cấu hình hỏng.

### TC-ADM-PERF-006 – Mất mạng khi lưu

- **Kết quả mong đợi:** Timeout hợp lý; loading kết thúc; báo lỗi; không báo thành công giả.

---

## 16. Ma trận API Admin

| STT | Method | Endpoint | Chức năng | Role | Thành công |
|---:|---|---|---|---|---:|
| 1 | POST | `/auth/login` | Đăng nhập | Public | 200 |
| 2 | GET | `/admin/dashboard-stats` | Dashboard | Admin | 200 |
| 3 | GET | `/admin/users` | Danh sách user | Admin | 200 |
| 4 | GET | `/admin/users/pending` | User pending | Admin | 200 |
| 5 | PUT | `/admin/users/{id}/approve` | Duyệt user | Admin | 200 |
| 6 | PUT | `/admin/users/{id}/status` | Khóa/mở khóa | Admin | 200 |
| 7 | GET | `/admin/companies` | Danh sách công ty | Admin | 200 |
| 8 | GET | `/admin/reports` | Báo cáo hệ thống | Admin | 200 |
| 9 | GET | `/admin/logs` | Audit log | Admin | 200 |
| 10 | GET | `/admin/settings` | Lấy cấu hình | Admin | 200 |
| 11 | PUT | `/admin/settings` | Cập nhật cấu hình | Admin | 200 |
| 12 | GET | `/admin/ai-prompts` | Danh sách Prompt | Admin | 200 |
| 13 | POST | `/admin/ai-prompts` | Tạo Prompt/version | Admin | 200/201 |

Quy tắc response chung:

| Trường hợp | HTTP mong đợi |
|---|---:|
| Không có token | 401 |
| Token sai/hết hạn | 401 |
| Candidate/Recruiter | 403 |
| Admin hợp lệ | 200/201 |
| Payload sai | 400 |
| Resource không tồn tại | 404 |
| Xung đột trạng thái | 409 |
| Lỗi ngoài dự kiến | 500, không lộ stack trace |

## 17. Tiêu chí nghiệm thu

Actor Admin đạt yêu cầu khi:

1. Admin hợp lệ đăng nhập được và non-admin bị từ chối.
2. Route và API Admin đều có RBAC phía backend.
3. Dashboard hiển thị số liệu đúng database.
4. Admin tìm kiếm, lọc, duyệt, khóa và mở khóa user được.
5. Không thể khóa Admin từ UI; backend bảo vệ Admin hiện tại/Admin cuối cùng.
6. Mọi hành động quan trọng tạo audit log đầy đủ.
7. Danh sách công ty, báo cáo và log hoạt động đúng.
8. Settings được lưu bền vững và thực sự tác động đến hệ thống.
9. Prompt AI được version/quản lý an toàn.
10. Không có XSS, SQL injection, IDOR hoặc lộ secret.
11. Token giả/hết hạn bị backend từ chối.
12. UI xử lý đúng khi API lỗi, timeout hoặc mất mạng.
13. Dữ liệu lớn được phân trang và có hiệu năng chấp nhận được.

## 18. Điểm cần chú ý từ code hiện tại

1. Frontend có route guard chặn non-admin, nhưng backend vẫn phải là lớp bảo vệ quyết định vì localStorage có thể sửa.
2. Backend đã gắn AuthMiddleware và RoleMiddleware cho các route trong UserHandler.
3. UI không hiển thị nút khóa/mở khóa với Admin, nhưng backend cũng nên cấm tự khóa hoặc khóa Admin cuối cùng.
4. Nút “Xem hồ sơ” tại trang công ty hiện cần xác minh vì chưa thấy hành động điều hướng.
5. Link giấy phép đang dùng `/files/{id}/download`; cần xác nhận backend có endpoint tương ứng, nếu không sẽ 404.
6. Trang settings cập nhật platform store trước khi API PUT hoàn tất; khi API lỗi cần rollback/đồng bộ lại.
7. Maintenance mode, JWT expiry và Admin 2FA phải được kiểm tra xem backend có thực thi hay mới chỉ lưu cấu hình.

---

## 19. Test case mở rộng – E2E, contract và tính toàn vẹn

| ID | Kịch bản mở rộng | Dữ liệu/bước chính | Kết quả mong đợi | Ưu tiên |
|---|---|---|---|---|
| TC-ADM-EXT-E2E-001 | Luồng duyệt Recruiter hoàn chỉnh | Login Admin → Pending Users → xem giấy phép → Approve → Recruiter login | Trạng thái chuyển `pending` sang `active`, Recruiter đăng nhập được và có audit log | Critical |
| TC-ADM-EXT-E2E-002 | Luồng khóa và mở khóa tài khoản | Khóa user active, thử dùng access/refresh token cũ, sau đó mở khóa | Phiên cũ bị từ chối theo chính sách; user chỉ đăng nhập lại sau khi được mở | Critical |
| TC-ADM-EXT-E2E-003 | Cấu hình hệ thống tác động thực tế | Thay upload limit/passing score rồi chạy luồng Candidate tương ứng | Giá trị mới được backend thực thi, không chỉ thay đổi giao diện | Critical |
| TC-ADM-EXT-E2E-004 | Maintenance mode | Bật maintenance, truy cập bằng Candidate/Recruiter/Admin rồi tắt | User thường bị chặn rõ ràng; Admin có đường phục hồi; tắt xong hệ thống hoạt động lại | Critical |
| TC-ADM-EXT-E2E-005 | Prompt AI version mới | Tạo prompt mới, chạy tác vụ AI, kiểm tra log/version | Tác vụ mới dùng đúng version; dữ liệu cũ vẫn truy vết được | High |
| TC-ADM-EXT-API-001 | Content-Type sai | Gửi `text/plain` tới endpoint PUT/POST Admin | Trả 400/415, không thay đổi dữ liệu | High |
| TC-ADM-EXT-API-002 | JSON lỗi cú pháp | Gửi body thiếu dấu ngoặc | Trả 400, không panic và không lộ parser stack | High |
| TC-ADM-EXT-API-003 | Field thừa/mass assignment | Thêm `role`, `is_super_admin`, `password_hash` vào payload settings/status | Field trái phép bị bỏ qua hoặc từ chối | Critical |
| TC-ADM-EXT-API-004 | Page/page_size biên | 0, âm, cực lớn, chuỗi, overflow | Chuẩn hóa hoặc từ chối; không full-scan ngoài ý muốn | High |
| TC-ADM-EXT-API-005 | Sort field tùy ý | Gửi sort chứa tên cột/SQL fragment | Chỉ whitelist field hợp lệ; không SQL injection | Critical |
| TC-ADM-EXT-API-006 | Request ID và lỗi chuẩn | Gây 400/403/404/409/500 có kiểm soát | Response theo schema chung và có request ID để truy vết | Medium |
| TC-ADM-EXT-DATA-001 | Approve tạo đúng một audit log | Gửi approve một lần và lặp lại | Không tạo chuyển trạng thái sai; log phản ánh before/after chính xác | Critical |
| TC-ADM-EXT-DATA-002 | Không log secret | Login, đổi settings, tạo prompt | Log không chứa password, token, OTP, API key hoặc prompt secret ngoài chính sách | Critical |
| TC-ADM-EXT-DATA-003 | Thống kê sau khóa/xóa mềm | Khóa user hoặc công ty theo chính sách | Dashboard/report nhất quán với quy tắc đếm đã công bố | High |
| TC-ADM-EXT-DATA-004 | Thời gian audit | Thực hiện action qua client ở timezone khác | Timestamp lưu UTC/chuẩn thống nhất, hiển thị đúng locale | Medium |
| TC-ADM-EXT-CONC-001 | Hai Admin duyệt cùng user | Gửi hai request approve đồng thời | Chỉ một transition hợp lệ; request còn lại idempotent hoặc 409 | Critical |
| TC-ADM-EXT-CONC-002 | Khóa/mở khóa đồng thời | Hai Admin gửi trạng thái trái ngược | Trạng thái cuối xác định được, không corrupt dữ liệu | Critical |
| TC-ADM-EXT-CONC-003 | Lưu settings đồng thời | Hai tab sửa các field khác nhau | Không mất cập nhật âm thầm; có versioning/last-write policy rõ | High |
| TC-ADM-EXT-CONC-004 | Tạo Prompt double-click | Gửi hai POST giống nhau | Không tạo version trùng ngoài ý muốn | High |

## 20. Test case mở rộng – Accessibility, responsive và phục hồi

| ID | Kịch bản | Kết quả mong đợi | Ưu tiên |
|---|---|---|---|
| TC-ADM-EXT-A11Y-001 | Điều hướng toàn bộ bằng bàn phím | Tab order hợp lý, focus visible, modal giữ/trả focus đúng | High |
| TC-ADM-EXT-A11Y-002 | Screen reader cho bảng và form | Label, heading, table header và lỗi validation được đọc rõ | Medium |
| TC-ADM-EXT-A11Y-003 | Tương phản màu và trạng thái | Đạt mức tương phản phù hợp; trạng thái không chỉ truyền đạt bằng màu | Medium |
| TC-ADM-EXT-A11Y-004 | Zoom 200% | Không mất nút hành động hoặc tràn nội dung nghiêm trọng | Medium |
| TC-ADM-EXT-UX-001 | Mobile/tablet Admin | Bảng cuộn/thu gọn hợp lý, modal và menu dùng được | Medium |
| TC-ADM-EXT-UX-002 | Back/forward browser | Filter, tab và phân trang không rơi vào trạng thái sai | Medium |
| TC-ADM-EXT-REC-001 | API mất kết nối khi đang lưu | Không báo thành công giả; form giữ dữ liệu và cho retry | High |
| TC-ADM-EXT-REC-002 | Refresh giữa thao tác | Không lặp approve/lock/create prompt; tải lại trạng thái từ backend | High |
| TC-ADM-EXT-REC-003 | Redis/AI service ngắt | Chức năng quản trị cốt lõi còn hoạt động; phần phụ báo degraded | High |
| TC-ADM-EXT-REC-004 | Database tạm thời lỗi | Trả 5xx an toàn, không ghi dữ liệu nửa chừng; phục hồi không nhân đôi | Critical |

## 21. Checklist regression Admin tối thiểu

- [ ] Login Admin, logout và token hết hạn.
- [ ] Candidate/Recruiter bị chặn ở cả route và API Admin.
- [ ] Dashboard totals đối chiếu được với dữ liệu nguồn.
- [ ] Approve, lock, unlock và bảo vệ Admin cuối cùng.
- [ ] Company list, report và audit log có phân trang.
- [ ] Settings lưu bền vững và rollback UI khi lỗi.
- [ ] Prompt AI có validation/version/audit.
- [ ] Không lộ token, password, OTP, secret hoặc stack trace.
- [ ] Các test Critical về concurrency và maintenance mode đạt.

