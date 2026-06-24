# H-S2-02 — LiveKit WebRTC Integration

| Field | Value |
|---|---|
| **Task ID** | H-S2-02 |
| **Sprint** | 2 — Chat + LiveKit WebRTC |
| **Độ khó** | Rất khó |
| **Trạng thái** | ⬜ Chưa bắt đầu |
| **Owner** | Hùng |

---

## Mục tiêu

Tích hợp LiveKit SFU cho audio/video WebRTC:
- Frontend dùng LiveKit SDK thay raw WebRTC
- Backend tạo LiveKit access token cho participant
- Recruiter và Candidate video/audio call qua LiveKit SFU
- Signaling qua LiveKit (không cần WebSocket signaling riêng)

---

## Yêu cầu chi tiết

### 1. LiveKit Server Setup

- Deploy LiveKit server (Docker hoặc cloud)
- Config LiveKit API key + secret
- Config TURN/STUN servers

### 2. Token Generation

- Backend tạo LiveKit JWT cho participant khi join room
- Token chứa: room_name, participant_identity, permissions
- Response trả `livekit_url` + `room_access_token`

### 3. Frontend Integration

- Dùng `livekit-client` SDK
- Connect LiveKit room với token
- Publish/subscribe audio/video tracks
- Handle track events (mute/unmute)

### 4. Kiến trúc

```
Frontend → LiveKit SDK → LiveKit SFU Server
                              ↓
                     Audio stream hook → AI Orchestrator
```

---

## Definition of Done

- [ ] Recruiter và Candidate video/audio call qua LiveKit
- [ ] Token generation từ backend hoạt động
- [ ] Mic/camera on/off đúng
- [ ] Reconnect media khi mạng chập chờn
- [ ] Nếu LiveKit lỗi → room WebSocket vẫn hoạt động

---

## Dependency

- **Sprint 0 + Sprint 1** hoàn thành
- **Khôi** — LiveKit server deployment

---

## Checklist test

- [ ] 2 peer connect qua LiveKit → thấy/nghe nhau
- [ ] Tắt mic → peer kia không nghe
- [ ] Tắt camera → peer kia không thấy
- [ ] Reload → reconnect media thành công
- [ ] LiveKit down → chat/room state vẫn hoạt động qua WebSocket
