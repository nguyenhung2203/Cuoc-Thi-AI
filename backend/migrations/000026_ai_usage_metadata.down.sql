DROP INDEX IF EXISTS ai_request_logs_operation_idx;
DROP INDEX IF EXISTS ai_request_logs_model_idx;
DROP INDEX IF EXISTS ai_request_logs_created_at_idx;

ALTER TABLE ai_request_logs
    DROP COLUMN IF EXISTS retry_count,
    DROP COLUMN IF EXISTS total_tokens,
    DROP COLUMN IF EXISTS operation,
    DROP COLUMN IF EXISTS model,
    DROP COLUMN IF EXISTS provider;
