# H-S4-03 — Role-based Event Visibility

| Field | Value |
|---|---|
| **Task ID** | H-S4-03 |
| **Sprint** | 4 — AI Realtime Bridge |
| **Độ khó** | Khó |
| **Trạng thái** | ⬜ Chưa bắt đầu |
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

- [ ] Visibility filter áp dụng cho TẤT CẢ event types
- [ ] Unit test cho mỗi event × role combination
- [ ] Candidate KHÔNG BAO GIỜ nhận event recruiter-only
- [ ] Mock interview exception: candidate thấy `ai:thinking`
- [ ] Security audit pass — không có bypass path

---

## Dependency

- **H-S4-01**, **H-S4-02** — AI events cần visibility filter

---

## Checklist test

- [ ] ai:suggestion → candidate KHÔNG nhận
- [ ] ai:score_update → candidate KHÔNG nhận
- [ ] ai:warning → candidate KHÔNG nhận
- [ ] note:create → candidate KHÔNG nhận
- [ ] report:ready → candidate KHÔNG nhận
- [ ] chat:message room → CẢ HAI nhận
- [ ] room:user_joined → CẢ HAI nhận
