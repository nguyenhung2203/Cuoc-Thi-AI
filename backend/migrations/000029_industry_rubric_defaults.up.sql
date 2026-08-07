-- Generic, industry-neutral rubric defaults used when a job has no custom rubric.
CREATE TABLE IF NOT EXISTS industry_rubric_defaults (
    industry_key VARCHAR(100) PRIMARY KEY,
    display_name VARCHAR(255) NOT NULL,
    criteria JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
INSERT INTO industry_rubric_defaults (industry_key, display_name, criteria) VALUES
('general','General','[{"name":"Domain knowledge","weight":70},{"name":"Communication","weight":15},{"name":"Professionalism","weight":15}]'),
('technology','Technology','[{"name":"Technical/domain knowledge","weight":70},{"name":"Communication","weight":15},{"name":"Problem solving","weight":15}]'),
('sales','Sales','[{"name":"Product/customer understanding","weight":50},{"name":"Communication","weight":25},{"name":"Persuasion and ownership","weight":25}]'),
('operations','Operations','[{"name":"Process knowledge","weight":60},{"name":"Communication","weight":20},{"name":"Reliability and problem solving","weight":20}]'),
('finance','Finance','[{"name":"Financial/domain knowledge","weight":70},{"name":"Communication","weight":15},{"name":"Accuracy and integrity","weight":15}]'),
('healthcare','Healthcare','[{"name":"Domain knowledge","weight":60},{"name":"Communication","weight":20},{"name":"Professionalism and empathy","weight":20}]')
ON CONFLICT (industry_key) DO UPDATE SET criteria = EXCLUDED.criteria, display_name = EXCLUDED.display_name;
