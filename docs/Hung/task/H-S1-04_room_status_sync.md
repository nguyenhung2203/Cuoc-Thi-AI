# H-S1-04 — Room Status Sync

| Field | Value |
|---|---|
| **Task ID** | H-S1-04 |
| **Sprint** | 1 — Interview Room MVP |
| **Độ khó** | Khó |
| **Trạng thái** | ✅ Hoàn thành |
| **Owner** | Hùng |

---

## Mục tiêu

Đồng bộ trạng thái room giữa server, DB và tất cả client:
- Room status được lưu trong DB và Redis
- Mọi thay đổi status đều broadcast cho client
- Client reload trang vẫn nhận đúng trạng thái hiện tại
- Status machine rõ ràng, không chuyển trạng thái sai

---

## Yêu cầu chi tiết

### 1. Room Status State Machine

```
[*] → scheduled
scheduled → waiting       (room opened / first participant join)
waiting → active          (recruiter starts interview)
active → paused           (recruiter pauses)
paused → active           (recruiter resumes)
active → completed        (recruiter ends interview)
scheduled → cancelled     (cancelled trước khi bắt đầu)
waiting → expired         (invite hết hạn, quá thời gian chờ)
```

### 2. Status Update Flow

1. Event trigger → validate transition hợp lệ
2. Cập nhật room status trong DB (`interview_rooms`)
3. Cập nhật interview status trong DB (`interviews`)
4. Cache status mới vào Redis
5. Broadcast status change cho tất cả participant

### 3. Client Reconnect Sync

Khi client reconnect (join lại room):
- Server gửi `room:joined` kèm `room_status` hiện tại
- Client tự điều chỉnh UI theo status
- Nếu `completed` → client hiển thị màn hình kết thúc

### 4. Auto-expire

- Nếu room ở `waiting` quá lâu (configurable, ví dụ 30 phút) → chuyển `expired`
- Dùng Redis TTL hoặc cron job

---

## Tài liệu tham chiếu

- [REALTIME_EVENTS.md](../../REALTIME_EVENTS.md) — Mục 3 (Room lifecycle state diagram)
- [DATABASE_DESIGN.md](../../DATABASE_DESIGN.md) — Bảng `interview_rooms` (status), `interviews` (status)

---

## Definition of Done

- [x] Room status machine đúng theo spec
- [x] Không có transition bất hợp lệ (ví dụ: completed → active)
- [x] Mọi status change broadcast cho client
- [x] Client reload → nhận đúng room status
- [x] Redis cache status đúng
- [x] Auto-expire room quá hạn
- [x] DB `interview_rooms.status` và `interviews.status` sync

---

## Dependency

- **H-S1-01** — Join/leave
- **H-S1-03** — Start/end interview

---

## Checklist test

- [x] Room mới → status = `waiting`
- [x] Start → status = `active`
- [x] End → status = `completed`
- [x] Client reload khi active → nhận `room_status: "active"`
- [x] Chuyển completed → active → bị reject
- [x] Room waiting quá lâu → auto expire
