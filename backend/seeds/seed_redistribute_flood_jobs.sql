-- ============================================================================
-- seed_redistribute_flood_jobs.sql
-- ----------------------------------------------------------------------------
-- Khôi phục 100 tin cccccccc-0000-4000-8000-* (đã soft-delete) và chia thành
-- 10 công ty × 10 tin, mỗi công ty tên + logo riêng.
-- An toàn chạy lại: ON CONFLICT DO UPDATE / gán lại company_id theo thứ tự id.
-- ============================================================================

BEGIN;

-- ---------------------------------------------------------------------------
-- 10 công ty mới (041 → 04a)
-- ---------------------------------------------------------------------------
-- Tên/logo theo TopCV/CareerViet (đa ngành, Đắk Lắk + tỉnh cạnh).
-- Chi tiết cập nhật: seed_update_group10_topcv.sql
INSERT INTO companies (id, name, slug, logo_url, website, industry, size, created_by, settings) VALUES
  ('aaaaaaaa-0000-0000-0000-000000000041', 'Ngân hàng TMCP Đông Nam Á (SeABank) - CN Đắk Lắk', 'demo-seabank-daklak',
   '/company-logos/seabank.png', 'https://www.seabank.com.vn', 'Ngân hàng / Tài chính', '1000+',
   'bbbbbbbb-0000-0000-0000-000000000002', '{"timezone":"Asia/Ho_Chi_Minh","locale":"vi","region":"DakLak","source":"TopCV"}'::jsonb),

  ('aaaaaaaa-0000-0000-0000-000000000042', 'WinMart / WinCommerce - Khu vực Lâm Đồng', 'demo-winmart-lamdong',
   '/company-logos/winmart.png', 'https://winmart.vn', 'Bán lẻ / Siêu thị', '1000+',
   'bbbbbbbb-0000-0000-0000-000000000002', '{"timezone":"Asia/Ho_Chi_Minh","locale":"vi","region":"LamDong","source":"CareerViet"}'::jsonb),

  ('aaaaaaaa-0000-0000-0000-000000000043', 'Công ty TNHH Shopee - Kho miền Trung & Tây Nguyên', 'demo-shopee-taynguyen',
   '/company-logos/shopee.png', 'https://shopee.vn', 'Logistics / E-commerce', '1000+',
   'bbbbbbbb-0000-0000-0000-000000000002', '{"timezone":"Asia/Ho_Chi_Minh","locale":"vi","region":"TayNguyen","source":"CareerViet"}'::jsonb),

  ('aaaaaaaa-0000-0000-0000-000000000044', 'Công ty Cổ phần Thực phẩm Dinh dưỡng NutiFood', 'demo-nutifood-daklak',
   '/company-logos/nutifood.svg', 'https://www.nutifood.com.vn', 'FMCG / Thực phẩm', '1000+',
   'bbbbbbbb-0000-0000-0000-000000000002', '{"timezone":"Asia/Ho_Chi_Minh","locale":"vi","region":"DakLak","source":"CareerViet"}'::jsonb),

  ('aaaaaaaa-0000-0000-0000-000000000045', 'Công ty CP Đầu tư & PT NLMT Bách Khoa (SolarBK)', 'demo-solarbk-taynguyen',
   '/company-logos/solarbk.png', 'https://solarbk.vn', 'Năng lượng / Môi trường', '200-500',
   'bbbbbbbb-0000-0000-0000-000000000002', '{"timezone":"Asia/Ho_Chi_Minh","locale":"vi","region":"TayNguyen","source":"CareerViet"}'::jsonb),

  ('aaaaaaaa-0000-0000-0000-000000000046', 'Công ty CP Giáo dục Đại Dương (Ocean Edu) - Gia Lai / BMT', 'demo-oceanedu-gialai',
   '/company-logos/oceanedu.svg', 'https://oceanedu.vn', 'Giáo dục / Đào tạo', '200-500',
   'bbbbbbbb-0000-0000-0000-000000000002', '{"timezone":"Asia/Ho_Chi_Minh","locale":"vi","region":"GiaLai","source":"TopCV/CareerViet"}'::jsonb),

  ('aaaaaaaa-0000-0000-0000-000000000047', 'Tổng Công ty CP Bảo hiểm BIDV (BIC) - Gia Lai / Khánh Hòa', 'demo-bic-gialai-kh',
   '/company-logos/bic.png', 'https://www.bic.vn', 'Bảo hiểm', '1000+',
   'bbbbbbbb-0000-0000-0000-000000000002', '{"timezone":"Asia/Ho_Chi_Minh","locale":"vi","region":"GiaLai-KhanhHoa","source":"CareerViet"}'::jsonb),

  ('aaaaaaaa-0000-0000-0000-000000000048', 'Công ty TNHH Kiểm soát Chất lượng Nông sản Xanh Việt Nam', 'demo-nongsan-xanh-vn',
   '/company-logos/nongsan-xanh.svg', NULL, 'Nông nghiệp / Kiểm định', '50-200',
   'bbbbbbbb-0000-0000-0000-000000000002', '{"timezone":"Asia/Ho_Chi_Minh","locale":"vi","region":"DakLak","source":"TopCV"}'::jsonb),

  ('aaaaaaaa-0000-0000-0000-000000000049', 'Abbott Việt Nam - Khu vực Phú Yên / Đắk Lắk / Khánh Hòa', 'demo-abbott-mientrung',
   '/company-logos/abbott.svg', 'https://www.abbott.com.vn', 'Dược / Y tế', '1000+',
   'bbbbbbbb-0000-0000-0000-000000000002', '{"timezone":"Asia/Ho_Chi_Minh","locale":"vi","region":"PhuYen-DakLak-KhanhHoa","source":"Vieclam24h/TopCV"}'::jsonb),

  ('aaaaaaaa-0000-0000-0000-00000000004a', 'Công ty TNHH Sơn TOA Việt Nam - Khu vực Tây Nguyên', 'demo-toa-taynguyen',
   '/company-logos/toa.svg', 'https://www.toagroup.com.vn', 'Vật liệu xây dựng', '1000+',
   'bbbbbbbb-0000-0000-0000-000000000002', '{"timezone":"Asia/Ho_Chi_Minh","locale":"vi","region":"TayNguyen","source":"CareerViet"}'::jsonb)
ON CONFLICT (id) DO UPDATE SET
  name = EXCLUDED.name,
  slug = EXCLUDED.slug,
  logo_url = EXCLUDED.logo_url,
  website = EXCLUDED.website,
  industry = EXCLUDED.industry,
  size = EXCLUDED.size,
  settings = EXCLUDED.settings,
  updated_at = NOW();

INSERT INTO company_members (id, company_id, user_id, role, status, joined_at) VALUES
  ('a1a1a1a1-0000-0000-0000-000000000041', 'aaaaaaaa-0000-0000-0000-000000000041', 'bbbbbbbb-0000-0000-0000-000000000002', 'owner', 'active', NOW()),
  ('a1a1a1a1-0000-0000-0000-000000000042', 'aaaaaaaa-0000-0000-0000-000000000042', 'bbbbbbbb-0000-0000-0000-000000000002', 'owner', 'active', NOW()),
  ('a1a1a1a1-0000-0000-0000-000000000043', 'aaaaaaaa-0000-0000-0000-000000000043', 'bbbbbbbb-0000-0000-0000-000000000002', 'owner', 'active', NOW()),
  ('a1a1a1a1-0000-0000-0000-000000000044', 'aaaaaaaa-0000-0000-0000-000000000044', 'bbbbbbbb-0000-0000-0000-000000000002', 'owner', 'active', NOW()),
  ('a1a1a1a1-0000-0000-0000-000000000045', 'aaaaaaaa-0000-0000-0000-000000000045', 'bbbbbbbb-0000-0000-0000-000000000002', 'owner', 'active', NOW()),
  ('a1a1a1a1-0000-0000-0000-000000000046', 'aaaaaaaa-0000-0000-0000-000000000046', 'bbbbbbbb-0000-0000-0000-000000000002', 'owner', 'active', NOW()),
  ('a1a1a1a1-0000-0000-0000-000000000047', 'aaaaaaaa-0000-0000-0000-000000000047', 'bbbbbbbb-0000-0000-0000-000000000002', 'owner', 'active', NOW()),
  ('a1a1a1a1-0000-0000-0000-000000000048', 'aaaaaaaa-0000-0000-0000-000000000048', 'bbbbbbbb-0000-0000-0000-000000000002', 'owner', 'active', NOW()),
  ('a1a1a1a1-0000-0000-0000-000000000049', 'aaaaaaaa-0000-0000-0000-000000000049', 'bbbbbbbb-0000-0000-0000-000000000002', 'owner', 'active', NOW()),
  ('a1a1a1a1-0000-0000-0000-00000000004a', 'aaaaaaaa-0000-0000-0000-00000000004a', 'bbbbbbbb-0000-0000-0000-000000000002', 'owner', 'active', NOW())
ON CONFLICT (id) DO NOTHING;

-- ---------------------------------------------------------------------------
-- Khôi phục + chia 100 tin: 10 tin / công ty theo thứ tự id
-- ---------------------------------------------------------------------------
WITH ranked AS (
  SELECT
    id,
    ROW_NUMBER() OVER (ORDER BY id) AS rn
  FROM jobs
  WHERE id >= 'cccccccc-0000-4000-8000-000000000001'
    AND id <= 'cccccccc-0000-4000-8000-000000000200'
),
mapped AS (
  SELECT
    id,
    rn,
    (ARRAY[
      'aaaaaaaa-0000-0000-0000-000000000041'::uuid,
      'aaaaaaaa-0000-0000-0000-000000000042'::uuid,
      'aaaaaaaa-0000-0000-0000-000000000043'::uuid,
      'aaaaaaaa-0000-0000-0000-000000000044'::uuid,
      'aaaaaaaa-0000-0000-0000-000000000045'::uuid,
      'aaaaaaaa-0000-0000-0000-000000000046'::uuid,
      'aaaaaaaa-0000-0000-0000-000000000047'::uuid,
      'aaaaaaaa-0000-0000-0000-000000000048'::uuid,
      'aaaaaaaa-0000-0000-0000-000000000049'::uuid,
      'aaaaaaaa-0000-0000-0000-00000000004a'::uuid
    ])[((rn - 1) / 10) + 1] AS new_company_id
  FROM ranked
  WHERE rn <= 100
)
UPDATE jobs j
SET
  company_id = m.new_company_id,
  deleted_at = NULL,
  status = 'open',
  created_at = NOW() - ((100 - m.rn) * INTERVAL '1 minute'),
  updated_at = NOW()
FROM mapped m
WHERE j.id = m.id;

COMMIT;
