ALTER TABLE ai_request_logs
    ADD COLUMN IF NOT EXISTS provider VARCHAR(50) NOT NULL DEFAULT 'gemini',
    ADD COLUMN IF NOT EXISTS model VARCHAR(100),
    ADD COLUMN IF NOT EXISTS operation VARCHAR(100),
    ADD COLUMN IF NOT EXISTS total_tokens INT,
    ADD COLUMN IF NOT EXISTS retry_count INT NOT NULL DEFAULT 0;

CREATE INDEX IF NOT EXISTS ai_request_logs_created_at_idx ON ai_request_logs (created_at);
CREATE INDEX IF NOT EXISTS ai_request_logs_model_idx ON ai_request_logs (model);
CREATE INDEX IF NOT EXISTS ai_request_logs_operation_idx ON ai_request_logs (operation);
