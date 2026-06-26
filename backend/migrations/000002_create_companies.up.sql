CREATE TABLE companies (
    id         UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       VARCHAR(255) NOT NULL,
    slug       VARCHAR(255) NOT NULL,
    logo_url   TEXT,
    website    TEXT,
    industry   VARCHAR(255),
    size       VARCHAR(100),
    created_by UUID         NOT NULL REFERENCES users (id),
    settings   JSONB,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX companies_slug_unique ON companies (slug) WHERE deleted_at IS NULL;
CREATE INDEX companies_created_by_idx    ON companies (created_by);

-- ---------------------------------------------------------------------------

CREATE TABLE company_members (
    id         UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID        NOT NULL REFERENCES companies (id),
    user_id    UUID        NOT NULL REFERENCES users (id),
    role       VARCHAR(50) NOT NULL DEFAULT 'member'
                           CHECK (role IN ('owner', 'admin', 'member', 'viewer')),
    status     VARCHAR(50) NOT NULL DEFAULT 'invited'
                           CHECK (status IN ('invited', 'active', 'removed')),
    invited_by UUID        REFERENCES users (id),
    joined_at  TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT company_members_company_user_unique UNIQUE (company_id, user_id)
);

CREATE INDEX company_members_company_id_idx ON company_members (company_id);
CREATE INDEX company_members_user_id_idx    ON company_members (user_id);
