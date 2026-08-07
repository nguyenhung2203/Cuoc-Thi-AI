# Nguồn dữ liệu tin tuyển dụng — Đắk Lắk

## Mục đích

Bổ sung tin việc làm gắn Đắk Lắk / Buôn Ma Thuột cho job board (tiêu chí giảng viên: nguồn dữ liệu rõ + ưu tiên địa phương).

## Nguồn tham khảo (sao chép thủ công, không bot cào hàng loạt)

| Nguồn | URL / ghi chú | Thời điểm tham chiếu |
|---|---|---|
| TopCV — Việc làm Đắk Lắk | https://www.topcv.vn/tim-viec-lam-moi-nhat-tai-dak-lak-l23 | 2026-08-07 |
| TopCV — Buôn Ma Thuột | https://www.topcv.vn/tim-viec-lam-tai-tp-buon-ma-thuot-l23d2302 | 2026-08-07 |
| CareerViet — Dak Lak | https://careerviet.vn/viec-lam/dak-lak-l50-vi.html | 2026-08-07 |
| Vieclam.net (Vinamilk BMT) | Tin bán hàng siêu thị Buôn Ma Thuột | 2026-08-07 |

## Cách xử lý trong dự án

1. **Copy thủ công** tiêu đề, công ty, mức lương, địa điểm, ngành từ trang công khai.
2. **Viết lại** mô tả / yêu cầu / quyền lợi ở dạng ngắn gọn phù hợp schema `jobs` (không copy nguyên văn JD dài nếu không cần).
3. **Logo:**
   - Brand thật: PNG/JPG/SVG trong `frontend/public/company-logos/` (favicon/site công khai hoặc Wikimedia).
   - Công ty demo địa phương: PNG minh họa trong cùng thư mục.
   - Cập nhật DB: `seed_daklak_logos.sql` + `seed_daklak_logos_batch2.sql`.
   - Không phụ thuộc CDN TopCV lúc runtime.
4. Nạp dữ liệu bằng file:
   - `backend/seeds/seed_daklak_jobs.sql` (10 công ty + 10 tin)
   - `backend/seeds/seed_daklak_jobs_batch2.sql` (+20 công ty + 40 tin)
5. **UTF-8 trên Windows:** không dùng `Get-Content … | docker … psql` (PowerShell sẽ làm mất dấu tiếng Việt thành `???`). Hãy mount thư mục seeds và chạy `-f`:

```bash
docker run --rm -v "%CD%/backend/seeds:/seeds:ro" -e PGCLIENTENCODING=UTF8 postgres:16-alpine \
  psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -f /seeds/seed_daklak_jobs_batch2.sql
```

Nếu batch2 đã bị lỗi ký tự, chạy wipe rồi import lại: `seed_daklak_jobs_batch2_wipe.sql` → `seed_daklak_jobs_batch2.sql`.

## Lưu ý pháp lý / học thuật

- Đây là **dataset demo phục vụ đồ án**, không phải đồng bộ realtime với TopCV.
- Khi demo trước giảng viên: nêu rõ nguồn tham khảo ở tài liệu này.
- Có thể thay bằng tin mới hơn bằng cách cập nhật SQL + bảng nguồn.

## Phạm vi module

Thuộc **dữ liệu tin tuyển dụng (Âu)** — phục vụ UI tìm việc của ứng viên (Khiến).
