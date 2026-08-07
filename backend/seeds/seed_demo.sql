-- ============================================================================
-- seed_demo.sql — Dữ liệu mẫu cho DEMO tuyển dụng (Recruiter + Admin)
-- ----------------------------------------------------------------------------
-- An toàn để chạy lại nhiều lần: mọi INSERT dùng UUID cố định + ON CONFLICT
-- DO NOTHING. Toàn bộ bọc trong 1 transaction.
--
-- Namespace: mọi user dùng email @demo.local, company slug 'demo-techviet'.
-- Gỡ sạch bằng: backend/seeds/seed_demo_teardown.sql
--
-- Mật khẩu đăng nhập cho MỌI tài khoản demo: demo1234
--   (hash bằng pgcrypto crypt()+bcrypt, tương thích bcrypt phía backend)
--
-- Tài khoản demo:
--   admin@demo.local      / demo1234  (admin)
--   recruiter@demo.local  / demo1234  (recruiter, owner công ty)
--   hr2@demo.local        / demo1234  (recruiter, member)
-- ============================================================================

BEGIN;

-- ---------------------------------------------------------------------------
-- 1) USERS (1 admin + 2 recruiter)
-- ---------------------------------------------------------------------------
INSERT INTO users (id, email, password_hash, full_name, role, status, email_verified_at) VALUES
  ('bbbbbbbb-0000-0000-0000-000000000001', 'admin@demo.local',     crypt('demo1234', gen_salt('bf', 10)), 'Demo Admin',        'admin',     'active', NOW()),
  ('bbbbbbbb-0000-0000-0000-000000000002', 'recruiter@demo.local', crypt('demo1234', gen_salt('bf', 10)), 'Trần Thu Hà',       'recruiter', 'active', NOW()),
  ('bbbbbbbb-0000-0000-0000-000000000003', 'hr2@demo.local',       crypt('demo1234', gen_salt('bf', 10)), 'Nguyễn Minh Quang', 'recruiter', 'active', NOW())
ON CONFLICT (id) DO NOTHING;

-- ---------------------------------------------------------------------------
-- 2) COMPANY + members
-- ---------------------------------------------------------------------------
INSERT INTO companies (id, name, slug, logo_url, website, industry, size, created_by, settings) VALUES
  ('aaaaaaaa-0000-0000-0000-000000000001', 'TechViet Solutions', 'demo-techviet',
   '/company-logos/techviet.svg', 'https://techviet.example.vn',
   'Công nghệ thông tin', '50-200', 'bbbbbbbb-0000-0000-0000-000000000002',
   '{"timezone":"Asia/Ho_Chi_Minh","locale":"vi"}'::jsonb)
ON CONFLICT (id) DO UPDATE SET
  logo_url = COALESCE(EXCLUDED.logo_url, companies.logo_url),
  updated_at = NOW();

INSERT INTO company_members (id, company_id, user_id, role, status, joined_at) VALUES
  ('a1a1a1a1-0000-0000-0000-000000000001', 'aaaaaaaa-0000-0000-0000-000000000001', 'bbbbbbbb-0000-0000-0000-000000000002', 'owner',  'active', NOW()),
  ('a1a1a1a1-0000-0000-0000-000000000002', 'aaaaaaaa-0000-0000-0000-000000000001', 'bbbbbbbb-0000-0000-0000-000000000003', 'member', 'active', NOW())
ON CONFLICT (id) DO NOTHING;

-- ---------------------------------------------------------------------------
-- 3) JOBS (4 vị trí, mix trạng thái)
-- ---------------------------------------------------------------------------
INSERT INTO jobs (id, company_id, title, department, level, location, employment_type,
                  salary_min, salary_max, currency, description, requirements, benefits, status,
                  ai_summary, created_by) VALUES
  ('cccccccc-0000-0000-0000-000000000001', 'aaaaaaaa-0000-0000-0000-000000000001',
   'Frontend Developer (VueJS)', 'Engineering', 'middle', 'Hà Nội', 'full_time',
   18000000, 30000000, 'VND',
   'Phát triển giao diện người dùng cho nền tảng SaaS bằng Vue 3 + Vite. Phối hợp với team backend và thiết kế để hoàn thiện trải nghiệm.',
   'Thành thạo Vue 3, JavaScript ES6+, HTML/CSS. Kinh nghiệm 2+ năm. Biết Pinia, Vue Router. Ưu tiên có kinh nghiệm TypeScript.',
   'Lương tháng 13, bảo hiểm đầy đủ, làm việc hybrid, review lương 2 lần/năm.',
   'open',
   'Vị trí Frontend middle tập trung Vue 3, phù hợp ứng viên 2-4 năm kinh nghiệm SPA.',
   'bbbbbbbb-0000-0000-0000-000000000002'),

  ('cccccccc-0000-0000-0000-000000000002', 'aaaaaaaa-0000-0000-0000-000000000001',
   'Backend Engineer (Go)', 'Engineering', 'senior', 'Hà Nội', 'full_time',
   30000000, 50000000, 'VND',
   'Thiết kế và xây dựng microservices bằng Go, PostgreSQL, Redis. Chịu trách nhiệm hiệu năng và độ tin cậy hệ thống realtime.',
   'Thành thạo Go, PostgreSQL, thiết kế API. Hiểu concurrency, message queue. Kinh nghiệm 4+ năm backend.',
   'ESOP, laptop cấu hình cao, ngân sách học tập 10 triệu/năm.',
   'open',
   'Backend senior Go cho hệ thống realtime, đòi hỏi kinh nghiệm concurrency và tối ưu hiệu năng.',
   'bbbbbbbb-0000-0000-0000-000000000002'),

  ('cccccccc-0000-0000-0000-000000000003', 'aaaaaaaa-0000-0000-0000-000000000001',
   'Data Scientist', 'Data', 'middle', 'TP. Hồ Chí Minh', 'full_time',
   25000000, 40000000, 'VND',
   'Xây dựng mô hình ML phục vụ gợi ý và phân tích dữ liệu tuyển dụng. Làm việc với Python, scikit-learn, và LLM.',
   'Thành thạo Python, pandas, scikit-learn. Kiến thức thống kê vững. Ưu tiên có kinh nghiệm NLP/LLM.',
   'Chế độ remote linh hoạt, thưởng theo dự án.',
   'open',
   'Data Scientist middle, trọng tâm ML + NLP cho sản phẩm tuyển dụng AI.',
   'bbbbbbbb-0000-0000-0000-000000000003'),

  ('cccccccc-0000-0000-0000-000000000004', 'aaaaaaaa-0000-0000-0000-000000000001',
   'DevOps Engineer', 'Infrastructure', 'senior', 'Hà Nội', 'full_time',
   28000000, 45000000, 'VND',
   'Quản lý hạ tầng cloud, CI/CD, giám sát và bảo mật hệ thống. Kubernetes, Docker, Terraform.',
   'Kinh nghiệm Kubernetes, Docker, CI/CD. Biết Terraform, AWS/GCP. 4+ năm kinh nghiệm.',
   'On-call phụ cấp, chứng chỉ cloud được tài trợ.',
   'paused',
   'DevOps senior cho hạ tầng cloud-native, cần kinh nghiệm K8s và IaC.',
   'bbbbbbbb-0000-0000-0000-000000000002')
ON CONFLICT (id) DO NOTHING;

-- ---------------------------------------------------------------------------
-- 4) QUESTION BANK (demo trang Ngân hàng câu hỏi)
-- ---------------------------------------------------------------------------
INSERT INTO question_bank (id, company_id, job_id, created_by, question_text, question_type, skill_tags, level, expected_signals, is_ai_generated) VALUES
  ('b0000000-0000-0000-0000-000000000001', 'aaaaaaaa-0000-0000-0000-000000000001', 'cccccccc-0000-0000-0000-000000000001', 'bbbbbbbb-0000-0000-0000-000000000002',
   'Giải thích sự khác nhau giữa ref và reactive trong Vue 3. Khi nào bạn dùng cái nào?', 'technical',
   '["vue","reactivity"]'::jsonb, 'middle', '["hiểu Composition API","phân biệt primitive vs object"]'::jsonb, false),
  ('b0000000-0000-0000-0000-000000000002', 'aaaaaaaa-0000-0000-0000-000000000001', 'cccccccc-0000-0000-0000-000000000001', 'bbbbbbbb-0000-0000-0000-000000000002',
   'Làm thế nào để tối ưu hiệu năng render một danh sách lớn trong Vue?', 'technical',
   '["vue","performance"]'::jsonb, 'middle', '["v-for key","virtual scroll","computed caching"]'::jsonb, true),
  ('b0000000-0000-0000-0000-000000000003', 'aaaaaaaa-0000-0000-0000-000000000001', 'cccccccc-0000-0000-0000-000000000002', 'bbbbbbbb-0000-0000-0000-000000000002',
   'Mô tả cách bạn xử lý goroutine leak trong một dịch vụ Go chạy dài hạn.', 'technical',
   '["go","concurrency"]'::jsonb, 'senior', '["context cancellation","channel close","pprof"]'::jsonb, false),
  ('b0000000-0000-0000-0000-000000000004', 'aaaaaaaa-0000-0000-0000-000000000001', 'cccccccc-0000-0000-0000-000000000002', 'bbbbbbbb-0000-0000-0000-000000000002',
   'Thiết kế schema và index cho bảng lưu 50 triệu bản ghi giao dịch, truy vấn theo user_id và thời gian.', 'system_design',
   '["postgresql","indexing"]'::jsonb, 'senior', '["composite index","partitioning","query plan"]'::jsonb, true),
  ('b0000000-0000-0000-0000-000000000005', 'aaaaaaaa-0000-0000-0000-000000000001', NULL, 'bbbbbbbb-0000-0000-0000-000000000002',
   'Kể về một lần bạn bất đồng quan điểm kỹ thuật với đồng nghiệp. Bạn đã xử lý thế nào?', 'behavioral',
   '["teamwork","communication"]'::jsonb, 'middle', '["lắng nghe","dữ liệu thuyết phục","tôn trọng"]'::jsonb, false)
ON CONFLICT (id) DO NOTHING;

-- ---------------------------------------------------------------------------
-- 5) RUBRIC + tiêu chí (dùng cho phỏng vấn Frontend)
-- ---------------------------------------------------------------------------
INSERT INTO rubrics (id, company_id, job_id, name, description, total_weight, created_by) VALUES
  ('77777777-0000-0000-0000-000000000001', 'aaaaaaaa-0000-0000-0000-000000000001', 'cccccccc-0000-0000-0000-000000000001',
   'Rubric Frontend Developer', 'Bộ tiêu chí đánh giá phỏng vấn kỹ thuật Frontend.', 100,
   'bbbbbbbb-0000-0000-0000-000000000002')
ON CONFLICT (id) DO NOTHING;

INSERT INTO rubric_criteria (id, rubric_id, name, description, weight, min_score, max_score, order_index) VALUES
  ('88888888-0000-0000-0000-000000000001', '77777777-0000-0000-0000-000000000001', 'Kiến thức kỹ thuật', 'Nắm vững Vue, JS, CSS.',            35, 1, 5, 1),
  ('88888888-0000-0000-0000-000000000002', '77777777-0000-0000-0000-000000000001', 'Giải quyết vấn đề',   'Tư duy phân tích, tiếp cận bài toán.', 30, 1, 5, 2),
  ('88888888-0000-0000-0000-000000000003', '77777777-0000-0000-0000-000000000001', 'Giao tiếp',           'Trình bày rõ ràng, mạch lạc.',       20, 1, 5, 3),
  ('88888888-0000-0000-0000-000000000004', '77777777-0000-0000-0000-000000000001', 'Phù hợp văn hoá',     'Tinh thần đồng đội, chủ động.',      15, 1, 5, 4)
ON CONFLICT (id) DO NOTHING;

-- ---------------------------------------------------------------------------
-- 6) FILES (metadata CV — trỏ tới storage key demo)
-- ---------------------------------------------------------------------------
INSERT INTO files (id, company_id, owner_user_id, original_name, storage_key, mime_type, size_bytes, file_type) VALUES
  ('eeeeeeee-0000-0000-0000-000000000001', 'aaaaaaaa-0000-0000-0000-000000000001', 'bbbbbbbb-0000-0000-0000-000000000002', 'CV_LeVanAn.pdf',     'demo/cv/le-van-an.pdf',     'application/pdf', 245000, 'cv'),
  ('eeeeeeee-0000-0000-0000-000000000002', 'aaaaaaaa-0000-0000-0000-000000000001', 'bbbbbbbb-0000-0000-0000-000000000002', 'CV_PhamThiBinh.pdf', 'demo/cv/pham-thi-binh.pdf', 'application/pdf', 198000, 'cv'),
  ('eeeeeeee-0000-0000-0000-000000000003', 'aaaaaaaa-0000-0000-0000-000000000001', 'bbbbbbbb-0000-0000-0000-000000000002', 'CV_TranQuocCuong.pdf','demo/cv/tran-quoc-cuong.pdf','application/pdf', 312000, 'cv'),
  ('eeeeeeee-0000-0000-0000-000000000004', 'aaaaaaaa-0000-0000-0000-000000000001', 'bbbbbbbb-0000-0000-0000-000000000003', 'CV_HoangThiDung.pdf','demo/cv/hoang-thi-dung.pdf','application/pdf', 267000, 'cv'),
  ('eeeeeeee-0000-0000-0000-000000000005', 'aaaaaaaa-0000-0000-0000-000000000001', 'bbbbbbbb-0000-0000-0000-000000000003', 'CV_VuMinhEm.pdf',    'demo/cv/vu-minh-em.pdf',    'application/pdf', 221000, 'cv'),
  ('eeeeeeee-0000-0000-0000-000000000006', 'aaaaaaaa-0000-0000-0000-000000000001', 'bbbbbbbb-0000-0000-0000-000000000002', 'CV_DangVanPhuc.pdf', 'demo/cv/dang-van-phuc.pdf', 'application/pdf', 289000, 'cv')
ON CONFLICT (id) DO NOTHING;

-- ---------------------------------------------------------------------------
-- 7) CANDIDATES (6 ứng viên, mix trạng thái pipeline)
-- ---------------------------------------------------------------------------
INSERT INTO candidates (id, company_id, full_name, email, phone, cv_file_id, parsed_cv_json, ai_cv_summary, source, status, tags, created_by) VALUES
  ('dddddddd-0000-0000-0000-000000000001', 'aaaaaaaa-0000-0000-0000-000000000001', 'Lê Văn An',      'an.le@example.com',      '0901000001', 'eeeeeeee-0000-0000-0000-000000000001',
   '{"skills":["Vue","JavaScript","CSS"],"years":3,"education":"ĐH Bách Khoa HN"}'::jsonb,
   'Frontend 3 năm, mạnh Vue 3 và CSS, từng làm dự án SaaS.', 'linkedin', 'interviewing', '["vue","frontend"]'::jsonb, 'bbbbbbbb-0000-0000-0000-000000000002'),
  ('dddddddd-0000-0000-0000-000000000002', 'aaaaaaaa-0000-0000-0000-000000000001', 'Phạm Thị Bình',  'binh.pham@example.com',  '0901000002', 'eeeeeeee-0000-0000-0000-000000000002',
   '{"skills":["Vue","TypeScript","Pinia"],"years":2,"education":"ĐH FPT"}'::jsonb,
   'Frontend 2 năm, thành thạo TypeScript và Pinia.', 'referral', 'screening', '["vue","typescript"]'::jsonb, 'bbbbbbbb-0000-0000-0000-000000000002'),
  ('dddddddd-0000-0000-0000-000000000003', 'aaaaaaaa-0000-0000-0000-000000000001', 'Trần Quốc Cường','cuong.tran@example.com', '0901000003', 'eeeeeeee-0000-0000-0000-000000000003',
   '{"skills":["Go","PostgreSQL","Redis"],"years":5,"education":"ĐH Công Nghệ"}'::jsonb,
   'Backend 5 năm Go, kinh nghiệm hệ thống realtime và tối ưu DB.', 'linkedin', 'interviewing', '["go","backend","senior"]'::jsonb, 'bbbbbbbb-0000-0000-0000-000000000002'),
  ('dddddddd-0000-0000-0000-000000000004', 'aaaaaaaa-0000-0000-0000-000000000001', 'Hoàng Thị Dung', 'dung.hoang@example.com', '0901000004', 'eeeeeeee-0000-0000-0000-000000000004',
   '{"skills":["Python","scikit-learn","NLP"],"years":3,"education":"ĐH KHTN"}'::jsonb,
   'Data Scientist 3 năm, mạnh NLP và mô hình gợi ý.', 'manual', 'new', '["data","python","ml"]'::jsonb, 'bbbbbbbb-0000-0000-0000-000000000003'),
  ('dddddddd-0000-0000-0000-000000000005', 'aaaaaaaa-0000-0000-0000-000000000001', 'Vũ Minh Em',     'em.vu@example.com',      '0901000005', 'eeeeeeee-0000-0000-0000-000000000005',
   '{"skills":["Kubernetes","Docker","Terraform"],"years":4,"education":"ĐH Bách Khoa HCM"}'::jsonb,
   'DevOps 4 năm, thành thạo K8s và IaC.', 'import', 'talent_pool', '["devops","k8s"]'::jsonb, 'bbbbbbbb-0000-0000-0000-000000000003'),
  ('dddddddd-0000-0000-0000-000000000006', 'aaaaaaaa-0000-0000-0000-000000000001', 'Đặng Văn Phúc',  'phuc.dang@example.com',  '0901000006', 'eeeeeeee-0000-0000-0000-000000000006',
   '{"skills":["Vue","Nuxt","JavaScript"],"years":1,"education":"ĐH FPT"}'::jsonb,
   'Frontend junior 1 năm, nền tảng tốt, ham học.', 'referral', 'rejected', '["vue","junior"]'::jsonb, 'bbbbbbbb-0000-0000-0000-000000000002')
ON CONFLICT (id) DO NOTHING;

-- ---------------------------------------------------------------------------
-- 8) JOB_CANDIDATES (ứng tuyển + fit_score AI)
-- ---------------------------------------------------------------------------
INSERT INTO job_candidates (id, company_id, job_id, candidate_id, pipeline_status, fit_score, ai_match_json, applied_at, created_by) VALUES
  ('ffffffff-0000-0000-0000-000000000001', 'aaaaaaaa-0000-0000-0000-000000000001', 'cccccccc-0000-0000-0000-000000000001', 'dddddddd-0000-0000-0000-000000000001', 'interviewing', 88.50, '{"matched":["Vue","CSS"],"missing":["TypeScript"]}'::jsonb, NOW() - INTERVAL '5 days',  'bbbbbbbb-0000-0000-0000-000000000002'),
  ('ffffffff-0000-0000-0000-000000000002', 'aaaaaaaa-0000-0000-0000-000000000001', 'cccccccc-0000-0000-0000-000000000001', 'dddddddd-0000-0000-0000-000000000002', 'screening',    76.00, '{"matched":["Vue","TypeScript"],"missing":["kinh nghiệm SaaS"]}'::jsonb, NOW() - INTERVAL '3 days', 'bbbbbbbb-0000-0000-0000-000000000002'),
  ('ffffffff-0000-0000-0000-000000000003', 'aaaaaaaa-0000-0000-0000-000000000001', 'cccccccc-0000-0000-0000-000000000001', 'dddddddd-0000-0000-0000-000000000006', 'rejected',     52.00, '{"matched":["Vue"],"missing":["kinh nghiệm","CSS nâng cao"]}'::jsonb, NOW() - INTERVAL '8 days', 'bbbbbbbb-0000-0000-0000-000000000002'),
  ('ffffffff-0000-0000-0000-000000000004', 'aaaaaaaa-0000-0000-0000-000000000001', 'cccccccc-0000-0000-0000-000000000002', 'dddddddd-0000-0000-0000-000000000003', 'interviewing', 91.00, '{"matched":["Go","PostgreSQL","Redis"],"missing":[]}'::jsonb, NOW() - INTERVAL '4 days', 'bbbbbbbb-0000-0000-0000-000000000002'),
  ('ffffffff-0000-0000-0000-000000000005', 'aaaaaaaa-0000-0000-0000-000000000001', 'cccccccc-0000-0000-0000-000000000003', 'dddddddd-0000-0000-0000-000000000004', 'new',          72.50, '{"matched":["Python","NLP"],"missing":["kinh nghiệm production"]}'::jsonb, NOW() - INTERVAL '1 days', 'bbbbbbbb-0000-0000-0000-000000000003'),
  ('ffffffff-0000-0000-0000-000000000006', 'aaaaaaaa-0000-0000-0000-000000000001', 'cccccccc-0000-0000-0000-000000000004', 'dddddddd-0000-0000-0000-000000000005', 'talent_pool',  80.00, '{"matched":["Kubernetes","Docker","Terraform"],"missing":[]}'::jsonb, NOW() - INTERVAL '10 days', 'bbbbbbbb-0000-0000-0000-000000000003')
ON CONFLICT (id) DO NOTHING;

-- ---------------------------------------------------------------------------
-- 9) INTERVIEWS (1 completed đầy đủ dữ liệu, 1 scheduled, 1 active)
-- ---------------------------------------------------------------------------
INSERT INTO interviews (id, company_id, job_id, candidate_id, recruiter_id, rubric_id, mode, title,
                        scheduled_at, started_at, ended_at, status, room_id, consent_recording, consent_ai, created_by) VALUES
  -- Completed: Lê Văn An (Frontend) — có report + scores + transcript
  ('11111111-0000-0000-0000-000000000001', 'aaaaaaaa-0000-0000-0000-000000000001', 'cccccccc-0000-0000-0000-000000000001', 'dddddddd-0000-0000-0000-000000000001',
   'bbbbbbbb-0000-0000-0000-000000000002', '77777777-0000-0000-0000-000000000001', 'real', 'Phỏng vấn kỹ thuật — Lê Văn An (Frontend)',
   NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days' + INTERVAL '5 min', NOW() - INTERVAL '2 days' + INTERVAL '50 min', 'completed',
   '22222222-0000-0000-0000-000000000001', true, true, 'bbbbbbbb-0000-0000-0000-000000000002'),
  -- Scheduled: Phạm Thị Bình (Frontend) — sắp diễn ra
  ('11111111-0000-0000-0000-000000000002', 'aaaaaaaa-0000-0000-0000-000000000001', 'cccccccc-0000-0000-0000-000000000001', 'dddddddd-0000-0000-0000-000000000002',
   'bbbbbbbb-0000-0000-0000-000000000002', '77777777-0000-0000-0000-000000000001', 'real', 'Phỏng vấn kỹ thuật — Phạm Thị Bình (Frontend)',
   NOW() + INTERVAL '1 day', NULL, NULL, 'scheduled',
   '22222222-0000-0000-0000-000000000002', false, true, 'bbbbbbbb-0000-0000-0000-000000000002'),
  -- Active: Trần Quốc Cường (Backend) — đang diễn ra
  ('11111111-0000-0000-0000-000000000003', 'aaaaaaaa-0000-0000-0000-000000000001', 'cccccccc-0000-0000-0000-000000000002', 'dddddddd-0000-0000-0000-000000000003',
   'bbbbbbbb-0000-0000-0000-000000000002', NULL, 'real', 'Phỏng vấn kỹ thuật — Trần Quốc Cường (Backend Go)',
   NOW() - INTERVAL '10 min', NOW() - INTERVAL '5 min', NULL, 'active',
   '22222222-0000-0000-0000-000000000003', true, true, 'bbbbbbbb-0000-0000-0000-000000000002')
ON CONFLICT (id) DO NOTHING;

-- ---------------------------------------------------------------------------
-- 10) INTERVIEW ROOMS
-- ---------------------------------------------------------------------------
INSERT INTO interview_rooms (id, interview_id, room_code, status, provider, opened_at, closed_at) VALUES
  ('22222222-0000-0000-0000-000000000001', '11111111-0000-0000-0000-000000000001', 'DEMO-ROOM-001', 'closed',  'livekit', NOW() - INTERVAL '2 days' + INTERVAL '5 min', NOW() - INTERVAL '2 days' + INTERVAL '50 min'),
  ('22222222-0000-0000-0000-000000000002', '11111111-0000-0000-0000-000000000002', 'DEMO-ROOM-002', 'waiting', 'livekit', NULL, NULL),
  ('22222222-0000-0000-0000-000000000003', '11111111-0000-0000-0000-000000000003', 'DEMO-ROOM-003', 'active',  'livekit', NOW() - INTERVAL '5 min', NULL)
ON CONFLICT (id) DO NOTHING;

-- ---------------------------------------------------------------------------
-- 11) PARTICIPANTS (cho interview completed)
-- ---------------------------------------------------------------------------
INSERT INTO interview_participants (id, interview_id, user_id, participant_type, display_name, joined_at, left_at, connection_state) VALUES
  ('33333333-0000-0000-0000-000000000001', '11111111-0000-0000-0000-000000000001', 'bbbbbbbb-0000-0000-0000-000000000002', 'recruiter', 'Trần Thu Hà', NOW() - INTERVAL '2 days' + INTERVAL '5 min', NOW() - INTERVAL '2 days' + INTERVAL '50 min', 'offline'),
  ('33333333-0000-0000-0000-000000000002', '11111111-0000-0000-0000-000000000001', NULL,                                   'candidate', 'Lê Văn An',   NOW() - INTERVAL '2 days' + INTERVAL '6 min', NOW() - INTERVAL '2 days' + INTERVAL '49 min', 'offline'),
  ('33333333-0000-0000-0000-000000000003', '11111111-0000-0000-0000-000000000001', NULL,                                   'ai',        'AI Trợ lý',   NOW() - INTERVAL '2 days' + INTERVAL '5 min', NOW() - INTERVAL '2 days' + INTERVAL '50 min', 'offline')
ON CONFLICT (id) DO NOTHING;

-- ---------------------------------------------------------------------------
-- 12) TRANSCRIPTS (cho interview completed)
-- ---------------------------------------------------------------------------
INSERT INTO interview_transcripts (id, interview_id, participant_id, speaker_type, speaker_name, content, language, source, is_final) VALUES
  ('44444444-0000-0000-0000-000000000001', '11111111-0000-0000-0000-000000000001', '33333333-0000-0000-0000-000000000001', 'recruiter', 'Trần Thu Hà', 'Chào An, cảm ơn em đã tham gia. Em có thể giới thiệu ngắn về bản thân không?', 'vi', 'audio', true),
  ('44444444-0000-0000-0000-000000000002', '11111111-0000-0000-0000-000000000001', '33333333-0000-0000-0000-000000000002', 'candidate', 'Lê Văn An',   'Dạ em là An, có 3 năm kinh nghiệm Frontend, chủ yếu làm với Vue 3 và các dự án SaaS.', 'vi', 'audio', true),
  ('44444444-0000-0000-0000-000000000003', '11111111-0000-0000-0000-000000000001', '33333333-0000-0000-0000-000000000001', 'recruiter', 'Trần Thu Hà', 'Em giải thích giúp chị sự khác nhau giữa ref và reactive trong Vue 3?', 'vi', 'audio', true),
  ('44444444-0000-0000-0000-000000000004', '11111111-0000-0000-0000-000000000001', '33333333-0000-0000-0000-000000000002', 'candidate', 'Lê Văn An',   'ref dùng cho giá trị đơn, truy cập qua .value; reactive dùng cho object. ref bọc primitive còn reactive theo dõi sâu các thuộc tính của object.', 'vi', 'audio', true),
  ('44444444-0000-0000-0000-000000000005', '11111111-0000-0000-0000-000000000001', '33333333-0000-0000-0000-000000000003', 'ai',        'AI Trợ lý',   'Gợi ý: có thể hỏi thêm ứng viên về cách tối ưu render danh sách lớn để đánh giá chiều sâu.', 'vi', 'ai', true),
  ('44444444-0000-0000-0000-000000000006', '11111111-0000-0000-0000-000000000001', '33333333-0000-0000-0000-000000000002', 'candidate', 'Lê Văn An',   'Với danh sách lớn em dùng key ổn định trong v-for, kết hợp virtual scrolling và computed để cache kết quả.', 'vi', 'audio', true)
ON CONFLICT (id) DO NOTHING;

-- ---------------------------------------------------------------------------
-- 13) SCORES (cho interview completed — khớp 4 tiêu chí rubric)
-- ---------------------------------------------------------------------------
INSERT INTO interview_scores (id, interview_id, rubric_criterion_id, criterion_name, score, max_score, weight, weighted_score, evidence, ai_comment, confidence, status, scored_by) VALUES
  ('55555555-0000-0000-0000-000000000001', '11111111-0000-0000-0000-000000000001', '88888888-0000-0000-0000-000000000001', 'Kiến thức kỹ thuật', 4.5, 5, 35, 31.50, 'Giải thích ref/reactive chính xác, hiểu tối ưu render.', 'Nền tảng Vue vững, trả lời sâu.',        0.9200, 'final', 'ai'),
  ('55555555-0000-0000-0000-000000000002', '11111111-0000-0000-0000-000000000001', '88888888-0000-0000-0000-000000000002', 'Giải quyết vấn đề',   4.0, 5, 30, 24.00, 'Tiếp cận bài toán render list có hệ thống.',            'Tư duy tốt, có thể sâu hơn về đo lường.', 0.8800, 'final', 'ai'),
  ('55555555-0000-0000-0000-000000000003', '11111111-0000-0000-0000-000000000001', '88888888-0000-0000-0000-000000000003', 'Giao tiếp',           4.5, 5, 20, 18.00, 'Trình bày mạch lạc, dễ hiểu.',                         'Giao tiếp rõ ràng, tự tin.',             0.9000, 'final', 'ai'),
  ('55555555-0000-0000-0000-000000000004', '11111111-0000-0000-0000-000000000001', '88888888-0000-0000-0000-000000000004', 'Phù hợp văn hoá',     4.0, 5, 15, 12.00, 'Thái độ cầu thị, chủ động.',                           'Phù hợp môi trường team.',               0.8500, 'final', 'ai')
ON CONFLICT (id) DO NOTHING;

-- ---------------------------------------------------------------------------
-- 14) REPORT (cho interview completed)
-- ---------------------------------------------------------------------------
INSERT INTO interview_reports (id, interview_id, summary, final_score, recommendation, strengths, weaknesses, risks,
                               ai_reasoning_summary, report_json, generated_by) VALUES
  ('66666666-0000-0000-0000-000000000001', '11111111-0000-0000-0000-000000000001',
   'Ứng viên Lê Văn An thể hiện nền tảng Frontend vững, đặc biệt là Vue 3 và tối ưu hiệu năng. Giao tiếp tốt, thái độ cầu thị. Đề xuất tuyển.',
   85.50, 'hire',
   '["Nắm chắc Vue 3 (ref/reactive, tối ưu render)","Giao tiếp rõ ràng","Thái độ cầu thị"]'::jsonb,
   '["Chưa nhiều kinh nghiệm TypeScript","Cần đo lường hiệu năng bài bản hơn"]'::jsonb,
   '["Mức lương kỳ vọng có thể cao hơn ngân sách"]'::jsonb,
   'Điểm tổng hợp có trọng số đạt 85.5/100, vượt ngưỡng tuyển. Các tiêu chí kỹ thuật và giao tiếp đều cao.',
   '{"final_score":85.5,"recommendation":"hire","criteria_count":4}'::jsonb,
   'ai')
ON CONFLICT (id) DO NOTHING;

-- ---------------------------------------------------------------------------
-- 15) AI SUGGESTIONS (demo AI Insights trên Dashboard)
-- ---------------------------------------------------------------------------
INSERT INTO ai_suggestions (id, interview_id, suggestion_type, title, content, reason, priority, confidence) VALUES
  ('99999999-0000-0000-0000-000000000001', '11111111-0000-0000-0000-000000000001', 'next_action', 'Ứng viên phù hợp cao', 'Lê Văn An đạt 85.5 điểm — nên chuyển sang vòng offer.', 'Điểm vượt ngưỡng tuyển và fit_score 88.5%.', 'high', 0.9100),
  ('99999999-0000-0000-0000-000000000002', '11111111-0000-0000-0000-000000000003', 'reminder',    'Phỏng vấn đang diễn ra', 'Buổi phỏng vấn Backend với Trần Quốc Cường đang diễn ra.', 'Interview status = active.', 'medium', 0.8000),
  ('99999999-0000-0000-0000-000000000003', '11111111-0000-0000-0000-000000000002', 'reminder',    'Nhắc lịch phỏng vấn',   'Phỏng vấn với Phạm Thị Bình diễn ra ngày mai.', 'Interview scheduled trong 24h tới.', 'medium', 0.8300)
ON CONFLICT (id) DO NOTHING;

COMMIT;

-- ============================================================================
-- Tổng kết dữ liệu demo:
--   3 users (1 admin, 2 recruiter) | 1 company | 4 jobs | 5 câu hỏi
--   1 rubric + 4 tiêu chí | 6 files | 6 candidates | 6 job_candidates
--   3 interviews (1 completed đầy đủ, 1 scheduled, 1 active)
--   3 rooms | 3 participants | 6 transcripts | 4 scores | 1 report | 3 AI suggestions
-- Đăng nhập demo: recruiter@demo.local / demo1234  (hoặc admin@demo.local)
-- ============================================================================
