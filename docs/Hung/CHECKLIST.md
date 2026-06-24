# CHECKLIST — Module Hùng: Realtime & Interview Room

> **Owner**: Hùng
> **Phạm vi**: Interview Room, Realtime Gateway, WebSocket, LiveKit WebRTC, Presence, Chat, Transcript Realtime, AI Bridge
> **Tài liệu gốc**: [MODULE_HUNG_REALTIME_AI.md](./MODULE_HUNG_REALTIME_AI.md)

---

## Sprint 0: Setup Realtime Foundation

- [ ] **H-S0-01** — Thiết kế event contract realtime → [task/H-S0-01_event_contract.md](./task/H-S0-01_event_contract.md)
- [ ] **H-S0-02** — Tạo WebSocket gateway skeleton → [task/H-S0-02_websocket_gateway.md](./task/H-S0-02_websocket_gateway.md)
- [ ] **H-S0-03** — Xác thực WebSocket bằng token → [task/H-S0-03_websocket_auth.md](./task/H-S0-03_websocket_auth.md)
- [ ] **H-S0-04** — Room channel architecture → [task/H-S0-04_room_channel.md](./task/H-S0-04_room_channel.md)

---

## Sprint 1: Interview Room MVP

- [ ] **H-S1-01** — Join/Leave room → [task/H-S1-01_join_leave_room.md](./task/H-S1-01_join_leave_room.md)
- [ ] **H-S1-02** — Presence realtime → [task/H-S1-02_presence_realtime.md](./task/H-S1-02_presence_realtime.md)
- [ ] **H-S1-03** — Start/End interview event → [task/H-S1-03_start_end_interview.md](./task/H-S1-03_start_end_interview.md)
- [ ] **H-S1-04** — Room status sync → [task/H-S1-04_room_status_sync.md](./task/H-S1-04_room_status_sync.md)

---

## Sprint 2: Chat + LiveKit WebRTC

- [ ] **H-S2-01** — Chat realtime trong room → [task/H-S2-01_chat_realtime.md](./task/H-S2-01_chat_realtime.md)
- [ ] **H-S2-02** — LiveKit WebRTC integration → [task/H-S2-02_livekit_webrtc.md](./task/H-S2-02_livekit_webrtc.md)
- [ ] **H-S2-03** — Audio stream hook cho AI → [task/H-S2-03_audio_stream_ai.md](./task/H-S2-03_audio_stream_ai.md)
- [ ] **H-S2-04** — Media status event → [task/H-S2-04_media_status.md](./task/H-S2-04_media_status.md)

---

## Sprint 3: Transcript Realtime

- [ ] **H-S3-01** — Transcript event pipeline → [task/H-S3-01_transcript_pipeline.md](./task/H-S3-01_transcript_pipeline.md)
- [ ] **H-S3-02** — Partial/Final transcript handling → [task/H-S3-02_partial_final_transcript.md](./task/H-S3-02_partial_final_transcript.md)
- [ ] **H-S3-03** — Speaker mapping → [task/H-S3-03_speaker_mapping.md](./task/H-S3-03_speaker_mapping.md)
- [ ] **H-S3-04** — Save transcript final → [task/H-S3-04_save_transcript.md](./task/H-S3-04_save_transcript.md)

---

## Sprint 4: AI Realtime Bridge

- [ ] **H-S4-01** — AI suggestion event → [task/H-S4-01_ai_suggestion.md](./task/H-S4-01_ai_suggestion.md)
- [ ] **H-S4-02** — AI score update event → [task/H-S4-02_ai_score_update.md](./task/H-S4-02_ai_score_update.md)
- [ ] **H-S4-03** — Role-based event visibility → [task/H-S4-03_role_visibility.md](./task/H-S4-03_role_visibility.md)
- [ ] **H-S4-04** — AI error event handling → [task/H-S4-04_ai_error_handling.md](./task/H-S4-04_ai_error_handling.md)

---

## Sprint 5: Hardening

- [ ] **H-S5-01** — Reconnect recovery → [task/H-S5-01_reconnect_recovery.md](./task/H-S5-01_reconnect_recovery.md)
- [ ] **H-S5-02** — Grace period disconnect → [task/H-S5-02_grace_period.md](./task/H-S5-02_grace_period.md)
- [ ] **H-S5-03** — Room event audit log → [task/H-S5-03_room_audit.md](./task/H-S5-03_room_audit.md)
- [ ] **H-S5-04** — Load test nhiều room → [task/H-S5-04_load_test.md](./task/H-S5-04_load_test.md)

---

## Tổng kết tiến độ

| Sprint | Tổng task | Hoàn thành | Trạng thái |
|---|---:|---:|---|
| Sprint 0: Foundation | 4 | 0 | ⬜ Chưa bắt đầu |
| Sprint 1: Room MVP | 4 | 0 | ⬜ Chưa bắt đầu |
| Sprint 2: Chat + WebRTC | 4 | 0 | ⬜ Chưa bắt đầu |
| Sprint 3: Transcript | 4 | 0 | ⬜ Chưa bắt đầu |
| Sprint 4: AI Bridge | 4 | 0 | ⬜ Chưa bắt đầu |
| Sprint 5: Hardening | 4 | 0 | ⬜ Chưa bắt đầu |
| **Tổng** | **24** | **0** | |

---

## Dependency chính

### Cần Khôi cung cấp trước
- [ ] Auth middleware + JWT validation
- [ ] Database schema: `interview_rooms`, `interview_participants`, `interview_transcripts`, `ai_suggestions`
- [ ] Interview API: tạo/lấy interview
- [ ] AI suggestion service endpoint
- [ ] AI scoring service endpoint
- [ ] Report generation trigger

### Cần Lai tích hợp sau
- [ ] UI Interview Room layout
- [ ] Chat panel component
- [ ] Transcript panel component
- [ ] AI suggestion panel component
- [ ] Score panel component
- [ ] Reconnect banner UI

---

## Quy tắc khi làm task

1. Đọc [REALTIME_EVENTS.md](../REALTIME_EVENTS.md) trước khi code event
2. Không sửa AI scoring/report schema của Khôi
3. Không redesign UI của Lai
4. Recruiter thấy AI suggestion/score — Candidate **KHÔNG** thấy
5. AI lỗi → room vẫn hoạt động bình thường
6. Mỗi task xong → ghi file đã sửa + checklist test
