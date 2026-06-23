# MODULE_LAI_PRODUCT_UI_BUSINESS.md
# Phân công Module cho Lai

## 1. Vai trò của Lai

Lai phụ trách nhóm module **còn lại nhưng rất quan trọng về sản phẩm, giao diện, trải nghiệm người dùng và các luồng nghiệp vụ CRUD**.

Lai là owner chính của các module:

1. UI Foundation.
2. Auth UI.
3. Recruiter Dashboard.
4. Job Management UI.
5. Candidate Management UI.
6. Interview Scheduling UI.
7. Interview Room UI integration.
8. Report UI.
9. Candidate Portal.
10. Mock Interview Portal.
11. Question Bank UI.
12. Settings/Workspace UI cơ bản.
13. Empty/loading/error state toàn hệ thống.

---

## 2. Mục tiêu module của Lai

Mục tiêu chính:

- Người dùng thao tác được hệ thống dễ dàng.
- Recruiter tạo job, thêm ứng viên, tạo lịch phỏng vấn và xem report được.
- Candidate tham gia phỏng vấn và luyện mock interview được.
- UI kết nối được API của Khôi.
- UI kết nối được realtime event của Hùng.
- Tất cả màn hình có loading/error/empty state.
- Giao diện rõ ràng, dễ hiểu, không rối.

---

## 3. Phạm vi công việc

## 3.1. UI Foundation

### Chức năng

- Layout chính.
- Sidebar/navbar.
- Route structure.
- Protected route.
- Role-based navigation.
- Design system component cơ bản:
  - Button
  - Input
  - Select
  - Modal
  - Table
  - Badge
  - Card
  - Tabs
  - Toast
  - Empty state
  - Loading state

### DoD

- App có layout thống nhất.
- Navigation phân biệt Recruiter/Candidate.
- Component tái sử dụng được.
- Responsive cơ bản.

---

## 3.2. Auth UI

### Chức năng

- Login page.
- Register page.
- Forgot password page nếu có.
- Accept interview invite page.
- Auth guard.
- Redirect theo role.

### DoD

- Login thành công chuyển đúng dashboard.
- Sai tài khoản hiển thị lỗi rõ.
- Loading khi submit.
- Không login thì không vào trang nội bộ.

---

## 3.3. Recruiter Dashboard

### Chức năng

Hiển thị tổng quan:

- Số job đang mở.
- Số ứng viên mới.
- Lịch phỏng vấn hôm nay.
- Report chờ xem.
- Candidate cần quyết định.
- Phím tắt tạo job/tạo lịch/thêm candidate.

### DoD

- Dashboard load được dữ liệu từ API hoặc mock data tạm.
- Có empty state.
- Có loading state.
- Có error state.
- Click sang đúng màn hình chi tiết.

---

## 3.4. Job Management UI

### Màn hình

- Job List.
- Job Create/Edit.
- Job Detail.
- JD editor.
- AI Analyze JD panel.
- Rubric panel.
- Question suggestions.
- Candidate list theo job.

### Chức năng

- Tạo job.
- Sửa job.
- Đóng/mở job.
- Xem danh sách candidate trong job.
- Bấm AI analyze JD.
- Xem rubric AI đề xuất.
- Xem câu hỏi AI đề xuất.

### DoD

- CRUD job chạy với API Khôi.
- Form validate rõ.
- AI analyze có loading/error/result.
- Không mất dữ liệu form khi API lỗi.

---

## 3.5. Candidate Management UI

### Màn hình

- Candidate List.
- Candidate Create/Edit.
- Candidate Detail.
- CV Preview.
- Parsed CV panel.
- Interview history.
- Candidate status pipeline.

### Chức năng

- Thêm candidate.
- Sửa candidate.
- Upload CV.
- Xem AI parsed CV.
- Gán candidate vào job.
- Cập nhật trạng thái candidate.
- Xem lịch sử phỏng vấn.

### DoD

- Candidate CRUD chạy với API.
- Upload CV có loading/progress nếu có.
- Parsed CV hiển thị rõ.
- Không leak candidate giữa company.

---

## 3.6. Interview Scheduling UI

### Màn hình

- Interview List.
- Create Interview Modal/Page.
- Interview Detail.
- Invite Link panel.
- Calendar/List view cơ bản.

### Chức năng

- Tạo lịch phỏng vấn.
- Chọn job.
- Chọn candidate.
- Chọn recruiter.
- Chọn thời gian.
- Copy invite link.
- Cancel/reschedule nếu có.

### DoD

- Tạo interview thành công.
- Validation thời gian.
- Link phòng hiển thị/copy được.
- Status interview rõ ràng.

---

## 3.7. Interview Room UI Integration

Lai không làm lõi realtime, nhưng làm UI để tích hợp module của Hùng.

### Thành phần UI

- Header: job, candidate, timer, status, end button.
- Video area.
- Chat panel.
- Transcript panel.
- AI Suggestion panel.
- Question checklist.
- Recruiter note panel.
- Score panel.
- Candidate profile sidebar.
- Reconnect banner.
- Media permission warning.

### DoD

- UI subscribe event từ Hùng.
- Chat gửi/nhận được.
- Transcript hiển thị realtime.
- AI suggestion chỉ hiển thị cho Recruiter.
- Score panel chỉ hiển thị cho Recruiter.
- Candidate UI không thấy note/score nội bộ.
- Có loading/reconnect/error state.

---

## 3.8. Report UI

### Màn hình

- Interview Report Detail.
- Candidate Report Summary.
- Job Candidate Comparison cơ bản.

### Thành phần report

- Tổng điểm.
- Điểm theo tiêu chí.
- Điểm mạnh.
- Điểm yếu.
- Rủi ro.
- AI recommendation.
- Recruiter final decision.
- Transcript.
- Recruiter notes.

### DoD

- Report hiển thị đúng dữ liệu từ Khôi.
- Có trạng thái pending/generating/ready/failed.
- Có nút retry nếu được phép.
- Candidate không xem report nội bộ.
- Recruiter có thể nhập final decision.

---

## 3.9. Candidate Portal

### Màn hình

- Candidate Home.
- My Interviews.
- Interview Invite Confirmation.
- Profile/CV page.
- Mock Interview Home.
- Mock Interview Result.

### Chức năng

- Candidate xem lịch phỏng vấn của mình.
- Candidate xác nhận thông tin trước khi vào room.
- Candidate vào room bằng link.
- Candidate upload/cập nhật CV nếu được cho phép.
- Candidate xem feedback mock interview.

### DoD

- Candidate chỉ xem dữ liệu của mình.
- Luồng vào phòng đơn giản.
- Có thông báo consent nếu có ghi âm/AI.
- UI thân thiện, ít gây áp lực.

---

## 3.10. Mock Interview Portal

### Màn hình

- Chọn vị trí luyện tập.
- Chọn level.
- Upload CV hoặc nhập profile.
- Mock Interview Room.
- Answer input bằng text/audio nếu có.
- Feedback từng câu.
- Final feedback report.
- History.

### Chức năng

- Tạo mock session.
- AI hỏi câu hỏi.
- Candidate trả lời.
- Hiển thị feedback.
- Lưu lịch sử luyện tập.

### DoD

- Candidate tạo mock interview được.
- AI question hiển thị rõ.
- Candidate gửi answer được.
- Feedback hiển thị dễ hiểu.
- Lịch sử luyện tập xem được.

---

## 3.11. Question Bank UI

### Chức năng

- Xem danh sách câu hỏi.
- Tạo/sửa/xóa câu hỏi.
- Lọc theo job/skill/level/type.
- Gán câu hỏi vào interview template.
- Xem câu hỏi AI đề xuất.

### DoD

- Danh sách câu hỏi dùng được.
- Filter/search cơ bản.
- Gán câu hỏi vào job/template được.
- Không làm rối UI job detail.

---

## 4. Sprint chia việc cho Lai

## Sprint 0: UI Foundation

| Task | Mô tả | Độ khó | DoD |
|---|---|---:|---|
| L-S0-01 | Setup layout chính | Trung bình | Sidebar/navbar/content layout |
| L-S0-02 | Tạo component base | Trung bình | Button/Input/Table/Modal/Card |
| L-S0-03 | Protected route UI | Trung bình | Chưa login bị redirect |
| L-S0-04 | Empty/loading/error components | Dễ | Dùng lại toàn app |

## Sprint 1: Auth + Dashboard

| Task | Mô tả | Độ khó | DoD |
|---|---|---:|---|
| L-S1-01 | Login/Register UI | Trung bình | Submit/loading/error |
| L-S1-02 | Role-based redirect | Trung bình | Recruiter/Candidate đúng portal |
| L-S1-03 | Recruiter dashboard | Trung bình | Card tổng quan + quick actions |
| L-S1-04 | Candidate home cơ bản | Dễ | Xem lời chào/lịch gần nhất |

## Sprint 2: Job UI

| Task | Mô tả | Độ khó | DoD |
|---|---|---:|---|
| L-S2-01 | Job list | Trung bình | Search/filter/status |
| L-S2-02 | Job create/edit form | Trung bình | Validate + submit API |
| L-S2-03 | Job detail | Trung bình | JD + candidates + interviews |
| L-S2-04 | AI analyze JD panel | Khó | Loading/error/result rõ |

## Sprint 3: Candidate UI

| Task | Mô tả | Độ khó | DoD |
|---|---|---:|---|
| L-S3-01 | Candidate list | Trung bình | Search/filter/status |
| L-S3-02 | Candidate create/edit | Trung bình | Validate + submit API |
| L-S3-03 | Candidate detail | Trung bình | CV/profile/history |
| L-S3-04 | CV upload/preview | Khó | Upload + xem metadata/preview |

## Sprint 4: Interview Scheduling UI

| Task | Mô tả | Độ khó | DoD |
|---|---|---:|---|
| L-S4-01 | Interview list | Trung bình | Status + filter |
| L-S4-02 | Create interview flow | Khó | Chọn job/candidate/time |
| L-S4-03 | Interview detail | Trung bình | Info + invite link |
| L-S4-04 | Copy invite link | Dễ | Copy thành công + toast |

## Sprint 5: Interview Room UI

| Task | Mô tả | Độ khó | DoD |
|---|---|---:|---|
| L-S5-01 | Room layout | Khó | Header/video/sidebar/panels |
| L-S5-02 | Chat panel integration | Trung bình | Gửi/nhận với Hùng |
| L-S5-03 | Transcript panel | Khó | Hiển thị realtime |
| L-S5-04 | AI suggestion/score panel | Khó | Chỉ Recruiter thấy |

## Sprint 6: Report UI

| Task | Mô tả | Độ khó | DoD |
|---|---|---:|---|
| L-S6-01 | Report detail page | Trung bình | Hiển thị summary/score |
| L-S6-02 | Rubric score visualization | Trung bình | Điểm theo tiêu chí dễ đọc |
| L-S6-03 | Final decision UI | Trung bình | Recruiter chọn decision |
| L-S6-04 | Report status handling | Dễ | pending/generating/failed |

## Sprint 7: Mock Interview Portal

| Task | Mô tả | Độ khó | DoD |
|---|---|---:|---|
| L-S7-01 | Mock setup page | Trung bình | Chọn role/level/CV |
| L-S7-02 | Mock interview UI | Khó | AI hỏi, candidate trả lời |
| L-S7-03 | Feedback result UI | Trung bình | Feedback từng câu + tổng kết |
| L-S7-04 | Mock history | Dễ | Xem lịch sử luyện tập |

## Sprint 8: Polish + UX

| Task | Mô tả | Độ khó | DoD |
|---|---|---:|---|
| L-S8-01 | Responsive pass | Trung bình | Desktop/tablet ổn |
| L-S8-02 | Empty state toàn app | Dễ | Không có màn hình trống xấu |
| L-S8-03 | Error state toàn app | Dễ | API lỗi hiển thị rõ |
| L-S8-04 | UX review toàn luồng | Trung bình | Recruiter/Candidate flow mượt |

---

## 5. Dependency với Hùng

Lai cần Hùng cung cấp:

- Client event contract.
- WebSocket connection method.
- Room join/leave events.
- Chat events.
- Transcript events.
- AI suggestion events.
- AI score events.
- Reconnect state.

Lai cần tích hợp:

- Chat panel.
- Transcript panel.
- AI suggestion panel.
- Score panel.
- Room presence UI.
- Media status UI.

---

## 6. Dependency với Khôi

Lai cần Khôi cung cấp:

- Auth API.
- Job API.
- Candidate API.
- Interview API.
- AI analyze JD/CV API.
- Report API.
- Mock Interview API.
- Permission response rõ ràng.

Lai cần thống nhất với Khôi:

- Field form.
- Validation rule.
- Report data shape.
- Candidate status.
- Job status.
- Interview status.

---

## 7. Checklist test riêng của Lai

- Login thành công/thất bại.
- Protected route hoạt động.
- Recruiter thấy dashboard.
- Candidate không vào được recruiter dashboard.
- Job list/create/edit/detail hoạt động.
- Candidate list/create/edit/detail hoạt động.
- Upload CV có loading/error.
- Tạo interview được.
- Copy invite link được.
- Vào room UI không vỡ layout.
- Chat hiển thị realtime.
- Transcript hiển thị realtime.
- Candidate không thấy score/note nội bộ.
- Report pending/ready/failed hiển thị đúng.
- Mock interview flow chạy được.
- Empty/loading/error state đầy đủ.

---

## 8. Prompt mẫu cho Lai dùng với AI coding assistant

```text
Bạn đang làm module của Lai trong dự án phỏng vấn cùng AI real-time.
Phạm vi của Lai là UI/UX, Recruiter Portal, Candidate Portal, Job UI, Candidate UI, Interview Scheduling UI, Interview Room UI integration, Report UI và Mock Interview Portal.

Yêu cầu task hiện tại:
[Điền task]

Trước khi code hãy:
1. Đọc cấu trúc frontend hiện tại.
2. Tận dụng component có sẵn, không tạo trùng lặp vô tội vạ.
3. Không tự ý sửa realtime core của Hùng.
4. Không tự ý sửa AI/backend schema của Khôi.
5. Mọi màn hình phải có loading/error/empty state.
6. Candidate không được thấy dữ liệu nội bộ của Recruiter.
7. Form phải có validation.
8. Sau khi code, liệt kê file đã sửa và checklist test.
```

---

## 9. Kết quả cuối cùng Lai cần bàn giao

Lai hoàn thành khi hệ thống có:

- UI foundation ổn định.
- Auth UI dùng được.
- Recruiter Dashboard dùng được.
- Job Management UI đầy đủ.
- Candidate Management UI đầy đủ.
- Interview Scheduling UI đầy đủ.
- Interview Room UI tích hợp realtime.
- Report UI hiển thị rõ ràng.
- Candidate Portal dùng được.
- Mock Interview Portal dùng được.
- Loading/error/empty state đầy đủ.
- UX dễ hiểu cho cả Recruiter và Candidate.
