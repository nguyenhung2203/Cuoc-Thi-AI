# H-S0-03 — Xác thực WebSocket bằng Token

| Field | Value |
|---|---|
| **Task ID** | H-S0-03 |
| **Sprint** | 0 — Setup Realtime Foundation |
| **Độ khó** | Trung bình |
| **Trạng thái** | ⬜ Chưa bắt đầu |
| **Owner** | Hùng |

---

## Mục tiêu

Xác thực user khi kết nối WebSocket bằng `room_access_token`, đảm bảo:
- Chỉ user có token hợp lệ mới connect được
- Token sai/hết hạn bị reject ngay
- Trích xuất user info (user_id, role, room_id) từ token

---

## Yêu cầu chi tiết

### 1. Token Flow

**Recruiter:**
```
POST /api/v1/companies/:company_id/interviews/:interview_id/room/token
→ Nhận room_access_token (LiveKit JWT)
→ Dùng token connect WebSocket
```

**Candidate:**
```
GET /api/v1/interviews/join/:invite_token
→ Nhận room_access_token (LiveKit JWT, TTL 4 giờ)
→ Dùng token connect WebSocket
```

### 2. Authentication khi Connect

Hỗ trợ 2 cách gửi token:

**Option 1 — Query string:**
```
ws://localhost:8080/ws/interview-room?token=<room_access_token>
```

**Option 2 — Header:**
```
Authorization: Bearer <room_access_token>
```

### 3. Validation Logic

```go
func ValidateRoomToken(token string) (*TokenClaims, error) {
    // 1. Parse JWT
    // 2. Verify signature
    // 3. Check expiration
    // 4. Extract claims: user_id, room_id, role, participant_type
    // 5. Return claims hoặc error
}
```

### 4. Reject Cases

| Trường hợp | Response |
|---|---|
| Không có token | HTTP 401 — close connection |
| Token format sai | HTTP 401 — close connection |
| Token hết hạn | HTTP 401 — close connection |
| Token không match room | HTTP 403 — close connection |

---

## Tài liệu tham chiếu

- [REALTIME_EVENTS.md](../../REALTIME_EVENTS.md) — Mục 2.2 (Authentication)
- [API_SPEC.md](../../API_SPEC.md) — Mục 6.7 (Join by invite token), Mục 7.2 (Room token)

---

## Definition of Done

- [ ] Token hợp lệ → connect thành công, trích xuất user info
- [ ] Token sai → bị reject với HTTP 401
- [ ] Token hết hạn → bị reject
- [ ] Không có token → bị reject
- [ ] User info (user_id, role, room_id) được gắn vào connection
- [ ] Log authentication success/failure

---

## Dependency

- **H-S0-02** — Cần WebSocket gateway skeleton
- **Khôi** — Cần JWT signing key và token format thống nhất

---

## Files dự kiến tạo/sửa

```
backend/
├── internal/realtime/
│   ├── auth.go                    # Token validation
│   ├── middleware.go              # Auth middleware cho WS upgrade
│   └── handler.go                # Cập nhật: thêm auth check
├── internal/auth/
│   └── jwt.go                    # Shared JWT utils (nếu cần)
```

---

## Checklist test

- [ ] Connect với token hợp lệ → thành công
- [ ] Connect không có token → bị reject 401
- [ ] Connect với token hết hạn → bị reject 401
- [ ] Connect với token giả → bị reject 401
- [ ] Connection có user_id và role đúng sau auth
- [ ] Log ghi nhận auth success/failure
