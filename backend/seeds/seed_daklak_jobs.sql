-- ============================================================================
-- seed_daklak_jobs.sql — Tin tuyển dụng gắn Đắk Lắk / Buôn Ma Thuột
-- ----------------------------------------------------------------------------
-- Nguồn tham khảo (copy thủ công, xem docs/SOURCE_JOB_DATA_DAKLAK.md):
--   - https://www.topcv.vn/tim-viec-lam-moi-nhat-tai-dak-lak-l23
--   - https://careerviet.vn/viec-lam/dak-lak-l50-vi.html
--
-- Yêu cầu: đã có user recruiter bbbbbbbb-0000-0000-0000-000000000002
--          (từ seed_demo.sql). An toàn chạy lại: ON CONFLICT DO NOTHING.
-- Logo: /company-logos/*.{png,jpg} phục vụ bởi frontend static (lấy từ TopCV CDN / site CT).
-- ============================================================================

BEGIN;

-- ---------------------------------------------------------------------------
-- Companies (Đắk Lắk / chi nhánh Tây Đắk Lắk) — đa ngành
-- ---------------------------------------------------------------------------
INSERT INTO companies (id, name, slug, logo_url, website, industry, size, created_by, settings) VALUES
  ('aaaaaaaa-0000-0000-0000-000000000011', 'Công ty TNHH Nông nghiệp Dế Mèn - VPĐD Đắk Lắk', 'demo-de-men-daklak',
   '/company-logos/de-men.png', NULL, 'Nông nghiệp', '50-200',
   'bbbbbbbb-0000-0000-0000-000000000002', '{"timezone":"Asia/Ho_Chi_Minh","locale":"vi","region":"DakLak"}'::jsonb),

  ('aaaaaaaa-0000-0000-0000-000000000012', 'Công ty Cổ phần MISA', 'demo-misa-daklak',
   '/company-logos/misa.png', 'https://www.misa.vn', 'Phần mềm / CNTT', '1000+',
   'bbbbbbbb-0000-0000-0000-000000000002', '{"timezone":"Asia/Ho_Chi_Minh","locale":"vi","region":"DakLak"}'::jsonb),

  ('aaaaaaaa-0000-0000-0000-000000000013', 'Công ty cổ phần Công nghệ Sapo', 'demo-sapo-daklak',
   '/company-logos/sapo.png', 'https://www.sapo.vn', 'Phần mềm / CNTT', '500-1000',
   'bbbbbbbb-0000-0000-0000-000000000002', '{"timezone":"Asia/Ho_Chi_Minh","locale":"vi","region":"DakLak"}'::jsonb),

  ('aaaaaaaa-0000-0000-0000-000000000014', 'Công ty Cổ phần Ong Mật Ban Mê Thuột', 'demo-ong-mat-bmt',
   '/company-logos/ong-mat.png', NULL, 'Thực phẩm / Nông sản', '50-200',
   'bbbbbbbb-0000-0000-0000-000000000002', '{"timezone":"Asia/Ho_Chi_Minh","locale":"vi","region":"DakLak"}'::jsonb),

  ('aaaaaaaa-0000-0000-0000-000000000015', 'Ngân hàng TMCP Lộc Phát Việt Nam (LPBank) - CN Đắk Lắk', 'demo-lpbank-daklak',
   '/company-logos/lpbank.png', 'https://lpbank.com.vn', 'Ngân hàng / Tài chính', '1000+',
   'bbbbbbbb-0000-0000-0000-000000000002', '{"timezone":"Asia/Ho_Chi_Minh","locale":"vi","region":"DakLak"}'::jsonb),

  ('aaaaaaaa-0000-0000-0000-000000000016', 'Vinamilk - Khu vực Buôn Ma Thuột', 'demo-vinamilk-bmt',
   '/company-logos/vinamilk.png', 'https://www.vinamilk.com.vn', 'FMCG / Thực phẩm', '1000+',
   'bbbbbbbb-0000-0000-0000-000000000002', '{"timezone":"Asia/Ho_Chi_Minh","locale":"vi","region":"DakLak"}'::jsonb),

  ('aaaaaaaa-0000-0000-0000-000000000017', 'Vietravel Đắk Lắk', 'demo-vietravel-daklak',
   '/company-logos/vietravel.png', 'https://www.vietravel.com', 'Du lịch / Dịch vụ', '200-500',
   'bbbbbbbb-0000-0000-0000-000000000002', '{"timezone":"Asia/Ho_Chi_Minh","locale":"vi","region":"DakLak"}'::jsonb),

  ('aaaaaaaa-0000-0000-0000-000000000018', 'Công ty TNHH Manulife Việt Nam', 'demo-manulife-daklak',
   '/company-logos/manulife.jpg', 'https://www.manulife.com.vn', 'Bảo hiểm', '1000+',
   'bbbbbbbb-0000-0000-0000-000000000002', '{"timezone":"Asia/Ho_Chi_Minh","locale":"vi","region":"DakLak"}'::jsonb),

  ('aaaaaaaa-0000-0000-0000-000000000019', 'Công ty Cổ phần Tập đoàn Hoa Sen', 'demo-hoa-sen-daklak',
   '/company-logos/hoa-sen.jpg', 'https://www.hoasengroup.vn', 'Vật liệu xây dựng', '1000+',
   'bbbbbbbb-0000-0000-0000-000000000002', '{"timezone":"Asia/Ho_Chi_Minh","locale":"vi","region":"DakLak"}'::jsonb),

  ('aaaaaaaa-0000-0000-0000-00000000001a', 'Công ty cổ phần Hai Bốn Bảy (247Express)', 'demo-247express-daklak',
   '/company-logos/express247.png', NULL, 'Logistics / Express', '500-1000',
   'bbbbbbbb-0000-0000-0000-000000000002', '{"timezone":"Asia/Ho_Chi_Minh","locale":"vi","region":"DakLak"}'::jsonb)
ON CONFLICT (id) DO NOTHING;

INSERT INTO company_members (id, company_id, user_id, role, status, joined_at) VALUES
  ('a1a1a1a1-0000-0000-0000-000000000011', 'aaaaaaaa-0000-0000-0000-000000000011', 'bbbbbbbb-0000-0000-0000-000000000002', 'owner', 'active', NOW()),
  ('a1a1a1a1-0000-0000-0000-000000000012', 'aaaaaaaa-0000-0000-0000-000000000012', 'bbbbbbbb-0000-0000-0000-000000000002', 'owner', 'active', NOW()),
  ('a1a1a1a1-0000-0000-0000-000000000013', 'aaaaaaaa-0000-0000-0000-000000000013', 'bbbbbbbb-0000-0000-0000-000000000002', 'owner', 'active', NOW()),
  ('a1a1a1a1-0000-0000-0000-000000000014', 'aaaaaaaa-0000-0000-0000-000000000014', 'bbbbbbbb-0000-0000-0000-000000000002', 'owner', 'active', NOW()),
  ('a1a1a1a1-0000-0000-0000-000000000015', 'aaaaaaaa-0000-0000-0000-000000000015', 'bbbbbbbb-0000-0000-0000-000000000002', 'owner', 'active', NOW()),
  ('a1a1a1a1-0000-0000-0000-000000000016', 'aaaaaaaa-0000-0000-0000-000000000016', 'bbbbbbbb-0000-0000-0000-000000000002', 'owner', 'active', NOW()),
  ('a1a1a1a1-0000-0000-0000-000000000017', 'aaaaaaaa-0000-0000-0000-000000000017', 'bbbbbbbb-0000-0000-0000-000000000002', 'owner', 'active', NOW()),
  ('a1a1a1a1-0000-0000-0000-000000000018', 'aaaaaaaa-0000-0000-0000-000000000018', 'bbbbbbbb-0000-0000-0000-000000000002', 'owner', 'active', NOW()),
  ('a1a1a1a1-0000-0000-0000-000000000019', 'aaaaaaaa-0000-0000-0000-000000000019', 'bbbbbbbb-0000-0000-0000-000000000002', 'owner', 'active', NOW()),
  ('a1a1a1a1-0000-0000-0000-00000000001a', 'aaaaaaaa-0000-0000-0000-00000000001a', 'bbbbbbbb-0000-0000-0000-000000000002', 'owner', 'active', NOW())
ON CONFLICT (id) DO NOTHING;

-- ---------------------------------------------------------------------------
-- Jobs (location = Đắk Lắk / Buôn Ma Thuột) — đa ngành
-- ---------------------------------------------------------------------------
-- Job IDs dải dddddddd-…0d1–0da — tách khỏi seed_jobs (cccccccc-…001–0xx).
INSERT INTO jobs (id, company_id, title, department, level, location, employment_type,
                  salary_min, salary_max, currency, description, requirements, benefits, status,
                  ai_summary, created_by) VALUES
  ('dddddddd-0000-4000-8000-0000000000d1', 'aaaaaaaa-0000-0000-0000-000000000011',
   'Nhân viên Chăm sóc khách hàng - Đắk Lắk',
   'Customer Service', 'junior', 'Buôn Ma Thuột, Đắk Lắk', 'full_time',
   7000000, 10000000, 'VND',
   'Tiếp nhận và xử lý yêu cầu khách hàng khu vực Đắk Lắk (nông nghiệp/phân phối). Hỗ trợ đơn hàng, phản hồi khiếu nại, cập nhật CRM.',
   'Tốt nghiệp Trung cấp trở lên. Kỹ năng giao tiếp tốt. Ưu tiên dưới 1 năm kinh nghiệm CSKH. Có xe máy, sẵn sàng làm việc tại Đắk Lắk.',
   'Lương 7-10 triệu, phụ cấp xăng xe, BHXH đầy đủ.',
   'open',
   'CSKH junior tại Đắk Lắk, phù hợp ứng viên mới vào nghề, ưu tiên giao tiếp và thái độ phục vụ.',
   'bbbbbbbb-0000-0000-0000-000000000002'),

  ('dddddddd-0000-4000-8000-0000000000d2', 'aaaaaaaa-0000-0000-0000-000000000012',
   'Nhân viên Kinh doanh giải pháp phần mềm cho hộ kinh doanh',
   'Sales', 'junior', 'Đắk Lắk', 'full_time',
   8000000, 20000000, 'VND',
   'Tư vấn và bán giải pháp phần mềm MISA cho hộ kinh doanh, cửa hàng, SME tại Đắk Lắk. Chăm sóc khách hàng sau bán.',
   'Không bắt buộc kinh nghiệm. Ham học hỏi về phần mềm kế toán/bán hàng. Kỹ năng tư vấn, chịu khó đi thị trường.',
   'Lương cứng + hoa hồng, đào tạo sản phẩm, thu nhập tới 20 triệu.',
   'open',
   'Sales phần mềm tại Đắk Lắk, phù hợp fresher/junior năng động.',
   'bbbbbbbb-0000-0000-0000-000000000002'),

  ('dddddddd-0000-4000-8000-0000000000d3', 'aaaaaaaa-0000-0000-0000-000000000013',
   'Chuyên viên Tư vấn phần mềm / Website (có nhận Fresher)',
   'Sales', 'junior', 'Đắk Lắk', 'full_time',
   12000000, 35000000, 'VND',
   'Tư vấn giải pháp website/phần mềm Sapo cho chủ shop, doanh nghiệp nhỏ. Open fresher, có đào tạo.',
   'Dưới 1 năm kinh nghiệm hoặc fresher. Kỹ năng nói chuyện thuyết phục. Biết sử dụng máy tính văn phòng.',
   'Lương cơ bản tới 13 triệu + hoa hồng cao, phụ cấp.',
   'open',
   'Telesales/tư vấn phần mềm toàn quốc có địa bàn Đắk Lắk, open fresher.',
   'bbbbbbbb-0000-0000-0000-000000000002'),

  ('dddddddd-0000-4000-8000-0000000000d4', 'aaaaaaaa-0000-0000-0000-000000000014',
   'Kế toán tổng hợp - Buôn Ma Thuột',
   'Accounting', 'middle', 'Buôn Ma Thuột, Đắk Lắk', 'full_time',
   10000000, 16000000, 'VND',
   'Phụ trách kế toán tổng hợp: hạch toán, lập báo cáo thuế, theo dõi công nợ, phối hợp kho/bán hàng nông sản mật ong Ban Mê.',
   'Tốt nghiệp Cao đẳng/ĐH Kế toán. Thành thạo Excel, phần mềm kế toán. Ưu tiên 1-3 năm kinh nghiệm.',
   'BHXH, thưởng lễ Tết theo doanh thu mùa vụ.',
   'open',
   'Kế toán middle tại doanh nghiệp nông sản Buôn Ma Thuột.',
   'bbbbbbbb-0000-0000-0000-000000000002'),

  ('dddddddd-0000-4000-8000-0000000000d5', 'aaaaaaaa-0000-0000-0000-000000000015',
   'Chuyên viên Khách hàng cá nhân - PGD Buôn Ma Thuột / Ea Kar / Cư M''gar',
   'Banking', 'middle', 'Buôn Ma Thuột, Đắk Lắk', 'full_time',
   12000000, 25000000, 'VND',
   'Tư vấn sản phẩm huy động, tín dụng, thẻ cho khách hàng cá nhân tại Đắk Lắk. Phát triển mạng lưới khách hàng địa phương.',
   '1-3 năm kinh nghiệm banking hoặc sales tài chính. Thành tích tốt. Có bằng CĐ/ĐH. Ưu tiên người địa phương.',
   'Lương cứng + KPI, phúc lợi ngân hàng.',
   'open',
   'Relationship banker cá nhân tại chi nhánh Đắk Lắk.',
   'bbbbbbbb-0000-0000-0000-000000000002'),

  ('dddddddd-0000-4000-8000-0000000000d6', 'aaaaaaaa-0000-0000-0000-000000000016',
   'Nhân viên bán hàng siêu thị - Buôn Ma Thuột',
   'Retail', 'junior', 'Buôn Ma Thuột, Đắk Lắk', 'full_time',
   10000000, 13000000, 'VND',
   'Tư vấn và bán sản phẩm sữa Vinamilk tại siêu thị khu vực BMT. Chăm sóc gian hàng, báo cáo doanh số.',
   'Nữ, 18-35 tuổi. THPT trở lên. Dưới 1 năm kinh nghiệm. Nhanh nhẹn, giao tiếp tốt.',
   'Lương 10-13 triệu, đồng phục, hỗ trợ đi lại.',
   'open',
   'Bán hàng FMCG tại siêu thị Buôn Ma Thuột.',
   'bbbbbbbb-0000-0000-0000-000000000002'),

  ('dddddddd-0000-4000-8000-0000000000d7', 'aaaaaaaa-0000-0000-0000-000000000017',
   'Thực tập sinh Kế toán - Vietravel Đắk Lắk',
   'Accounting', 'intern', 'Buôn Ma Thuột, Đắk Lắk', 'intern',
   3000000, 5000000, 'VND',
   'Hỗ trợ bộ phận kế toán chi nhánh du lịch: nhập chứng từ, đối soát tour, hỗ trợ báo cáo.',
   'Sinh viên năm cuối ngành Kế toán/Tài chính. Biết Excel cơ bản. Chăm chỉ, đúng giờ.',
   'Trợ cấp thực tập, được hướng dẫn, cơ hội nhận việc chính thức.',
   'open',
   'Intern kế toán ngành du lịch tại Đắk Lắk.',
   'bbbbbbbb-0000-0000-0000-000000000002'),

  ('dddddddd-0000-4000-8000-0000000000d8', 'aaaaaaaa-0000-0000-0000-000000000018',
   'Tư vấn viên bảo hiểm nhân thọ - Đắk Lắk',
   'Insurance', 'junior', 'Đắk Lắk', 'full_time',
   10000000, 50000000, 'VND',
   'Tư vấn giải pháp bảo hiểm Manulife cho khách hàng cá nhân/gia đình tại Đắk Lắk. Xây dựng đội nhóm (nếu phù hợp).',
   'Ưu tiên 2 năm kinh nghiệm sales/tư vấn. Giao tiếp tốt. Có xe máy. Không yêu cầu bằng cấp bảo hiểm lúc ứng tuyển.',
   'Thu nhập 10-50 triệu theo doanh số, đào tạo chứng chỉ.',
   'open',
   'Insurance advisor tại Đắk Lắk, thu nhập theo hiệu suất.',
   'bbbbbbbb-0000-0000-0000-000000000002'),

  ('dddddddd-0000-4000-8000-0000000000d9', 'aaaaaaaa-0000-0000-0000-000000000019',
   'Nhân viên Kinh doanh vật liệu xây dựng - Đắk Lắk',
   'Sales', 'junior', 'Đắk Lắk', 'full_time',
   12000000, 18000000, 'VND',
   'Phát triển đại lý, tư vấn tôn thép và vật liệu Hoa Sen tại thị trường Tây Nguyên / Đắk Lắk.',
   'Không bắt buộc kinh nghiệm ngành. Chịu khó đi thị trường. Ưu tiên nam, khỏe mạnh.',
   'Lương 12-18 triệu + thưởng doanh số.',
   'open',
   'Sales vật liệu xây dựng khu vực Đắk Lắk.',
   'bbbbbbbb-0000-0000-0000-000000000002'),

  ('dddddddd-0000-4000-8000-0000000000da', 'aaaaaaaa-0000-0000-0000-00000000001a',
   'Nhân viên Kinh doanh B2B - Logistic / Express',
   'Sales', 'junior', 'Đắk Lắk', 'full_time',
   20000000, 30000000, 'VND',
   'Kinh doanh dịch vụ chuyển phát B2B cho SME tại Đắk Lắk và vùng phụ cận. Chăm sóc tài khoản doanh nghiệp.',
   'Dưới 1 năm kinh nghiệm sales B2B hoặc logistics. Kỹ năng đàm phán. Có thể đi công tác tỉnh.',
   'Thu nhập 20-30 triệu, thưởng theo KPI.',
   'open',
   'Sales logistics B2B có địa bàn Đắk Lắk.',
   'bbbbbbbb-0000-0000-0000-000000000002')
ON CONFLICT (id) DO NOTHING;

-- Dọn bản seed thử nghiệm cũ (trùng ID seed_jobs hoặc dải 02a tạm).
DELETE FROM jobs WHERE id IN (
  'cccccccc-0000-0000-0000-00000000001a',
  'cccccccc-0000-0000-0000-00000000002a'
);

COMMIT;
