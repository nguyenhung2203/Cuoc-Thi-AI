# H-S5-01 — Reconnect Recovery

| Field | Value |
|---|---|
| **Task ID** | H-S5-01 |
| **Sprint** | 5 — Hardening |
| **Độ khó** | Rất khó |
| **Trạng thái** | ✅ Hoàn thành |
| **Owner** | Hùng |

---

## Mục tiêu

Khi user mất kết nối ngắn rồi reconnect, khôi phục state đầy đủ:
- Giữ participant state trong room (không bị xóa ngay)
- Khi reconnect thành công → gửi current state (room status, participants, recent events)
- Client sync lại transcript/chat bằng REST hoặc event backlog

---

## Definition of Done

- [x] Mất mạng < 2 phút → reconnect không mất state (Đã test integration: `TestReconnectRecovery/Disconnect_30s...`)
- [x] Mất mạng > 2 phút → user_left broadcast, vẫn join lại được (Đã test integration: `TestReconnectRecovery/Disconnect_3_phút...`)
- [x] Reconnect → nhận đúng room state hiện tại (Gửi ACK `room:joined` với `MediaStatus`, `SyncFromTimestamp`)
- [x] Không tạo duplicate participant khi reconnect (Tái sử dụng `Participant` struct và cập nhật `ConnectionID` trong track map)
- [x] Transcript không bị mất khi reconnect (Hệ thống tự động phát lại buffered events từ `GetMissedEvents`)

---

## Dependency

- **H-S1-02** — Presence realtime
- **H-S1-04** — Room status sync

---

## Checklist test

- [x] Disconnect 30s → reconnect → vẫn trong room (Bằng chứng: `reconnected to room=room1 (recovered 1 missed events)`)
- [x] Disconnect 30s → người kia KHÔNG nhận user_left (Bằng chứng: Recruiter nhận `room:presence_update` trạng thái reconnecting thay vì `room:user_left`)
- [x] Disconnect 3 phút → user_left broadcast → join lại thành công (Bằng chứng: Broadcast `user_left` reason `reconnect_timeout`, ứng viên dial & join lại thành công)
- [x] Reconnect → nhận room_status, participants đúng (Bằng chứng: Kiểm tra payload `RoomJoinedPayload.MissedEventsCount >= 1`)
- [x] Reconnect → transcript panel sync lại (Bằng chứng: Nhận lại đủ missed chat/transcript events ngay sau ACK join)
