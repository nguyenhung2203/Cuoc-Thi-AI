DROP INDEX IF EXISTS interviews_report_status_idx;
DROP INDEX IF EXISTS mock_interviews_ai_status_idx;
DROP INDEX IF EXISTS candidates_cv_ai_status_idx;

ALTER TABLE interviews
    DROP COLUMN IF EXISTS report_updated_at,
    DROP COLUMN IF EXISTS report_attempts,
    DROP COLUMN IF EXISTS report_error;

ALTER TABLE mock_interviews
    DROP COLUMN IF EXISTS ai_attempts,
    DROP COLUMN IF EXISTS ai_error,
    DROP COLUMN IF EXISTS ai_scoring_status,
    DROP COLUMN IF EXISTS ai_question_status;

ALTER TABLE candidates
    DROP COLUMN IF EXISTS cv_ai_updated_at,
    DROP COLUMN IF EXISTS cv_ai_attempts,
    DROP COLUMN IF EXISTS cv_ai_error,
    DROP COLUMN IF EXISTS cv_ai_status;
