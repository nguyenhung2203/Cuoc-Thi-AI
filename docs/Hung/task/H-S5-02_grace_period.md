# H-S5-02 — Grace Period Disconnect

| Field | Value |
|---|---|
| **Task ID** | H-S5-02 |
| **Sprint** | 5 — Hardening |
| **Độ khó** | Trung bình |
| **Trạng thái** | ⬜ Chưa bắt đầu |
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

- [ ] Grace period hoạt động đúng thời gian
- [ ] Reconnect trong grace period → cancel timer
- [ ] Grace period hết → room cleanup
- [ ] Không memory leak từ orphaned rooms/timers
- [ ] DB cập nhật đúng khi cleanup

---

## Dependency

- **H-S5-01** — Reconnect recovery logic

---

## Checklist test

- [ ] Disconnect 1 phút → room vẫn tồn tại
- [ ] Disconnect 3 phút → room cleanup
- [ ] Reconnect 90s → timer cancelled, room hoạt động
- [ ] End interview → 30s sau room close
- [ ] 100 rooms → không memory leak từ timers
