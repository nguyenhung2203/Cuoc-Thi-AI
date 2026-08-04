# Hướng dẫn triển khai production — ViecLamAI

Tài liệu này hướng dẫn triển khai hệ thống trên VPS Ubuntu 22.04 bằng Docker Compose.

## 1. Thông tin triển khai hiện tại

- Hệ điều hành: Ubuntu 24.04.4 LTS x64 (AWS)
- Tài khoản SSH: `ubuntu`
- VPS public IPv4: `52.77.252.250`
- Website: `https://vieclamai.top`
- LiveKit: `wss://livekit.vieclamai.top`
- Email ACME/SSL: `nguyenhung22032006@gmail.com`
- PostgreSQL: Supabase Cloud
- Database production không chạy container PostgreSQL.

## 2. Kiến trúc

- Caddy phục vụ frontend, API và WebSocket qua HTTPS.
- LiveKit self-hosted xử lý audio/video WebRTC.
- Go backend gồm API và realtime gateway.
- Python AI service gọi Gemini.
- Redis chạy nội bộ.
- Supabase Cloud cung cấp PostgreSQL.
- Uploads dùng Docker volume `uploads_data`.

## 3. DNS bắt buộc

Trong trang quản lý DNS của `vieclamai.top`, tạo:

| Type | Name | Value |
|---|---|---|
| A | `@` | `52.77.252.250` |
| A | `livekit` | `52.77.252.250` |

Kết quả cần có:

```text
vieclamai.top          A  52.77.252.250
livekit.vieclamai.top  A  52.77.252.250
```

Kiểm tra từ máy local:

```bash
nslookup vieclamai.top
nslookup livekit.vieclamai.top
```

Chỉ khởi động Caddy sau khi DNS đã phân giải đúng; nếu không, Caddy không thể cấp SSL.

## 4. Chuẩn bị VPS

Đăng nhập bằng SSH:

```bash
ssh root@52.77.252.250
```

Cài Docker và công cụ cần thiết:

```bash
apt update
apt install -y ca-certificates curl git
curl -fsSL https://get.docker.com | sh
docker compose version
```

Khuyến nghị dùng user `ubuntu` hiện có để triển khai; không chạy ứng dụng bằng root. Thêm user này vào nhóm Docker:

```bash
sudo usermod -aG docker ubuntu
# đăng xuất rồi SSH lại để quyền nhóm có hiệu lực
```

## 5. Lấy mã nguồn

Thay URL bên dưới bằng repository Git thật:

```bash
git clone <REPOSITORY_URL> CuocThiAI
cd CuocThiAI
```

Nếu repository private, dùng SSH deploy key hoặc GitHub token; không ghi token vào file và không commit token.

## 6. Cấu hình secrets

```bash
cp .env.production.example .env.production
chmod 600 .env.production
nano .env.production
```

Các giá trị production hiện tại cần đặt:

```env
APP_DOMAIN=vieclamai.top
LIVEKIT_DOMAIN=livekit.vieclamai.top
ACME_EMAIL=nguyenhung22032006@gmail.com
FRONTEND_URL=https://vieclamai.top
PUBLIC_BASE_URL=https://vieclamai.top
LIVEKIT_URL=wss://livekit.vieclamai.top
VITE_LIVEKIT_URL=wss://livekit.vieclamai.top
DB_SSL_MODE=require
AI_MOCK=false
```

Điền các secrets thật:

- `DATABASE_URL`, `DB_HOST`, `DB_USER`, `DB_PASSWORD` của Supabase.
- `GEMINI_API_KEY`.
- `JWT_SECRET`, `FILE_URL_SECRET`.
- `LIVEKIT_API_KEY`, `LIVEKIT_API_SECRET`.
- `AI_SETTINGS_ENCRYPTION_KEY`, `AI_INTERNAL_SERVICE_TOKEN`.
- `REDIS_PASSWORD`.

Sinh secret mới:

```bash
openssl rand -hex 32
```

Không commit `.env.production`, API key hoặc mật khẩu vào Git. Mật khẩu VPS tạm phải đổi sau lần đăng nhập đầu tiên:

```bash
passwd
```

## 7. Firewall và cổng

```bash
apt install -y ufw
ufw allow OpenSSH
ufw allow 80/tcp
ufw allow 443/tcp
ufw allow 7881/tcp
ufw allow 50000:50100/udp
ufw enable
ufw status
```

Ngoài UFW, cần mở các cổng tương tự trong Cloud Security Group của nhà cung cấp VPS nếu có.

## 8. Kiểm tra và khởi động

Từ thư mục root repository:

```bash
docker compose --env-file .env.production \
  -f docker-compose.prod.yml config

docker compose --env-file .env.production \
  -f docker-compose.prod.yml config --services

docker compose --env-file .env.production \
  -f docker-compose.prod.yml up -d --build --remove-orphans
```

`migrate` là job chạy một lần. Trạng thái `Exited (0)` là thành công và bình thường.

Trong Compose hiện tại, service `migrate` phải tham gia cả hai network để vừa truy cập Redis nội bộ vừa truy cập Supabase Cloud:

```yaml
migrate:
  networks: [internal, edge]
```

## 9. Kiểm tra sau triển khai

```bash
docker compose --env-file .env.production -f docker-compose.prod.yml ps
docker compose --env-file .env.production -f docker-compose.prod.yml logs --tail=100 migrate
docker compose --env-file .env.production -f docker-compose.prod.yml logs --tail=100 backend-api
docker compose --env-file .env.production -f docker-compose.prod.yml logs --tail=100 caddy
curl -fsS https://vieclamai.top/health
```

Các service thường phải ở trạng thái `Up` hoặc `healthy`:

- `redis`
- `migrate` — `Exited (0)`
- `backend-api`
- `backend-realtime`
- `ai-service`
- `livekit`
- `frontend`
- `caddy`

## 10. Cập nhật phiên bản

```bash
cd /path/to/CuocThiAI
git pull
docker compose --env-file .env.production \
  -f docker-compose.prod.yml up -d --build --remove-orphans
docker image prune -f
```

Migration mới tự chạy trước khi API và realtime khởi động.

## 11. Xem log và dừng hệ thống

```bash
docker compose --env-file .env.production -f docker-compose.prod.yml logs -f backend-api
docker compose --env-file .env.production -f docker-compose.prod.yml logs -f backend-realtime
docker compose --env-file .env.production -f docker-compose.prod.yml logs -f ai-service
```

Dừng stack:

```bash
docker compose --env-file .env.production -f docker-compose.prod.yml down
```

Không dùng `down -v` nếu không muốn xóa volume uploads/Redis/Caddy.

## 12. Xử lý lỗi thường gặp

### Migration lỗi DNS hoặc không kết nối Supabase

Kiểm tra `DATABASE_URL`, DNS và network. Service `migrate` phải có `[internal, edge]`. Sau khi sửa:

```bash
docker compose --env-file .env.production -f docker-compose.prod.yml up -d --force-recreate
```

### Backend restart với lỗi LiveKit URL

Không dùng `ws://livekit:7880` trong production. Dùng:

```env
LIVEKIT_URL=wss://livekit.vieclamai.top
```

Kiểm tra DNS, cổng TCP `7881` và UDP `50000-50100`.

### Caddy không cấp SSL

1. Kiểm tra hai bản ghi DNS.
2. Đảm bảo cổng 80/443 không bị dịch vụ khác chiếm.
3. Kiểm tra `ACME_EMAIL` là email thật.
4. Xem log:

```bash
docker compose --env-file .env.production -f docker-compose.prod.yml logs -f caddy
```

### Frontend không gọi được API

Kiểm tra `FRONTEND_URL`, `PUBLIC_BASE_URL`, Caddyfile và log `backend-api`. Không dùng localhost trong production.

## 13. Bảo mật

- Không public Redis, API nội bộ hoặc AI service.
- Không dùng `devkey`, `devsecret`, `minioadmin` trong production.
- Đổi mật khẩu VPS tạm.
- Rotate mọi API key/secret đã lộ.
- Không bật `AI_MOCK` production.
- Sao lưu Supabase bằng Dashboard/PITR.
- Backup volume `uploads_data` nếu file CV quan trọng.
