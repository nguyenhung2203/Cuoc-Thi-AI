ALTER TABLE ai_prompt_templates ADD COLUMN company_id UUID REFERENCES companies(id);
ALTER TABLE ai_prompt_templates ADD COLUMN deleted_at TIMESTAMPTZ;

ALTER TABLE ai_request_logs ADD COLUMN company_id UUID REFERENCES companies(id);
