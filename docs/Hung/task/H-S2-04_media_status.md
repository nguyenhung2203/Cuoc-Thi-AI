# H-S2-04 — Media Status Event

| Field | Value |
|---|---|
| **Task ID** | H-S2-04 |
| **Sprint** | 2 — Chat + LiveKit WebRTC |
| **Độ khó** | Trung bình |
| **Trạng thái** | ✅ Đã hoàn thành |
| **Owner** | Hùng |

---

## Mục tiêu

Đồng bộ trạng thái mic/camera/screen share giữa các participant:
- Client gửi `media:status` khi bật/tắt mic/camera
- Server broadcast `media:status_changed` cho tất cả participant
- Lưu `media_status` trong participant data

---

## Yêu cầu chi tiết

### 1. `media:status` (Client → Server)

```json
{
  "event": "media:status",
  "payload": {
    "mic_enabled": true,
    "camera_enabled": false,
    "screen_sharing": false
  }
}
```

### 2. `media:status_changed` (Server → Client)

```json
{
  "event": "media:status_changed",
  "payload": {
    "participant_id": "uuid",
    "mic_enabled": true,
    "camera_enabled": false,
    "screen_sharing": false
  }
}
```

### 3. Lưu vào Participant

- Cập nhật `media_status` JSONB trong `interview_participants`
- Client mới join nhận media status của participant hiện tại

---

## Definition of Done

- [x] Bật/tắt mic → broadcast đúng
- [x] Bật/tắt camera → broadcast đúng
- [x] Client mới join → nhận media status hiện tại
- [x] DB `media_status` cập nhật đúng
- [x] Rate limit: max 30 lần/phút/participant

---

## Dependency

- **H-S1-01** — Join room (cần participant data)

---

## Checklist test

- [x] Tắt mic → peer nhận `media:status_changed` mic_enabled=false
- [x] Join room → nhận media status của người đã có
- [x] Spam toggle → bị rate limit
