ALTER TABLE ai_request_logs DROP COLUMN IF EXISTS company_id;
ALTER TABLE ai_prompt_templates DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE ai_prompt_templates DROP COLUMN IF EXISTS company_id;
