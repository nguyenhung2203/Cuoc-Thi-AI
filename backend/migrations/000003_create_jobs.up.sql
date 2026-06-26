CREATE TABLE jobs (
    id                UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id        UUID         NOT NULL REFERENCES companies (id),
    title             VARCHAR(255) NOT NULL,
    department        VARCHAR(255),
    level             VARCHAR(100) CHECK (level IN ('intern', 'junior', 'middle', 'senior', 'lead')),
    location          VARCHAR(255),
    employment_type   VARCHAR(100) CHECK (employment_type IN ('full_time', 'part_time', 'contract', 'intern')),
    salary_min        NUMERIC,
    salary_max        NUMERIC,
    currency          VARCHAR(20),
    description       TEXT         NOT NULL,
    requirements      TEXT,
    benefits          TEXT,
    status            VARCHAR(50)  NOT NULL DEFAULT 'draft'
                                   CHECK (status IN ('draft', 'open', 'paused', 'closed')),
    ai_summary        TEXT,
    ai_analysis_json  JSONB,
    created_by        UUID         NOT NULL REFERENCES users (id),
    created_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ
);

CREATE INDEX jobs_company_status_idx ON jobs (company_id, status);
CREATE INDEX jobs_company_title_idx  ON jobs (company_id, title);
