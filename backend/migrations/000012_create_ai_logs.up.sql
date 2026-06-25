CREATE TABLE ai_prompt_templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    version VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    variables_schema JSONB,
    model VARCHAR(100) NOT NULL,
    params JSONB,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX ai_prompt_templates_name_version_idx ON ai_prompt_templates (name, version);

CREATE TABLE ai_request_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    template_id UUID REFERENCES ai_prompt_templates(id),
    template_version VARCHAR(50),
    interview_id UUID REFERENCES interviews(id),
    job_id UUID REFERENCES jobs(id),
    candidate_id UUID REFERENCES candidates(id),
    input_json JSONB,
    output_json JSONB,
    latency_ms INT,
    tokens_in INT,
    tokens_out INT,
    cost NUMERIC,
    status VARCHAR(50) NOT NULL,
    error TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);\n