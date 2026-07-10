-- interview_templates table is created in 000005. Add a company_id index for
-- tenant-scoped list queries backing the recruiter "Mẫu AI" page.
CREATE INDEX IF NOT EXISTS interview_templates_company_id_idx ON interview_templates (company_id);
