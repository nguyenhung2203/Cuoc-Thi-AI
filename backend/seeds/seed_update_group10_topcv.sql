-- ============================================================================
-- seed_update_group10_topcv.sql
-- ----------------------------------------------------------------------------
-- Đổi 10 công ty nhóm flood (041–04a) sang tên/logo tham khảo TopCV/CareerViet
-- đa ngành, khu vực Đắk Lắk + tỉnh cạnh (Gia Lai, Lâm Đồng, Khánh Hòa, Phú Yên).
-- Nguồn: docs/SOURCE_JOB_DATA_DAKLAK.md (cập nhật 2026-08-07).
-- ============================================================================

BEGIN;

UPDATE companies SET
  name = 'Ngân hàng TMCP Đông Nam Á (SeABank) - CN Đắk Lắk',
  slug = 'demo-seabank-daklak',
  logo_url = '/company-logos/seabank.png',
  website = 'https://www.seabank.com.vn',
  industry = 'Ngân hàng / Tài chính',
  size = '1000+',
  settings = '{"timezone":"Asia/Ho_Chi_Minh","locale":"vi","region":"DakLak","source":"TopCV"}'::jsonb,
  updated_at = NOW()
WHERE id = 'aaaaaaaa-0000-0000-0000-000000000041';

UPDATE companies SET
  name = 'WinMart / WinCommerce - Khu vực Lâm Đồng',
  slug = 'demo-winmart-lamdong',
  logo_url = '/company-logos/winmart.png',
  website = 'https://winmart.vn',
  industry = 'Bán lẻ / Siêu thị',
  size = '1000+',
  settings = '{"timezone":"Asia/Ho_Chi_Minh","locale":"vi","region":"LamDong","source":"CareerViet"}'::jsonb,
  updated_at = NOW()
WHERE id = 'aaaaaaaa-0000-0000-0000-000000000042';

UPDATE companies SET
  name = 'Công ty TNHH Shopee - Kho miền Trung & Tây Nguyên',
  slug = 'demo-shopee-taynguyen',
  logo_url = '/company-logos/shopee.png',
  website = 'https://shopee.vn',
  industry = 'Logistics / E-commerce',
  size = '1000+',
  settings = '{"timezone":"Asia/Ho_Chi_Minh","locale":"vi","region":"TayNguyen","source":"CareerViet"}'::jsonb,
  updated_at = NOW()
WHERE id = 'aaaaaaaa-0000-0000-0000-000000000043';

UPDATE companies SET
  name = 'Công ty Cổ phần Thực phẩm Dinh dưỡng NutiFood',
  slug = 'demo-nutifood-daklak',
  logo_url = '/company-logos/nutifood.svg',
  website = 'https://www.nutifood.com.vn',
  industry = 'FMCG / Thực phẩm',
  size = '1000+',
  settings = '{"timezone":"Asia/Ho_Chi_Minh","locale":"vi","region":"DakLak","source":"CareerViet"}'::jsonb,
  updated_at = NOW()
WHERE id = 'aaaaaaaa-0000-0000-0000-000000000044';

UPDATE companies SET
  name = 'Công ty CP Đầu tư & PT NLMT Bách Khoa (SolarBK)',
  slug = 'demo-solarbk-taynguyen',
  logo_url = '/company-logos/solarbk.png',
  website = 'https://solarbk.vn',
  industry = 'Năng lượng / Môi trường',
  size = '200-500',
  settings = '{"timezone":"Asia/Ho_Chi_Minh","locale":"vi","region":"TayNguyen","source":"CareerViet"}'::jsonb,
  updated_at = NOW()
WHERE id = 'aaaaaaaa-0000-0000-0000-000000000045';

UPDATE companies SET
  name = 'Công ty CP Giáo dục Đại Dương (Ocean Edu) - Gia Lai / BMT',
  slug = 'demo-oceanedu-gialai',
  logo_url = '/company-logos/oceanedu.svg',
  website = 'https://oceanedu.vn',
  industry = 'Giáo dục / Đào tạo',
  size = '200-500',
  settings = '{"timezone":"Asia/Ho_Chi_Minh","locale":"vi","region":"GiaLai","source":"TopCV/CareerViet"}'::jsonb,
  updated_at = NOW()
WHERE id = 'aaaaaaaa-0000-0000-0000-000000000046';

UPDATE companies SET
  name = 'Tổng Công ty CP Bảo hiểm BIDV (BIC) - Gia Lai / Khánh Hòa',
  slug = 'demo-bic-gialai-kh',
  logo_url = '/company-logos/bic.png',
  website = 'https://www.bic.vn',
  industry = 'Bảo hiểm',
  size = '1000+',
  settings = '{"timezone":"Asia/Ho_Chi_Minh","locale":"vi","region":"GiaLai-KhanhHoa","source":"CareerViet"}'::jsonb,
  updated_at = NOW()
WHERE id = 'aaaaaaaa-0000-0000-0000-000000000047';

UPDATE companies SET
  name = 'Công ty TNHH Kiểm soát Chất lượng Nông sản Xanh Việt Nam',
  slug = 'demo-nongsan-xanh-vn',
  logo_url = '/company-logos/nongsan-xanh.svg',
  website = NULL,
  industry = 'Nông nghiệp / Kiểm định',
  size = '50-200',
  settings = '{"timezone":"Asia/Ho_Chi_Minh","locale":"vi","region":"DakLak","source":"TopCV"}'::jsonb,
  updated_at = NOW()
WHERE id = 'aaaaaaaa-0000-0000-0000-000000000048';

UPDATE companies SET
  name = 'Abbott Việt Nam - Khu vực Phú Yên / Đắk Lắk / Khánh Hòa',
  slug = 'demo-abbott-mientrung',
  logo_url = '/company-logos/abbott.svg',
  website = 'https://www.abbott.com.vn',
  industry = 'Dược / Y tế',
  size = '1000+',
  settings = '{"timezone":"Asia/Ho_Chi_Minh","locale":"vi","region":"PhuYen-DakLak-KhanhHoa","source":"Vieclam24h/TopCV"}'::jsonb,
  updated_at = NOW()
WHERE id = 'aaaaaaaa-0000-0000-0000-000000000049';

UPDATE companies SET
  name = 'Công ty TNHH Sơn TOA Việt Nam - Khu vực Tây Nguyên',
  slug = 'demo-toa-taynguyen',
  logo_url = '/company-logos/toa.svg',
  website = 'https://www.toagroup.com.vn',
  industry = 'Vật liệu xây dựng',
  size = '1000+',
  settings = '{"timezone":"Asia/Ho_Chi_Minh","locale":"vi","region":"TayNguyen","source":"CareerViet"}'::jsonb,
  updated_at = NOW()
WHERE id = 'aaaaaaaa-0000-0000-0000-00000000004a';

COMMIT;
