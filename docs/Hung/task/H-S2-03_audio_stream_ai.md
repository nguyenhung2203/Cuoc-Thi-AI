# H-S2-03 — Audio Stream Hook cho AI

| Field | Value |
|---|---|
| **Task ID** | H-S2-03 |
| **Sprint** | 2 — Chat + LiveKit WebRTC |
| **Độ khó** | Rất khó |
| **Trạng thái** | ✅ Hoàn thành |
| **Owner** | Hùng |

---

## Mục tiêu

Hook audio stream từ LiveKit SFU để gửi cho AI Orchestrator (Whisper STT):
- LiveKit server-side SDK subscribe audio tracks
- Forward audio chunks → AI Orchestrator (Python)
- AI Orchestrator chạy STT → trả transcript về WebSocket gateway

---

## Yêu cầu chi tiết

### 1. Kiến trúc

```
LiveKit SFU → Audio stream hook (Go service) → AI Orchestrator (Python)
                                                      ↓
                                              Whisper STT → transcript
                                                      ↓
                                              WebSocket Gateway → broadcast
```

### 2. Audio Hook Options

- **Option A**: LiveKit Egress API — record audio và xử lý
- **Option B**: LiveKit server-side SDK — subscribe room audio tracks
- **Option C**: LiveKit webhook — nhận audio events

### 3. Forward to AI

- Gửi audio chunks qua internal REST hoặc gRPC đến Python AI Orchestrator
- AI Orchestrator chạy Whisper → trả transcript text
- Transcript gửi lại WebSocket gateway → broadcast `transcript:update`

---

## Definition of Done

- [x] Audio từ phòng được capture qua LiveKit
- [x] Audio forward đến AI Orchestrator thành công
- [x] AI trả transcript → broadcast qua WebSocket
- [x] Nếu AI/STT lỗi → audio vẫn hoạt động bình thường

---

## Dependency

- **H-S2-02** — LiveKit integration
- **Khôi** — AI Orchestrator STT endpoint

---

## Checklist test

- [x] Nói trong phòng → audio được capture
- [x] Audio gửi đến AI → nhận transcript text
- [x] Transcript broadcast đúng room
- [x] STT lỗi → room vẫn hoạt động
