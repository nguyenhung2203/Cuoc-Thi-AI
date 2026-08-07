-- Persist portal CV AI parse on the file row so results survive when the user
-- has no candidates row yet (common for portal-only uploads).
ALTER TABLE files
    ADD COLUMN IF NOT EXISTS parsed_json JSONB,
    ADD COLUMN IF NOT EXISTS ai_summary TEXT,
    ADD COLUMN IF NOT EXISTS parse_status VARCHAR(50),
    ADD COLUMN IF NOT EXISTS parse_error TEXT,
    ADD COLUMN IF NOT EXISTS parsed_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS files_owner_cv_parse_idx
    ON files (owner_user_id, created_at DESC)
    WHERE file_type = 'cv';
