-- files must be created before candidates because candidates.cv_file_id references it
CREATE TABLE files (
    id             UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id     UUID         REFERENCES companies (id),
    owner_user_id  UUID         REFERENCES users (id),
    original_name  VARCHAR(255) NOT NULL,
    storage_key    TEXT         NOT NULL,
    mime_type      VARCHAR(255) NOT NULL,
    size_bytes     BIGINT       NOT NULL,
    file_type      VARCHAR(100) NOT NULL
                                CHECK (file_type IN ('cv', 'audio_recording', 'avatar', 'attachment')),
    checksum       VARCHAR(255),
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX files_company_type_idx    ON files (company_id, file_type);
CREATE INDEX files_owner_user_id_idx   ON files (owner_user_id);

-- ---------------------------------------------------------------------------

CREATE TABLE candidates (
    id              UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id      UUID         NOT NULL REFERENCES companies (id),
    user_id         UUID         REFERENCES users (id),
    full_name       VARCHAR(255) NOT NULL,
    email           VARCHAR(255) NOT NULL,
    phone           VARCHAR(50),
    avatar_url      TEXT,
    cv_file_id      UUID         REFERENCES files (id),
    parsed_cv_json  JSONB,
    ai_cv_summary   TEXT,
    source          VARCHAR(100) CHECK (source IN ('linkedin', 'referral', 'import', 'manual')),
    status          VARCHAR(50)  NOT NULL DEFAULT 'new'
                                 CHECK (status IN (
                                     'new', 'screening', 'invited', 'interviewing',
                                     'completed', 'passed', 'rejected', 'talent_pool'
                                 )),
    tags            JSONB        NOT NULL DEFAULT '[]',
    created_by      UUID         REFERENCES users (id),
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX candidates_company_email_idx     ON candidates (company_id, email);
CREATE INDEX candidates_company_status_idx    ON candidates (company_id, status);
CREATE INDEX candidates_company_full_name_idx ON candidates (company_id, full_name);

-- ---------------------------------------------------------------------------

CREATE TABLE job_candidates (
    id              UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id      UUID         NOT NULL REFERENCES companies (id),
    job_id          UUID         NOT NULL REFERENCES jobs (id),
    candidate_id    UUID         NOT NULL REFERENCES candidates (id),
    pipeline_status VARCHAR(50)  NOT NULL DEFAULT 'new',
    fit_score       NUMERIC(5,2),
    ai_match_json   JSONB,
    applied_at      TIMESTAMPTZ,
    created_by      UUID         REFERENCES users (id),
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT job_candidates_job_candidate_unique UNIQUE (job_id, candidate_id)
);

CREATE INDEX job_candidates_job_id_idx               ON job_candidates (job_id);
CREATE INDEX job_candidates_candidate_id_idx         ON job_candidates (candidate_id);
CREATE INDEX job_candidates_company_pipeline_idx     ON job_candidates (company_id, pipeline_status);
