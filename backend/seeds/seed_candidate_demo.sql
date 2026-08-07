-- Demo candidate account for portal testing.
-- Password: demo1234
-- Safe to re-run (ON CONFLICT DO NOTHING).

BEGIN;

INSERT INTO users (id, email, password_hash, full_name, role, status, email_verified_at) VALUES
  (
    'bbbbbbbb-0000-0000-0000-000000000010',
    'candidate@demo.local',
    crypt('demo1234', gen_salt('bf', 10)),
    'Nguyễn Văn Ứng Viên',
    'candidate',
    'active',
    NOW()
  )
ON CONFLICT (id) DO NOTHING;

COMMIT;
