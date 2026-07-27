-- ============================================================================
-- seed_demo_teardown.sql — Gỡ sạch dữ liệu DEMO do seed_demo.sql tạo ra
-- ----------------------------------------------------------------------------
-- Chỉ xóa các bản ghi có UUID cố định của bộ demo (prefix đã đặt trong seed).
-- KHÔNG đụng tới dữ liệu thật của team.
-- Xóa theo thứ tự ngược phụ thuộc (con trước, cha sau) để không vỡ foreign key.
-- An toàn chạy lại nhiều lần.
-- ============================================================================

BEGIN;

-- 15) AI suggestions
DELETE FROM ai_suggestions       WHERE id IN (
  '99999999-0000-0000-0000-000000000001',
  '99999999-0000-0000-0000-000000000002',
  '99999999-0000-0000-0000-000000000003');

-- 14) Report
DELETE FROM interview_reports    WHERE id = '66666666-0000-0000-0000-000000000001';

-- 13) Scores
DELETE FROM interview_scores     WHERE interview_id = '11111111-0000-0000-0000-000000000001';

-- 12) Transcripts
DELETE FROM interview_transcripts WHERE interview_id = '11111111-0000-0000-0000-000000000001';

-- 11) Participants
DELETE FROM interview_participants WHERE interview_id = '11111111-0000-0000-0000-000000000001';

-- 10) Rooms
DELETE FROM interview_rooms      WHERE id IN (
  '22222222-0000-0000-0000-000000000001',
  '22222222-0000-0000-0000-000000000002',
  '22222222-0000-0000-0000-000000000003');

-- 9) Interviews
DELETE FROM interviews           WHERE id IN (
  '11111111-0000-0000-0000-000000000001',
  '11111111-0000-0000-0000-000000000002',
  '11111111-0000-0000-0000-000000000003');

-- 8) Job-candidates
DELETE FROM job_candidates       WHERE company_id = 'aaaaaaaa-0000-0000-0000-000000000001';

-- 7) Candidates
DELETE FROM candidates           WHERE company_id = 'aaaaaaaa-0000-0000-0000-000000000001';

-- 6) Files
DELETE FROM files                WHERE company_id = 'aaaaaaaa-0000-0000-0000-000000000001';

-- 5) Rubric criteria + rubric
DELETE FROM rubric_criteria      WHERE rubric_id = '77777777-0000-0000-0000-000000000001';
DELETE FROM rubrics              WHERE id = '77777777-0000-0000-0000-000000000001';

-- 4) Question bank
DELETE FROM question_bank        WHERE company_id = 'aaaaaaaa-0000-0000-0000-000000000001';

-- 3) Jobs
DELETE FROM jobs                 WHERE company_id = 'aaaaaaaa-0000-0000-0000-000000000001';

-- 2) Company members + company
DELETE FROM company_members      WHERE company_id = 'aaaaaaaa-0000-0000-0000-000000000001';
DELETE FROM companies            WHERE id = 'aaaaaaaa-0000-0000-0000-000000000001';

-- 1) Users
DELETE FROM users                WHERE id IN (
  'bbbbbbbb-0000-0000-0000-000000000001',
  'bbbbbbbb-0000-0000-0000-000000000002',
  'bbbbbbbb-0000-0000-0000-000000000003');

COMMIT;

-- ============================================================================
-- Sau khi chạy: toàn bộ dữ liệu demo bị gỡ, dữ liệu thật của team giữ nguyên.
-- ============================================================================
