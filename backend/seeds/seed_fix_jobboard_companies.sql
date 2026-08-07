-- ============================================================================
-- seed_fix_jobboard_companies.sql
-- ----------------------------------------------------------------------------
-- Job board đang bị "ngập" tin gắn TechViet (logo mặc định), đẩy tin Đắk Lắk
-- đa công ty xuống trang sau.
--
-- 1) Ẩn (soft-delete) các tin flood dưới TechViet:
--    - cccccccc-0000-4000-8000-* : bản trùng tiêu đề Đắk Lắk, sai company
--    - cccccccc-0000-0000-0000-0005..0036 : bulk seed_jobs cùng 1 công ty
-- 2) Giữ demo TechViet 0001–0004.
-- 3) Gắn logo cho TechViet + đồng bộ lại logo công ty Đắk Lắk.
-- ============================================================================

BEGIN;

-- Flood bản sao Đắk Lắk gắn nhầm TechViet
UPDATE jobs
SET deleted_at = NOW(), updated_at = NOW()
WHERE deleted_at IS NULL
  AND company_id = 'aaaaaaaa-0000-0000-0000-000000000001'
  AND id >= 'cccccccc-0000-4000-8000-000000000001'
  AND id <= 'cccccccc-0000-4000-8000-000000000200';

-- Bulk IT seed cùng TechViet (giữ 0001–0004 của seed_demo)
UPDATE jobs
SET deleted_at = NOW(), updated_at = NOW()
WHERE deleted_at IS NULL
  AND company_id = 'aaaaaaaa-0000-0000-0000-000000000001'
  AND id >= 'cccccccc-0000-0000-0000-000000000005'
  AND id <= 'cccccccc-0000-0000-0000-000000000099';

-- Đưa tin Đắk Lắk lên đầu danh sách (ORDER BY created_at DESC)
UPDATE jobs
SET created_at = NOW(), updated_at = NOW()
WHERE deleted_at IS NULL
  AND id >= 'dddddddd-0000-4000-8000-0000000000d1'
  AND id <= 'dddddddd-0000-4000-8000-000000000228';

-- Logo TechViet (tránh fallback /images/logo.png)
UPDATE companies
SET logo_url = '/company-logos/techviet.svg', updated_at = NOW()
WHERE id = 'aaaaaaaa-0000-0000-0000-000000000001';

-- Logo công ty Demo Recruiter (nếu còn tin open)
UPDATE companies
SET logo_url = '/company-logos/demo-recruiter.svg', updated_at = NOW()
WHERE id = '0cd4d1ac-2ef2-42e5-a100-e83ffb9091ca'
   OR (logo_url IS NULL AND name ILIKE '%Demo Recruiter%');

COMMIT;
