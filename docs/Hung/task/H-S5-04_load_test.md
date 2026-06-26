# H-S5-04 — Load Test nhiều Room

| Field | Value |
|---|---|
| **Task ID** | H-S5-04 |
| **Sprint** | 5 — Hardening |
| **Độ khó** | Khó |
| **Trạng thái** | ✅ Đã hoàn thành |
| **Owner** | Hùng |

---

## Mục tiêu

Kiểm tra hiệu năng hệ thống realtime với nhiều room đồng thời:
- Xác định bottleneck
- Đảm bảo không memory leak
- Đo latency WebSocket ở tải cao
- Validate connection limit

---

## Yêu cầu chi tiết

### Test Scenarios

| Scenario | Target |
|---|---|
| 10 rooms × 2 participants | Baseline — phải hoạt động mượt |
| 50 rooms × 2 participants | Medium load |
| 100 rooms × 2 participants | High load |
| 1 room × 10 chat messages/giây | Chat burst |
| 1 room × continuous transcript | Transcript throughput |

### Metrics cần đo

- WebSocket message latency (p50, p95, p99)
- Memory usage per room / per connection
- CPU usage at peak
- Connection establish time
- Broadcast time (1 event → N participants)
- DB write latency (transcript, audit)

### Tools

- Custom Go test client (mô phỏng participant)
- `k6` hoặc `artillery` cho WebSocket load test
- `pprof` cho Go profiling
- Grafana/Prometheus dashboard (nếu có)

### Acceptance Criteria

| Metric | Target |
|---|---|
| Message latency p95 | < 100ms |
| Memory per room | < 10MB |
| Connection time | < 500ms |
| Broadcast (10 participants) | < 50ms |
| No goroutine leak | ✅ |
| No memory leak (1 giờ) | ✅ |

---

## Definition of Done

- [x] Load test script hoạt động
- [x] 10 rooms đồng thời → pass
- [x] 50 rooms đồng thời → pass hoặc biết bottleneck
- [x] Không goroutine leak
- [x] Không memory leak sau 1 giờ
- [x] Report kết quả load test

---

## Dependency

- **Tất cả sprint trước** hoàn thành
- Tất cả features cần ổn định trước khi load test

---

## Checklist test

- [x] 10 rooms — latency < 100ms
- [x] 50 rooms — hệ thống vẫn responsive
- [x] Chat burst — không mất tin
- [x] Transcript continuous — không delay quá 2s
- [x] goroutine count ổn định sau 30 phút
- [x] Memory ổn định sau 30 phút

---

## Output

- File: `docs/Hung/load_test_report.md`
- Ghi lại kết quả, bottleneck phát hiện, hướng tối ưu
