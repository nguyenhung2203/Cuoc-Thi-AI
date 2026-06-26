# Báo Cáo Kết Quả Load Test Realtime Gateway (H-S5-04)

## 1. Môi trường & Kịch bản thử nghiệm
- **Tool**: Custom Go Test Client (`backend/cmd/loadtest/main.go`).
- **Hình thức**: Chạy Local Server (cổng 8081).
- **Setup**: Mỗi Room bao gồm 2 WebSocket connection (1 Recruiter, 1 Candidate).
- **Traffic**: Recruiter gửi `chat:send` liên tục (1 msg/sec) và đo thời gian nhận bản tin broadcast trả về `chat:message`.
- **Thời lượng**: 10 giây mỗi mức tải.

## 2. Kết quả đo lường (Metrics)

| Mức tải (Scenario) | Total Messages | p50 Latency | p95 Latency | p99 Latency | Goroutines | Total Memory |
|:---|---:|---:|---:|---:|---:|---:|
| **Baseline (10 Rooms)** | 90 | 0ms | 0ms | 0ms | 40 | 3.33 MB |
| **Medium (50 Rooms)** | 450 | 0ms | 0ms | 0.38ms | 193 | 14.31 MB |
| **High (100 Rooms)** | 900 | 0ms | 0ms | 0.51ms | 392 | 28.05 MB |

## 3. Phân tích Acceptance Criteria
- **Message latency p95 < 100ms**: ✅ **PASS** (Luôn duy trì ở mức < 1ms ngay cả với 100 rooms).
- **Memory per room < 10MB**: ✅ **PASS** (~ 0.28 MB / room).
- **No goroutine leak**: ✅ **PASS** (Số lượng Goroutine tỉ lệ tuyến tính ổn định ~ 4 goroutines cho mỗi phòng, gồm 2 readPump, 2 writePump, không bị spike).

## 4. Kết luận
- Hệ thống Realtime Gateway hiện tại xử lý kết nối WebSocket và message routing rất nhanh, không có dấu hiệu bị thắt cổ chai (bottleneck) ở mức tải 100 phòng đồng thời.
- **Hướng tối ưu (nếu scale lớn hơn)**:
  - Khi triển khai production, cần cấu hình load balancer (Nginx/HAProxy) hỗ trợ sticky session hoặc sử dụng Redis Pub/Sub để scale ra nhiều node Gateway.
  - Theo dõi thêm số liệu I/O database khi bỏ chế độ mock DB.

*Ngày thực hiện: 2026-06-26*
