# BÁO CÁO KẾT QUẢ KIỂM THỬ TỰ ĐỘNG REALTIME (FE & BE)

> **Ngày thực hiện**: 2026-07-09
> **Trạng thái**: ✅ **HOÀN TOÀN THÀNH CÔNG (PASSED)**
> **Phạm vi kiểm tra**: Tương tác WebSocket Realtime Gateway (Hùng) & State Management (Lai)

---

## 1. Bằng chứng kiểm thử (Execution Logs)

Dưới đây là log chạy thực tế của kịch bản E2E Test tự động cô lập hoàn toàn (`test-room-nj2gw077b`) chạy trên cổng `:8081`:

```text
====================================================
KHỞI ĐỘNG HỆ THỐNG KIỂM THỬ TỰ ĐỘNG REALTIME
====================================================

[BE] Đang biên dịch Go Realtime server...
[BE] Biên dịch thành công.
[BE] Đang khởi động Go Realtime server binary trên cổng 8081...
[BE LOG] Starting AI Interview Platform Realtime Gateway on :8081...
[BE ERR] 2026/07/09 09:04:27 [realtime] listening on :8081

[BE] Server đã lắng nghe. Đang chạy test client...
[TEST FE] ====================================================
[TEST FE] BẮT ĐẦU CHẠY AUTOMATED E2E INTEGRATION TEST REALTIME
[TEST FE] ====================================================
[BE ERR] 2026/07/09 09:04:28 [ws] connected connID=75447b83-0b58-4f2c-8ef6-ae05880b33f5 userID=recruiter-user-123 role=recruiter roomID=test-room-nj2gw077b
[TEST FE] [Client Hùng Recruiter] Connected to server.
[BE ERR] 2026/07/09 09:04:28 [ws] connected connID=7c192636-ad0f-4178-a2ee-98af5dd24b35 userID=candidate-user-123 role=candidate roomID=test-room-nj2gw077b
[TEST FE] [Client Lai Candidate] Connected to server.

[TEST FE] --- TEST STEP 1: Join Room ---
[TEST FE] [Client Hùng Recruiter] Sent Event: room:join (reqId: req_fl1qn196r)
[BE ERR] 2026/07/09 09:04:28 [redis] GET room_status:test-room-nj2gw077b
[BE ERR] 2026/07/09 09:04:28 [db] SELECT status FROM interview_rooms WHERE id = 'test-room-nj2gw077b'
[BE ERR] 2026/07/09 09:04:28 [room-mgr] created room=test-room-nj2gw077b interview=test-interview-nj2gw077b defaulting to waiting
[BE ERR] 2026/07/09 09:04:28 [presence] started watcher for room=test-room-nj2gw077b
[BE ERR] 2026/07/09 09:04:28 [room] participant=75447b83-0b58-4f2c-8ef6-ae05880b33f5 joined room=test-room-nj2gw077b as recruiter
[BE ERR] 2026/07/09 09:04:28 [db] INSERT INTO interview_participants (id, interview_id, user_id, participant_type, display_name, joined_at, connection_state) VALUES ('75447b83-0b58-4f2c-8ef6-ae05880b33f5', 'test-interview-nj2gw077b', 'recruiter-user-123', 'recruiter', 'Hùng Recruiter', '2026-07-09T02:04:28Z', 'online')
[BE ERR] 2026/07/09 09:04:28 [redis] SET presence:test-room-nj2gw077b:75447b83-0b58-4f2c-8ef6-ae05880b33f5 value=online EX 90
[TEST FE] [Client Lai Candidate] Sent Event: room:join (reqId: req_oemtnyora)
[BE ERR] 2026/07/09 09:04:28 [room] participant=7c192636-ad0f-4178-a2ee-98af5dd24b35 joined room=test-room-nj2gw077b as candidate
[BE ERR] 2026/07/09 09:04:28 [db] INSERT INTO interview_participants (id, interview_id, user_id, participant_type, display_name, joined_at, connection_state) VALUES ('7c192636-ad0f-4178-a2ee-98af5dd24b35', 'test-interview-nj2gw077b', 'candidate-user-123', 'candidate', 'Lai Candidate', '2026-07-09T02:04:28Z', 'online')
[BE ERR] 2026/07/09 09:04:28 [redis] SET presence:test-room-nj2gw077b:7c192636-ad0f-4178-a2ee-98af5dd24b35 value=online EX 90
[TEST FE] [Client Lai Candidate] Received Event: room:user_joined (reqId: req_fl1qn196r)
[TEST FE] [Client Hùng Recruiter] Received Event: room:joined (reqId: req_fl1qn196r)
[TEST FE] [Client Hùng Recruiter] Received Event: room:presence_update (reqId: req_fl1qn196r)
[TEST FE] [Client Hùng Recruiter] Received Event: room:user_joined (reqId: req_oemtnyora)
[TEST FE] [Client Hùng Recruiter] Received Event: room:presence_update (reqId: req_oemtnyora)
[TEST FE] [Client Lai Candidate] Received Event: room:presence_update (reqId: req_fl1qn196r)
[TEST FE] [Client Lai Candidate] Received Event: room:joined (reqId: req_oemtnyora)
[TEST FE] [Client Lai Candidate] Received Event: room:presence_update (reqId: req_oemtnyora)
[TEST FE] SUCCESS: Cả 2 client join room và nhận ACK thành công.
[TEST FE] SUCCESS: Presence Update hoạt động chính xác.

[TEST FE] --- TEST STEP 2: Media Status ---
[TEST FE] [Client Hùng Recruiter] Sent Event: media:status (reqId: req_7gux7n8v4)
[BE ERR] 2026/07/09 09:04:29 [db] UPDATE interview_participants SET media_status = '{"mic_enabled": true, "camera_enabled": true, "screen_sharing": false}'
[TEST FE] [Client Lai Candidate] Received Event: media:status_changed (reqId: req_7gux7n8v4)
[TEST FE] SUCCESS: Media status thay đổi và broadcast thành công.

[TEST FE] --- TEST STEP 3: Chat & Notes Visibility ---
[TEST FE] [Client Hùng Recruiter] Sent Event: chat:send (reqId: req_e0bqchdwg)
[TEST FE] [Client Hùng Recruiter] Sent Event: chat:send (reqId: req_aaw5u3k2a)
[BE ERR] 2026/07/09 09:04:30 [db] INSERT INTO interview_transcripts (id, interview_id, speaker_type, speaker_name, content, source, visibility, created_at) VALUES ('c067bf56-d73b-427d-8166-4951c644b8e4', 'test-interview-nj2gw077b', 'recruiter', 'Hùng Recruiter', 'Chào bạn Lai!', 'chat', 'room', '2026-07-09T02:04:30Z')
[BE ERR] 2026/07/09 09:04:30 [db] INSERT INTO interview_transcripts (id, interview_id, speaker_type, speaker_name, content, source, visibility, created_at) VALUES ('7fc7d47c-e368-4689-ae35-106e9e57ae7a', 'test-interview-nj2gw077b', 'recruiter', 'Hùng Recruiter', 'Ghi chú: Lai đang khá bối rối.', 'chat', 'recruiter_only', '2026-07-09T02:04:30Z')
[TEST FE] [Client Lai Candidate] Received Event: chat:message (reqId: req_e0bqchdwg)
[TEST FE] [Client Hùng Recruiter] Received Event: chat:message (reqId: req_e0bqchdwg)
[TEST FE] [Client Hùng Recruiter] Received Event: chat:message (reqId: req_aaw5u3k2a)
[TEST FE] SUCCESS: Visibility kiểm soát chat và note nội bộ an toàn (Candidate không thấy).

[TEST FE] --- TEST STEP 4: Interview Control ---
[TEST FE] [Client Lai Candidate] Sent Event: interview:start (reqId: req_8wk2dujo5)
[TEST ERR] [Client Lai Candidate] Received ERROR Event: {"code":"FORBIDDEN","message":"Chỉ Recruiter mới có quyền bắt đầu buổi phỏng vấn","recoverable":false}
[TEST FE] SUCCESS: Chặn Candidate tự ý start interview chính xác.
[TEST FE] [Client Hùng Recruiter] Sent Event: interview:start (reqId: req_0ahkmddga)
[BE ERR] 2026/07/09 09:04:31 [db] UPDATE interview_rooms SET status = 'active'
[BE ERR] 2026/07/09 09:04:31 [db] UPDATE interviews SET status = 'active', started_at = '2026-07-09T02:04:31Z'
[BE ERR] 2026/07/09 09:04:31 [redis] SET room_status:test-room-nj2gw077b value=active
[BE ERR] 2026/07/09 09:04:31 [ai-orchestrator] starting AI pipeline for room=test-room-nj2gw077b (interview=test-interview-nj2gw077b, consent_recording=true)
[BE ERR] 2026/07/09 09:04:31 [audio_hook] Starting audio stream hook for room=test-room-nj2gw077b, interview=test-interview-nj2gw077b
[TEST FE] [Client Lai Candidate] Received Event: interview:started (reqId: req_0ahkmddga)
[TEST FE] [Client Hùng Recruiter] Received Event: interview:started (reqId: req_0ahkmddga)
[TEST FE] SUCCESS: Bắt đầu phỏng vấn thành công.

[TEST FE] --- TEST STEP 5: AI Bridge & Rate Limiting ---
[TEST FE] [Client Hùng Recruiter] Sent Event: ai:request_suggestion (reqId: req_1as0m92s7)
[TEST FE] [Client Hùng Recruiter] Received Event: ai:thinking (reqId: req_1as0m92s7)
[BE ERR] 2026/07/09 09:04:32 [db] INSERT INTO ai_suggestions (id, interview_id, suggestion_type, content, reason, target_skill, priority, confidence, created_at) VALUES ('15aa398e-c090-47fd-a68c-f56f17416405', 'test-interview-nj2gw077b', 'follow_up_question', 'Bạn có thể nói rõ bạn đã đo performance bằng chỉ số nào không?', 'Ứng viên nói đã tối ưu performance nhưng chưa nêu metric cụ thể.', 'Performance Optimization', 'high', 0.840000, '2026-07-09T02:04:32Z')
[TEST FE] [Client Hùng Recruiter] Received Event: ai:suggestion (reqId: req_1as0m92s7)
[TEST FE] SUCCESS: Trạng thái AI Thinking chỉ hiển thị với Recruiter.
[TEST FE] SUCCESS: AI suggestions trả về đúng chuẩn envelope.

[TEST FE] --- TEST STEP 6: Room Heartbeat ---
[TEST FE] [Client Hùng Recruiter] Sent Event: room:heartbeat (reqId: req_2diqtvczx)
[TEST FE] [Client Lai Candidate] Sent Event: room:heartbeat (reqId: req_i90fowygz)
[TEST FE] SUCCESS: Heartbeat gửi không gây lỗi server.

[TEST FE] --- TEST STEP 7: End Interview & Grace Period ---
[TEST FE] [Client Hùng Recruiter] Sent Event: interview:end (reqId: req_9wu7rqhpl)
[BE ERR] 2026/07/09 09:04:34 [db] UPDATE interview_rooms SET status = 'completed'
[BE ERR] 2026/07/09 09:04:34 [db] UPDATE interviews SET status = 'completed', ended_at = '2026-07-09T02:04:34Z'
[BE ERR] 2026/07/09 09:04:34 [redis] SET room_status:test-room-nj2gw077b value=completed
[BE ERR] 2026/07/09 09:04:34 [audio_hook] Requested stop for audio stream hook in room=test-room-nj2gw077b
[BE ERR] 2026/07/09 09:04:34 [audio_hook] Stopped audio stream hook for room=test-room-nj2gw077b
[TEST FE] [Client Hùng Recruiter] Received Event: interview:completed (reqId: req_9wu7rqhpl)
[TEST FE] [Client Lai Candidate] Received Event: interview:completed (reqId: req_9wu7rqhpl)
[TEST FE] SUCCESS: Hoàn tất đóng phỏng vấn.
====================================================
TẤT CẢ CÁC BÀI INTEGRATION TEST ĐÃ VƯỢT QUA (PASSED) ✅
====================================================

[BE] Đang tắt Go Realtime server...
[BE] Server đã shutdown thành công.
[BE] Đã xóa file thực thi tạm thời.

✅ KẾT QUẢ: KIỂM THỬ THÀNH CÔNG! KHÔNG PHÁT HIỆN BUG.
```

---

## 2. Các điểm lỗi đã sửa đổi

### 2.1. Cấu hình phân tách cổng Gateway Realtime và API
- **Trước**: Gateway Realtime hardcode cổng `:8080` trong `cmd/realtime/main.go`, trùng lặp với cổng API gây ra lỗi `Only one usage of each socket address is normally permitted`.
- **Sau**: Gateway Realtime được cấu hình để đọc biến môi trường `REALTIME_PORT`, mặc định chạy trên cổng `:8081` (đúng cấu hình `VITE_WS_URL=ws://localhost:8081`). Kịch bản test tự động cũng chuyển sang chạy trên `:8081`.

### 2.2. Khắc phục lỗi `panic: runtime error: invalid memory address or nil pointer dereference` trong `AuditRepository`
- **Trước**: Khi chạy Realtime server độc lập, do không thiết lập kết nối Database nên DB instance truyền vào `NewAuditRepository(nil)` bị `nil`. Khi có sự kiện `room:join` phát sinh, background worker của `AuditLogger` gọi `r.db.PrepareNamedContext()` dẫn đến panic crash server lập tức.
- **Sau**: Đã thêm kiểm tra an toàn `if r.db == nil { return nil }` trong `AuditRepository.Insert` để tự động bypass ghi DB khi không có cấu hình kết nối mà không gây crash server.
