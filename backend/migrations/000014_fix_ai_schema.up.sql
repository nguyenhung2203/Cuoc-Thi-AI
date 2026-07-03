-- Drop the restrictive global unique index
DROP INDEX IF EXISTS ai_prompt_templates_name_version_idx;

-- Alter version from VARCHAR to INT
ALTER TABLE ai_prompt_templates ALTER COLUMN version TYPE INT USING version::integer;
ALTER TABLE ai_request_logs ALTER COLUMN template_version TYPE INT USING template_version::integer;

-- Create correct scoped unique constraints
-- 1. System templates (company_id is NULL) must be unique by name and version
CREATE UNIQUE INDEX ai_prompt_templates_sys_idx 
ON ai_prompt_templates (name, version) 
WHERE company_id IS NULL;

-- 2. Company templates (company_id is NOT NULL) must be unique by name, version, and company_id
CREATE UNIQUE INDEX ai_prompt_templates_com_idx 
ON ai_prompt_templates (name, version, company_id) 
WHERE company_id IS NOT NULL;
