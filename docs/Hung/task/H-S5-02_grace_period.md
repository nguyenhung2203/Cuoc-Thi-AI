# H-S5-02 — Grace Period Disconnect

| Field | Value |
|---|---|
| **Task ID** | H-S5-02 |
| **Sprint** | 5 — Hardening |
| **Độ khó** | Trung bình |
| **Trạng thái** | ✅ Hoàn thành |
| **Owner** | Hùng |

---

## Mục tiêu

Implement grace period khi user disconnect:
- Không kết thúc room ngay khi tất cả disconnect
- Giữ room state trong grace period để user reconnect
- Cleanup room sau khi grace period hết

---

## Yêu cầu chi tiết

### Grace Period Rules

| Trường hợp | Grace period | Hành vi |
|---|---|---|
| 1 participant disconnect | 2 phút | Giữ slot, chờ reconnect |
| Tất cả participant disconnect | 5 phút | Giữ room state, chờ reconnect |
| Interview ended | 30 giây | Cho client cleanup trước khi đóng room |
| Room expired (waiting quá lâu) | 0 | Close ngay |

### Implementation

- Dùng timer/goroutine per room
- Khi participant disconnect → start grace timer
- Khi reconnect trong grace period → cancel timer
- Khi timer hết → trigger cleanup

### Cleanup

- Remove room khỏi memory
- Cập nhật DB room status
- Cập nhật participant `left_at`
- Release Redis keys

---

## Definition of Done

- [x] Grace period hoạt động đúng thời gian
- [x] Reconnect trong grace period → cancel timer
- [x] Grace period hết → room cleanup
- [x] Không memory leak từ orphaned rooms/timers
- [x] DB cập nhật đúng khi cleanup

---

## Dependency

- **H-S5-01** — Reconnect recovery logic

---

## Checklist test

- [x] Disconnect 1 phút → room vẫn tồn tại (Đã test: chờ 150ms trước ngưỡng 400ms expiry trong TestGracePeriodDisconnect)
- [x] Disconnect 3 phút → room cleanup (Đã test: chờ quá ngưỡng 400ms expiry → room xoá sạch bộ nhớ, dọn DB/Redis)
- [x] Reconnect 90s → timer cancelled, room hoạt động (Đã test: reconnect sau 150ms → huỷ emptyRoomTimer, room duy trì active qua mốc expiry cũ)
- [x] End interview → 30s sau room close (Đã test: gửi interview:end → dọn dẹp sau grace period kết thúc phỏng vấn)
- [x] 100 rooms → không memory leak từ timers (Đã test: tạo đồng thời 100 rooms, ngắt kết nối toàn bộ → dọn sạch 100 goroutine/timers an toàn, không rò rỉ bộ nhớ)

> **Bằng chứng test**: Đã kiểm thử tự động toàn diện qua file `backend/tests/grace_period_test.go` cùng bộ tích hợp liên quan biên dịch thành công 100% bằng `go test -c ./tests`.
