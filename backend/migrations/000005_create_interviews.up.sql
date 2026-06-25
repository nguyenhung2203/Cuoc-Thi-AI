CREATE TABLE interview_templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(100) NOT NULL,
    duration_minutes INT NOT NULL,
    description TEXT,
    config_json JSONB,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE question_bank (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID REFERENCES companies(id),
    job_id UUID REFERENCES jobs(id),
    created_by UUID REFERENCES users(id),
    question_text TEXT NOT NULL,
    question_type VARCHAR(100) NOT NULL,
    skill_tags JSONB,
    level VARCHAR(100),
    expected_signals JSONB,
    is_ai_generated BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX question_bank_company_job_idx ON question_bank (company_id, job_id);
CREATE INDEX question_bank_type_level_idx ON question_bank (question_type, level);

CREATE TABLE rubrics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id),
    job_id UUID REFERENCES jobs(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    total_weight NUMERIC NOT NULL DEFAULT 100,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE rubric_criteria (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    rubric_id UUID NOT NULL REFERENCES rubrics(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    weight NUMERIC NOT NULL,
    min_score INT NOT NULL DEFAULT 1,
    max_score INT NOT NULL DEFAULT 5,
    scoring_guide JSONB,
    order_index INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE interviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID REFERENCES companies(id),
    job_id UUID REFERENCES jobs(id),
    candidate_id UUID NOT NULL REFERENCES candidates(id),
    recruiter_id UUID REFERENCES users(id),
    template_id UUID REFERENCES interview_templates(id),
    rubric_id UUID REFERENCES rubrics(id),
    mode VARCHAR(50) NOT NULL CHECK(mode IN ('real', 'mock')),
    title VARCHAR(255) NOT NULL,
    scheduled_at TIMESTAMPTZ,
    started_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ,
    status VARCHAR(50) NOT NULL DEFAULT 'scheduled',
    room_id UUID,
    invite_token_hash TEXT,
    invite_expires_at TIMESTAMPTZ,
    consent_recording BOOLEAN NOT NULL DEFAULT false,
    consent_ai BOOLEAN NOT NULL DEFAULT false,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX interviews_company_status_idx ON interviews (company_id, status);
CREATE INDEX interviews_job_candidate_idx ON interviews (job_id, candidate_id);
CREATE INDEX interviews_scheduled_at_idx ON interviews (scheduled_at);