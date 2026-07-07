# AI Interview Platform

## Tổng quan

Nền tảng phỏng vấn trực tuyến có AI real-time hỗ trợ nhà tuyển dụng và ứng viên.

## Tech Stack

| Layer | Technology |
|---|---|
| Frontend | Vue 3 + TailwindCSS + Pinia |
| Backend | Golang (Gin/Fiber) |
| AI Service | Python (FastAPI) + Gemini + Whisper |
| Database | PostgreSQL + Redis |
| Media | LiveKit SFU |
| Storage | S3-compatible |
| Deploy | Docker + Docker Compose |

## Quick Start

### Yêu cầu

- Go 1.21+
- Node.js 20+
- Python 3.11+
- Docker & Docker Compose
- PostgreSQL 15+
- Redis 7+

### Chạy toàn bộ bằng Docker

```bash
docker-compose up -d
```

### Chạy lại toàn bộ bằng Docker

```bash
docker-compose up -d --build
```

### Chạy từng service

```bash
# Frontend
cd frontend && npm install && npm run dev

# Backend API
cd backend && go run cmd/api/main.go

# Backend Realtime
cd backend && go run cmd/realtime/main.go

# AI Service
cd ai-service && pip install -r requirements.txt && uvicorn app.main:app --reload
```

## Cấu trúc dự án

Xem chi tiết tại [PROJECT_STRUCTURE.md](./PROJECT_STRUCTURE.md)

## Tài liệu

Xem tại [docs/README.md](./docs/README.md)

## Team

| Thành viên | Phạm vi |
|---|---|
| Hùng | Realtime, Room, WebSocket, Transcript |
| Khôi | Backend Core, DB, AI Engine, Scoring, Report |
| Lai | UI/UX, Portal, Business flow |
