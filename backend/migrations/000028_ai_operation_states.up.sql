-- Persist explicit AI processing state so failed provider calls are retryable
-- and never indistinguishable from an unprocessed record.
ALTER TABLE candidates
    ADD COLUMN IF NOT EXISTS cv_ai_status VARCHAR(30) NOT NULL DEFAULT 'pending',
    ADD COLUMN IF NOT EXISTS cv_ai_error TEXT,
    ADD COLUMN IF NOT EXISTS cv_ai_attempts INT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS cv_ai_updated_at TIMESTAMPTZ;

ALTER TABLE mock_interviews
    ADD COLUMN IF NOT EXISTS ai_question_status VARCHAR(30) NOT NULL DEFAULT 'pending',
    ADD COLUMN IF NOT EXISTS ai_scoring_status VARCHAR(30) NOT NULL DEFAULT 'pending',
    ADD COLUMN IF NOT EXISTS ai_error TEXT,
    ADD COLUMN IF NOT EXISTS ai_attempts INT NOT NULL DEFAULT 0;

ALTER TABLE interviews
    ADD COLUMN IF NOT EXISTS report_error TEXT,
    ADD COLUMN IF NOT EXISTS report_attempts INT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS report_updated_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS candidates_cv_ai_status_idx ON candidates (cv_ai_status);
CREATE INDEX IF NOT EXISTS mock_interviews_ai_status_idx ON mock_interviews (ai_question_status, ai_scoring_status);
CREATE INDEX IF NOT EXISTS interviews_report_status_idx ON interviews (report_status);
