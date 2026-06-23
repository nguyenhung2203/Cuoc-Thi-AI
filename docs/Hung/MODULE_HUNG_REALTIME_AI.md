# MODULE_HUNG_REALTIME_AI.md
# Phân công Module cho Hùng

## 1. Vai trò của Hùng

Hùng phụ trách nhóm module **khó và phức tạp nhất liên quan đến realtime, phòng phỏng vấn, đồng bộ trạng thái, transcript live và kết nối AI realtime**.

Hùng là owner chính của các module:

1. Interview Room Core.
2. Realtime Gateway.
3. WebSocket Event System.
4. WebRTC Signaling.
5. Room Presence.
6. Chat trong phòng.
7. Transcript realtime.
8. AI Suggestion realtime bridge.
9. AI Score realtime bridge.
10. Reconnect và recovery state.

---

## 2. Mục tiêu module của Hùng

Mục tiêu chính:

- Recruiter và Candidate vào được cùng một phòng phỏng vấn.
- Phòng có trạng thái rõ ràng: waiting, active, paused, completed, cancelled.
- Có realtime event cho join/leave/chat/start/end/transcript/AI suggestion.
- Có cơ chế reconnect khi mất mạng.
- Có phân quyền event theo role.
- Có nền tảng để Khôi gắn AI engine vào.
- Có nền tảng để Lai gắn UI room vào.

---

## 3. Phạm vi công việc

## 3.1. Interview Room Core

### Chức năng

- Tạo room cho mỗi interview.
- Join room bằng interview link.
- Kiểm tra quyền trước khi vào room.
- Lưu trạng thái room.
- Start interview.
- End interview.
- Cancel interview.
- Auto expire room nếu quá hạn.

### Trạng thái room

| Status | Ý nghĩa |
|---|---|
| scheduled | Đã lên lịch |
| waiting | Đang chờ người vào |
| candidate_joined | Candidate đã vào |
| recruiter_joined | Recruiter đã vào |
| active | Đang phỏng vấn |
| paused | Tạm dừng |
| completed | Đã kết thúc |
| cancelled | Đã hủy |
| expired | Link hết hạn |

### DoD

- Room tạo được từ interview.
- Recruiter và Candidate join đúng quyền.
- Không cho user lạ vào room.
- Start/end chỉ Recruiter được thao tác.
- Room status lưu đúng.
- Có event realtime khi status đổi.

---

## 3.2. Realtime Gateway

### Chức năng

- Quản lý WebSocket connection.
- Xác thực token khi connect.
- Map user vào room.
- Broadcast event theo room.
- Gửi event riêng theo role.
- Handle disconnect.
- Handle reconnect.

### Event format chuẩn

```json
{
  "event": "room:join",
  "room_id": "uuid",
  "sender_id": "uuid",
  "sender_role": "recruiter",
  "payload": {},
  "timestamp": "2026-06-23T10:00:00Z"
}
```

### DoD

- Connect WebSocket thành công.
- Reject connection nếu token sai.
- Join room đúng.
- Broadcast không gửi nhầm room.
- Không gửi dữ liệu nội bộ cho Candidate.
- Có log lỗi connection.

---

## 3.3. WebRTC Signaling

### Chức năng

- Gửi offer.
- Gửi answer.
- Gửi ICE candidate.
- Sync trạng thái mic/camera.
- Cho phép reconnect media.

### Event

| Event | Mô tả |
|---|---|
| webrtc:offer | Gửi offer |
| webrtc:answer | Gửi answer |
| webrtc:ice_candidate | Gửi ICE candidate |
| media:status | Cập nhật mic/camera |
| media:error | Lỗi media |

### DoD

- Recruiter và Candidate call được peer-to-peer ở MVP.
- Trạng thái mic/camera sync đúng.
- Reload có thể reconnect lại.
- Nếu WebRTC lỗi, room không crash.

---

## 3.4. Presence trong phòng

### Chức năng

- Hiển thị ai đang online.
- Hiển thị ai rời phòng.
- Hiển thị ai reconnect.
- Lưu last_seen.
- Phân biệt recruiter/candidate.

### Event

- room:user_joined
- room:user_left
- room:user_reconnected
- room:presence_update

### DoD

- Join/leave update realtime.
- Mất mạng tạm thời không end room ngay.
- Có grace period cho reconnect.
- Presence không bị nhân đôi user.

---

## 3.5. Chat trong phòng

### Chức năng

- Recruiter và Candidate gửi tin nhắn.
- Lưu chat message.
- Broadcast realtime.
- Có timestamp.
- Có sender role.
- Có trạng thái gửi thành công/thất bại.

### DoD

- Gửi/nhận chat realtime.
- Lưu chat vào DB.
- Reload vẫn xem được lịch sử chat.
- Không cho người ngoài room gửi chat.

---

## 3.6. Transcript realtime

### Chức năng

- Nhận transcript từ STT service hoặc AI service của Khôi.
- Broadcast transcript update vào room.
- Phân biệt speaker.
- Có confidence.
- Có partial/final transcript.
- Lưu transcript theo interview.

### Event

```json
{
  "event": "transcript:update",
  "room_id": "uuid",
  "payload": {
    "transcript_id": "uuid",
    "speaker_type": "candidate",
    "content": "Tôi có 2 năm kinh nghiệm làm Vue...",
    "is_final": false,
    "confidence": 0.86,
    "start_time": 12.4,
    "end_time": 16.9
  }
}
```

### DoD

- Transcript update realtime.
- Không bị duplicate quá nhiều.
- Có speaker rõ ràng.
- Có partial và final.
- Lưu được transcript final.

---

## 3.7. AI Suggestion realtime bridge

### Chức năng

- Nhận suggestion từ AI engine của Khôi.
- Gửi suggestion cho Recruiter.
- Không gửi suggestion nội bộ cho Candidate.
- Có trạng thái loading/thinking/done/error.

### Event

- ai:suggestion_loading
- ai:suggestion
- ai:suggestion_error

### DoD

- Recruiter nhận được gợi ý realtime.
- Candidate không thấy gợi ý nội bộ.
- Nếu AI lỗi, UI vẫn hoạt động.
- Có retry event nếu cần.

---

## 3.8. AI Score realtime bridge

### Chức năng

- Nhận score update từ AI engine.
- Gửi score panel cho Recruiter.
- Không gửi score tuyển dụng cho Candidate.
- Có confidence/evidence nếu có.

### Event

- ai:score_update
- ai:score_error

### DoD

- Score update đúng room.
- Chỉ Recruiter thấy.
- Không block room nếu scoring lỗi.

---

## 4. Sprint chia việc cho Hùng

## Sprint 0: Setup realtime foundation

| Task | Mô tả | Độ khó | DoD |
|---|---|---:|---|
| H-S0-01 | Thiết kế event contract realtime | Khó | Có danh sách event, payload, role visibility |
| H-S0-02 | Tạo WebSocket gateway skeleton | Khó | Client connect/disconnect được |
| H-S0-03 | Xác thực WebSocket bằng token | Trung bình | Token sai bị reject |
| H-S0-04 | Room channel architecture | Khó | Broadcast theo room đúng |

## Sprint 1: Interview room MVP

| Task | Mô tả | Độ khó | DoD |
|---|---|---:|---|
| H-S1-01 | Join/leave room | Khó | Recruiter/Candidate join đúng |
| H-S1-02 | Presence realtime | Khó | Online/offline/reconnect update đúng |
| H-S1-03 | Start/end interview event | Trung bình | Chỉ Recruiter được start/end |
| H-S1-04 | Room status sync | Khó | UI nhận đúng trạng thái |

## Sprint 2: Chat + WebRTC signaling

| Task | Mô tả | Độ khó | DoD |
|---|---|---:|---|
| H-S2-01 | Chat realtime trong room | Trung bình | Gửi/nhận/lưu chat được |
| H-S2-02 | WebRTC signaling offer/answer | Rất khó | 2 peer connect được |
| H-S2-03 | ICE candidate exchange | Rất khó | Media ổn định ở local/dev |
| H-S2-04 | Media status event | Trung bình | Mic/camera sync đúng |

## Sprint 3: Transcript realtime

| Task | Mô tả | Độ khó | DoD |
|---|---|---:|---|
| H-S3-01 | Transcript event pipeline | Rất khó | Nhận/broadcast transcript được |
| H-S3-02 | Partial/final transcript handling | Khó | Không duplicate final text |
| H-S3-03 | Speaker mapping | Khó | Biết ai đang nói |
| H-S3-04 | Save transcript final | Trung bình | Lưu DB đúng interview |

## Sprint 4: AI realtime bridge

| Task | Mô tả | Độ khó | DoD |
|---|---|---:|---|
| H-S4-01 | AI suggestion event | Khó | Recruiter nhận gợi ý |
| H-S4-02 | AI score update event | Khó | Recruiter nhận score |
| H-S4-03 | Role-based event visibility | Khó | Candidate không thấy data nội bộ |
| H-S4-04 | AI error event handling | Trung bình | AI lỗi không crash room |

## Sprint 5: Hardening

| Task | Mô tả | Độ khó | DoD |
|---|---|---:|---|
| H-S5-01 | Reconnect recovery | Rất khó | Mất mạng vào lại được |
| H-S5-02 | Grace period disconnect | Khó | Không kick ngay khi mạng chập chờn |
| H-S5-03 | Room event audit | Trung bình | Log sự kiện quan trọng |
| H-S5-04 | Load test nhiều room | Khó | Có kết quả test cơ bản |

---

## 5. Dependency với Khôi

Hùng cần Khôi cung cấp:

- API tạo/lấy interview.
- DB schema room/interview/transcript.
- AI suggestion service.
- AI scoring service.
- Report generation trigger.
- Auth/RBAC middleware.

Các điểm cần thống nhất:

- Event payload cho transcript.
- Event payload cho AI suggestion.
- Event payload cho AI score.
- Cách lưu transcript.
- Cách mapping user/participant trong room.

---

## 6. Dependency với Lai

Hùng cần Lai cung cấp:

- UI Interview Room.
- Component chat.
- Component transcript panel.
- Component AI suggestion panel.
- Component score panel.
- Trạng thái loading/error trên UI.

Các điểm cần thống nhất:

- Client WebSocket wrapper.
- Cách frontend subscribe/unsubscribe event.
- Role visibility trên UI.
- UI reconnect state.
- UI báo lỗi media.

---

## 7. Checklist test riêng của Hùng

- Recruiter join room thành công.
- Candidate join room thành công.
- User không có quyền bị chặn.
- Recruiter start interview được.
- Candidate không start interview được.
- Chat realtime hoạt động.
- Presence hoạt động.
- Reload room không mất trạng thái nghiêm trọng.
- Disconnect/reconnect hoạt động.
- Transcript update đúng.
- AI suggestion chỉ gửi Recruiter.
- AI score chỉ gửi Recruiter.
- End interview đóng room đúng.
- AI service lỗi nhưng room vẫn dùng được.

---

## 8. Prompt mẫu cho Hùng dùng với AI coding assistant

```text
Bạn đang làm module của Hùng trong dự án phỏng vấn cùng AI real-time.
Phạm vi của Hùng là Interview Room, Realtime Gateway, WebSocket/WebRTC, Presence, Chat, Transcript realtime và AI realtime bridge.

Yêu cầu task hiện tại:
[Điền task]

Trước khi code hãy:
1. Đọc các file realtime/interview hiện có.
2. Xác định event contract liên quan.
3. Không sửa AI scoring/report schema của Khôi nếu không cần.
4. Không redesign UI của Lai nếu không cần.
5. Đảm bảo role visibility: Recruiter thấy AI suggestion/score, Candidate không thấy.
6. Có xử lý reconnect/error.
7. Sau khi code, liệt kê file đã sửa và checklist test.
```

---

## 9. Kết quả cuối cùng Hùng cần bàn giao

Hùng hoàn thành khi hệ thống có:

- Phòng phỏng vấn realtime hoạt động.
- Recruiter/Candidate join được.
- Presence ổn định.
- Chat realtime hoạt động.
- WebRTC signaling cơ bản hoạt động.
- Transcript realtime hoạt động.
- AI suggestion/score realtime bridge hoạt động.
- Reconnect không làm hỏng room.
- Role-based event visibility an toàn.
- Có checklist test rõ ràng.
