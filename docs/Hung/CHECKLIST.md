# CHECKLIST — Module Hùng: Realtime & Interview Room

> **Owner**: Hùng
> **Phạm vi**: Interview Room, Realtime Gateway, WebSocket, LiveKit WebRTC, Presence, Chat, Transcript Realtime, AI Bridge
> **Tài liệu gốc**: [MODULE_HUNG_REALTIME_AI.md](./MODULE_HUNG_REALTIME_AI.md)

---

## Sprint 0: Setup Realtime Foundation

- [x] **H-S0-01** — Thiết kế event contract realtime → [task/H-S0-01_event_contract.md](./task/H-S0-01_event_contract.md)
  - `events/event_types.go` — tất cả event name constants, enums (ParticipantType, RoomStatus, SpeakerType, AIErrorType, …)
  - `events/event_envelope.go` — Envelope struct, NewEnvelope, ParsePayload, ToJSON
  - `events/room_events.go` — RoomJoinPayload, RoomJoinedPayload, RoomUserJoinedPayload, RoomUserLeftPayload, RoomPresenceUpdatePayload, ParticipantInfo
  - `events/interview_events.go` — InterviewStartPayload, InterviewEndPayload, InterviewStartedPayload, InterviewCompletedPayload
  - `events/chat_events.go` — ChatSendPayload, NoteCreatePayload, QuestionMarkAskedPayload, ChatMessagePayload
  - `events/media_events.go` — MediaStatusPayload, MediaStatusChangedPayload
  - `events/transcript_events.go` — TranscriptPartialPayload, TranscriptFinalPayload, TranscriptUpdatePayload
  - `events/ai_events.go` — AIThinkingPayload, AISuggestionPayload, AIScoreUpdatePayload, AIWarningPayload, AIErrorPayload, ReportReadyPayload, ErrorPayload
  - `events/visibility.go` — VisibilityRule, ShouldSendToParticipant, ShouldSendChatToParticipant

- [x] **H-S0-02** — Tạo WebSocket gateway skeleton → [task/H-S0-02_websocket_gateway.md](./task/H-S0-02_websocket_gateway.md)
  - `server.go` — Server struct, Start, Shutdown, handleInternal stub
  - `handler.go` — handleUpgrade, readPump, writePump, onDisconnect, sendError
  - `connection.go` — ClientConnection struct, WriteMessage (nil-safe), Close
  - `connection_manager.go` — Add/Remove/Get/Count/ByRoom/CloseAll, newClientConnection
  - `message_router.go` — MessageRouter, register, Route, stub handlers cho 14 events
  - `config.go` — writeWait, pongWait, pingPeriod, maxMessageSize, sendBufferSize, grace periods, rate limits
  - `presence.go` — PresenceManager, OnHeartbeat, StartWatcher, broadcastPresence
  - `room.go` — Room struct, AddParticipant, BroadcastAll/Except/RecruitersOnly/WithVisibility/SendTo
  - `room_manager.go` — RoomManager, GetOrCreate, OnParticipantDisconnect, SetRoomStatus, isValidTransition
  - `participant.go` — Participant struct (ConnectionState, MediaStatus, LastSeenAt), ToInfo
  - 12 tests PASS: 6 unit (connection_manager_test.go) + 6 integration (websocket_integration_test.go)

- [x] **H-S0-03** — Xác thực WebSocket bằng token → [task/H-S0-03_websocket_auth.md](./task/H-S0-03_websocket_auth.md)
  - `auth.go` — TokenClaims struct, extractAndValidateToken, tokenFromRequest (query + header), validateRoomToken stub (LiveKit JWT — TODO H-S0-03 full impl)

- [x] **H-S0-04** — Room channel architecture → [task/H-S0-04_room_channel.md](./task/H-S0-04_room_channel.md)
  - `room_handler.go` — handleRoomJoin, handleRoomLeave, handleRoomHeartbeat
  - `message_router.go` — register handlers, add PresenceManager
  - `handler.go` — update onDisconnect cleanup
  - `room_channel_test.go` — join, presence, visibility, leave/cleanup, multi-room, concurrency tests

---

## Sprint 1: Interview Room MVP

- [x] **H-S1-01** — Join/Leave room → [task/H-S1-01_join_leave_room.md](./task/H-S1-01_join_leave_room.md)
  - `auth.go` — restrict guest roles
  - `room_handler.go` — room ID validation, simulated DB logging for interview_participants
  - `handler.go` — simulated DB log for disconnects
- [x] **H-S1-02** — Presence realtime → [task/H-S1-02_presence_realtime.md](./task/H-S1-02_presence_realtime.md)
  - `presence.go` — start/stop background watcher registry, reconnecting & offline state calculations, simulated cache/DB updates
  - `room_handler.go` — manage room presence watcher lifecycle, route room:heartbeat
  - `config.go` — set heartbeat thresholds and reconnecting grace period
  - `room_channel_test.go` — verify presence transition states (online → reconnecting → offline → online)
- [x] **H-S1-03** — Start/End interview event → [task/H-S1-03_start_end_interview.md](./task/H-S1-03_start_end_interview.md)
  - `interview_handler.go` [NEW] — recruiter-only start/end handlers, DB/AI/report logging, grace-period socket closure and room cleanup
  - `message_router.go` — wire start/end routes to actual handler functions
  - `room.go` — add processed request IDs and RecordRequest helper to enforce idempotency
  - `room_channel_test.go` — verify recruiter start, candidate forbidden check, recruiter end, idempotency, and room teardown
- [x] **H-S1-04** — Room status sync → [task/H-S1-04_room_status_sync.md](./task/H-S1-04_room_status_sync.md)
  - `event_types.go` & `interview_events.go` — add pause, resume, cancel, and room expiration events and payloads
  - `room_manager.go` — simulate DB/Redis status persistence and restore status on room recreation/reload
  - `interview_handler.go` & `message_router.go` — recruiter-only pause, resume, cancel handlers, sync status with DB/Redis
  - `presence.go` — check waiting room auto-expiration and trigger room:expired with connection cleanup
  - `room_status_test.go` [NEW] — verify room state transitions, invalid transition rejection, persistence on reload, and auto-expiry

---

## Sprint 2: Chat + LiveKit WebRTC

- [x] **H-S2-01** — Chat realtime trong room → [task/H-S2-01_chat_realtime.md](./task/H-S2-01_chat_realtime.md)
  - `chat_handler.go` [NEW] — handle chat:send, DB INSERT query logging, broadcast based on visibility payload
  - `message_router.go` — register handleChatSend route
  - `room_manager.go` — add simulatedChat map, implement thread-safe SaveChatMessage and GetChatHistory
  - `server.go` — implement GET /api/v1/rooms/{room_id}/chat secure REST history endpoint with recruiter-only filtering
  - `chat_realtime_test.go` [NEW] — verify room visibility, recruiter-only visibility, and REST history endpoint logic
- [x] **H-S2-02** — LiveKit WebRTC integration → [task/H-S2-02_livekit_webrtc.md](./task/H-S2-02_livekit_webrtc.md)
  - `token.go` — token generator signed with HS256 containing room access claims and map grants
  - `auth.go` — add extractAndValidateRecruiterToken helper to validate recruiter token prior to room generation
  - `server.go` — register POST recruiter token endpoint and GET candidate invite join endpoint, generating and returning LiveKit tokens
  - `livekit_integration_test.go` [NEW] — verify token generation, recruiter endpoint token creation, and candidate join invite verification/exchange
- [x] **H-S2-03** — Audio stream hook cho AI → [task/H-S2-03_audio_stream_ai.md](./task/H-S2-03_audio_stream_ai.md)
  - `audio_hook.go` — thiết lập `AudioHookService` mô phỏng capture audio gửi cho Python AI Orchestrator
  - `server.go` — thêm `handleInternal` nhận webhook transcript từ AI và broadcast qua websocket
  - `message_router.go` — khởi tạo `AudioHookService` tại router
  - `interview_handler.go` — tích hợp `StartHook` vào `interview:start` và `StopHook` vào `interview:end`/`cancel`
  - `audio_hook_integration_test.go` [NEW] — verify webhook integration và logic broadcast transcript qua bài test (chạy bằng `go test`)
- [x] **H-S2-04** — Media status event → [task/H-S2-04_media_status.md](./task/H-S2-04_media_status.md)
  - `participant.go` — thêm `MediaRateLimiter` để xử lý rate limit
  - `room_handler.go` — khởi tạo rate limiter trong sự kiện join room
  - `room.go` — thêm method `UpdateMediaStatus` thread-safe
  - `message_router.go` — đăng ký sự kiện `media:status`
  - `media_handler.go` [NEW] — xử lý cập nhật trạng thái media, rate limiting và broadcast cho cả phòng
  - `media_status_test.go` [NEW] — verify broadcast `media:status_changed`, client join muộn nhận status cũ, spam block

---

## Sprint 3: Transcript Realtime

- [x] **H-S3-01** — Transcript event pipeline → [task/H-S3-01_transcript_pipeline.md](./task/H-S3-01_transcript_pipeline.md)
  - `transcript_pipeline.go` [NEW] — xây dựng `TranscriptPipeline` async buffer queue (cap 1000) giúp decouple HTTP worker và WebSocket broadcast
  - `server.go` — tích hợp `TranscriptPipeline` vào endpoint `POST /internal/rooms/:room_id/transcript`
  - `transcript_pipeline_test.go` [NEW] — kiểm thử tích hợp xác nhận recruiter luôn nhận transcript, candidate nhận/không nhận theo cấu hình room

- [x] **H-S3-02** — Partial/Final transcript handling → [task/H-S3-02_partial_final_transcript.md](./task/H-S3-02_partial_final_transcript.md)
  - `message_router.go` — đăng ký WebSocket handler cho `transcript:partial` và `transcript:final`
  - `transcript_handler.go` [NEW] — xử lý tạo deterministic `TranscriptID` (phục vụ dedup partial→final), ghi log cảnh báo confidence < 0.5 và đẩy vào pipeline
  - `partial_final_transcript_test.go` [NEW] — kiểm thử tự động xác minh partial cập nhật thay thế nhau, final chốt text không duplicate và cảnh báo confidence thấp

- [x] **H-S3-03** — Speaker mapping → [task/H-S3-03_speaker_mapping.md](./task/H-S3-03_speaker_mapping.md)
  - `transcript_events.go` — thêm `participant_id`, `track_id`, `identity`, `speaker_name` vào struct payload
  - `room.go` — thêm `TrackParticipantMap`, `RegisterTrack`, `ResolveSpeaker`
  - `transcript_handler.go` & `server.go` — phân giải speaker mapping cho cả sự kiện WebSocket và Webhook push
  - `transcript_test.go` [NEW] — kiểm thử tự động `ResolveSpeaker`, fallback và role visibility
- [x] **H-S3-04** — Save transcript final → [task/H-S3-04_save_transcript.md](./task/H-S3-04_save_transcript.md)
  - `transcript_saver.go` [NEW] — tạo struct `TranscriptRecord` & `TranscriptBatchSaver` async queue gom lô mỗi 5s hoặc 100 items
  - `message_router.go` — gắn `transcriptSaver` vào router, cung cấp getter
  - `transcript_handler.go`, `chat_handler.go` & `server.go` — đẩy bản ghi final vào queue gom lô; đóng queue an toàn khi Shutdown server
  - `transcript_saver_test.go` [NEW] — bộ kiểm thử tự động kiểm chứng tra cứu thứ tự thời gian, gộp lô định kỳ/dung lượng và chặn partial transcript

---

## Sprint 4: AI Realtime Bridge

- [x] **H-S4-01** — AI suggestion event → [task/H-S4-01_ai_suggestion.md](./task/H-S4-01_ai_suggestion.md)
  - `ai_handler.go` [NEW] — tạo `AIRateLimiter` (cap 10 reqs/10m) & `handleAIRequestSuggestion` gửi `ai:thinking` lập tức, phát gợi ý cho riêng Recruiter
  - `message_router.go` — tích hợp `aiRateLimiter` & đăng ký route handler
  - `server.go` — thêm HTTP Webhook `handleAISuggestionPush` tiếp nhận từ AI Orchestrator ngoài
  - `ai_suggestion_test.go` [NEW] — bộ kiểm thử tự động xác minh phân quyền role, rate limit và phản hồi đúng envelope chuẩn
- [x] **H-S4-02** — AI score update event → [task/H-S4-02_ai_score_update.md](./task/H-S4-02_ai_score_update.md)
  - `ai_handler.go` — thêm `handleAIRequestScoreUpdate` chấm điểm Rubric kèm bằng chứng (`evidence`), tự động gán `insufficient_evidence` khi thiếu dữ liệu
  - `message_router.go` — tích hợp `scoreRateLimiter` (10 reqs/10m) & định tuyến websocket event
  - `server.go` — bổ sung Webhook intake `handleAIScoreUpdatePush` tiếp nhận điểm từ Python AI Orchestrator
  - `ai_score_update_test.go` [NEW] — kiểm thử tự động xác minh phân quyền Recruiter only và kiểm chứng bằng chứng chấm điểm
- [x] **H-S4-03** — Role-based event visibility → [task/H-S4-03_role_visibility.md](./task/H-S4-03_role_visibility.md)
  - `visibility.go` — bổ sung quy tắc `EventNoteCreate: VisibleRecruitersOnly`
  - `ai_handler.go` & `server.go` — chuyển toàn bộ lệnh broadcast sang bộ lọc phân quyền tập trung `BroadcastWithVisibility`
  - `role_visibility_test.go` [NEW] — bộ kiểm thử tự động xác minh ma trận phân quyền 7 loại sự kiện giữa Recruiter/Candidate cùng ngoại lệ Mock Interview
- [x] **H-S4-04** — AI error event handling → [task/H-S4-04_ai_error_handling.md](./task/H-S4-04_ai_error_handling.md)
  - `ai_handler.go` & `message_router.go` — xử lý gửi `ai:error`, theo dõi số lần thử lại (`aiRetries`), phân biệt lỗi tạm thời/vĩnh viễn và tự động thử lại
  - `server.go` — thêm Webhook intake `handleAIErrorPush`
  - `ai_error_handling_test.go` [NEW] — kiểm thử tự động xác minh gửi lỗi degraded/critical, tự thử lại, ẩn với Candidate và duy trì hoạt động phòng/chat/end interview khi AI down

---

## Sprint 5: Hardening

- [x] **H-S5-01** — Reconnect recovery → [task/H-S5-01_reconnect_recovery.md](./task/H-S5-01_reconnect_recovery.md)
  - `room.go` — Thêm `EventHistory` circular buffer (cap 500) lưu trữ sự kiện phát lại, hàm lọc `GetMissedEvents` chuẩn role visibility
  - `room_events.go` — Mở rộng `RoomJoinedPayload` với `MediaStatus`, `MissedEventsCount`, `SyncFromTimestamp`
  - `handler.go` & `presence.go` — Xử lý ngắt kết nối chuyển sang `reconnecting` thay vì `left`, watcher dọn dẹp sau `offline threshold`
  - `room_handler.go` — Khôi phục `Participant` trong `handleRoomJoin`, phát lại `missed events` lập tức cho client tái kết nối
  - `reconnect_recovery_test.go` [NEW] — Kiểm thử tự động chứng minh giữ kết nối ngắn hạn, ẩn với người khác, timeout chuyển offline và khôi phục backlog chat/transcript
- [x] **H-S5-02** — Grace period disconnect → [task/H-S5-02_grace_period.md](./task/H-S5-02_grace_period.md)
  - `config.go` — Chuyển `gracePeriodRoom`, `gracePeriodShort` sang biến `var` hỗ trợ can thiệp thời gian kiểm thử tự động tốc độ cao `SetGracePeriodsForTest`
  - `room.go` — Thêm hàm kiểm tra nhanh `HasOnlineParticipants` xác định phòng không còn kết nối trực tuyến
  - `room_manager.go` — Quản lý `emptyTimers` map, triển khai goroutine hẹn giờ dọn dẹp phòng `StartEmptyRoomTimer`, huỷ hẹn giờ `CancelEmptyRoomTimer` và thực thi dọn bộ nhớ/DB/Redis `CleanupExpiredRoom`
  - `handler.go` & `room_handler.go` — Khởi động room grace timer khi người tham gia cuối cùng ngắt kết nối/leave; huỷ hẹn giờ ngay lập tức khi có người tham gia mới/reconnect joined room
  - `grace_period_test.go` [NEW] — Bộ kiểm thử tự động chứng minh phòng duy trì trạng thái khi disconnect tạm thời, dọn sạch khi hết hạn, huỷ dọn dẹp khi reconnect và kiểm tra tải 100 phòng không rò rỉ bộ nhớ/goroutines
- [x] **H-S5-03** — Room event audit log → [task/H-S5-03_room_audit.md](./task/H-S5-03_room_audit.md)
  - `models/audit_log.go` — Cập nhật schema đầy đủ các field.
  - `repository/audit_repo.go` — Thêm hàm `Insert` mô phỏng DB.
  - `audit_logger.go` [NEW] — Async logger sử dụng buffered channel để không chặn luồng chính.
  - `server.go` & `handler.go` — Lấy `IPAddress` từ request HTTP, lưu log `disconnect`.
  - `room_handler.go`, `interview_handler.go`, `ai_handler.go` — Gắn log các event tương ứng.
- [ ] **H-S5-04** — Load test nhiều room → [task/H-S5-04_load_test.md](./task/H-S5-04_load_test.md)

---

## Tổng kết tiến độ

| Sprint | Tổng task | Hoàn thành | Trạng thái |
|---|---:|---:|---|
| Sprint 0: Foundation | 4 | 4 | ✅ Hoàn thành |
| Sprint 1: Room MVP | 4 | 4 | ✅ Hoàn thành |
| Sprint 2: Chat + WebRTC | 4 | 4 | ✅ Hoàn thành |
| Sprint 3: Transcript | 4 | 4 | ✅ Hoàn thành |
| Sprint 4: AI Bridge | 4 | 4 | ✅ Hoàn thành |
| Sprint 5: Hardening | 4 | 3 | 🔄 Đang thực hiện |
| **Tổng** | **24** | **23** | |

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
6. Mỗi task xong → cập nhật CHECKLIST.md + ghi file đã sửa + checklist test
