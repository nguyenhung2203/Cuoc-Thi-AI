DROP INDEX IF EXISTS files_owner_cv_parse_idx;

ALTER TABLE files
    DROP COLUMN IF EXISTS parsed_json,
    DROP COLUMN IF EXISTS ai_summary,
    DROP COLUMN IF EXISTS parse_status,
    DROP COLUMN IF EXISTS parse_error,
    DROP COLUMN IF EXISTS parsed_at;
