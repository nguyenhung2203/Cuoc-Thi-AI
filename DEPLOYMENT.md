# Hướng dẫn triển khai production bằng Docker

Tài liệu này hướng dẫn chạy toàn bộ hệ thống trên server Linux có domain và HTTPS.

## 1. Kiến trúc

- `https://APP_DOMAIN`: giao diện, API và WebSocket.
- `wss://LIVEKIT_DOMAIN`: LiveKit signaling.
- Caddy tự cấp và gia hạn SSL Let's Encrypt.
- PostgreSQL dùng Supabase Cloud qua `DATABASE_URL`; production không chạy container PostgreSQL.
- Redis, API, realtime, AI và LiveKit chạy trong Docker network nội bộ.
- LiveKit self-hosted, dùng UDP media ports `50000-50100`.

## 2. Điều kiện server

Khuyến nghị Ubuntu 22.04/24.04, tối thiểu 4 GB RAM và 2 vCPU.
Cài Docker Engine + Docker Compose plugin:

```bash
sudo apt update
sudo apt install -y ca-certificates curl git
curl -fsSL https://get.docker.com | sudo sh
sudo usermod -aG docker "$USER"
# đăng xuất/đăng nhập lại sau lệnh trên
 docker compose version
```

Mở firewall:

```bash
sudo ufw allow OpenSSH
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 7881/tcp
sudo ufw allow 50000:50100/udp
sudo ufw enable
```

## 3. DNS

Tạo hai bản ghi `A` trỏ về public IP server:

- `APP_DOMAIN` — website
- `LIVEKIT_DOMAIN` — LiveKit

Ví dụ:

```text
app.example.com       A    SERVER_PUBLIC_IP
livekit.example.com   A    SERVER_PUBLIC_IP
```

DNS phải phân giải đúng trước khi chạy Caddy, nếu không Let's Encrypt không cấp được certificate.

## 4. Cấu hình secrets

```bash
git clone <REPOSITORY_URL>
cd CuocThiAI
cp .env.production.example .env.production
chmod 600 .env.production
```

Sửa `.env.production`:

- `APP_DOMAIN`, `LIVEKIT_DOMAIN`, `ACME_EMAIL`
- `DATABASE_URL` của Supabase Cloud và `DB_SSL_MODE=require`
- `JWT_SECRET`, `FILE_URL_SECRET`
- `LIVEKIT_API_KEY`, `LIVEKIT_API_SECRET`
- `GEMINI_API_KEY`
- SMTP/OAuth nếu dùng

Tạo secret ngẫu nhiên:

```bash
openssl rand -hex 32
```

Không commit `.env.production`, API key hoặc private key vào Git.

## 5. Cấu hình LiveKit

LiveKit được tạo cấu hình runtime từ `deploy/livekit/livekit.prod.yaml.template`; API key và secret lấy từ `.env.production`, vì vậy không cần và không được ghi secret trực tiếp vào YAML.

LiveKit cần public UDP range `50000-50100`; nếu server nằm sau NAT, cấu hình public IP/NAT theo nhà cung cấp server. Với mạng client bị hạn chế, nên bổ sung TURN/TLS.

## 6. Khởi chạy một lệnh

Từ thư mục root repository:

```bash
docker compose --env-file .env.production \
  -f docker-compose.prod.yml \
  up -d --build --remove-orphans
```

Kiểm tra trạng thái:

```bash
docker compose --env-file .env.production -f docker-compose.prod.yml ps
docker compose --env-file .env.production -f docker-compose.prod.yml logs -f migrate
```

Kiểm tra website và HTTPS:

```bash
curl -I https://app.example.com
curl -fsS https://app.example.com/health
```

Certificate có thể mất vài phút trong lần đầu. Xem log Caddy:

```bash
docker compose --env-file .env.production -f docker-compose.prod.yml logs -f caddy
```

## 7. Cập nhật phiên bản

```bash
git pull
docker compose --env-file .env.production -f docker-compose.prod.yml up -d --build --remove-orphans
docker image prune -f
```

Migration mới sẽ tự chạy một lần trước khi API/realtime khởi động.

## 8. Dừng và xem log

```bash
docker compose --env-file .env.production -f docker-compose.prod.yml logs -f backend-api

docker compose --env-file .env.production -f docker-compose.prod.yml logs -f backend-realtime

docker compose --env-file .env.production -f docker-compose.prod.yml down
```

Không dùng `down -v` trừ khi muốn xóa vĩnh viễn database, Redis và uploads.

## 9. Backup PostgreSQL (Supabase Cloud)

Dùng Supabase Dashboard để cấu hình backup/PITR, hoặc chạy `pg_dump` từ máy đã cài PostgreSQL client:

```bash
set -a
. ./.env.production
set +a
pg_dump "$DATABASE_URL" > "backups/db-$(date +%Y%m%d-%H%M%S).sql"
```

Không có volume PostgreSQL cục bộ trong stack. Volume `uploads_data` vẫn cần backup nếu file upload quan trọng.

## 10. Kiểm tra cấu hình trước khi chạy

```bash
docker compose --env-file .env.production -f docker-compose.prod.yml config

docker compose --env-file .env.production -f docker-compose.prod.yml build
```

Nếu lỗi certificate:

1. kiểm tra DNS `dig +short APP_DOMAIN` và `dig +short LIVEKIT_DOMAIN`;
2. kiểm tra cổng 80/443 không bị Nginx/Apache khác chiếm;
3. xem `logs caddy`;
4. chắc chắn `ACME_EMAIL` hợp lệ.

Nếu video không hoạt động:

1. kiểm tra UDP `50000-50100` trên firewall/cloud security group;
2. kiểm tra TCP 7881;
3. kiểm tra browser đang dùng `wss://LIVEKIT_DOMAIN`;
4. kiểm tra credentials backend và LiveKit giống nhau;
5. cân nhắc TURN cho mạng doanh nghiệp/VPN.

## 11. Lưu ý bảo mật

- Không public Redis, API hoặc AI; PostgreSQL nằm trên Supabase Cloud.
- Không dùng `devkey`, `devsecret`, `minioadmin` trong production.
- Đổi toàn bộ placeholder `REPLACE_WITH_*` trước khi chạy.
- Không bật `AI_MOCK` trong production.
- Theo dõi dung lượng PostgreSQL, uploads và Caddy certificates.
