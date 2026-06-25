CREATE TABLE interview_reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    interview_id UUID NOT NULL UNIQUE REFERENCES interviews(id),
    summary TEXT NOT NULL,
    final_score NUMERIC(6,2),
    recommendation VARCHAR(50) NOT NULL,
    strengths JSONB,
    weaknesses JSONB,
    risks JSONB,
    evidence_json JSONB,
    ai_reasoning_summary TEXT,
    recruiter_decision VARCHAR(50),
    recruiter_comment TEXT,
    report_json JSONB NOT NULL,
    generated_by VARCHAR(50) NOT NULL,
    generated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);\n