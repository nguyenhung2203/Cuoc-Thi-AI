DROP INDEX IF EXISTS ai_prompt_templates_sys_idx;
DROP INDEX IF EXISTS ai_prompt_templates_com_idx;

ALTER TABLE ai_request_logs ALTER COLUMN template_version TYPE VARCHAR(50);
ALTER TABLE ai_prompt_templates ALTER COLUMN version TYPE VARCHAR(50);

CREATE UNIQUE INDEX ai_prompt_templates_name_version_idx ON ai_prompt_templates (name, version);
