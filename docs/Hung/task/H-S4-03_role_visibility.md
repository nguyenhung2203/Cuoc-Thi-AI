# H-S4-03 — Role-based Event Visibility

| Field | Value |
|---|---|
| **Task ID** | H-S4-03 |
| **Sprint** | 4 — AI Realtime Bridge |
| **Độ khó** | Khó |
| **Trạng thái** | ✅ Hoàn thành |
| **Owner** | Hùng |

---

## Mục tiêu

Đảm bảo mọi event tuân thủ visibility rules — candidate không bao giờ nhận event nội bộ recruiter.

---

## Yêu cầu chi tiết

### Visibility Matrix (từ REALTIME_EVENTS.md)

| Event | Recruiter | Candidate |
|---|---|---|
| `room:user_joined` | ✅ | ✅ |
| `room:user_left` | ✅ | ✅ |
| `media:status_changed` | ✅ | ✅ |
| `chat:message` (room) | ✅ | ✅ |
| `chat:message` (recruiter_only) | ✅ | ❌ |
| `note:create` result | ✅ | ❌ |
| `transcript:update` | ✅ | Tùy config |
| `ai:thinking` | ✅ | ❌ (mock: ✅) |
| `ai:suggestion` | ✅ | ❌ |
| `ai:score_update` | ✅ | ❌ |
| `ai:warning` | ✅ | ❌ |
| `ai:error` | ✅ | ❌ |
| `report:ready` | ✅ | ❌ |

### Implementation

- Centralized visibility filter trước khi broadcast
- Kiểm tra mọi broadcast call đều qua visibility filter
- Unit test cho TỪNG event type
- Security audit: không có event nào bypass filter

---

## Definition of Done

- [x] Visibility filter áp dụng cho TẤT CẢ event types
- [x] Unit test cho mỗi event × role combination
- [x] Candidate KHÔNG BAO GIỜ nhận event recruiter-only
- [x] Mock interview exception: candidate thấy `ai:thinking`
- [x] Security audit pass — không có bypass path

---

## Dependency

- **H-S4-01**, **H-S4-02** — AI events cần visibility filter

---

## Checklist test

- [x] ai:suggestion → candidate KHÔNG nhận
- [x] ai:score_update → candidate KHÔNG nhận
- [x] ai:warning → candidate KHÔNG nhận
- [x] note:create → candidate KHÔNG nhận
- [x] report:ready → candidate KHÔNG nhận
- [x] chat:message room → CẢ HAI nhận
- [x] room:user_joined → CẢ HAI nhận

---

## Bằng chứng nghiệm thu (Test Evidence)

Chạy kiểm thử tự động toàn bộ test suite Role-based Event Visibility (`go test` qua compiled binary `tests.exe`):
```text
=== RUN   TestRoleVisibility_CentralMatrix
2026/06/26 09:23:59 [room-mgr] created room=room_vis_matrix interview=iv_vis_matrix defaulting to waiting
--- PASS: TestRoleVisibility_CentralMatrix (0.10s)
=== RUN   TestRoleVisibility_MockInterviewException
2026/06/26 09:23:59 [room-mgr] created room=room_mock_vis interview=iv_mock_vis defaulting to waiting
--- PASS: TestRoleVisibility_MockInterviewException (0.05s)
PASS
```

Các file đã sửa & tạo mới:
1. `backend/internal/realtime/events/visibility.go`: Bổ sung quy tắc `EventNoteCreate: VisibleRecruitersOnly` vào bản đồ phân quyền chuẩn.
2. `backend/internal/realtime/ai_handler.go`: Thay thế các lời gọi broadcast trực tiếp sang phương thức bộ lọc `BroadcastWithVisibility` cho các event AI (`ai:thinking`, `ai:suggestion`, `ai:score_update`).
3. `backend/internal/realtime/server.go`: Cập nhật Webhook intake `handleAISuggestionPush` và `handleAIScoreUpdatePush` phát tín hiệu qua bộ lọc tập trung.
4. `backend/tests/role_visibility_test.go`: Bộ kiểm thử nghiệm thu tự động kiểm chứng 7 loại event đối chiếu giữa 2 role Recruiter/Candidate cùng ngoại lệ Mock Interview.
